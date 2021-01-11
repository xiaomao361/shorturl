package redisClient

import (
	"fmt"

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
