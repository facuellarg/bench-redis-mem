package main

type MemoryRedis struct {
	data map[string]string
}

func NewMemoryRedis() *MemoryRedis {
	return &MemoryRedis{data: make(map[string]string)}
}

func (r *MemoryRedis) Get(key string) (string, error) {
	return r.data[key], nil
}

func (r *MemoryRedis) Set(key string, value string) error {
	r.data[key] = value
	return nil
}

func (r *MemoryRedis) Clean() error {
	r.data = make(map[string]string)
	return nil
}
