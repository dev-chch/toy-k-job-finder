package helper

import (
	"fmt"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"unicode"
)

// 에러를 출력한다.
func CheckErr(err error) {
	if err != nil {
		fmt.Println(err)
	}
}

// 연결을 확인한다.
func CheckConnect(res *http.Response) {
	if res.StatusCode != 200 {
		fmt.Println("Request failed with Status: ", res.StatusCode)
	}
}

// URL 경로를 붙인다
func UrlDirectoryBuild(directorys []string, siteUrl string) string {
	var UrlDirectoryBuild strings.Builder

	UrlDirectoryBuild.WriteString(siteUrl)
	separator := "/"
	for idx, _ := range directorys {
		UrlDirectoryBuild.WriteString(separator)
		UrlDirectoryBuild.WriteString(directorys[idx])
	}

	return UrlDirectoryBuild.String()
}

// URL 파라미터를 붙인다.
func UrlParamBuild(params map[string]string, siteUrl string) string {
	var UrlParamBuild strings.Builder

	UrlParamBuild.WriteString(siteUrl)
	separator := "?"
	i := 0
	for param, val := range params {
		if i > 0 {
			separator = "&"
		}
		UrlParamBuild.WriteString(separator)
		UrlParamBuild.WriteString(param)
		UrlParamBuild.WriteString("=")
		UrlParamBuild.WriteString(val)
		i++
	}

	return UrlParamBuild.String()
}

// 문자열에서 띄워쓰기를 제거한다.
func StripSpace(str string) string {
	var b strings.Builder
	b.Grow(len(str))
	for _, ch := range str {
		if !unicode.IsSpace(ch) {
			b.WriteRune(ch)
		}
	}
	return b.String()
}

// 문자열에서 태그를 제거한다.
func CleanString(str string) string {
	return strings.Join(strings.Fields(strings.TrimSpace(str)), " ")
}

// 문자열에서 숫자를 추출해 순서대로 배열로 반환한다.
func GetNumFromString(strings string) []string {
	numbers := []string{}
	extractedNumbers := regexp.MustCompile(`[-]?\d[\d]*[\]?[\d{2}]*`).FindAllString(strings, -1)
	for _, number := range extractedNumbers {
		numbers = append(numbers, number)
	}
	return numbers
}

// 1000 -> 1,000
func NumberFormat(n int64) string {
	in := strconv.FormatInt(n, 10)
	numOfDigits := len(in)
	if n < 0 {
		numOfDigits-- // First character is the - sign (not a digit)
	}
	numOfCommas := (numOfDigits - 1) / 3

	out := make([]byte, len(in)+numOfCommas)
	if n < 0 {
		in, out[0] = in[1:], '-'
	}

	for i, j, k := len(in)-1, len(out)-1, 0; ; i, j = i-1, j-1 {
		out[j] = in[i]
		if i == 0 {
			return string(out)
		}
		if k++; k == 3 {
			j, k = j-1, 0
			out[j] = ','
		}
	}
}
