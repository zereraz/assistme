package policy

import "github.com/teris-io/shortid"

var DefaultPolicy *Policy

func init() {
	var err error
	DefaultPolicy, err = NewPolicy("Default")
	if err != nil {
		panic("Could not create Default policy")
	}
}

type Policy struct {
	Id   string `json:"id"`
	Name string `json:"name"`
}

func NewPolicy(name string) (*Policy, error) {
	policyId, err := shortid.Generate()
	if err != nil {
		return nil, err
	}
	return &Policy{policyId, name}, nil
}
