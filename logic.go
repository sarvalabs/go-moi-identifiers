package identifiers

import (
	"encoding/binary"
	"encoding/hex"
	"errors"
	"fmt"
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
// The given value is trimmed/padded to 32 bytes and validated into a LogicID.
func NewLogicIDFromBytes(data []byte) (LogicID, error) {
	return NewLogicID(trim32(data))
}

// NewLogicIDFromHex creates a new LogicID from the given hex string.
// The given value is hex-decoded (must contain 0x prefix),
// trimmed/padded to 32 bytes and validated into a LogicID.
func NewLogicIDFromHex(data string) (LogicID, error) {
	// Decode the given hex string into []byte
	decoded, err := decodeHexString(data)
	if err != nil {
		return Nil, err
	}

	// Create a new LogicID from the decoded value
	return NewLogicIDFromBytes(decoded)
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

// LogicFlag is a flag specifier for LogicID's Flags.
// Returns the minimum supported version and the bit location of the flag.
//
// All LogicID flags are located in the [1] index.
// Use the LogicID.Flag() method to check if a specific flag is set.
type LogicFlag func() (version, location uint8)

// todo: add logic flags for hasLogicState, hasActorState, isAssetLogic

// Flag returns if the given LogicFlag is set on the LogicID.
//
// If the specified flag is not supported by the LogicID's version,
// it will return FALSE, regardless of the actual flag value.
func (logic LogicID) Flag(flag LogicFlag) bool {
	ver, loc := flag()

	// Check if the flag is supported by the logic ID version
	// If not supported, return FALSE, regardless of the actual flag value
	if logic.Tag().Version() > ver {
		return false
	}

	// Check if flag bit is set
	return isBitSet(logic[1], loc)
}

// Validate returns an error if the LogicID is invalid.
// An error is returned if the AssetID has an invalid tag or contains unsupported flags.
func (logic LogicID) Validate() error {
	// Check basic validity of the identifier tag
	if err := logic.Tag().Validate(); err != nil {
		return fmt.Errorf("invalid tag: %w", err)
	}

	// Check if the tag is an asset tag
	if logic.Tag().Kind() != KindLogic {
		return errors.New("invalid tag: not a logic identifier")
	}

	// Perform checks based on the asset ID version
	switch logic.Tag().Version() {
	case 0:
		// Check if the flags for any position apart from 0, 1 & 2 are set
		if (logic[1] & byte(0b00011111)) != 0 {
			return errors.New("invalid flags: malformed flags for logic identifier v0")
		}

	default:
		return ErrUnsupportedVersion
	}

	return nil
}

// Bytes returns the LogicID as a []byte
func (logic LogicID) Bytes() []byte { return logic[:] }

// String returns the LogicID as a hex-encoded string.
// This is identical to LogicID.Hex() but is required for the fmt.Stringer interface
func (logic LogicID) String() string { return logic.Hex() }

// Hex returns the LogicID as a hex-encoded string with the 0x prefix
func (logic LogicID) Hex() string {
	return prefix0xString + hex.EncodeToString(logic[:])
}

// IsNil returns if the LogicID is nil, i.e., 0x000..000
func (logic LogicID) IsNil() bool {
	return logic == Nil
}

// AsIdentifier returns the LogicID as an Identifier.
func (logic LogicID) AsIdentifier() Identifier {
	return Identifier(logic)
}

// MarshalText implements the encoding.TextMarshaler interface for LogicID
func (logic *LogicID) MarshalText() ([]byte, error) {
	return marshal32(*logic)
}

// UnmarshalText implements the encoding.TextUnmarshaler interface for LogicID
func (logic *LogicID) UnmarshalText(data []byte) error {
	decoded, err := unmarshal32(data)
	if err != nil {
		return err
	}

	*logic = decoded
	return nil
}

// GenerateLogicIDv0 creates a new LogicID for v0 with the given parameters.
// Returns an error if unsupported flags are used.
//
// [tag:1][{hasLogicState}{hasActorState}{isAssetLogic}{reserved:5}][standard:2][account:24][variant:4]
func GenerateLogicIDv0(account [24]byte, variant uint32, flagOpts ...LogicFlag) (LogicID, error) {
	// Create the metadata buffer
	// [tag][flags][standard]
	metadata := make([]byte, 4)

	// Attach the tag for LogicID v0
	metadata[0] = byte(TagLogicV0)
	// Attach the flags to the metadata
	for _, opt := range flagOpts {
		// Get the minimum supported version and the location of the flag
		minver, loc := opt()
		// Check that flag is supported by version 0
		if 0 < minver {
			return Nil, ErrUnsupportedFlag
		}

		metadata[1] |= 1 << loc
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

// RandomLogicID generates a random LogicID.
func RandomLogicID() LogicID {
	// todo
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
