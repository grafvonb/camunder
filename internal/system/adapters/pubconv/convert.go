package pubconv

import (
	d "github.com/grafvonb/camunder/internal/system/domain"
	p "github.com/grafvonb/camunder/pkg/system/publicv1"
)

func ToPublic(s d.System) (p.System, error) {
	return p.System{
		ID:                 s.ID,
		DisplayName:        s.DisplayName,
		CI:                 s.CI,
		LeanixID:           s.LeanixID,
		EnvironmentDisplay: s.EnvironmentDisplay,
		EnvironmentType:    s.EnvironmentType,
		Modulfilter:        s.Modulfilter,
		Status:             s.Status,
		EntraAD:            s.EntraAD,
		ADS:                s.ADS,
		EntraPrefix:        s.EntraPrefix,
		ADSPrefix:          s.ADSPrefix,
		Interface:          s.Interface,
		AppType:            s.AppType,
		APIPermissions:     s.APIPermissions,
	}, nil
}

func ToPublicSlice(in []d.System) ([]p.System, error) {
	var errs []error
	if in == nil {
		return nil, nil
	}

	out := make([]p.System, len(in))
	var err error
	for i := range in {
		out[i], err = ToPublic(in[i])
		if err != nil {
			errs = append(errs, err)
		}
	}
	return out, nil
}
