package identifiers

import (
	"encoding/binary"
	"encoding/hex"

	"github.com/pkg/errors"
)

// LogicIDV0Length is the length of the v0 specification of the LogicID Standard
const LogicIDV0Length = 35

// LogicIdentifierV0 is an implementation of v0 specification
// of the LogicID Standard and implements the LogicIdentifier
type LogicIdentifierV0 [LogicIDV0Length]byte

// NewLogicIDv0 generates a new LogicID with the v0 specification. The LogicID v0 Form is defined as follows:
// [version(4bits)|persistent(1bit)|ephemeral(1bit)|interactable(1bit)|asset(1bit)][edition(16bits)][address(256bits)]
func NewLogicIDv0(persistent, ephemeral, interactable, assetlogic bool, edition uint16, addr Address) LogicID {
	// The 4 MSB bits of the head are set the
	// version of the Logic ID Form (v0)
	var head uint8 = 0x00 << 4

	// If persistent stateful flag is on, the 5th MSB is set
	if persistent {
		head |= 0x8
	}

	// If ephemeral stateful flag is on, the 6th MSB is set
	if ephemeral {
		head |= 0x4
	}

	// If interactable flag is on, the 7th MSB is set
	if interactable {
		head |= 0x2
	}

	// If asset logic flag is on, the 8th MSB is set
	if assetlogic {
		head |= 0x1
	}

	// Encode the 16-bit edition into its BigEndian bytes
	editionBuf := make([]byte, 2)
	binary.BigEndian.PutUint16(editionBuf, edition)

	// Order the logic ID buffer [head][edition][address]
	buf := make([]byte, 0, 35)
	buf = append(buf, head)
	buf = append(buf, editionBuf...)
	buf = append(buf, addr[:]...)

	return LogicID(hex.EncodeToString(buf))
}

// decodeLogicIDv0 can be used to decode some data into a LogicIdentifierV0.
func decodeLogicIDv0(data []byte) (LogicIdentifierV0, error) {
	// Check if data is the correct length for v0
	if len(data) != LogicIDV0Length {
		return LogicIdentifierV0{}, errors.New("invalid logic ID: insufficient length for v0")
	}

	// Create an LogicIdentifierV0 and copy the data into it
	identifier := LogicIdentifierV0{}
	copy(identifier[:], data)

	return identifier, nil
}

// LogicID returns the LogicIdentifierV0 as a LogicID
func (logic LogicIdentifierV0) LogicID() LogicID {
	return LogicID(hex.EncodeToString(logic[:]))
}

// Version returns the version of the LogicIdentifierV0.
func (logic LogicIdentifierV0) Version() int { return 0 }

// HasPersistentState returns whether the persistent state flag is set for the LogicIdentifierV0.
func (logic LogicIdentifierV0) HasPersistentState() bool {
	// Determine the 5th LSB of the first byte (v0)
	bit := (logic[0] >> 3) & 0x1
	// Return true if bit is set
	return bit != 0
}

// HasEphemeralState returns whether the ephemeral state flag is set for the LogicIdentifierV0.
func (logic LogicIdentifierV0) HasEphemeralState() bool {
	// Determine the 6th LSB of the first byte (v0)
	bit := (logic[0] >> 2) & 0x1
	// Return true if bit is set
	return bit != 0
}

// HasInteractableSites returns whether the interactable flag is set for the LogicIdentifierV0.
func (logic LogicIdentifierV0) HasInteractableSites() bool {
	// Determine the 7th LSB of the first byte (v0)
	bit := (logic[0] >> 1) & 0x1
	// Return true if bit is set
	return bit != 0
}

// AssetLogic returns whether the asset logic flag is set for the LogicIdentifierV0.
func (logic LogicIdentifierV0) AssetLogic() bool {
	// Determine the 8th LSB of the first byte (v0)
	bit := logic[0] & 0x1
	// Return true if bit is set
	return bit != 0
}

// Edition returns the edition number of the LogicIdentifierV0.
func (logic LogicIdentifierV0) Edition() uint64 {
	// Decode the edition data from the second and third byte of
	// the LogicID (v0). We decode it as 16-bit number and convert
	edition := binary.BigEndian.Uint16(logic[1:3])

	return uint64(edition)
}

// Address returns the Logic Address of the LogicIdentifierV0.
func (logic LogicIdentifierV0) Address() Address {
	// Address data is everything after the third byte (v0)
	// We know it will be 32 bytes, because of the validity check
	address := logic[3:]
	// Convert address data into an Address and return
	return NewAddressFromBytes(address)
}
