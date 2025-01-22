package identifiers

import (
	"encoding/json"
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

		// Fingerprint
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

	// Test Fingerprint
	assert.Equal(t, [24]byte{
		0x10, 0x11, 0x12, 0x13, 0x14, 0x15, 0x16, 0x17, 0x18, 0x19, 0x1A, 0x1B,
		0x20, 0x21, 0x22, 0x23, 0x24, 0x25, 0x26, 0x27, 0x28, 0x29, 0x2A, 0x2B,
	}, id.Fingerprint())

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

func TestIdentifier_FromHex(t *testing.T) {
	t.Run("ValidHex", func(t *testing.T) {
		_, err := NewIdentifierFromHex(RandomAssetIDv0().AsIdentifier().Hex())
		require.NoError(t, err)

		_, err = NewIdentifierFromHex(trim0xPrefixString(RandomAssetIDv0().AsIdentifier().Hex()))
		require.NoError(t, err)
	})

	t.Run("InvalidHex", func(t *testing.T) {
		_, err := NewIdentifierFromHex("invalid-hex")
		require.EqualError(t, err, "encoding/hex: invalid byte: U+0069 'i'")

		_, err = NewIdentifierFromHex("0xf") // odd length
		require.EqualError(t, err, "encoding/hex: odd length hex string")
	})

	t.Run("MustFromHex", func(t *testing.T) {
		assert.Panics(t, func() { _ = MustIdentifierFromHex("0xFF") })
	})
}

func TestIdentifier_DeriveVariant(t *testing.T) {
	t.Run("SimpleDerivation", func(t *testing.T) {
		// Generate an asset ID with a zero variant (and standard = 0)
		identifier, err := GenerateAssetIDv0(RandomFingerprint(), 0, 0)
		require.NoError(t, err)

		// Attempt derivation without changing any flags
		derived, err := identifier.AsIdentifier().DeriveVariant(100, nil, nil)
		require.NoError(t, err)

		// Verify that the derived identifier has the new variant
		assert.Equal(t, uint32(100), derived.Variant())
		// Verify that the derived identifier has the same flags
		assert.Equal(t, identifier.AsIdentifier().Flags(), derived.Flags())
	})

	t.Run("DeriveWithSetFlag", func(t *testing.T) {
		// Generate an asset ID with a zero variant (and standard = 0)
		identifier, err := GenerateAssetIDv0(RandomFingerprint(), 0, 0)
		require.NoError(t, err)

		// Attempt derivation with a new variant and a flag set
		derived, err := identifier.AsIdentifier().DeriveVariant(100, []Flag{AssetStateful}, nil)
		require.NoError(t, err)

		// Verify that the derived identifier has the new variant
		assert.Equal(t, uint32(100), derived.Variant())
		// Verify that the derived identifier has the new flag set
		assert.True(t, must(derived.AsAssetID()).Flag(AssetStateful))
	})

	t.Run("DeriveWithUnsupportedSet", func(t *testing.T) {
		// Generate a logic ID with a zero variant
		identifier, err := GenerateLogicIDv0(RandomFingerprint(), 0)
		require.NoError(t, err)

		// Attempt derivation with a new variant and an unsupported flag set
		_, err = identifier.AsIdentifier().DeriveVariant(100, []Flag{AssetStateful}, nil)
		require.EqualError(t, err, ErrUnsupportedFlag.Error())
	})

	t.Run("DeriveWithUnsetFlag", func(t *testing.T) {
		// Generate an asset ID with a zero variant (and standard = 0)
		// Set the AssetStateful flag
		identifier, err := GenerateAssetIDv0(RandomFingerprint(), 0, 0, AssetStateful)
		require.NoError(t, err)
		require.True(t, identifier.Flag(AssetStateful))

		// Attempt derivation with a new variant and a flag unset
		derived, err := identifier.AsIdentifier().DeriveVariant(100, nil, []Flag{AssetStateful})
		require.NoError(t, err)

		// Verify that the derived identifier has the new variant
		assert.Equal(t, uint32(100), derived.Variant())
		// Verify that the derived identifier has the new flag unset
		assert.False(t, must(derived.AsAssetID()).Flag(AssetStateful))
	})

	t.Run("DeriveWithUnsupportedUnset", func(t *testing.T) {
		// Generate a logic ID with a zero variant
		identifier, err := GenerateLogicIDv0(RandomFingerprint(), 0, LogicIntrinsic)
		require.NoError(t, err)
		require.True(t, identifier.Flag(LogicIntrinsic))

		// Attempt derivation with a new variant and an unsupported flag unset
		_, err = identifier.AsIdentifier().DeriveVariant(100, nil, []Flag{AssetStateful})
		require.EqualError(t, err, ErrUnsupportedFlag.Error())
	})
}

func TestIdentifier_TextMarshal(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		identifier := RandomAssetIDv0().AsIdentifier()
		encoded, err := json.Marshal(identifier)
		require.NoError(t, err)

		var decoded Identifier

		require.NoError(t, json.Unmarshal(encoded, &decoded))
		require.Equal(t, identifier, decoded)
	})

	t.Run("MissingPrefix", func(t *testing.T) {
		var decoded Identifier

		require.Equal(t, json.Unmarshal([]byte(`"invalid-json"`), &decoded), ErrMissingHexPrefix)
	})

	t.Run("InvalidLength", func(t *testing.T) {
		var decoded Identifier

		require.Equal(t, json.Unmarshal([]byte(`"0xffabcd"`), &decoded), ErrInvalidLength)
	})

	t.Run("HexError", func(t *testing.T) {
		var decoded Identifier

		require.EqualError(t,
			json.Unmarshal([]byte(`"0xYY01001001020304050607081112131415161718212223242526272800000042"`), &decoded),
			"encoding/hex: invalid byte: U+0059 'Y'",
		)
	})
}
