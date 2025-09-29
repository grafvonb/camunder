package wire

import "encoding/json"

type SystemWire struct {
	UIDSystem          ColStr  `json:"UID_System"`
	DisplayName        ColStr  `json:"DisplayName"`
	CI                 ColStr  `json:"CI"`
	LeanixID           ColStr  `json:"LeanixID"`
	EnvironmentDisplay ColStr  `json:"EnvironmentDisplay"`
	EnvironmentType    ColInt  `json:"EnvironmentType"`
	Modulfilter        ColStr  `json:"Modulfilter"`
	Status             ColStr  `json:"Status"`
	EntraAD            ColBool `json:"EntraAD"`
	ADS                ColBool `json:"ADS"`
	EntraPrefix        ColStr  `json:"EntraPrefix"`
	ADSPrefix          ColStr  `json:"ADSPrefix"`
	Interface          ColStr  `json:"Interface"`
	AppType            ColStr  `json:"AppType"`
	APIPermissions     ColStr  `json:"APIPermissions"`
}

type envelope struct {
	Entities []struct {
		Columns SystemWire `json:"Columns"`
	} `json:"Entities"`
}

func DecodeEntities(b []byte) ([]SystemWire, error) {
	var e envelope
	if err := json.Unmarshal(b, &e); err != nil {
		return nil, err
	}
	system := make([]SystemWire, 0, len(e.Entities))
	for _, entity := range e.Entities {
		system = append(system, entity.Columns)
	}
	return system, nil
}

type ColStr struct {
	Value *string `json:"Value"`
}
type ColInt struct {
	Value *int `json:"Value"`
}
type ColBool struct {
	Value *bool `json:"Value"`
}
