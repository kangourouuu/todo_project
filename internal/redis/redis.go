package redis

import (
	"context"
	"todo_project/common/log"
	"time"

	"github.com/redis/go-redis/v9"
    "github.com/sirupsen/logrus"
)

var ctx = context.Background()

type IRedis interface {
	Connect() error
	Ping() error
	GetClient() *redis.Client
	Set(key string, value interface{}) (string, error)
	Get(key string) (string, error)
	Delete(key string) (int64, error)
}

var Redis IRedis

type RedisClient struct {
	Client   *redis.Client
	rdConfig RedisConfig
}

type RedisConfig struct {
	Addr         string
	Password     string
	DB           int
	Poolsize     int
	PoolTimeOut  int
	IdleTimeOut  int
	ReadTimeOut  int
	WriteTimeOut int
}

func NewRedis(config RedisConfig) (IRedis, error) {
	r := &RedisClient{
		rdConfig: config,
	}
	err := r.Connect()
	if err != nil {
		return nil, err
	}
	return r, nil
}

func (r *RedisClient) Connect() error {
	Client := redis.NewClient(&redis.Options{
		Addr:         r.rdConfig.Addr,
		Password:     r.rdConfig.Password,
		DB:           r.rdConfig.DB,
		PoolSize:     15,
		PoolTimeout:  time.Duration(r.rdConfig.PoolTimeOut) * time.Second,
		ReadTimeout:  time.Duration(r.rdConfig.ReadTimeOut) * time.Second,
		WriteTimeout: time.Duration(r.rdConfig.WriteTimeOut) * time.Second,
	})
	str, err := Client.Ping(ctx).Result()
	if err != nil {
		log.Fatal(err)
		return err
	}
	logrus.Info(str)
	r.Client = Client
	return nil
}

func (r *RedisClient) GetClient() *redis.Client {
	return r.Client
}

func (r *RedisClient) Ping() error {
	_, err := r.Client.Ping(ctx).Result()
	return err
}

func (r *RedisClient) Set(key string, value interface{}) (string, error) {
	ret, err := r.Client.Set(ctx, key, value, 0).Result()
	return ret, err
}

func (r *RedisClient) Get(key string) (string, error) {
	ret, err := r.Client.Get(ctx, key).Result()
	return ret, err
}

func (r *RedisClient) Delete(key string) (int64, error) {
	ret, err := r.Client.Del(ctx, key).Result()
	return ret, err
}
