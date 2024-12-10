package identifiers

import (
	"encoding/binary"
	"encoding/hex"
	"errors"
	"fmt"
	"math"
	"math/rand/v2"
)

// AssetID is a unique identifier for an asset in the MOI Protocol.
// It is 32 bytes long and its metadata is structured as follows:
//   - Tag: The first byte contains the tag for the asset identifier.
//   - Flags: The second byte contains flags for the asset identifier.
//   - Standard: The next 2 bytes contain the standard for the asset.
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
// The given value is trimmed/padded to 32 bytes and validated into an AssetID.
func NewAssetIDFromBytes(data []byte) (AssetID, error) {
	return NewAssetID(trim32(data))
}

// NewAssetIDFromHex creates a new AssetID from the given hex string.
// The given value is hex-decoded (must contain 0x prefix),
// trimmed/padded to 32 bytes and validated into an AssetID.
func NewAssetIDFromHex(data string) (AssetID, error) {
	// Decode the given hex string into []byte
	decoded, err := decodeHexString(data)
	if err != nil {
		return Nil, err
	}

	// Create a new AssetID from the decoded value
	return NewAssetIDFromBytes(decoded)
}

// Tag returns the IdentifierTag for the AssetID.
func (asset AssetID) Tag() IdentifierTag {
	return IdentifierTag(asset[0])
}

// AccountID returns the 24-byte account ID from the AssetID.
func (asset AssetID) AccountID() [24]byte {
	return trimMid24(asset)
}

// Variant returns the 32-bit variant ID from the AssetID.
func (asset AssetID) Variant() uint32 {
	low4 := trimLow4(asset)
	return binary.BigEndian.Uint32(low4[:])
}

// Standard returns the 16-bit standard for the AssetID.
func (asset AssetID) Standard() uint16 {
	// get the standard from the 2nd and 3rd bytes
	return binary.BigEndian.Uint16(asset[2:4])
}

// AssetFlag is a flag specifier for AssetID's Flags.
// Returns the minimum supported version and the bit location of the flag.
//
// All AssetID flags are located in the [1] index.
// Use the AssetID.Flag() method to check if a specific flag is set.
type AssetFlag func() (version, location uint8)

// AssetFlagLogical is an AssetFlag for the IsLogical flag on AssetID
// It is supported from version 0 and is located at the 0th flag bit.
//
// The IsLogical flag indicates that the asset has some logic associated with it.
func AssetFlagLogical() (uint8, uint8) { return 0, 0x0 }

// AssetFlagStateful is an AssetFlag for the IsStateful flag on AssetID
// It is supported from version 0 and is located at the 1st flag bit.
//
// The IsStateful flag indicates that the asset has some stateful information such as its supply.
func AssetFlagStateful() (uint8, uint8) { return 0, 0x1 }

// Flag returns if the given AssetFlag is set on the AssetID.
//
// If the specified flag is not supported by the AssetID's version,
// it will return FALSE, regardless of the actual flag value.
func (asset AssetID) Flag(flag AssetFlag) bool {
	ver, loc := flag()

	// Check if the flag is supported by the asset ID version
	// If not supported, return FALSE, regardless of the actual flag value
	if asset.Tag().Version() > ver {
		return false
	}

	// Check if flag bit is set
	return isBitSet(asset[1], loc)
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
		return errors.New("invalid tag: not an asset identifier")
	}

	// Perform checks based on the asset ID version
	switch asset.Tag().Version() {
	case 0:
		// Check if the flags for any position apart from 0 and 1 are set
		if (asset[1] & byte(0b00111111)) != 0 {
			return errors.New("invalid flags: malformed flags for asset identifier v0")
		}

	default:
		return ErrUnsupportedVersion
	}

	return nil
}

// Bytes returns the AssetID as a []byte
func (asset AssetID) Bytes() []byte { return asset[:] }

// String returns the AssetID as a hex-encoded string.
// This is identical to AssetID.Hex() but is required for the fmt.Stringer interface
func (asset AssetID) String() string { return asset.Hex() }

// Hex returns the AssetID as a hex-encoded string with the 0x prefix
func (asset AssetID) Hex() string {
	return prefix0xString + hex.EncodeToString(asset[:])
}

// IsNil returns if the AssetID is nil, i.e., 0x000..000
func (asset AssetID) IsNil() bool {
	return asset == Nil
}

// AsIdentifier returns the AssetID as an AssetID.
func (asset AssetID) AsIdentifier() Identifier {
	return Identifier(asset)
}

// MarshalText implements the encoding.TextMarshaler interface for AssetID
func (asset *AssetID) MarshalText() ([]byte, error) {
	return marshal32(*asset)
}

// UnmarshalText implements the encoding.TextUnmarshaler interface for AssetID
func (asset *AssetID) UnmarshalText(data []byte) error {
	decoded, err := unmarshal32(data)
	if err != nil {
		return err
	}

	*asset = decoded
	return nil
}

// GenerateAssetIDv0 creates a new AssetID for v0 with the given parameters.
// Returns an error if unsupported flags are used.
//
// [tag:1][{isLogical}{isStateful}{reserved:6}][standard:2][account:24][variant:4]
func GenerateAssetIDv0(account [24]byte, variant uint32, standard uint16, flagOpts ...AssetFlag) (AssetID, error) {
	// Create the metadata buffer
	// [tag][flags][standard]
	metadata := make([]byte, 4)

	// Attach the tag for AssetID v0
	metadata[0] = byte(TagAssetV0)
	// Attach the flags to the metadata
	for _, opt := range flagOpts {
		// Get the minimum supported version and the location of the flag
		minver, loc := opt()
		// Check that flag is supported by version 0
		if 0 < minver {
			return Nil, ErrUnsupportedFlag
		}

		metadata[1] |= 1 << loc
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
// random account ID, variant ID, standard, and flags.
//   - There is a 50% chance that the IsLogical flag will be set.
//   - There is a 50% chance that the IsStateful flag will be set.
func RandomAssetIDv0() AssetID {
	flags := make([]AssetFlag, 0, 2)

	if rand.Int64() > 0 {
		flags = append(flags, AssetFlagLogical)
	}

	if rand.Int64() > 0 {
		flags = append(flags, AssetFlagStateful)
	}

	// Safe to ignore error as the flags are supported
	asset, _ := GenerateAssetIDv0(RandomAccountID(), rand.Uint32(), uint16(rand.UintN(math.MaxUint16)), flags...)

	return asset
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
