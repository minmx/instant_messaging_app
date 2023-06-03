package main

import (
    "context"
    "encoding/json"
    "github.com/redis/go-redis/v9"
    "time"
)

type RedisClient struct {
    client *redis.Client
}

func (c *RedisClient) InitClient(ctx context.Context, address, password string) error {
    r := redis.NewClient(&redis.Options{
        Addr: address,
        Password: password,
        DB: 0,
    })

    // test connection
    if err := r.Ping(ctx).Err(); err != nil {
        return err
    }

    c.client = r
    return nil
}

type Message struct {
    Sender string `json:"sender"`
    Message string `json:"message"`
    Timestamp int64 `json:"timestamp"`
}

func (c *RedisClient) SaveMessage(ctx context.Context, roomID string, message *Message) error {
    text, err := json.Marshal(message)
    if err != nil {
       return err
    }

    member := &redis.Z{
        Score:  float64(message.Timestamp,)
        Member: text,
    }

    e, err = c.client.ZAdd(ctx, roomID, *member).Result()
    if err != nil {
        return err
    }

    return nil
}

func (c *RedisClient) GetMessage(ctx context.Context, roomID string, start, end int64, reverse bool) ([]*Message, error ) {
    var (
        strMessages []string
        messages []*Message
        err error
    )
    if reverse {
        strMessages, err = c.client.ZRevRange(ctx, roomID, start, end).Result()
        if err != nil {
            return nil, err
        }
    }
    else {
        strMessages, err = c.client.ZRange(ctx, roomID, start, end).Result()
        if err != nil {
            return nil, err
        }
    }
    for _, msg := range strMessages {
        temp := &Message{}
        err := json.Unmarshal([]byte(msg), temp)
        if err != nil {
            return nil, err
        }
        messages = append(messages, temp)
    }
    return messages, nil
}