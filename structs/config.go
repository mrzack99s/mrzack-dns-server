package structs

type SystemConfig struct {
	SConfig struct {
		ForwarderAddress string `yaml:"forwarderAddress"`
		LogPath          string `yaml:"logPath"`
	} `yaml:"sconfig"`
}
