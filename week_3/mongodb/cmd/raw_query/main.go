package main

import (
	"context"
	"log"
	"os"
	"time"

	"github.com/brianvoe/gofakeit"
	"github.com/joho/godotenv"
	"github.com/mbakhodurov/examples/week_3/mongo/internal/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func main() {
	ctx := context.Background()

	err := godotenv.Load("../../.env")
	if err != nil {
		log.Printf("Не удалось загрузить файл .env: %v\n", err)
		return
	}

	dbURI := os.Getenv("MONGO_URI")
	if dbURI == "" {
		log.Println("Ошибка: переменная окружения MONGO_URI не установлена")
		return
	}

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(dbURI))
	if err != nil {
		log.Printf("Ошибка подключения к MongoDB: %v\n", err)
		return
	}

	defer func() {
		if cerr := client.Disconnect(ctx); cerr != nil {
			log.Printf("Ошибка при отключении от MongoDB: %v\n", cerr)
		}
	}()

	err = client.Ping(ctx, nil)
	if err != nil {
		log.Printf("MongoDB недоступна, ошибка ping: %v\n", err)
		return
	}
	log.Println("Успешное подключение к MongoDB")

	collection := client.Database("example").Collection("notes")

	indexModel := mongo.IndexModel{
		Keys:    bson.D{{Key: "title", Value: 1}},
		Options: options.Index().SetUnique(false),
	}

	indexName, err := collection.Indexes().CreateOne(ctx, indexModel)
	if err != nil {
		log.Printf("Ошибка создания индекса: %v\n", err)
		return
	}
	log.Printf("Создан индекс: %s\n", indexName)

	// note := model.Note{
	// 	Title:     gofakeit.City(),
	// 	Body:      gofakeit.Address().Street,
	// 	CreatedAt: time.Now(),
	// }

	// insertResult, err := collection.InsertOne(ctx, note)
	// if err != nil {
	// 	log.Printf("Ошибка вставки заметки: %v\n", err)
	// 	return
	// }

	// insertedID := insertResult.InsertedID
	// log.Printf("Заметка успешно добавлена, ID: %s\n", insertedID)

	cursor, err := collection.Find(ctx, bson.M{})
	if err != nil {
		log.Printf("Ошибка при получении заметок: %v\n", err)
		return
	}

	defer func() {
		if cerr := cursor.Close(ctx); cerr != nil {
			log.Printf("Ошибка при закрытии курсора: %v\n", cerr)
		}
	}()

	var notes []model.Note
	if err = cursor.All(ctx, &notes); err != nil {
		log.Printf("Ошибка декодирования заметок: %v\n", err)
		return
	}
	log.Printf("Найдено заметок: %d\n", len(notes))

	for _, n := range notes {
		updateStr := "не обновлялась"
		if n.UpdatedAt != nil {
			updateStr = n.UpdatedAt.String()
		}
		log.Printf("\n\nID: %s, Заголовок: %s, Содержание: %s, Создана: %v, Обновлена: %s\n",
			n.ID.Hex(), n.Title, n.Body, n.CreatedAt, updateStr)
	}

	now := time.Now()

	id, err := primitive.ObjectIDFromHex("696af78669142a2ab2732fc0")
	if err != nil {
		log.Fatal(err)
	}
	updateResult, err := collection.UpdateOne(ctx,
		bson.M{"_id": id},
		bson.M{
			"$set": bson.M{
				"title":      gofakeit.City(),           // Новый случайный заголовок
				"body":       gofakeit.Address().Street, // Новое случайное содержание
				"updated_at": now,
			},
		},
	)

	if err != nil {
		log.Printf("Ошибка обновления заметки: %v\n", err)
		return
	}
	log.Printf("Обновлено заметок: %d\n", updateResult.ModifiedCount)

	var updatedNote model.Note
	err = collection.FindOne(ctx, bson.M{"_id": id}).Decode(&updatedNote)
	if err != nil {
		log.Printf("Ошибка получения обновленной заметки: %v\n", err)
		return
	}

	// Проверяем, есть ли значение в поле UpdatedAt
	updatedAtStr := "не обновлялась"
	if updatedNote.UpdatedAt != nil {
		updatedAtStr = updatedNote.UpdatedAt.Format(time.RFC3339)
	}

	log.Printf("Обновленная заметка - ID: %s, Заголовок: %s, Содержание: %s, Создана: %v, Обновлена: %s\n",
		updatedNote.ID.Hex(), updatedNote.Title, updatedNote.Body,
		updatedNote.CreatedAt.Format(time.RFC3339), updatedAtStr)

	deleteResult, err := collection.DeleteOne(ctx, bson.M{"_id": id})
	if err != nil {
		log.Printf("Ошибка удаления заметки: %v\n", err)
		return
	}
	log.Printf("Удалено заметок: %d, ID удаленной заметки: %s\n",
		deleteResult.DeletedCount, id)

	log.Println("Демонстрация операций с MongoDB успешно завершена!")

}
