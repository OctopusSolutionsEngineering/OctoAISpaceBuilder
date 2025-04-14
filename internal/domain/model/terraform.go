package model

type Terraform struct {
	ID            string `jsonapi:"primary,terraform" json:"id"`
	Configuration string `jsonapi:"attr,configuration" json:"configuration"`
}
