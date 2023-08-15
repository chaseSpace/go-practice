package main

import (
	"strings"
	"testing"
)

func TestCheckParam(t *testing.T) {
	tt := []struct {
		val interface{}
		ParamType
		empty          bool
		errmsgContains string
	}{
		{
			val:            "",
			ParamType:      ParamTypeDate,
			empty:          true,
			errmsgContains: "",
		},
		{
			val:            "",
			ParamType:      ParamTypeDate,
			empty:          false,
			errmsgContains: "cannot be empty",
		},
		{
			val:            "x",
			ParamType:      ParamTypeDate,
			empty:          false,
			errmsgContains: "not a Date",
		},
		{
			val:            "2023-01-01",
			ParamType:      ParamTypeDate,
			empty:          false,
			errmsgContains: "",
		},
		{
			val:            "2023-01-01 01:00",
			ParamType:      ParamTypeDateTime,
			empty:          false,
			errmsgContains: "not a DateTime",
		},
		{
			val:            "2023-01-01 01:00:00",
			ParamType:      ParamTypeDateTime,
			empty:          false,
			errmsgContains: "",
		},
		{
			val:            nil,
			ParamType:      ParamTypeDateTime,
			empty:          false,
			errmsgContains: "not a kind:",
		},
	}

	for _, s := range tt {
		err := CheckParam(s.ParamType, "field", s.val, s.empty)
		if err != nil && !strings.Contains(err.Error(), s.errmsgContains) {
			t.Fatal(err, "--", s)
		}
	}
}

func TestParamsCheck(t *testing.T) {
	type xReq struct {
		XDate     string
		XDateTime string
	}

	tt := []struct {
		val            interface{}
		rules          []*TypeRule
		errmsgContains string
	}{
		{
			val: xReq{
				XDate:     "",
				XDateTime: "",
			},
			rules: []*TypeRule{
				{FieldName: "XDate", AllowEmpty: false, ParamType: ParamTypeDate},
				{FieldName: "XDateTime", AllowEmpty: false, ParamType: ParamTypeDateTime},
			},
			errmsgContains: "cannot be empty",
		},
		{
			val: xReq{
				XDate:     "x",
				XDateTime: "",
			},
			rules: []*TypeRule{
				{FieldName: "XDate", AllowEmpty: false, ParamType: ParamTypeDate},
				{FieldName: "XDateTime", AllowEmpty: true, ParamType: ParamTypeDateTime},
			},
			errmsgContains: "not a Date",
		},
		{
			val: xReq{
				XDate:     "",
				XDateTime: "x",
			},
			rules: []*TypeRule{
				{FieldName: "XDate", AllowEmpty: true, ParamType: ParamTypeDate},
				{FieldName: "XDateTime", AllowEmpty: false, ParamType: ParamTypeDateTime},
			},
			errmsgContains: "not a Datetime",
		},
		{
			val: xReq{
				XDate:     "2020-01-01",
				XDateTime: "x",
			},
			rules: []*TypeRule{
				{FieldName: "XDate", AllowEmpty: false, ParamType: ParamTypeDate},
				{FieldName: "XDateTime", AllowEmpty: false, ParamType: ParamTypeDateTime},
			},
			errmsgContains: "not a Datetime",
		},
		{
			val: xReq{
				XDate:     "2020",
				XDateTime: "2020-01-01 01:01:01",
			},
			rules: []*TypeRule{
				{FieldName: "XDate", AllowEmpty: false, ParamType: ParamTypeDate},
				{FieldName: "XDateTime", AllowEmpty: false, ParamType: ParamTypeDateTime},
			},
			errmsgContains: "not a Date",
		},
		{
			val: xReq{
				XDate:     "2020-01-01",
				XDateTime: "2020-01-01 01:01:01",
			},
			rules: []*TypeRule{
				{FieldName: "XDate", AllowEmpty: false, ParamType: ParamTypeDate},
				{FieldName: "XDateTime", AllowEmpty: false, ParamType: ParamTypeDateTime},
			},
			errmsgContains: "",
		},
		{
			val: &xReq{
				XDate:     "2020-01-01",
				XDateTime: "2020-01-01 01:01:01",
			},
			rules: []*TypeRule{
				{FieldName: "XDate", AllowEmpty: false, ParamType: ParamTypeDate},
				{FieldName: "XDateTime", AllowEmpty: false, ParamType: ParamTypeDateTime},
			},
			errmsgContains: "",
		},
	}

	for i, s := range tt {
		err := ParamsCheck(s.val, s.rules...)
		if err != nil && !strings.Contains(err.Error(), s.errmsgContains) {
			t.Fatal(err, "--", i, s.errmsgContains)
		}
	}
}
