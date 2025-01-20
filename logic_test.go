package identifiers

import (
	"encoding/hex"
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLogicID(t *testing.T) {
	data := [32]byte{
		byte(TagLogicV0), // Tag
		0b00000001,       // Flags
		0x00, 0x10,       // Standard

		// AccountID
		0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08,
		0x11, 0x12, 0x13, 0x14, 0x15, 0x16, 0x17, 0x18,
		0x21, 0x22, 0x23, 0x24, 0x25, 0x26, 0x27, 0x28,

		0x00, 0x00, 0x00, 0x42, // Variant
	}

	// Create a test LogicID
	logicID, err := NewLogicID(data)
	require.NoError(t, err)

	// Test Tag
	assert.Equal(t, TagLogicV0, logicID.Tag())

	// Test AccountID
	assert.Equal(t, [24]byte{
		0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08,
		0x11, 0x12, 0x13, 0x14, 0x15, 0x16, 0x17, 0x18,
		0x21, 0x22, 0x23, 0x24, 0x25, 0x26, 0x27, 0x28,
	}, logicID.AccountID())

	// Test Variant
	assert.Equal(t, uint32(0x42), logicID.Variant())
	// Test IsVariant
	assert.True(t, logicID.IsVariant())

	// Test Flags
	assert.True(t, logicID.Flag(LogicIntrinsic))
	assert.False(t, logicID.Flag(LogicExtrinsic))
	assert.False(t, logicID.Flag(LogicAuxiliary))
	assert.False(t, logicID.Flag(Systemic))
	assert.False(t, logicID.Flag(AssetStateful)) // unsupported flag on set bit

	// Test AsIdentifier
	identifier := Identifier(data[:])
	assert.Equal(t, identifier, logicID.AsIdentifier())

	// Test From Identifier
	converted, err := identifier.AsLogicID()
	require.NoError(t, err)
	require.Equal(t, logicID, converted)

	// Test Bytes
	assert.Equal(t, data[:], logicID.Bytes())

	// Test String & Hex
	expectedHex := "0x2001001001020304050607081112131415161718212223242526272800000042"
	assert.Equal(t, expectedHex, logicID.String())
	assert.Equal(t, expectedHex, logicID.Hex())
}

//nolint:dupl // similar functions
func TestLogicID_Constructor(t *testing.T) {
	t.Run("NewLogicID", func(t *testing.T) {
		t.Run("Valid", func(t *testing.T) {
			logicID, err := NewLogicID([32]byte{
				byte(TagLogicV0), // Tag
				0b00000000,       // Flags
				0x00, 0x01,       // Standard
				// Empty bytes for account and variant
			})

			require.NoError(t, err)
			require.NoError(t, logicID.Validate())
		})

		t.Run("InvalidTag", func(t *testing.T) {
			_, err := NewLogicID([32]byte{0xF0}) // Invalid tag kind
			require.EqualError(t, err, "invalid tag: unsupported tag kind")

			_, err = NewLogicID([32]byte{0x0F}) // Invalid tag version
			require.EqualError(t, err, "invalid tag: unsupported tag version")

			_, err = NewLogicID([32]byte{byte(TagAssetV0)}) // Invalid tag
			require.EqualError(t, err, "invalid tag: not a logic id")
		})

		t.Run("InvalidFlags", func(t *testing.T) {
			_, err := NewLogicID([32]byte{
				byte(TagLogicV0), // Tag
				0b11111111,       // Invalid flags
			})
			require.EqualError(t, err, "invalid flags: unsupported flags for logic id")
		})
	})

	t.Run("NewLogicIDFromBytes", func(t *testing.T) {
		// Less than 32 bytes
		t.Run("< 32 bytes", func(t *testing.T) {
			_, err := NewLogicIDFromBytes([]byte{byte(TagLogicV0), 0x00, 0x00, 0x01})
			require.EqualError(t, err, "invalid length: logic id must be 32 bytes")
		})

		// Exactly 32 bytes
		t.Run("= 32 bytes", func(t *testing.T) {
			logicID, err := NewLogicIDFromBytes([]byte{
				byte(TagLogicV0), // Tag
				0b00000000,       // Flags
				0x00, 0x01,       // Standard

				// AccountID
				0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08,
				0x11, 0x12, 0x13, 0x14, 0x15, 0x16, 0x17, 0x18,
				0x21, 0x22, 0x23, 0x24, 0x25, 0x26, 0x27, 0x28,

				0x00, 0x00, 0x00, 0x01, // Variant
			})

			require.NoError(t, err)
			require.NoError(t, logicID.Validate())
			require.Equal(t, LogicID{
				byte(TagLogicV0), 0x00, 0x00, 0x01,
				0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08,
				0x11, 0x12, 0x13, 0x14, 0x15, 0x16, 0x17, 0x18,
				0x21, 0x22, 0x23, 0x24, 0x25, 0x26, 0x27, 0x28,
				0x00, 0x00, 0x00, 0x01,
			}, logicID)
		})

		// More than 32 bytes
		t.Run("> 32 bytes", func(t *testing.T) {
			_, err := NewLogicIDFromBytes([]byte{
				byte(TagLogicV0), // Tag
				0b00000000,       // Flags
				0x00, 0x01,       // Standard

				// AccountID
				0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08,
				0x11, 0x12, 0x13, 0x14, 0x15, 0x16, 0x17, 0x18,
				0x21, 0x22, 0x23, 0x24, 0x25, 0x26, 0x27, 0x28,

				0x00, 0x00, 0x00, 0x01, // Variant
				0xFF, 0xFF, 0xFF, 0xFF, // Extra bytes
			})
			require.EqualError(t, err, "invalid length: logic id must be 32 bytes")
		})
	})

	t.Run("NewLogicIDFromHex", func(t *testing.T) {
		t.Run("ValidHex", func(t *testing.T) {
			logicID, err := NewLogicIDFromHex("0x" + hex.EncodeToString([]byte{
				byte(TagLogicV0), // Tag
				0b00000000,       // Flags
				0x00, 0x01,       // Standard

				// AccountID
				0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08,
				0x11, 0x12, 0x13, 0x14, 0x15, 0x16, 0x17, 0x18,
				0x21, 0x22, 0x23, 0x24, 0x25, 0x26, 0x27, 0x28,

				0x00, 0x00, 0x00, 0x01, // Variant
			}))

			require.NoError(t, err)
			require.NoError(t, logicID.Validate())
		})

		t.Run("ValidHexNoPrefix", func(t *testing.T) {
			logicID, err := NewLogicIDFromHex(hex.EncodeToString([]byte{
				byte(TagLogicV0), // Tag
				0b00000000,       // Flags
				0x00, 0x01,       // Standard

				// AccountID
				0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08,
				0x11, 0x12, 0x13, 0x14, 0x15, 0x16, 0x17, 0x18,
				0x21, 0x22, 0x23, 0x24, 0x25, 0x26, 0x27, 0x28,

				0x00, 0x00, 0x00, 0x01, // Variant
			}))

			require.NoError(t, err)
			require.NoError(t, logicID.Validate())
		})

		t.Run("InvalidHex", func(t *testing.T) {
			_, err := NewLogicIDFromHex("invalid-hex")
			require.EqualError(t, err, "encoding/hex: invalid byte: U+0069 'i'")

			_, err = NewLogicIDFromHex("0xf") // odd length
			require.EqualError(t, err, "encoding/hex: odd length hex string")
		})
	})

	t.Run("MustLogicID", func(t *testing.T) {
		t.Run("MustLogicID", func(t *testing.T) {
			assert.Panics(t, func() { _ = MustLogicID([32]byte{0xFF}) })
		})

		t.Run("MustLogicIDFromBytes", func(t *testing.T) {
			assert.Panics(t, func() { _ = MustLogicIDFromBytes([]byte{0xFF}) })
		})

		t.Run("MustLogicIDFromHex", func(t *testing.T) {
			assert.Panics(t, func() { _ = MustLogicIDFromHex("0xFF") })
		})
	})
}

func TestLogicID_TextMarshal(t *testing.T) {
	// Create a test LogicID
	data := [32]byte{
		byte(TagLogicV0), // Tag
		0b00000001,       // Flags
		0x00, 0x10,       // Standard

		// AccountID
		0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08,
		0x11, 0x12, 0x13, 0x14, 0x15, 0x16, 0x17, 0x18,
		0x21, 0x22, 0x23, 0x24, 0x25, 0x26, 0x27, 0x28,

		0x00, 0x00, 0x00, 0x42, // Variant
	}

	logicID, err := NewLogicID(data)
	require.NoError(t, err)

	encoded, err := json.Marshal(logicID)
	require.NoError(t, err)
	require.Equal(t, `"0x2001001001020304050607081112131415161718212223242526272800000042"`, string(encoded))

	t.Run("Unmarshal_Success", func(t *testing.T) {
		var decoded LogicID

		require.NoError(t, json.Unmarshal(encoded, &decoded))
		require.Equal(t, logicID, decoded)
	})

	t.Run("Unmarshal_MissingPrefix", func(t *testing.T) {
		var decoded LogicID

		require.Equal(t, json.Unmarshal([]byte(`"invalid-json"`), &decoded), ErrMissingHexPrefix)
	})

	t.Run("Unmarshal_InvalidLength", func(t *testing.T) {
		var decoded LogicID

		require.Equal(t, json.Unmarshal([]byte(`"0xffabcd"`), &decoded), ErrInvalidLength)
	})

	t.Run("Unmarshal_HexError", func(t *testing.T) {
		var decoded LogicID

		require.EqualError(t,
			json.Unmarshal([]byte(`"0xYY01001001020304050607081112131415161718212223242526272800000042"`), &decoded),
			"encoding/hex: invalid byte: U+0059 'Y'",
		)
	})

	t.Run("Unmarshal_Invalid", func(t *testing.T) {
		var decoded LogicID

		require.EqualError(t,
			json.Unmarshal([]byte(`"0xFF01001001020304050607081112131415161718212223242526272800000042"`), &decoded),
			"invalid tag: unsupported tag kind",
		)
	})
}

func TestLogicID_Generation(t *testing.T) {
	t.Run("v0", func(t *testing.T) {
		t.Run("Generate", func(t *testing.T) {
			account := RandomAccountID()
			logicID, err := GenerateLogicIDv0(
				account,
				42,
				LogicIntrinsic,
				LogicExtrinsic,
			)
			require.NoError(t, err)

			assert.Equal(t, TagLogicV0, logicID.Tag())
			assert.Equal(t, uint32(42), logicID.Variant())
			assert.True(t, logicID.Flag(LogicIntrinsic))
			assert.True(t, logicID.Flag(LogicExtrinsic))
			assert.False(t, logicID.Flag(LogicAuxiliary))

			// Test unsupported flags
			_, err = GenerateLogicIDv0(account, 42, AssetLogical)
			assert.Equal(t, err, ErrUnsupportedFlag)
		})

		t.Run("Random", func(t *testing.T) {
			logicID := RandomLogicIDv0()

			assert.NoError(t, logicID.Validate())
			assert.Equal(t, TagLogicV0, logicID.Tag())
		})
	})
}
