package model

import "time"

type TerraformPlan struct {
	ID      string    `jsonapi:"primary,terraform" json:"id"`
	Plan    string    `jsonapi:"attr,plan" json:"plan"`
	Server  string    `jsonapi:"attr,server" json:"server"`
	Created time.Time `jsonapi:"attr,created" json:"created"`
}
