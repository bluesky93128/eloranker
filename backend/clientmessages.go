package main

import (
	"encoding/json"

	"errors"
	"github.com/ark120202/easy-elo-ranker/backend/variant"
	"github.com/garyburd/redigo/redis"
	"sync"
)

// Send JSON data to client
func (c *Client) Send(data interface{}) {
	message, _ := json.Marshal(data)
	c.sent <- message
}

func (c *Client) Error(message string, event string) {
	c.Send(map[string]interface{}{
		"event": event,
		"error": message,
	})
}

type requestMessage struct {
	Type string          `json:"type"`
	Data json.RawMessage `json:"data"`
}

type requestMessageDataNewRoom struct {
	Title string `json:"title"`
}

type requestMessageDataJoinRoom struct {
	Name   string `json:"name"`
	Secret string `json:"secret"`
}

type requestMessageDataUpdateVariant struct {
	UUID  string `json:"uuid"`
	Text  string `json:"text"`
	Image string `json:"image"`
}

type requestMessageDataSetVariantIgnored struct {
	UUID    string `json:"uuid"`
	Ignored bool   `json:"ignored"`
}

type requestMessageDataRemoveVariant struct {
	ID string `json:"id"`
}

type requestMessageDataSubmitVoting struct {
	UUID string `json:"uuid"`
}

func (c *Client) handleMessage(request *requestMessage) {
	switch request.Type {
	case "room:new":
		var message requestMessageDataNewRoom
		json.Unmarshal(request.Data, &message)
		c.newRoom(message)
	case "room:join":
		var message requestMessageDataJoinRoom
		json.Unmarshal(request.Data, &message)
		c.joinRoom(message)
	case "room:leave":
		c.leaveRoom()
	case "variant:allocate":
		c.allocateVariant()
	case "variant:update":
		var message requestMessageDataUpdateVariant
		json.Unmarshal(request.Data, &message)
		c.updateVariant(message)
	case "variant:setIgnored":
		var message requestMessageDataSetVariantIgnored
		json.Unmarshal(request.Data, &message)
		c.setVariantIgnored(message)
	case "variant:remove":
		var message requestMessageDataRemoveVariant
		json.Unmarshal(request.Data, &message)
		c.removeVariant(message.ID)
	case "voting:get":
		c.getVoting()
	case "voting:submit":
		var message requestMessageDataSubmitVoting
		json.Unmarshal(request.Data, &message)
		c.submitVote(message.UUID)
	case "settings:title":
		var message struct {
			Value string `json:"value"`
		}
		json.Unmarshal(request.Data, &message)
		c.setTitle(message.Value)
	case "settings:quotaEnabled":
		var message *struct {
			Value bool `json:"value"`
		}
		json.Unmarshal(request.Data, &message)
		c.setQuotaEnabled(message.Value)
	case "settings:editMode":
		var message *struct {
			Value EditMode `json:"value"`
		}
		json.Unmarshal(request.Data, &message)
		if message.Value.IsValid() {
			c.setEditMode(message.Value)
		} else {
			c.Error("invalid edit mode value", "settings:editMode")
		}
	default:
		c.Error("Incorrect type of message: "+request.Type, "")
	}
}

func (c *Client) newRoom(message requestMessageDataNewRoom) {
	room, secret, err := CreateRandomRoom()
	if err != nil {
		c.Error(err.Error(), "room:new")
		return
	}

	room.SetTitle(message.Title)

	c.Send(map[string]interface{}{
		"event":  "room:new",
		"name":   room.name,
		"secret": secret,
	})
}

func runGroup(functions ...func()) {
	var waitGroup sync.WaitGroup
	waitGroup.Add(len(functions))

	defer waitGroup.Wait()

	for _, function := range functions {
		go func(copy func()) {
			defer waitGroup.Done()
			copy()
		}(function)
	}
}

func (c *Client) joinRoom(message requestMessageDataJoinRoom) {
	room, err := RoomByName(message.Name)
	if err != nil {
		c.Error(err.Error(), "room:join")
		return
	}

	variants, err := room.GetVariants()
	if err != nil {
		c.Error(err.Error(), "room:join")
		return
	}

	c.room = room
	c.secret = message.Secret

	var isAdmin bool
	var title string
	var quotaEnabled bool
	var editMode EditMode
	var ignoredVariants map[string]bool
	var err1, err2, err3, err4, err5 error

	runGroup(
		func() { isAdmin, err1 = c.IsAdmin() },
		func() { title, err2 = room.GetTitle() },
		func() { quotaEnabled, err3 = room.IsQuotaEnabled() },
		func() { editMode, err4 = room.GetEditMode() },
		func() { ignoredVariants, err5 = c.getIgnoredVariants() },
	)

	if err1 != nil || err2 != nil || err3 != nil || err4 != nil || err5 != nil {
		c.Error("internal error: couldn't read room data", "room:join")
		c.room = nil
		c.secret = ""
		return
	}

	room.register <- c

	c.Send(map[string]interface{}{
		"event":           "room:join",
		"variants":        variants,
		"identifier":      c.getUniqueIdentifier(),
		"ignoredVariants": ignoredVariants,

		"isAdmin":      isAdmin,
		"title":        title,
		"quotaEnabled": quotaEnabled,
		"editMode":     editMode,
	})
}

func (c *Client) leaveRoom() {
	if c.room != nil {
		c.room.unregister <- c
	}

	c.room = nil
	c.secret = ""
}

func (c *Client) allocateVariant() {
	if c.room == nil {
		c.Error("you should be in room to allocate variant", "variant:allocate")
		return
	}

	canWrite, err := c.hasWriteAccess("")
	if err != nil {
		c.Error(err.Error(), "variant:update")
		return
	}
	if !canWrite {
		c.Error("not enough permissions to allocate variant", "variant:update")
		return
	}

	_, err = c.room.AllocateNewVariant(c)
	if err != nil {
		c.Error(err.Error(), "variant:allocate")
		return
	}
}

func (c *Client) updateVariant(message requestMessageDataUpdateVariant) {
	if c.room == nil {
		c.Error("you should be in room to update variants", "variant:update")
		return
	}

	canWrite, err := c.hasWriteAccess(message.UUID)
	if err != nil {
		c.Error(err.Error(), "variant:update")
		return
	}
	if !canWrite {
		c.Error("not enough permissions to edit variant "+message.UUID, "variant:update")
		return
	}

	err = c.room.UpdateVariant(c, message)
	if err != nil {
		c.Error(err.Error(), "variant:update")
	}
}

func (c *Client) setVariantIgnored(message requestMessageDataSetVariantIgnored) {
	if c.room == nil {
		c.Error("you should be in room to ignore variants", "variant:setIgnored")
		return
	}
	if message.UUID == "" {
		c.Error("invalid variant id", "variant:setIgnored")
		return
	}

	conn := pool.Get()
	defer conn.Close()

	ignoredKey := "room:" + c.room.name + ":ignored:" + c.getUniqueIdentifier()

	isMember, err := redis.Bool(conn.Do("SISMEMBER", ignoredKey, message.UUID))
	if err != nil {
		c.Error(err.Error(), "variant:setIgnored")
		return
	}
	if isMember == message.Ignored {
		c.Error(message.UUID+" is already in the same state", "variant:setIgnored")
		return
	}

	if message.Ignored {
		variants, err := c.room.getVariantsLength()
		if err != nil {
			c.Error(err.Error(), "variant:setIgnored")
			return
		}
		if variants < 6 {
			c.Error("6 variants or more are required to ignore", "variant:setIgnored")
			return
		}

		ignoredVariants, err := c.getIgnoredVariantsLength()
		if err != nil {
			c.Error(err.Error(), "variant:setIgnored")
			return
		}

		if ignoredVariants*2 >= variants {
			c.Error("you can't ignore more variants", "variant:setIgnored")
			return
		}
	}

	if message.Ignored {
		conn.Do("SADD", ignoredKey, message.UUID)
	} else {
		conn.Do("SREM", ignoredKey, message.UUID)
	}
}

func (c *Client) removeVariant(id string) {
	if c.room == nil {
		c.Error("you should be in room to allocate variant", "variant:remove")
		return
	}

	canWrite, err := c.hasWriteAccess(id)
	if err != nil {
		c.Error(err.Error(), "variant:remove")
		return
	}
	if !canWrite {
		c.Error("not enough permissions to remove variant", "variant:remove")
		return
	}

	c.room.RemoveVariant(c, id)
}

func (c *Client) getVoting() {
	if c.room == nil {
		c.Error("you should be in room to vote", "voting:get")
		return
	}

	variants, err := c.SelectRandomPair()
	if err != nil {
		c.Error(err.Error(), "voting:get")
		return
	}

	c.voting[0] = variants[0]
	c.voting[1] = variants[1]

	c.Send(map[string]interface{}{
		"event":    "voting:get",
		"variants": variants,
	})
}

func (c *Client) submitVote(uuid string) {
	quotaOk, err := c.useVotingQuota()
	if err != nil {
		c.Error(err.Error(), "voting:submit")
		return
	}
	if !quotaOk {
		c.Error("quota reached", "voting:submit")
		return
	}

	var variantWin string
	var variantLose string
	switch {
	case c.voting[0] == uuid:
		variantWin = c.voting[0]
		variantLose = c.voting[1]
	case c.voting[1] == uuid:
		variantWin = c.voting[1]
		variantLose = c.voting[0]
	default:
		c.Error("incorrect uuid", "voting:submit")
		return
	}

	ratingWin, ratingLose, err := variant.Win(pool, variantWin, variantLose)
	if err != nil {
		c.Error(err.Error(), "voting:submit")
		return
	}

	c.room.SendToEveryone(nil, map[string]interface{}{
		"event":  "variant:update",
		"uuid":   variantWin,
		"rating": ratingWin,
	})
	c.room.SendToEveryone(nil, map[string]interface{}{
		"event":  "variant:update",
		"uuid":   variantLose,
		"rating": ratingLose,
	})

	// Make sure client can't vote multiple times for this pair
	c.voting[0] = ""
	c.voting[1] = ""
}

func (c *Client) canUpdateSettings() error {
	if c.room == nil {
		return errors.New("you should be in room to update settings")
	}

	isAdmin, err := c.IsAdmin()
	if err != nil {
		return err
	}
	if !isAdmin {
		return errors.New("you should be admin to update room settings")
	}

	return nil
}

func (c *Client) setTitle(value string) {
	err := c.canUpdateSettings()
	if err != nil {
		c.Error(err.Error(), "settings:title")
		return
	}

	err = c.room.SetTitle(value)
	if err != nil {
		c.Error(err.Error(), "settings:title")
		return
	}

	c.room.SendToEveryone(c, map[string]interface{}{
		"event": "settings:title",
		"value": value,
	})
}

func (c *Client) setQuotaEnabled(value bool) {
	err := c.canUpdateSettings()
	if err != nil {
		c.Error(err.Error(), "settings:quotaEnabled")
		return
	}

	err = c.room.SetQuotaEnabled(value)
	if err != nil {
		c.Error(err.Error(), "settings:quotaEnabled")
		return
	}

	c.room.SendToEveryone(c, map[string]interface{}{
		"event": "settings:quotaEnabled",
		"value": value,
	})
}

func (c *Client) setEditMode(value EditMode) {
	err := c.canUpdateSettings()
	if err != nil {
		c.Error(err.Error(), "settings:editMode")
		return
	}

	err = c.room.SetEditMode(value)
	if err != nil {
		c.Error(err.Error(), "settings:editMode")
		return
	}

	c.room.SendToEveryone(c, map[string]interface{}{
		"event": "settings:editMode",
		"value": value,
	})
}
