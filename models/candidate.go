package models

type Candidate struct {
	Education string   `json:"education"`
	Skills    []string `json:"skills"`
	Interests []string `json:"interests"`
	Location  string   `json:"location"`
	Phone     string   `json:"phone"`
}
