package eastmoney

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"
)

var client = &http.Client{Timeout: 10 * time.Second}

const (
	baseSearch   = "https://searchapi.eastmoney.com/api/suggest/get"
	baseNews     = "https://search-api-web.eastmoney.com/search/jsonp"
	baseReport   = "https://reportapi.eastmoney.com/report/list"
	basePush2Ex  = "https://push2ex.eastmoney.com"
)

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

// ── push2 单股查询 (secid+fields → {data:{...}}) ──

func Push2StockGet(path, secid, fields string) (map[string]json.RawMessage, error) {
	u := fmt.Sprintf("https://push2.eastmoney.com/api/qt/%s?secid=%s&fields=%s", path, secid, fields)
	body, err := DoGet(u)
	if err != nil {
		return nil, err
	}
	var raw struct {
		Data map[string]json.RawMessage `json:"data"`
	}
	if err := json.Unmarshal(body, &raw); err != nil {
		return nil, err
	}
	return raw.Data, nil
}

// ── push2 diff 查询 (secids → {data:{diff:[...]}}) ──

func Push2DiffGet(path, params, fields string) ([]map[string]json.RawMessage, error) {
	u := fmt.Sprintf("https://push2.eastmoney.com/api/qt/%s?%s&fields=%s", path, params, fields)
	body, err := DoGet(u)
	if err != nil {
		return nil, err
	}
	var raw struct {
		Data struct {
			Diff []map[string]json.RawMessage `json:"diff"`
		} `json:"data"`
	}
	if err := json.Unmarshal(body, &raw); err != nil {
		return nil, err
	}
	return raw.Data.Diff, nil
}

// ── push2his K线查询 (klines 逗号分隔) ──

type Push2HisQuery struct {
	Path   string
	SecID  string
	Params string
}

func Push2HisGet(q Push2HisQuery) ([]string, error) {
	u := fmt.Sprintf("https://push2his.eastmoney.com/api/qt/%s?secid=%s&%s", q.Path, q.SecID, q.Params)
	body, err := DoGet(u)
	if err != nil {
		return nil, err
	}
	var raw struct {
		Data struct {
			Klines []string `json:"klines"`
		} `json:"data"`
	}
	if err := json.Unmarshal(body, &raw); err != nil {
		return nil, err
	}
	return raw.Data.Klines, nil
}

// ── datacenter-web 通用查询 ──

type DatacenterQuery struct {
	ReportName  string
	Columns     string
	SortColumns string
	SortTypes   string
	PageSize    int
	Filter      string
	Extra       string
}

func DatacenterGet(q DatacenterQuery) ([]map[string]any, error) {
	cols := q.Columns
	if cols == "" {
		cols = "ALL"
	}
	u := fmt.Sprintf(
		"https://datacenter-web.eastmoney.com/api/data/v1/get?sortColumns=%s&sortTypes=%s&pageSize=%d&pageNumber=1&reportName=%s&columns=%s",
		url.QueryEscape(q.SortColumns),
		url.QueryEscape(q.SortTypes),
		q.PageSize,
		url.QueryEscape(q.ReportName),
		url.QueryEscape(cols),
	)
	if q.Filter != "" {
		u += "&filter=" + url.QueryEscape(q.Filter)
	}
	if q.Extra != "" {
		u += "&" + q.Extra
	}
	body, err := DoGet(u)
	if err != nil {
		return nil, err
	}
	var raw struct {
		Result struct {
			Data []map[string]any `json:"data"`
		} `json:"result"`
	}
	if err := json.Unmarshal(body, &raw); err != nil {
		return nil, err
	}
	return raw.Result.Data, nil
}

// ── push2 clist 通用查询 ──

type ClistQuery struct {
	FS        string
	Fields    string
	Limit     int
	SortField string
	SortOrder string // "1" desc, "0" asc
}

func ClistGet(q ClistQuery) ([]map[string]json.RawMessage, error) {
	sortField := q.SortField
	if sortField == "" {
		sortField = "f3"
	}
	sortOrder := q.SortOrder
	if sortOrder == "" {
		sortOrder = "1"
	}
	u := fmt.Sprintf(
		"https://push2.eastmoney.com/api/qt/clist/get?pn=1&pz=%d&po=%s&np=1&fltt=2&invt=2&fid=%s&fs=%s&fields=%s",
		q.Limit, sortOrder, sortField, url.QueryEscape(q.FS), url.QueryEscape(q.Fields),
	)
	body, err := DoGet(u)
	if err != nil {
		return nil, err
	}
	var raw struct {
		Data struct {
			Diff []map[string]json.RawMessage `json:"diff"`
		} `json:"data"`
	}
	if err := json.Unmarshal(body, &raw); err != nil {
		return nil, err
	}
	return raw.Data.Diff, nil
}
