package main

import (
	"time"
	"github.com/garyburd/redigo/redis"
)

// GetVotesInLast24Hours returns count of votes made by client in latest 24 hours
func (r *Room) GetVotesInLast24Hours(token string) (int, error) {
	conn := pool.Get()
	defer conn.Close()

	return redis.Int(conn.Do("ZCOUNT", "room:"+r.name+":votes:"+token,
		time.Now().AddDate(0, 0, -1).Unix(),
		time.Now().Unix(),
	))
}