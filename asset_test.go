package identifiers

import (
	"encoding/hex"
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAssetID(t *testing.T) {
	data := [32]byte{
		byte(TagAssetV0), // Tag
		0b00000001,       // Flags
		0x00, 0x10,       // Standard

		// Fingerprint
		0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08,
		0x11, 0x12, 0x13, 0x14, 0x15, 0x16, 0x17, 0x18,
		0x21, 0x22, 0x23, 0x24, 0x25, 0x26, 0x27, 0x28,

		0x00, 0x00, 0x00, 0x42, // Variant
	}

	// Create a test AssetID
	assetID, err := NewAssetID(data)
	require.NoError(t, err)

	// Test Tag
	assert.Equal(t, TagAssetV0, assetID.Tag())

	// Test Fingerprint
	assert.Equal(t, [24]byte{
		0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08,
		0x11, 0x12, 0x13, 0x14, 0x15, 0x16, 0x17, 0x18,
		0x21, 0x22, 0x23, 0x24, 0x25, 0x26, 0x27, 0x28,
	}, assetID.Fingerprint())

	// Test Variant
	assert.Equal(t, uint32(0x42), assetID.Variant())
	// Test IsVariant
	assert.True(t, assetID.IsVariant())

	// Test Standard
	assert.Equal(t, uint16(0x10), assetID.Standard())

	// Test Flags
	assert.True(t, assetID.Flag(AssetStateful))
	assert.False(t, assetID.Flag(AssetLogical))
	assert.False(t, assetID.Flag(Systemic))
	assert.False(t, assetID.Flag(LogicIntrinsic)) // unsupported flag on set bit

	// Test AsIdentifier
	identifier := Identifier(data[:])
	assert.Equal(t, identifier, assetID.AsIdentifier())

	// Test From Identifier
	converted, err := identifier.AsAssetID()
	require.NoError(t, err)
	require.Equal(t, assetID, converted)

	// Test Bytes
	assert.Equal(t, data[:], assetID.Bytes())

	// Test String & Hex
	expectedHex := "0x1001001001020304050607081112131415161718212223242526272800000042"
	assert.Equal(t, expectedHex, assetID.String())
	assert.Equal(t, expectedHex, assetID.Hex())
}

//nolint:dupl // similar functions
func TestAssetID_Constructor(t *testing.T) {
	t.Run("NewAssetID", func(t *testing.T) {
		t.Run("Valid", func(t *testing.T) {
			assetID, err := NewAssetID([32]byte{
				byte(TagAssetV0), // Tag
				0b00000000,       // Flags
				0x00, 0x01,       // Standard
				// Empty bytes for fingerprint and variant
			})

			require.NoError(t, err)
			require.NoError(t, assetID.Validate())
		})

		t.Run("InvalidTag", func(t *testing.T) {
			_, err := NewAssetID([32]byte{0xF0}) // Invalid tag kind
			require.EqualError(t, err, "invalid tag: unsupported tag kind")

			_, err = NewAssetID([32]byte{0x0F}) // Invalid tag version
			require.EqualError(t, err, "invalid tag: unsupported tag version")

			_, err = NewAssetID([32]byte{byte(TagLogicV0)}) // Invalid tag
			require.EqualError(t, err, "invalid tag: not an asset id")
		})

		t.Run("InvalidFlags", func(t *testing.T) {
			_, err := NewAssetID([32]byte{
				byte(TagAssetV0), // Tag
				0b11111111,       // Invalid flags
			})
			require.EqualError(t, err, "invalid flags: unsupported flags for asset id")
		})
	})

	t.Run("NewAssetIDFromBytes", func(t *testing.T) {
		// Less than 32 bytes
		t.Run("< 32 bytes", func(t *testing.T) {
			_, err := NewAssetIDFromBytes([]byte{byte(TagAssetV0), 0x00, 0x00, 0x01})
			require.EqualError(t, err, "invalid length: asset id must be 32 bytes")
		})

		// Exactly 32 bytes
		t.Run("= 32 bytes", func(t *testing.T) {
			assetID, err := NewAssetIDFromBytes([]byte{
				byte(TagAssetV0), // Tag
				0b00000000,       // Flags
				0x00, 0x01,       // Standard

				// Fingerprint
				0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08,
				0x11, 0x12, 0x13, 0x14, 0x15, 0x16, 0x17, 0x18,
				0x21, 0x22, 0x23, 0x24, 0x25, 0x26, 0x27, 0x28,

				0x00, 0x00, 0x00, 0x01, // Variant
			})

			require.NoError(t, err)
			require.NoError(t, assetID.Validate())
			require.Equal(t, AssetID{
				byte(TagAssetV0), 0x00, 0x00, 0x01,
				0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08,
				0x11, 0x12, 0x13, 0x14, 0x15, 0x16, 0x17, 0x18,
				0x21, 0x22, 0x23, 0x24, 0x25, 0x26, 0x27, 0x28,
				0x00, 0x00, 0x00, 0x01,
			}, assetID)
		})

		// More than 32 bytes
		t.Run("> 32 bytes", func(t *testing.T) {
			_, err := NewAssetIDFromBytes([]byte{
				byte(TagAssetV0), // Tag
				0b00000000,       // Flags
				0x00, 0x01,       // Standard

				// Fingerprint
				0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08,
				0x11, 0x12, 0x13, 0x14, 0x15, 0x16, 0x17, 0x18,
				0x21, 0x22, 0x23, 0x24, 0x25, 0x26, 0x27, 0x28,

				0x00, 0x00, 0x00, 0x01, // Variant
				0xFF, 0xFF, 0xFF, 0xFF, // Extra bytes
			})
			require.EqualError(t, err, "invalid length: asset id must be 32 bytes")
		})
	})

	t.Run("NewAssetIDFromHex", func(t *testing.T) {
		t.Run("ValidHex", func(t *testing.T) {
			assetID, err := NewAssetIDFromHex("0x" + hex.EncodeToString([]byte{
				byte(TagAssetV0), // Tag
				0b00000000,       // Flags
				0x00, 0x01,       // Standard

				// Fingerprint
				0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08,
				0x11, 0x12, 0x13, 0x14, 0x15, 0x16, 0x17, 0x18,
				0x21, 0x22, 0x23, 0x24, 0x25, 0x26, 0x27, 0x28,

				0x00, 0x00, 0x00, 0x01, // Variant
			}))

			require.NoError(t, err)
			require.NoError(t, assetID.Validate())
		})

		t.Run("ValidHexNoPrefix", func(t *testing.T) {
			assetID, err := NewAssetIDFromHex(hex.EncodeToString([]byte{
				byte(TagAssetV0), // Tag
				0b00000000,       // Flags
				0x00, 0x01,       // Standard

				// Fingerprint
				0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08,
				0x11, 0x12, 0x13, 0x14, 0x15, 0x16, 0x17, 0x18,
				0x21, 0x22, 0x23, 0x24, 0x25, 0x26, 0x27, 0x28,

				0x00, 0x00, 0x00, 0x01, // Variant
			}))

			require.NoError(t, err)
			require.NoError(t, assetID.Validate())
		})

		t.Run("InvalidHex", func(t *testing.T) {
			_, err := NewAssetIDFromHex("invalid-hex")
			require.EqualError(t, err, "encoding/hex: invalid byte: U+0069 'i'")

			_, err = NewAssetIDFromHex("0xf") // odd length
			require.EqualError(t, err, "encoding/hex: odd length hex string")
		})
	})

	t.Run("MustAssetID", func(t *testing.T) {
		t.Run("MustAssetID", func(t *testing.T) {
			assert.Panics(t, func() { _ = MustAssetID([32]byte{0xFF}) })
		})

		t.Run("MustAssetIDFromBytes", func(t *testing.T) {
			assert.Panics(t, func() { _ = MustAssetIDFromBytes([]byte{0xFF}) })
		})

		t.Run("MustAssetIDFromHex", func(t *testing.T) {
			assert.Panics(t, func() { _ = MustAssetIDFromHex("0xFF") })
		})
	})
}

func TestAssetID_TextMarshal(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		// Create a test AssetID
		data := [32]byte{
			byte(TagAssetV0), // Tag
			0b00000001,       // Flags
			0x00, 0x10,       // Standard

			// Fingerprint
			0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08,
			0x11, 0x12, 0x13, 0x14, 0x15, 0x16, 0x17, 0x18,
			0x21, 0x22, 0x23, 0x24, 0x25, 0x26, 0x27, 0x28,

			0x00, 0x00, 0x00, 0x42, // Variant
		}

		assetID, err := NewAssetID(data)
		require.NoError(t, err)

		encoded, err := json.Marshal(assetID)
		require.NoError(t, err)
		require.Equal(t, `"0x1001001001020304050607081112131415161718212223242526272800000042"`, string(encoded))

		var decoded AssetID

		require.NoError(t, json.Unmarshal(encoded, &decoded))
		require.Equal(t, assetID, decoded)
	})

	t.Run("MissingPrefix", func(t *testing.T) {
		var decoded AssetID

		require.Equal(t, json.Unmarshal([]byte(`"invalid-json"`), &decoded), ErrMissingHexPrefix)
	})

	t.Run("InvalidLength", func(t *testing.T) {
		var decoded AssetID

		require.Equal(t, json.Unmarshal([]byte(`"0xffabcd"`), &decoded), ErrInvalidLength)
	})

	t.Run("HexError", func(t *testing.T) {
		var decoded AssetID

		require.EqualError(t,
			json.Unmarshal([]byte(`"0xYY01001001020304050607081112131415161718212223242526272800000042"`), &decoded),
			"encoding/hex: invalid byte: U+0059 'Y'",
		)
	})
}

func TestAssetID_Generation(t *testing.T) {
	t.Run("v0", func(t *testing.T) {
		t.Run("Generate", func(t *testing.T) {
			fingerprint := RandomFingerprint()
			assetID, err := GenerateAssetIDv0(
				fingerprint,
				42,
				1,
				AssetLogical,
				AssetStateful,
			)
			require.NoError(t, err)

			assert.Equal(t, TagAssetV0, assetID.Tag())
			assert.Equal(t, uint32(42), assetID.Variant())
			assert.Equal(t, uint16(1), assetID.Standard())
			assert.True(t, assetID.Flag(AssetLogical))
			assert.True(t, assetID.Flag(AssetStateful))

			// Test unsupported flags
			_, err = GenerateAssetIDv0(fingerprint, 42, 1, LogicAuxiliary)
			assert.Equal(t, err, ErrUnsupportedFlag)
		})

		t.Run("Random", func(t *testing.T) {
			assetID := RandomAssetIDv0()

			assert.NoError(t, assetID.Validate())
			assert.Equal(t, TagAssetV0, assetID.Tag())
		})
	})
}
