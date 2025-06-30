package inits

import (
	"context"
	"fmt"

	"github.com/redis/go-redis/v9"
)

var Client *redis.Client

func ExampleClient() {
	Client = redis.NewClient(&redis.Options{
		Addr:     "14.103.243.149:6379",
		Password: "2003225zyh", // no password set
		DB:       0,            // use default DB
	})

	err := Client.Set(context.Background(), "key", "value", 0).Err()
	if err != nil {
		panic(err)
	}

	val, err := Client.Get(context.Background(), "key").Result()
	if err != nil {
		panic(err)
	}
	fmt.Println("key", val)

	val2, err := Client.Get(context.Background(), "key2").Result()
	if err == redis.Nil {
		fmt.Println("key2 does not exist")
	} else if err != nil {
		panic(err)
	} else {
		fmt.Println("key2", val2)
	}

}
