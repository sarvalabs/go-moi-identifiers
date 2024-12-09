package identifiers

import (
	"bytes"
	"encoding/hex"
	"errors"
	"strings"
)

// Nil is a nil [32]byte value.
// Can be used to represent any nil identifier.
var Nil [32]byte

var (
	prefix0xString = "0x"
	prefix0xBytes  = []byte(prefix0xString)
)

var (
	ErrMissingHexPrefix = errors.New("missing '0x' prefix")
	ErrInvalidLength    = errors.New("invalid length")
	ErrUnsupportedFlag  = errors.New("unsupported flag")
)

func trim0xPrefixString(value string) string {
	return strings.TrimPrefix(value, prefix0xString)
}

func trim0xPrefixBytes(value []byte) []byte {
	return bytes.TrimPrefix(value, prefix0xBytes)
}

func has0xPrefixBytes(value []byte) bool {
	return bytes.HasPrefix(value, prefix0xBytes)
}

func decodeHexString(str string) ([]byte, error) {
	// Trim the 0x prefix from the string (if it exists)
	str = trim0xPrefixString(str)

	decoded, err := hex.DecodeString(str)
	if err != nil {
		return nil, err
	}

	return decoded, nil
}

func trim32(data []byte) (trim [32]byte) {
	// Trim the data if it is longer than 32 bytes
	if len(data) > 32 {
		data = data[len(data)-32:]
	}

	// Copy the data into the trimmed array
	// This automatically pads the array if the given data is less than 32 bytes
	copy(trim[32-len(data):], data)

	return
}

func trimHigh4(bytes [32]byte) [4]byte {
	return [4]byte(bytes[:4])
}

func trimMid24(bytes [32]byte) [24]byte {
	return [24]byte(bytes[4:28])
}

func trimLow4(bytes [32]byte) [4]byte {
	return [4]byte(bytes[28:])
}

func isBitSet(value byte, loc uint8) bool {
	if loc > 7 {
		panic("invalid flag location: must be between 0 and 7")
	}

	// Determine the bit value at the given location
	bit := (value >> loc) & 0x1
	// Check if bit is set
	return bit != 0
}

func marshal32(data [32]byte) ([]byte, error) {
	buffer := make([]byte, 32*2+2)

	// Copy the 0x prefix into the buffer
	copy(buffer[:2], prefix0xString)
	// Hex-encode the copied value into the buffer
	hex.Encode(buffer[2:], data[:])

	return buffer, nil
}

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

func must[T any](t T, err error) T {
	if err != nil {
		panic(err)
	}

	return t
}
