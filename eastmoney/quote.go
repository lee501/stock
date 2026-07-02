package eastmoney

import (
	"encoding/json"
	"fmt"
	"net/url"
	"strings"
)

// ════════════════════════════════════════
// 股票搜索
// ════════════════════════════════════════

type SearchResult struct {
	Code     string `json:"code"`
	Name     string `json:"name"`
	Market   string `json:"market"`
	Category string `json:"category"`
}

func SearchStock(keyword string) ([]SearchResult, error) {
	u := fmt.Sprintf("%s?input=%s&type=14&token=D43BF722C8E33BDC906FB84D85E326E8&count=10",
		baseSearch, url.QueryEscape(keyword))
	body, err := DoGet(u)
	if err != nil {
		return nil, err
	}

	var raw struct {
		QuotationCodeTable struct {
			Data []struct {
				Code             string `json:"Code"`
				Name             string `json:"Name"`
				MktNum           string `json:"MktNum"`
				SecurityTypeName string `json:"SecurityTypeName"`
			} `json:"Data"`
		} `json:"QuotationCodeTable"`
	}
	if err := json.Unmarshal(body, &raw); err != nil {
		return nil, err
	}

	var results []SearchResult
	for _, d := range raw.QuotationCodeTable.Data {
		mkt := "sz"
		if d.MktNum == "1" {
			mkt = "sh"
		}
		results = append(results, SearchResult{
			Code:     d.Code,
			Name:     d.Name,
			Market:   mkt,
			Category: d.SecurityTypeName,
		})
	}
	return results, nil
}

// ════════════════════════════════════════
// 实时行情
// ════════════════════════════════════════

type Quote struct {
	Code      string  `json:"code"`
	Name      string  `json:"name"`
	Price     float64 `json:"price"`
	Change    float64 `json:"change"`
	ChangePct float64 `json:"change_pct"`
	Open      float64 `json:"open"`
	High      float64 `json:"high"`
	Low       float64 `json:"low"`
	PreClose  float64 `json:"pre_close"`
	Volume    float64 `json:"volume"`
	Amount    float64 `json:"amount"`
	TurnOver  float64 `json:"turnover"`
	PE        float64 `json:"pe"`
	PB        float64 `json:"pb"`
	MarketCap float64 `json:"market_cap"`
	FloatCap  float64 `json:"float_cap"`
}

func GetQuote(code string) (*Quote, error) {
	m, err := Push2StockGet("stock/get", ToSecID(code),
		"f43,f44,f45,f46,f47,f48,f50,f51,f52,f55,f57,f58,f60,f71,f116,f117,f162,f167,f168,f169,f170")
	if err != nil {
		return nil, err
	}
	if m == nil {
		return nil, fmt.Errorf("stock %s not found", code)
	}
	div := func(key string) float64 { return GetFloat(m, key) / 1000 }

	return &Quote{
		Code:      GetStr(m, "f57"),
		Name:      GetStr(m, "f58"),
		Price:     div("f43"),
		Open:      div("f46"),
		High:      div("f44"),
		Low:       div("f45"),
		PreClose:  div("f60"),
		Change:    div("f169"),
		ChangePct: GetFloat(m, "f170") / 100,
		Volume:    GetFloat(m, "f47"),
		Amount:    GetFloat(m, "f48"),
		TurnOver:  GetFloat(m, "f168") / 100,
		PE:        GetFloat(m, "f162") / 100,
		PB:        GetFloat(m, "f167") / 100,
		MarketCap: GetFloat(m, "f116"),
		FloatCap:  GetFloat(m, "f117"),
	}, nil
}

// ════════════════════════════════════════
// K线数据
// ════════════════════════════════════════

type KLine struct {
	Date   string  `json:"date"`
	Open   float64 `json:"open"`
	Close  float64 `json:"close"`
	High   float64 `json:"high"`
	Low    float64 `json:"low"`
	Volume float64 `json:"volume"`
	Amount float64 `json:"amount"`
	Change float64 `json:"change"`
}

func GetKLine(code, period string, limit int) ([]KLine, error) {
	klt := "101"
	switch period {
	case "weekly":
		klt = "102"
	case "monthly":
		klt = "103"
	}
	if limit <= 0 || limit > 120 {
		limit = 30
	}

	lines, err := Push2HisGet(Push2HisQuery{
		Path:  "stock/kline/get",
		SecID: ToSecID(code),
		Params: fmt.Sprintf("klt=%s&fqt=1&beg=0&end=20500101&lmt=%d&fields1=f1,f2,f3,f4,f5,f6&fields2=f51,f52,f53,f54,f55,f56,f57,f58,f59,f60,f61",
			klt, limit),
	})
	if err != nil {
		return nil, err
	}

	var klines []KLine
	for _, line := range lines {
		parts := strings.Split(line, ",")
		if len(parts) < 11 {
			continue
		}
		klines = append(klines, KLine{
			Date:   parts[0],
			Open:   ParseFloat(parts[1]),
			Close:  ParseFloat(parts[2]),
			High:   ParseFloat(parts[3]),
			Low:    ParseFloat(parts[4]),
			Volume: ParseFloat(parts[5]),
			Amount: ParseFloat(parts[6]),
			Change: ParseFloat(parts[8]),
		})
	}
	return klines, nil
}

// ════════════════════════════════════════
// 指数行情
// ════════════════════════════════════════

type IndexQuote struct {
	Code      string  `json:"code"`
	Name      string  `json:"name"`
	Price     float64 `json:"price"`
	ChangePct float64 `json:"change_pct"`
	Change    float64 `json:"change"`
	Volume    float64 `json:"volume"`
	Amount    float64 `json:"amount"`
	Advance   int     `json:"advance"`
	Decline   int     `json:"decline"`
}

func GetIndexQuotes() ([]IndexQuote, error) {
	codes := []string{
		"1.000001", "0.399001", "0.399006",
		"1.000688", "1.000300", "1.000905",
		"0.399673",
	}

	diff, err := Push2DiffGet("ulist/get",
		"secids="+strings.Join(codes, ","),
		"f2,f3,f4,f5,f6,f12,f14,f104,f105")
	if err != nil {
		return nil, err
	}

	var indexes []IndexQuote
	for _, m := range diff {
		indexes = append(indexes, IndexQuote{
			Code:      GetStr(m, "f12"),
			Name:      GetStr(m, "f14"),
			Price:     GetFloat(m, "f2"),
			ChangePct: GetFloat(m, "f3"),
			Change:    GetFloat(m, "f4"),
			Volume:    GetFloat(m, "f5"),
			Amount:    GetFloat(m, "f6"),
			Advance:   int(GetFloat(m, "f104")),
			Decline:   int(GetFloat(m, "f105")),
		})
	}
	return indexes, nil
}
