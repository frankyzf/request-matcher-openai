package util

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"math/rand"
	"regexp"
	"strconv"
	"strings"
	"time"

	guid "github.com/bsm/go-guid"
)

func GetMd5String(s string) string {
	h := md5.New()
	h.Write([]byte(s))
	return hex.EncodeToString(h.Sum(nil))
}

//生成Guid字串
func GetGuid() string {
	id := guid.New96()
	return hex.EncodeToString(id.Bytes())
}

func GenerateOrderId() string {
	id := strconv.FormatInt(time.Now().UTC().UnixNano()+int64(rand.Intn(1000)), 10)
	return id
}

func GetUUID() string {
	// id, _ := uuid.NewUUID()
	// return id.String()
	return GetUUID_12Bytes()
}

func GetUUID_12Bytes() string {
	id := guid.New96()
	return hex.EncodeToString(id.Bytes())
}

func ParseFloat(f string) float64 {
	r, _ := strconv.ParseFloat(f, 64)
	return r
}

func ParseInt(f string) int64 {
	r, _ := strconv.ParseInt(f, 10, 64)
	return r
}

func GetTimeDuration(sec float64) time.Duration {
	d, _ := time.ParseDuration(fmt.Sprintf("%fs", sec))
	return d
}

func GetRandomValue(rg float64) float64 {
	return float64((rand.Float64() - 0.5) * rg)
}

func GetMax(a, b int64) int64 {
	if a > b {
		return a
	}
	return b
}

func IsZero(quantity float64) bool {
	return quantity == float64(0)
}

func GetTodayDate() string {
	day := time.Now().String()[0:10]
	return day
}

func TrimHtml(src string) string {
	//将HTML标签全转换成小写
	re, _ := regexp.Compile("\\<[\\S\\s]+?\\>")
	src = re.ReplaceAllStringFunc(src, strings.ToLower)
	//去除STYLE
	re, _ = regexp.Compile("\\<style[\\S\\s]+?\\</style\\>")
	src = re.ReplaceAllString(src, "")
	//去除SCRIPT
	re, _ = regexp.Compile("\\<script[\\S\\s]+?\\</script\\>")
	src = re.ReplaceAllString(src, "")
	//去除所有尖括号内的HTML代码，并换成换行符
	re, _ = regexp.Compile("\\<[\\S\\s]+?\\>")
	src = re.ReplaceAllString(src, "\n")
	//去除连续的换行符
	re, _ = regexp.Compile("\\s{2,}")
	src = re.ReplaceAllString(src, "\n")
	return strings.TrimSpace(src)
}

func In(target string, str_array []string) bool {
	for _, element := range str_array {
		if target == element {
			return true
		}
	}
	return false
}
