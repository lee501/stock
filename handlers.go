package main

import (
	"context"
	"encoding/json"

	"github.com/mark3labs/mcp-go/mcp"

	em "stock-mcp/eastmoney"
)

// ── 基础行情 ──

func handleSearchStock(_ context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	results, err := em.SearchStock(mustStr(req, "keyword"))
	if err != nil {
		return errResult(err), nil
	}
	return jsonResult(results)
}

func handleGetQuote(_ context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	quote, err := em.GetQuote(mustStr(req, "code"))
	if err != nil {
		return errResult(err), nil
	}
	return jsonResult(quote)
}

func handleGetKLine(_ context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	klines, err := em.GetKLine(mustStr(req, "code"), optStr(req, "period", "daily"), optInt(req, "limit", 30))
	if err != nil {
		return errResult(err), nil
	}
	return jsonResult(klines)
}

func handleIndex(_ context.Context, _ mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	indexes, err := em.GetIndexQuotes()
	if err != nil {
		return errResult(err), nil
	}
	return jsonResult(indexes)
}

// ── 资金流向 ──

func handleMoneyFlow(_ context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	flow, err := em.GetMoneyFlow(mustStr(req, "code"))
	if err != nil {
		return errResult(err), nil
	}
	return jsonResult(flow)
}

func handleMoneyFlowHistory(_ context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	days, err := em.GetMoneyFlowHistory(mustStr(req, "code"), optInt(req, "limit", 10))
	if err != nil {
		return errResult(err), nil
	}
	return jsonResult(days)
}

func handleNorthFlow(_ context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	flows, err := em.GetNorthFlow(optInt(req, "limit", 10))
	if err != nil {
		return errResult(err), nil
	}
	return jsonResult(flows)
}

func handleNorthHoldings(_ context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	holdings, err := em.GetNorthHoldings(optStr(req, "market", "sh"), optInt(req, "limit", 20))
	if err != nil {
		return errResult(err), nil
	}
	return jsonResult(holdings)
}

func handleMarginData(_ context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	data, err := em.GetMarginData(optInt(req, "limit", 10))
	if err != nil {
		return errResult(err), nil
	}
	return jsonResult(data)
}

func handleStockMargin(_ context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	data, err := em.GetStockMargin(mustStr(req, "code"), optInt(req, "limit", 10))
	if err != nil {
		return errResult(err), nil
	}
	return jsonResult(data)
}

// ── 基本面 ──

func handleFinancial(_ context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	fin, err := em.GetFinancial(mustStr(req, "code"))
	if err != nil {
		return errResult(err), nil
	}
	return jsonResult(fin)
}

func handleTopHolders(_ context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	holders, err := em.GetTopHolders(mustStr(req, "code"))
	if err != nil {
		return errResult(err), nil
	}
	return jsonResult(holders)
}

func handleDividend(_ context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	divs, err := em.GetDividendHistory(mustStr(req, "code"), optInt(req, "limit", 10))
	if err != nil {
		return errResult(err), nil
	}
	return jsonResult(divs)
}

func handleAnalystRatings(_ context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	ratings, err := em.GetAnalystRatings(mustStr(req, "code"), optInt(req, "limit", 10))
	if err != nil {
		return errResult(err), nil
	}
	return jsonResult(ratings)
}

// ── 板块 ──

func handleSectors(_ context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	sectors, err := em.GetSectors(optStr(req, "type", "industry"), optInt(req, "limit", 20))
	if err != nil {
		return errResult(err), nil
	}
	return jsonResult(sectors)
}

func handleSectorStocks(_ context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	stocks, err := em.GetSectorStocks(mustStr(req, "sector_code"), optInt(req, "limit", 20))
	if err != nil {
		return errResult(err), nil
	}
	return jsonResult(stocks)
}

// ── 排行/异动 ──

func handleRanking(_ context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	stocks, err := em.GetRanking(optStr(req, "type", "top"), optInt(req, "limit", 20))
	if err != nil {
		return errResult(err), nil
	}
	return jsonResult(stocks)
}

func handleLimitStocks(_ context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	stocks, err := em.GetLimitStocks(optStr(req, "type", "up"), optInt(req, "limit", 30))
	if err != nil {
		return errResult(err), nil
	}
	return jsonResult(stocks)
}

func handleDragonTiger(_ context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	items, err := em.GetDragonTiger(optInt(req, "limit", 20))
	if err != nil {
		return errResult(err), nil
	}
	return jsonResult(items)
}

// ── 杠杆/风控 ──

func handleBlockTrades(_ context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	trades, err := em.GetBlockTrades(optInt(req, "limit", 20))
	if err != nil {
		return errResult(err), nil
	}
	return jsonResult(trades)
}

func handleLockupExpiry(_ context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	items, err := em.GetLockupExpiry(optInt(req, "limit", 20))
	if err != nil {
		return errResult(err), nil
	}
	return jsonResult(items)
}

// ── 资讯 ──

func handleStockNews(_ context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	news, err := em.GetStockNews(mustStr(req, "code"), optInt(req, "limit", 10))
	if err != nil {
		return errResult(err), nil
	}
	return jsonResult(news)
}

// ── 辅助函数 ──

func jsonResult(v any) (*mcp.CallToolResult, error) {
	data, _ := json.MarshalIndent(v, "", "  ")
	return mcp.NewToolResultText(string(data)), nil
}

func errResult(err error) *mcp.CallToolResult {
	return mcp.NewToolResultError(err.Error())
}

func mustStr(req mcp.CallToolRequest, key string) string {
	s, _ := req.Params.Arguments[key].(string)
	return s
}

func optStr(req mcp.CallToolRequest, key, fallback string) string {
	if v, ok := req.Params.Arguments[key].(string); ok && v != "" {
		return v
	}
	return fallback
}

func optInt(req mcp.CallToolRequest, key string, fallback int) int {
	if v, ok := req.Params.Arguments[key].(float64); ok && v > 0 {
		return int(v)
	}
	return fallback
}
