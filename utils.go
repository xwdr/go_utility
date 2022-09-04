package utils

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"reflect"
	"strconv"
	"strings"
	"time"
	"unicode/utf8"
	"unsafe"

	"github.com/spf13/cast"
)

// StringToBytes converts string to byte slice without a memory allocation.
func StringToBytes(s string) (b []byte) {
	sh := *(*reflect.StringHeader)(unsafe.Pointer(&s))
	bh := (*reflect.SliceHeader)(unsafe.Pointer(&b))
	bh.Data, bh.Len, bh.Cap = sh.Data, sh.Len, sh.Len
	return b
}

// BytesToString converts byte slice to string without a memory allocation.
func BytesToString(b []byte) string {
	return *(*string)(unsafe.Pointer(&b))
}

// string to int64 slice
func SplitToInt64s(str string) (result []int64) {
	if len(str) == 0 {
		return
	}
	strArr := strings.Split(str, ",")
	result = make([]int64, 0, len(strArr))
	for _, v := range strArr {
		result = append(result, cast.ToInt64(v))
	}
	return
}

// string to int32 slice
func SplitToInt32s(str string) (result []int32) {
	if len(str) == 0 {
		return
	}
	strArr := strings.Split(str, ",")
	result = make([]int32, 0, len(strArr))
	for _, v := range strArr {
		result = append(result, cast.ToInt32(v))
	}
	return
}

// string to int slice
func SplitToInts(str string) (result []int) {
	if len(str) == 0 {
		return
	}
	strArr := strings.Split(str, ",")
	result = make([]int, 0, len(strArr))
	for _, v := range strArr {
		result = append(result, cast.ToInt(v))
	}
	return
}

// HidePhone 手机号屏蔽  屏蔽中间4位
func HidePhone(mobile string) string {
	if len(mobile) == 0 || utf8.RuneCountInString(mobile) != 11 {
		return mobile
	}
	mobileRune := []rune(mobile)
	begin := string(mobileRune[0:3])
	end := string(mobileRune[7:])
	return begin + "****" + end
}

// StructToMap
func StructToMap(m interface{}) map[string]interface{} {
	str, err := json.Marshal(m)
	if err != nil {
		return map[string]interface{}{}
	}
	res := map[string]interface{}{}
	_ = json.Unmarshal(str, &res)
	return res
}

// StructToMap
func StructToSliceMap(m interface{}) []map[string]interface{} {
	str, err := json.Marshal(m)
	if err != nil {
		return []map[string]interface{}{}
	}
	res := []map[string]interface{}{}
	_ = json.Unmarshal(str, &res)
	return res
}

// StringToTime 时间格式转换
func StringToTime(timeLayout, strTime string) time.Time {
	t, _ := time.ParseInLocation(timeLayout, strTime, time.Local)
	return t
}

// 根据tagName的值做可以将struct转map
func StructToMapByTagName(obj interface{}, tagName string) map[string]interface{} {
	t := reflect.TypeOf(obj)
	v := reflect.ValueOf(obj)
	data := make(map[string]interface{})
	for i := 0; i < t.NumField(); i++ {
		data[t.Field(i).Tag.Get(tagName)] = v.Field(i).Interface()
	}
	return data
}

// map key to array
func MapKeyToArray(m map[int]interface{}) []int {
	if len(m) == 0 {
		return nil
	}
	arr := make([]int, 0, len(m))
	for k, _ := range m {
		arr = append(arr, int(k))
	}
	return arr
}

// 获取年龄
// @param birthday string datetime 格式 yyyy-mm-dd
func GetAge(birthday string) (int, error) {
	if len(birthday) == 0 {
		return 0, fmt.Errorf("日期不能为空")
	}
	var age int
	_, err := time.ParseInLocation("2006-01-02", birthday, time.Local)
	if err != nil {
		return age, fmt.Errorf("日期格式不正确")
	}
	curTime := time.Now().Format("2006-01-02")
	bt := strings.Split(birthday, "-")
	ct := strings.Split(curTime, "-")
	by, _ := strconv.Atoi(bt[0])
	cy, _ := strconv.Atoi(ct[0])
	age = cy - by
	bmd, _ := strconv.Atoi(bt[1] + bt[2])
	cmd, _ := strconv.Atoi(ct[1] + ct[2])
	if cmd < bmd {
		age -= 1
	}
	return age, nil
}

// SliceUniqueString 字符串切片去重
func SliceUniqueString(strs []string) []string {
	if len(strs) == 0 {
		return strs
	}
	result := make([]string, 0, len(strs))
	temp := map[string]bool{}
	for _, str := range strs {
		if _, ok := temp[str]; ok {
			continue
		}
		temp[str] = true
		result = append(result, str)
	}
	return result
}

// SliceUniqueInt int切片去重
func SliceUniqueInt(ints []int) []int {
	if len(ints) == 0 {
		return ints
	}
	result := make([]int, 0, len(ints))
	temp := map[int]bool{}
	for _, v := range ints {
		if _, ok := temp[v]; ok {
			continue
		}
		temp[v] = true
		result = append(result, v)
	}
	return result
}

// string to md5
func ToMd5(str string) string {
	hash := md5.New()
	hash.Write([]byte(str))
	return hex.EncodeToString(hash.Sum(nil))
}

// InArray 是否包含在数组中
func InArray(needle interface{}, haystack interface{}) bool {
	val := reflect.ValueOf(haystack)
	switch val.Kind() {
	case reflect.Slice, reflect.Array:
		for i := 0; i < val.Len(); i++ {
			if reflect.DeepEqual(needle, val.Index(i).Interface()) {
				return true
			}
		}
	case reflect.Map:
		for _, k := range val.MapKeys() {
			if reflect.DeepEqual(needle, val.MapIndex(k).Interface()) {
				return true
			}
		}
	default:
		panic("haystack: haystack type muset be slice, array or map")
	}
	return false
}

// MergeMap 两个map合并
func MergeMap(m1 map[string]interface{}, m2 map[string]interface{}) map[string]interface{} {
	if m1 == nil {
		return m2
	}
	if m2 == nil {
		return m1
	}
	for k, v := range m2 {
		m1[k] = v
	}
	return m1
}
