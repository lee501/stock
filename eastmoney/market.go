package eastmoney

import (
	"encoding/json"
	"fmt"
	"time"
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

	diff, err := ClistGet(ClistQuery{
		FS:     fs,
		Fields: "f2,f3,f4,f12,f14,f128,f136,f140,f141,f62",
		Limit:  limit,
	})
	if err != nil {
		return nil, err
	}

	var sectors []Sector
	for _, m := range diff {
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

	diff, err := ClistGet(ClistQuery{
		FS:     fmt.Sprintf("b:%s+f:!50", sectorCode),
		Fields: "f2,f3,f4,f5,f6,f12,f14,f62",
		Limit:  limit,
	})
	if err != nil {
		return nil, err
	}

	var stocks []SectorStock
	for _, m := range diff {
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

	diff, err := ClistGet(ClistQuery{
		FS:        "m:0+t:6,m:0+t:80,m:1+t:2,m:1+t:23+f:!50",
		Fields:    "f2,f3,f4,f5,f6,f8,f9,f12,f14,f20",
		Limit:     limit,
		SortField: sortField,
		SortOrder: po,
	})
	if err != nil {
		return nil, err
	}

	var stocks []RankStock
	for _, m := range diff {
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
	if limit <= 0 || limit > 200 {
		limit = 30
	}
	date := time.Now().Format("20060102")

	endpoint := "getTopicZTPool"
	sort := "fbt:asc"
	if limitType == "down" {
		endpoint = "getTopicDTPool"
		sort = "fund:asc"
	}

	u := fmt.Sprintf("%s/%s?ut=7eea3edcaed734bea9cbfc24409ed989&dpt=wz.ztzt&Pageindex=0&pagesize=%d&sort=%s&date=%s",
		basePush2Ex, endpoint, limit, sort, date)
	body, err := DoGet(u)
	if err != nil {
		return nil, err
	}

	var raw struct {
		Data struct {
			Pool []map[string]json.RawMessage `json:"pool"`
		} `json:"data"`
	}
	if err := json.Unmarshal(body, &raw); err != nil {
		return nil, err
	}

	var stocks []LimitStock
	for _, m := range raw.Data.Pool {
		days := int(GetFloat(m, "lbc"))
		if limitType == "down" {
			days = int(GetFloat(m, "days"))
		}
		var zttj struct {
			Days int `json:"days"`
			Ct   int `json:"ct"`
		}
		if v, ok := m["zttj"]; ok {
			json.Unmarshal(v, &zttj)
			if zttj.Days > 0 {
				days = zttj.Days
			}
		}

		stocks = append(stocks, LimitStock{
			Code:       GetStr(m, "c"),
			Name:       GetStr(m, "n"),
			Price:      GetFloat(m, "p") / 1000,
			ChangePct:  GetFloat(m, "zdp"),
			TurnOver:   GetFloat(m, "hs"),
			Amount:     GetFloat(m, "amount"),
			LimitTimes: days,
			FirstTime:  fmtTime(GetFloat(m, "fbt")),
			LastTime:   fmtTime(GetFloat(m, "lbt")),
			OpenCount:  int(GetFloat(m, "zbc")),
			Theme:      GetStr(m, "hybk"),
		})
	}
	return stocks, nil
}

func fmtTime(v float64) string {
	if v <= 0 {
		return ""
	}
	t := int(v)
	return fmt.Sprintf("%02d:%02d:%02d", t/10000, (t/100)%100, t%100)
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
	data, err := DatacenterGet(DatacenterQuery{
		ReportName:  "RPT_DAILYBILLBOARD_DETAILSNEW",
		SortColumns: "TRADE_DATE,SECURITY_CODE",
		SortTypes:   "-1,1",
		PageSize:    limit,
	})
	if err != nil {
		return nil, err
	}

	var items []DragonTiger
	for _, d := range data {
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
	data, err := DatacenterGet(DatacenterQuery{
		ReportName:  "RPT_DATA_BLOCKTRADE",
		Columns:     "TRADE_DATE,SECURITY_CODE,SECURITY_NAME_ABBR,DEAL_PRICE,DEAL_VOLUME,DEAL_AMT,PREMIUM_RATIO,BUYER_NAME,SELLER_NAME",
		SortColumns: "TRADE_DATE",
		SortTypes:   "-1",
		PageSize:    limit,
	})
	if err != nil {
		return nil, err
	}

	var trades []BlockTrade
	for _, d := range data {
		trades = append(trades, BlockTrade{
			Date:    ToStr(d["TRADE_DATE"]),
			Code:    ToStr(d["SECURITY_CODE"]),
			Name:    ToStr(d["SECURITY_NAME_ABBR"]),
			Price:   ToFloat(d["DEAL_PRICE"]),
			Volume:  ToFloat(d["DEAL_VOLUME"]),
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
	data, err := DatacenterGet(DatacenterQuery{
		ReportName:  "RPT_LIFT_STAGE",
		Columns:     "FREE_DATE,SECURITY_CODE,SECURITY_NAME_ABBR,CURRENT_FREE_SHARES,LIFT_MARKET_CAP,FREE_RATIO,FREE_SHARES_TYPE",
		SortColumns: "FREE_DATE,CURRENT_FREE_SHARES",
		SortTypes:   "1,1",
		PageSize:    limit,
	})
	if err != nil {
		return nil, err
	}

	var items []LockupExpiry
	for _, d := range data {
		items = append(items, LockupExpiry{
			Date:         ToStr(d["FREE_DATE"]),
			Code:         ToStr(d["SECURITY_CODE"]),
			Name:         ToStr(d["SECURITY_NAME_ABBR"]),
			UnlockShares: ToFloat(d["CURRENT_FREE_SHARES"]),
			UnlockValue:  ToFloat(d["LIFT_MARKET_CAP"]),
			UnlockRatio:  ToFloat(d["FREE_RATIO"]),
			LockupType:   ToStr(d["FREE_SHARES_TYPE"]),
		})
	}
	return items, nil
}
