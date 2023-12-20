package identifiers

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/google/gofuzz"
	"github.com/stretchr/testify/require"
)

func TestLogicIDv0(t *testing.T) {
	// Setup fuzzer
	f := fuzz.New().NilChance(0)

	type params struct {
		Persistent  bool
		Ephemeral   bool
		Interactive bool
		AssetLogic  bool

		Edition uint16
		Address Address
	}

	var x params

	for i := 0; i < 1000; i++ {
		// Fuzz Parameters
		f.Fuzz(&x)

		// Create new LogicID
		logicID := NewLogicIDv0(x.Persistent, x.Ephemeral, x.Interactive, x.AssetLogic, x.Edition, x.Address)
		require.Equal(t, x.Address, logicID.Address())

		// Check type conversions
		require.Equal(t, "0x"+string(logicID), logicID.String())
		require.Equal(t, must(decodeHexString(string(logicID))), logicID.Bytes())

		// Check identifier conversion
		identifier, err := logicID.Identifier()
		require.NoError(t, err)

		// Check encoded parameter accessors
		require.Equal(t, 0, identifier.Version())
		require.Equal(t, logicID, identifier.LogicID())

		require.Equal(t, x.Persistent, identifier.HasPersistentState())
		require.Equal(t, x.Ephemeral, identifier.HasEphemeralState())
		require.Equal(t, x.Interactive, identifier.HasInteractableSites())
		require.Equal(t, x.AssetLogic, identifier.AssetLogic())

		require.Equal(t, uint64(x.Edition), identifier.Edition())
		require.Equal(t, x.Address, identifier.Address())

		// Check serialization
		encoded, err := json.Marshal(logicID)
		require.NoError(t, err)
		require.Equal(t, fmt.Sprintf(`"0x%v"`, string(logicID)), string(encoded))

		// Check deserialization
		var decoded LogicID
		err = json.Unmarshal(encoded, &decoded)
		require.NoError(t, err)
		require.Equal(t, logicID, decoded)
	}
}
