package models

type TruenasApp struct {
	Id               string            `json:"id"`
	Name             string            `json:"name"`
	State            string            `json:"state"`
	UpgradeAvailable bool              `json:"upgrade_available"`
	LatestVersion    string            `json:"latest_version"`
	Version          string            `json:"version"`
	HumanVersion     string            `json:"human_version"`
	Portals          map[string]string `json:"portals"`
}
