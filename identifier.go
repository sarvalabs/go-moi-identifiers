package identifiers

import (
	"encoding/binary"
	"encoding/hex"
)

// Identifier represents a unique 32-byte (256-bit) identifier
// This is the base type for all identifiers in the MOI Protocol.
//
// Every identifier is composed of 3 parts:
//   - Metadata: The 4 most-significant bytes
//   - AccountID: The 24 middle bytes
//   - Variant: The 4 least-significant bytes
//
// The first byte of the metadata contain a tag represented by IdentifierTag,
// which itself comprises a kind and version, each represented by 4 bits on the tag.
// Apart from the tag, the metadata contains 1 byte for flags and upto 2 bytes of
// additional data that can be used by different kinds of identifiers as required.
//
// The next 24 bytes represent the account ID, which is unique to each kind of identifier.
// The last 4 bytes represent a 32-bit variant number, which can be used for sub-identifiers.
type Identifier [32]byte

// IdentifierKind represents the kinds of recognized identifiers.
// The identifiers specification allows for upto 16 (2^4) kinds,
// but only the following are currently supported.
type IdentifierKind byte

const (
	KindParticipant IdentifierKind = iota
	KindGuardian
	KindAsset
	KindLogic
)

// maxIdentifierKind represents the maximum supported IdentifierKind value
const maxIdentifierKind = KindLogic

// IdentifierTag represents the tag of an identifier.
// The first 4-bit nibble represents the kind of the identifier (IdentifierKind),
// and the second 4-bit nibble represents the version for that identifier kind.
type IdentifierTag byte

const (
	TagParticipantV0 = IdentifierTag((KindParticipant << 4) | 0)
	TagGuardianV0    = IdentifierTag((KindGuardian << 4) | 0)
	TagAssetV0       = IdentifierTag((KindAsset << 4) | 0)
	TagLogicV0       = IdentifierTag((KindLogic << 4) | 0)
)

// Kind returns the IdentifierKind from the IdentifierTag
func (tag IdentifierTag) Kind() IdentifierKind {
	// Determine the kind from the upper 4 bits
	return IdentifierKind(tag >> 4)
}

// Version returns the version from the IdentifierTag
func (tag IdentifierTag) Version() uint8 {
	// Determine the version from the lower 4 bits
	return uint8(tag & 0x0F)
}

// Validate checks if the IdentifierTag is valid and returns an error if not.
// An error is returned if the version is not supported or the kind is invalid
func (tag IdentifierTag) Validate() error {
	if tag.Version() != 0 {
		return ErrUnsupportedVersion
	}

	if tag.Kind() > maxIdentifierKind {
		return ErrUnsupportedKind
	}

	return nil
}

// Tag returns the IdentifierTag from the Identifier
func (id Identifier) Tag() IdentifierTag { return IdentifierTag(id[0]) }

// Metadata returns the 4 most-significant bytes of the Identifier
func (id Identifier) Metadata() [4]byte { return trimHigh4(id) }

// AccountID returns 24-byte account ID from the Identifier
func (id Identifier) AccountID() [24]byte { return trimMid24(id) }

// Variant returns the 32-bit variant ID from the Identifier
func (id Identifier) Variant() uint32 {
	low4 := trimLow4(id)
	return binary.BigEndian.Uint32(low4[:])
}

// Bytes returns the Identifier as a []byte
func (id Identifier) Bytes() []byte { return id[:] }

// String returns the Identifier as a hex-encoded string.
// This is identical to Identifier.Hex() but is required for the fmt.Stringer interface
func (id Identifier) String() string { return id.Hex() }

// Hex returns the Identifier as a hex-encoded string with the 0x prefix
func (id Identifier) Hex() string { return prefix0xString + hex.EncodeToString(id[:]) }

// IsNil returns if the Identifier is nil, i.e., 0x000..000
func (id Identifier) IsNil() bool { return id == Nil }

// AsAssetID returns the Identifier as an AssetID.
// Returns an error if the Identifier is not a valid AssetID
func (id Identifier) AsAssetID() (AssetID, error) { return NewAssetID(id) }

// MarshalText implements the encoding.TextMarshaler interface for Identifier
func (id *Identifier) MarshalText() ([]byte, error) {
	return marshal32(*id)
}

// UnmarshalText implements the encoding.TextUnmarshaler interface for Identifier
func (id *Identifier) UnmarshalText(data []byte) error {
	decoded, err := unmarshal32(data)
	if err != nil {
		return err
	}

	*id = decoded
	return nil
}
