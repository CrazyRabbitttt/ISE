package redis_lock

import (
	"fmt"
	"github.com/gomodule/redigo/redis"
	"time"
)

func TestSetAndGet() {
	c, err := redis.Dial("tcp", "127.0.0.1:6379")
	if err != nil {
		fmt.Println("Connect to redis err")
		return
	}
	defer c.Close()

	_, err = c.Do("SET", "name", "shaoguixin", "EX", "5")
	if err != nil {
		fmt.Println("redis set failed:", err)
	}

	username, err := redis.String(c.Do("GET", "name"))
	if err != nil {
		fmt.Println("redis get failed:", err)
	} else {
		fmt.Printf("Get mykey: %v \n", username)
	}

	time.Sleep(8 * time.Second)

	username, err = redis.String(c.Do("GET", "name"))
	if err != nil {
		fmt.Println("redis get failed:", err)
	} else {
		fmt.Println("Get name: %v \n", username)
	}
}

func TestMultiSetGet() {
	c, err := redis.Dial("tcp", "127.0.0.1:6379")
	if err != nil {
		fmt.Println("Connect to redis error")
		return
	}
	defer c.Close()

	_, err = c.Do("MSET", "name", "shaoguixin", "sex", "male")
	if err != nil {
		fmt.Println("redis set failed:", err)
	}

	keyExist, err := redis.Bool(c.Do("EXISTS", "name"))
	if err != nil {
		fmt.Println("error:", err)
	} else {
		fmt.Printf("Is name key exists ? : %v \n", keyExist)
	}
	reply, err := redis.Values(c.Do("MGET", "name", "sex"))
	if err != nil {
		fmt.Println("multi get error:", err)
	} else {
		var name string
		var sex string
		_, err := redis.Scan(reply, &name, &sex)
		if err != nil {
			fmt.Printf("Scan error \n。")
		} else {
			fmt.Printf("The name is %v, sex is %v \n", name, sex)
		}
	}

}
