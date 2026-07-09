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

	api.PUT("/entries/:id", func(c *gin.Context) {
		idStr := c.Param("id")
		id, err := strconv.Atoi(idStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
			return
		}

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

		if err := app.UpdateEntry(book, id, req.Amount, req.Category, req.Date, req.Note); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "updated"})
	})

	api.DELETE("/entries/:id", func(c *gin.Context) {
		idStr := c.Param("id")
		id, err := strconv.Atoi(idStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
			return
		}

		book := c.DefaultQuery("book", app.GetDefaultBook())

		if err := app.DeleteEntry(book, id); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "deleted"})
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
