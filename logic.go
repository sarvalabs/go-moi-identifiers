package identifiers

import (
	"encoding/binary"
	"encoding/hex"
	"errors"
	"fmt"
	"math/rand/v2"
)

// LogicID is a unique identifier for a logic in the MOI Protocol.
// It is 32 bytes long and its metadata is structured as follows:
//   - Tag: The first byte contains the tag for the logic identifier.
//   - Flags: The second byte contains flags for the logic identifier.
//
// Like all identifiers, the LogicID also contains an AccountID and a Variant ID.
// Flags of a LogicID are specific to a version and are invalid if set in an unsupported version.
type LogicID [32]byte

// NewLogicID creates a new LogicID from the 32-byte value.
// It returns an error if the given data is not a valid LogicID.
func NewLogicID(data [32]byte) (LogicID, error) {
	// Convert the data into a LogicID
	logicID := LogicID(data)
	// Validate the LogicID
	if err := logicID.Validate(); err != nil {
		return Nil, err
	}

	return logicID, nil
}

// NewLogicIDFromBytes creates a new LogicID from the given byte slice.
// The given value must have a length of 32 and validate into an LogicID.
func NewLogicIDFromBytes(data []byte) (LogicID, error) {
	// Check length of the data
	if len(data) != 32 {
		return Nil, errors.New("invalid length: logic id must be 32 bytes")
	}

	return NewLogicID([32]byte(data))
}

// NewLogicIDFromHex creates a new LogicID from the given hex string.
// The given value must decode as hexadecimal string (0x prefix is optional),
// with a length of 64 characters (32 bytes) and validate into an LogicID.
func NewLogicIDFromHex(data string) (LogicID, error) {
	// Decode the given hex string into []byte
	decoded, err := decodeHexString(data)
	if err != nil {
		return Nil, err
	}

	// Create a new LogicID from the decoded value
	// Length check is performed in NewLogicIDFromBytes
	return NewLogicIDFromBytes(decoded)
}

// MustLogicID is an enforced version of NewLogicID.
// Panics if an error occurs. Use with caution.
func MustLogicID(data [32]byte) LogicID { return must(NewLogicID(data)) }

// MustLogicIDFromBytes is an enforced version of NewLogicIDFromBytes.
// Panics if an error occurs. Use with caution.
func MustLogicIDFromBytes(data []byte) LogicID { return must(NewLogicIDFromBytes(data)) }

// MustLogicIDFromHex is an enforced version of NewLogicIDFromHex.
// Panics if an error occurs. Use with caution.
func MustLogicIDFromHex(data string) LogicID { return must(NewLogicIDFromHex(data)) }

// Bytes returns the LogicID as a []byte
func (logic LogicID) Bytes() []byte { return logic[:] }

// String returns the LogicID as a hex-encoded string.
// This is identical to LogicID.Hex() but is required for the fmt.Stringer interface
func (logic LogicID) String() string { return logic.Hex() }

// Hex returns the LogicID as a hex-encoded string with the 0x prefix
func (logic LogicID) Hex() string {
	return prefix0xString + hex.EncodeToString(logic[:])
}

// AsIdentifier returns the LogicID as an Identifier.
func (logic LogicID) AsIdentifier() Identifier {
	return Identifier(logic)
}

// Tag returns the IdentifierTag for the LogicID.
func (logic LogicID) Tag() IdentifierTag {
	return IdentifierTag(logic[0])
}

// AccountID returns the 24-byte account ID from the LogicID.
func (logic LogicID) AccountID() [24]byte {
	return trimMid24(logic)
}

// Variant returns the 32-bit variant ID from the LogicID.
func (logic LogicID) Variant() uint32 {
	low4 := trimLow4(logic)
	return binary.BigEndian.Uint32(low4[:])
}

// IsVariant returns if the LogicID has a non-zero variant ID
func (logic LogicID) IsVariant() bool {
	low4 := trimLow4(logic)
	return !(low4[0] == 0 && low4[1] == 0 && low4[2] == 0 && low4[3] == 0)
}

// Flag returns if the given Flag is set on the LogicID.
//
// If the specified flag is not supported by the LogicID,
// it will return False, regardless of the actual flag value.
func (logic LogicID) Flag(flag Flag) bool {
	// Check if the flag is supported by LogicID.
	// If not supported, return FALSE, regardless of the actual flag value
	if !flag.Supports(logic.Tag()) {
		return false
	}

	return getFlag(logic[1], flag.index)
}

// Validate returns an error if the LogicID is invalid.
// An error is returned if the LogicID has an invalid tag or contains unsupported flags.
func (logic LogicID) Validate() error {
	// Check basic validity of the identifier tag
	if err := logic.Tag().Validate(); err != nil {
		return fmt.Errorf("invalid tag: %w", err)
	}

	// Check if the tag is a logic tag
	if logic.Tag().Kind() != KindLogic {
		return errors.New("invalid tag: not a logic id")
	}

	// Check that there are no unsupported flags set
	if (logic[1] & flagMasks[logic.Tag()]) != 0 {
		return errors.New("invalid flags: unsupported flags for logic id")
	}

	return nil
}

// MarshalText implements the encoding.TextMarshaler interface for LogicID
func (logic LogicID) MarshalText() ([]byte, error) {
	return marshal32(logic)
}

// UnmarshalText implements the encoding.TextUnmarshaler interface for LogicID
func (logic *LogicID) UnmarshalText(data []byte) error {
	decoded, err := unmarshal32(data)
	if err != nil {
		return err
	}

	if err = LogicID(decoded).Validate(); err != nil {
		return err
	}

	*logic = decoded
	return nil
}

// GenerateLogicIDv0 creates a new LogicID for v0 with the given parameters.
// Returns an error if unsupported flags are used.
//
// [tag:1][{systemic}{reserved:4}{auxiliary}{extrinsic}{intrinsic}][standard:2][account:24][variant:4]
func GenerateLogicIDv0(account [24]byte, variant uint32, flags ...Flag) (LogicID, error) {
	// Create the metadata buffer
	// [tag][flags][standard]
	metadata := make([]byte, 4)
	// Attach the tag for LogicID v0
	metadata[0] = byte(TagLogicV0)

	// Attach the flags to the metadata
	for _, flag := range flags {
		// Check if the given flag is supported by LogicID v0
		if !flag.Supports(TagLogicV0) {
			return Nil, ErrUnsupportedFlag
		}

		// Set the flag in the metadata
		metadata[1] = setFlag(metadata[1], flag.index, true)
	}

	// Order the logic ID buffer
	// [metadata][account][variant]
	buffer := make([]byte, 0, 32)
	buffer = append(buffer, metadata...)
	buffer = append(buffer, account[:]...)
	// Append 4 bytes for the variant and encode the value into it
	buffer = append(buffer, make([]byte, 4)...)
	binary.BigEndian.PutUint32(buffer[28:], variant)

	return LogicID(buffer), nil
}

// RandomLogicIDv0 creates a random v0 LogicID
// with a random account ID, variant ID and flags.
//   - There is a 50% chance that the LogicIntrinsic flag will be set.
//   - There is a 50% chance that the LogicExtrinsic flag will be set.
//   - There is a 50% chance that the LogicAuxiliary flag will be set.
//   - There is a 0% chance that the Systemic flag will be set.
func RandomLogicIDv0() LogicID {
	flags := make([]Flag, 0, 3)

	if rand.Int64() > 0 {
		flags = append(flags, LogicIntrinsic)
	}

	if rand.Int64() > 0 {
		flags = append(flags, LogicExtrinsic)
	}

	if rand.Int64() > 0 {
		flags = append(flags, LogicAuxiliary)
	}

	// Safe to ignore error as the flags are supported
	logic, _ := GenerateLogicIDv0(RandomAccountID(), rand.Uint32(), flags...)

	return logic
}
