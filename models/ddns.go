package models

type DDNS struct {
	V4API      []string   `yaml:"v4api"`
	V6API      []string   `yaml:"v6api"`
	Resolver   map[string]map[string]string `yaml:"resolver"`
}
type Domain struct {
	Domain  string	`json:"domain"`
	Type	string	`json:"type"`
	Record	string	`json:"record"`
	Value	string	`json:"value"`
	Result  string	`json:"result,omitempty"`
	Status  int		`json:"status,omitempty"`
	TTL		int		`json:"ttl"`
}
type Cloudflare struct {
	API		string `yaml:"api"`
}

type JobForm struct{
	V4method	int	`json:"v4method,omitempty"`
	V6method	int	`json:"v6method,omitempty"`
	IPV4	string	`json:"ipv4,omitempty"`
	IPV6	string	`json:"ipv6,omitempty"`
	Rev4	string	`json:"rev4,omitempty"`
	Rev6	string	`json:"rev6,omitempty"`
	V4script	string	`json:"v4script,omitempty"`
	V6script	string	`json:"v6script,omitempty"`
	V4interface	string	`json:"v4interface,omitempty"`
	V6interface	string	`json:"v6interface,omitempty"`
	Domains	[]Domain	`json:"domains"`
	Config	map[string]string	`json:"config"`
}