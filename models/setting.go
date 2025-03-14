package models

type Setting struct {
	User   User   `yaml:"user"`
	Server Server `yaml:"server"`
	DDNS   DDNS   `yaml:"ddns"`
}