package eastmoney

import (
	"fmt"
	"strings"
)

// ════════════════════════════════════════
// 个股资金流向(当日)
// ════════════════════════════════════════

type MoneyFlow struct {
	Code       string  `json:"code"`
	Name       string  `json:"name"`
	MainIn     float64 `json:"main_in"`
	MainOut    float64 `json:"main_out"`
	MainNet    float64 `json:"main_net"`
	MainNetPct float64 `json:"main_net_pct"`
	SuperIn    float64 `json:"super_in"`
	SuperOut   float64 `json:"super_out"`
	BigIn      float64 `json:"big_in"`
	BigOut     float64 `json:"big_out"`
	MidIn      float64 `json:"mid_in"`
	MidOut     float64 `json:"mid_out"`
	SmallIn    float64 `json:"small_in"`
	SmallOut   float64 `json:"small_out"`
}

func GetMoneyFlow(code string) (*MoneyFlow, error) {
	m, err := Push2StockGet("stock/fflow/current", ToSecID(code),
		"f1,f2,f3,f62,f63,f64,f65,f66,f67,f68,f69,f70,f71,f72,f73,f74,f75,f76,f77,f78,f79,f80,f81,f82,f83,f84,f85,f86")
	if err != nil {
		return nil, err
	}
	if m == nil {
		return nil, fmt.Errorf("no money flow data for %s", code)
	}
	return &MoneyFlow{
		Code:       GetStr(m, "f57"),
		Name:       GetStr(m, "f58"),
		MainIn:     GetFloat(m, "f62"),
		MainOut:    GetFloat(m, "f63"),
		MainNet:    GetFloat(m, "f64"),
		MainNetPct: GetFloat(m, "f65") / 100,
		SuperIn:    GetFloat(m, "f66"),
		SuperOut:   GetFloat(m, "f67"),
		BigIn:      GetFloat(m, "f72"),
		BigOut:     GetFloat(m, "f73"),
		MidIn:      GetFloat(m, "f78"),
		MidOut:     GetFloat(m, "f79"),
		SmallIn:    GetFloat(m, "f84"),
		SmallOut:   GetFloat(m, "f85"),
	}, nil
}

// ════════════════════════════════════════
// 个股历史资金流向
// ════════════════════════════════════════

type MoneyFlowDay struct {
	Date     string  `json:"date"`
	MainNet  float64 `json:"main_net"`
	SuperNet float64 `json:"super_net"`
	BigNet   float64 `json:"big_net"`
	MidNet   float64 `json:"mid_net"`
	SmallNet float64 `json:"small_net"`
}

func GetMoneyFlowHistory(code string, limit int) ([]MoneyFlowDay, error) {
	if limit <= 0 || limit > 60 {
		limit = 10
	}
	lines, err := Push2HisGet(Push2HisQuery{
		Path:  "stock/fflow/daykline/get",
		SecID: ToSecID(code),
		Params: fmt.Sprintf("lmt=%d&klt=101&fields1=f1,f2,f3,f7&fields2=f51,f52,f53,f54,f55,f56,f57,f58,f59,f60,f61,f62,f63,f64,f65",
			limit),
	})
	if err != nil {
		return nil, err
	}

	var days []MoneyFlowDay
	for _, line := range lines {
		parts := strings.Split(line, ",")
		if len(parts) < 6 {
			continue
		}
		days = append(days, MoneyFlowDay{
			Date:     parts[0],
			MainNet:  ParseFloat(parts[1]),
			SmallNet: ParseFloat(parts[2]),
			MidNet:   ParseFloat(parts[3]),
			BigNet:   ParseFloat(parts[4]),
			SuperNet: ParseFloat(parts[5]),
		})
	}
	return days, nil
}

// ════════════════════════════════════════
// 北向资金(沪深股通)每日净买入
// ════════════════════════════════════════

type NorthFlow struct {
	Date      string  `json:"date"`
	HKToSH    float64 `json:"hk_to_sh"`
	HKToSZ    float64 `json:"hk_to_sz"`
	Total     float64 `json:"total"`
	HKToSHAcc float64 `json:"hk_to_sh_acc"`
	HKToSZAcc float64 `json:"hk_to_sz_acc"`
}

func GetNorthFlow(limit int) ([]NorthFlow, error) {
	if limit <= 0 || limit > 60 {
		limit = 10
	}
	data, err := DatacenterGet(DatacenterQuery{
		ReportName:  "RPT_MUTUAL_DEAL_HISTORY",
		SortColumns: "TRADE_DATE",
		SortTypes:   "-1",
		PageSize:    limit,
	})
	if err != nil {
		return nil, err
	}

	var flows []NorthFlow
	for _, d := range data {
		sh := ToFloat(d["MUTUAL_A_DEAL_FIN"])
		sz := ToFloat(d["MUTUAL_D_DEAL_FIN"])
		flows = append(flows, NorthFlow{
			Date:      ToStr(d["TRADE_DATE"]),
			HKToSH:    sh,
			HKToSZ:    sz,
			Total:     sh + sz,
			HKToSHAcc: ToFloat(d["MUTUAL_A_ACCUM_DEAL"]),
			HKToSZAcc: ToFloat(d["MUTUAL_D_ACCUM_DEAL"]),
		})
	}
	return flows, nil
}

// ════════════════════════════════════════
// 北向资金个股持仓排行
// ════════════════════════════════════════

type NorthHolding struct {
	Code       string  `json:"code"`
	Name       string  `json:"name"`
	HoldShares float64 `json:"hold_shares"`
	HoldRatio  float64 `json:"hold_ratio"`
	HoldValue  float64 `json:"hold_value"`
	ChangePct  float64 `json:"change_pct"`
}

func GetNorthHoldings(market string, limit int) ([]NorthHolding, error) {
	if limit <= 0 || limit > 50 {
		limit = 20
	}
	filter := `(MUTUAL_TYPE="001")`
	if market == "sz" {
		filter = `(MUTUAL_TYPE="003")`
	}
	data, err := DatacenterGet(DatacenterQuery{
		ReportName:  "RPT_MUTUAL_HOLDSTOCKNORTH_STA",
		SortColumns: "HOLD_MARKET_CAP",
		SortTypes:   "-1",
		PageSize:    limit,
		Filter:      filter,
		Extra:       "source=WEB&client=WEB",
	})
	if err != nil {
		return nil, err
	}

	var holdings []NorthHolding
	for _, d := range data {
		holdings = append(holdings, NorthHolding{
			Code:       ToStr(d["SECURITY_CODE"]),
			Name:       ToStr(d["SECURITY_NAME"]),
			HoldShares: ToFloat(d["HOLD_SHARES"]),
			HoldRatio:  ToFloat(d["FREE_SHARES_RATIO"]),
			HoldValue:  ToFloat(d["HOLD_MARKET_CAP"]),
			ChangePct:  ToFloat(d["CHANGE_RATE"]),
		})
	}
	return holdings, nil
}

// ════════════════════════════════════════
// 两市融资融券余额汇总
// ════════════════════════════════════════

type MarginData struct {
	Date       string  `json:"date"`
	FinBuy     float64 `json:"fin_buy"`
	FinBalance float64 `json:"fin_balance"`
	SecSell    float64 `json:"sec_sell"`
	SecBalance float64 `json:"sec_balance"`
	Total      float64 `json:"total"`
}

func GetMarginData(limit int) ([]MarginData, error) {
	if limit <= 0 || limit > 60 {
		limit = 10
	}
	data, err := DatacenterGet(DatacenterQuery{
		ReportName:  "RPTA_WEB_MARGIN_DAILYTRADE",
		Columns:     "STATISTICS_DATE,FIN_BUY_AMT,FIN_BALANCE,LOAN_SELL_AMT,LOAN_BALANCE,MARGIN_BALANCE",
		SortColumns: "STATISTICS_DATE",
		SortTypes:   "-1",
		PageSize:    limit,
	})
	if err != nil {
		return nil, err
	}

	var items []MarginData
	for _, d := range data {
		items = append(items, MarginData{
			Date:       ToStr(d["STATISTICS_DATE"]),
			FinBuy:     ToFloat(d["FIN_BUY_AMT"]),
			FinBalance: ToFloat(d["FIN_BALANCE"]),
			SecSell:    ToFloat(d["LOAN_SELL_AMT"]),
			SecBalance: ToFloat(d["LOAN_BALANCE"]),
			Total:      ToFloat(d["MARGIN_BALANCE"]),
		})
	}
	return items, nil
}

// ════════════════════════════════════════
// 个股融资融券
// ════════════════════════════════════════

type StockMargin struct {
	Date       string  `json:"date"`
	Code       string  `json:"code"`
	Name       string  `json:"name"`
	FinBuy     float64 `json:"fin_buy"`
	FinBalance float64 `json:"fin_balance"`
	SecSell    float64 `json:"sec_sell"`
	SecBalance float64 `json:"sec_balance"`
}

func GetStockMargin(code string, limit int) ([]StockMargin, error) {
	if limit <= 0 || limit > 30 {
		limit = 10
	}
	data, err := DatacenterGet(DatacenterQuery{
		ReportName:  "RPTA_WEB_RZRQ_GGMX",
		SortColumns: "DATE",
		SortTypes:   "-1",
		PageSize:    limit,
		Filter:      fmt.Sprintf(`(SCODE="%s")`, code),
	})
	if err != nil {
		return nil, err
	}

	var items []StockMargin
	for _, d := range data {
		items = append(items, StockMargin{
			Date:       ToStr(d["DATE"]),
			Code:       ToStr(d["SCODE"]),
			Name:       ToStr(d["SECNAME"]),
			FinBuy:     ToFloat(d["RZJME"]),
			FinBalance: ToFloat(d["RZYE"]),
			SecSell:    ToFloat(d["RQMCL"]),
			SecBalance: ToFloat(d["RQYE"]),
		})
	}
	return items, nil
}
