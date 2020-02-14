package main

import (
	"fmt"
	"github.com/gomodule/redigo/redis"
	"sync"
	"time"
)

var wait sync.WaitGroup

func main() {
	food := NewFoods()
	start:=time.Now().UnixNano() / 1e6
	for i := 0; i <= 10000; i++ {
		wait.Add(1)
		go func() {
			startbuy(&food)
		}()
	}
	wait.Wait()
	fmt.Println("消耗时间：",time.Now().UnixNano() / 1e6-start)
	fmt.Println(food.surple())
}
func startbuy(food *Food) {
	food.deal <- 1 //阻塞
	defer func() { <-food.deal;wait.Done() }()
	Alock := NewGetRedis(food.name, Redispool.Get())
	err := Alock.Lock()
	if err != nil {
		return
	}
	food.get()
	err = Alock.Unlock() //想删除的是Alock锁，但是Alock 已经被自动删除 ,Block由于value 不一样，所以也不会删除
	if err != nil {
		fmt.Println(err)
	}

}

var Redispool *redis.Pool

func init() {
	Redispool = &redis.Pool{
		MaxIdle:     16,
		MaxActive:   0,
		IdleTimeout: 300 * time.Second,
		Dial: func() (redis.Conn, error) {
			tcp := fmt.Sprintf("%s:%d", "127.0.0.1", 6379)
			c, err := redis.Dial("tcp", tcp)
			if err != nil {
				return nil, err
			}
			//fmt.Println("connect redis success!")
			return c, err
		},
		TestOnBorrow: func(c redis.Conn, t time.Time) error {
			_, err := c.Do("PING")
			return err
		},
	}

}
