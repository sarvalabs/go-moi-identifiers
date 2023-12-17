package identifiers

import (
	"bytes"
	"encoding/hex"
	"errors"
	"strings"
)

var (
	prefix0xString = "0x"
	prefix0xBytes  = []byte(prefix0xString)
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

var (
	ErrMissing0xPrefix = errors.New("missing '0x' prefix")
	ErrInvalidLength   = errors.New("invalid length")
)

func decodeHexString(str string) ([]byte, error) {
	// Trim the 0x prefix from the string (if it exists)
	str = trim0xPrefixString(str)

	decoded, err := hex.DecodeString(str)
	if err != nil {
		return nil, err
	}

	return decoded, nil
}

func must[T any](t T, err error) T {
	if err != nil {
		panic(err)
	}

	return t
}
