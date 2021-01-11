package main

import (
	"flag"
	"fmt"
	"net/http"
	"os"

	"zhouwei/shorturl/lib"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis"
)

// Client reids
var Client *redis.Client

// get redis keys value with hash keys
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

// set redis keys
func setKey(key string, val string) {
	err := Client.Set(key, val, 0).Err()
	if err != nil {
		panic(err)
	} else {
		fmt.Println("key set success")
	}
}

// del redis keys
func delKey(key string) {
	err := Client.Del(key).Err()
	if err != nil {
		panic(err)
	} else {
		fmt.Println("key del success")
	}
}

func main() {

	// static vars
	var redisHost string
	var redisPassWord string
	var redisPort string
	var baseURL string
	var port string

	/// flags
	flag.StringVar(&redisHost, "h", "127.0.0.1", "Redis host, default is 127.0.0.1")
	flag.StringVar(&redisPassWord, "pw", "", "Redis password, default is null")
	flag.StringVar(&redisPort, "port", "6379", "Redis port, default is 6370")
	flag.StringVar(&baseURL, "url", "", "base url")
	flag.StringVar(&port, "p", "3070", "server web port")
	flag.Parse()

	// check input vars
	if baseURL == "" {
		fmt.Println("you need input a baseURL")
		fmt.Println("usage ./shortUrl --help check input params")
		os.Exit(0)
	}

	// print the input params
	fmt.Println("redis host:", redisHost)
	fmt.Println("redis password:", redisPassWord)
	fmt.Println("redis port:", redisPassWord)
	fmt.Println("base url:", baseURL)
	fmt.Println("server port:", port)

	Client = redis.NewClient(&redis.Options{
		Addr:     redisHost + ":" + redisPort,
		Password: redisPassWord,
		DB:       13,
	})

	gin.SetMode(gin.ReleaseMode)
	r := gin.Default()

	// health check
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})

	r.Use(func(c *gin.Context) {
		url := c.Request.RequestURI
		// todo: need to check url with or without http:// or  https://
		if len(url) == 7 {
			val := getKey(url[1:7])
			c.Redirect(http.StatusMovedPermanently, val)
		} else {
			fmt.Println("err")
		}

	})

	r.GET("/makeUrl", func(c *gin.Context) {
		url := c.Query("url")
		if url[0:4] != "http" {
			c.JSON(200, gin.H{
				"error": "url need begin with http:// or https://",
			})
		} else {
			result, _ := lib.Transform(url)

			c.JSON(200, "https://"+baseURL+"/"+result[0])
			setKey(result[0], url)
			fmt.Println("shortUrl is", "https://"+baseURL+"/"+url)
		}
	})

	r.GET("/delUrl", func(c *gin.Context) {
		url := c.Query("url")
		if url[0:4] != "http" {
			c.JSON(200, gin.H{
				"error": "url need begin with http:// or https://",
			})
		} else {
			delKey(url)
			c.JSON(200, gin.H{
				"message": "key delete suuccess",
			})
		}
	})

	fmt.Println("server start on port ", port)
	r.Run(":" + port) // default listen and serve on 0.0.0.0:3070
}
