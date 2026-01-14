// Package redis implements locator plugin with Redis as backend.
package redis

import (
	"context"
	"fmt"

	"github.com/redis/go-redis/v9"

	"github.com/byteweap/wukong/component/locator"
)

// ID is Redis locator implementation identifier.
const ID = "redis(hash)"

// Locator implements locator.Locator using Redis hash structures.
type Locator struct {
	rc                redis.UniversalClient // Redis client for hash operations
	keyFormat         string                // Format string for constructing Redis keys
	gateNodeFieldName string                // Field name for gate node in hash structure
	gameNodeFieldName string                // Field name for game node in hash structure
}

// Ensure Locator implements locator.Locator interface
var _ locator.Locator = (*Locator)(nil)

// New creates Redis locator with redis client configuration.
func New(opts redis.UniversalOptions, keyFormat, gateNodeFieldName, gameNodeFieldName string) *Locator {
	return newWith(redis.NewUniversalClient(&opts), keyFormat, gateNodeFieldName, gameNodeFieldName)
}

// newWith creates Redis locator with redis client.
func newWith(rc redis.UniversalClient, keyFormat, gateNodeFieldName, gameNodeFieldName string) *Locator {
	return &Locator{
		rc:                rc,
		keyFormat:         keyFormat,
		gateNodeFieldName: gateNodeFieldName,
		gameNodeFieldName: gameNodeFieldName,
	}
}

// ID returns locator implementation identifier.
func (l *Locator) ID() string {
	return ID
}

// Gate returns gate node for user ID.
func (l *Locator) Gate(ctx context.Context, uid int64) (string, error) {

	key := fmt.Sprintf(l.keyFormat, uid)
	return l.rc.HGet(ctx, key, l.gateNodeFieldName).Result()
}

// BindGate associates user ID with gate node.
func (l *Locator) BindGate(ctx context.Context, uid int64, node string) error {

	key := fmt.Sprintf(l.keyFormat, uid)
	return l.rc.HMSet(ctx, key, l.gateNodeFieldName, node).Err()
}

// UnBindGate removes user ID's gate node association if node matches.
func (l *Locator) UnBindGate(ctx context.Context, uid int64, node string) error {

	current, err := l.Gate(ctx, uid)
	if err != nil {
		return err
	}

	if current == node {
		key := fmt.Sprintf(l.keyFormat, uid)
		if err = l.rc.HMSet(ctx, key, l.gateNodeFieldName, "").Err(); err != nil {
			return err
		}
	}
	return nil
}

// Game returns game node for user ID.
func (l *Locator) Game(ctx context.Context, uid int64) (string, error) {

	key := fmt.Sprintf(l.keyFormat, uid)
	return l.rc.HGet(ctx, key, l.gameNodeFieldName).Result()
}

// BindGame associates user ID with game node.
func (l *Locator) BindGame(ctx context.Context, uid int64, node string) error {

	key := fmt.Sprintf(l.keyFormat, uid)
	return l.rc.HMSet(ctx, key, l.gameNodeFieldName, node).Err()
}

// UnBindGame removes user ID's game node association if node matches.
func (l *Locator) UnBindGame(ctx context.Context, uid int64, node string) error {

	current, err := l.Game(ctx, uid)
	if err != nil {
		return err
	}

	if current == node {
		key := fmt.Sprintf(l.keyFormat, uid)
		if err = l.rc.HMSet(ctx, key, l.gameNodeFieldName, "").Err(); err != nil {
			return err
		}
	}
	return nil
}

func (l *Locator) Close() error {
	return l.rc.Close()
}
