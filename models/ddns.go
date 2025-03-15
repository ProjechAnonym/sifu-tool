package models

type DDNS struct {
	V4API      []string   `yaml:"v4api"`
	V6API      []string   `yaml:"v6api"`
	Cloudflare Cloudflare `yaml:"cloudflare"`
}
type Domain struct {
	Domain  string	`json:"domain"`
	Type	string	`json:"type"`
	Value	string	`json:"value"`
}
type Cloudflare struct {
	ZoneAPI		string `yaml:"zoneAPI"`
	RecordAPI	string `yaml:"recordAPI"`
}