package identifiers

import (
	"encoding"
	"encoding/binary"
	"encoding/hex"
	"errors"
	"fmt"
	"math/rand/v2"
)

// ParticipantID is a unique identifier for a participant in the MOI Protocol.
// It is 32 bytes long and its first 4 bytes are structured as follows:
//   - Tag: The first byte contains the tag for the participant identifier.
//   - Flags: The second byte contains flags for the participant identifier.
//   - Metadata: As of v0, ParticipantID has no metadata.
//
// Like all identifiers, the ParticipantID also contains an AccountID and a Variant ID.
// Flags of a ParticipantID are specific to a version and are invalid if set in an unsupported version.
type ParticipantID [32]byte

// NewParticipantID creates a new ParticipantID from the 32-byte value.
// It returns an error if the given data is not a valid ParticipantID.
func NewParticipantID(data [32]byte) (ParticipantID, error) {
	// Convert the data into a ParticipantID
	participantID := ParticipantID(data)
	// Validate the ParticipantID
	if err := participantID.Validate(); err != nil {
		return Nil, err
	}

	return participantID, nil
}

// NewParticipantIDFromBytes creates a new ParticipantID from the given byte slice.
// The given value must have a length of 32 and validate into a ParticipantID.
func NewParticipantIDFromBytes(data []byte) (ParticipantID, error) {
	// Check length of the data
	if len(data) != 32 {
		return Nil, errors.New("invalid length: participant id must be 32 bytes")
	}

	return NewParticipantID([32]byte(data))
}

// NewParticipantIDFromHex creates a new ParticipantID from the given hex string.
// The given value must decode as hexadecimal string (0x prefix is optional),
// with a length of 64 characters (32 bytes) and validate into a ParticipantID.
func NewParticipantIDFromHex(data string) (ParticipantID, error) {
	// Decode the given hex string into []byte
	decoded, err := decodeHexString(data)
	if err != nil {
		return Nil, err
	}

	// Create a new ParticipantID from the decoded value
	// Length check is performed in NewParticipantIDFromBytes
	return NewParticipantIDFromBytes(decoded)
}

// MustParticipantID is an enforced version of NewParticipantID.
// Panics if an error occurs. Use with caution.
func MustParticipantID(data [32]byte) ParticipantID { return must(NewParticipantID(data)) }

// MustParticipantIDFromBytes is an enforced version of NewParticipantIDFromBytes.
// Panics if an error occurs. Use with caution.
func MustParticipantIDFromBytes(data []byte) ParticipantID {
	return must(NewParticipantIDFromBytes(data))
}

// MustParticipantIDFromHex is an enforced version of NewParticipantIDFromHex.
// Panics if an error occurs. Use with caution.
func MustParticipantIDFromHex(data string) ParticipantID {
	return must(NewParticipantIDFromHex(data))
}

// Bytes returns the ParticipantID as a []byte.
func (participant ParticipantID) Bytes() []byte { return participant[:] }

// String returns the ParticipantID as a hex-encoded string.
// This is identical to ParticipantID.Hex() but is required for the fmt.Stringer interface.
func (participant ParticipantID) String() string { return participant.Hex() }

// Hex returns the ParticipantID as a hex-encoded string with the 0x prefix.
func (participant ParticipantID) Hex() string {
	return prefix0xString + hex.EncodeToString(participant[:])
}

// AsIdentifier returns the ParticipantID as an Identifier.
func (participant ParticipantID) AsIdentifier() Identifier {
	return Identifier(participant)
}

// Tag returns the IdentifierTag for the ParticipantID.
func (participant ParticipantID) Tag() IdentifierTag {
	return IdentifierTag(participant[0])
}

// AccountID returns the 24-byte account ID from the ParticipantID.
func (participant ParticipantID) AccountID() [24]byte {
	return trimAccount(participant)
}

// Variant returns the 32-bit variant ID from the ParticipantID.
func (participant ParticipantID) Variant() uint32 {
	variant := trimVariant(participant)
	return binary.BigEndian.Uint32(variant[:])
}

// IsVariant returns if the ParticipantID has a non-zero variant ID.
func (participant ParticipantID) IsVariant() bool {
	variant := trimVariant(participant)
	return !(variant[0] == 0 && variant[1] == 0 && variant[2] == 0 && variant[3] == 0)
}

// Flag returns if the given Flag is set on the ParticipantID.
//
// If the specified flag is not supported by the ParticipantID,
// it will return False, regardless of the actual flag value.
func (participant ParticipantID) Flag(flag Flag) bool {
	// Check if the flag is supported by ParticipantID.
	// If not supported, return FALSE, regardless of the actual flag value
	if !flag.Supports(participant.Tag()) {
		return false
	}

	return getFlag(participant[1], flag.index)
}

// Validate returns an error if the ParticipantID is invalid.
// An error is returned if the ParticipantID has an invalid tag or contains unsupported flags.
func (participant ParticipantID) Validate() error {
	// Check basic validity of the identifier tag
	if err := participant.Tag().Validate(); err != nil {
		return fmt.Errorf("invalid tag: %w", err)
	}

	// Check if the tag is a participant tag
	if participant.Tag().Kind() != KindParticipant {
		return errors.New("invalid tag: not a participant id")
	}

	// Check that there are no unsupported flags set
	if (participant[1] & flagMasks[participant.Tag()]) != 0 {
		return errors.New("invalid flags: unsupported flags for participant id")
	}

	return nil
}

var (
	// Ensure ParticipantID implements text marshaling interfaces
	_ encoding.TextMarshaler   = (*ParticipantID)(nil)
	_ encoding.TextUnmarshaler = (*ParticipantID)(nil)
)

// MarshalText implements the encoding.TextMarshaler interface for ParticipantID
func (participant ParticipantID) MarshalText() ([]byte, error) {
	return marshal32(participant)
}

// UnmarshalText implements the encoding.TextUnmarshaler interface for ParticipantID
func (participant *ParticipantID) UnmarshalText(data []byte) error {
	decoded, err := unmarshal32(data)
	if err != nil {
		return err
	}

	if err = ParticipantID(decoded).Validate(); err != nil {
		return err
	}

	*participant = decoded
	return nil
}

// GenerateParticipantIDv0 creates a new ParticipantID for v0 with the given parameters.
// Returns an error if unsupported flags are used.
//
// [tag:1][{systemic}{reserved:7}][standard:2][account:24][variant:4]
func GenerateParticipantIDv0(account [24]byte, variant uint32, flags ...Flag) (ParticipantID, error) {
	// Create the metadata buffer
	// [tag][flags][standard]
	metadata := make([]byte, 4)
	// Attach the tag for ParticipantID v0
	metadata[0] = byte(TagParticipantV0)

	// Attach the flags to the metadata
	for _, flag := range flags {
		// Check if the given flag is supported by ParticipantID v0
		if !flag.Supports(TagParticipantV0) {
			return Nil, ErrUnsupportedFlag
		}

		// Set the flag in the metadata
		metadata[1] = setFlag(metadata[1], flag.index, true)
	}

	// Order the participant ID buffer
	// [metadata][account][variant]
	buffer := make([]byte, 0, 32)
	buffer = append(buffer, metadata...)
	buffer = append(buffer, account[:]...)
	// Append 4 bytes for the variant and encode the value into it
	buffer = append(buffer, make([]byte, 4)...)
	binary.BigEndian.PutUint32(buffer[28:], variant)

	return ParticipantID(buffer), nil
}

// RandomParticipantIDv0 creates a random v0 ParticipantID
// with a random account ID, variant ID and flags.
//   - There is a 0% chance that the Systemic flag will be set.
func RandomParticipantIDv0() ParticipantID {
	// Safe to ignore error as the flags are supported
	participant, _ := GenerateParticipantIDv0(RandomAccountID(), rand.Uint32())
	return participant
}
