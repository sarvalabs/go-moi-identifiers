package identifiers

// Every identifier reserves its second byte (index 1) for some bit flags.
// These flags are used to provide additional information about the identifier.
// The flag indices start at 7 for the MSB and end at 0 for the LSB.

var (
	// Systemic is a Flag for the MSB on all identifiers flags regardless of the kind.
	// It indicates that the account associated with identifier belongs to the system.
	// Supported from v0 for all identifiers
	Systemic = Flag{
		index: 7,
		support: map[IdentifierKind]uint8{
			KindParticipant: 0,
			KindAsset:       0,
			KindLogic:       0,
		},
	}

	// AssetStateful is a Flag on AssetID for the Stateful flag on its 0th bit.
	// It indicates that the asset has some stateful information such as its supply.
	// Supported from v0 of AssetID
	AssetStateful = makeFlag(KindAsset, 0, 0)
	// AssetLogical is a Flag on AssetID for the Logical flag on its 1st bit.
	// It indicates that the asset has some logic associated with it.
	// Supported from v0 of AssetID
	AssetLogical = makeFlag(KindAsset, 1, 0)

	// LogicIntrinsic is a Flag on LogicID for the Intrinsic flag on its 0th bit.
	// It indicates that the logic manages some intrinsic state
	// Supported from v0 of LogicID
	LogicIntrinsic = makeFlag(KindLogic, 0, 0)
	// LogicExtrinsic is a Flag on LogicID for the Extrinsic flag on its 1st bit.
	// It indicates that the logic manages some extrinsic state
	// Supported from v0 of LogicID
	LogicExtrinsic = makeFlag(KindLogic, 1, 0)
	// LogicAuxiliary is a Flag on LogicID for the Auxiliary flag on its 2nd bit.
	// It indicates that the logic is attached as an auxiliary to another object.
	// Supported from v0 of LogicID
	LogicAuxiliary = makeFlag(KindLogic, 2, 0)
)

// Flag represents a flag specifier for an identifier.
type Flag struct {
	// the bit index of the flag
	index uint8
	// the supported identifier kinds mapped to minimum supported version
	support map[IdentifierKind]uint8
}

// Supports returns if the flag is supported by the given kind.
func (flag Flag) Supports(tag IdentifierTag) bool {
	// Check if the kind is supported by the flag & obtain version
	version, ok := flag.support[tag.Kind()]
	if !ok {
		return false
	}

	// Check if the version is supported by flag
	return tag.Version() >= version
}

// getFlag retrieves a flag value from a given flag set and an index.
func getFlag(value byte, index uint8) bool {
	// Determine the bit value at the given index
	bit := value & (1 << index)
	// Check if flag is set
	return bit != 0
}

// setFlag updates the flag value for a given flag
// set and index and returns the modified flag set.
//
// Whether to set/unset the flag is determined by the flag input.
func setFlag(value byte, index uint8, flag bool) byte {
	if flag {
		value |= 1 << index // Set the bit value at the given index
	} else {
		value &^= 1 << index // Unset the bit value at the given index
	}

	return value
}

// makeFlag is used to construct a valid Flag object
// which is only supported by a single IdentifierKind
func makeFlag(kind IdentifierKind, index uint8, version uint8) Flag {
	if index > 7 {
		panic("invalid flag location: must be between 0 and 7")
	}

	if version > 7 {
		panic("invalid flag version: must be between 0 and 7")
	}

	return Flag{
		index:   index,
		support: map[IdentifierKind]uint8{kind: version},
	}
}

// flagMasks represent the mask of supported flags for an IdentifierTag.
// Can be accessed with IdentifierTag.FlagMask().
//
// A set bit indicates that position is not allowed for the tag,
// While an unset bit indicates it is a supported flag for the tag.
var flagMasks = map[IdentifierTag]byte{
	TagParticipantV0: 0b01111111,
	TagLogicV0:       0b01111000,
	TagAssetV0:       0b01111100,
}
