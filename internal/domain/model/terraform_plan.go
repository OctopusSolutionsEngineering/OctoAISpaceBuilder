package model

import "time"

type TerraformPlan struct {
	ID               string    `jsonapi:"primary,terraformplan" json:"id"`
	PlanBinaryBase64 *string   `jsonapi:"attr,plan_binary" json:"plan"`
	PlanText         *string   `jsonapi:"attr,plan_text" json:"plan_text"`
	Server           string    `jsonapi:"attr,server" json:"server"`
	Created          time.Time `jsonapi:"attr,created" json:"created"`
	SpaceId          string    `jsonapi:"attr,space_id" json:"space_id"`
	Configuration    string    `jsonapi:"attr,configuration" json:"configuration"`
}
