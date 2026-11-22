package config

import (
	"errors"
	"fmt"
	"os"
	"strconv"
)

type Config struct {
	DBURL     string
	Port      int
	JWTSecret string
}

func Load() (Config, error) {
	dbURL := os.Getenv("DB_URL")
	jwtSecret := os.Getenv("JWT_SECRET")
	port := os.Getenv("PORT")

	var err error

	if dbURL == "" {
		err = errors.Join(err, errors.New("DB_URL is empty"))
	}

	if jwtSecret == "" {
		err = errors.Join(err, errors.New("JWT_SECRET is empty"))
	} else if len(jwtSecret) < 32 {
		err = errors.Join(err, errors.New("JWT_SECRET is not secure: must be at least 32 bytes"))
	}

	var portNumeric int

	if port == "" {
		err = errors.Join(err, errors.New("PORT is empty"))
	} else {
		portNumeric, errNumeric := strconv.Atoi(port)
		if errNumeric != nil {
			err = errors.Join(err, fmt.Errorf("PORT must be integer: %w", errNumeric))
		} else if portNumeric < 1 || portNumeric > 65535 {
			err = errors.Join(err, errors.New("PORT must be between 1 and 65535"))
		}
	}

	if err != nil {
		return Config{}, err
	}

	return Config{
		DBURL:     dbURL,
		Port:      portNumeric,
		JWTSecret: jwtSecret,
	}, nil
}
