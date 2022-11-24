package utils

import (
	"errors"
	"reflect"
	"regexp"
	"strconv"
	"strings"
)

// 	Rule 规则信息,包括规则和、错误提示
type Rule struct {
	Condition string `json:"condition"` // 条件
	Tip       string `json:"tip"`       // 提示
}

// 规则Rules
type Rules map[string][]*Rule

var CustomizeMap = make(map[string]Rules)

// compareMap 校验规则
var compareMap = map[string]bool{
	"lt": true,
	"le": true,
	"eq": true,
	"ne": true,
	"ge": true,
	"gt": true,
}

// ValidateFunc 选项函数定义
type ValidateFunc func() (err error)

// ValidateExec 函数执行
func ValidateExec(opts ...ValidateFunc) error {
	for _, opt := range opts {
		if err := opt(); err != nil {
			return err
		}
	}
	return nil
}

// Handler 校验函数
func Handler(st interface{}, roleMap Rules) ValidateFunc {
	return func() error {
		return Validate(st, roleMap)
	}
}

// RegisterRule 注册自定义规则方案建议在路由初始化层即注册
func RegisterRule(key string, rule Rules) (err error) {
	if CustomizeMap[key] != nil {
		return errors.New(key + "已注册,无法重复注册")
	} else {
		CustomizeMap[key] = rule
		return nil
	}
}

// Validate 校验方法
func Validate(st interface{}, roleMap Rules) error {
	typ := reflect.TypeOf(st)
	val := reflect.ValueOf(st) // 获取reflect.Type类型

	kd := val.Kind() // 获取到st对应的类别
	if kd != reflect.Struct {
		return errors.New("expect struct")
	}
	// 遍历结构体的所有字段
	for i := 0; i < val.NumField(); i++ {
		name := typ.Field(i).Tag.Get("json")
		val := val.Field(i)
		if len(roleMap[name]) == 0 {
			continue
		}
		// 校验参数表规则
		for _, rule := range roleMap[name] {
			if err := execRule(val, name, rule); err != nil {
				return err
			}
		}
	}
	return nil
}

// execRule 执行规则
func execRule(val reflect.Value, name string, rule *Rule) error {
	switch {
	case rule.Condition == "notEmpty":
		if isBlank(val) {
			return msgTip(errors.New(name+"值不能为空"), rule)
		}
	case strings.Split(rule.Condition, "=")[0] == "regexp":
		if !regexpMatch(strings.Split(rule.Condition, "=")[1], val.String()) {
			return msgTip(errors.New(name+"格式校验不通过"), rule)
		}
	case compareMap[strings.Split(rule.Condition, "=")[0]]:
		if !compareVerify(val, rule.Condition) {
			return msgTip(errors.New(name+"长度或值不在合法范围"), rule)
		}
	}
	return nil
}

// msgTip 消息提示
func msgTip(orgErr error, rule *Rule) error {
	return If(len(rule.Tip) == 0, orgErr, errors.New(rule.Tip)).(error)
}

// compareVerify 长度和数字的校验方法 根据类型自动校验
func compareVerify(value reflect.Value, VerifyStr string) bool {
	switch value.Kind() {
	case reflect.String:
		return compare(value.String(), VerifyStr)
	case reflect.Slice, reflect.Array:
		return compare(value.Len(), VerifyStr)
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		return compare(value.Uint(), VerifyStr)
	case reflect.Float32, reflect.Float64:
		return compare(value.Float(), VerifyStr)
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return compare(value.Int(), VerifyStr)
	default:
		return false
	}
}

// isBlank 非空校验
func isBlank(value reflect.Value) bool {
	switch value.Kind() {
	case reflect.String, reflect.Slice:
		return value.Len() == 0
	case reflect.Bool:
		return !value.Bool()
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return value.Int() == 0
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		return value.Uint() == 0
	case reflect.Float32, reflect.Float64:
		return value.Float() == 0
	case reflect.Interface, reflect.Ptr:
		return value.IsNil()
	}
	return reflect.DeepEqual(value.Interface(), reflect.Zero(value.Type()).Interface())
}

// compare 比较函数
func compare(value interface{}, VerifyStr string) bool {
	VerifyStrArr := strings.Split(VerifyStr, "=")
	val := reflect.ValueOf(value)
	switch val.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		VInt, VErr := strconv.ParseInt(VerifyStrArr[1], 10, 64)
		if VErr != nil {
			return false
		}
		switch {
		case VerifyStrArr[0] == "lt":
			return val.Int() < VInt
		case VerifyStrArr[0] == "le":
			return val.Int() <= VInt
		case VerifyStrArr[0] == "eq":
			return val.Int() == VInt
		case VerifyStrArr[0] == "ne":
			return val.Int() != VInt
		case VerifyStrArr[0] == "ge":
			return val.Int() >= VInt
		case VerifyStrArr[0] == "gt":
			return val.Int() > VInt
		default:
			return false
		}
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		VInt, VErr := strconv.Atoi(VerifyStrArr[1])
		if VErr != nil {
			return false
		}
		switch {
		case VerifyStrArr[0] == "lt":
			return val.Uint() < uint64(VInt)
		case VerifyStrArr[0] == "le":
			return val.Uint() <= uint64(VInt)
		case VerifyStrArr[0] == "eq":
			return val.Uint() == uint64(VInt)
		case VerifyStrArr[0] == "ne":
			return val.Uint() != uint64(VInt)
		case VerifyStrArr[0] == "ge":
			return val.Uint() >= uint64(VInt)
		case VerifyStrArr[0] == "gt":
			return val.Uint() > uint64(VInt)
		default:
			return false
		}
	case reflect.Float32, reflect.Float64:
		VFloat, VErr := strconv.ParseFloat(VerifyStrArr[1], 64)
		if VErr != nil {
			return false
		}
		switch {
		case VerifyStrArr[0] == "lt":
			return val.Float() < VFloat
		case VerifyStrArr[0] == "le":
			return val.Float() <= VFloat
		case VerifyStrArr[0] == "eq":
			return val.Float() == VFloat
		case VerifyStrArr[0] == "ne":
			return val.Float() != VFloat
		case VerifyStrArr[0] == "ge":
			return val.Float() >= VFloat
		case VerifyStrArr[0] == "gt":
			return val.Float() > VFloat
		default:
			return false
		}
	case reflect.String:
		return VerifyStrArr[1] == value
	default:
		return false
	}
}

// regexpMatch 正则匹配
func regexpMatch(rule, matchStr string) bool {
	return regexp.MustCompile(rule).MatchString(matchStr)
}
