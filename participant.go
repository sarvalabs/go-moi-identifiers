package identifiers

import "encoding/binary"

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

// ParticipantFlag is a flag specifier for ParticipantID's Flags.
// Returns the minimum supported version and the bit location of the flag.
//
// All ParticipantID flags are located in the [1] index.
// Use the ParticipantID.Flag() method to check if a specific flag is set.
type ParticipantFlag func() (version, location uint8)

// ParticipantFlagMaster is a ParticipantFlag for the IsMaster flag on ParticipantID.
// It is supported from version 0 and located at the 0th flag bit.
//
// The IsMaster flag indicates if the identifier is for the master account of the participant.
func ParticipantFlagMaster() (uint8, uint8) { return 0, 0x0 }
