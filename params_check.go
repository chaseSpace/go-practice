package main

import (
	"errors"
	"fmt"
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
	ParamTypeDateTime
)

var (
	_ParamDate     = _Param{Typ: ParamTypeDate, Desc: "Date => {YYYY-MM-DD}", ReflectKind: reflect.String}
	_ParamDatetime = _Param{Typ: ParamTypeDate, Desc: "Datetime => {YYYY-mm-DD HH:MM:SS}", ReflectKind: reflect.String}
)

func getParam(paramType ParamType) _Param {
	switch paramType {
	case ParamTypeDate:
		return _ParamDate
	case ParamTypeDateTime:
		return _ParamDatetime
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

func CheckParam(paramTyp ParamType, fieldName string, val interface{}, allowEmpty bool) (err error) {
	param := getParam(paramTyp)

	if val == "" {
		if !allowEmpty {
			return param.NotAllowEmpty(fieldName)
		}
		return
	}
	if reflect.TypeOf(val).Kind() != param.ReflectKind {
		return param.WrongKind(fieldName)
	}
	switch param.Typ {
	case ParamTypeDate:
		_, err = time.ParseInLocation("2006-01-02", val.(string), time.Local)
		if err != nil {
			return param.WrongFormat(fieldName)
		}
	case ParamTypeDateTime:
		_, err = time.ParseInLocation("2006-01-02 15:04:05", val.(string), time.Local)
		if err != nil {
			return param.WrongFormat(fieldName)
		}
	}
	return
}

type TypeRule struct {
	FieldName  string
	AllowEmpty bool
	ParamType
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
			err = CheckParam(rule.ParamType, rule.FieldName, rVal.Field(i).Interface(), rule.AllowEmpty)
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

	return errors.New(errmsg)
}
