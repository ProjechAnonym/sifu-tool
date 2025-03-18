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
	Result  string	`json:"result"`
	Status  int		`json:"status"`
}
type Cloudflare struct {
	API		string `yaml:"api"`
}