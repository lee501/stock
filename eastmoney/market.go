package eastmoney

import (
	"encoding/json"
	"fmt"
	"net/url"
)

// ════════════════════════════════════════
// 板块行情排行
// ════════════════════════════════════════

type Sector struct {
	Code       string  `json:"code"`
	Name       string  `json:"name"`
	ChangePct  float64 `json:"change_pct"`
	MainNet    float64 `json:"main_net"`
	LeaderCode string  `json:"leader_code"`
	LeaderName string  `json:"leader_name"`
	LeaderPct  float64 `json:"leader_pct"`
}

func GetSectors(sectorType string, limit int) ([]Sector, error) {
	fs := "m:90+t:2"
	if sectorType == "concept" {
		fs = "m:90+t:3"
	}
	if limit <= 0 || limit > 50 {
		limit = 20
	}

	u := fmt.Sprintf(
		"https://push2.eastmoney.com/api/qt/clist/get?pn=1&pz=%d&po=1&np=1&fltt=2&invt=2&fs=%s&fields=f2,f3,f4,f12,f14,f128,f136,f140,f141,f62",
		limit, url.QueryEscape(fs),
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

	var sectors []Sector
	for _, m := range raw.Data.Diff {
		sectors = append(sectors, Sector{
			Code:       GetStr(m, "f12"),
			Name:       GetStr(m, "f14"),
			ChangePct:  GetFloat(m, "f3"),
			MainNet:    GetFloat(m, "f62"),
			LeaderCode: GetStr(m, "f140"),
			LeaderName: GetStr(m, "f128"),
			LeaderPct:  GetFloat(m, "f136"),
		})
	}
	return sectors, nil
}

// ════════════════════════════════════════
// 板块成分股
// ════════════════════════════════════════

type SectorStock struct {
	Code      string  `json:"code"`
	Name      string  `json:"name"`
	Price     float64 `json:"price"`
	ChangePct float64 `json:"change_pct"`
	Volume    float64 `json:"volume"`
	Amount    float64 `json:"amount"`
	MainNet   float64 `json:"main_net"`
}

func GetSectorStocks(sectorCode string, limit int) ([]SectorStock, error) {
	if limit <= 0 || limit > 50 {
		limit = 20
	}
	fs := fmt.Sprintf("b:%s+f:!50", sectorCode)

	u := fmt.Sprintf(
		"https://push2.eastmoney.com/api/qt/clist/get?pn=1&pz=%d&po=1&np=1&fltt=2&invt=2&fs=%s&fields=f2,f3,f4,f5,f6,f12,f14,f62",
		limit, url.QueryEscape(fs),
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

	var stocks []SectorStock
	for _, m := range raw.Data.Diff {
		stocks = append(stocks, SectorStock{
			Code:      GetStr(m, "f12"),
			Name:      GetStr(m, "f14"),
			Price:     GetFloat(m, "f2"),
			ChangePct: GetFloat(m, "f3"),
			Volume:    GetFloat(m, "f5"),
			Amount:    GetFloat(m, "f6"),
			MainNet:   GetFloat(m, "f62"),
		})
	}
	return stocks, nil
}

// ════════════════════════════════════════
// 涨跌排行榜
// ════════════════════════════════════════

type RankStock struct {
	Code      string  `json:"code"`
	Name      string  `json:"name"`
	Price     float64 `json:"price"`
	ChangePct float64 `json:"change_pct"`
	TurnOver  float64 `json:"turnover"`
	Amount    float64 `json:"amount"`
	PE        float64 `json:"pe"`
	MarketCap float64 `json:"market_cap"`
}

func GetRanking(rankType string, limit int) ([]RankStock, error) {
	if limit <= 0 || limit > 50 {
		limit = 20
	}
	po, sortField := "1", "f3"
	switch rankType {
	case "bottom":
		po = "0"
	case "volume":
		sortField = "f5"
	case "amount":
		sortField = "f6"
	case "turnover":
		sortField = "f8"
	}

	fs := "m:0+t:6,m:0+t:80,m:1+t:2,m:1+t:23+f:!50"
	u := fmt.Sprintf(
		"https://push2.eastmoney.com/api/qt/clist/get?pn=1&pz=%d&po=%s&np=1&fltt=2&invt=2&fid=%s&fs=%s&fields=f2,f3,f4,f5,f6,f8,f9,f12,f14,f20",
		limit, po, sortField, url.QueryEscape(fs),
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

	var stocks []RankStock
	for _, m := range raw.Data.Diff {
		stocks = append(stocks, RankStock{
			Code:      GetStr(m, "f12"),
			Name:      GetStr(m, "f14"),
			Price:     GetFloat(m, "f2"),
			ChangePct: GetFloat(m, "f3"),
			TurnOver:  GetFloat(m, "f8"),
			Amount:    GetFloat(m, "f6"),
			PE:        GetFloat(m, "f9"),
			MarketCap: GetFloat(m, "f20"),
		})
	}
	return stocks, nil
}

// ════════════════════════════════════════
// 涨停/跌停分析
// ════════════════════════════════════════

type LimitStock struct {
	Code       string  `json:"code"`
	Name       string  `json:"name"`
	Price      float64 `json:"price"`
	ChangePct  float64 `json:"change_pct"`
	TurnOver   float64 `json:"turnover"`
	Amount     float64 `json:"amount"`
	LimitTimes int     `json:"limit_times"`
	FirstTime  string  `json:"first_time"`
	LastTime   string  `json:"last_time"`
	OpenCount  int     `json:"open_count"`
	Theme      string  `json:"theme"`
}

func GetLimitStocks(limitType string, limit int) ([]LimitStock, error) {
	if limit <= 0 || limit > 50 {
		limit = 30
	}
	reportName := "RPT_LIMITUP_BASICINFO"
	if limitType == "down" {
		reportName = "RPT_LIMITDOWN_BASICINFO"
	}

	u := fmt.Sprintf(
		"https://datacenter-web.eastmoney.com/api/data/v1/get?sortColumns=LATEST_BOARD_TIME&sortTypes=-1&pageSize=%d&pageNumber=1&reportName=%s&columns=ALL",
		limit, reportName,
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

	var stocks []LimitStock
	for _, d := range raw.Result.Data {
		stocks = append(stocks, LimitStock{
			Code:       ToStr(d["SECURITY_CODE"]),
			Name:       ToStr(d["SECURITY_NAME_ABBR"]),
			Price:      ToFloat(d["LATEST"]),
			ChangePct:  ToFloat(d["CHANGE_RATE"]),
			TurnOver:   ToFloat(d["TURNOVERRATE"]),
			Amount:     ToFloat(d["DEAL_AMOUNT"]),
			LimitTimes: int(ToFloat(d["CONTINUOUS_DAYS"])),
			FirstTime:  ToStr(d["FIRST_BOARD_TIME"]),
			LastTime:   ToStr(d["LATEST_BOARD_TIME"]),
			OpenCount:  int(ToFloat(d["OPEN_TIMES"])),
			Theme:      ToStr(d["REASON_TYPE"]),
		})
	}
	return stocks, nil
}

// ════════════════════════════════════════
// 龙虎榜
// ════════════════════════════════════════

type DragonTiger struct {
	Code      string  `json:"code"`
	Name      string  `json:"name"`
	Date      string  `json:"date"`
	ChangePct float64 `json:"change_pct"`
	Close     float64 `json:"close"`
	NetBuy    float64 `json:"net_buy"`
	BuyTotal  float64 `json:"buy_total"`
	SellTotal float64 `json:"sell_total"`
	Reason    string  `json:"reason"`
	TurnOver  float64 `json:"turnover"`
	Amount    float64 `json:"amount"`
}

func GetDragonTiger(limit int) ([]DragonTiger, error) {
	if limit <= 0 || limit > 50 {
		limit = 20
	}
	u := fmt.Sprintf(
		"https://datacenter-web.eastmoney.com/api/data/v1/get?sortColumns=SECURITY_CODE&sortTypes=1&pageSize=%d&pageNumber=1&reportName=RPT_DAILYBILLBOARD_DETAILSNEW&columns=ALL&filter=",
		limit,
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

	var items []DragonTiger
	for _, d := range raw.Result.Data {
		items = append(items, DragonTiger{
			Code:      ToStr(d["SECURITY_CODE"]),
			Name:      ToStr(d["SECURITY_NAME_ABBR"]),
			Date:      ToStr(d["TRADE_DATE"]),
			ChangePct: ToFloat(d["CHANGE_RATE"]),
			Close:     ToFloat(d["CLOSE_PRICE"]),
			NetBuy:    ToFloat(d["NET_BUY_AMT"]),
			BuyTotal:  ToFloat(d["BUY_TOTAL_AMT"]),
			SellTotal: ToFloat(d["SELL_TOTAL_AMT"]),
			Reason:    ToStr(d["EXPLANATION"]),
			TurnOver:  ToFloat(d["TURNOVERRATE"]),
			Amount:    ToFloat(d["DEAL_AMT"]),
		})
	}
	return items, nil
}

// ════════════════════════════════════════
// 大宗交易
// ════════════════════════════════════════

type BlockTrade struct {
	Date    string  `json:"date"`
	Code    string  `json:"code"`
	Name    string  `json:"name"`
	Price   float64 `json:"price"`
	Volume  float64 `json:"volume"`
	Amount  float64 `json:"amount"`
	Premium float64 `json:"premium"`
	Buyer   string  `json:"buyer"`
	Seller  string  `json:"seller"`
}

func GetBlockTrades(limit int) ([]BlockTrade, error) {
	if limit <= 0 || limit > 50 {
		limit = 20
	}
	u := fmt.Sprintf(
		"https://datacenter-web.eastmoney.com/api/data/v1/get?sortColumns=TRADE_DATE&sortTypes=-1&pageSize=%d&pageNumber=1&reportName=RPT_BLOCKTRADE_DETAILSINFO&columns=TRADE_DATE,SECURITY_CODE,SECURITY_NAME_ABBR,DEAL_PRICE,DEAL_VOL,DEAL_AMT,PREMIUM_RATIO,BUYER_NAME,SELLER_NAME",
		limit,
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

	var trades []BlockTrade
	for _, d := range raw.Result.Data {
		trades = append(trades, BlockTrade{
			Date:    ToStr(d["TRADE_DATE"]),
			Code:    ToStr(d["SECURITY_CODE"]),
			Name:    ToStr(d["SECURITY_NAME_ABBR"]),
			Price:   ToFloat(d["DEAL_PRICE"]),
			Volume:  ToFloat(d["DEAL_VOL"]),
			Amount:  ToFloat(d["DEAL_AMT"]),
			Premium: ToFloat(d["PREMIUM_RATIO"]),
			Buyer:   ToStr(d["BUYER_NAME"]),
			Seller:  ToStr(d["SELLER_NAME"]),
		})
	}
	return trades, nil
}

// ════════════════════════════════════════
// 限售解禁日历
// ════════════════════════════════════════

type LockupExpiry struct {
	Date         string  `json:"date"`
	Code         string  `json:"code"`
	Name         string  `json:"name"`
	UnlockShares float64 `json:"unlock_shares"`
	UnlockValue  float64 `json:"unlock_value"`
	UnlockRatio  float64 `json:"unlock_ratio"`
	LockupType   string  `json:"lockup_type"`
}

func GetLockupExpiry(limit int) ([]LockupExpiry, error) {
	if limit <= 0 || limit > 50 {
		limit = 20
	}
	u := fmt.Sprintf(
		"https://datacenter-web.eastmoney.com/api/data/v1/get?sortColumns=FREE_DATE&sortTypes=1&pageSize=%d&pageNumber=1&reportName=RPT_LIFT_STAGE&columns=FREE_DATE,SECURITY_CODE,SECURITY_NAME_ABBR,FREE_SHARES_NUM,FREE_MARKET_CAP,FREE_RATIO,LIFT_TYPE&filter=",
		limit,
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

	var items []LockupExpiry
	for _, d := range raw.Result.Data {
		items = append(items, LockupExpiry{
			Date:         ToStr(d["FREE_DATE"]),
			Code:         ToStr(d["SECURITY_CODE"]),
			Name:         ToStr(d["SECURITY_NAME_ABBR"]),
			UnlockShares: ToFloat(d["FREE_SHARES_NUM"]),
			UnlockValue:  ToFloat(d["FREE_MARKET_CAP"]),
			UnlockRatio:  ToFloat(d["FREE_RATIO"]),
			LockupType:   ToStr(d["LIFT_TYPE"]),
		})
	}
	return items, nil
}
