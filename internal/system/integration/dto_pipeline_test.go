package integration_test

import (
	"testing"

	"github.com/grafvonb/camunder/internal/system/adapters/pubconv"
	"github.com/grafvonb/camunder/internal/system/adapters/wireconv"
	d "github.com/grafvonb/camunder/internal/system/domain"
	"github.com/grafvonb/camunder/internal/system/wire"
	"github.com/stretchr/testify/require"
)

const payload = `{
  "TotalCount": 2,
  "Entities": [
    {
      "Keys": ["476d87d3-6b8d-4fb8-e8d4-4b7343717c76"],
      "Columns": {
        "UID_System": { "Value": "809601fb-05e1-49a4-a64b-b92a7418d8eb", "IsReadOnly": true },
        "DisplayName": { "Value": "Atruvia", "IsReadOnly": true },
        "CI": { "Value": "000000001", "IsReadOnly": true },
        "LeanixID": { "Value": "NULL", "IsReadOnly": true },
        "EnvironmentDisplay": { "Value": "Produktionsumgebung", "IsReadOnly": true },
        "EnvironmentType": { "Value": 0, "IsReadOnly": true, "DisplayValue": "0" }
      }
    },
    {
      "Keys": ["ffb00c62-e7b5-afc5-03ee-f42d09aa1f52"],
      "Columns": {
        "UID_System": { "Value": "dc0fc84f-5a74-4eb0-80a8-6cf777f2c929", "IsReadOnly": true },
        "DisplayName": { "Value": "AADUMGPrefixTest001", "IsReadOnly": true },
        "CI": { "Value": "9987654233665", "IsReadOnly": true },
        "LeanixID": { "Value": "NULL", "IsReadOnly": true },
        "EnvironmentDisplay": { "Value": "Entwicklungsumgebung", "IsReadOnly": true },
        "EnvironmentType": { "Value": 2, "IsReadOnly": true, "DisplayValue": "2" }
      }
    }
  ]
}`

func Test_FromWireOverDomainToPublic(t *testing.T) {
	ws, err := wire.DecodeEntities([]byte(payload))
	require.NoError(t, err)
	require.Len(t, ws, 2)

	s0, _ := wireconv.FromWire(ws[0])
	require.Equal(t, "809601fb-05e1-49a4-a64b-b92a7418d8eb", s0.ID)
	require.Equal(t, "Atruvia", s0.DisplayName)
	s1, _ := wireconv.FromWire(ws[1])
	require.Equal(t, "dc0fc84f-5a74-4eb0-80a8-6cf777f2c929", s1.ID)
	require.Equal(t, "AADUMGPrefixTest001", s1.DisplayName)

	ds := []d.System{s0, s1}
	for i := range ds {
		require.NoError(t, ds[i].Validate(), "idx %d", i)
	}

	ps, _ := pubconv.ToPublicSlice(ds)
	require.Len(t, ps, 2)

	require.Equal(t, "809601fb-05e1-49a4-a64b-b92a7418d8eb", ps[0].ID)
	require.Equal(t, "Atruvia", ps[0].DisplayName)
	require.Equal(t, "000000001", ps[0].CI)
	require.Equal(t, "", ps[0].LeanixID)
	require.Equal(t, "Produktionsumgebung", ps[0].EnvironmentDisplay)
	require.Equal(t, 0, ps[0].EnvironmentType)

	require.Equal(t, "dc0fc84f-5a74-4eb0-80a8-6cf777f2c929", ps[1].ID)
	require.Equal(t, "AADUMGPrefixTest001", ps[1].DisplayName)
	require.Equal(t, "9987654233665", ps[1].CI)
	require.Equal(t, "", ps[1].LeanixID)
	require.Equal(t, "Entwicklungsumgebung", ps[1].EnvironmentDisplay)
	require.Equal(t, 2, ps[1].EnvironmentType)
}
