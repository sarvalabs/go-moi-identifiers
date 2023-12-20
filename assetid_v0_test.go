package identifiers

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/google/gofuzz"
	"github.com/stretchr/testify/require"
)

func TestAssetIDv0(t *testing.T) {
	// Setup fuzzer
	f := fuzz.New().NilChance(0)

	type params struct {
		Logical   bool
		Stateful  bool
		Dimension uint8
		Standard  uint16
		Address   Address
	}

	var x params

	for i := 0; i < 1000; i++ {
		// Fuzz Parameters
		f.Fuzz(&x)

		// Create new AssetID
		assetID := NewAssetIDv0(x.Logical, x.Stateful, x.Dimension, x.Standard, x.Address)
		require.Equal(t, x.Address, assetID.Address())

		// Check type conversions
		require.Equal(t, "0x"+string(assetID), assetID.String())
		require.Equal(t, must(decodeHexString(string(assetID))), assetID.Bytes())

		// Check identifier conversion
		identifier, err := assetID.Identifier()
		require.NoError(t, err)

		// Check encoded parameter accessors
		require.Equal(t, 0, identifier.Version())
		require.Equal(t, assetID, identifier.AssetID())

		require.Equal(t, x.Logical, identifier.IsLogical())
		require.Equal(t, x.Stateful, identifier.IsStateful())

		require.Equal(t, x.Dimension, identifier.Dimension())
		require.Equal(t, uint64(x.Standard), identifier.Standard())
		require.Equal(t, x.Address, identifier.Address())

		// Check serialization
		encoded, err := json.Marshal(assetID)
		require.NoError(t, err)
		require.Equal(t, fmt.Sprintf(`"0x%v"`, string(assetID)), string(encoded))

		// Check deserialization
		var decoded AssetID
		err = json.Unmarshal(encoded, &decoded)
		require.NoError(t, err)
		require.Equal(t, assetID, decoded)
	}
}
