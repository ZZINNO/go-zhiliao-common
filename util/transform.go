package util

import (
	"bytes"
	"encoding/json"
	"fmt"
	"golang.org/x/text/encoding/simplifiedchinese"
	"golang.org/x/text/transform"
	"io/ioutil"
	"reflect"
	"strconv"
	"time"
	"unicode"
)

// S 字符串类型转换
type S string

func (s S) String() string {
	return string(s)
}

// Bytes 转换为[]byte
func (s S) Bytes() []byte {
	return []byte(s)
}

// Int64 转换为int64
func (s S) Int64() int64 {
	i, err := strconv.ParseInt(s.String(), 10, 64)
	if err != nil {
		return 0
	}
	return i
}

// Int 转换为int
func (s S) Int() int {
	return int(s.Int64())
}

// Uint 转换为uint
func (s S) Uint() uint {
	return uint(s.Uint64())
}

// Uint64 转换为uint64
func (s S) Uint64() uint64 {
	i, err := strconv.ParseUint(s.String(), 10, 64)
	if err != nil {
		return 0
	}
	return i
}

// Float64 转换为float64
func (s S) Float64() float64 {
	f, err := strconv.ParseFloat(s.String(), 64)
	if err != nil {
		return 0
	}
	return f
}

// ToJSON 转换为JSON
func (s S) ToJSON(v interface{}) error {
	return json.Unmarshal(s.Bytes(), v)
}

//Struct 转 map
func StructToMap(obj interface{}) map[string]interface{} {
	t := reflect.TypeOf(obj)
	v := reflect.ValueOf(obj)

	var data = make(map[string]interface{})
	for i := 0; i < t.NumField(); i++ {
		data[t.Field(i).Name] = v.Field(i).Interface()
	}
	return data
}

//interface{} 转换为 bool
func Bool(inter interface{}) bool {
	switch t := inter.(type) {
	case bool:
		return t
	case string:
		if t == "true" || t == "1" {
			return true
		}
		return false
	case int, int8, int16, int32, int64, float32, float64:
		if Int(t) == 1 {
			return true
		}
		return false
	default:
		//db.Engine.Logger().Infof("%s", reflect.TypeOf(t))
		return false
	}
}

// interface{} 转换为 int
func Int(num interface{}) int {
	switch t := num.(type) {
	case int:
		return t
	case int16:
		return int(t)
	case int32:
		return int(t)
	case int64:
		return int(t)
	case float32:
		return int(t)
	case float64:
		return int(t)
	case string:
		ret, _ := strconv.Atoi(t)
		return ret
	case bool:
		if t {
			return 1
		}
		return 0
	default:
		return -1
	}
}

func Int64(num interface{}) int64 {
	switch t := num.(type) {
	case int:
		return int64(t)
	case int16:
		return int64(t)
	case int32:
		return int64(t)
	case int64:
		return t
	case float32:
		return int64(t)
	case float64:
		return int64(t)
	case string:
		ret, _ := strconv.Atoi(t)
		return int64(ret)
	default:
		return -1
	}
}

// interface{} 转换为 []int{}
func IntSlice(inter interface{}) []int {
	switch t := inter.(type) {
	case []interface{}:
		sli2 := make([]int, 0)
		for _, v := range t {
			sli2 = append(sli2, Int(v))
		}
		return sli2
	case []int:
		return t
	case []string:
		sli2 := make([]int, 0)
		for _, i := range t {
			a, _ := strconv.Atoi(i)
			sli2 = append(sli2, a)
		}
		return sli2
	case int:
		return []int{t}
	case int16:
		return []int{int(t)}
	case int32:
		return []int{int(t)}
	case int64:
		return []int{int(t)}
	case float32:
		return []int{int(t)}
	case float64:
		return []int{int(t)}
	default:
		fmt.Println(reflect.TypeOf(t))
		return nil
	}
}

func Int64Slice(inter interface{}) []int64 {
	switch t := inter.(type) {
	case []interface{}:
		sli2 := make([]int64, 0)
		for _, v := range t {
			sli2 = append(sli2, Int64(v))
		}
		return sli2
	case []string:
		sli2 := make([]int64, 0)
		for _, i := range t {
			a, _ := strconv.Atoi(i)
			sli2 = append(sli2, int64(a))
		}
		return sli2
	case int:
		return []int64{int64(t)}
	case int16:
		return []int64{int64(t)}
	case int32:
		return []int64{int64(t)}
	case int64:
		return []int64{t}
	case float32:
		return []int64{int64(t)}
	case float64:
		return []int64{int64(t)}
	default:
		return nil
	}
}

// interface{} 转换为 float64
func Float64(num interface{}) float64 {
	switch t := num.(type) {
	case int:
		return float64(t)
	case int16:
		return float64(t)
	case int32:
		return float64(t)
	case int64:
		return float64(t)
	case float32:
		return float64(t)
	case float64:
		return t
	case string:
		f, err := strconv.ParseFloat(t, 64)
		if err != nil {
			return -1
		}
		return f
	default:
		return -1
	}
}

func Float64Slice(inter interface{}) []float64 {
	switch t := inter.(type) {
	case []interface{}:
		sli2 := make([]float64, 0)
		for _, v := range t {
			sli2 = append(sli2, Float64(v))
		}
		return sli2
	case []float64:
		return t
	case []string:
		sli2 := make([]float64, 0)
		for _, i := range t {
			a, _ := strconv.ParseFloat(i, 64)
			sli2 = append(sli2, a)
		}
		return sli2
	case int:
		return []float64{float64(t)}
	case int16:
		return []float64{float64(t)}
	case int32:
		return []float64{float64(t)}
	case int64:
		return []float64{float64(t)}
	case float32:
		return []float64{float64(t)}
	case float64:
		return []float64{t}
	default:
		fmt.Println(reflect.TypeOf(t))
		return nil
	}
}

func IntMap(inter interface{}) map[int]bool {
	switch t := inter.(type) {
	case []int:
		sli2 := make(map[int]bool, 0)
		for _, v := range t {
			sli2[Int(v)] = true
		}
		return sli2
	case []interface{}:
		sli2 := make(map[int]bool, 0)
		for _, v := range t {
			sli2[Int(v)] = true
		}
		return sli2
	case interface{}:
		return map[int]bool{
			Int(inter): true,
		}
	default:
		return nil
	}
}

// interface{} 转换为 string
func String(str interface{}) string {
	switch t := str.(type) {
	case json.Number:
		return t.String()
	case string:
		return t
	case []uint8:
		return string(t)
	case map[string]interface{}:
		value_new, _ := json.Marshal(t)
		return string(value_new[:])
	case int:
		return strconv.Itoa(t)
	case int16:
		return strconv.Itoa(int(t))
	case int32:
		return strconv.Itoa(int(t))
	case int64:
		return strconv.Itoa(int(t))
	case float32:
		//return strconv.FormatFloat(float64(t), 'E', -1, 32)
		return strconv.FormatFloat(float64(t), 'f', 2, 32)
	case float64:
		return strconv.FormatFloat(t, 'f', 2, 64)
		//return strconv.FormatFloat(t, 'E', -1, 64)
	default:
		return ""
	}
}

// interface{} 转换为 []string{}
func StringSlice(inter interface{}) []string {
	switch t := inter.(type) {
	case []interface{}:
		sli2 := make([]string, 0)
		for _, v := range t {
			sli2 = append(sli2, String(v))
		}
		return sli2
	case []int:
		sli2 := make([]string, 0)
		for _, v := range t {
			sli2 = append(sli2, String(v))
		}
		return sli2
	case string:
		if t == "" {
			return nil
		}
		return []string{t}
	case int, float64:
		return []string{String(t)}
	case interface{}:
		ret, ok := t.([]string)
		if ok {
			return ret
		}
		return nil
	default:
		return []string{}
	}
}

// string 转换为 map
func StringToMap(str string) map[string]interface{} {
	ret := make(map[string]interface{})
	if err := json.Unmarshal([]byte(str), &ret); err != nil {
		return nil
	}
	return ret
}

// interface{} 转换为 map
func MapStr(inter interface{}) map[string]interface{} {
	switch t := inter.(type) {
	case map[string]interface{}:
		return t
	case []uint8:
		s := String(t)
		smap := StringToMap(s)
		return smap
	default:
		return nil
	}
}

// interface{} 转换为 []map
func MapSlice(inter interface{}) []map[string]interface{} {
	switch t := inter.(type) {
	case []interface{}:
		sli := make([]map[string]interface{}, 0)
		for _, v := range t {
			sli = append(sli, MapStr(v))
		}
		return sli
	case map[string]interface{}:
		return []map[string]interface{}{t}
	case interface{}:
		ret, ok := inter.([]map[string]interface{})
		if ok {
			return ret
		}
		return nil
	default:
		return nil
	}
}

// interface{} 转换为 []interface{}
func ToInterfaceSlice(inter interface{}) []interface{} {
	switch t := inter.(type) {
	case []int:
		ret := make([]interface{}, len(t))
		for i, v := range t {
			ret[i] = v
		}
		return ret
	case []int8:
		ret := make([]interface{}, len(t))
		for i, v := range t {
			ret[i] = v
		}
		return ret
	case []int16:
		ret := make([]interface{}, len(t))
		for i, v := range t {
			ret[i] = v
		}
		return ret
	case []int32:
		ret := make([]interface{}, len(t))
		for i, v := range t {
			ret[i] = v
		}
		return ret
	case []int64:
		ret := make([]interface{}, len(t))
		for i, v := range t {
			ret[i] = v
		}
		return ret
	case []float32:
		ret := make([]interface{}, len(t))
		for i, v := range t {
			ret[i] = v
		}
		return ret
	case []float64:
		ret := make([]interface{}, len(t))
		for i, v := range t {
			ret[i] = v
		}
		return ret
	case []string:
		ret := make([]interface{}, len(t))
		for i, v := range t {
			ret[i] = v
		}
		return ret
	case int, int8, int16, int32, int64, float32, float64, string:
		return []interface{}{t}
	case []interface{}:
		return t
	case []map[string]interface{}:
		ret := make([]interface{}, len(t))
		for i, v := range t {
			ret[i] = v
		}
		return ret
	default:
		return nil
	}
}

// 数据库不能存储map，将map的value从interface{}转为string
func MapValueMarshal(body interface{}) map[string]interface{} {
	ret := make(map[string]interface{})
	switch t := body.(type) {
	case map[string]interface{}:
		for key, value := range t {
			switch value.(type) {
			case []map[string]interface{}:
				value_new, _ := json.Marshal(value)
				ret[key] = string(value_new[:])
			case map[string]interface{}:
				value_new, _ := json.Marshal(value)
				ret[key] = string(value_new[:])

			default:
				ret[key] = value
			}
		}
		return ret
	default:
		return nil
	}
}

// 删除slice指定索引位
func SliceRemove(slice []map[string]interface{}, index int) []map[string]interface{} {
	return append(slice[:index], slice[index+1:]...)
}

// 数据库提取的map[string]interface{}中字段值类型为[]unit8时, 转换为map
//func RealResult(result map[string]interface{}) map[string]interface{} {
//	for key, value := range result {
//		switch t := value.(type) {
//		case []byte:
//			//ret := new(interface{})
//			//err := json.Unmarshal(t, ret)
//			//if err != nil {
//			//	result[key] = string(t)
//			//} else {
//			//	result[key] = *ret
//			//}
//			ret := make(map[string]interface{})
//			err := json.Unmarshal(t, &ret)
//			if err != nil {
//				ret2 := make([]map[string]interface{}, 0)
//				err := json.Unmarshal(t, &ret2)
//				if err != nil {
//					ret3 := make([]string, 0)
//					err := json.Unmarshal(t, &ret3)
//					if err != nil {
//						result[key] = string(t)
//					} else {
//						result[key] = ret3
//					}
//				} else {
//					result[key] = ret2
//				}
//			} else {
//				result[key] = ret
//			}
//		}
//	}
//	return result
//}

func RealResult(result map[string]interface{}) map[string]interface{} {
	for key, value := range result {
		switch t := value.(type) {
		case []byte:
			ret := new(interface{})
			err := json.Unmarshal(t, ret)
			if err == nil {
				switch t1 := (*ret).(type) {
				case map[string]interface{}:
					result[key] = t1
				case []interface{}:
					if len(t1) > 0 {
						switch t1[0].(type) {
						case string:
							sli2 := make([]string, 0)
							for _, v := range t1 {
								sli2 = append(sli2, v.(string))
							}
							result[key] = sli2
							//result[key] = util.StringSlice(*ret)
						case map[string]interface{}:
							sli2 := make([]map[string]interface{}, 0)
							for _, v := range t1 {
								sli2 = append(sli2, v.(map[string]interface{}))
							}
							result[key] = sli2
							//result[key] = util.MapSlice(*ret)
						default:
							result[key] = string(t)
						}
					} else {
						result[key] = []map[string]interface{}{}
					}

				default:
					result[key] = string(t)
				}
			} else {
				result[key] = string(t)
			}

		}
	}
	return result
}

// 数据库提取的[]map[string]interface{}
func RealResultSlice(resultSlice []map[string]interface{}) []map[string]interface{} {
	resultSlice_new := make([]map[string]interface{}, 0)
	for _, result := range resultSlice {
		resultSlice_new = append(resultSlice_new, RealResult(result))
	}
	return resultSlice_new
}

func CopyMap(m map[string]interface{}) map[string]interface{} {
	cp := make(map[string]interface{})
	for k, v := range m {
		vm, ok := v.(map[string]interface{})
		if ok {
			cp[k] = CopyMap(vm)
		} else {
			cp[k] = v
		}
	}
	return cp
}

func FmtDuration(duration time.Duration) string {
	duration = duration.Round(time.Second)
	switch {
	case duration < time.Minute: // xx秒
		s := duration / time.Second
		return fmt.Sprintf("%02d秒", s)
	case duration >= time.Minute && duration < time.Hour: // xx分xx秒
		m := duration / time.Minute
		duration = duration - m*time.Minute
		s := duration / time.Second
		return fmt.Sprintf("%02d分%02d秒", m, s)
	case duration >= time.Hour && duration < 24*time.Hour: // xx小时xx分
		h := duration / time.Hour
		duration = duration - h*time.Hour
		m := duration / time.Minute
		return fmt.Sprintf("%02d小时%02d分", h, m)
	case duration >= 24*time.Hour: // xx天xx小时
		d := duration / (24 * time.Hour)
		duration = duration - d*time.Hour*24
		h := duration / time.Hour
		return fmt.Sprintf("%02d天%02d小时", d, h)
	default:
		return ""
	}
}

func RemoveDuplicateString(strs []string) []string {
	result := make([]string, 0, len(strs))
	temp := map[string]struct{}{}
	for _, item := range strs {
		if _, ok := temp[item]; !ok {
			temp[item] = struct{}{}
			result = append(result, item)
		}
	}
	return result
}

func RemoveDuplicateInt(ints []int) []int {
	result := make([]int, 0, len(ints))
	temp := map[int]struct{}{}
	for _, item := range ints {
		if _, ok := temp[item]; !ok {
			temp[item] = struct{}{}
			result = append(result, item)
		}
	}
	return result
}

func IsChinese(str string) bool {
	var count int
	for _, v := range str {
		if unicode.Is(unicode.Han, v) {
			count++
			break
		}
	}
	return count > 0
}

type MyStringList []string

func (msl MyStringList) Len() int { return len(msl) }

func (msl MyStringList) Less(i, j int) bool {
	// 数字>字母>汉字

	a, _ := UTF82GBK(msl[i])
	b, _ := UTF82GBK(msl[j])
	bLen := len(b)
	for idx, chr := range a {
		if idx%2 == 0 {
			c, err := GBK2UTF8(a[idx : idx+2])
			if err != nil {
				fmt.Println(err.Error())
			}
			d, err := GBK2UTF8(b[idx : idx+2])
			if err != nil {
				fmt.Println(err.Error())
			}
			if IsChinese(c) && IsChinese(d) {
				var c1, d1 int
				switch c {
				case "一":
					c1 = 1
				case "二":
					c1 = 2
				case "三":
					c1 = 3
				case "四":
					c1 = 4
				case "五":
					c1 = 5
				case "六":
					c1 = 6
				case "七":
					c1 = 7
				case "八":
					c1 = 8
				case "九":
					c1 = 9
				case "十":
					c1 = 10
				default:
					c1 = 0
				}
				switch d {
				case "一":
					d1 = 1
				case "二":
					d1 = 2
				case "三":
					d1 = 3
				case "四":
					d1 = 4
				case "五":
					d1 = 5
				case "六":
					d1 = 6
				case "七":
					d1 = 7
				case "八":
					d1 = 8
				case "九":
					d1 = 9
				case "十":
					d1 = 10
				default:
					d1 = 0
				}
				if c1 != 0 && d1 != 0 {
					return c1 < d1
				}
			}
		}

		if idx > bLen-1 {
			return false
		}

		if chr != b[idx] {
			return chr < b[idx]
		}
	}
	return true
}

func (msl MyStringList) Swap(i, j int) { msl[i], msl[j] = msl[j], msl[i] }

//UTF82GBK : transform UTF8 rune into GBK byte array
func UTF82GBK(src string) ([]byte, error) {
	GB18030 := simplifiedchinese.All[0]
	return ioutil.ReadAll(transform.NewReader(bytes.NewReader([]byte(src)), GB18030.NewEncoder()))
}

//GBK2UTF8 : transform  GBK byte array into UTF8 string
func GBK2UTF8(src []byte) (string, error) {
	GB18030 := simplifiedchinese.All[0]
	bytes1, err := ioutil.ReadAll(transform.NewReader(bytes.NewReader(src), GB18030.NewDecoder()))
	return string(bytes1), err
}
