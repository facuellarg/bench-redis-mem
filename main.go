package main

import (
	"log"
	"net/http"
	"strconv"
	"testing"
)

func benchmarkHandler(w http.ResponseWriter, r *http.Request) {
	// Run the benchmark command

	loadFactorString := r.URL.Query().Get("load-factor")
	goroutinesString := r.URL.Query().Get("goroutines")
	loadFactor, _ = strconv.ParseFloat(loadFactorString, 64)
	goroutines, _ = strconv.Atoi(goroutinesString)
	w.Header().Set("Content-Type", "text/plain")

	redisResult := testing.Benchmark(BenchmarkGetFromRedis)
	w.Write([]byte("redis:\n" + redisResult.String() + "\n"))

	memoryResult := testing.Benchmark(BenchmarkGetFromMemory)
	w.Write([]byte("memory:\n" + memoryResult.String() + "\n"))

	redisConcurrentResult := testing.Benchmark(BenchmarkGetConcurrentFromRedis)
	w.Write([]byte("redis concurrent:\n" + redisConcurrentResult.String() + "\n"))
	memoryConcurrentResult := testing.Benchmark(BenchmarkGetConcurrentFromMemory)
	w.Write([]byte("memory concurrent:\n" + memoryConcurrentResult.String() + "\n"))
	// output := []byte("redis:\n" + redisResult.String() + "\n\nmemory:\n" + memoryResult.String() + "\n\nredis concurrent:\n" + redisConcurrentResult.String() + "\n\nmemory concurrent:\n" + memoryConcurrentResult.String())

	// Write the benchmark output to the response
	// w.Write(output)
}

func main() {
	http.HandleFunc("/benchmark", benchmarkHandler)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
