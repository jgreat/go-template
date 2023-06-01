package main

import (
	"os"

	"github.com/sirupsen/logrus"
)

func setupLogging() {
	logLevel, err := logrus.ParseLevel(getEnv("LOG_LEVEL", "info"))
	if err != nil {
		log.Fatalf("Invalid LOG_LEVEL: %v\n", err)
	}
	log.SetLevel(logLevel)
	log.SetFormatter(&logrus.JSONFormatter{})
}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}
