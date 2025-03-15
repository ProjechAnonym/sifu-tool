package models

type DDNS struct {
	V4API      []string   `yaml:"v4api"`
	V6API      []string   `yaml:"v6api"`
	Cloudflare Cloudflare `yaml:"cloudflare"`
}
type Cloudflare struct {
	ZoneAPI string `yaml:"zoneAPI"`
	RecordAPI string `yaml:"recordAPI"`
}