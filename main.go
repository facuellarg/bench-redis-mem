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

	redisResult := testing.Benchmark(BenchmarkGetFromRedis)
	memoryResult := testing.Benchmark(BenchmarkGetFromMemory)

	output := []byte("redis:\n" + redisResult.String() + "\n\nmemory:\n" + memoryResult.String())
	// os.Setenv("LOAD_FACTOR", loadFactorString)
	// cmd := exec.Command("go", "test", "-bench=.")
	// fmt.Printf("Running: %s\n", cmd)
	// output, err := cmd.CombinedOutput()
	// if err != nil {
	// 	http.Error(w, fmt.Sprintf("Failed to run benchmark: %s", err), http.StatusInternalServerError)
	// 	return
	// }

	// Write the benchmark output to the response
	w.Header().Set("Content-Type", "text/plain")
	w.Write(output)
}

func main() {
	http.HandleFunc("/benchmark", benchmarkHandler)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
