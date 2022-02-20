package main

import (
	"encoding/xml"
	"github.com/gin-gonic/gin"
)

func IndexHandler(c *gin.Context) {
	name := c.Params.ByName("name")
	c.JSON(200, gin.H{
		"message": "hello " + name,
	})
}

type Person struct {
	XMLName xml.Name `xml:"person"`
	Name    string   `xml:"name,attr"`
}

func XmlHandler(c *gin.Context) {
	name := c.Params.ByName("name")
	c.XML(200, Person{Name: name})
}

func main() {
	router := gin.Default()
	router.GET("/hello/:name", IndexHandler)
	router.GET("/xml/:name", XmlHandler)
	router.Run()
}
