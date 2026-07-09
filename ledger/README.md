# Finance Book

轻量级个人记账工具，基于 Go 语言开发，支持 CLI 和 Web 双端操作。

## 特性

- **轻量级**：单一二进制文件，零依赖部署
- **数据易迁移**：SQLite 单文件存储，复制即可迁移
- **Agent 友好**：CLI 命令行接口，便于 AI Agent 调用
- **多账本管理**：支持多个独立账本，完全隔离
- **Web 界面**：提供可视化查账和手动记账功能
- **主题切换**：支持亮色/暗色主题
- **Docker 支持**：容器化部署

## 技术栈

- **后端**：Go 1.21, SQLite, Gin, Cobra
- **前端**：Vue 3, Vite, Tailwind CSS, Chart.js

## 快速开始

### 直接运行

```bash
# 构建
go build -o ledger ./cmd/ledger

# 启动 Web 服务
./ledger serve --port 8080

# 访问 http://localhost:8080
```

### Docker

```bash
docker build -t finance-book .

docker run -d \
  -p 8080:8080 \
  -v ~/.ledger:/root/.ledger \
  finance-book
```

### CLI 记账

```bash
# 查看账本列表
ledger book list

# 添加记账记录（支出）
ledger add --amount -100 --category 餐饮招待 --note 午餐

# 添加记账记录（收入）
ledger add --amount 5000 --category 销售收入 --note 订单收入

# 查看记录
ledger list

# 查看收支统计
ledger balance

# 按月份筛选
ledger list --month 2024-01
ledger balance --month 2024-01

# 创建新账本
ledger book create 项目A

# 使用指定账本
ledger add 项目A --amount -200 --category 差旅费
```

## CLI 命令

### 账本管理

| 命令 | 说明 |
|------|------|
| `ledger book list` | 列出所有账本 |
| `ledger book create <名称>` | 创建新账本 |
| `ledger book delete <名称>` | 删除账本 |

### 记账操作

| 命令 | 说明 |
|------|------|
| `ledger add [账本] --amount <金额> --category <分类>` | 添加记账记录 |
| `ledger list [账本] [--month <月份>] [--category <分类>]` | 列出记录 |
| `ledger balance [账本] [--month <月份>]` | 查看收支统计 |
| `ledger serve [--port <端口>]` | 启动 Web 服务 |

### 参数说明

| 参数 | 简写 | 说明 |
|------|------|------|
| `--amount` | `-a` | 金额（正数收入，负数支出） |
| `--category` | `-c` | 分类 |
| `--date` | `-d` | 日期 (YYYY-MM-DD)，默认当天 |
| `--note` | `-n` | 备注 |
| `--month` | `-m` | 月份筛选 (YYYY-MM) |
| `--port` | `-p` | 服务端口，默认 8080 |

## 分类预设

### 支出类

- 办公用品
- 差旅费
- 餐饮招待
- 采购
- 通讯费
- 水电物业
- 租金
- 工资薪酬
- 营销推广
- 其他支出

### 收入类

- 销售收入
- 投资收益
- 拨款补贴
- 其他收入

## API 接口

### 账本管理

| 方法 | 路径 | 说明 |
|------|------|------|
| GET | `/api/books` | 获取账本列表 |
| POST | `/api/books` | 创建新账本 |
| DELETE | `/api/books/:name` | 删除账本 |

### 记账记录

| 方法 | 路径 | 说明 |
|------|------|------|
| GET | `/api/entries` | 查询记录 |
| POST | `/api/entries` | 新增记录 |
| PUT | `/api/entries/:id` | 修改记录 |
| DELETE | `/api/entries/:id` | 删除记录 |

### 统计

| 方法 | 路径 | 说明 |
|------|------|------|
| GET | `/api/stats/balance` | 收支统计 |
| GET | `/api/categories` | 获取分类列表 |

## 项目结构

```
finance-book/
├── cmd/ledger/main.go          # 入口
├── internal/
│   ├── app/app.go              # 应用初始化
│   ├── book/manager.go         # 账本管理
│   ├── entry/                  # 记账服务 + 数据访问
│   ├── stats/calculator.go     # 统计计算
│   ├── server/                 # REST API + Web 服务器
│   ├── storage/sqlite.go       # SQLite 封装
│   └── cli/                    # CLI 命令
├── web/                        # Vue 3 前端
├── Dockerfile
├── go.mod
└── README.md
```

## 数据存储

- 账本数据存储在 `~/.ledger/books/` 目录
- 每个账本对应一个 `.db` SQLite 文件
- 配置文件存储在 `~/.ledger/config.json`

## 开发

### 后端开发

```bash
go run ./cmd/ledger/ serve --port 8080
```

### 前端开发

```bash
cd web
npm install
npm run dev
```

### 构建

```bash
# 构建前端
cd web && npm run build

# 构建后端
go build -o ledger ./cmd/ledger
```

## License

MIT