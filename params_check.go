package main

import (
	"errors"
	"fmt"
	"math"
	"reflect"
	"strings"
	"time"
)

type ParamType int8

type _Param struct {
	Typ         ParamType
	Desc        string
	ReflectKind reflect.Kind
}

const (
	ParamTypeDate ParamType = iota
	ParamTypeDatetime
	ParamTypeInt  // 所有int类型共用
	ParamTypeUInt // 所有uint类型共用
)

var (
	_ParamDate     = _Param{Typ: ParamTypeDate, Desc: "Date => {YYYY-MM-DD}", ReflectKind: reflect.String}
	_ParamDatetime = _Param{Typ: ParamTypeDatetime, Desc: "Datetime => {YYYY-mm-DD HH:MM:SS}", ReflectKind: reflect.String}
	_ParamInt      = _Param{Typ: ParamTypeInt, Desc: "Integer", ReflectKind: reflect.Int64}
	_ParamUInt     = _Param{Typ: ParamTypeUInt, Desc: "Unsigned Integer", ReflectKind: reflect.Uint64}
)

func getParam(paramType ParamType) _Param {
	switch paramType {
	case ParamTypeDate:
		return _ParamDate
	case ParamTypeDatetime:
		return _ParamDatetime
	case ParamTypeInt:
		return _ParamInt
	case ParamTypeUInt:
		return _ParamUInt
	default:
		panic(fmt.Sprintf("unknown paramType: %d", paramType))
	}
}

func (pm _Param) NotAllowEmpty(fieldName string) error {
	return fmt.Errorf("%s cannot be empty", fieldName)
}
func (pm _Param) WrongFormat(fieldName string) error {
	return fmt.Errorf("%s is not a %s", fieldName, pm.Desc)
}
func (pm _Param) WrongKind(fieldName string) error {
	return fmt.Errorf("%s is not a kind:%s", fieldName, pm.ReflectKind)
}
func (pm _Param) WrongIntRange(fieldName string, val, left, right int64) error {
	var __right interface{} = right
	if right == math.MaxInt64 {
		__right = "MaxInt64"
	}
	return fmt.Errorf("int %s:%d is not in range[%d, %v]", fieldName, val, left, __right)
}
func (pm _Param) WrongUIntRange(fieldName string, val, left, right uint64) error {
	var __right interface{} = right
	if right == math.MaxUint64 {
		__right = "MaxUint64"
	}
	return fmt.Errorf("uint %s:%d is not in range[%d, %v]", fieldName, val, left, __right)
}

type TypeRule struct {
	FieldName  string
	AllowEmpty bool // 针对string有效
	ParamType
	IntGreatEq, IntLessEq   int64
	UIntGreatEq, UIntLessEq uint64
}

func CheckParam(val interface{}, rule TypeRule) (err error) {
	if val == nil {
		return fmt.Errorf("val cannot be nil")
	}
	param := getParam(rule.ParamType)

	fieldName := rule.FieldName
	if val == "" {
		if !rule.AllowEmpty {
			return param.NotAllowEmpty(fieldName)
		}
		return
	}

	kind := reflect.TypeOf(val).Kind()
	if param.ReflectKind == reflect.Int64 {
		if kind < reflect.Int || kind > reflect.Int64 {
			return param.WrongKind(fieldName)
		}
	} else if param.ReflectKind == reflect.Uint64 {
		if kind < reflect.Uint || kind > reflect.Uint64 {
			return param.WrongKind(fieldName)
		}
	} else if kind != param.ReflectKind {
		return param.WrongKind(fieldName)
	}

	switch param.Typ {
	case ParamTypeDate:
		_, err = time.ParseInLocation("2006-01-02", val.(string), time.Local)
		if err != nil {
			return param.WrongFormat(fieldName)
		}
	case ParamTypeDatetime:
		_, err = time.ParseInLocation("2006-01-02 15:04:05", val.(string), time.Local)
		if err != nil {
			return param.WrongFormat(fieldName)
		}
	case ParamTypeInt:
		var v int64
		switch val.(type) {
		case int:
			v = int64(val.(int))
		case int8:
			v = int64(val.(int8))
		case int16:
			v = int64(val.(int16))
		case int32:
			v = int64(val.(int32))
		case int64:
			v = int64(val.(int64))
		}
		if v < rule.IntGreatEq || v > rule.IntLessEq {
			return param.WrongIntRange(rule.FieldName, v, rule.IntGreatEq, rule.IntLessEq)
		}
	case ParamTypeUInt:
		var v uint64
		switch val.(type) {
		case uint8:
			v = uint64(val.(uint8))
		case uint:
			v = uint64(val.(uint))
		case uint32:
			v = uint64(val.(uint32))
		case uint64:
			v = uint64(val.(uint64))
		}
		if v < rule.UIntGreatEq || v > rule.UIntLessEq {
			return param.WrongUIntRange(rule.FieldName, v, rule.UIntGreatEq, rule.UIntLessEq)
		}
	}
	return
}

// ParamsCheck 接受结构体或其指针，对成员字段做参数格式检验
// 不能对 非导出字段 检验，会panic
func ParamsCheck(val interface{}, typeRule ...*TypeRule) (err error) {
	if val == nil {
		return fmt.Errorf("ParamsCheck: val is nil")
	}
	if len(typeRule) == 0 {
		return fmt.Errorf("ParamsCheck: need a TypeRule")
	}
	rVal := reflect.ValueOf(val)
	if rVal.Kind() == reflect.Ptr {
		rVal = rVal.Elem() // 最多解一层指针
	}
	if rVal.Kind() != reflect.Struct {
		return fmt.Errorf("ParamsCheck: val must be a struct or struct pointer")
	}
	ruleMap := make(map[string]*TypeRule)
	for _, i := range typeRule {
		ruleMap[i.FieldName] = i
	}

	rTyp := rVal.Type() // 用于获取字段名

	var errmsg string
	for i := 0; i < rVal.NumField(); i++ {
		if rule := ruleMap[rTyp.Field(i).Name]; rule != nil {
			err = CheckParam(rVal.Field(i).Interface(), *rule)
			if err != nil {
				errmsg += err.Error() + "; "
			}
			delete(ruleMap, rule.FieldName)
		}
	}

	for field := range ruleMap {
		errmsg += fmt.Sprintf("field %s not found; ", field)
	}

	if errmsg != "" {
		errmsg = strings.TrimRight(errmsg, "; ")
	}

	if errmsg == "" {
		return nil
	}
	return errors.New(errmsg)
}
