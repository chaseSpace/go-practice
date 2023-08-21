package main

import (
	"math"
	"strings"
	"testing"
)

func TestCheckParam(t *testing.T) {
	tt := []struct {
		ID  int
		val interface{}
		ParamType
		empty                   bool
		errmsgContains          string
		intgreateq, intlesseq   int64
		uintgreateq, uintlesseq uint64
	}{
		{
			ID:             1,
			val:            "",
			ParamType:      ParamTypeDate,
			empty:          true,
			errmsgContains: "",
		},
		{
			ID:             2,
			val:            "",
			ParamType:      ParamTypeDate,
			empty:          false,
			errmsgContains: "cannot be empty",
		},
		{
			ID:             3,
			val:            "x",
			ParamType:      ParamTypeDate,
			empty:          false,
			errmsgContains: "not a Date",
		},
		{
			ID:             4,
			val:            "2023-01-01",
			ParamType:      ParamTypeDate,
			empty:          false,
			errmsgContains: "",
		},
		{
			ID:             5,
			val:            "2023-01-01 01:00",
			ParamType:      ParamTypeDatetime,
			empty:          false,
			errmsgContains: "not a Datetime",
		},
		{
			ID:             6,
			val:            "2023-01-01 01:00:00",
			ParamType:      ParamTypeDatetime,
			empty:          false,
			errmsgContains: "",
		},
		{
			ID:             7,
			val:            nil,
			ParamType:      ParamTypeDatetime,
			empty:          false,
			errmsgContains: "cannot be nil",
		},
		{
			ID:             8,
			val:            int32(-3),
			ParamType:      ParamTypeInt,
			empty:          false,
			errmsgContains: "not in range[-2, 3]",
			intgreateq:     -2,
			intlesseq:      3,
		},
		{
			ID:             9,
			val:            int32(-1),
			ParamType:      ParamTypeInt,
			empty:          false,
			errmsgContains: "",
			intgreateq:     -2,
			intlesseq:      3,
		},
		{
			ID:             10,
			val:            uint32(1),
			ParamType:      ParamTypeUInt,
			empty:          false,
			errmsgContains: "not in range[2, 3]",
			uintgreateq:    2,
			uintlesseq:     3,
		},
		{
			ID:             11,
			val:            uint32(1),
			ParamType:      ParamTypeUInt,
			empty:          false,
			errmsgContains: "",
			uintgreateq:    1,
			uintlesseq:     2,
		},
	}

	for _, s := range tt {
		err := CheckParam(s.val, TypeRule{
			FieldName:   "field",
			AllowEmpty:  s.empty,
			ParamType:   s.ParamType,
			IntGreatEq:  s.intgreateq,
			IntLessEq:   s.intlesseq,
			UIntGreatEq: s.uintgreateq,
			UIntLessEq:  s.uintlesseq,
		})
		if err != nil && (s.errmsgContains == "" || !strings.Contains(err.Error(), s.errmsgContains)) {
			t.Fatal(err, "--ID:", s.ID, s.errmsgContains)
		}
	}
}

func TestParamsCheck(t *testing.T) {
	type xReq struct {
		XDate     string
		XDatetime string
		XInt      int8
		XUInt     uint32
	}

	tt := []struct {
		ID             int
		val            interface{}
		rules          []*TypeRule
		errmsgContains string
	}{
		{
			ID: 1,
			val: xReq{
				XDate:     "",
				XDatetime: "",
			},
			rules: []*TypeRule{
				{FieldName: "XDate", AllowEmpty: false, ParamType: ParamTypeDate},
				{FieldName: "XDatetime", AllowEmpty: false, ParamType: ParamTypeDatetime},
			},
			errmsgContains: "cannot be empty",
		},
		{
			ID: 2,
			val: xReq{
				XDate:     "x",
				XDatetime: "",
			},
			rules: []*TypeRule{
				{FieldName: "XDate", AllowEmpty: false, ParamType: ParamTypeDate},
				{FieldName: "XDatetime", AllowEmpty: true, ParamType: ParamTypeDatetime},
			},
			errmsgContains: "not a Date",
		},
		{
			ID: 3,
			val: xReq{
				XDate:     "",
				XDatetime: "x",
			},
			rules: []*TypeRule{
				{FieldName: "XDate", AllowEmpty: true, ParamType: ParamTypeDate},
				{FieldName: "XDatetime", AllowEmpty: false, ParamType: ParamTypeDatetime},
			},
			errmsgContains: "not a Datetime",
		},
		{
			ID: 4,
			val: xReq{
				XDate:     "2020-01-01",
				XDatetime: "x",
			},
			rules: []*TypeRule{
				{FieldName: "XDate", AllowEmpty: false, ParamType: ParamTypeDate},
				{FieldName: "XDatetime", AllowEmpty: false, ParamType: ParamTypeDatetime},
			},
			errmsgContains: "not a Datetime",
		},
		{
			ID: 5,
			val: xReq{
				XDate:     "2020",
				XDatetime: "2020-01-01 01:01:01",
			},
			rules: []*TypeRule{
				{FieldName: "XDate", AllowEmpty: false, ParamType: ParamTypeDate},
				{FieldName: "XDatetime", AllowEmpty: false, ParamType: ParamTypeDatetime},
			},
			errmsgContains: "not a Date",
		},
		{
			ID: 6,
			val: xReq{
				XDate:     "2020-01-01",
				XDatetime: "2020-01-01 01:01:01",
			},
			rules: []*TypeRule{
				{FieldName: "XDate", AllowEmpty: false, ParamType: ParamTypeDate},
				{FieldName: "XDatetime", AllowEmpty: false, ParamType: ParamTypeDatetime},
			},
			errmsgContains: "",
		},
		{
			ID: 7,
			val: &xReq{
				XDate:     "2020-01-01",
				XDatetime: "2020-01-01 01:01:01",
			},
			rules: []*TypeRule{
				{FieldName: "XDate", AllowEmpty: false, ParamType: ParamTypeDate},
				{FieldName: "XDatetime", AllowEmpty: false, ParamType: ParamTypeDatetime},
			},
			errmsgContains: "",
		},
		{
			ID: 8,
			val: &xReq{
				XInt: int8(-1),
			},
			rules: []*TypeRule{
				{FieldName: "XInt", AllowEmpty: false, ParamType: ParamTypeInt, IntGreatEq: 2, IntLessEq: 3},
			},
			errmsgContains: "not in range[2, 3]",
		},
		{
			ID: 9,
			val: &xReq{
				XInt:  int8(1),
				XUInt: uint32(11),
			},
			rules: []*TypeRule{
				{FieldName: "XUInt", AllowEmpty: false, ParamType: ParamTypeUInt, UIntGreatEq: 2, UIntLessEq: 3},
			},
			errmsgContains: "not in range[2, 3]",
		},
		{
			ID: 10,
			val: &xReq{
				XInt:  int8(1),
				XUInt: uint32(2),
			},
			rules: []*TypeRule{
				{FieldName: "XInt", AllowEmpty: false, ParamType: ParamTypeInt, IntGreatEq: 1, IntLessEq: 3},
				{FieldName: "XUInt", AllowEmpty: false, ParamType: ParamTypeUInt, UIntGreatEq: 2, UIntLessEq: 3},
			},
			errmsgContains: "",
		},
		{
			ID: 11,
			val: &xReq{
				XInt:  int8(0),
				XUInt: uint32(2),
			},
			rules: []*TypeRule{
				{FieldName: "XInt", AllowEmpty: false, ParamType: ParamTypeInt, IntGreatEq: 1, IntLessEq: math.MaxInt64},
				{FieldName: "XUInt", AllowEmpty: false, ParamType: ParamTypeUInt, UIntGreatEq: 2, UIntLessEq: 3},
			},
			errmsgContains: "not in range[1, MaxInt64]",
		},
		{
			ID: 12,
			val: &xReq{
				XInt:  int8(0),
				XUInt: uint32(0),
			},
			rules: []*TypeRule{
				{FieldName: "XInt", AllowEmpty: false, ParamType: ParamTypeInt, IntGreatEq: 1, IntLessEq: math.MaxInt64},
				{FieldName: "XUInt", AllowEmpty: false, ParamType: ParamTypeUInt, UIntGreatEq: 1, UIntLessEq: math.MaxUint64},
			},
			errmsgContains: "not in range[1, MaxUint64]",
		},
	}

	for _, s := range tt {
		err := ParamsCheck(s.val, s.rules...)
		if err != nil && (s.errmsgContains == "" || !strings.Contains(err.Error(), s.errmsgContains)) {
			t.Fatal(err, "--ID:", s.ID, s.errmsgContains)
		}
	}
}
