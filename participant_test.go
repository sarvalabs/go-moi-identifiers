package identifiers

import (
	"encoding/hex"
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestParticipantID(t *testing.T) {
	data := [32]byte{
		byte(TagParticipantV0), // Tag
		0b10000000,             // Flags
		0x00, 0x10,             // Standard

		// Fingerprint
		0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08,
		0x11, 0x12, 0x13, 0x14, 0x15, 0x16, 0x17, 0x18,
		0x21, 0x22, 0x23, 0x24, 0x25, 0x26, 0x27, 0x28,

		0x00, 0x00, 0x00, 0x42, // Variant
	}

	// Create a test ParticipantID
	participantID, err := NewParticipantID(data)
	require.NoError(t, err)

	// Test Tag
	assert.Equal(t, TagParticipantV0, participantID.Tag())

	// Test Fingerprint
	assert.Equal(t, [24]byte{
		0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08,
		0x11, 0x12, 0x13, 0x14, 0x15, 0x16, 0x17, 0x18,
		0x21, 0x22, 0x23, 0x24, 0x25, 0x26, 0x27, 0x28,
	}, participantID.Fingerprint())

	// Test Variant
	assert.Equal(t, uint32(0x42), participantID.Variant())
	// Test IsVariant
	assert.True(t, participantID.IsVariant())

	// Test Flags
	assert.True(t, participantID.Flag(Systemic))
	assert.False(t, participantID.Flag(LogicIntrinsic)) // unsupported flag on set bit

	// Test AsIdentifier
	identifier := Identifier(data[:])
	assert.Equal(t, identifier, participantID.AsIdentifier())

	// Test From Identifier
	converted, err := identifier.AsParticipantID()
	require.NoError(t, err)
	require.Equal(t, participantID, converted)

	// Test Bytes
	assert.Equal(t, data[:], participantID.Bytes())

	// Test String & Hex
	expectedHex := "0x0080001001020304050607081112131415161718212223242526272800000042"
	assert.Equal(t, expectedHex, participantID.String())
	assert.Equal(t, expectedHex, participantID.Hex())
}

//nolint:dupl // similar functions
func TestParticipantID_Constructor(t *testing.T) {
	t.Run("NewParticipantID", func(t *testing.T) {
		t.Run("Valid", func(t *testing.T) {
			participantID, err := NewParticipantID([32]byte{
				byte(TagParticipantV0), // Tag
				0b00000000,             // Flags
				0x00, 0x01,             // Standard
				// Empty bytes for fingerprint and variant
			})

			require.NoError(t, err)
			require.NoError(t, participantID.Validate())
		})

		t.Run("InvalidTag", func(t *testing.T) {
			_, err := NewParticipantID([32]byte{0xF0}) // Invalid tag kind
			require.EqualError(t, err, "invalid tag: unsupported tag kind")

			_, err = NewParticipantID([32]byte{0x0F}) // Invalid tag version
			require.EqualError(t, err, "invalid tag: unsupported tag version")

			_, err = NewParticipantID([32]byte{byte(TagLogicV0)}) // Invalid tag
			require.EqualError(t, err, "invalid tag: not a participant id")
		})

		t.Run("InvalidFlags", func(t *testing.T) {
			_, err := NewParticipantID([32]byte{
				byte(TagParticipantV0), // Tag
				0b11111111,             // Invalid flags
			})
			require.EqualError(t, err, "invalid flags: unsupported flags for participant id")
		})
	})

	t.Run("NewParticipantIDFromBytes", func(t *testing.T) {
		// Less than 32 bytes
		t.Run("< 32 bytes", func(t *testing.T) {
			_, err := NewParticipantIDFromBytes([]byte{byte(TagParticipantV0), 0x00, 0x00, 0x01})
			require.EqualError(t, err, "invalid length: participant id must be 32 bytes")
		})

		// Exactly 32 bytes
		t.Run("= 32 bytes", func(t *testing.T) {
			participantID, err := NewParticipantIDFromBytes([]byte{
				byte(TagParticipantV0), // Tag
				0b00000000,             // Flags
				0x00, 0x01,             // Standard

				// Fingerprint
				0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08,
				0x11, 0x12, 0x13, 0x14, 0x15, 0x16, 0x17, 0x18,
				0x21, 0x22, 0x23, 0x24, 0x25, 0x26, 0x27, 0x28,

				0x00, 0x00, 0x00, 0x01, // Variant
			})

			require.NoError(t, err)
			require.NoError(t, participantID.Validate())
			require.Equal(t, ParticipantID{
				byte(TagParticipantV0), 0x00, 0x00, 0x01,
				0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08,
				0x11, 0x12, 0x13, 0x14, 0x15, 0x16, 0x17, 0x18,
				0x21, 0x22, 0x23, 0x24, 0x25, 0x26, 0x27, 0x28,
				0x00, 0x00, 0x00, 0x01,
			}, participantID)
		})

		// More than 32 bytes
		t.Run("> 32 bytes", func(t *testing.T) {
			_, err := NewParticipantIDFromBytes([]byte{
				byte(TagParticipantV0), // Tag
				0b00000000,             // Flags
				0x00, 0x01,             // Standard

				// Fingerprint
				0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08,
				0x11, 0x12, 0x13, 0x14, 0x15, 0x16, 0x17, 0x18,
				0x21, 0x22, 0x23, 0x24, 0x25, 0x26, 0x27, 0x28,

				0x00, 0x00, 0x00, 0x01, // Variant
				0xFF, 0xFF, 0xFF, 0xFF, // Extra bytes
			})
			require.EqualError(t, err, "invalid length: participant id must be 32 bytes")
		})
	})

	t.Run("NewParticipantIDFromHex", func(t *testing.T) {
		t.Run("ValidHex", func(t *testing.T) {
			participantID, err := NewParticipantIDFromHex("0x" + hex.EncodeToString([]byte{
				byte(TagParticipantV0), // Tag
				0b00000000,             // Flags
				0x00, 0x01,             // Standard

				// Fingerprint
				0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08,
				0x11, 0x12, 0x13, 0x14, 0x15, 0x16, 0x17, 0x18,
				0x21, 0x22, 0x23, 0x24, 0x25, 0x26, 0x27, 0x28,

				0x00, 0x00, 0x00, 0x01, // Variant
			}))

			require.NoError(t, err)
			require.NoError(t, participantID.Validate())
		})

		t.Run("ValidHexNoPrefix", func(t *testing.T) {
			participantID, err := NewParticipantIDFromHex(hex.EncodeToString([]byte{
				byte(TagParticipantV0), // Tag
				0b00000000,             // Flags
				0x00, 0x01,             // Standard

				// Fingerprint
				0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08,
				0x11, 0x12, 0x13, 0x14, 0x15, 0x16, 0x17, 0x18,
				0x21, 0x22, 0x23, 0x24, 0x25, 0x26, 0x27, 0x28,

				0x00, 0x00, 0x00, 0x01, // Variant
			}))

			require.NoError(t, err)
			require.NoError(t, participantID.Validate())
		})

		t.Run("InvalidHex", func(t *testing.T) {
			_, err := NewParticipantIDFromHex("invalid-hex")
			require.EqualError(t, err, "encoding/hex: invalid byte: U+0069 'i'")

			_, err = NewParticipantIDFromHex("0xf") // odd length
			require.EqualError(t, err, "encoding/hex: odd length hex string")
		})
	})

	t.Run("MustParticipantID", func(t *testing.T) {
		t.Run("MustParticipantID", func(t *testing.T) {
			assert.Panics(t, func() { _ = MustParticipantID([32]byte{0xFF}) })
		})

		t.Run("MustParticipantIDFromBytes", func(t *testing.T) {
			assert.Panics(t, func() { _ = MustParticipantIDFromBytes([]byte{0xFF}) })
		})

		t.Run("MustParticipantIDFromHex", func(t *testing.T) {
			assert.Panics(t, func() { _ = MustParticipantIDFromHex("0xFF") })
		})
	})
}

func TestParticipantID_TextMarshal(t *testing.T) {
	// Create a test ParticipantID
	data := [32]byte{
		byte(TagParticipantV0), // Tag
		0b10000000,             // Flags
		0x00, 0x10,             // Standard

		// Fingerprint
		0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08,
		0x11, 0x12, 0x13, 0x14, 0x15, 0x16, 0x17, 0x18,
		0x21, 0x22, 0x23, 0x24, 0x25, 0x26, 0x27, 0x28,

		0x00, 0x00, 0x00, 0x42, // Variant
	}

	participantID, err := NewParticipantID(data)
	require.NoError(t, err)

	encoded, err := json.Marshal(participantID)
	require.NoError(t, err)
	require.Equal(t, `"0x0080001001020304050607081112131415161718212223242526272800000042"`, string(encoded))

	t.Run("Unmarshal_Success", func(t *testing.T) {
		var decoded ParticipantID

		require.NoError(t, json.Unmarshal(encoded, &decoded))
		require.Equal(t, participantID, decoded)
	})

	t.Run("Unmarshal_MissingPrefix", func(t *testing.T) {
		var decoded ParticipantID

		require.Equal(t, json.Unmarshal([]byte(`"invalid-json"`), &decoded), ErrMissingHexPrefix)
	})

	t.Run("Unmarshal_InvalidLength", func(t *testing.T) {
		var decoded ParticipantID

		require.Equal(t, json.Unmarshal([]byte(`"0xffabcd"`), &decoded), ErrInvalidLength)
	})

	t.Run("Unmarshal_HexError", func(t *testing.T) {
		var decoded ParticipantID

		require.EqualError(t,
			json.Unmarshal([]byte(`"0xYY01001001020304050607081112131415161718212223242526272800000042"`), &decoded),
			"encoding/hex: invalid byte: U+0059 'Y'",
		)
	})

	t.Run("Unmarshal_Invalid", func(t *testing.T) {
		var decoded ParticipantID

		require.EqualError(t,
			json.Unmarshal([]byte(`"0xFF01001001020304050607081112131415161718212223242526272800000042"`), &decoded),
			"invalid tag: unsupported tag kind",
		)
	})
}

func TestParticipantID_Generation(t *testing.T) {
	t.Run("v0", func(t *testing.T) {
		t.Run("Generate", func(t *testing.T) {
			fingerprint := RandomFingerprint()
			participantID, err := GenerateParticipantIDv0(
				fingerprint,
				42,
				Systemic,
			)
			require.NoError(t, err)

			assert.Equal(t, TagParticipantV0, participantID.Tag())
			assert.Equal(t, uint32(42), participantID.Variant())
			assert.True(t, participantID.Flag(Systemic))

			// Test unsupported flags
			_, err = GenerateParticipantIDv0(fingerprint, 42, LogicAuxiliary)
			assert.Equal(t, err, ErrUnsupportedFlag)
		})

		t.Run("Random", func(t *testing.T) {
			participantID := RandomParticipantIDv0()

			assert.NoError(t, participantID.Validate())
			assert.Equal(t, TagParticipantV0, participantID.Tag())
		})
	})
}
