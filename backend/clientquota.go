package main

import (
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"github.com/ark120202/easy-elo-ranker/backend/redisutil"
	"github.com/garyburd/redigo/redis"
	"math/rand"
	"strings"
	"time"
)

const sessionLetters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

func generateSession() string {
	b := make([]byte, 32)
	for i := range b {
		b[i] = sessionLetters[rand.Intn(len(sessionLetters))]
	}
	return string(b)
}

func (c *Client) getUniqueIdentifier() string {
	hash := sha256.New()
	hash.Write([]byte(c.ip))

	// TODO: https://security.stackexchange.com/a/74115
	return base64.URLEncoding.EncodeToString(hash.Sum(nil))
}

func (c *Client) getIgnoredVariants() (map[string]bool, error) {
	conn := pool.Get()
	defer conn.Close()

	return map[string]bool{}, nil // redisutil.ArrayToMap(conn.Do("SMEMBERS", ""))
}

// SelectRandomPair returns a pair of variants for voting
func (c *Client) SelectRandomPair() ([]string, error) {
	conn := pool.Get()
	defer conn.Close()

	seenPairsKey := "room:" + c.room.name + ":seen:" + c.getUniqueIdentifier()

	// TODO: Benchmark it
	possibleVariants, err := c.room.GetVariantIDs()
	if err != nil {
		return nil, err
	}
	if len(possibleVariants) <= 1 {
		return nil, errors.New("not enough variants to vote")
	}

	seenPairs, err := redisutil.ArrayToMap(conn.Do("SMEMBERS", seenPairsKey))
	if err != nil {
		return nil, err
	}

	ignoredVariants, err := c.getIgnoredVariants()
	if err != nil {
		return nil, err
	}

	possiblePermutations := makePermutations(possibleVariants, seenPairs, ignoredVariants)

	if len(possiblePermutations) == 0 {
		err := redisutil.Error(conn.Do("DEL", seenPairsKey))
		if err != nil {
			return nil, err
		}

		// Pass empty table, since it was removed
		var seenPairs map[string]bool
		possiblePermutations = makePermutations(possibleVariants, seenPairs, ignoredVariants)
	}

	pair := possiblePermutations[rand.Intn(len(possiblePermutations))]

	err = redisutil.Error(conn.Do("SADD", seenPairsKey, pair))
	if err != nil {
		return nil, err
	}

	variantsInPair := strings.SplitN(pair, "+", 2)

	return variantsInPair, nil
}

func makePermutations(elements []string, ignoredPairs map[string]bool, ignoredVariants map[string]bool) []string {
	var permutations []string
	for i, variant := range elements {
		if _, ok := ignoredVariants[variant]; ok {
			continue
		}

		for j := i + 1; j < len(elements); j++ {
			if _, ok := ignoredVariants[elements[j]]; ok {
				continue
			}

			pair := variant + "+" + elements[j]

			if _, ok := ignoredPairs[pair]; !ok {
				permutations = append(permutations, pair)
			}
		}
	}
	return permutations
}

func (c *Client) useVotingQuota() (bool, error) {
	if c.room == nil {
		return false, errors.New("you should be in room to vote")
	}

	quotaEnabled, err := c.room.IsQuotaEnabled()
	if err != nil {
		return false, err
	}

	if !quotaEnabled {
		return true, nil
	}

	identifier := c.getUniqueIdentifier()

	votes, err := c.room.GetVotesInLast24Hours(identifier)
	if err != nil {
		return false, err
	}

	max, err := c.room.GetMaxQuota()
	if err != nil {
		return false, err
	}

	if votes >= max {
		return false, nil
	}

	conn := pool.Get()
	defer conn.Close()

	timestamp := time.Now().Unix()
	redisutil.Error(conn.Do("ZADD", "room:"+c.room.name+":votes:"+identifier, timestamp, timestamp))

	return true, nil
}

// IsAdmin returns true when client's secret is valid
func (c *Client) IsAdmin() (bool, error) {
	if c.room == nil {
		return false, errors.New("client should be in the room")
	}
	secret, err := c.room.GetSecret()
	if err != nil {
		return false, err
	}

	return c.secret == secret, nil
}

func (c *Client) hasWriteAccess(id string) (bool, error) {
	connection := pool.Get()
	defer connection.Close()

	editMode, err := c.room.GetEditMode()
	if err != nil {
		return false, err
	}

	if editMode == EditModeTrust {
		return true, nil
	}

	isAdmin, err := c.IsAdmin()
	if err != nil {
		return false, err
	}

	// Admins can write in any case
	if isAdmin {
		return true, nil
	}

	if editMode == EditModeNormal {
		// Allow to create variants to everyone
		if id == "" {
			return true, nil
		}

		variantAuthor, err := redis.String(connection.Do("HGET", "variants:"+id, "author"))
		if err != nil {
			return false, err
		}

		identifier := c.getUniqueIdentifier()

		return variantAuthor == identifier, nil
	}

	return false, nil
}
