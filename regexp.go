package utils

import (
	"regexp"
)

// 正则表达式相关方法
const (
	regMobile = "^((13[0-9])|(14[5,7])|(15[0-3,5-9])|(17[0,3,5-8])|(18[0-9])|166|198|199|(147))\\d{8}$"
	regEmail  = "^([a-zA-Z0-9_.-])+@([a-zA-Z0-9_-])+(.[a-zA-Z0-9_-])+"
	regIdCard = "^[1-9]\\d{5}(18|19|20)\\d{2}((0[1-9])|(1[0-2]))(([0-2][1-9])|10|20|30|31)\\d{3}[0-9Xx]$"
)

// 判断手机号是否符合规则
func IsValidMobile(mobile string) bool {
	if m, _ := regexp.MatchString(regMobile, mobile); !m {
		return false
	}
	return true
}

// 判断邮箱是否符合规则
func IsValidEmail(email string) bool {
	if m, _ := regexp.MatchString(regEmail, email); !m {
		return false
	}
	return true
}

// 判断省份证号是否符合规则
func IsValidIdCard(id string) bool {
	if m, _ := regexp.MatchString(regIdCard, id); !m {
		return false
	}
	return true
}
