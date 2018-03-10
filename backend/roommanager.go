package main

import (
	"errors"
	"math/rand"

	"github.com/dchest/uniuri"

	"github.com/ark120202/easy-elo-ranker/backend/redisutil"
	"github.com/garyburd/redigo/redis"
)

var runningRooms = make(map[string]*Room)

const roomLetters = "qwerasdfzxcv1234567890"

func randomRoomName() string {
	b := make([]byte, 8)
	for i := range b {
		b[i] = roomLetters[rand.Intn(len(roomLetters))]
	}
	return string(b)
}

func randomRoomSecret() string {
	return uniuri.NewLen(16)
}

func createRoom(name string, secret string) (*Room, error) {
	if _, err := RoomByName(name); err == nil {
		return nil, errors.New("room already exists")
	}

	conn := pool.Get()
	defer conn.Close()

	err := redisutil.Error(conn.Do(
		"HSET", "room:"+name,
		"secret", secret,
		"editMode", EditModeNormal,
		"quotaEnabled", false,
		"title", "",
	))
	if err != nil {
		return nil, err
	}

	room := &Room{}
	runningRooms[name] = room
	room.Make(name)

	return room, nil
}

// CreateRandomRoom creates room with random unique name
func CreateRandomRoom() (*Room, string, error) {
	secret := randomRoomSecret()

	// Repeat creating until roll unique name
	for {
		name := randomRoomName()
		room, err := createRoom(name, secret)

		if err == nil {
			return room, secret, nil
		}

		if err.Error() != "room already exists" {
			return nil, "", err
		}
	}
}

// RoomByName returns room by it's name. If room is not running, but exists in db it runs it.
func RoomByName(name string) (*Room, error) {
	// Room is not running, try to find it in database
	if runningRooms[name] == nil {
		conn := pool.Get()
		defer conn.Close()

		exists, err := redis.Bool(conn.Do("EXISTS", "room:"+name))
		if err != nil {
			return nil, err
		}
		if !exists {
			return nil, errors.New("room not exists")
		}

		room := &Room{}
		runningRooms[name] = room
		room.Make(name)
	}

	return runningRooms[name], nil
}
