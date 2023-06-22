package main

import "fmt"

type RedisTester struct {
	loadFactor float64
	redisI     RedisInterface
}

func NewRedisTester(loadFactor float64, redisI RedisInterface) *RedisTester {
	return &RedisTester{
		loadFactor,
		redisI,
	}
}

func (r *RedisTester) LoadALotOfData() {
	for i := 0.0; i < r.loadFactor; i++ {
		if err := r.redisI.Set(fmt.Sprint(i), fmt.Sprint(i)); err != nil {
			panic(err)
		}
	}
}

func (r *RedisTester) GetData(key string) (string, error) {
	return r.redisI.Get(key)
}
