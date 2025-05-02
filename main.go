package main

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"

	"github.com/DincerY/blogging-api/db"
	"github.com/DincerY/blogging-api/models"
	"github.com/DincerY/blogging-api/repositories"
)

func createBlog(c *gin.Context) {
	var blog models.Blog
	if err := c.ShouldBindJSON(&blog); err != nil {
		c.JSON(http.StatusBadRequest, "invalid value format")
		return
	}

	blogRepo := repositories.NewBlogRepository(db.DB)
	createdId, err := blogRepo.Create(blog)

	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}
	blog.ID = createdId
	c.JSON(http.StatusCreated, "Adding blog operation was success")
}

func updateBlog(c *gin.Context) {
	paramId := c.Param("id")
	id, err := strconv.Atoi(paramId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, "id was not valid")
		return
	}
	blogRepo := repositories.NewBlogRepository(db.DB)
	var b models.Blog
	err = c.ShouldBindJSON(&b)
	if err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}
	err = blogRepo.Update(b, id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, "blog was updated")
}

func deleteBlog(c *gin.Context) {
	paramId := c.Param("id")
	id, err := strconv.Atoi(paramId)
	if err != nil {
		c.JSON(http.StatusBadRequest, "invalid id")
		return
	}

	blogRepo := repositories.NewBlogRepository(db.DB)
	err = blogRepo.Delete(id)
	if err != nil {
		c.JSON(http.StatusNotFound, err.Error())
		return
	}
	c.JSON(http.StatusNoContent, nil)
}
func getAll(c *gin.Context) {
	category := c.Query("term")
	blogRepo := repositories.NewBlogRepository(db.DB)
	var blogs []models.Blog
	var err error
	if category != "" {
		blogs, err = blogRepo.GetByCategory(category)
	} else {
		blogs, err = blogRepo.GetAll()
	}
	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, blogs)
}

func getById(c *gin.Context) {
	paramId := c.Param("id")
	id, err := strconv.Atoi(paramId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, "id is invalid")
		return
	}
	blogRepo := repositories.NewBlogRepository(db.DB)
	blogs, err := blogRepo.GetById(id)

	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, blogs)
}

func main() {

	err := db.Init()
	if err != nil {
		panic(err)
	}

	r := gin.Default()

	r.POST("/posts", createBlog)
	r.PUT("/posts/:id", updateBlog)
	r.DELETE("posts/:id", deleteBlog)
	r.GET("/posts/:id", getById)
	r.GET("/posts", getAll)

	r.Run(":8080")
}
