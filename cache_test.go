package utils

import (
	"context"
	"fmt"
	"testing"
)

func TestStruct(t *testing.T) {
	redisClient := NewRedisClient(context.Background())
	type Page struct {
		Name string
		Age  int
	}
	data := Page{
		Name: "mongo",
		Age:  12,
	}

	k := "pluging.redis.test2."
	err := redisClient.SaveStruct(k, data, 300)
	if err != nil {
		t.Errorf("savestruct(%s, %v, %d) error:%s", k, data, 300, err.Error())
	}

	dt := Page{}
	err = redisClient.GetStruct(k, &dt)
	if err != nil {
		t.Errorf("cache.GetStruct(%s, %v) error:%s", k, &dt, err)
	} else {
		fmt.Println(dt)
	}
}
