package model

type TerraformApply struct {
	ID     string `jsonapi:"primary,terraformapply" json:"id"`
	PlanId string `jsonapi:"attr,plan_id" json:"plan_id"`
	Server string `jsonapi:"attr,server" json:"server"`
}
