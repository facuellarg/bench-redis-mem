package main

import (
	"fmt"
	"runtime"
	"sync"
	"testing"
)

var loadFactor float64
var goroutines int

func BenchmarkGetFromMemory(b *testing.B) {
	redis := NewMemoryRedis()
	RunSerialTester(b, redis, loadFactor)
}

func BenchmarkGetFromRedis(b *testing.B) {
	client, err := GetRedisCon()
	redis := NewRedis(client)

	if err != nil {
		panic(err)
	}
	RunSerialTester(b, redis, loadFactor)
}

func BenchmarkGetConcurrentFromMemory(b *testing.B) {
	redis := NewMemoryRedis()
	RunConcurrentTester(b, redis, loadFactor, goroutines)
}

func BenchmarkGetConcurrentFromRedis(b *testing.B) {
	client, err := GetRedisCon()
	redis := NewRedis(client)

	if err != nil {
		panic(err)
	}
	RunConcurrentTester(b, redis, loadFactor, goroutines)
}

func RunSerialTester(b *testing.B, redis RedisInterface, loadFactor float64) {
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

func RunConcurrentTester(b *testing.B, redis RedisInterface, loadFactor float64, goroutines int) {
	redisTester := NewRedisTester(loadFactor, redis)

	redisTester.LoadALotOfData()

	b.ResetTimer()
	b.ReportAllocs()
	b.StartTimer()

	allocs := GetAllocs(func() {
		for i := 0; i < b.N; i++ {
			for j := 0.0; j < loadFactor; j++ {
				wg := sync.WaitGroup{}
				wg.Add(goroutines)
				for k := 0; k < goroutines; k++ {
					go func(j float64) {
						defer wg.Done()
						if _, err := redis.Get(fmt.Sprint(j)); err != nil {
							println(j)
							panic(fmt.Sprintf("getting: %s", err))
						}
					}(j)
					wg.Wait()
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
