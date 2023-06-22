package main

import "sync"

type MemoryRedis struct {
	data  map[string]string
	mutex *sync.Mutex
}

func NewMemoryRedis() *MemoryRedis {
	return &MemoryRedis{
		data:  make(map[string]string),
		mutex: &sync.Mutex{},
	}
}

func (r *MemoryRedis) Get(key string) (string, error) {
	return r.data[key], nil
}

func (r *MemoryRedis) Set(key string, value string) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()
	r.data[key] = value
	return nil
}

func (r *MemoryRedis) Clean() error {
	r.data = make(map[string]string)
	return nil
}
