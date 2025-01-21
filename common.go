package identifiers

import (
	"bytes"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"strings"

	"github.com/sarvalabs/go-polo"
)

// Nil is a nil [32]byte value.
// Can be used to represent any nil identifier.
var Nil [32]byte

// RandomFingerprint generates a random 24-byte fingerprint
func RandomFingerprint() (fingerprint [24]byte) {
	_, _ = rand.Read(fingerprint[:])
	return fingerprint
}

var (
	prefix0xString = "0x"
	prefix0xBytes  = []byte(prefix0xString)
)

var (
	ErrMissingHexPrefix = errors.New("missing '0x' prefix")
	ErrInvalidLength    = errors.New("invalid length")

	ErrUnsupportedFlag    = errors.New("unsupported flag")
	ErrUnsupportedVersion = errors.New("unsupported tag version")
	ErrUnsupportedKind    = errors.New("unsupported tag kind")
)

// trim0xPrefixString trims the 0x prefix from the given string (if it exists).
func trim0xPrefixString(value string) string {
	return strings.TrimPrefix(value, prefix0xString)
}

// trim0xPrefixBytes trims the 0x prefix from the given byte slice (if it exists).
func trim0xPrefixBytes(value []byte) []byte {
	return bytes.TrimPrefix(value, prefix0xBytes)
}

// has0xPrefixBytes checks if the given byte slice has a 0x prefix.
func has0xPrefixBytes(value []byte) bool {
	return bytes.HasPrefix(value, prefix0xBytes)
}

// decodeHexString decodes the given hex string into a byte slice.
// It trims the 0x prefix (if found) from the string before decoding.
func decodeHexString(str string) ([]byte, error) {
	// Trim the 0x prefix from the string (if it exists)
	str = trim0xPrefixString(str)

	decoded, err := hex.DecodeString(str)
	if err != nil {
		return nil, err
	}

	return decoded, nil
}

// trimFingerprint returns the 24 bytes in the middle of the given 32-byte array.
func trimFingerprint(bytes [32]byte) [24]byte {
	return [24]byte(bytes[4:28])
}

// trimVariant returns the 4 least-significant bytes of the given 32-byte array.
func trimVariant(bytes [32]byte) [4]byte {
	return [4]byte(bytes[28:])
}

// marshal32 is a generic marshal function for 32-byte identifiers.
// To be used in conjunction with MarshalText
func marshal32(data [32]byte) ([]byte, error) {
	buffer := make([]byte, 32*2+2)

	// Copy the 0x prefix into the buffer
	copy(buffer[:2], prefix0xString)
	// Hex-encode the copied value into the buffer
	hex.Encode(buffer[2:], data[:])

	return buffer, nil
}

// unmarshal32 is generic unmarshal function for 32-byte identifiers.
// To be used in conjunction with UnmarshalText
func unmarshal32(data []byte) ([32]byte, error) {
	// Assert that the 0x prefix exists
	if !has0xPrefixBytes(data) {
		return Nil, ErrMissingHexPrefix
	}

	// Trim the 0x prefix
	data = trim0xPrefixBytes(data)

	// Check that the data has enough length for the identifier data
	if len(data) != 32*2 {
		return Nil, ErrInvalidLength
	}

	// Decode the hex-encoded data
	decoded, err := decodeHexString(string(data))
	if err != nil {
		return Nil, err
	}

	return [32]byte(decoded), nil
}

func polorize32(data [32]byte) (*polo.Polorizer, error) {
	// Create a new polorizer
	polorizer := polo.NewPolorizer()
	// Encode the given 32-byte value
	polorizer.PolorizeBytes(data[:])

	return polorizer, nil
}

func depolorize32(buffer *polo.Depolorizer) ([32]byte, error) {
	// Check if buffer contains a null tag
	if buffer.IsNull() {
		return Nil, nil
	}

	// Attempt to decode a 32-byte value
	decoded, err := buffer.DepolorizeBytes32()
	if err != nil {
		return Nil, err
	}

	return decoded, nil
}

// must is correctness enforcer for error handling.
// For use in functions that should never return an error.
// Panics if an error is encountered.
func must[T any](t T, err error) T {
	if err != nil {
		panic(err)
	}

	return t
}
