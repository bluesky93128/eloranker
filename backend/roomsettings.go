package main

import (
	"errors"
	"github.com/ark120202/easy-elo-ranker/backend/redisutil"
	"github.com/garyburd/redigo/redis"
)

func (r *Room) getSetting(name string) (interface{}, error) {
	conn := pool.Get()
	defer conn.Close()

	return conn.Do("HGET", "room:"+r.name, name)
}

func (r *Room) setSetting(name string, value interface{}) error {
	conn := pool.Get()
	defer conn.Close()

	return redisutil.Error(conn.Do("HSET", "room:"+r.name, name, value))
}

// EditMode is a room setting, which controls main rules of room
type EditMode int8

// EditMode is a room setting, which controls main rules of room
const (
	EditModeTrust EditMode = iota
	EditModeNormal
	EditModeRestricted
)

// IsValid makes sure that EditMode has valid value
func (e EditMode) IsValid() bool {
	return e == EditModeTrust || e == EditModeNormal || e == EditModeRestricted
}

// GetEditMode returns edit mode for this room
func (r *Room) GetEditMode() (EditMode, error) {
	value, err := redis.Int(r.getSetting("editMode"))
	return EditMode(value), err
}

// SetEditMode updates edit mode setting of the room
func (r *Room) SetEditMode(editMode EditMode) error {
	if !editMode.IsValid() {
		return errors.New("invalid edit mode")
	}

	return r.setSetting("editMode", editMode)
}

// IsQuotaEnabled returns true if quota is enabled for this room
func (r *Room) IsQuotaEnabled() (bool, error) {
	return redis.Bool(r.getSetting("quotaEnabled"))
}

// SetQuotaEnabled updates quota setting of the room
func (r *Room) SetQuotaEnabled(value bool) error {
	return r.setSetting("quotaEnabled", value)
}

// GetTitle returns current room title
func (r *Room) GetTitle() (string, error) {
	return redis.String(r.getSetting("title"))
}

// SetTitle updates room title
func (r *Room) SetTitle(value string) error {
	return r.setSetting("title", value)
}
