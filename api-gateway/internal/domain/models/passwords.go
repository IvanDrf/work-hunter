package models

type Password struct {
	Old string `json:"old"`
	New string `json:"new"`
}

func (p *Password) IsValid() bool {
	return p != nil && p.Old != "" && p.New != ""
}
