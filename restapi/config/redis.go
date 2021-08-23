package config

import (
	"context"
	jwt "github.com/dgrijalva/jwt-go"
	"github.com/go-redis/redis/v8"
	"github.com/google/uuid"
	"log"
	"os"
	"time"
)

type redisClient struct {
	c *redis.Client
}

type TokenDetails struct {
	AccessToken  string
	RefreshToken string
	AccessUuid   string
	RefreshUuid  string
	AtExpires    int64
	RtExpires    int64
}

var (
	ctx    = context.Background()
	client = &redisClient{}

	JWT_SECRET_KEY = []byte(os.Getenv("JWT_SECRET_KEY"))

	redisAddr = os.Getenv("REDIS_DB_HOST") + ":6379"
)

func InitRedis() *redisClient {
	// initalizing redis
	c := redis.NewClient(&redis.Options{
		Addr:     redisAddr,
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	_, err := c.Ping(ctx).Result()
	if err != nil {
		log.Fatal("Failed connect to redis.")
	}

	client.c = c
	return client
}

func CreateToken(user_id uint64) (*TokenDetails, error) {
	td := &TokenDetails{}

	td.AtExpires = time.Now().Add(time.Minute * 15).Unix() // set token expires after 15 minutes.
	td.AccessUuid = uuid.New().String()

	td.RtExpires = time.Now().Add(time.Hour * 24 * 7).Unix() // set token expires after 7 days.
	td.RefreshUuid = uuid.New().String()

	var err error

	// access token
	atClaims := jwt.MapClaims{}
	atClaims["jti"] = td.RefreshUuid
	atClaims["identity"] = user_id
	atClaims["exp"] = td.AtExpires
	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, atClaims)
	td.AccessToken, err = accessToken.SignedString(JWT_SECRET_KEY)
	if err != nil {
		return nil, err
	}

	// refresh token
	rtClaims := jwt.MapClaims{}
	rtClaims["jti"] = td.RefreshUuid
	rtClaims["identity"] = user_id
	rtClaims["exp"] = td.RtExpires
	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, rtClaims)
	td.RefreshToken, err = refreshToken.SignedString(JWT_SECRET_KEY)
	if err != nil {
		return nil, err
	}

	return td, nil
}
