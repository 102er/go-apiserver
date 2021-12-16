package validator

import "testing"

func TestValidateData(t *testing.T) {
	type testStruct struct {
		Name      string `json:"name" validator:"required" label:"函数名"`
		StartTime int    `json:"start_time" validator:"required,gt=0" label:"开始时间"`
		EndTime   int    `json:"end_time"`
		Interval  string `json:"interval" validator:"required" label:"时间间隔"`
	}
	d := testStruct{
		Name: "1",
	}
	err := ValidateData("en", d)
	if err != nil {
		return
	}
}

func TestValidateData2(t *testing.T) {
	type testStruct struct {
		Name      string `json:"name" validator:"required" label:"函数名"`
		StartTime int    `json:"start_time" validator:"required,gt=0" label:"开始时间"`
		EndTime   int    `json:"end_time"`
		Interval  string `json:"interval" validator:"required" label:"时间间隔"`
	}
	d := testStruct{
		Name: "1",
	}
	err := ValidateData("en", d)
	if err != nil {
		return
	}
}
