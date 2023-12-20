package identifiers

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNewAddressFromBytes(t *testing.T) {
	tests := []struct {
		name string
		data []byte
		addr Address
	}{
		{
			name: "Exactly 32 Bytes",
			data: []byte{
				222, 236, 155, 94, 169, 166, 126, 223, 253, 36, 250, 127, 164, 129, 129,
				6, 247, 23, 60, 60, 3, 98, 11, 214, 171, 84, 241, 232, 31, 251, 112, 152,
			},
			addr: Address{
				222, 236, 155, 94, 169, 166, 126, 223, 253, 36, 250, 127, 164, 129, 129,
				6, 247, 23, 60, 60, 3, 98, 11, 214, 171, 84, 241, 232, 31, 251, 112, 152,
			},
		},
		{
			name: "Less Than 32 Bytes",
			data: []byte{
				222, 236, 155, 94, 169, 166, 126, 223, 253, 36, 250, 127, 164, 129, 129,
				6, 247, 23, 60, 60, 3, 98, 11, 214, 171,
			},
			addr: Address{
				0, 0, 0, 0, 0, 0, 0, 222, 236, 155, 94, 169, 166, 126, 223, 253, 36,
				250, 127, 164, 129, 129, 6, 247, 23, 60, 60, 3, 98, 11, 214, 171,
			},
		},
		{
			name: "Greater Than 32 Bytes",
			data: []byte{
				97, 56, 45, 32, 222, 236, 155, 94, 169, 166, 126, 223, 253, 36, 250, 127, 164,
				129, 129, 6, 247, 23, 60, 60, 3, 98, 11, 214, 171, 84, 241, 232, 31, 251, 112, 152,
			},
			addr: Address{
				222, 236, 155, 94, 169, 166, 126, 223, 253, 36, 250, 127, 164, 129, 129,
				6, 247, 23, 60, 60, 3, 98, 11, 214, 171, 84, 241, 232, 31, 251, 112, 152,
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			require.Equal(t, test.addr, NewAddressFromBytes(test.data))
		})
	}
}

func TestAddressString(t *testing.T) {
	tests := []struct {
		addr Address
		out  string
	}{
		{
			Address{
				222, 236, 155, 94, 169, 166, 126, 223, 253, 36, 250, 127, 164, 129, 129,
				6, 247, 23, 60, 60, 3, 98, 11, 214, 171, 84, 241, 232, 31, 251, 112, 152,
			},
			"0xdeec9b5ea9a67edffd24fa7fa4818106f7173c3c03620bd6ab54f1e81ffb7098",
		},
		{
			NilAddress,
			"0x0000000000000000000000000000000000000000000000000000000000000000",
		},
	}

	for _, test := range tests {
		require.Equal(t, test.out, test.addr.String())
	}
}
