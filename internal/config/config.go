package config

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"log"
	"os"
	"strconv"
)

const (
	// env names for db config
	dbPortEnv = "DATABASE_PORT"
	dbUserEnv = "DATABASE_USER"
	dbPassEnv = "DATABASE_PASSWORD"
	dbNameEnv = "DATABASE_NAME"
	dbHostEnv = "DATABASE_HOST"

	// env names for srv config
	srvPortEnv = "SERVER_PORT" 

	secretKeyLen = 16
)

type Config struct {
	Db ConfigDB
	Srv ConfigSrv
	JWT ConfigJWT
}

type ConfigSrv struct {
	SrvPort int
}

type ConfigDB struct {
	DbPort int
	DbUser string
	DbPass string
	DbName string
	DbHost string
}

type ConfigJWT struct {
	Secret string
}

func MustLoad() *Config {
	dbPortStr := getStringOrDefault(dbPortEnv, "5432")
	dbPort, err := strconv.Atoi(dbPortStr)
	if err != nil {
		log.Fatalf("invalid database port: %s", err)
	}

	dbUser := getStringOrDefault(dbUserEnv, "postgres")
	dbPass := getStringOrDefault(dbPassEnv, "password")
	dbName := getStringOrDefault(dbNameEnv, "shop")
	dbHost := getStringOrDefault(dbHostEnv, "db")
	
	srvPortStr := getStringOrDefault(srvPortEnv, "8080")
	srvPort, err := strconv.Atoi(srvPortStr)
	if err != nil {
		log.Fatalf("invalid server port: %s", err)
	}

	secret, err := generateSecretKey()
	if err != nil {
		log.Fatal(err)
	}
	
	return &Config{
		Db: ConfigDB{
			DbPort: dbPort,
			DbUser: dbUser,
			DbPass: dbPass,
			DbName: dbName,
			DbHost: dbHost,
		},	
		Srv: ConfigSrv{
			SrvPort: srvPort,
		},	
		JWT: ConfigJWT{
			Secret: secret,
		},
	}
}

func getStringOrDefault(name, defaultVal string) string {
	res := os.Getenv(name)
	if res == "" {
		return defaultVal
	}

	return res
}

func generateSecretKey() (string, error) {
	bytes := make([]byte, secretKeyLen)
	_, err := rand.Read(bytes)
	if err != nil {
		return "", fmt.Errorf("cannot generate secret: %s", err)
	}
	return base64.URLEncoding.EncodeToString(bytes), nil
}