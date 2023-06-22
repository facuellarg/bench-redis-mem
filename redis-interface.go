package main

type RedisInterface interface {
	Get(key string) (string, error)
	Set(key string, value string) error
	Clean() error
}
