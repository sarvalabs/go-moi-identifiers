package identifiers

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMust(t *testing.T) {
	sample := func(a bool) (struct{}, error) {
		if a {
			return struct{}{}, errors.New("error")
		}

		return struct{}{}, nil
	}

	require.NotPanics(t, func() {
		must(sample(false))
	})

	require.Panics(t, func() {
		must(sample(true))
	})
}

func TestGetFlag(t *testing.T) {
	tests := []struct {
		value byte
		index uint8
		want  bool
	}{
		{0b00000001, 0, true},
		{0b00000000, 0, false},
		{0b10000000, 7, true},
		{0b00000000, 7, false},
		{0b10101010, 1, true},
		{0b10101010, 2, false},
	}

	for _, tt := range tests {
		assert.Equal(t, tt.want, getFlag(tt.value, tt.index))
	}
}

func TestSetFlag(t *testing.T) {
	tests := []struct {
		value byte
		index uint8
		set   bool
		want  byte
	}{
		{0b00000000, 0, true, 0b00000001},
		{0b00000001, 0, false, 0b00000000},
		{0b00000000, 7, true, 0b10000000},
		{0b10000000, 7, false, 0b00000000},
		{0b10101010, 1, true, 0b10101010},
		{0b10101010, 2, true, 0b10101110},
		{0b10101010, 1, false, 0b10101000},
	}

	for _, tt := range tests {
		assert.Equal(t, tt.want, setFlag(tt.value, tt.index, tt.set))
	}
}

func TestMakeFlag(t *testing.T) {
	tests := []struct {
		kind    IdentifierKind
		index   uint8
		version uint8
		want    Flag
	}{
		{
			KindParticipant, 0, 0,
			Flag{index: 0, support: map[IdentifierKind]uint8{KindParticipant: 0}},
		},
		{
			KindAsset, 1, 1,
			Flag{index: 1, support: map[IdentifierKind]uint8{KindAsset: 1}},
		},
		{
			KindLogic, 10, 1,
			Flag{index: 1, support: map[IdentifierKind]uint8{KindAsset: 1}},
		},
		{
			KindLogic, 1, 20,
			Flag{index: 1, support: map[IdentifierKind]uint8{KindAsset: 1}},
		},
	}

	for _, tt := range tests {
		if tt.index > 7 || tt.version > 15 {
			require.Panics(t, func() {
				makeFlag(tt.kind, tt.index, tt.version)
			})
		} else {
			assert.Equal(t, tt.want, makeFlag(tt.kind, tt.index, tt.version))
		}
	}
}
