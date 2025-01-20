package identifiers

import (
	"encoding"
	"encoding/binary"
	"encoding/hex"
)

// IdentifierKind represents the kinds of recognized identifiers.
type IdentifierKind byte

const (
	KindParticipant IdentifierKind = iota
	KindAsset
	KindLogic
)

const (
	maxIdentifierKind = KindLogic
	identifierV0      = 0
)

// kindSupport is a map of IdentifierKind to the maximum supported version.
var kindSupport = map[IdentifierKind]uint8{
	KindParticipant: 0,
	KindAsset:       0,
	KindLogic:       0,
}

// IdentifierTag represents the tag of an identifier.
// The first 4-bit nibble represents the kind of the identifier (IdentifierKind),
// and the second 4-bit nibble represents the version for that identifier kind.
//
// While the version is currently set to 0 for all kinds, it allows for future
// changes to the identifier format while maintaining backward compatibility.
//
// This format allows for up to 16 different kinds of identifiers and 16 different
// versions for each kind. While this headroom is excessive for current requirements,
// and could be optimized further, using the nibble as the smallest unit, allows for
// easily recognizing the kind and version of an identifier in its hexadecimal format.
type IdentifierTag byte

const (
	TagParticipantV0 = IdentifierTag((KindParticipant << 4) | identifierV0)
	TagAssetV0       = IdentifierTag((KindAsset << 4) | identifierV0)
	TagLogicV0       = IdentifierTag((KindLogic << 4) | identifierV0)
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
	// Check if the kind is under the maximum supported kind
	if tag.Kind() > maxIdentifierKind {
		return ErrUnsupportedKind
	}

	// Check if the version is supported for the kind
	if tag.Version() > kindSupport[tag.Kind()] {
		return ErrUnsupportedVersion
	}

	return nil
}

// Identifier represents a unique 32-byte (256-bit) identifier
// This is the base type for all identifiers in the MOI Protocol.
//
// Every identifier is composed of 5 parts:
//   - Tag: The most-significant byte
//   - Flags: The second byte
//   - Metadata: The 3rd & 4th byte
//   - AccountID: The next 24 middle bytes
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

// Bytes returns the Identifier as a []byte
func (id Identifier) Bytes() []byte { return id[:] }

// String returns the Identifier as a hex-encoded string.
// This is identical to Identifier.Hex() but is required for the fmt.Stringer interface
func (id Identifier) String() string { return id.Hex() }

// Hex returns the Identifier as a hex-encoded string with the 0x prefix
func (id Identifier) Hex() string { return prefix0xString + hex.EncodeToString(id[:]) }

// IsNil returns if the Identifier is nil, i.e., 0x000..000
func (id Identifier) IsNil() bool { return id == Nil }

// Tag returns the IdentifierTag from the Identifier
func (id Identifier) Tag() IdentifierTag { return IdentifierTag(id[0]) }

// Flags returns the byte of flag bits from the Identifier
func (id Identifier) Flags() byte { return id[1] }

// Metadata returns the 3rd & 4th bytes of the Identifier
func (id Identifier) Metadata() [2]byte { return [2]byte{id[2], id[3]} }

// AccountID returns 24-byte account ID from the Identifier
func (id Identifier) AccountID() [24]byte { return trimAccount(id) }

// Variant returns the 32-bit variant ID from the Identifier
func (id Identifier) Variant() uint32 {
	variant := trimVariant(id)
	return binary.BigEndian.Uint32(variant[:])
}

// IsVariant returns if the Identifier has a non-zero variant ID
func (id Identifier) IsVariant() bool {
	variant := trimVariant(id)
	return !(variant[0] == 0 && variant[1] == 0 && variant[2] == 0 && variant[3] == 0)
}

// DeriveVariant returns a new Identifier with the given variant ID and specified flags set/unset.
// Returns an error if the given flags are not supported for the Identifier tag.
func (id Identifier) DeriveVariant(variant uint32, set []Flag, unset []Flag) (Identifier, error) {
	var derived Identifier

	// Copy the original identifier
	copy(derived[:], id[:])
	// Encode the new variant ID
	binary.BigEndian.PutUint32(derived[28:], variant)

	for _, flag := range set {
		// Check if the given flag is supported by identifier tag
		if !flag.Supports(derived.Tag()) {
			return Nil, ErrUnsupportedFlag
		}

		// Set the flag
		derived[1] = setFlag(derived[1], flag.index, true)
	}

	for _, flag := range unset {
		// Check if the given flag is supported by identifier tag
		if !flag.Supports(derived.Tag()) {
			return Nil, ErrUnsupportedFlag
		}

		// Unset the flag
		derived[1] = setFlag(derived[1], flag.index, false)
	}

	return derived, nil
}

// AsParticipantID returns the Identifier as a ParticipantID.
// Returns an error if the Identifier is not a valid ParticipantID
func (id Identifier) AsParticipantID() (ParticipantID, error) { return NewParticipantID(id) }

// AsAssetID returns the Identifier as an AssetID.
// Returns an error if the Identifier is not a valid AssetID
func (id Identifier) AsAssetID() (AssetID, error) { return NewAssetID(id) }

// AsLogicID returns the Identifier as a LogicID.
// Returns an error if the Identifier is not a valid LogicID
func (id Identifier) AsLogicID() (LogicID, error) { return NewLogicID(id) }

var (
	// Ensure Identifier implements text marshaling interfaces
	_ encoding.TextMarshaler   = (*Identifier)(nil)
	_ encoding.TextUnmarshaler = (*Identifier)(nil)
)

// MarshalText implements the encoding.TextMarshaler interface for Identifier
func (id Identifier) MarshalText() ([]byte, error) {
	return marshal32(id)
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
