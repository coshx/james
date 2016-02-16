package main

import (
    "strings"
    "github.com/gin-gonic/gin"
)

func main() {
    router := gin.Default()
    router.Static("/assets", "./assets")
    router.LoadHTMLGlob("templates/*")

    router.GET("/", func(c *gin.Context) {
        c.HTML(200, "index.tmpl", gin.H {})
    })

    router.GET("/search", func(c *gin.Context) {
        rawKeywords := strings.Split(c.Query("keywords"), "+")
        keywords := make([]string, len(rawKeywords))
        whenSearchIsDone := make(chan []WrappedBlogPost)

        for _, s := range rawKeywords {
            keywords = append(keywords, strings.ToLower(s))
        }

        go search(keywords, whenSearchIsDone)

        posts := <- whenSearchIsDone

        c.JSON(200, gin.H {"posts": toJSON(posts)})
    })

    router.POST("/reset", func(c *gin.Context) {
        resetIndexation()
        indexBlog()
        c.JSON(200, gin.H {})
    })

    router.Run()
}