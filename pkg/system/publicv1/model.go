package publicv1

// System mirrors domain.System exactly (same fields, same types).
type System struct {
	ID                 string `json:"id"`
	DisplayName        string `json:"displayName"`
	CI                 string `json:"ci"`
	LeanixID           string `json:"leanixId"`
	EnvironmentDisplay string `json:"environmentDisplay"`
	EnvironmentType    int    `json:"environmentType"`
	Modulfilter        string `json:"modulfilter"`
	Status             string `json:"status"`
	EntraAD            bool   `json:"entraAD"`
	ADS                bool   `json:"ads"`
	EntraPrefix        string `json:"entraPrefix"`
	ADSPrefix          string `json:"adsPrefix"`
	Interface          string `json:"interface"`
	AppType            string `json:"appType"`
	APIPermissions     string `json:"apiPermissions"`
}
