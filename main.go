package main

import "github.com/gin-gonic/gin"

func setupRouter() *gin.Engine {
	r := gin.Default()
	r.POST("/posts", CreatePost)
	r.GET("/posts", GetPosts)
	r.GET("/posts/:id", GetPost)
	r.PATCH("/posts/:id", UpdatePost)
	r.DELETE("/posts/:id", DeletePost)
	return r
}


func main() {
	r := setupRouter()
	r.Run(":8080")
}

func GetPosts(c *gin.Context) {}
func GetPost(c *gin.Context) {}
func CreatePost(c *gin.Context) {}
func UpdatePost(c *gin.Context) {}
func DeletePost(c *gin.Context) {}
