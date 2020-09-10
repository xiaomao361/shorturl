package main

import (
	"fmt"
	"net/http"

	"zhouwei/shorturl/lib"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis"
)

// Client reids
var Client *redis.Client

func getKey(key string) string {
	val, err := Client.Get(key).Result()
	if err == redis.Nil {
		fmt.Println("keys doe not exist")
	} else if err != nil {
		panic(err)
	} else {
		fmt.Println(key, val)
	}
	return val
}

func setKey(key string, val string) {
	err := Client.Set(key, val, 0).Err()
	if err != nil {
		panic(err)
	} else {
		fmt.Println("key set success")
	}
}

func init() {
	Client = redis.NewClient(&redis.Options{
		Addr:     "test.yjwh.shop:6379",
		Password: "123456",

		DB: 13,
	})
}

func main() {
	gin.SetMode(gin.ReleaseMode)
	r := gin.Default()
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})

	r.Use(func(c *gin.Context) {
		url := c.Request.RequestURI
		if len(url) == 7 {
			val := getKey(url[1:7])
			c.Redirect(http.StatusMovedPermanently, val)
		} else {
			fmt.Println("err")
		}

	})

	r.GET("/makeUrl", func(c *gin.Context) {
		url := c.Query("url")
		result, _ := lib.Transform(url)

		c.JSON(200, "https://d.au32.cn/"+result[0])
		setKey(result[0], url)
		fmt.Println("shortUrl is https://d.au32.cn/", url)
	})

	r.Run(":3070") // listen and serve on 0.0.0.0:8080
}
