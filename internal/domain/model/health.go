package model

type Health struct {
	ID                                string `jsonapi:"primary,health" json:"id"`
	Status                            string `jsonapi:"attr,status" json:"status"`
	TerraformPlanSuccessLast1Minute   int64  `jsonapi:"attr,terraformPlanSuccessLast1Minute" json:"terraformPlanSuccessLast1Minute"`
	TerraformPlanFailedLast1Minute    int64  `jsonapi:"attr,terraformPlanFailedLast1Minute" json:"terraformPlanFailedLast1Minute"`
	TerraformPlanSuccessLast15Minute  int64  `jsonapi:"attr,terraformPlanSuccessLast15Minute" json:"terraformPlanSuccessLast15Minute"`
	TerraformPlanFailedLast15Minute   int64  `jsonapi:"attr,terraformPlanFailedLast15Minute" json:"terraformPlanFailedLast15Minute"`
	TerraformPlanSuccessLast1Hour     int64  `jsonapi:"attr,terraformPlanSuccessLast1Hour" json:"terraformPlanSuccessLast1Hour"`
	TerraformPlanFailedLast1Hour      int64  `jsonapi:"attr,terraformPlanFailedLast1Hour" json:"terraformPlanFailedLast1Hour"`
	TerraformPlanSuccessAllTime       int64  `jsonapi:"attr,terraformPlanSuccessAllTime" json:"terraformPlanSuccessAllTime"`
	TerraformPlanFailedAllTime        int64  `jsonapi:"attr,terraformPlanFailedAllTime" json:"terraformPlanFailedAllTime"`
	TerraformApplySuccessLast1Minute  int64  `jsonapi:"attr,terraformApplySuccessLast1Minute" json:"terraformApplySuccessLast1Minute"`
	TerraformApplyFailedLast1Minute   int64  `jsonapi:"attr,terraformApplyFailedLast1Minute" json:"terraformApplyFailedLast1Minute"`
	TerraformApplySuccessLast15Minute int64  `jsonapi:"attr,terraformApplySuccessLast15Minute" json:"terraformApplySuccessLast15Minute"`
	TerraformApplyFailedLast15Minute  int64  `jsonapi:"attr,terraformApplyFailedLast15Minute" json:"terraformApplyFailedLast15Minute"`
	TerraformApplySuccessLast1Hour    int64  `jsonapi:"attr,terraformApplySuccessLast1Hour" json:"terraformApplySuccessLast1Hour"`
	TerraformApplyFailedLast1Hour     int64  `jsonapi:"attr,terraformApplyFailedLast1Hour" json:"terraformApplyFailedLast1Hour"`
	TerraformApplySuccessAllTime      int64  `jsonapi:"attr,terraformApplySuccessAllTime" json:"terraformApplySuccessAllTime"`
	TerraformApplyFailedAllTime       int64  `jsonapi:"attr,terraformApplyFailedAllTime" json:"terraformApplyFailedAllTime"`
}
