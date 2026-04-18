package main

import (
	"fmt"
	"log"
	"strconv"

	"github.com/brianvoe/gofakeit/v7"
	"github.com/gomodule/redigo/redis"
)

const (
	fieldName     = "name"
	fieldLastName = "last_name"
	fieldAge      = "age"
	fieldEmail    = "email"
)

type User struct {
	Name     string `redis:"name"`
	LastName string `redis:"last_name"`
	Age      int    `redis:"age"`
	Email    string `redis:"email"`
}

func main() {
	conn, err := redis.Dial("tcp", "localhost:6379")
	if err != nil {
		log.Printf("failed to connect to Redis: %v\n", err)
		return
	}

	defer func() {
		if err := conn.Close(); err != nil {
			log.Printf("failed to close Redis connection: %v\n", err)
		}
	}()

	setAndGet(conn)
}

func setAndGet(conn redis.Conn) {
	key := gofakeit.UUID()
	value := gofakeit.FirstName()
	_, err := conn.Do("SET", key, value)
	if err != nil {
		log.Printf("failed to set key: %v\n", err)
		return
	}

	value, err = redis.String(conn.Do("GET", key))
	if err != nil {
		log.Printf("failed to get key: %v\n", err)
		return
	}

	fmt.Printf("Пара ключ-значение (%s: %s)\n\n", key, value)
	hsetAndGet(conn)
}

func hsetAndGet(conn redis.Conn) {
	hashKey := gofakeit.UUID()
	fields := map[string]string{
		fieldName:     gofakeit.FirstName(),
		fieldLastName: gofakeit.LastName(),
		fieldAge:      strconv.FormatInt(int64(gofakeit.IntRange(0, 100)), 10),
		fieldEmail:    gofakeit.Email(),
	}

	var err error
	for field, value := range fields {
		_, err = conn.Do("HSET", hashKey, field, value)
		if err != nil {
			log.Printf("failed to set hash field: %v\n", err)
			return
		}
	}

	// Получаем значения из хеш-таблицы разными способами
	printMapFieldsByOne(conn, hashKey)
	fmt.Println()
	printMapFields(conn, hashKey)
	fmt.Println()
}

func printMapFieldsByStruct(conn redis.Conn, hash string) {
	values, err := redis.Values(conn.Do("HGETALL", hash))
	if err != nil {
		log.Printf("failed to get all hash fields: %v\n", err)
		return
	}

	var user User
	err = redis.ScanStruct(values, &user)
	if err != nil {
		log.Printf("failed to scan hash fields to struct: %v\n", err)
		return
	}
	fmt.Printf("Данные пользователя с идентифкатором (распаршенные в структуру) %s:\n", hash)
	fmt.Printf("%#v\n", user)
}

func printMapFields(conn redis.Conn, hashKey string) {
	hashMap, err := redis.StringMap(conn.Do("HGETALL", hashKey))
	if err != nil {
		log.Printf("failed to get all hash fields: %v\n", err)
		return
	}
	fmt.Printf("Данные пользователя с идентифкатором (полученные разом) %s:\n", hashKey)
	fmt.Printf("%#v\n", hashMap)
}

func printMapFieldsByOne(conn redis.Conn, hashKey string) {
	name, err := redis.String(conn.Do("HGET", hashKey, fieldName))
	if err != nil {
		log.Printf("failed to get hash field \"%v\": %v\n", fieldName, err)
		return
	}

	lastName, err := redis.String(conn.Do("HGET", hashKey, fieldLastName))
	if err != nil {
		log.Printf("failed to get hash field \"%v\": %v\n", fieldLastName, err)
		return
	}

	age, err := redis.String(conn.Do("HGET", hashKey, fieldAge))
	if err != nil {
		log.Printf("failed to get hash field \"%v\": %v\n", fieldAge, err)
		return
	}

	email, err := redis.String(conn.Do("HGET", hashKey, fieldEmail))
	if err != nil {
		log.Printf("failed to get hash field \"%v\": %v\n", fieldEmail, err)
		return
	}

	fmt.Printf("Данные пользователя с идентифкатором %s:\n", hashKey)
	fmt.Printf("Имя: %s\n", name)
	fmt.Printf("Фамилия: %s\n", lastName)
	fmt.Printf("Возраст: %s\n", age)
	fmt.Printf("Email: %s\n", email)
}
