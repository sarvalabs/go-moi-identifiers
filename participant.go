package identifiers

import (
	"encoding/binary"
	"encoding/hex"
	"errors"
	"fmt"
	"math/rand/v2"
)

// ParticipantID is a unique identifier for a participant in the MOI Protocol.
// It is 32 bytes long and its metadata is structured as follows:
//   - Tag: The first byte contains the tag for the participant identifier.
//   - Flags: The second byte contains flags for the participant identifier.
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
// The given value is trimmed/padded to 32 bytes and validated into a ParticipantID.
func NewParticipantIDFromBytes(data []byte) (ParticipantID, error) {
	return NewParticipantID(trim32(data))
}

// NewParticipantIDFromHex creates a new ParticipantID from the given hex string.
// The given value is hex-decoded (must contain 0x prefix),
// trimmed/padded to 32 bytes and validated into a ParticipantID.
func NewParticipantIDFromHex(data string) (ParticipantID, error) {
	// Decode the given hex string into []byte
	decoded, err := decodeHexString(data)
	if err != nil {
		return Nil, err
	}

	// Create a new ParticipantID from the decoded value
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

// IsNil returns if the ParticipantID is nil, i.e., 0x000..000.
func (participant ParticipantID) IsNil() bool {
	return participant == Nil
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
	return trimMid24(participant)
}

// Variant returns the 32-bit variant ID from the ParticipantID.
func (participant ParticipantID) Variant() uint32 {
	low4 := trimLow4(participant)
	return binary.BigEndian.Uint32(low4[:])
}

// IsVariant returns if the ParticipantID has a non-zero variant ID.
func (participant ParticipantID) IsVariant() bool {
	return participant.Variant() != 0
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
	if participant.Tag().Kind() != KindLogic {
		return errors.New("invalid tag: not a participant identifier")
	}

	// Check that there are no unsupported flags set
	if (participant[1] & participant.Tag().FlagMask()) != 0 {
		return errors.New("invalid flags: malformed flags for participant identifier")
	}

	return nil
}

// MarshalText implements the encoding.TextMarshaler interface for ParticipantID
func (participant *ParticipantID) MarshalText() ([]byte, error) {
	return marshal32(*participant)
}

// UnmarshalText implements the encoding.TextUnmarshaler interface for ParticipantID
func (participant *ParticipantID) UnmarshalText(data []byte) error {
	decoded, err := unmarshal32(data)
	if err != nil {
		return err
	}

	*participant = decoded
	return nil
}

// GenerateParticipantIDv0 creates a new ParticipantID for v0 with the given parameters.
// Returns an error if unsupported flags are used.
//
// [tag:1][{reserved:7}{systemic}][standard:2][account:24][variant:4]
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

// RandomParticipantID creates a random v0 ParticipantID
// with a random account ID, variant ID and flags.
//   - There is a 0% chance that the Systemic flag will be set.
func RandomParticipantID() ParticipantID {
	// Safe to ignore error as the flags are supported
	participant, _ := GenerateParticipantIDv0(RandomAccountID(), rand.Uint32())
	return participant
}
