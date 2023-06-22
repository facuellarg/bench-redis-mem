package main

import (
	"context"
	"crypto/tls"
	"fmt"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/memorydb"
	"github.com/go-redis/redis/v8"
)

var redisConSingleton *redis.ClusterClient

func GetRedisCon() (*redis.ClusterClient, error) {
	if redisConSingleton != nil {
		return redisConSingleton, nil
	}
	redisAddress := "localhost:6379"
	if os.Getenv("ENV") != "" {
		addr, port := memoryDB()
		redisAddress = fmt.Sprintf("%s:%d", addr, port)
	}
	fmt.Println(redisAddress)
	clusterOptions := &redis.ClusterOptions{
		Addrs: []string{
			redisAddress,
		},
		NewClient: func(opt *redis.Options) *redis.Client {
			opt.DB = 0
			opt.Password = ""
			opt.TLSConfig = &tls.Config{
				InsecureSkipVerify: false,
			}
			return redis.NewClient(opt)
		},
		TLSConfig: &tls.Config{
			InsecureSkipVerify: false,
		},
	}

	// Create a Redis redisConSingleton using the go-redis library
	// redisConSingleton = redis.NewClient(&redis.Options{
	redisConSingleton = redis.NewClusterClient(clusterOptions)
	// TLSConfig: &tls.Config{
	// 	InsecureSkipVerify: false,
	// },
	// })

	// Ping the Redis server to check the connectivity
	_, err := redisConSingleton.Ping(redisConSingleton.Context()).Result()
	if err != nil {
		return nil, err
	}

	return redisConSingleton, nil
}

func memoryDB() (string, int64) {
	// Create a session using your AWS credentials
	sess, err := session.NewSessionWithOptions(session.Options{
		Config: aws.Config{
			Region: aws.String("us-east-1"),
		},
		Profile: "personal",
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

	fmt.Printf("result: %v\n", result)
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
