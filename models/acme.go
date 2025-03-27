package models

type AcmeForm struct {
	Email	string	`json:"email"`
	Domains	[]string	`json:"domains"`
	Config	map[string]string	`json:"config"`
	Auto	bool	`json:"auto"`
	Result	string	`json:"result"`
}

