package identifiers

import (
	"encoding"
	"encoding/binary"
	"encoding/hex"
	"errors"
	"fmt"
	"math"
	"math/rand/v2"
)

// AssetID is a unique identifier for an asset in the MOI Protocol.
// It is 32 bytes long and its first 4 bytes are structured as follows:
//   - Tag: The first byte contains the tag for the asset identifier.
//   - Flags: The second byte contains flags for the asset identifier.
//   - Metadata: The next 2 bytes contain the standard for the asset.
//
// Like all identifiers, the AssetID also contains an AccountID and a Variant ID.
// Flags of an AssetID are specific to a version and are invalid if set in an unsupported version.
type AssetID [32]byte

// NewAssetID creates a new AssetID from the 32-byte value.
// It returns an error if the given data is not a valid AssetID.
func NewAssetID(data [32]byte) (AssetID, error) {
	// Convert the data into an AssetID
	assetID := AssetID(data)
	// Validate the AssetID
	if err := assetID.Validate(); err != nil {
		return Nil, err
	}

	return assetID, nil
}

// NewAssetIDFromBytes creates a new AssetID from the given byte slice.
// The given value must have a length of 32 and validate into an AssetID.
func NewAssetIDFromBytes(data []byte) (AssetID, error) {
	// Check length of the data
	if len(data) != 32 {
		return Nil, errors.New("invalid length: asset id must be 32 bytes")
	}

	return NewAssetID([32]byte(data))
}

// NewAssetIDFromHex creates a new AssetID from the given hex string.
// The given value must decode as hexadecimal string (0x prefix is optional),
// with a length of 64 characters (32 bytes) and validate into an AssetID.
func NewAssetIDFromHex(data string) (AssetID, error) {
	// Decode the given hex string into []byte
	decoded, err := decodeHexString(data)
	if err != nil {
		return Nil, err
	}

	// Create a new AssetID from the decoded value
	// Length check is performed in NewAssetIDFromBytes
	return NewAssetIDFromBytes(decoded)
}

// MustAssetID is an enforced version of NewAssetID.
// Panics if an error occurs. Use with caution.
func MustAssetID(data [32]byte) AssetID { return must(NewAssetID(data)) }

// MustAssetIDFromBytes is an enforced version of NewAssetIDFromBytes.
// Panics if an error occurs. Use with caution.
func MustAssetIDFromBytes(data []byte) AssetID { return must(NewAssetIDFromBytes(data)) }

// MustAssetIDFromHex is an enforced version of NewAssetIDFromHex.
// Panics if an error occurs. Use with caution.
func MustAssetIDFromHex(data string) AssetID { return must(NewAssetIDFromHex(data)) }

// Bytes returns the AssetID as a []byte
func (asset AssetID) Bytes() []byte { return asset[:] }

// String returns the AssetID as a hex-encoded string.
// This is identical to AssetID.Hex() but is required for the fmt.Stringer interface
func (asset AssetID) String() string { return asset.Hex() }

// Hex returns the AssetID as a hex-encoded string with the 0x prefix
func (asset AssetID) Hex() string {
	return prefix0xString + hex.EncodeToString(asset[:])
}

// AsIdentifier returns the AssetID as an AssetID.
func (asset AssetID) AsIdentifier() Identifier {
	return Identifier(asset)
}

// Tag returns the IdentifierTag for the AssetID.
func (asset AssetID) Tag() IdentifierTag {
	return IdentifierTag(asset[0])
}

// AccountID returns the 24-byte account ID from the AssetID.
func (asset AssetID) AccountID() [24]byte {
	return trimAccount(asset)
}

// Variant returns the 32-bit variant ID from the AssetID.
func (asset AssetID) Variant() uint32 {
	variant := trimVariant(asset)
	return binary.BigEndian.Uint32(variant[:])
}

// IsVariant returns if the AssetID has a non-zero variant ID
func (asset AssetID) IsVariant() bool {
	variant := trimVariant(asset)
	return !(variant[0] == 0 && variant[1] == 0 && variant[2] == 0 && variant[3] == 0)
}

// Standard returns the 16-bit standard for the AssetID.
func (asset AssetID) Standard() uint16 {
	// get the standard from the 2nd and 3rd bytes
	return binary.BigEndian.Uint16(asset[2:4])
}

// Flag returns if the given Flag is set on the AssetID.
//
// If the specified flag is not supported by the AssetID,
// it will return False, regardless of the actual flag value.
func (asset AssetID) Flag(flag Flag) bool {
	// Check if the flag is supported by AssetID.
	// If not supported, return FALSE, regardless of the actual flag value
	if !flag.Supports(asset.Tag()) {
		return false
	}

	return getFlag(asset[1], flag.index)
}

// Validate checks if the AssetID is valid.
// An error is returned if the AssetID has an invalid tag or contains unsupported flags.
func (asset AssetID) Validate() error {
	// Check basic validity of the identifier tag
	if err := asset.Tag().Validate(); err != nil {
		return fmt.Errorf("invalid tag: %w", err)
	}

	// Check if the tag is an asset tag
	if asset.Tag().Kind() != KindAsset {
		return errors.New("invalid tag: not an asset id")
	}

	// Check that there are no unsupported flags set
	if (asset[1] & flagMasks[asset.Tag()]) != 0 {
		return errors.New("invalid flags: unsupported flags for asset id")
	}

	return nil
}

var (
	// Ensure AssetID implements text marshaling interfaces
	_ encoding.TextMarshaler   = (*AssetID)(nil)
	_ encoding.TextUnmarshaler = (*AssetID)(nil)
)

// MarshalText implements the encoding.TextMarshaler interface for AssetID
func (asset AssetID) MarshalText() ([]byte, error) {
	return marshal32(asset)
}

// UnmarshalText implements the encoding.TextUnmarshaler interface for AssetID
func (asset *AssetID) UnmarshalText(data []byte) error {
	decoded, err := unmarshal32(data)
	if err != nil {
		return err
	}

	if err = AssetID(decoded).Validate(); err != nil {
		return err
	}

	*asset = decoded
	return nil
}

// GenerateAssetIDv0 creates a new AssetID for v0 with the given parameters.
// Returns an error if unsupported flags are used.
//
// [tag:1][{systemic}{reserved:5}{logical}{stateful}][standard:2][account:24][variant:4]
func GenerateAssetIDv0(account [24]byte, variant uint32, standard uint16, flags ...Flag) (AssetID, error) {
	// Create the metadata buffer
	// [tag][flags][standard]
	metadata := make([]byte, 4)
	// Attach the tag for AssetID v0
	metadata[0] = byte(TagAssetV0)

	// Attach the flags to the metadata
	for _, flag := range flags {
		// Check if the given flag is supported by AssetID v0
		if !flag.Supports(TagAssetV0) {
			return Nil, ErrUnsupportedFlag
		}

		// Set the flag in the metadata
		metadata[1] = setFlag(metadata[1], flag.index, true)
	}

	// Encode and attach the standard to the metadata
	binary.BigEndian.PutUint16(metadata[2:], standard)

	// Order the asset ID buffer
	// [metadata][account][variant]
	buffer := make([]byte, 0, 32)
	buffer = append(buffer, metadata...)
	buffer = append(buffer, account[:]...)
	// Append 4 bytes for the variant and encode the value into it
	buffer = append(buffer, make([]byte, 4)...)
	binary.BigEndian.PutUint32(buffer[28:], variant)

	return AssetID(buffer), nil
}

// RandomAssetIDv0 creates a random v0 AssetID with a
// random account ID, variant ID, standard and flags.
//   - There is a 50% chance that the AssetLogical flag will be set.
//   - There is a 50% chance that the AssetStateful flag will be set.
//   - There is a 0% chance that the Systemic flag will be set.
func RandomAssetIDv0() AssetID {
	flags := make([]Flag, 0, 2)

	if rand.Int64() > 0 {
		flags = append(flags, AssetLogical)
	}

	if rand.Int64() > 0 {
		flags = append(flags, AssetStateful)
	}

	// Safe to ignore error as the flags are supported
	asset, _ := GenerateAssetIDv0(RandomAccountID(), rand.Uint32(), uint16(rand.UintN(math.MaxUint16)), flags...)

	return asset
}
