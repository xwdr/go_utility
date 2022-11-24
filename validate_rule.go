package utils

// 参考gin_vue_admin项目
// 定义参数校验规则
var (
	IdRule            = Rules{"id": {NotEmpty()}}
	UserRule          = Rules{"user_id": {NotEmpty()}, "user_name": {NotEmpty()}}
	ResetPasswordRule = Rules{"new_password": {NotEmpty()}, "old_password": {NotEmpty()}}
	MobileRule        = Rules{"mobile": {NotEmpty(), Eq("11")}}
	EmailRule         = Rules{"email": {RegexpMatch(regEmail)}}
)

// NotEmpty 非空 不能为其对应类型的0值
func (r *Rule) SetTip(tip string) *Rule {
	r.Tip = tip
	return r
}

// NotEmpty 非空 不能为其对应类型的0值
func NotEmpty() *Rule {
	return &Rule{Condition: "notEmpty"}
}

// RegexpMatch 正则校验 校验输入项是否满足正则表达式
func RegexpMatch(rule string) *Rule {
	return &Rule{Condition: "regexp=" + rule}
}

// Lt 小于入参(<)
// 如果为string array Slice则为长度比较
// 如果是 int uint float 则为数值比较
func Lt(mark string) *Rule {
	return &Rule{Condition: "lt=" + mark}
}

// Le 小于等于入参(<=)
// 如果为string array Slice则为长度比较
// 如果是 int uint float 则为数值比较
func Le(mark string) *Rule {
	return &Rule{Condition: "le=" + mark}
}

// Eq 等于入参(==)
// 如果为string array Slice则为长度比较
// 如果是 int uint float 则为数值比较
func Eq(mark string) *Rule {
	return &Rule{Condition: "eq=" + mark}
}

// Ne 不等于入参(!=)
// 如果为string array Slice则为长度比较
// 如果是 int uint float 则为数值比较
func Ne(mark string) *Rule {
	return &Rule{Condition: "ne=" + mark}
}

// Ge 大于等于入参(>=)
// 如果为string array Slice则为长度比较
// 如果是 int uint float 则为数值比较
func Ge(mark string) *Rule {
	return &Rule{Condition: "ge=" + mark}
}

// Gt 大于入参(>)
// 如果为string array Slice则为长度比较
// 如果是 int uint float 则为数值比较
func Gt(mark string) *Rule {
	return &Rule{Condition: "gt=" + mark}
}
