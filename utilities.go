package main

import (
	"log"
	"math/rand"
	"os"
	"strconv"
	"strings"
)

func CheckErr(err error, args ...string) {
	if err != nil {
		// log.Println("Error")
		// log.Println("%q: %s", err, args)
		log.Println(err)
	}
}

func CheckErrPanic(err error) {
	if err != nil {
		log.Fatalln(err)
	}
}

// ###########################################################
// https://dev.to/craicoverflow/a-no-nonsense-guide-to-environment-variables-in-go-a2f

// Simple helper function to read an environment or return a default value
func GetEnv(key string, defaultVal string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}

	return defaultVal
}

// Simple helper function to read an environment variable into integer or return a default value
func GetEnvAsInt(name string, defaultVal int) int {
	valueStr := GetEnv(name, "")
	if value, err := strconv.Atoi(valueStr); err == nil {
		return value
	}

	return defaultVal
}

func GetEnvAsInt64(name string, defaultVal int64) int64 {
	valueStr := GetEnv(name, "")
	if value, err := strconv.ParseInt(valueStr, 10, 64); err == nil {
		return value
	}

	return defaultVal
}

// Helper to read an environment variable into a bool or return default value
func GetEnvAsBool(name string, defaultVal bool) bool {
	valStr := GetEnv(name, "")
	if val, err := strconv.ParseBool(valStr); err == nil {
		return val
	}

	return defaultVal
}

// Helper to read an environment variable into a string slice or return default value
func GetEnvAsSlice(name string, defaultVal []string, sep string) []string {
	valStr := GetEnv(name, "")

	if valStr == "" {
		return defaultVal
	}

	val := strings.Split(valStr, sep)

	return val
}

// https://henvic.dev/posts/uuid/
func NewID() string {
	const (
		alphabet = "123456789ABCDEFGHJKLMNPQRSTUVWXYZabcdefghijkmnopqrstuvwxyz" // base58
		size     = 11
	)

	var id = make([]byte, size)
	if _, err := rand.Read(id); err != nil {
		panic(err)
	}
	for i, p := range id {
		id[i] = alphabet[int(p)%len(alphabet)] // discard everything but the least significant bits
	}
	return string(id)
}
