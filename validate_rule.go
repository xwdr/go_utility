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
func NotEmpty() string {
	return "notEmpty"
}

// RegexpMatch 正则校验 校验输入项是否满足正则表达式
func RegexpMatch(rule string) string {
	return "regexp=" + rule
}

// Lt 小于入参(<)
// 如果为string array Slice则为长度比较
// 如果是 int uint float 则为数值比较
func Lt(mark string) string {
	return "lt=" + mark
}

// Le 小于等于入参(<=)
// 如果为string array Slice则为长度比较
// 如果是 int uint float 则为数值比较
func Le(mark string) string {
	return "le=" + mark
}

// Eq 等于入参(==)
// 如果为string array Slice则为长度比较
// 如果是 int uint float 则为数值比较
func Eq(mark string) string {
	return "eq=" + mark
}

// Ne 不等于入参(!=)
// 如果为string array Slice则为长度比较
// 如果是 int uint float 则为数值比较
func Ne(mark string) string {
	return "ne=" + mark
}

// Ge 大于等于入参(>=)
// 如果为string array Slice则为长度比较
// 如果是 int uint float 则为数值比较
func Ge(mark string) string {
	return "ge=" + mark
}

// Gt 大于入参(>)
// 如果为string array Slice则为长度比较
// 如果是 int uint float 则为数值比较
func Gt(mark string) string {
	return "gt=" + mark
}
