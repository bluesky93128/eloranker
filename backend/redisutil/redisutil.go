package redisutil

import (
	"errors"
	"github.com/garyburd/redigo/redis"
)

// Error drops reply from redis and returns only error
func Error(_ interface{}, err error) error {
	return err
}

// HSetArgs is a helper function that converts map to HSET args slice
func HSetArgs(key string, fields map[string]interface{}) []interface{} {
	args := make([]interface{}, 1+len(fields)*2)
	args[0] = key
	i := 1

	for k, v := range fields {
		args[i] = k
		args[i+1] = v
		i += 2
	}

	return args
}

// ArrayToMap is a helper that converts an array of strings into a map[string]bool (value is always true).
func ArrayToMap(result interface{}, err error) (map[string]bool, error) {
	values, err := redis.Values(result, err)
	if err != nil {
		return nil, err
	}

	m := make(map[string]bool, len(values))
	for _, value := range values {
		key, okKey := value.([]byte)
		if !okKey {
			return nil, errors.New("SetToMap key not a bulk string value")
		}
		m[string(key)] = true
	}

	return m, nil
}
