package main

import (
	"net/http"
	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
)

type person struct {
	Id string `json:"id" form:"id"`
	FirstName string `json:"first_name" form:"first_name"`
	LastName string `json:"last_name" form:"last_name"`
}

var people = []person {
	{Id: "user01", FirstName: "first", LastName: "one"},
	{Id: "user02", FirstName: "second", LastName: "two"},
	{Id: "user03", FirstName: "third", LastName: "three"},
}

var baseUrlPath = "/api" // Reverse Proxy 연동을 위한 base경로 

func getPeople(c *gin.Context) {
	c.IndentedJSON(http.StatusOK,people)
}

// API 서버 실행 : go run .
func main() {
	
	router := gin.Default()
	router.GET(baseUrlPath + "/people", getPeople)

	router.Run("localhost:8989")
	
	// db, err := sql.Open("mysql", "root:password1004!@tcp(127.0.0.1:3306)/user?parseTime=true")
	
	// if err != nil {
	// 	log.Fatalln(err)
	// }
	// db.SetMaxIdleConns(20)
	// db.SetMaxOpenConns(20)
	// if err := db.Ping(); err != nil {
	// 	log.Fatalln(err)
	// }
	// defer db.Close()
	// router := gin.Default()
	// router.GET("/", func(c *gin.Context) {
	// 	c.String(http.StatusOK, "API works!!!")
	// })

	// router.Run(":8000");
	// server.ReverseProxy("8088");


	// r := gin.Default()
	// r.GET("/ping", func(c *gin.Context) {
	// 	c.JSON(200, gin.H{
	// 		"message" : "gin framework installation succeed",
	// 	})
	// })
	// r.Run()
}