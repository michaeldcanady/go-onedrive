package domain

type Contents struct {
	CTag string `json:"ctag"`
	Data []byte `json:"data"`
}
