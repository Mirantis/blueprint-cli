package types

type Metadata struct {
	Name string `yaml:"name"`
}

type Host struct {
	SSH       SSHHost    `yaml:"ssh"`
	Role      string     `yaml:"role"`
	LocalHost *LocalHost `yaml:"localhost"`
}

type SSHHost struct {
	Address string `yaml:"address"`
	KeyPath string `yaml:"keyPath"`
	Port    int    `yaml:"port"`
	User    string `yaml:"user"`
}

type LocalHost struct {
	Enabled bool `yaml:"enabled"`
}
