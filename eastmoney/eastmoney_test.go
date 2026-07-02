package eastmoney

import (
	"testing"
	"time"
)

const testStock = "000001" // 平安银行

// throttle avoids push2 rate-limiting between tests.
func throttle() { time.Sleep(350 * time.Millisecond) }

func skipIfPush2Unavailable(t *testing.T, name string, err error) {
	t.Helper()
	if err != nil {
		t.Skipf("%s: %v (push2 may be unavailable outside trading hours)", name, err)
	}
}

// ── 基础行情 ──

func TestSearchStock(t *testing.T) {
	throttle()
	results, err := SearchStock("平安")
	if err != nil {
		t.Fatalf("SearchStock: %v", err)
	}
	if len(results) == 0 {
		t.Fatal("SearchStock returned empty")
	}
	for _, r := range results {
		if r.Code == "" || r.Name == "" {
			t.Errorf("empty code or name: %+v", r)
		}
	}
}

func TestGetQuote(t *testing.T) {
	throttle()
	q, err := GetQuote(testStock)
	skipIfPush2Unavailable(t, "GetQuote", err)
	if q.Code == "" && q.Name == "" {
		t.Skip("empty response (push2 unavailable)")
	}
	if q.Code != testStock {
		t.Errorf("expected code %s, got %s", testStock, q.Code)
	}
}

func TestGetKLine(t *testing.T) {
	throttle()
	klines, err := GetKLine(testStock, "daily", 5)
	skipIfPush2Unavailable(t, "GetKLine", err)
	if len(klines) == 0 {
		t.Skip("empty (push2his unavailable)")
	}
	for _, k := range klines {
		if k.Date == "" {
			t.Error("empty date")
		}
	}
}

func TestGetKLineWeekly(t *testing.T) {
	throttle()
	klines, err := GetKLine(testStock, "weekly", 3)
	skipIfPush2Unavailable(t, "GetKLine weekly", err)
	if len(klines) == 0 {
		t.Skip("empty")
	}
}

func TestGetIndexQuotes(t *testing.T) {
	throttle()
	indexes, err := GetIndexQuotes()
	skipIfPush2Unavailable(t, "GetIndexQuotes", err)
	if len(indexes) == 0 {
		t.Skip("empty (push2 unavailable)")
	}
	for _, idx := range indexes {
		if idx.Code == "" {
			t.Error("empty code")
		}
	}
}

// ── 资金流向 ──

func TestGetMoneyFlow(t *testing.T) {
	throttle()
	flow, err := GetMoneyFlow(testStock)
	skipIfPush2Unavailable(t, "GetMoneyFlow", err)
	if flow.Code == "" {
		t.Skip("empty response (push2 unavailable)")
	}
	if flow.Code != testStock {
		t.Errorf("expected code %s, got %s", testStock, flow.Code)
	}
}

func TestGetMoneyFlowHistory(t *testing.T) {
	throttle()
	days, err := GetMoneyFlowHistory(testStock, 5)
	skipIfPush2Unavailable(t, "GetMoneyFlowHistory", err)
	if len(days) == 0 {
		t.Skip("empty (push2his unavailable)")
	}
	if days[0].Date == "" {
		t.Error("empty date")
	}
}

func TestGetNorthFlow(t *testing.T) {
	throttle()
	flows, err := GetNorthFlow(3)
	if err != nil {
		t.Fatalf("GetNorthFlow: %v", err)
	}
	if len(flows) == 0 {
		t.Fatal("returned empty")
	}
	if flows[0].Date == "" {
		t.Error("empty date")
	}
}

func TestGetNorthHoldings(t *testing.T) {
	throttle()
	holdings, err := GetNorthHoldings("sh", 3)
	if err != nil {
		t.Fatalf("GetNorthHoldings: %v", err)
	}
	if len(holdings) == 0 {
		t.Skip("returned empty (server may be busy)")
	}
	for _, h := range holdings {
		if h.Code == "" || h.Name == "" {
			t.Errorf("empty code or name: %+v", h)
		}
	}
}

func TestGetMarginData(t *testing.T) {
	throttle()
	data, err := GetMarginData(3)
	if err != nil {
		t.Fatalf("GetMarginData: %v", err)
	}
	if len(data) == 0 {
		t.Fatal("returned empty")
	}
	if data[0].Total <= 0 {
		t.Errorf("invalid total %.2f", data[0].Total)
	}
}

func TestGetStockMargin(t *testing.T) {
	throttle()
	data, err := GetStockMargin(testStock, 3)
	if err != nil {
		t.Fatalf("GetStockMargin: %v", err)
	}
	if len(data) == 0 {
		t.Fatal("returned empty")
	}
	if data[0].Code == "" {
		t.Error("empty code")
	}
}

// ── 基本面 ──

func TestGetFinancial(t *testing.T) {
	throttle()
	fin, err := GetFinancial(testStock)
	skipIfPush2Unavailable(t, "GetFinancial", err)
	if fin.Code == "" {
		t.Skip("empty response (push2 unavailable)")
	}
	if fin.Code != testStock {
		t.Errorf("expected code %s, got %s", testStock, fin.Code)
	}
}

func TestGetTopHolders(t *testing.T) {
	throttle()
	holders, err := GetTopHolders(testStock)
	if err != nil {
		t.Fatalf("GetTopHolders: %v", err)
	}
	if len(holders) == 0 {
		t.Fatal("returned empty")
	}
	if holders[0].Name == "" {
		t.Error("empty holder name")
	}
}

func TestGetDividendHistory(t *testing.T) {
	throttle()
	divs, err := GetDividendHistory(testStock, 3)
	if err != nil {
		t.Fatalf("GetDividendHistory: %v", err)
	}
	if len(divs) == 0 {
		t.Fatal("returned empty")
	}
}

func TestGetAnalystRatings(t *testing.T) {
	throttle()
	ratings, err := GetAnalystRatings(testStock, 3)
	if err != nil {
		t.Fatalf("GetAnalystRatings: %v", err)
	}
	if len(ratings) == 0 {
		t.Fatal("returned empty")
	}
	if ratings[0].Broker == "" {
		t.Error("empty broker")
	}
}

// ── 板块 ──

func TestGetSectorsIndustry(t *testing.T) {
	throttle()
	sectors, err := GetSectors("industry", 5)
	skipIfPush2Unavailable(t, "GetSectors industry", err)
	if len(sectors) == 0 {
		t.Skip("empty (push2 unavailable)")
	}
	for _, s := range sectors {
		if s.Code == "" || s.Name == "" {
			t.Errorf("empty code or name: %+v", s)
		}
	}
}

func TestGetSectorsConcept(t *testing.T) {
	throttle()
	sectors, err := GetSectors("concept", 5)
	skipIfPush2Unavailable(t, "GetSectors concept", err)
	if len(sectors) == 0 {
		t.Skip("empty (push2 unavailable)")
	}
}

func TestGetSectorStocks(t *testing.T) {
	throttle()
	sectors, err := GetSectors("industry", 1)
	if err != nil || len(sectors) == 0 {
		t.Skip("no sectors available (push2 unavailable)")
	}
	throttle()
	stocks, err := GetSectorStocks(sectors[0].Code, 5)
	if err != nil {
		t.Fatalf("GetSectorStocks: %v", err)
	}
	if len(stocks) == 0 {
		t.Skip("empty (push2 unavailable)")
	}
	for _, s := range stocks {
		if s.Code == "" || s.Name == "" {
			t.Errorf("empty code or name: %+v", s)
		}
	}
}

// ── 排行/异动 ──

func TestGetRankingTop(t *testing.T) {
	throttle()
	stocks, err := GetRanking("top", 5)
	skipIfPush2Unavailable(t, "GetRanking top", err)
	if len(stocks) == 0 {
		t.Skip("empty (push2 unavailable)")
	}
}

func TestGetRankingBottom(t *testing.T) {
	throttle()
	stocks, err := GetRanking("bottom", 5)
	skipIfPush2Unavailable(t, "GetRanking bottom", err)
	if len(stocks) == 0 {
		t.Skip("empty (push2 unavailable)")
	}
}

func TestGetRankingVolume(t *testing.T) {
	throttle()
	stocks, err := GetRanking("volume", 5)
	skipIfPush2Unavailable(t, "GetRanking volume", err)
	if len(stocks) == 0 {
		t.Skip("empty (push2 unavailable)")
	}
}

func TestGetLimitStocksUp(t *testing.T) {
	throttle()
	stocks, err := GetLimitStocks("up", 5)
	if err != nil {
		t.Fatalf("GetLimitStocks up: %v", err)
	}
	if len(stocks) == 0 {
		t.Skip("no limit-up stocks (market may be closed)")
	}
	for _, s := range stocks {
		if s.Code == "" || s.Name == "" {
			t.Errorf("empty code or name: %+v", s)
		}
		if s.Price <= 0 {
			t.Errorf("invalid price for %s: %.2f", s.Code, s.Price)
		}
	}
}

func TestGetLimitStocksDown(t *testing.T) {
	throttle()
	stocks, err := GetLimitStocks("down", 5)
	if err != nil {
		t.Fatalf("GetLimitStocks down: %v", err)
	}
	for _, s := range stocks {
		if s.Code == "" || s.Name == "" {
			t.Errorf("empty code or name: %+v", s)
		}
	}
}

func TestGetDragonTiger(t *testing.T) {
	throttle()
	items, err := GetDragonTiger(5)
	if err != nil {
		t.Fatalf("GetDragonTiger: %v", err)
	}
	if len(items) == 0 {
		t.Fatal("returned empty")
	}
	if items[0].Date == "" {
		t.Error("empty date")
	}
	if items[0].Code == "" {
		t.Error("empty code")
	}
}

// ── 杠杆/风控 ──

func TestGetBlockTrades(t *testing.T) {
	throttle()
	trades, err := GetBlockTrades(5)
	if err != nil {
		t.Fatalf("GetBlockTrades: %v", err)
	}
	if len(trades) == 0 {
		t.Fatal("returned empty")
	}
	if trades[0].Code == "" || trades[0].Name == "" {
		t.Error("empty code or name")
	}
}

func TestGetLockupExpiry(t *testing.T) {
	throttle()
	items, err := GetLockupExpiry(5)
	if err != nil {
		t.Fatalf("GetLockupExpiry: %v", err)
	}
	if len(items) == 0 {
		t.Fatal("returned empty")
	}
	if items[0].Code == "" {
		t.Error("empty code")
	}
}

// ── 资讯 ──

func TestGetStockNews(t *testing.T) {
	throttle()
	news, err := GetStockNews(testStock, 3)
	if err != nil {
		t.Fatalf("GetStockNews: %v", err)
	}
	if len(news) == 0 {
		t.Fatal("returned empty")
	}
	if news[0].Title == "" {
		t.Error("empty title")
	}
}

// ── 辅助函数 ──

func TestToSecID(t *testing.T) {
	tests := []struct {
		code string
		want string
	}{
		{"600000", "1.600000"},
		{"000001", "0.000001"},
		{"300750", "0.300750"},
		{"688981", "1.688981"},
		{"900901", "1.900901"},
	}
	for _, tt := range tests {
		got := ToSecID(tt.code)
		if got != tt.want {
			t.Errorf("ToSecID(%s) = %s, want %s", tt.code, got, tt.want)
		}
	}
}
