package variant

import (
	"math"
	"strconv"

	"github.com/garyburd/redigo/redis"
)

// Variant is a struct containing variant info. All data is immutable.
type Variant struct {
	UUID      string `json:"uuid"`
	Text      string `json:"text"`
	Image     string `json:"image"`
	Rating    uint32 `json:"rating"`
	CreatedAt int64  `json:"createdAt"`
	Author    string `json:"author"`
}

// redisVariant converts hash reply from redis to Variant struct (similar to redis.* methods)
func redisVariant(reply interface{}, v *Variant) error {
	fields, err := redis.StringMap(reply, nil)
	if err != nil {
		return err
	}

	rating, err := strconv.ParseUint(fields["rating"], 10, 32)
	if err != nil {
		return err
	}

	createdAt, err := strconv.ParseInt(fields["createdAt"], 10, 64)
	if err != nil {
		return err
	}

	v.Text = fields["text"]
	v.Image = fields["image"]
	v.Author = fields["author"]
	v.Rating = uint32(rating)
	v.CreatedAt = createdAt

	return nil
}

// GetByIds is a helper function which reads all variants by their UUIDs
func GetByIds(pool *redis.Pool, uuids []string) ([]*Variant, error) {
	conn := pool.Get()
	defer conn.Close()

	conn.Send("MULTI")
	for _, uuid := range uuids {
		conn.Send("HGETALL", "variants:"+uuid)
	}
	resp, err := redis.Values(conn.Do("EXEC"))
	if err != nil {
		return nil, err
	}

	variants := make([]*Variant, len(uuids))
	for i, reply := range resp {
		variant := &Variant{UUID: uuids[i]}
		err := redisVariant(reply, variant)
		if err != nil {
			return nil, err
		}

		variants[i] = variant
	}

	return variants, nil
}

// Win changes rating of two variants
func Win(pool *redis.Pool, win string, lose string) (int, int, error) {
	conn := pool.Get()
	defer conn.Close()

	conn.Send("MULTI")
	conn.Send("HGET", "variants:"+win, "rating")
	conn.Send("HGET", "variants:"+lose, "rating")
	ratings, err := redis.Values(conn.Do("EXEC"))
	if err != nil {
		return 0, 0, err
	}

	var oldWin int
	var oldLose int
	if _, err := redis.Scan(ratings, &oldWin, &oldLose); err != nil {
		return 0, 0, err
	}
	newWin, newLose := calculateElo(oldWin, oldLose)


	conn.Send("MULTI")
	conn.Send("HSET", "variants:"+win, "rating", newWin)
	conn.Send("HSET", "variants:"+lose, "rating", newLose)
	_, err = redis.Values(conn.Do("EXEC"))
	if err != nil {
		return 0, 0, nil
	}

	return newWin, newLose, nil
}

func calculateElo(oldWin int, oldLose int) (int, int) {
	chance := 1 / (1 + math.Pow(10, float64(oldLose-oldWin)/400))
	delta := int(math.Floor(32 * (1 - chance)))

	newWin := oldWin + delta
	newLose := int(math.Max(float64(oldLose-delta), 0))
	return newWin, newLose
}
