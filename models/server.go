package models

type Tls struct {
	Key  string `yaml:"key"`
	Cert string `yaml:"cert"`
	Port int    `yaml:"port"`
}
type Server struct {
	Tls      *Tls   `yaml:"tls,omitempty"`
	Entrance string `yaml:"entrance,omitempty"`
}