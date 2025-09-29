package domain

import "errors"

type System struct {
	ID                 string
	DisplayName        string
	CI                 string
	LeanixID           string
	EnvironmentDisplay string
	EnvironmentType    int
	Modulfilter        string
	Status             string
	EntraAD            bool
	ADS                bool
	EntraPrefix        string
	ADSPrefix          string
	Interface          string
	AppType            string
	APIPermissions     string
}

func (s System) Validate() error {
	var errs []error
	if s.ID == "" {
		errs = append(errs, errors.New("id is required"))
	}
	if s.DisplayName == "" {
		errs = append(errs, errors.New("display name is required"))
	}
	return errors.Join(errs...)
}
