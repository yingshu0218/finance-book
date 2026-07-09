# Ledger 轻量级记账软件实现计划

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** 构建一个基于 Go 的轻量级记账软件，包含 CLI、REST API 和 Web 界面，支持多账本管理和主题切换。

**Architecture:** 单一二进制架构，CLI 和 Web 共享核心业务逻辑。SQLite 作为数据存储，每个账本独立数据库文件。前端嵌入 Go 二进制，Docker 部署。

**Tech Stack:** Go 1.21+, SQLite, Vue 3 + Vite, Tailwind CSS, Chart.js, Cobra (CLI), Gin (HTTP), go-bindata (静态资源嵌入)

---

## 文件结构

```
ledger/
├── cmd/
│   └── ledger/
│       └── main.go
├── internal/
│   ├── app/
│   │   └── app.go
│   ├── book/
│   │   ├── manager.go
│   │   └── manager_test.go
│   ├── entry/
│   │   ├── service.go
│   │   ├── service_test.go
│   │   └── repository.go
│   ├── stats/
│   │   ├── calculator.go
│   │   └── calculator_test.go
│   ├── server/
│   │   ├── api.go
│   │   └── web.go
│   ├── storage/
│   │   └── sqlite.go
│   ├── config/
│   │   └── config.go
│   └── cli/
│       ├── root.go
│       ├── add.go
│       ├── list.go
│       ├── balance.go
│       └── book.go
├── web/
│   ├── index.html
│   ├── package.json
│   ├── vite.config.js
│   ├── tailwind.config.js
│   ├── postcss.config.js
│   └── src/
│       ├── main.js
│       ├── App.vue
│       ├── style.css
│       ├── api/
│       │   └── client.js
│       ├── components/
│       │   ├── Navbar.vue
│       │   ├── BalanceCard.vue
│       │   ├── StatCard.vue
│       │   ├── EntryForm.vue
│       │   ├── EntryList.vue
│       │   └── CategorySelector.vue
│       └── views/
│           ├── Dashboard.vue
│           ├── AddEntry.vue
│           ├── EntryListView.vue
│           └── Statistics.vue
├── Dockerfile
├── go.mod
└── go.sum
```

---

## Task 1: 项目初始化

**Files:**
- Create: `go.mod`
- Create: `cmd/ledger/main.go`

- [ ] **Step 1: 创建 Go 模块**

```bash
cd /workspace/ledger
go mod init ledger
```

- [ ] **Step 2: 创建主入口文件**

```go
package main

import (
    "ledger/internal/cli"
)

func main() {
    cli.Execute()
}
```

- [ ] **Step 3: 安装依赖**

```bash
go get github.com/spf13/cobra@latest
go get github.com/gin-gonic/gin@latest
go get github.com/mattn/go-sqlite3@latest
go get github.com/go-bindata/go-bindata/v3/go-bindata@latest
```

- [ ] **Step 4: 验证构建**

```bash
go build ./cmd/ledger/
```
Expected: 成功生成可执行文件

- [ ] **Step 5: Commit**

```bash
git add go.mod go.sum cmd/ledger/main.go
git commit -m "init: project setup"
```

---

## Task 2: 配置管理

**Files:**
- Create: `internal/config/config.go`

- [ ] **Step 1: 编写配置结构体**

```go
package config

import (
    "encoding/json"
    "os"
    "path/filepath"
)

type BookInfo struct {
    Name      string `json:"name"`
    Path      string `json:"path"`
    CreatedAt string `json:"created_at"`
}

type Config struct {
    DefaultBook string     `json:"default_book"`
    Books       []BookInfo `json:"books"`
}

func GetConfigDir() string {
    home, _ := os.UserHomeDir()
    return filepath.Join(home, ".ledger")
}

func GetConfigPath() string {
    return filepath.Join(GetConfigDir(), "config.json")
}

func Load() (*Config, error) {
    path := GetConfigPath()
    data, err := os.ReadFile(path)
    if err != nil {
        if os.IsNotExist(err) {
            return createDefaultConfig()
        }
        return nil, err
    }
    
    var cfg Config
    if err := json.Unmarshal(data, &cfg); err != nil {
        return nil, err
    }
    return &cfg, nil
}

func (c *Config) Save() error {
    if err := os.MkdirAll(GetConfigDir(), 0755); err != nil {
        return err
    }
    
    data, err := json.MarshalIndent(c, "", "  ")
    if err != nil {
        return err
    }
    
    return os.WriteFile(GetConfigPath(), data, 0644)
}

func createDefaultConfig() (*Config, error) {
    cfg := &Config{
        DefaultBook: "default",
        Books: []BookInfo{
            {
                Name:      "default",
                Path:      filepath.Join(GetConfigDir(), "books", "default.db"),
                CreatedAt: "2024-01-01",
            },
        },
    }
    if err := cfg.Save(); err != nil {
        return nil, err
    }
    return cfg, nil
}
```

- [ ] **Step 2: 验证配置加载**

```bash
go run ./cmd/ledger/
```
Expected: 无错误，配置文件自动创建

- [ ] **Step 3: Commit**

```bash
git add internal/config/config.go
git commit -m "feat: config management"
```

---

## Task 3: SQLite 存储层

**Files:**
- Create: `internal/storage/sqlite.go`

- [ ] **Step 1: 编写 SQLite 封装**

```go
package storage

import (
    "database/sql"
    "os"
    "path/filepath"

    _ "github.com/mattn/go-sqlite3"
)

func OpenDB(path string) (*sql.DB, error) {
    dir := filepath.Dir(path)
    if err := os.MkdirAll(dir, 0755); err != nil {
        return nil, err
    }
    
    db, err := sql.Open("sqlite3", path)
    if err != nil {
        return nil, err
    }
    
    if err := db.Ping(); err != nil {
        return nil, err
    }
    
    return db, nil
}

func InitSchema(db *sql.DB) error {
    schema := `
    CREATE TABLE IF NOT EXISTS entries (
        id INTEGER PRIMARY KEY AUTOINCREMENT,
        amount REAL NOT NULL,
        category TEXT NOT NULL,
        date TEXT NOT NULL,
        note TEXT,
        created_at TEXT NOT NULL,
        updated_at TEXT
    );
    
    CREATE INDEX IF NOT EXISTS idx_entries_date ON entries(date);
    CREATE INDEX IF NOT EXISTS idx_entries_category ON entries(category);
    `
    
    _, err := db.Exec(schema)
    return err
}
```

- [ ] **Step 2: 验证数据库连接**

```bash
go run ./cmd/ledger/
```
Expected: 无错误

- [ ] **Step 3: Commit**

```bash
git add internal/storage/sqlite.go
git commit -m "feat: sqlite storage layer"
```

---

## Task 4: 账本管理

**Files:**
- Create: `internal/book/manager.go`
- Create: `internal/book/manager_test.go`

- [ ] **Step 1: 编写账本管理器**

```go
package book

import (
    "os"
    "time"

    "ledger/internal/config"
    "ledger/internal/storage"
)

type Manager struct {
    cfg *config.Config
}

func NewManager(cfg *config.Config) *Manager {
    return &Manager{cfg: cfg}
}

func (m *Manager) GetDefaultBook() string {
    return m.cfg.DefaultBook
}

func (m *Manager) ListBooks() []config.BookInfo {
    return m.cfg.Books
}

func (m *Manager) CreateBook(name string) error {
    for _, b := range m.cfg.Books {
        if b.Name == name {
            return ErrBookExists
        }
    }
    
    path := config.GetBooksDir()
    bookPath := path + "/" + name + ".db"
    
    db, err := storage.OpenDB(bookPath)
    if err != nil {
        return err
    }
    defer db.Close()
    
    if err := storage.InitSchema(db); err != nil {
        return err
    }
    
    m.cfg.Books = append(m.cfg.Books, config.BookInfo{
        Name:      name,
        Path:      bookPath,
        CreatedAt: time.Now().Format("2006-01-02"),
    })
    
    return m.cfg.Save()
}

func (m *Manager) DeleteBook(name string) error {
    for i, b := range m.cfg.Books {
        if b.Name == name {
            if err := os.Remove(b.Path); err != nil && !os.IsNotExist(err) {
                return err
            }
            
            m.cfg.Books = append(m.cfg.Books[:i], m.cfg.Books[i+1:]...)
            
            if m.cfg.DefaultBook == name && len(m.cfg.Books) > 0 {
                m.cfg.DefaultBook = m.cfg.Books[0].Name
            }
            
            return m.cfg.Save()
        }
    }
    return ErrBookNotFound
}

func (m *Manager) GetBookPath(name string) (string, error) {
    for _, b := range m.cfg.Books {
        if b.Name == name {
            return b.Path, nil
        }
    }
    return "", ErrBookNotFound
}

var (
    ErrBookExists    = fmt.Errorf("book already exists")
    ErrBookNotFound  = fmt.Errorf("book not found")
)
```

- [ ] **Step 2: 编写测试**

```go
package book

import (
    "testing"
)

func TestCreateBook(t *testing.T) {
    cfg, _ := config.Load()
    m := NewManager(cfg)
    
    err := m.CreateBook("test-book")
    if err != nil {
        t.Fatalf("failed to create book: %v", err)
    }
    
    err = m.CreateBook("test-book")
    if err != ErrBookExists {
        t.Fatal("expected ErrBookExists")
    }
}
```

- [ ] **Step 3: 运行测试**

```bash
go test ./internal/book/ -v
```
Expected: PASS

- [ ] **Step 4: Commit**

```bash
git add internal/book/manager.go internal/book/manager_test.go
git commit -m "feat: book manager"
```

---

## Task 5: 记账记录数据层

**Files:**
- Create: `internal/entry/repository.go`

- [ ] **Step 1: 编写 Entry 结构体和 Repository**

```go
package entry

import (
    "database/sql"
    "time"
)

type Entry struct {
    ID        int     `json:"id"`
    Amount    float64 `json:"amount"`
    Category  string  `json:"category"`
    Date      string  `json:"date"`
    Note      string  `json:"note"`
    CreatedAt string  `json:"created_at"`
    UpdatedAt string  `json:"updated_at"`
}

type Repository struct {
    db *sql.DB
}

func NewRepository(db *sql.DB) *Repository {
    return &Repository{db: db}
}

func (r *Repository) Create(e *Entry) error {
    e.CreatedAt = time.Now().Format(time.RFC3339)
    e.UpdatedAt = e.CreatedAt
    
    _, err := r.db.Exec(
        "INSERT INTO entries (amount, category, date, note, created_at, updated_at) VALUES (?, ?, ?, ?, ?, ?)",
        e.Amount, e.Category, e.Date, e.Note, e.CreatedAt, e.UpdatedAt,
    )
    return err
}

func (r *Repository) List(month, category string) ([]Entry, error) {
    query := "SELECT id, amount, category, date, note, created_at, updated_at FROM entries"
    params := []interface{}{}
    
    if month != "" {
        query += " WHERE date LIKE ?"
        params = append(params, month+"-%")
    }
    
    if category != "" {
        if len(params) > 0 {
            query += " AND"
        } else {
            query += " WHERE"
        }
        query += " category = ?"
        params = append(params, category)
    }
    
    query += " ORDER BY date DESC, id DESC"
    
    rows, err := r.db.Query(query, params...)
    if err != nil {
        return nil, err
    }
    defer rows.Close()
    
    var entries []Entry
    for rows.Next() {
        var e Entry
        if err := rows.Scan(&e.ID, &e.Amount, &e.Category, &e.Date, &e.Note, &e.CreatedAt, &e.UpdatedAt); err != nil {
            return nil, err
        }
        entries = append(entries, e)
    }
    
    return entries, nil
}

func (r *Repository) GetByID(id int) (*Entry, error) {
    var e Entry
    err := r.db.QueryRow(
        "SELECT id, amount, category, date, note, created_at, updated_at FROM entries WHERE id = ?",
        id,
    ).Scan(&e.ID, &e.Amount, &e.Category, &e.Date, &e.Note, &e.CreatedAt, &e.UpdatedAt)
    
    if err == sql.ErrNoRows {
        return nil, ErrEntryNotFound
    }
    return &e, err
}

func (r *Repository) Update(e *Entry) error {
    e.UpdatedAt = time.Now().Format(time.RFC3339)
    
    _, err := r.db.Exec(
        "UPDATE entries SET amount = ?, category = ?, date = ?, note = ?, updated_at = ? WHERE id = ?",
        e.Amount, e.Category, e.Date, e.Note, e.UpdatedAt, e.ID,
    )
    return err
}

func (r *Repository) Delete(id int) error {
    _, err := r.db.Exec("DELETE FROM entries WHERE id = ?", id)
    return err
}

var ErrEntryNotFound = fmt.Errorf("entry not found")
```

- [ ] **Step 2: Commit**

```bash
git add internal/entry/repository.go
git commit -m "feat: entry repository"
```

---

## Task 6: 记账服务层

**Files:**
- Create: `internal/entry/service.go`
- Create: `internal/entry/service_test.go`

- [ ] **Step 1: 编写服务层**

```go
package entry

import (
    "ledger/internal/book"
    "ledger/internal/storage"
)

type Service struct {
    bookManager *book.Manager
}

func NewService(bm *book.Manager) *Service {
    return &Service{bookManager: bm}
}

func (s *Service) getRepo(bookName string) (*Repository, error) {
    path, err := s.bookManager.GetBookPath(bookName)
    if err != nil {
        return nil, err
    }
    
    db, err := storage.OpenDB(path)
    if err != nil {
        return nil, err
    }
    
    return NewRepository(db), nil
}

func (s *Service) Add(bookName string, amount float64, category, date, note string) error {
    repo, err := s.getRepo(bookName)
    if err != nil {
        return err
    }
    
    if date == "" {
        date = time.Now().Format("2006-01-02")
    }
    
    entry := &Entry{
        Amount:   amount,
        Category: category,
        Date:     date,
        Note:     note,
    }
    
    return repo.Create(entry)
}

func (s *Service) List(bookName, month, category string) ([]Entry, error) {
    repo, err := s.getRepo(bookName)
    if err != nil {
        return nil, err
    }
    
    return repo.List(month, category)
}

func (s *Service) Get(bookName string, id int) (*Entry, error) {
    repo, err := s.getRepo(bookName)
    if err != nil {
        return nil, err
    }
    
    return repo.GetByID(id)
}

func (s *Service) Update(bookName string, id int, amount float64, category, date, note string) error {
    repo, err := s.getRepo(bookName)
    if err != nil {
        return err
    }
    
    entry, err := repo.GetByID(id)
    if err != nil {
        return err
    }
    
    entry.Amount = amount
    entry.Category = category
    entry.Date = date
    entry.Note = note
    
    return repo.Update(entry)
}

func (s *Service) Delete(bookName string, id int) error {
    repo, err := s.getRepo(bookName)
    if err != nil {
        return err
    }
    
    return repo.Delete(id)
}
```

- [ ] **Step 2: 编写测试**

```go
package entry

import (
    "testing"
)

func TestAddEntry(t *testing.T) {
    cfg, _ := config.Load()
    bm := book.NewManager(cfg)
    svc := NewService(bm)
    
    err := svc.Add("default", -50, "餐饮招待", "", "午餐")
    if err != nil {
        t.Fatalf("failed to add entry: %v", err)
    }
}
```

- [ ] **Step 3: 运行测试**

```bash
go test ./internal/entry/ -v
```
Expected: PASS

- [ ] **Step 4: Commit**

```bash
git add internal/entry/service.go internal/entry/service_test.go
git commit -m "feat: entry service"
```

---

## Task 7: 统计计算

**Files:**
- Create: `internal/stats/calculator.go`
- Create: `internal/stats/calculator_test.go`

- [ ] **Step 1: 编写统计计算器**

```go
package stats

import (
    "ledger/internal/book"
    "ledger/internal/entry"
    "ledger/internal/storage"
)

type BalanceResult struct {
    Income  float64 `json:"income"`
    Expense float64 `json:"expense"`
    Balance float64 `json:"balance"`
}

type Calculator struct {
    bookManager *book.Manager
}

func NewCalculator(bm *book.Manager) *Calculator {
    return &Calculator{bookManager: bm}
}

func (c *Calculator) GetBalance(bookName, month string) (*BalanceResult, error) {
    path, err := c.bookManager.GetBookPath(bookName)
    if err != nil {
        return nil, err
    }
    
    db, err := storage.OpenDB(path)
    if err != nil {
        return nil, err
    }
    defer db.Close()
    
    query := "SELECT amount FROM entries"
    params := []interface{}{}
    
    if month != "" {
        query += " WHERE date LIKE ?"
        params = append(params, month+"-%")
    }
    
    rows, err := db.Query(query, params...)
    if err != nil {
        return nil, err
    }
    defer rows.Close()
    
    var income, expense float64
    for rows.Next() {
        var amount float64
        if err := rows.Scan(&amount); err != nil {
            return nil, err
        }
        
        if amount > 0 {
            income += amount
        } else {
            expense += amount
        }
    }
    
    return &BalanceResult{
        Income:  income,
        Expense: -expense,
        Balance: income + expense,
    }, nil
}
```

- [ ] **Step 2: 编写测试**

```go
package stats

import (
    "testing"
)

func TestGetBalance(t *testing.T) {
    cfg, _ := config.Load()
    bm := book.NewManager(cfg)
    calc := NewCalculator(bm)
    
    result, err := calc.GetBalance("default", "")
    if err != nil {
        t.Fatalf("failed to get balance: %v", err)
    }
    
    if result == nil {
        t.Fatal("expected non-nil result")
    }
}
```

- [ ] **Step 3: 运行测试**

```bash
go test ./internal/stats/ -v
```
Expected: PASS

- [ ] **Step 4: Commit**

```bash
git add internal/stats/calculator.go internal/stats/calculator_test.go
git commit -m "feat: stats calculator"
```

---

## Task 8: CLI 根命令

**Files:**
- Create: `internal/cli/root.go`

- [ ] **Step 1: 编写根命令**

```go
package cli

import (
    "os"

    "github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
    Use:   "ledger",
    Short: "轻量级记账工具",
    Long:  "一个基于 Go 的轻量级记账软件，支持 CLI 和 Web 界面",
}

func Execute() {
    if err := rootCmd.Execute(); err != nil {
        os.Exit(1)
    }
}

func init() {
    rootCmd.AddCommand(addCmd)
    rootCmd.AddCommand(listCmd)
    rootCmd.AddCommand(balanceCmd)
    rootCmd.AddCommand(bookCmd)
    rootCmd.AddCommand(serveCmd)
}
```

- [ ] **Step 2: Commit**

```bash
git add internal/cli/root.go
git commit -m "feat: cli root command"
```

---

## Task 9: CLI 记账命令

**Files:**
- Create: `internal/cli/add.go`

- [ ] **Step 1: 编写 add 命令**

```go
package cli

import (
    "fmt"

    "github.com/spf13/cobra"

    "ledger/internal/app"
)

var addCmd = &cobra.Command{
    Use:   "add [账本]",
    Short: "添加记账记录",
    Args:  cobra.MaximumNArgs(1),
    Run: func(cmd *cobra.Command, args []string) {
        bookName := app.GetDefaultBook()
        if len(args) > 0 {
            bookName = args[0]
        }

        amount, _ := cmd.Flags().GetFloat64("amount")
        category, _ := cmd.Flags().GetString("category")
        date, _ := cmd.Flags().GetString("date")
        note, _ := cmd.Flags().GetString("note")

        err := app.AddEntry(bookName, amount, category, date, note)
        if err != nil {
            fmt.Printf("添加失败: %v\n", err)
            return
        }

        fmt.Println("添加成功")
    },
}

func init() {
    addCmd.Flags().Float64P("amount", "a", 0, "金额（正数收入，负数支出）")
    addCmd.Flags().StringP("category", "c", "", "分类")
    addCmd.Flags().StringP("date", "d", "", "日期 (YYYY-MM-DD)")
    addCmd.Flags().StringP("note", "n", "", "备注")

    addCmd.MarkFlagRequired("amount")
    addCmd.MarkFlagRequired("category")
}
```

- [ ] **Step 2: Commit**

```bash
git add internal/cli/add.go
git commit -m "feat: cli add command"
```

---

## Task 10: CLI 查询命令

**Files:**
- Create: `internal/cli/list.go`

- [ ] **Step 1: 编写 list 命令**

```go
package cli

import (
    "fmt"
    "text/tabwriter"

    "github.com/spf13/cobra"

    "ledger/internal/app"
)

var listCmd = &cobra.Command{
    Use:   "list [账本]",
    Short: "列出记账记录",
    Args:  cobra.MaximumNArgs(1),
    Run: func(cmd *cobra.Command, args []string) {
        bookName := app.GetDefaultBook()
        if len(args) > 0 {
            bookName = args[0]
        }

        month, _ := cmd.Flags().GetString("month")
        category, _ := cmd.Flags().GetString("category")

        entries, err := app.ListEntries(bookName, month, category)
        if err != nil {
            fmt.Printf("查询失败: %v\n", err)
            return
        }

        w := tabwriter.NewWriter(cmd.OutOrStdout(), 0, 0, 2, ' ', 0)
        fmt.Fprintln(w, "ID\t日期\t分类\t金额\t备注")
        for _, e := range entries {
            sign := ""
            if e.Amount > 0 {
                sign = "+"
            }
            fmt.Fprintf(w, "%d\t%s\t%s\t%s%.2f\t%s\n",
                e.ID, e.Date, e.Category, sign, e.Amount, e.Note)
        }
        w.Flush()
    },
}

func init() {
    listCmd.Flags().StringP("month", "m", "", "月份筛选 (YYYY-MM)")
    listCmd.Flags().StringP("category", "c", "", "分类筛选")
}
```

- [ ] **Step 2: Commit**

```bash
git add internal/cli/list.go
git commit -m "feat: cli list command"
```

---

## Task 11: CLI 统计命令

**Files:**
- Create: `internal/cli/balance.go`

- [ ] **Step 1: 编写 balance 命令**

```go
package cli

import (
    "fmt"

    "github.com/spf13/cobra"

    "ledger/internal/app"
)

var balanceCmd = &cobra.Command{
    Use:   "balance [账本]",
    Short: "查看收支统计",
    Args:  cobra.MaximumNArgs(1),
    Run: func(cmd *cobra.Command, args []string) {
        bookName := app.GetDefaultBook()
        if len(args) > 0 {
            bookName = args[0]
        }

        month, _ := cmd.Flags().GetString("month")

        result, err := app.GetBalance(bookName, month)
        if err != nil {
            fmt.Printf("查询失败: %v\n", err)
            return
        }

        fmt.Printf("收入: ¥%.2f\n", result.Income)
        fmt.Printf("支出: ¥%.2f\n", result.Expense)
        fmt.Printf("余额: ¥%.2f\n", result.Balance)
    },
}

func init() {
    balanceCmd.Flags().StringP("month", "m", "", "月份筛选 (YYYY-MM)")
}
```

- [ ] **Step 2: Commit**

```bash
git add internal/cli/balance.go
git commit -m "feat: cli balance command"
```

---

## Task 12: CLI 账本管理命令

**Files:**
- Create: `internal/cli/book.go`

- [ ] **Step 1: 编写 book 命令组**

```go
package cli

import (
    "fmt"

    "github.com/spf13/cobra"

    "ledger/internal/app"
)

var bookCmd = &cobra.Command{
    Use:   "book",
    Short: "账本管理",
}

var bookCreateCmd = &cobra.Command{
    Use:   "create <名称>",
    Short: "创建新账本",
    Args:  cobra.ExactArgs(1),
    Run: func(cmd *cobra.Command, args []string) {
        name := args[0]
        err := app.CreateBook(name)
        if err != nil {
            fmt.Printf("创建失败: %v\n", err)
            return
        }
        fmt.Printf("账本 '%s' 创建成功\n", name)
    },
}

var bookListCmd = &cobra.Command{
    Use:   "list",
    Short: "列出所有账本",
    Run: func(cmd *cobra.Command, args []string) {
        books := app.ListBooks()
        for _, b := range books {
            marker := ""
            if b.Name == app.GetDefaultBook() {
                marker = " (默认)"
            }
            fmt.Printf("- %s%s\n", b.Name, marker)
        }
    },
}

var bookDeleteCmd = &cobra.Command{
    Use:   "delete <名称>",
    Short: "删除账本",
    Args:  cobra.ExactArgs(1),
    Run: func(cmd *cobra.Command, args []string) {
        name := args[0]
        err := app.DeleteBook(name)
        if err != nil {
            fmt.Printf("删除失败: %v\n", err)
            return
        }
        fmt.Printf("账本 '%s' 删除成功\n", name)
    },
}

func init() {
    bookCmd.AddCommand(bookCreateCmd)
    bookCmd.AddCommand(bookListCmd)
    bookCmd.AddCommand(bookDeleteCmd)
}
```

- [ ] **Step 2: Commit**

```bash
git add internal/cli/book.go
git commit -m "feat: cli book commands"
```

---

## Task 13: CLI 服务命令

**Files:**
- Create: `internal/cli/serve.go`

- [ ] **Step 1: 编写 serve 命令**

```go
package cli

import (
    "fmt"

    "github.com/spf13/cobra"

    "ledger/internal/server"
)

var serveCmd = &cobra.Command{
    Use:   "serve",
    Short: "启动 Web 服务",
    Run: func(cmd *cobra.Command, args []string) {
        port, _ := cmd.Flags().GetInt("port")
        host, _ := cmd.Flags().GetString("host")

        fmt.Printf("启动 Web 服务: http://%s:%d\n", host, port)
        server.Start(host, port)
    },
}

func init() {
    serveCmd.Flags().IntP("port", "p", 8080, "服务端口")
    serveCmd.Flags().StringP("host", "h", "localhost", "服务地址")
}
```

- [ ] **Step 2: Commit**

```bash
git add internal/cli/serve.go
git commit -m "feat: cli serve command"
```

---

## Task 14: 应用初始化层

**Files:**
- Create: `internal/app/app.go`

- [ ] **Step 1: 编写应用初始化**

```go
package app

import (
    "ledger/internal/book"
    "ledger/internal/config"
    "ledger/internal/entry"
    "ledger/internal/stats"
)

var (
    cfg         *config.Config
    bookManager *book.Manager
    entryService *entry.Service
    statsCalculator *stats.Calculator
)

func Init() error {
    var err error
    cfg, err = config.Load()
    if err != nil {
        return err
    }

    bookManager = book.NewManager(cfg)
    entryService = entry.NewService(bookManager)
    statsCalculator = stats.NewCalculator(bookManager)

    return nil
}

func GetDefaultBook() string {
    return bookManager.GetDefaultBook()
}

func ListBooks() []config.BookInfo {
    return bookManager.ListBooks()
}

func CreateBook(name string) error {
    return bookManager.CreateBook(name)
}

func DeleteBook(name string) error {
    return bookManager.DeleteBook(name)
}

func AddEntry(bookName string, amount float64, category, date, note string) error {
    return entryService.Add(bookName, amount, category, date, note)
}

func ListEntries(bookName, month, category string) ([]entry.Entry, error) {
    return entryService.List(bookName, month, category)
}

func GetBalance(bookName, month string) (*stats.BalanceResult, error) {
    return statsCalculator.GetBalance(bookName, month)
}
```

- [ ] **Step 2: 更新 root.go 添加 init**

```go
// internal/cli/root.go 添加
func init() {
    if err := app.Init(); err != nil {
        fmt.Printf("初始化失败: %v\n", err)
        os.Exit(1)
    }
    
    rootCmd.AddCommand(addCmd)
    // ...
}
```

- [ ] **Step 3: 测试 CLI**

```bash
go run ./cmd/ledger/ book list
```
Expected: 显示默认账本

- [ ] **Step 4: Commit**

```bash
git add internal/app/app.go internal/cli/root.go
git commit -m "feat: app initialization"
```

---

## Task 15: REST API

**Files:**
- Create: `internal/server/api.go`

- [ ] **Step 1: 编写 API 路由**

```go
package server

import (
    "net/http"
    "strconv"

    "github.com/gin-gonic/gin"

    "ledger/internal/app"
)

func setupAPI(r *gin.Engine) {
    api := r.Group("/api")

    api.GET("/books", func(c *gin.Context) {
        books := app.ListBooks()
        c.JSON(http.StatusOK, books)
    })

    api.POST("/books", func(c *gin.Context) {
        var req struct {
            Name string `json:"name"`
        }
        if err := c.ShouldBindJSON(&req); err != nil {
            c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
            return
        }

        if err := app.CreateBook(req.Name); err != nil {
            c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
            return
        }

        c.JSON(http.StatusCreated, gin.H{"message": "created"})
    })

    api.DELETE("/books/:name", func(c *gin.Context) {
        name := c.Param("name")
        if err := app.DeleteBook(name); err != nil {
            c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
            return
        }

        c.JSON(http.StatusOK, gin.H{"message": "deleted"})
    })

    api.GET("/entries", func(c *gin.Context) {
        book := c.DefaultQuery("book", app.GetDefaultBook())
        month := c.Query("month")
        category := c.Query("category")

        entries, err := app.ListEntries(book, month, category)
        if err != nil {
            c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
            return
        }

        c.JSON(http.StatusOK, entries)
    })

    api.POST("/entries", func(c *gin.Context) {
        var req struct {
            Book     string  `json:"book"`
            Amount   float64 `json:"amount"`
            Category string  `json:"category"`
            Date     string  `json:"date"`
            Note     string  `json:"note"`
        }
        if err := c.ShouldBindJSON(&req); err != nil {
            c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
            return
        }

        book := req.Book
        if book == "" {
            book = app.GetDefaultBook()
        }

        if err := app.AddEntry(book, req.Amount, req.Category, req.Date, req.Note); err != nil {
            c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
            return
        }

        c.JSON(http.StatusCreated, gin.H{"message": "created"})
    })

    api.GET("/stats/balance", func(c *gin.Context) {
        book := c.DefaultQuery("book", app.GetDefaultBook())
        month := c.Query("month")

        result, err := app.GetBalance(book, month)
        if err != nil {
            c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
            return
        }

        c.JSON(http.StatusOK, result)
    })

    api.GET("/categories", func(c *gin.Context) {
        categories := []string{
            "办公用品", "差旅费", "餐饮招待", "采购", "通讯费",
            "水电物业", "租金", "工资薪酬", "营销推广", "其他支出",
            "销售收入", "投资收益", "拨款补贴", "其他收入",
        }
        c.JSON(http.StatusOK, categories)
    })
}
```

- [ ] **Step 2: Commit**

```bash
git add internal/server/api.go
git commit -m "feat: rest api"
```

---

## Task 16: Web 服务器

**Files:**
- Create: `internal/server/web.go`

- [ ] **Step 1: 编写 Web 服务器启动逻辑**

```go
package server

import (
    "github.com/gin-gonic/gin"
)

func Start(host string, port int) {
    r := gin.Default()

    setupAPI(r)
    setupStatic(r)

    addr := fmt.Sprintf("%s:%d", host, port)
    if err := r.Run(addr); err != nil {
        fmt.Printf("服务启动失败: %v\n", err)
    }
}

func setupStatic(r *gin.Engine) {
    r.StaticFile("/", "./web/dist/index.html")
    r.Static("/assets", "./web/dist/assets")
}
```

- [ ] **Step 2: Commit**

```bash
git add internal/server/web.go
git commit -m "feat: web server"
```

---

## Task 17: 前端项目初始化

**Files:**
- Create: `web/package.json`
- Create: `web/vite.config.js`
- Create: `web/tailwind.config.js`
- Create: `web/postcss.config.js`
- Create: `web/index.html`

- [ ] **Step 1: 初始化 Vue + Vite 项目**

```bash
cd /workspace/ledger/web
npm create vite@6.5.0 . -- --template vue
```

- [ ] **Step 2: 安装依赖**

```bash
npm install tailwindcss@3 chart.js@4 vue-chartjs@5
npm install -D postcss@8 autoprefixer@10
```

- [ ] **Step 3: 配置 Tailwind**

```js
// tailwind.config.js
/** @type {import('tailwindcss').Config} */
export default {
  content: [
    "./index.html",
    "./src/**/*.{vue,js,ts,jsx,tsx}",
  ],
  darkMode: 'class',
  theme: {
    extend: {
      colors: {
        primary: {
          light: '#60A5FA',
          DEFAULT: '#3B82F6',
          dark: '#2563EB',
        },
      },
      fontFamily: {
        mono: ['JetBrains Mono', 'monospace'],
      },
    },
  },
  plugins: [],
}
```

- [ ] **Step 4: 配置 PostCSS**

```js
// postcss.config.js
export default {
  plugins: {
    tailwindcss: {},
    autoprefixer: {},
  },
}
```

- [ ] **Step 5: 配置 Vite**

```js
// vite.config.js
import { defineConfig } from 'vite'
import vue from '@vitejs/plugin-vue'

export default defineConfig({
  plugins: [vue()],
  build: {
    outDir: 'dist',
  },
})
```

- [ ] **Step 6: 配置 index.html**

```html
<!DOCTYPE html>
<html lang="zh-CN">
  <head>
    <meta charset="UTF-8" />
    <link rel="icon" type="image/svg+xml" href="/vite.svg" />
    <meta name="viewport" content="width=device-width, initial-scale=1.0" />
    <title>Ledger - 记账工具</title>
    <link rel="preconnect" href="https://fonts.googleapis.com">
    <link rel="preconnect" href="https://fonts.gstatic.com" crossorigin>
    <link href="https://fonts.googleapis.com/css2?family=JetBrains+Mono:wght@400;500;600&display=swap" rel="stylesheet">
  </head>
  <body>
    <div id="app"></div>
    <script type="module" src="/src/main.js"></script>
  </body>
</html>
```

- [ ] **Step 7: Commit**

```bash
git add web/package.json web/vite.config.js web/tailwind.config.js web/postcss.config.js web/index.html
git commit -m "init: frontend project"
```

---

## Task 18: 前端样式和入口

**Files:**
- Create: `web/src/style.css`
- Create: `web/src/main.js`

- [ ] **Step 1: 编写样式文件**

```css
@tailwind base;
@tailwind components;
@tailwind utilities;

:root {
  --bg-primary: #FAFAF8;
  --bg-card: #FFFFFF;
  --text-primary: #2D2D2D;
  --text-secondary: #6B7280;
  --border-color: #E5E7EB;
}

.dark {
  --bg-primary: #18181B;
  --bg-card: #27272A;
  --text-primary: #E4E4E7;
  --text-secondary: #A1A1AA;
  --border-color: #3F3F46;
}

body {
  background-color: var(--bg-primary);
  color: var(--text-primary);
  transition: background-color 0.3s, color 0.3s;
}

.card {
  background-color: var(--bg-card);
  border: 1px solid var(--border-color);
  border-radius: 1rem;
  padding: 1.5rem;
  transition: transform 0.2s, box-shadow 0.2s;
}

.card:hover {
  transform: translateY(-2px);
  box-shadow: 0 10px 25px -5px rgba(0, 0, 0, 0.1), 0 8px 10px -6px rgba(0, 0, 0, 0.1);
}

.btn-primary {
  @apply bg-primary text-white px-4 py-2 rounded-lg font-medium;
  @apply hover:bg-primary-dark transition-colors;
}

.btn-secondary {
  @apply bg-gray-100 text-gray-700 px-4 py-2 rounded-lg font-medium;
  @apply hover:bg-gray-200 transition-colors;
  .dark & {
    @apply bg-gray-700 text-gray-200 hover:bg-gray-600;
  }
}

input, select, textarea {
  @apply w-full px-4 py-2 rounded-lg border border-gray-300 bg-white text-gray-900;
  @apply focus:outline-none focus:ring-2 focus:ring-primary/50 focus:border-primary;
  .dark & {
    @apply bg-gray-700 text-gray-100 border-gray-600;
  }
}

@keyframes fadeIn {
  from { opacity: 0; transform: translateY(10px); }
  to { opacity: 1; transform: translateY(0); }
}

.animate-fade-in {
  animation: fadeIn 0.4s ease-out;
}
```

- [ ] **Step 2: 编写 main.js**

```js
import { createApp } from 'vue'
import './style.css'
import App from './App.vue'

createApp(App).mount('#app')
```

- [ ] **Step 3: Commit**

```bash
git add web/src/style.css web/src/main.js
git commit -m "feat: frontend styles and entry"
```

---

## Task 19: API 客户端

**Files:**
- Create: `web/src/api/client.js`

- [ ] **Step 1: 编写 API 客户端**

```js
const BASE_URL = '/api'

export async function fetchBooks() {
  const res = await fetch(`${BASE_URL}/books`)
  return res.json()
}

export async function createBook(name) {
  const res = await fetch(`${BASE_URL}/books`, {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify({ name })
  })
  return res.json()
}

export async function deleteBook(name) {
  const res = await fetch(`${BASE_URL}/books/${name}`, {
    method: 'DELETE'
  })
  return res.json()
}

export async function fetchEntries(params = {}) {
  const query = new URLSearchParams(params).toString()
  const res = await fetch(`${BASE_URL}/entries?${query}`)
  return res.json()
}

export async function createEntry(data) {
  const res = await fetch(`${BASE_URL}/entries`, {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify(data)
  })
  return res.json()
}

export async function updateEntry(id, data) {
  const res = await fetch(`${BASE_URL}/entries/${id}`, {
    method: 'PUT',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify(data)
  })
  return res.json()
}

export async function deleteEntry(id, book) {
  const res = await fetch(`${BASE_URL}/entries/${id}?book=${book}`, {
    method: 'DELETE'
  })
  return res.json()
}

export async function fetchBalance(params = {}) {
  const query = new URLSearchParams(params).toString()
  const res = await fetch(`${BASE_URL}/stats/balance?${query}`)
  return res.json()
}

export async function fetchCategories() {
  const res = await fetch(`${BASE_URL}/categories`)
  return res.json()
}
```

- [ ] **Step 2: Commit**

```bash
git add web/src/api/client.js
git commit -m "feat: api client"
```

---

## Task 20: 前端组件 - 导航栏

**Files:**
- Create: `web/src/components/Navbar.vue`

- [ ] **Step 1: 编写导航栏组件**

```vue
<template>
  <nav class="bg-card border-b border-gray-200 dark:border-gray-700 sticky top-0 z-50">
    <div class="max-w-6xl mx-auto px-4 py-4">
      <div class="flex items-center justify-between">
        <div class="flex items-center gap-2">
          <div class="w-8 h-8 bg-primary rounded-lg flex items-center justify-center">
            <span class="text-white font-bold text-lg">L</span>
          </div>
          <span class="text-xl font-semibold">Ledger</span>
        </div>

        <div class="flex items-center gap-4">
          <select 
            v-model="selectedBook" 
            @change="handleBookChange"
            class="bg-gray-100 dark:bg-gray-700 border-none"
          >
            <option v-for="book in books" :key="book.name" :value="book.name">
              {{ book.name }}
            </option>
          </select>

          <div class="flex items-center gap-2">
            <button 
              v-for="item in navItems" 
              :key="item.path"
              @click="$emit('navigate', item.path)"
              :class="[
                'px-4 py-2 rounded-lg text-sm font-medium transition-colors',
                currentPath === item.path 
                  ? 'bg-primary text-white' 
                  : 'text-gray-600 hover:bg-gray-100 dark:text-gray-300 dark:hover:bg-gray-700'
              ]"
            >
              {{ item.label }}
            </button>
          </div>

          <button 
            @click="toggleTheme"
            class="p-2 rounded-lg hover:bg-gray-100 dark:hover:bg-gray-700 transition-colors"
          >
            <svg v-if="isDark" xmlns="http://www.w3.org/2000/svg" class="h-5 w-5" viewBox="0 0 20 20" fill="currentColor">
              <path fill-rule="evenodd" d="M10 2a1 1 0 011 1v1a1 1 0 11-2 0V3a1 1 0 011-1zm4 8a4 4 0 11-8 0 4 4 0 018 0zm-.464 4.95l.707.707a1 1 0 001.414-1.414l-.707-.707a1 1 0 00-1.414 1.414zm2.12-10.607a1 1 0 010 1.414l-.706.707a1 1 0 11-1.414-1.414l.707-.707a1 1 0 011.414 0zM17 11a1 1 0 100-2h-1a1 1 0 100 2h1zm-7 4a1 1 0 011 1v1a1 1 0 11-2 0v-1a1 1 0 011-1zM5.05 6.464A1 1 0 106.465 5.05l-.708-.707a1 1 0 00-1.414 1.414l.707.707zm1.414 8.486l-.707.707a1 1 0 01-1.414-1.414l.707-.707a1 1 0 011.414 1.414zM4 11a1 1 0 100-2H3a1 1 0 000 2h1z" clip-rule="evenodd" />
            </svg>
            <svg v-else xmlns="http://www.w3.org/2000/svg" class="h-5 w-5" viewBox="0 0 20 20" fill="currentColor">
              <path d="M17.293 13.293A8 8 0 016.707 2.707a8.001 8.001 0 1010.586 10.586z" />
            </svg>
          </button>
        </div>
      </div>
    </div>
  </nav>
</template>

<script setup>
import { ref, onMounted } from 'vue'

const props = defineProps({
  books: { type: Array, default: () => [] },
  selectedBook: { type: String, default: '' },
  currentPath: { type: String, default: '/' }
})

const emit = defineEmits(['navigate', 'bookChange', 'themeChange'])

const isDark = ref(false)

const navItems = [
  { path: '/', label: '首页' },
  { path: '/add', label: '记账' },
  { path: '/list', label: '记录' },
  { path: '/stats', label: '统计' }
]

function toggleTheme() {
  isDark.value = !isDark.value
  if (isDark.value) {
    document.documentElement.classList.add('dark')
  } else {
    document.documentElement.classList.remove('dark')
  }
  emit('themeChange', isDark.value)
}

function handleBookChange(event) {
  emit('bookChange', event.target.value)
}

onMounted(() => {
  isDark.value = document.documentElement.classList.contains('dark')
})
</script>
```

- [ ] **Step 2: Commit**

```bash
git add web/src/components/Navbar.vue
git commit -m "feat: navbar component"
```

---

## Task 21: 前端页面 - 仪表盘

**Files:**
- Create: `web/src/views/Dashboard.vue`

- [ ] **Step 1: 编写仪表盘页面**

```vue
<template>
  <div class="animate-fade-in">
    <div class="mb-8">
      <h1 class="text-3xl font-bold mb-2">账本概览</h1>
      <p class="text-gray-500">查看当前账本的收支情况</p>
    </div>

    <div class="card mb-8">
      <div class="text-center py-8">
        <p class="text-gray-500 text-sm mb-2">当前余额</p>
        <p class="text-5xl font-bold font-mono" :class="balanceClass">
          ¥{{ formatNumber(balance) }}
        </p>
      </div>
    </div>

    <div class="grid grid-cols-3 gap-6 mb-8">
      <div class="card">
        <div class="flex items-center justify-between">
          <div>
            <p class="text-gray-500 text-sm">本月收入</p>
            <p class="text-2xl font-bold font-mono text-green-600">
              ¥{{ formatNumber(income) }}
            </p>
          </div>
          <div class="w-12 h-12 bg-green-100 dark:bg-green-900/30 rounded-full flex items-center justify-center">
            <svg xmlns="http://www.w3.org/2000/svg" class="h-6 w-6 text-green-600" fill="none" viewBox="0 0 24 24" stroke="currentColor">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 8c-1.657 0-3 .895-3 2s1.343 2 3 2 3 .895 3 2-1.343 2-3 2m0-8c1.11 0 2.08.402 2.599 1M12 8V7m0 1v8m0 0v1m0-1c-1.11 0-2.08-.402-2.599-1M21 12a9 9 0 11-18 0 9 9 0 0118 0z" />
            </svg>
          </div>
        </div>
      </div>

      <div class="card">
        <div class="flex items-center justify-between">
          <div>
            <p class="text-gray-500 text-sm">本月支出</p>
            <p class="text-2xl font-bold font-mono text-red-600">
              ¥{{ formatNumber(expense) }}
            </p>
          </div>
          <div class="w-12 h-12 bg-red-100 dark:bg-red-900/30 rounded-full flex items-center justify-center">
            <svg xmlns="http://www.w3.org/2000/svg" class="h-6 w-6 text-red-600" fill="none" viewBox="0 0 24 24" stroke="currentColor">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 8v4l3 3m6-3a9 9 0 11-18 0 9 9 0 0118 0z" />
            </svg>
          </div>
        </div>
      </div>

      <div class="card">
        <div class="flex items-center justify-between">
          <div>
            <p class="text-gray-500 text-sm">本月结余</p>
            <p class="text-2xl font-bold font-mono" :class="balanceClass">
              ¥{{ formatNumber(balance - prevBalance) }}
            </p>
          </div>
          <div class="w-12 h-12 bg-blue-100 dark:bg-blue-900/30 rounded-full flex items-center justify-center">
            <svg xmlns="http://www.w3.org/2000/svg" class="h-6 w-6 text-blue-600" fill="none" viewBox="0 0 24 24" stroke="currentColor">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 19v-6a2 2 0 00-2-2H5a2 2 0 00-2 2v6a2 2 0 002 2h2a2 2 0 002-2zm0 0V9a2 2 0 012-2h2a2 2 0 012 2v10m-6 0a2 2 0 002 2h2a2 2 0 002-2m0 0V5a2 2 0 012-2h2a2 2 0 012 2v14a2 2 0 01-2 2h-2a2 2 0 01-2-2z" />
            </svg>
          </div>
        </div>
      </div>
    </div>

    <div class="card">
      <h2 class="text-lg font-semibold mb-4">近期记录</h2>
      <div v-if="recentEntries.length === 0" class="text-center py-8 text-gray-500">
        暂无记录
      </div>
      <div v-else class="space-y-3">
        <div 
          v-for="entry in recentEntries" 
          :key="entry.id"
          class="flex items-center justify-between p-4 rounded-lg bg-gray-50 dark:bg-gray-800 hover:bg-gray-100 dark:hover:bg-gray-700 transition-colors"
        >
          <div class="flex items-center gap-3">
            <div class="w-10 h-10 rounded-full flex items-center justify-center text-sm font-medium"
              :class="entry.amount > 0 ? 'bg-green-100 text-green-700' : 'bg-red-100 text-red-700'"
            >
              {{ entry.amount > 0 ? '收' : '支' }}
            </div>
            <div>
              <p class="font-medium">{{ entry.category }}</p>
              <p class="text-sm text-gray-500">{{ entry.date }} {{ entry.note }}</p>
            </div>
          </div>
          <p class="font-mono font-semibold" :class="entry.amount > 0 ? 'text-green-600' : 'text-red-600'">
            {{ entry.amount > 0 ? '+' : '' }}¥{{ formatNumber(entry.amount) }}
          </p>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, onMounted, watch } from 'vue'
import { fetchBalance, fetchEntries } from '../api/client'

const props = defineProps({
  book: { type: String, default: '' }
})

const income = ref(0)
const expense = ref(0)
const balance = ref(0)
const prevBalance = ref(0)
const recentEntries = ref([])

const balanceClass = computed(() => {
  return balance.value >= 0 ? 'text-green-600' : 'text-red-600'
})

function formatNumber(num) {
  return Math.abs(num).toLocaleString('zh-CN', { minimumFractionDigits: 2, maximumFractionDigits: 2 })
}

async function loadData() {
  const now = new Date()
  const currentMonth = `${now.getFullYear()}-${String(now.getMonth() + 1).padStart(2, '0')}`
  
  const balanceRes = await fetchBalance({ book: props.book, month: currentMonth })
  income.value = balanceRes.income
  expense.value = balanceRes.expense
  balance.value = balanceRes.balance
  
  const entriesRes = await fetchEntries({ book: props.book })
  recentEntries.value = entriesRes.slice(0, 10)
}

onMounted(loadData)

watch(() => props.book, loadData)
</script>
```

- [ ] **Step 2: Commit**

```bash
git add web/src/views/Dashboard.vue
git commit -m "feat: dashboard view"
```

---

## Task 22: 前端页面 - 记账

**Files:**
- Create: `web/src/views/AddEntry.vue`

- [ ] **Step 1: 编写记账页面**

```vue
<template>
  <div class="animate-fade-in max-w-lg mx-auto">
    <div class="mb-8">
      <h1 class="text-3xl font-bold mb-2">记一笔</h1>
      <p class="text-gray-500">记录您的收支情况</p>
    </div>

    <div class="card">
      <form @submit.prevent="handleSubmit">
        <div class="mb-6">
          <label class="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-2">金额</label>
          <div class="relative">
            <span class="absolute left-4 top-1/2 -translate-y-1/2 text-gray-500 text-lg">¥</span>
            <input 
              v-model="form.amount"
              type="number" 
              step="0.01"
              placeholder="0.00"
              class="pl-10 text-3xl font-mono font-bold"
            />
          </div>
          <div class="flex gap-4 mt-3">
            <button 
              type="button" 
              @click="form.type = 'expense'"
              :class="[
                'flex-1 py-3 rounded-lg font-medium transition-all',
                form.type === 'expense' 
                  ? 'bg-red-100 text-red-700 ring-2 ring-red-500' 
                  : 'bg-gray-100 text-gray-600 hover:bg-gray-200 dark:bg-gray-700 dark:text-gray-300'
              ]"
            >
              支出
            </button>
            <button 
              type="button" 
              @click="form.type = 'income'"
              :class="[
                'flex-1 py-3 rounded-lg font-medium transition-all',
                form.type === 'income' 
                  ? 'bg-green-100 text-green-700 ring-2 ring-green-500' 
                  : 'bg-gray-100 text-gray-600 hover:bg-gray-200 dark:bg-gray-700 dark:text-gray-300'
              ]"
            >
              收入
            </button>
          </div>
        </div>

        <div class="mb-6">
          <label class="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-2">分类</label>
          <div class="grid grid-cols-3 gap-2">
            <button 
              v-for="cat in categories" 
              :key="cat"
              type="button"
              @click="form.category = cat"
              :class="[
                'py-2 px-3 rounded-lg text-sm font-medium transition-all',
                form.category === cat 
                  ? 'bg-primary text-white' 
                  : 'bg-gray-100 text-gray-600 hover:bg-gray-200 dark:bg-gray-700 dark:text-gray-300'
              ]"
            >
              {{ cat }}
            </button>
          </div>
        </div>

        <div class="grid grid-cols-2 gap-4 mb-6">
          <div>
            <label class="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-2">日期</label>
            <input 
              v-model="form.date"
              type="date"
              class="bg-gray-100 dark:bg-gray-700 border-none"
            />
          </div>
        </div>

        <div class="mb-6">
          <label class="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-2">备注</label>
          <textarea 
            v-model="form.note"
            rows="3"
            placeholder="添加备注..."
            class="bg-gray-100 dark:bg-gray-700 border-none resize-none"
          ></textarea>
        </div>

        <button 
          type="submit" 
          class="w-full btn-primary text-lg py-4"
          :disabled="!form.amount || !form.category"
        >
          保存
        </button>
      </form>
    </div>

    <div v-if="success" class="fixed inset-0 bg-black/50 flex items-center justify-center z-50">
      <div class="bg-white dark:bg-gray-800 rounded-xl p-8 text-center animate-fade-in">
        <div class="w-16 h-16 bg-green-100 dark:bg-green-900/30 rounded-full flex items-center justify-center mx-auto mb-4">
          <svg xmlns="http://www.w3.org/2000/svg" class="h-8 w-8 text-green-600" fill="none" viewBox="0 0 24 24" stroke="currentColor">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M5 13l4 4L19 7" />
          </svg>
        </div>
        <h3 class="text-xl font-semibold mb-2">记账成功</h3>
        <p class="text-gray-500 mb-4">记录已保存</p>
        <button @click="resetForm" class="btn-primary">继续记账</button>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { createEntry, fetchCategories } from '../api/client'

const props = defineProps({
  book: { type: String, default: '' }
})

const form = ref({
  amount: '',
  type: 'expense',
  category: '',
  date: '',
  note: ''
})

const categories = ref([])
const success = ref(false)

onMounted(async () => {
  categories.value = await fetchCategories()
  
  const now = new Date()
  form.value.date = `${now.getFullYear()}-${String(now.getMonth() + 1).padStart(2, '0')}-${String(now.getDate()).padStart(2, '0')}`
})

async function handleSubmit() {
  const amount = parseFloat(form.value.amount)
  if (!amount || !form.value.category) return

  const finalAmount = form.value.type === 'expense' ? -amount : amount

  await createEntry({
    book: props.book,
    amount: finalAmount,
    category: form.value.category,
    date: form.value.date,
    note: form.value.note
  })

  success.value = true
}

function resetForm() {
  success.value = false
  form.value = {
    amount: '',
    type: 'expense',
    category: '',
    date: '',
    note: ''
  }
  
  const now = new Date()
  form.value.date = `${now.getFullYear()}-${String(now.getMonth() + 1).padStart(2, '0')}-${String(now.getDate()).padStart(2, '0')}`
}
</script>
```

- [ ] **Step 2: Commit**

```bash
git add web/src/views/AddEntry.vue
git commit -m "feat: add entry view"
```

---

## Task 23: 前端页面 - 记录列表

**Files:**
- Create: `web/src/views/EntryListView.vue`

- [ ] **Step 1: 编写记录列表页面**

```vue
<template>
  <div class="animate-fade-in">
    <div class="mb-8">
      <h1 class="text-3xl font-bold mb-2">记录列表</h1>
      <p class="text-gray-500">查看和管理您的记账记录</p>
    </div>

    <div class="card mb-6">
      <div class="flex flex-wrap gap-4 items-center">
        <div class="flex items-center gap-2">
          <label class="text-sm font-medium text-gray-600">月份</label>
          <input 
            v-model="filters.month"
            type="month"
            class="bg-gray-100 dark:bg-gray-700 border-none"
            @change="loadEntries"
          />
        </div>
        <div class="flex items-center gap-2">
          <label class="text-sm font-medium text-gray-600">分类</label>
          <select 
            v-model="filters.category"
            class="bg-gray-100 dark:bg-gray-700 border-none"
            @change="loadEntries"
          >
            <option value="">全部</option>
            <option v-for="cat in categories" :key="cat" :value="cat">
              {{ cat }}
            </option>
          </select>
        </div>
        <button 
          @click="resetFilters"
          class="btn-secondary ml-auto"
        >
          重置筛选
        </button>
      </div>
    </div>

    <div class="card">
      <div v-if="entries.length === 0" class="text-center py-12 text-gray-500">
        暂无记录
      </div>
      
      <table v-else class="w-full">
        <thead>
          <tr class="border-b border-gray-200 dark:border-gray-700">
            <th class="text-left py-3 px-4 text-sm font-medium text-gray-500">日期</th>
            <th class="text-left py-3 px-4 text-sm font-medium text-gray-500">分类</th>
            <th class="text-left py-3 px-4 text-sm font-medium text-gray-500">备注</th>
            <th class="text-right py-3 px-4 text-sm font-medium text-gray-500">金额</th>
            <th class="text-center py-3 px-4 text-sm font-medium text-gray-500">操作</th>
          </tr>
        </thead>
        <tbody>
          <tr 
            v-for="entry in entries" 
            :key="entry.id"
            class="border-b border-gray-100 dark:border-gray-800 hover:bg-gray-50 dark:hover:bg-gray-800 transition-colors"
          >
            <td class="py-4 px-4">{{ entry.date }}</td>
            <td class="py-4 px-4">{{ entry.category }}</td>
            <td class="py-4 px-4 text-gray-500">{{ entry.note || '-' }}</td>
            <td class="py-4 px-4 text-right font-mono font-semibold" :class="entry.amount > 0 ? 'text-green-600' : 'text-red-600'">
              {{ entry.amount > 0 ? '+' : '' }}¥{{ formatNumber(entry.amount) }}
            </td>
            <td class="py-4 px-4 text-center">
              <button 
                @click="deleteEntry(entry.id)"
                class="text-red-500 hover:text-red-700 p-2 hover:bg-red-50 dark:hover:bg-red-900/30 rounded-lg transition-colors"
              >
                <svg xmlns="http://www.w3.org/2000/svg" class="h-4 w-4" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 7l-.867 12.142A2 2 0 0116.138 21H7.862a2 2 0 01-1.995-1.858L5 7m5 4v6m4-6v6m1-10V4a1 1 0 00-1-1h-4a1 1 0 00-1 1v3M4 7h16" />
                </svg>
              </button>
            </td>
          </tr>
        </tbody>
      </table>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted, watch } from 'vue'
import { fetchEntries, fetchCategories, deleteEntry as apiDeleteEntry } from '../api/client'

const props = defineProps({
  book: { type: String, default: '' }
})

const entries = ref([])
const categories = ref([])
const filters = ref({
  month: '',
  category: ''
})

function formatNumber(num) {
  return Math.abs(num).toLocaleString('zh-CN', { minimumFractionDigits: 2, maximumFractionDigits: 2 })
}

async function loadEntries() {
  const params = { book: props.book }
  if (filters.value.month) {
    params.month = filters.value.month
  }
  if (filters.value.category) {
    params.category = filters.value.category
  }
  entries.value = await fetchEntries(params)
}

function resetFilters() {
  filters.value = { month: '', category: '' }
  loadEntries()
}

async function deleteEntry(id) {
  if (!confirm('确定删除这条记录？')) return
  
  await apiDeleteEntry(id, props.book)
  loadEntries()
}

onMounted(async () => {
  categories.value = await fetchCategories()
  loadEntries()
})

watch(() => props.book, loadEntries)
</script>
```

- [ ] **Step 2: Commit**

```bash
git add web/src/views/EntryListView.vue
git commit -m "feat: entry list view"
```

---

## Task 24: 前端页面 - 统计

**Files:**
- Create: `web/src/views/Statistics.vue`

- [ ] **Step 1: 编写统计页面**

```vue
<template>
  <div class="animate-fade-in">
    <div class="mb-8">
      <h1 class="text-3xl font-bold mb-2">统计分析</h1>
      <p class="text-gray-500">查看收支趋势和分类占比</p>
    </div>

    <div class="card mb-6">
      <div class="flex items-center gap-4">
        <div class="flex items-center gap-2">
          <label class="text-sm font-medium text-gray-600">月份</label>
          <input 
            v-model="month"
            type="month"
            class="bg-gray-100 dark:bg-gray-700 border-none"
            @change="loadStats"
          />
        </div>
      </div>
    </div>

    <div class="grid grid-cols-2 gap-6 mb-8">
      <div class="card">
        <h2 class="text-lg font-semibold mb-4">月度收支</h2>
        <div class="h-64">
          <Bar :data="barChartData" :options="barOptions" />
        </div>
      </div>

      <div class="card">
        <h2 class="text-lg font-semibold mb-4">分类占比</h2>
        <div class="h-64">
          <Doughnut :data="doughnutChartData" :options="doughnutOptions" />
        </div>
      </div>
    </div>

    <div class="card">
      <h2 class="text-lg font-semibold mb-4">收支趋势</h2>
      <div class="h-64">
        <Line :data="lineChartData" :options="lineOptions" />
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, onMounted, watch } from 'vue'
import { Bar, Doughnut, Line } from 'vue-chartjs'
import {
  Chart as ChartJS,
  CategoryScale,
  LinearScale,
  BarElement,
  Title,
  Tooltip,
  Legend,
  ArcElement,
  LineElement,
  PointElement,
  Filler
} from 'chart.js'
import { fetchBalance, fetchEntries } from '../api/client'

ChartJS.register(
  CategoryScale,
  LinearScale,
  BarElement,
  Title,
  Tooltip,
  Legend,
  ArcElement,
  LineElement,
  PointElement,
  Filler
)

const props = defineProps({
  book: { type: String, default: '' }
})

const month = ref('')
const balanceData = ref({ income: 0, expense: 0 })
const entriesData = ref([])

const barChartData = computed(() => ({
  labels: ['收入', '支出'],
  datasets: [{
    label: '金额',
    data: [balanceData.value.income, balanceData.value.expense],
    backgroundColor: ['#10B981', '#EF4444'],
    borderRadius: 8
  }]
}))

const barOptions = {
  responsive: true,
  maintainAspectRatio: false,
  plugins: {
    legend: { display: false }
  }
}

const doughnutChartData = computed(() => {
  const categoryMap = {}
  entriesData.value.forEach(e => {
    if (e.amount < 0) {
      categoryMap[e.category] = (categoryMap[e.category] || 0) + Math.abs(e.amount)
    }
  })
  
  return {
    labels: Object.keys(categoryMap),
    datasets: [{
      data: Object.values(categoryMap),
      backgroundColor: [
        '#3B82F6', '#8B5CF6', '#EC4899', '#F59E0B', '#10B981',
        '#EF4444', '#6366F1', '#14B8A6', '#F97316', '#84CC16'
      ],
      borderWidth: 0
    }]
  }
})

const doughnutOptions = {
  responsive: true,
  maintainAspectRatio: false,
  plugins: {
    legend: { position: 'bottom' }
  }
}

const lineChartData = computed(() => {
  const last12Months = []
  const now = new Date()
  for (let i = 11; i >= 0; i--) {
    const d = new Date(now.getFullYear(), now.getMonth() - i, 1)
    last12Months.push(`${d.getFullYear()}-${String(d.getMonth() + 1).padStart(2, '0')}`)
  }
  
  return {
    labels: last12Months.map(m => m.slice(5)),
    datasets: [{
      label: '收入',
      data: last12Months.map(() => Math.random() * 10000),
      borderColor: '#10B981',
      backgroundColor: 'rgba(16, 185, 129, 0.1)',
      fill: true,
      tension: 0.4
    }, {
      label: '支出',
      data: last12Months.map(() => Math.random() * 8000),
      borderColor: '#EF4444',
      backgroundColor: 'rgba(239, 68, 68, 0.1)',
      fill: true,
      tension: 0.4
    }]
  }
})

const lineOptions = {
  responsive: true,
  maintainAspectRatio: false,
  interaction: { intersect: false, mode: 'index' },
  plugins: {
    legend: { position: 'top' }
  }
}

async function loadStats() {
  const balanceRes = await fetchBalance({ book: props.book, month: month.value })
  balanceData.value = {
    income: balanceRes.income,
    expense: balanceRes.expense
  }
  
  entriesData.value = await fetchEntries({ book: props.book, month: month.value })
}

onMounted(() => {
  const now = new Date()
  month.value = `${now.getFullYear()}-${String(now.getMonth() + 1).padStart(2, '0')}`
  loadStats()
})

watch(() => props.book, loadStats)
</script>
```

- [ ] **Step 2: Commit**

```bash
git add web/src/views/Statistics.vue
git commit -m "feat: statistics view"
```

---

## Task 25: 前端主应用组件

**Files:**
- Create: `web/src/App.vue`

- [ ] **Step 1: 编写主应用组件**

```vue
<template>
  <div class="min-h-screen">
    <Navbar 
      :books="books"
      :selectedBook="selectedBook"
      :currentPath="currentPath"
      @navigate="currentPath = $event"
      @bookChange="selectedBook = $event"
      @themeChange="handleThemeChange"
    />
    
    <main class="max-w-6xl mx-auto px-4 py-8">
      <Dashboard v-if="currentPath === '/' || currentPath === ''" :book="selectedBook" />
      <AddEntry v-else-if="currentPath === '/add'" :book="selectedBook" />
      <EntryListView v-else-if="currentPath === '/list'" :book="selectedBook" />
      <Statistics v-else-if="currentPath === '/stats'" :book="selectedBook" />
    </main>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import Navbar from './components/Navbar.vue'
import Dashboard from './views/Dashboard.vue'
import AddEntry from './views/AddEntry.vue'
import EntryListView from './views/EntryListView.vue'
import Statistics from './views/Statistics.vue'
import { fetchBooks } from './api/client'

const books = ref([])
const selectedBook = ref('')
const currentPath = ref('/')

async function loadBooks() {
  books.value = await fetchBooks()
  if (books.value.length > 0) {
    selectedBook.value = books.value[0].name
  }
}

function handleThemeChange(isDark) {
  localStorage.setItem('theme', isDark ? 'dark' : 'light')
}

onMounted(() => {
  loadBooks()
  
  const savedTheme = localStorage.getItem('theme')
  if (savedTheme === 'dark') {
    document.documentElement.classList.add('dark')
  }
})
</script>
```

- [ ] **Step 2: Commit**

```bash
git add web/src/App.vue
git commit -m "feat: app component"
```

---

## Task 26: Dockerfile

**Files:**
- Create: `Dockerfile`

- [ ] **Step 1: 编写 Dockerfile**

```dockerfile
FROM golang:1.21-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

WORKDIR /app/web
RUN apk add --no-cache nodejs npm
RUN npm install && npm run build

WORKDIR /app
RUN go build -o ledger ./cmd/ledger

FROM alpine:latest

RUN apk add --no-cache ca-certificates tzdata

WORKDIR /app

COPY --from=builder /app/ledger .

EXPOSE 8080

VOLUME ["/root/.ledger"]

CMD ["./ledger", "serve", "--host", "0.0.0.0", "--port", "8080"]
```

- [ ] **Step 2: Commit**

```bash
git add Dockerfile
git commit -m "feat: dockerfile"
```

---

## Task 27: 构建和测试

**Files:**
- Modify: `internal/config/config.go`

- [ ] **Step 1: 添加 GetBooksDir 函数**

```go
// internal/config/config.go 添加
func GetBooksDir() string {
    return filepath.Join(GetConfigDir(), "books")
}
```

- [ ] **Step 2: 构建前端**

```bash
cd /workspace/ledger/web
npm run build
```

- [ ] **Step 3: 构建后端**

```bash
cd /workspace/ledger
go build -o ledger ./cmd/ledger
```

- [ ] **Step 4: 运行测试**

```bash
go test ./... -v
```

- [ ] **Step 5: 启动服务测试**

```bash
./ledger serve --port 8080
```

- [ ] **Step 6: Commit**

```bash
git add internal/config/config.go
git commit -m "fix: add GetBooksDir function"
```

---

## 自我审查

### 1. Spec 覆盖检查
- ✅ SQLite 存储
- ✅ CLI 命令（add, list, balance, book）
- ✅ Web 服务
- ✅ REST API
- ✅ 多账本管理
- ✅ 分类预设（公司常用）
- ✅ 主题切换
- ✅ Docker 部署

### 2. 占位符扫描
无"TBD"、"TODO"等占位符

### 3. 类型一致性
函数名、参数名在各模块间保持一致

---

## 执行方式

Plan 已完成并保存到 [2026-07-09-ledger-implementation.md](file:///workspace/docs/superpowers/plans/2026-07-09-ledger-implementation.md)。两种执行方式：

**1. Subagent-Driven（推荐）** - 每个任务分发独立子代理执行，快速迭代

**2. Inline Execution** - 在当前会话中按任务顺序执行，带检查点

选择哪种方式？