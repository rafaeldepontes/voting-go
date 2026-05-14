package repository

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/rafaeldepontes/voting-go/internal/poll"
	"github.com/rafaeldepontes/voting-go/internal/poll/model"
	rdb "github.com/rafaeldepontes/voting-go/pkg/cache/redis"
	"github.com/redis/go-redis/v9"
)

type repository struct {
	db *redis.Client
}

func NewRepository() poll.Repository {
	return &repository{
		db: rdb.GetCache(),
	}
}

func pollKey(id string) string {
	return fmt.Sprintf("poll:%s", id)
}

// FindPollByID implements [poll.Repository].
func (r *repository) FindPollByID(ctx context.Context, pollID string) (model.Poll, error) {
	val, err := r.db.Get(ctx, pollKey(pollID)).Result()
	if err != nil {
		return model.Poll{}, err
	}

	var p model.Poll
	if err := json.Unmarshal([]byte(val), &p); err != nil {
		return model.Poll{}, err
	}

	return p, nil
}

// Insert implements [poll.Repository].
func (r *repository) Insert(ctx context.Context, p model.Poll) error {
	data, err := json.Marshal(p)
	if err != nil {
		return err
	}
	return r.db.Set(ctx, pollKey(p.ID), data, p.Duration).Err()
}

// ListPolls implements [poll.Repository].
func (r *repository) ListPolls(ctx context.Context) ([]model.Poll, error) {
	const batchSize = 100
	const match = "poll:*"

	var (
		cursor uint64
		polls  []model.Poll
	)

	for {
		keys, nextCursor, err := r.db.Scan(ctx, cursor, match, batchSize).Result()
		if err != nil {
			return nil, err
		}

		cursor = nextCursor

		if len(keys) > 0 {
			val, err := r.db.MGet(ctx, keys...).Result()
			if err != nil {
				return nil, err
			}

			for i := range val {
				if val[i] == nil {
					continue
				}

				str, ok := val[i].(string)
				if !ok {
					continue
				}

				var p model.Poll
				if err := json.Unmarshal([]byte(str), &p); err != nil {
					return nil, err
				}

				polls = append(polls, p)
			}
		}

		if cursor == 0 {
			break
		}
	}

	return polls, nil
}

// Update implements [poll.Repository].
func (r *repository) Update(ctx context.Context, p model.Poll) error {
	return r.Insert(ctx, p)
}
