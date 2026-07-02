package eastmoney

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

var client = &http.Client{Timeout: 10 * time.Second}

// DoGet 发起带 Referer 的 GET 请求，对空响应自动重试
func DoGet(u string) ([]byte, error) {
	delays := []time.Duration{time.Second, 2 * time.Second}
	for _, d := range delays {
		body, err := doGetOnce(u)
		if err == nil && len(body) > 0 {
			return body, nil
		}
		time.Sleep(d)
	}
	return doGetOnce(u)
}

func doGetOnce(u string) ([]byte, error) {
	req, _ := http.NewRequest("GET", u, nil)
	req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36")
	req.Header.Set("Referer", "https://quote.eastmoney.com/")
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	return io.ReadAll(resp.Body)
}

// ToSecID 将6位股票代码转为东方财富 secid 格式
// 上海(6/9开头) → "1.code", 深圳/北交(0/3/8/4) → "0.code"
func ToSecID(code string) string {
	code = strings.TrimSpace(code)
	if strings.HasPrefix(code, "6") || strings.HasPrefix(code, "9") {
		return "1." + code
	}
	return "0." + code
}

// ── JSON 解析辅助 ──

// GetFloat 从 map[string]json.RawMessage 中取 float64
func GetFloat(m map[string]json.RawMessage, key string) float64 {
	v, ok := m[key]
	if !ok {
		return 0
	}
	var f float64
	json.Unmarshal(v, &f)
	return f
}

// GetStr 从 map[string]json.RawMessage 中取 string
func GetStr(m map[string]json.RawMessage, key string) string {
	v, ok := m[key]
	if !ok {
		return ""
	}
	var s string
	json.Unmarshal(v, &s)
	return s
}

// ToFloat 将 any 转为 float64 (适配 datacenter 接口返回)
func ToFloat(v any) float64 {
	if v == nil {
		return 0
	}
	switch t := v.(type) {
	case float64:
		return t
	case json.Number:
		f, _ := t.Float64()
		return f
	case string:
		return ParseFloat(t)
	}
	return 0
}

// ToStr 将 any 转为 string，自动截断日期时间部分
func ToStr(v any) string {
	if v == nil {
		return ""
	}
	switch t := v.(type) {
	case string:
		if len(t) > 10 && t[4] == '-' && t[7] == '-' {
			if idx := strings.Index(t, " "); idx > 0 {
				return t[:idx]
			}
		}
		return t
	case float64:
		return fmt.Sprintf("%.0f", t)
	}
	return fmt.Sprintf("%v", v)
}

// ParseFloat 安全解析字符串为 float64
func ParseFloat(s string) float64 {
	var f float64
	fmt.Sscanf(s, "%f", &f)
	return f
}
