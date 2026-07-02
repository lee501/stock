# A股 MCP Server

基于东方财富免费 API 的 A 股全维度行情 MCP Server，提供 22 个 Tool，覆盖行情、资金、基本面、板块、异动、风控六大场景。无需 API Key，编译即用。

## 快速开始

### 构建

```bash
cd stock
go mod tidy
go build -o stock .
```

### 配置 Claude Desktop

编辑配置文件：

- macOS: `~/Library/Application Support/Claude/claude_desktop_config.json`
- Windows: `%APPDATA%\Claude\claude_desktop_config.json`
- Linux: `~/.config/Claude/claude_desktop_config.json`

```json
{
  "mcpServers": {
    "a-stock": {
      "command": "/absolute/path/to/stock"
    }
  }
}
```

重启 Claude Desktop，左下角出现锤子图标即表示连接成功。

### 配置 Claude Code

```bash
claude mcp add a-stock /absolute/path/to/stock
```

## 项目结构

```
stock/
├── go.mod
├── README.md
├── main.go                    # 入口: server 初始化 + tool 注册
├── handlers.go                # 所有 handler 函数 + 参数解析辅助
└── eastmoney/                 # 东方财富 API 客户端(按领域拆分)
    ├── client.go              # HTTP client + 6个通用查询函数 + URL常量 + 解析辅助
    ├── quote.go               # 行情: 搜索、实时报价、K线、指数
    ├── flow.go                # 资金: 个股资金流、北向资金、融资融券
    ├── fundamental.go         # 基本面: 财务指标、十大股东、分红、研报
    ├── market.go              # 市场: 板块、排行、涨停跌停、龙虎榜、大宗、解禁
    ├── news.go                # 资讯: 个股新闻
    └── eastmoney_test.go      # 全量集成测试(22个Tool)
```

**设计原则：**

- `main.go` 只做 server 启动和 tool 元数据注册，不含业务逻辑
- `handlers.go` 负责 MCP 协议层的参数解析和结果序列化，调用 `eastmoney` 包
- `eastmoney/` 是纯粹的 API 客户端，不依赖 MCP，可独立复用
- `eastmoney/` 内部按领域拆文件，每个文件职责单一
- API 的 URL、参数构建、响应解析全部集中在 `client.go`，业务函数只声明查询参数和字段映射

## Tool 一览

### 基础行情

| Tool | 功能 | 必填参数 | 可选参数 |
|------|------|---------|---------|
| `search_stock` | 模糊搜索股票 | `keyword` | — |
| `get_quote` | 实时行情(价格/涨跌/PE/PB/市值) | `code` | — |
| `get_kline` | K线历史数据(前复权) | `code` | `period`(daily/weekly/monthly), `limit` |
| `get_index` | 主要指数行情 + 涨跌家数 | — | — |

### 资金流向

| Tool | 功能 | 必填参数 | 可选参数 |
|------|------|---------|---------|
| `get_money_flow` | 个股当日资金流向(主力/大单/中单/小单) | `code` | — |
| `get_money_flow_history` | 个股近N日资金流向趋势 | `code` | `limit` |
| `get_north_flow` | 北向资金每日净买入 | — | `limit` |
| `get_north_holdings` | 北向资金个股持仓排行 | — | `market`(sh/sz), `limit` |
| `get_margin_data` | 两市融资融券余额汇总 | — | `limit` |
| `get_stock_margin` | 个股融资融券数据 | `code` | `limit` |

### 基本面

| Tool | 功能 | 必填参数 | 可选参数 |
|------|------|---------|---------|
| `get_financial` | 核心财务指标(ROE/毛利率/净利率/营收增速) | `code` | — |
| `get_top_holders` | 十大流通股东(持股/占比/增减/性质) | `code` | — |
| `get_dividend_history` | 分红送配历史 | `code` | `limit` |
| `get_analyst_ratings` | 机构评级/研报(券商/目标价) | `code` | `limit` |

### 板块

| Tool | 功能 | 必填参数 | 可选参数 |
|------|------|---------|---------|
| `get_sectors` | 行业/概念板块涨跌排行 | `type`(industry/concept) | `limit` |
| `get_sector_stocks` | 板块成分股列表 | `sector_code` | `limit` |

> 使用流程：先调 `get_sectors` 获取板块代码，再用 `get_sector_stocks` 查成分股。

### 排行与异动

| Tool | 功能 | 必填参数 | 可选参数 |
|------|------|---------|---------|
| `get_ranking` | 涨跌/成交排行榜 | `type`(top/bottom/volume/amount/turnover) | `limit` |
| `get_limit_stocks` | 涨停/跌停分析(连板/封板/题材) | `type`(up/down) | `limit` |
| `get_dragon_tiger` | 龙虎榜(净买入/上榜原因) | — | `limit` |

### 杠杆与风控

| Tool | 功能 | 必填参数 | 可选参数 |
|------|------|---------|---------|
| `get_block_trades` | 大宗交易(折溢价/买卖方) | — | `limit` |
| `get_lockup_expiry` | 限售解禁日历 | — | `limit` |

### 资讯

| Tool | 功能 | 必填参数 | 可选参数 |
|------|------|---------|---------|
| `get_stock_news` | 个股相关新闻 | `code` | `limit` |

## 使用示例

对 Claude 说自然语言，它会自动组合调用合适的 Tool。

### 查询类

**搜索股票**
```
"帮我找一下和锂电池相关的股票"
→ search_stock("锂电池")
```

**实时行情**
```
"茅台现在什么价？"
→ search_stock("茅台") → get_quote("600519")
返回: 价格、涨跌幅、成交量、换手率、PE、PB、总市值、流通市值
```

**K线走势**
```
"看看比亚迪最近一个月的日K"
→ get_kline("002594", period="daily", limit=22)

"宁德时代的周K走势"
→ get_kline("300750", period="weekly", limit=20)

"贵州茅台近一年月K"
→ get_kline("600519", period="monthly", limit=12)
```

**大盘指数**
```
"今天大盘怎么样"
→ get_index()
返回: 上证指数、深证成指、创业板指、科创50、沪深300、中证500，含涨跌家数
```

### 资金流向

**个股资金流**
```
"宁德时代今天资金流向怎么样"
→ get_money_flow("300750")
返回: 主力/超大单/大单/中单/小单的流入流出及净额

"平安银行最近10天的资金流向趋势"
→ get_money_flow_history("000001", limit=10)
```

**北向资金**
```
"北向资金最近在买什么"
→ get_north_flow(limit=10) → get_north_holdings(market="sh") → get_north_holdings(market="sz")

"北向资金最近5天净流入多少"
→ get_north_flow(limit=5)
返回: 沪股通/深股通每日净买入、累计净买入
```

**融资融券**
```
"最近杠杆资金情绪怎么样"
→ get_margin_data(limit=10)
返回: 两市融资融券余额汇总趋势

"看看平安银行的融资融券数据"
→ get_stock_margin("000001", limit=10)
返回: 融资买入额/余额、融券卖出量/余额
```

### 基本面分析

**财务指标**
```
"比亚迪的财务状况怎么样"
→ get_financial("002594")
返回: ROE、毛利率、净利率、营收同比、净利润同比、资产负债率、EPS
```

**股东分析**
```
"宁德时代的十大流通股东是谁"
→ get_top_holders("300750")
返回: 股东名称、持股数量、占比、增减变动、股东性质(基金/个人/QFII)
```

**分红历史**
```
"茅台历年分红情况"
→ get_dividend_history("600519", limit=10)
返回: 分配方案、除权除息日、股息率
```

**机构评级**
```
"券商怎么看宁德时代"
→ get_analyst_ratings("300750", limit=10)
返回: 券商、分析师、评级(买入/增持/中性)、目标价、研报标题
```

### 板块分析

**板块排行**
```
"今天哪些行业板块涨得好"
→ get_sectors(type="industry", limit=20)

"最近有什么热门概念"
→ get_sectors(type="concept", limit=20)
返回: 板块名称、涨跌幅、主力资金净流入、领涨股
```

**板块成分股**
```
"固态电池板块有哪些值得关注的票"
→ get_sectors(type="concept") → 找到板块代码 → get_sector_stocks("BK0477", limit=20)
返回: 成分股列表，含价格、涨跌幅、成交额、主力资金
```

### 排行与异动

**涨跌排行**
```
"今天涨幅最大的股票"
→ get_ranking(type="top", limit=20)

"今天跌幅最大的股票"
→ get_ranking(type="bottom", limit=20)

"今天成交额最大的股票"
→ get_ranking(type="amount", limit=20)

"今天换手率最高的股票"
→ get_ranking(type="turnover", limit=20)
```

**涨停跌停**
```
"今天有哪些涨停板"
→ get_limit_stocks(type="up", limit=30)
返回: 涨停股列表，含连板天数、封板时间、开板次数、所属行业

"今天有哪些跌停"
→ get_limit_stocks(type="down", limit=30)
```

**龙虎榜**
```
"今天龙虎榜有什么"
→ get_dragon_tiger(limit=20)
返回: 上榜股票、净买入、买卖总额、上榜原因、换手率
```

### 风控排查

**大宗交易**
```
"最近有哪些大宗交易"
→ get_block_trades(limit=20)
返回: 成交价、折溢价率、成交量、买卖方营业部
```

**限售解禁**
```
"最近有哪些股票要解禁"
→ get_lockup_expiry(limit=20)
返回: 解禁日期、解禁数量、占总股本比例、限售类型
```

### 资讯

**个股新闻**
```
"平安银行最近有什么新闻"
→ get_stock_news("000001", limit=10)
返回: 标题、来源、日期、摘要
```

### 综合场景

**个股全面分析**
```
"帮我全面分析一下宁德时代"
→ search_stock → get_quote → get_financial → get_money_flow → get_kline
→ get_top_holders → get_analyst_ratings → get_stock_news
```

**大盘情绪判断**
```
"今天市场情绪怎么样"
→ get_index → get_north_flow → get_margin_data → get_limit_stocks(up) → get_sectors(concept)
```

**短线机会扫描**
```
"今天有什么短线机会"
→ get_limit_stocks(up) → get_dragon_tiger → get_ranking(turnover) → get_sectors(concept)
```

**个股排雷**
```
"帮我看看这只票有没有雷"
→ get_lockup_expiry → get_stock_margin → get_top_holders → get_block_trades → get_dividend_history
```

## 参数说明

| 参数 | 格式 | 说明 |
|------|------|------|
| `code` | 6位数字 | 股票代码，如 `600519`(茅台)、`300750`(宁德时代) |
| `keyword` | 字符串 | 支持名称、代码、拼音首字母模糊匹配 |
| `period` | 枚举 | K线周期: `daily` / `weekly` / `monthly` |
| `limit` | 正整数 | 返回条数，各接口有不同上限，超出自动截断 |
| `sector_code` | 字符串 | 板块代码(如 `BK0477`)，从 `get_sectors` 返回值获取 |
| `market` | 枚举 | 北向持仓市场: `sh`(沪股通) / `sz`(深股通) |
| `type` | 枚举 | 各 Tool 各有不同枚举值，见上表 |

## 数据源

所有数据来自东方财富公开行情接口，免费无需注册。域名和基础 URL 统一定义在 `client.go`：

| 域名 | 用途 | 对应查询函数 |
|------|------|-------------|
| `push2.eastmoney.com` | 实时行情、资金流、板块、排行、财务指标 | `Push2StockGet` / `ClistGet` / `Push2DiffGet` |
| `push2his.eastmoney.com` | 历史K线、历史资金流 | `Push2HisGet` |
| `push2ex.eastmoney.com` | 涨停/跌停股池 | `basePush2Ex` 常量 |
| `searchapi.eastmoney.com` | 股票搜索 | `baseSearch` 常量 |
| `datacenter-web.eastmoney.com` | 数据中心(龙虎榜/融资/大宗/解禁/北向持仓/股东/分红) | `DatacenterGet` |
| `reportapi.eastmoney.com` | 机构评级/研报 | `baseReport` 常量 |
| `search-api-web.eastmoney.com` | 新闻搜索 | `baseNews` 常量 |

**已知限制：**

- 盘后数据通常 15:30 后更新完毕
- 财务指标仅在财报发布后更新
- 高频调用可能触发限流，建议控制合理频率
- 数据仅供参考，不构成投资建议

## 扩展指南

**添加新 Tool 三步走：**

1. `eastmoney/` 下对应文件新增 API 函数（结构体定义 + 调用通用查询 + 字段映射）
2. `handlers.go` 新增 handler 函数（参数解析 → 调 API → 返回 JSON）
3. `main.go` 的 `registerTools()` 新增 `s.AddTool(...)` 注册

**`client.go` 通用查询函数：**

| 函数 | 用途 | 适用接口 |
|------|------|---------|
| `DatacenterGet` | datacenter-web 报表查询 | 龙虎榜、大宗、解禁、北向持仓、融资融券、股东、分红等 |
| `ClistGet` | push2 列表排行查询 | 板块排行、成分股、涨跌排行 |
| `Push2StockGet` | push2 单股数据查询 | 实时行情、财务指标、资金流向 |
| `Push2DiffGet` | push2 多标的 diff 查询 | 指数行情 |
| `Push2HisGet` | push2his K线类查询 | K线、历史资金流向 |

添加新的 datacenter 类接口只需声明 `DatacenterQuery` 即可，无需手写 URL：

```go
data, err := DatacenterGet(DatacenterQuery{
    ReportName:  "RPT_SOME_REPORT",
    SortColumns: "TRADE_DATE",
    SortTypes:   "-1",
    PageSize:    limit,
    Filter:      fmt.Sprintf(`(SECURITY_CODE="%s")`, code),
})
```

**可扩展方向：**

- 财务报表明细（资产负债表 / 利润表 / 现金流量表）
- ETF / 可转债行情
- 港股通标的行情
- 期货 / 期权数据
- Level 2 盘口数据（需付费数据源）

## License

MIT
