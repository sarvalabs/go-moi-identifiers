package identifiers

import (
	"encoding"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestIdentifierTag(t *testing.T) {
	tests := []struct {
		name            string
		tag             IdentifierTag
		expectedKind    IdentifierKind
		expectedVersion uint8
		expectedValid   bool
	}{
		{
			name:            "Participant V0",
			tag:             TagParticipantV0,
			expectedKind:    KindParticipant,
			expectedVersion: 0,
			expectedValid:   true,
		},
		{
			name:            "Asset V0",
			tag:             TagAssetV0,
			expectedKind:    KindAsset,
			expectedVersion: 0,
			expectedValid:   true,
		},
		{
			name:            "Logic V0",
			tag:             TagLogicV0,
			expectedKind:    KindLogic,
			expectedVersion: 0,
			expectedValid:   true,
		},
		{
			name:            "Invalid Version",
			tag:             IdentifierTag((KindParticipant << 4) | 1),
			expectedKind:    KindParticipant,
			expectedVersion: 1,
			expectedValid:   false,
		},
		{
			name:            "Invalid Kind",
			tag:             IdentifierTag((0x0F << 4) | identifierV0),
			expectedKind:    IdentifierKind(0x0F),
			expectedVersion: 0,
			expectedValid:   false,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			// Test Kind
			assert.Equal(t, tc.expectedKind, tc.tag.Kind())

			// Test Version
			assert.Equal(t, tc.expectedVersion, tc.tag.Version())

			// Test Validation
			err := tc.tag.Validate()
			if tc.expectedValid {
				assert.NoError(t, err)
			} else {
				assert.Error(t, err)
			}
		})
	}
}

func TestIdentifier(t *testing.T) {
	data := [32]byte{
		byte(TagParticipantV0), // Tag
		0b00000001,             // Flags
		0x02, 0x03,             // Metadata

		// Account ID
		0x10, 0x11, 0x12, 0x13, 0x14, 0x15, 0x16, 0x17, 0x18, 0x19, 0x1A, 0x1B,
		0x20, 0x21, 0x22, 0x23, 0x24, 0x25, 0x26, 0x27, 0x28, 0x29, 0x2A, 0x2B,

		0x30, 0x31, 0x32, 0x33, // Variant
	}

	// Create a test Identifier
	id := Identifier(data)

	// Test Tag, Flags & Metadata
	assert.Equal(t, TagParticipantV0, id.Tag())
	assert.Equal(t, byte(0b00000001), id.Flags())
	assert.Equal(t, [2]byte{0x02, 0x03}, id.Metadata())

	// Test AccountID
	assert.Equal(t, [24]byte{
		0x10, 0x11, 0x12, 0x13, 0x14, 0x15, 0x16, 0x17, 0x18, 0x19, 0x1A, 0x1B,
		0x20, 0x21, 0x22, 0x23, 0x24, 0x25, 0x26, 0x27, 0x28, 0x29, 0x2A, 0x2B,
	}, id.AccountID())

	// Test Variant
	assert.Equal(t, uint32(0x30313233), id.Variant())
	// Test IsVariant
	assert.True(t, id.IsVariant())

	// Test IsNil
	assert.False(t, id.IsNil())
	assert.True(t, Identifier(Nil).IsNil())

	// Test Bytes method
	assert.Equal(t, id[:], id.Bytes())

	// Test String & Hex method
	expectedHex := "0x00010203101112131415161718191a1b202122232425262728292a2b30313233"
	assert.Equal(t, expectedHex, id.String())
	assert.Equal(t, expectedHex, id.Hex())
}

func TestIdentifier_DeriveVariant(t *testing.T) {
	// Create a test Identifier
	identifier := RandomAssetIDv0().AsIdentifier()
	variant := identifier.Variant()

	// Derive a variant without changing any flags
	derivedOne, err := identifier.DeriveVariant(variant+100, nil, nil)
	require.NoError(t, err)

	// Verify that the derived identifier has the new variant
	assert.Equal(t, variant+100, derivedOne.Variant())
	assert.Equal(t, identifier.Metadata(), derivedOne.Metadata())

	// Derive another variant and set a flag
	derivedTwo, err := derivedOne.DeriveVariant(variant+200, []Flag{AssetStateful}, nil)
	require.NoError(t, err)

	// Verify that the derived identifier has the new variant
	assert.Equal(t, variant+200, derivedTwo.Variant())
	// Verify that the derived identifier has the new flag set
	assert.True(t, must(derivedTwo.AsAssetID()).Flag(AssetStateful))

	// Derive another variant and unset a flag
	derivedThree, err := derivedTwo.DeriveVariant(variant+300, nil, []Flag{AssetStateful})
	require.NoError(t, err)

	// Verify that the derived identifier has the new variant
	assert.Equal(t, variant+300, derivedThree.Variant())
	// Verify that the derived identifier has the new flag unset
	assert.False(t, must(derivedThree.AsAssetID()).Flag(AssetStateful))

	// Test DeriveVariant with unsupported flag set
	_, err = derivedThree.DeriveVariant(variant+400, []Flag{LogicExtrinsic}, nil)
	require.EqualError(t, err, ErrUnsupportedFlag.Error())

	// Test DeriveVariant with unsupported flag unset
	_, err = derivedThree.DeriveVariant(variant+400, nil, []Flag{LogicExtrinsic})
	require.EqualError(t, err, ErrUnsupportedFlag.Error())

	// Test DeriveVariant with same variant
	_, err = derivedThree.DeriveVariant(derivedThree.Variant(), nil, nil)
	require.EqualError(t, err, "cannot derive with the same variant")
}

func TestIdentifier_TextMarshal(t *testing.T) {
	// Ensure Identifier implements text marshaling interfaces
	var _ encoding.TextMarshaler = (*Identifier)(nil)
	var _ encoding.TextUnmarshaler = (*Identifier)(nil)

	// Create a sample identifier
	var original Identifier
	for i := range original {
		original[i] = byte(i)
	}

	// Test MarshalText
	marshaledText, err := original.MarshalText()
	require.NoError(t, err, "MarshalText should not return an error")

	// Test UnmarshalText
	var unmarshaled Identifier
	err = unmarshaled.UnmarshalText(marshaledText)
	require.NoError(t, err, "UnmarshalText should not return an error")

	// Verify that unmarshaled matches original
	assert.Equal(t, original, unmarshaled, "Unmarshaled identifier should match original")

	// Test UnmarshalText with invalid data
	err = unmarshaled.UnmarshalText([]byte("invalid"))
	require.Error(t, err, "UnmarshalText should return an error for invalid data")
}
