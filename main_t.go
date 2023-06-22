package main

import (
	"fmt"
	"os"
	"runtime"
	"strconv"
	"testing"
)

var loadFactor float64

func init() {
	loadFactorFlag := os.Getenv("LOAD_FACTOR")
	var err error
	loadFactor, err = strconv.ParseFloat(loadFactorFlag, 64)
	if err != nil {
		loadFactor = 1e3
	}
}

func BenchmarkGetFromMemory(b *testing.B) {
	redis := NewMemoryRedis()
	RunTester(b, redis, loadFactor)
}

func BenchmarkGetFromRedis(b *testing.B) {
	client, err := GetRedisCon()
	redis := NewRedis(client)

	if err != nil {
		panic(err)
	}
	RunTester(b, redis, loadFactor)
}

func RunTester(b *testing.B, redis RedisInterface, loadFactor float64) {
	redisTester := NewRedisTester(loadFactor, redis)

	redisTester.LoadALotOfData()

	b.ResetTimer()
	b.ReportAllocs()
	b.StartTimer()

	allocs := GetAllocs(func() {
		for i := 0; i < b.N; i++ {
			for j := 0.0; j < loadFactor; j++ {
				if _, err := redis.Get(fmt.Sprint(j)); err != nil {
					println(j)
					panic(fmt.Sprintf("getting: %s", err))
				}
			}
		}
	})

	_ = allocs
	b.ReportMetric(float64(allocs), "alloc-isolated")
	b.StopTimer()

	if err := redis.Clean(); err != nil {
		panic(fmt.Sprintf("cleaning: %s", err))
	}

}

func GetAllocs(fn func()) uint64 {
	var memStats runtime.MemStats
	runtime.ReadMemStats(&memStats)
	allocs := memStats.Mallocs
	fn()
	runtime.ReadMemStats(&memStats)
	return memStats.Mallocs - allocs
}
