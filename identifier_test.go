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
			tag:             IdentifierTag((0x0F << 4) | 0),
			expectedKind:    IdentifierKind(0x0F),
			expectedVersion: 0,
			expectedValid:   false,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			// Test Kind extraction
			assert.Equal(t, tc.expectedKind, tc.tag.Kind(),
				"Incorrect kind extraction")

			// Test Version extraction
			assert.Equal(t, tc.expectedVersion, tc.tag.Version(),
				"Incorrect version extraction")

			// Test Validation
			err := tc.tag.Validate()
			if tc.expectedValid {
				assert.NoError(t, err, "Expected valid tag")
			} else {
				assert.Error(t, err, "Expected invalid tag")
			}
		})
	}
}

func TestIdentifier(t *testing.T) {
	// Create a sample identifier
	id := Identifier{
		// Tag (first byte)
		byte(TagParticipantV0),
		// Flags (next 1 byte)
		0b00000001,
		// Auxiliary (next 2 bytes)
		0x02, 0x03,
		// Account ID (next 24 bytes)
		0x10, 0x11, 0x12, 0x13, 0x14, 0x15, 0x16, 0x17, 0x18, 0x19, 0x1A, 0x1B,
		0x20, 0x21, 0x22, 0x23, 0x24, 0x25, 0x26, 0x27, 0x28, 0x29, 0x2A, 0x2B,
		// Variant (last 4 bytes)
		0x30, 0x31, 0x32, 0x33,
	}

	// Test Tag extraction
	assert.Equal(t, TagParticipantV0, id.Tag(), "Incorrect tag extraction")

	// Test Metadata extraction
	expectedMetadata := [4]byte{byte(TagParticipantV0), 0x01, 0x02, 0x03}
	assert.Equal(t, expectedMetadata, id.Metadata(), "Incorrect metadata extraction")

	// Test AccountID extraction
	expectedAccountID := [24]byte{
		0x10, 0x11, 0x12, 0x13, 0x14, 0x15, 0x16, 0x17, 0x18, 0x19, 0x1A, 0x1B,
		0x20, 0x21, 0x22, 0x23, 0x24, 0x25, 0x26, 0x27, 0x28, 0x29, 0x2A, 0x2B,
	}
	assert.Equal(t, expectedAccountID, id.AccountID(), "Incorrect account ID extraction")

	// Test Variant extraction
	expectedVariant := uint32(0x30313233)
	assert.Equal(t, expectedVariant, id.Variant(), "Incorrect variant extraction")

	// Test IsNil
	assert.False(t, id.IsNil(), "Non-zero identifier should not be nil")
	assert.True(t, Identifier(Nil).IsNil(), "Zero identifier should be nil")

	// Test Bytes method
	assert.Equal(t, id[:], id.Bytes(), "Bytes method should return full identifier")

	// Test String method
	expectedHex := "0x00010203101112131415161718191a1b202122232425262728292a2b30313233"
	assert.Equal(t, expectedHex, id.String(), "Incorrect hex representation")
	assert.Equal(t, expectedHex, id.Hex(), "String and Hex methods should match")
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

// Benchmark different Identifier methods
func BenchmarkIdentifier(b *testing.B) {
	// Create a sample identifier
	id := Identifier{
		// Tag (first byte)
		byte(TagParticipantV0),
		// Flags (next 1 byte)
		0b00000001,
		// Auxiliary (next 2 bytes)
		0x02, 0x03,
		// Account ID (next 24 bytes)
		0x10, 0x11, 0x12, 0x13, 0x14, 0x15, 0x16, 0x17, 0x18, 0x19, 0x1A, 0x1B,
		0x20, 0x21, 0x22, 0x23, 0x24, 0x25, 0x26, 0x27, 0x28, 0x29, 0x2A, 0x2B,
		// Variant (last 4 bytes)
		0x30, 0x31, 0x32, 0x33,
	}

	b.Run("Tag", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_ = id.Tag()
		}
	})

	b.Run("Metadata", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_ = id.Metadata()
		}
	})

	b.Run("AccountID", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_ = id.AccountID()
		}
	})

	b.Run("Variant", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_ = id.Variant()
		}
	})

	b.Run("Hex", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_ = id.Hex()
		}
	})
}
