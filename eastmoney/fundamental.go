package eastmoney

import (
	"encoding/json"
	"fmt"
	"net/url"
)

// ════════════════════════════════════════
// 核心财务指标
// ════════════════════════════════════════

type Financial struct {
	Code         string  `json:"code"`
	Name         string  `json:"name"`
	ROE          float64 `json:"roe"`
	GrossMargin  float64 `json:"gross_margin"`
	NetMargin    float64 `json:"net_margin"`
	RevenueYoY   float64 `json:"revenue_yoy"`
	NetProfitYoY float64 `json:"net_profit_yoy"`
	DebtRatio    float64 `json:"debt_ratio"`
	EPS          float64 `json:"eps"`
	BPS          float64 `json:"bps"`
	CashFlowPS   float64 `json:"cashflow_ps"`
}

func GetFinancial(code string) (*Financial, error) {
	u := fmt.Sprintf(
		"https://push2.eastmoney.com/api/qt/stock/get?secid=%s&fields=f57,f58,f173,f183,f184,f185,f186,f187,f188,f190,f191,f192,f198,f199",
		ToSecID(code),
	)
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
	if raw.Data == nil {
		return nil, fmt.Errorf("no financial data for %s", code)
	}

	m := raw.Data
	return &Financial{
		Code:         GetStr(m, "f57"),
		Name:         GetStr(m, "f58"),
		ROE:          GetFloat(m, "f173") / 100,
		GrossMargin:  GetFloat(m, "f186") / 100,
		NetMargin:    GetFloat(m, "f184") / 100,
		RevenueYoY:   GetFloat(m, "f183") / 100,
		NetProfitYoY: GetFloat(m, "f185") / 100,
		DebtRatio:    GetFloat(m, "f188") / 100,
		EPS:          GetFloat(m, "f183") / 1000,
		BPS:          GetFloat(m, "f192") / 1000,
		CashFlowPS:   GetFloat(m, "f190") / 1000,
	}, nil
}

// ════════════════════════════════════════
// 十大流通股东
// ════════════════════════════════════════

type TopHolder struct {
	Rank        int     `json:"rank"`
	Name        string  `json:"name"`
	HoldCount   float64 `json:"hold_count"`
	HoldRatio   float64 `json:"hold_ratio"`
	Change      string  `json:"change"`
	ChangeCount float64 `json:"change_count"`
	HolderType  string  `json:"holder_type"`
	ReportDate  string  `json:"report_date"`
}

func GetTopHolders(code string) ([]TopHolder, error) {
	filter := fmt.Sprintf(`(SECURITY_CODE="%s")`, code)
	u := fmt.Sprintf(
		"https://datacenter-web.eastmoney.com/api/data/v1/get?sortColumns=RANK&sortTypes=1&pageSize=10&pageNumber=1&reportName=RPT_F10_EH_FREEHOLDERS&columns=ALL&filter=%s",
		url.QueryEscape(filter),
	)
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

	var holders []TopHolder
	for _, d := range raw.Result.Data {
		holders = append(holders, TopHolder{
			Rank:        int(ToFloat(d["RANK"])),
			Name:        ToStr(d["FREE_HOLDNUM_NAME"]),
			HoldCount:   ToFloat(d["FREE_HOLDNUM"]),
			HoldRatio:   ToFloat(d["FREE_RATIO_QSZ"]),
			Change:      ToStr(d["IS_HOLDORG"]),
			ChangeCount: ToFloat(d["HOLD_NUM_CHANGE"]),
			HolderType:  ToStr(d["HOLDER_TYPE"]),
			ReportDate:  ToStr(d["END_DATE"]),
		})
	}
	return holders, nil
}

// ════════════════════════════════════════
// 分红送配历史
// ════════════════════════════════════════

type Dividend struct {
	ReportDate string  `json:"report_date"`
	Plan       string  `json:"plan"`
	ExDate     string  `json:"ex_date"`
	RecordDate string  `json:"record_date"`
	PayDate    string  `json:"pay_date"`
	BonusRatio float64 `json:"bonus_ratio"`
	TransRatio float64 `json:"trans_ratio"`
	CashDiv    float64 `json:"cash_div"`
}

func GetDividendHistory(code string, limit int) ([]Dividend, error) {
	if limit <= 0 || limit > 20 {
		limit = 10
	}
	filter := fmt.Sprintf(`(SECURITY_CODE="%s")`, code)
	u := fmt.Sprintf(
		"https://datacenter-web.eastmoney.com/api/data/v1/get?sortColumns=EX_DIVIDEND_DATE&sortTypes=-1&pageSize=%d&pageNumber=1&reportName=RPT_SHAREBONUS_DET&columns=ALL&filter=%s",
		limit, url.QueryEscape(filter),
	)
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

	var divs []Dividend
	for _, d := range raw.Result.Data {
		divs = append(divs, Dividend{
			ReportDate: ToStr(d["REPORT_DATE"]),
			Plan:       ToStr(d["ASSIGN_DETAIL"]),
			ExDate:     ToStr(d["EX_DIVIDEND_DATE"]),
			RecordDate: ToStr(d["EQUITY_RECORD_DATE"]),
			PayDate:    ToStr(d["PAY_CASH_DATE"]),
			BonusRatio: ToFloat(d["BONUS_IT_RATIO"]),
			TransRatio: ToFloat(d["TRANSFER_IT_RATIO"]),
			CashDiv:    ToFloat(d["PRETAX_BONUS_RMB"]),
		})
	}
	return divs, nil
}

// ════════════════════════════════════════
// 机构评级/研报
// ════════════════════════════════════════

type AnalystRating struct {
	Date        string  `json:"date"`
	Code        string  `json:"code"`
	Name        string  `json:"name"`
	Broker      string  `json:"broker"`
	Analyst     string  `json:"analyst"`
	Rating      string  `json:"rating"`
	TargetPrice float64 `json:"target_price"`
	Title       string  `json:"title"`
}

func GetAnalystRatings(code string, limit int) ([]AnalystRating, error) {
	if limit <= 0 || limit > 30 {
		limit = 10
	}
	filter := fmt.Sprintf(`(SECURITY_CODE="%s")`, code)
	u := fmt.Sprintf(
		"https://datacenter-web.eastmoney.com/api/data/v1/get?sortColumns=REPORT_DATE&sortTypes=-1&pageSize=%d&pageNumber=1&reportName=RPT_RATINGCHANGE_DET&columns=ALL&filter=%s",
		limit, url.QueryEscape(filter),
	)
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

	var ratings []AnalystRating
	for _, d := range raw.Result.Data {
		ratings = append(ratings, AnalystRating{
			Date:        ToStr(d["REPORT_DATE"]),
			Code:        ToStr(d["SECURITY_CODE"]),
			Name:        ToStr(d["SECURITY_NAME_ABBR"]),
			Broker:      ToStr(d["ORG_NAME"]),
			Analyst:     ToStr(d["RESEARCHER"]),
			Rating:      ToStr(d["RATING_NAME"]),
			TargetPrice: ToFloat(d["AIM_PRICE"]),
			Title:       ToStr(d["TITLE"]),
		})
	}
	return ratings, nil
}
