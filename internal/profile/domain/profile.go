package domain

type Profile struct {
	Name              string `json:"name"`
	Path              string `json:"path"`
	ConfigurationPath string `json:"configurationPath"`
}
