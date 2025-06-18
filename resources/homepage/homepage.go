package homepage

import "github.com/gin-gonic/gin"

func HomePage(c *gin.Context) {
	c.Header("Content-Type", "text/plain")
	message := "http://example.com/books?title=Sagar+Parmar+Golang+Practice&author=Sagar+Parmar"
	c.String(200, message)
}

func Test() {}
