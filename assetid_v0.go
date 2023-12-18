package identifiers

import (
	"encoding/binary"
	"encoding/hex"

	"github.com/pkg/errors"
)

// AssetIDV0Length is the length of the v0 specification of the AssetID Standard
const AssetIDV0Length = 36

// AssetIdentifierV0 is an implementation of v0 specification
// of the AssetID Standard and implements the AssetIdentifier
type AssetIdentifierV0 [AssetIDV0Length]byte

// NewAssetIDv0 generates a new AssetID with the v0 specification. The AssetID v0 Form is defined as follows:
// [version(4bits)|logical(1bit)|stateful(1bit)|reserved(2bits)][dimension(8bits)][standard(16bits)][address(256bits)]
func NewAssetIDv0(logical, stateful bool, dimension uint8, standard uint16, addr Address) AssetID {
	// The 4 MSB bits of the head are set the
	// version of the Asset ID Form (v0)
	var head uint8 = 0x00 << 4

	// If logical flag is on, the 5th MSB is set
	if logical {
		head |= 0x8
	}

	// If stateful flag is on, the 6th MSB is set
	if stateful {
		head |= 0x4
	}

	// Encode the 16-bit standard into its BigEndian bytes
	standardBuf := make([]byte, 2)
	binary.BigEndian.PutUint16(standardBuf, uint16(standard))

	// Order the asset ID buffer [head][dimension][standard][address]
	buf := make([]byte, 0, 36)
	buf = append(buf, head)
	buf = append(buf, dimension)
	buf = append(buf, standardBuf...)
	buf = append(buf, addr[:]...)

	return AssetID(hex.EncodeToString(buf))
}

// decodeAssetIDv0 can be used to decode some data into a AssetIdentifierV0.
func decodeAssetIDv0(data []byte) (AssetIdentifierV0, error) {
	// Check if data is the correct length for v0
	if len(data) != AssetIDV0Length {
		return AssetIdentifierV0{}, errors.New("invalid asset ID: insufficient length for v0")
	}

	// Create an AssetIdentifierV0 and copy the data into it
	identifier := AssetIdentifierV0{}
	copy(identifier[:], data)

	return identifier, nil
}

// AssetID returns the AssetIdentifierV0 as an AssetID
func (asset AssetIdentifierV0) AssetID() AssetID {
	return AssetID(hex.EncodeToString(asset[:]))
}

// Version returns the version of the AssetIdentifierV0.
func (asset AssetIdentifierV0) Version() int { return 0 }

// IsLogical returns whether the logical flag is set for the AssetIdentifierV0.
func (asset AssetIdentifierV0) IsLogical() bool {
	// Determine the 5th LSB of the first byte (v0)
	bit := (asset[0] >> 3) & 0x1
	// Return true if bit is set
	return bit != 0
}

// IsStateful returns whether the stateful flag is set for the AssetIdentifierV0.
func (asset AssetIdentifierV0) IsStateful() bool {
	// Determine the 6th LSB of the first byte (v0)
	bit := (asset[0] >> 2) & 0x1
	// Return true if bit is set
	return bit != 0
}

// Dimension returns the dimension of the AssetIdentifierV0.
func (asset AssetIdentifierV0) Dimension() uint8 {
	// Dimension data is in the second byte of the AssetID (v0)
	return asset[1]
}

// Standard returns the standard of the AssetIdentifierV0.
func (asset AssetIdentifierV0) Standard() uint64 {
	// Decode the edition data from the third and fourth byte of
	// the AssetID (v0). We decode it as 16-bit number and convert
	standard := binary.BigEndian.Uint16(asset[2:4])

	return uint64(standard)
}

// Address returns the Asset Address of the AssetIdentifier.
func (asset AssetIdentifierV0) Address() Address {
	// Address data is everything after the fourth byte (v0)
	// We know it will be 32 bytes, because of the validity check
	address := asset[4:]
	// Convert address data into an Address and return
	return NewAddressFromBytes(address)
}
