package main

import (
	"context"
	"fmt"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/memorydb"
	"github.com/go-redis/redis/v8"
)

func GetRedisCon() (*redis.Client, error) {

	redisAddress := "host.docker.internal:6379"
	if os.Getenv("ENV") == "" {
		addr, port := memoryDB()
		redisAddress = fmt.Sprintf("%s:%d", addr, port)
	}
	fmt.Println(redisAddress)

	// Create a Redis client using the go-redis library
	client := redis.NewClient(&redis.Options{
		Addr: redisAddress,
	})

	// Ping the Redis server to check the connectivity
	_, err := client.Ping(client.Context()).Result()
	if err != nil {
		return nil, err
	}

	return client, nil
}

func memoryDB() (string, int64) {
	// Create a session using your AWS credentials
	sess, err := session.NewSessionWithOptions(session.Options{
		Config: aws.Config{
			Region: aws.String("us-east-1"),
		},
	})
	if err != nil {

		fmt.Println("Error creating session:", err)
		return "", 0
	}

	// Create a MemoryDB service client
	svc := memorydb.New(sess)

	// Retrieve information about your MemoryDB clusters
	input := &memorydb.DescribeClustersInput{
		ClusterName: aws.String("redis-test"),
	}

	result, err := svc.DescribeClustersWithContext(context.Background(), input)
	if err != nil {
		fmt.Println("Error describing clusters:", err)
		return "", 0
	}

	// Access cluster details from the result variable
	for _, cluster := range result.Clusters {
		fmt.Println("Cluster name", *cluster.Name)
		fmt.Println("Cluster DP:", *cluster.ClusterEndpoint)
		fmt.Println("Status:", *cluster.Status)
		return *cluster.ClusterEndpoint.Address, *cluster.ClusterEndpoint.Port
		// Add any other desired cluster information
	}
	return "", 0
}
