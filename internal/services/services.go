package services

import (
	"context"
	"crypto/ed25519"
	"crypto/x509"
	"errors"
	"fmt"
	"net/url"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
)

var (
	ErrKeyNotFound = errors.New("key not found")
)

var Pool *pgxpool.Pool
var RDB *redis.Client
var IsTesting bool
var PrivateKey ed25519.PrivateKey
var PublicKey ed25519.PublicKey

func getRedisUrl() string {
	if os.Getenv("CONTAINERIZED") == "true" {
		return "redis:6379"
	} else {
		return "localhost:6379"
	}
}

func GetPGConnectionString() string {
	var domain string
	if os.Getenv("CONTAINERIZED") == "true" {
		domain = "db"
	} else {
		domain = "localhost"
	}

	return fmt.Sprintf("postgresql://%s:%s@%s:5432/%s?pool_max_conns=100",
		os.Getenv("POSTGRES_USER"), url.QueryEscape(os.Getenv("POSTGRES_PASSWORD")), domain, os.Getenv("POSTGRES_DB"))
}

func getPrivateKey() (key ed25519.PrivateKey, err error) {
	rawKey, err := os.ReadFile("./id_ed25519")
	if err != nil {
		return
	}
	privateKey, err := x509.ParsePKCS8PrivateKey(rawKey)
	if err != nil {
		return
	}
	key = privateKey.(ed25519.PrivateKey)
	return
}

func getPublicKey() (key ed25519.PublicKey, err error) {
	rawKey, err := os.ReadFile("./id_ed25519.pub")
	if err != nil {
		return
	}
	publicKey, err := x509.ParsePKIXPublicKey(rawKey)
	if err != nil {
		return
	}
	key = publicKey.(ed25519.PublicKey)
	return
}

func GenerateKeys() error {
	publicKeyRaw, privateKeyRaw, err := ed25519.GenerateKey(nil)
	if err != nil {
		return err
	}
	privateKey, err := x509.MarshalPKCS8PrivateKey(privateKeyRaw)
	if err != nil {
		return err
	}
	publicKey, err := x509.MarshalPKIXPublicKey(publicKeyRaw)
	if err != nil {
		return err
	}

	err = os.WriteFile("./id_ed25519", privateKey, 0666)
	if err != nil {
		return err
	}
	err = os.WriteFile("./id_ed25519.pub", publicKey, 0666)
	if err != nil {
		return err
	}
	return nil
}

func Start() (err error) {
	if os.Getenv("TESTING") == "true" {
		IsTesting = true
	}

	PublicKey, err = getPublicKey()
	if err != nil {
		return
	}

	PrivateKey, err = getPrivateKey()
	if err != nil {
		return
	}

	Pool, err = pgxpool.New(context.Background(), GetPGConnectionString())
	if err != nil {
		return
	}

	RDB = redis.NewClient(&redis.Options{
		Addr:     getRedisUrl(),
		Password: os.Getenv("REDIS_PASSWORD"),
		DB:       0,
	})
	return
}

func Stop() {
	if Pool != nil {
		Pool.Close()
	}
}
