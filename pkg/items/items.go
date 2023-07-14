package items

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/litepubl/test-game/pkg/entity"
	"github.com/nats-io/nats.go"
	"github.com/redis/go-redis/v9"
	"github.com/rs/zerolog/log"
)

const (
	redisKey      = "items.list"
	redisLifetime = time.Minute
)

// Service Items provide access items in database
type Items struct {
	repo  ItemsRepo
	redis *redis.Client
	nc    *nats.Conn
}

var _ CRUDItems = (*Items)(nil)

// Items constructor
func New(r ItemsRepo, redis *redis.Client, nc *nats.Conn) *Items {
	return &Items{
		repo:  r,
		redis: redis,
		nc:    nc,
	}
}

// Implementation of CRUDItems interface
func (it *Items) List(ctx context.Context) ([]byte, error) {
	v, err := it.redis.Get(ctx, redisKey).Result()
	if err == nil && err != redis.Nil && v != "" {
		return []byte(v), nil
	}

	if err != nil && err != redis.Nil {
		log.Warn().Err(err).Msg("items.list error read from redis")
	}

	list, err := it.repo.List(ctx)
	if err != nil {
		return nil, fmt.Errorf("error %w read items list from database", err)
	}

	json, err := json.Marshal(list)
	if err != nil {
		return nil, fmt.Errorf("items.list error encoding into json %w", err)
	}

	err = it.redis.Set(ctx, redisKey, json, redisLifetime).Err()
	if err != nil {
		log.Warn().Err(err).Msg("items.list error store into redis")
	}

	return json, nil
}

func (it *Items) Create(ctx context.Context, campaignId int, name string) (entity.Item, error) {
	p, err := it.repo.MaxPriority(ctx)
	if err != nil {
		return entity.Item{}, fmt.Errorf("items.create error get MaxPriority %w", err)
	}

	item := entity.Item{
		CampaignId:  campaignId,
		Name:        name,
		Description: "",
		Priority:    p + 1,
		Removed:     false,
		CreatedAt:   time.Now(),
	}

	err = it.repo.Create(ctx, &item)
	if err != nil {
		return item, fmt.Errorf("items.create error save to database %w", err)
	}

	it.invalidateRedis(ctx)
	it.publishNats(item)
	return item, nil
}

func (it *Items) Update(ctx context.Context, u entity.UpdateData) (entity.Item, error) {
	p, err := it.repo.MaxPriority(ctx)
	if err != nil {
		return entity.Item{}, fmt.Errorf("items.Update error get MaxPriority %w", err)
	}

	u.Priority = p + 1
	item, err := it.repo.Update(ctx, u)
	if err != nil {
		return item, fmt.Errorf("items.Update error from database %w", err)
	}

	it.invalidateRedis(ctx)
	it.publishNats(item)
	return item, nil
}

func (it *Items) Delete(ctx context.Context, id int) error {
	item, err := it.repo.Delete(ctx, id)
	if err != nil {
		return fmt.Errorf("items.Delete error from database %w", err)
	}

	it.invalidateRedis(ctx)
	it.publishNats(item)
	return nil
}

func (it *Items) invalidateRedis(ctx context.Context) {
	pipe := it.redis.Pipeline()
	pipe.Del(ctx, redisKey)
	pipe.Exec(ctx)
}

func (it *Items) publishNats(e entity.Item) {
	d := entity.NatsItem{
		Id:          e.Id,
		CampaignId:  e.CampaignId,
		Name:        e.Name,
		Description: e.Description,
		Priority:    e.Priority,
		Removed:     e.Removed,
		EventTime:   time.Now(),
	}

	json, err := json.Marshal(d)
	if err != nil {
		return
	}

	it.nc.Publish("item", json)
}
