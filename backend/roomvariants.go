package main

import (
	"errors"

	"github.com/ark120202/easy-elo-ranker/backend/redisutil"
	"github.com/ark120202/easy-elo-ranker/backend/variant"
	"github.com/garyburd/redigo/redis"
	"github.com/satori/go.uuid"
	"math"
	"time"
)

// AllocateNewVariant creates new entry in database and returns it's id
func (r *Room) AllocateNewVariant(author *Client) (string, error) {
	conn := pool.Get()
	defer conn.Close()

	UUID := uuid.NewV4().String()

	fields := map[string]interface{}{
		"rating":    1000,
		"createdAt": time.Now().UnixNano() / int64(time.Millisecond),
		"author":    author.getUniqueIdentifier(),
	}

	conn.Send("MULTI")
	conn.Send("SADD", "room:"+r.name+":variants", UUID)
	conn.Send("HSET", redisutil.HSetArgs("variants:"+UUID, fields)...)
	err := redisutil.Error(conn.Do("EXEC"))
	if err != nil {
		return "", err
	}

	fields["event"] = "variant:allocate"
	fields["uuid"] = UUID
	r.SendToEveryone(nil, fields)

	return UUID, nil
}

// RemoveVariant removes variant from room
func (r *Room) RemoveVariant(initiator *Client, id string) error {
	conn := pool.Get()
	defer conn.Close()

	conn.Send("MULTI")
	conn.Send("SREM", "room:"+r.name+":variants", id)
	conn.Send("DEL", "variants:"+id)
	err := redisutil.Error(conn.Do("EXEC"))
	if err != nil {
		return err
	}

	fields := map[string]interface{}{
		"event": "variant:remove",
		"id":    id,
	}
	r.SendToEveryone(initiator, fields)

	return nil
}

// UpdateVariant updates variant based on user input
func (r *Room) UpdateVariant(initiator *Client, message requestMessageDataUpdateVariant) error {
	conn := pool.Get()
	defer conn.Close()

	exists, err := redis.Bool(conn.Do("EXISTS", "variants:"+message.UUID))
	if err != nil {
		return err
	}
	if !exists {
		return errors.New("variant not exists")
	}

	fields := map[string]interface{}{
		"image": message.Image,
		"text":  message.Text,
	}
	err = redisutil.Error(conn.Do("HSET", redisutil.HSetArgs("variants:"+message.UUID, fields)...))
	if err != nil {
		return err
	}

	fields["event"] = "variant:update"
	fields["uuid"] = message.UUID
	r.SendToEveryone(initiator, fields)

	return err
}

func (r *Room) getVariantsLength() (int, error) {
	conn := pool.Get()
	defer conn.Close()

	return redis.Int(conn.Do("SCARD", "room:"+r.name+":variants"))
}

// GetMaxQuota returns current quota limit, based on variants amount
func (r *Room) GetMaxQuota() (int, error) {
	variants, err := r.getVariantsLength()
	if err != nil {
		return 0, err
	}
	return int(math.Pow(float64(variants-1), 2)), err
}

// GetVariantIDs returns all variants in the room
func (r *Room) GetVariantIDs() ([]string, error) {
	conn := pool.Get()
	defer conn.Close()

	return redis.Strings(conn.Do("SMEMBERS", "room:"+r.name+":variants"))
}

// GetVariants returns all variants in the room
func (r *Room) GetVariants() ([]*variant.Variant, error) {
	uuids, err := r.GetVariantIDs()
	if err != nil {
		return nil, err
	}

	return variant.GetByIds(pool, uuids)
}
