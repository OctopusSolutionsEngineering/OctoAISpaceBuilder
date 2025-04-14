package model

import "time"

type TerraformPlan struct {
	ID               string    `jsonapi:"primary,terraformplan" json:"id"`
	PlanBinaryBase64 *string   `jsonapi:"attr,plan_binary" json:"plan"`
	PlanText         *string   `jsonapi:"attr,plan_text" json:"plan_text"`
	Server           string    `jsonapi:"attr,server" json:"server"`
	Created          time.Time `jsonapi:"attr,created" json:"created"`
}
