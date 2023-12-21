package identifiers

import (
	"encoding/json"

	"github.com/pkg/errors"
)

// AssetID is a unique hex-encoded string identifier for assets in the MOI Protocol.
// It encodes within itself key properties about the asset itself such as the address
// at which the asset is deployed, the nature of it logic control and others.
//
// The Asset ID Standard is an extensible standard is compliant with the protocol specification
// for it at https://sarvalabs.notion.site/Asset-ID-Standard-e4fcd9151e7d4e7eb2447f1d8edf4672?pvs=4
type AssetID string

// AssetIdentifier is a representation of AssetID allows access to the
// values encoded within it such as the address or the asset standard.
//
// Each version of the AssetID standard is available
// as its own type that implements this interface
type AssetIdentifier interface {
	// Version returns the version of the AssetID standard for the AssetIdentifier
	Version() int
	// Address returns the 32-byte address associated with the AssetID
	Address() Address
	// AssetID returns the AssetIdentifier in its encoded representation as a AssetID
	AssetID() AssetID

	// Standard returns the asset standard number from the AssetID
	Standard() uint64
	// Dimension returns the asset dimension value from the AssetID
	Dimension() uint8

	// IsLogical returns whether the asset has some logic associated with it
	IsLogical() bool
	// IsStateful returns whether the asset has some stateful information such as its supply
	IsStateful() bool
}

// NewAssetID generates an AssetID from some arbitrary string,
// validating it in the process. It is version agnostic.
func NewAssetID(id string) (AssetID, error) {
	asset := AssetID(id)

	// Attempt to generate an identifier from the AssetID
	// This will fail if the AssetID is invalid in any way
	if _, err := asset.Identifier(); err != nil {
		return "", err
	}

	return asset, nil
}

// Bytes returns the AssetID as a []byte after being
// decoded from its hexadecimal string representation
// Panics if the AssetID is not a valid hex string
func (asset AssetID) Bytes() []byte {
	return must(decodeHexString(string(asset)))
}

// String returns the AssetID as a string.
// Implements the fmt.Stringer interface for AssetID.
func (asset AssetID) String() string {
	return "0x" + string(asset)
}

// Address returns the Address of the AssetID.
// Returns NilAddress if the AssetID is invalid.
// The AssetID standard expects the address to ALWAYS be the last 32 bytes
func (asset AssetID) Address() Address {
	// Error if length is too short
	if len(asset) < 64 {
		return NilAddress
	}

	// Trim the last 64 characters (32 bytes)
	addr := string(asset[len(asset)-64:])
	// Assertively decode into an Address
	return must(NewAddressFromHex(addr))
}

// Identifier returns a AssetIdentifier for the AssetID.
//
// This decodes the AssetID from its simple hex-encoded string
// format into a representation appropriate for the version of
// the AssetID to allow access to all the encoded fields within it.
//
// It can also be used to verify the integrity of the AssetID
func (asset AssetID) Identifier() (AssetIdentifier, error) {
	id, err := decodeHexString(string(asset))
	if err != nil {
		return nil, errors.Wrap(err, "invalid asset ID")
	}

	// We verify that there is at least 1 byte, so that
	// we can safely access the 0th index in the byte slice
	if len(id) < 1 {
		return nil, errors.New("invalid asset ID: missing version prefix")
	}

	// Determine the version of the AssetID and check if there are enough bytes
	switch version := int(id[0] & 0xF0); version {
	case 0:
		return decodeAssetIDv0(id)
	default:
		return nil, errors.Errorf("invalid asset ID: unsupported version: %v", version)
	}
}

// MarshalText implements the encoding.TextMarshaler interface for AssetID
func (asset AssetID) MarshalText() ([]byte, error) {
	return []byte(asset.String()), nil
}

// UnmarshalText implements the encoding.TextUnmarshaler interface for AssetID
func (asset *AssetID) UnmarshalText(text []byte) error {
	// Assert that the 0x prefix exists
	if !has0xPrefixBytes(text) {
		return ErrMissing0xPrefix
	}

	// Trim the 0x prefix
	text = trim0xPrefixBytes(text)
	// Generate an identifier for the AssetID
	if _, err := AssetID(text).Identifier(); err != nil {
		return err
	}

	*asset = AssetID(text)

	return nil
}

// MarshalJSON implements the json.Marshaler interface for AssetID
func (asset AssetID) MarshalJSON() ([]byte, error) {
	return json.Marshal(asset.String())
}

// UnmarshalJSON implements the json.Unmarshaler interface for AssetID
func (asset *AssetID) UnmarshalJSON(data []byte) error {
	var decoded string

	// Decode the JSON data into a string
	if err := json.Unmarshal(data, &decoded); err != nil {
		return err
	}

	// Assert that the 0x prefix exists
	if !has0xPrefixString(decoded) {
		return ErrMissing0xPrefix
	}

	// Trim the 0x prefix
	decoded = trim0xPrefixString(decoded)
	// Generate an identifier for the AssetID
	if _, err := AssetID(decoded).Identifier(); err != nil {
		return err
	}

	*asset = AssetID(decoded)

	return nil
}
