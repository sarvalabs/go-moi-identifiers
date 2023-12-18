package identifiers

import (
	"crypto/rand"
	"encoding/hex"
)

// AddressLength represents the number of bytes (32) for an Address
const AddressLength = 32

// Address represents a unique 32-byte (256-bit) identifier
type Address [AddressLength]byte

// NilAddress represents a nil Address value
// which can also be represented as 0x0000...
var NilAddress Address

// NewRandomAddress generates a random Address
func NewRandomAddress() Address {
	// Generate a random [32]byte value
	address := make([]byte, 32)
	_, _ = rand.Read(address)

	// Create a new Address from the random value
	return NewAddressFromBytes(address)
}

// NewAddressFromBytes creates a new Address from the given []byte data.
// The length of data is expected be 32, but if it less than or greater
// than that, we trim (from the right) or pad (to the right) appropriately.
func NewAddressFromBytes(data []byte) (addr Address) {
	// Trim the data from if it exceeds 32 bytes
	if len(data) > AddressLength {
		data = data[len(data)-AddressLength:]
	}

	// Copy the given bytes into the address
	// This automatically pads the Address if the given data is less than 32 bytes
	copy(addr[AddressLength-len(data):], data)

	return addr
}

// NewAddressFromHex creates a new Address from the given hex-encoded string.
// Ignores the 0x prefix in the string if it exists and returns errors if the
// data is not correctly hex-encoded (has invalid or uneven number of characters)
func NewAddressFromHex(data string) (Address, error) {
	// Decode the given hex string into []byte
	decoded, err := decodeHexString(data)
	if err != nil {
		return NilAddress, err
	}

	// Create a new Address from decoded value
	return NewAddressFromBytes(decoded), nil
}

// IsNil returns if the Address is nil i.e, 0x0000...
func (addr Address) IsNil() bool {
	return addr == NilAddress
}

// Bytes returns the Address as a []byte
func (addr Address) Bytes() []byte {
	return addr[:]
}

// String returns the Address as a hex-encoded string with the 0x prefix.
// Implements the fmt.Stringer interface for Address.
func (addr Address) String() string {
	return addr.Hex()
}

// Hex return the hex-encoded representation of the Address with the 0x prefix.
func (addr Address) Hex() string {
	return "0x" + hex.EncodeToString(addr[:])
}

// MarshalText implements the encoding.TextMarshaler interface for Address
func (addr Address) MarshalText() ([]byte, error) {
	result := make([]byte, len(addr)*2+2)

	// Copy the 0x into the buffer
	copy(result[:2], "0x")
	// Hex-encode the copy the address value into the buffer
	hex.Encode(result[2:], addr.Bytes())

	return result, nil
}

// UnmarshalText implements the encoding.TextUnmarshaler interface for Address
func (addr *Address) UnmarshalText(text []byte) error {
	// Assert that the 0x prefix exists
	if !has0xPrefixBytes(text) {
		return ErrMissing0xPrefix
	}

	// Trim the 0x prefix
	text = trim0xPrefixBytes(text)
	// Check that text has enough length for the address data
	if len(text) != AddressLength*2 {
		return ErrInvalidLength
	}

	// Decode the hex-encoded text into the address
	_, err := hex.Decode(addr[:], text)

	return err
}
