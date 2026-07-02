# A股 MCP Server

基于东方财富免费 API 的 A 股全维度行情 MCP Server，提供 22 个 Tool，覆盖行情、资金、基本面、板块、异动、风控六大场景。无需 API Key，编译即用。

## 快速开始

### 构建

```bash
cd stock-mcp
go mod tidy
go build -o stock-mcp .
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
      "command": "/absolute/path/to/stock-mcp"
    }
  }
}
```

重启 Claude Desktop，左下角出现锤子图标即表示连接成功。

### 配置 Claude Code

```bash
claude mcp add a-stock /absolute/path/to/stock-mcp
```

## 项目结构

```
stock-mcp/
├── go.mod
├── README.md
├── main.go                    # 入口: server 初始化 + tool 注册
├── handlers.go                # 所有 handler 函数 + 参数解析辅助
└── eastmoney/                 # 东方财富 API 客户端(按领域拆分)
    ├── client.go              # HTTP client + 通用解析辅助(ToSecID/GetFloat/ToStr...)
    ├── quote.go               # 行情: 搜索、实时报价、K线、指数
    ├── flow.go                # 资金: 个股资金流、北向资金、融资融券
    ├── fundamental.go         # 基本面: 财务指标、十大股东、分红、研报
    ├── market.go              # 市场: 板块、排行、涨停跌停、龙虎榜、大宗、解禁
    └── news.go                # 资讯: 个股新闻
```

**设计原则：**

- `main.go` 只做 server 启动和 tool 元数据注册，不含业务逻辑
- `handlers.go` 负责 MCP 协议层的参数解析和结果序列化，调用 `eastmoney` 包
- `eastmoney/` 是纯粹的 API 客户端，不依赖 MCP，可独立复用
- `eastmoney/` 内部按领域拆文件，每个文件职责单一

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

对 Claude 说自然语言，它会自动组合调用合适的 Tool：

**个股全面分析**
```
"帮我全面分析一下宁德时代"
→ get_quote → get_financial → get_money_flow → get_kline → get_top_holders → get_analyst_ratings → get_stock_news
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

**北向资金追踪**
```
"北向资金最近在买什么"
→ get_north_flow → get_north_holdings(sh) → get_north_holdings(sz)
```

**板块挖掘**
```
"固态电池板块有哪些值得关注的票"
→ get_sectors(concept) → get_sector_stocks → 逐一 get_quote + get_financial
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

所有数据来自东方财富公开行情接口，免费无需注册:

| 域名 | 用途 |
|------|------|
| `push2.eastmoney.com` | 实时行情、资金流、板块、排行 |
| `push2his.eastmoney.com` | 历史K线、历史资金流 |
| `searchapi.eastmoney.com` | 股票搜索 |
| `datacenter-web.eastmoney.com` | 数据中心(龙虎榜/融资/大宗/解禁/研报/股东/分红) |
| `search-api-web.eastmoney.com` | 新闻搜索 |

**已知限制：**

- 盘后数据通常 15:30 后更新完毕
- 财务指标仅在财报发布后更新
- 高频调用可能触发限流，建议控制合理频率
- 数据仅供参考，不构成投资建议

## 扩展指南

**添加新 Tool 三步走：**

1. `eastmoney/` 下对应文件新增 API 函数（结构体定义 + HTTP 请求 + JSON 解析）
2. `handlers.go` 新增 handler 函数（参数解析 → 调 API → 返回 JSON）
3. `main.go` 的 `registerTools()` 新增 `s.AddTool(...)` 注册

**可扩展方向：**

- 财务报表明细（资产负债表 / 利润表 / 现金流量表）
- ETF / 可转债行情
- 港股通标的行情
- 期货 / 期权数据
- Level 2 盘口数据（需付费数据源）

## License

MIT
