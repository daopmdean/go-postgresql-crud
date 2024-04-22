package main

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/daopmdean/go-postgresql-crud/model"
	"github.com/gin-gonic/gin"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	dsn := "host=localhost user=postgres password=mystrongpassword dbname=postgres port=5432 sslmode=disable TimeZone=Asia/Ho_Chi_Minh"

	db, err := gorm.Open(postgres.Open(dsn))
	if err != nil {
		panic(err)
	}

	db.AutoMigrate(model.Book{})

	r := gin.Default()

	r.GET("/books", func(ctx *gin.Context) {
		var books []*model.Book
		db.Find(&books)
		ctx.JSON(http.StatusOK, gin.H{
			"books": books,
		})
	})

	r.GET("/books/:id", func(ctx *gin.Context) {
		idStr := ctx.Param("id")
		id, err := strconv.Atoi(idStr)
		if err != nil {
			ctx.JSON(http.StatusOK, gin.H{
				"error": "invalid id",
			})
			return
		}

		var book model.Book
		db.First(&book, id)
		ctx.JSON(http.StatusOK, gin.H{
			"book": book,
		})
	})

	r.POST("/books", func(ctx *gin.Context) {
		book := model.Book{}
		if err := ctx.ShouldBindJSON(&book); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"error": "failed binding json",
			})
			return
		}

		if book.Name == "" {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"error": "book name required",
			})
			return
		}

		tx := db.Create(&book)
		if tx.Error != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"error": "failed create book",
			})
			return
		}

		ctx.JSON(http.StatusOK, gin.H{
			"id":     book.ID,
			"create": "ok",
		})
	})

	r.PUT("/books", func(ctx *gin.Context) {
		book := model.Book{}
		if err := ctx.ShouldBindJSON(&book); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"error": "failed binding json",
			})
			return
		}

		if book.ID == 0 {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"error": "book id required",
			})
			return
		}

		if book.Name == "" {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"error": "book name required",
			})
			return
		}

		tx := db.Model(&book).Updates(model.Book{
			Name:      book.Name,
			Author:    book.Author,
			Publisher: book.Publisher,
		})
		if tx.Error != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"error": fmt.Sprintf("failed update book: %s", tx.Error),
			})
			return
		}

		ctx.JSON(http.StatusOK, gin.H{
			"id":     book.ID,
			"update": "ok",
		})
	})

	r.DELETE("/books/:id", func(ctx *gin.Context) {
		idStr := ctx.Param("id")
		id, err := strconv.Atoi(idStr)
		if err != nil {
			ctx.JSON(http.StatusOK, gin.H{
				"error": "invalid id",
			})
			return
		}

		if id == 0 {
			ctx.JSON(http.StatusOK, gin.H{
				"error": "invalid id",
			})
			return
		}

		tx := db.Delete(&model.Book{}, id)
		if tx.Error != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"error": fmt.Sprintf("failed delete book: %s", tx.Error),
			})
			return
		}

		ctx.JSON(http.StatusOK, gin.H{
			"id":     id,
			"delete": "ok",
		})
	})

	r.Run()
}
