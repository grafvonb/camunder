package systemwireconv

import (
	conv "github.com/grafvonb/camunder/internal/api/convert"
	d "github.com/grafvonb/camunder/internal/system/domain"
	w "github.com/grafvonb/camunder/internal/system/wire"
)

func FromWire(wire w.SystemWire) (d.System, error) {
	return d.System{
		ID:                 conv.DerefClean(wire.UIDSystem.Value),
		DisplayName:        conv.DerefClean(wire.DisplayName.Value),
		CI:                 conv.DerefClean(wire.CI.Value),
		LeanixID:           conv.DerefClean(wire.LeanixID.Value),
		EnvironmentDisplay: conv.DerefClean(wire.EnvironmentDisplay.Value),
		EnvironmentType:    conv.Deref(wire.EnvironmentType.Value, 0),
		Modulfilter:        conv.DerefClean(wire.Modulfilter.Value),
		Status:             conv.DerefClean(wire.Status.Value),
		EntraAD:            conv.Deref(wire.EntraAD.Value, false),
		ADS:                conv.Deref(wire.ADS.Value, false),
		EntraPrefix:        conv.DerefClean(wire.EntraPrefix.Value),
		ADSPrefix:          conv.DerefClean(wire.ADSPrefix.Value),
		Interface:          conv.DerefClean(wire.Interface.Value),
		AppType:            conv.DerefClean(wire.AppType.Value),
		APIPermissions:     conv.DerefClean(wire.APIPermissions.Value),
	}, nil
}
