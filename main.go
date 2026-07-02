package main

import (
	"fmt"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

func main() {
	s := server.NewMCPServer(
		"a-stock",
		"3.0.0",
		server.WithToolCapabilities(false),
	)

	registerTools(s)

	if err := server.ServeStdio(s); err != nil {
		fmt.Printf("server error: %v\n", err)
	}
}

func registerTools(s *server.MCPServer) {

	// ═══════ 基础行情 ═══════

	s.AddTool(mcp.NewTool("search_stock",
		mcp.WithDescription("根据关键词搜索A股股票,返回代码、名称、市场"),
		mcp.WithString("keyword", mcp.Required(), mcp.Description("股票名称或代码关键词")),
	), handleSearchStock)

	s.AddTool(mcp.NewTool("get_quote",
		mcp.WithDescription("获取A股实时行情: 价格、涨跌幅、成交量、换手率、PE、PB、市值"),
		mcp.WithString("code", mcp.Required(), mcp.Description("6位股票代码")),
	), handleGetQuote)

	s.AddTool(mcp.NewTool("get_kline",
		mcp.WithDescription("获取K线历史数据(前复权): 日线/周线/月线"),
		mcp.WithString("code", mcp.Required(), mcp.Description("6位股票代码")),
		mcp.WithString("period", mcp.Description("daily(默认)/weekly/monthly"), mcp.Enum("daily", "weekly", "monthly")),
		mcp.WithNumber("limit", mcp.Description("条数,默认30,最大120")),
	), handleGetKLine)

	s.AddTool(mcp.NewTool("get_index",
		mcp.WithDescription("获取主要指数行情(上证/深证/创业板/科创50/沪深300/中证500),含涨跌家数"),
	), handleIndex)

	// ═══════ 资金流向 ═══════

	s.AddTool(mcp.NewTool("get_money_flow",
		mcp.WithDescription("个股当日资金流向: 主力/超大单/大单/中单/小单的流入流出及净额"),
		mcp.WithString("code", mcp.Required(), mcp.Description("6位股票代码")),
	), handleMoneyFlow)

	s.AddTool(mcp.NewTool("get_money_flow_history",
		mcp.WithDescription("个股近N日资金流向趋势,观察主力连续流入/流出"),
		mcp.WithString("code", mcp.Required(), mcp.Description("6位股票代码")),
		mcp.WithNumber("limit", mcp.Description("天数,默认10,最大60")),
	), handleMoneyFlowHistory)

	s.AddTool(mcp.NewTool("get_north_flow",
		mcp.WithDescription("北向资金(沪深股通)近N日净买入数据,观察外资动向"),
		mcp.WithNumber("limit", mcp.Description("天数,默认10")),
	), handleNorthFlow)

	s.AddTool(mcp.NewTool("get_north_holdings",
		mcp.WithDescription("北向资金个股持仓排行,按持仓市值排序"),
		mcp.WithString("market", mcp.Description("sh(沪股通,默认)/sz(深股通)"), mcp.Enum("sh", "sz")),
		mcp.WithNumber("limit", mcp.Description("条数,默认20")),
	), handleNorthHoldings)

	s.AddTool(mcp.NewTool("get_margin_data",
		mcp.WithDescription("两市融资融券余额汇总,观察杠杆情绪"),
		mcp.WithNumber("limit", mcp.Description("天数,默认10")),
	), handleMarginData)

	s.AddTool(mcp.NewTool("get_stock_margin",
		mcp.WithDescription("个股融资融券数据: 融资买入/余额、融券卖出/余额"),
		mcp.WithString("code", mcp.Required(), mcp.Description("6位股票代码")),
		mcp.WithNumber("limit", mcp.Description("天数,默认10")),
	), handleStockMargin)

	// ═══════ 基本面 ═══════

	s.AddTool(mcp.NewTool("get_financial",
		mcp.WithDescription("个股核心财务指标: ROE、毛利率、净利率、营收同比、净利润同比、资产负债率、EPS"),
		mcp.WithString("code", mcp.Required(), mcp.Description("6位股票代码")),
	), handleFinancial)

	s.AddTool(mcp.NewTool("get_top_holders",
		mcp.WithDescription("十大流通股东: 持股数量、占比、增减变动、股东性质(基金/个人/QFII)"),
		mcp.WithString("code", mcp.Required(), mcp.Description("6位股票代码")),
	), handleTopHolders)

	s.AddTool(mcp.NewTool("get_dividend_history",
		mcp.WithDescription("分红送配历史: 分配方案、除权除息日、股息率"),
		mcp.WithString("code", mcp.Required(), mcp.Description("6位股票代码")),
		mcp.WithNumber("limit", mcp.Description("条数,默认10")),
	), handleDividend)

	s.AddTool(mcp.NewTool("get_analyst_ratings",
		mcp.WithDescription("机构评级/研报: 券商、分析师、评级(买入/增持/中性)、目标价"),
		mcp.WithString("code", mcp.Required(), mcp.Description("6位股票代码")),
		mcp.WithNumber("limit", mcp.Description("条数,默认10")),
	), handleAnalystRatings)

	// ═══════ 板块 ═══════

	s.AddTool(mcp.NewTool("get_sectors",
		mcp.WithDescription("行业/概念板块涨跌排行,含领涨股和主力资金"),
		mcp.WithString("type", mcp.Required(), mcp.Description("industry(行业)/concept(概念)"), mcp.Enum("industry", "concept")),
		mcp.WithNumber("limit", mcp.Description("条数,默认20")),
	), handleSectors)

	s.AddTool(mcp.NewTool("get_sector_stocks",
		mcp.WithDescription("板块成分股列表(按涨跌幅排序),需先用 get_sectors 获取板块代码"),
		mcp.WithString("sector_code", mcp.Required(), mcp.Description("板块代码,如'BK0477'")),
		mcp.WithNumber("limit", mcp.Description("条数,默认20")),
	), handleSectorStocks)

	// ═══════ 排行/异动 ═══════

	s.AddTool(mcp.NewTool("get_ranking",
		mcp.WithDescription("A股涨跌/成交排行榜"),
		mcp.WithString("type", mcp.Required(), mcp.Description("top/bottom/volume/amount/turnover"), mcp.Enum("top", "bottom", "volume", "amount", "turnover")),
		mcp.WithNumber("limit", mcp.Description("条数,默认20")),
	), handleRanking)

	s.AddTool(mcp.NewTool("get_limit_stocks",
		mcp.WithDescription("涨停/跌停股票列表: 连板天数、封板时间、开板次数、涨停原因/题材"),
		mcp.WithString("type", mcp.Required(), mcp.Description("up(涨停)/down(跌停)"), mcp.Enum("up", "down")),
		mcp.WithNumber("limit", mcp.Description("条数,默认30")),
	), handleLimitStocks)

	s.AddTool(mcp.NewTool("get_dragon_tiger",
		mcp.WithDescription("龙虎榜: 上榜股票、机构/游资席位、净买入、上榜原因"),
		mcp.WithNumber("limit", mcp.Description("条数,默认20")),
	), handleDragonTiger)

	// ═══════ 杠杆/风控 ═══════

	s.AddTool(mcp.NewTool("get_block_trades",
		mcp.WithDescription("大宗交易: 成交价、折溢价率、买卖方营业部"),
		mcp.WithNumber("limit", mcp.Description("条数,默认20")),
	), handleBlockTrades)

	s.AddTool(mcp.NewTool("get_lockup_expiry",
		mcp.WithDescription("限售解禁日历: 解禁日期、数量、占总股本比例、限售类型"),
		mcp.WithNumber("limit", mcp.Description("条数,默认20")),
	), handleLockupExpiry)

	// ═══════ 资讯 ═══════

	s.AddTool(mcp.NewTool("get_stock_news",
		mcp.WithDescription("个股相关新闻资讯: 标题、来源、摘要"),
		mcp.WithString("code", mcp.Required(), mcp.Description("6位股票代码")),
		mcp.WithNumber("limit", mcp.Description("条数,默认10")),
	), handleStockNews)
}
