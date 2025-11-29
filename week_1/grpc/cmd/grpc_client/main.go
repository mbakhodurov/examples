package main

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/brianvoe/gofakeit/v7"
	ufo_v1 "github.com/mbakhodurov/examples/week_1/grpc/pkg/proto/ufo/v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/protobuf/types/known/timestamppb"
	"google.golang.org/protobuf/types/known/wrapperspb"
)

const serverAddress = "localhost:50051"

func UpdateSight(ctx context.Context, client ufo_v1.UFOServiceClient, uuid string) error {
	if uuid == "" {
		return fmt.Errorf("UpdateSight: uuid не может быть пустым")
	}
	updateInfo := &ufo_v1.SightingUpdateInfo{
		Description: wrapperspb.String("Обновленное описание наблюдения"),
		Color:       wrapperspb.String("Red"),
	}
	_, err := client.Update(ctx, &ufo_v1.UpdateRequest{
		Uuid:       uuid,
		UpdateInfo: updateInfo,
	})
	if err != nil {
		return fmt.Errorf("UpdateSight: ошибка при обновлении наблюдения с UUID %s: %w", uuid, err)
	}
	log.Printf("Наблюдение с UUID %s успешно обновлено", uuid)
	return nil
}

func CreateSight(ctx context.Context, client ufo_v1.UFOServiceClient) (string, error) {
	observedAt := gofakeit.DateRange(
		time.Now().AddDate(-3, 0, 0), // за последние 3 года
		time.Now(),
	)
	location := gofakeit.City() + ", " + gofakeit.StreetName()
	description := gofakeit.Sentence(gofakeit.Number(5, 15))

	info := &ufo_v1.SightingInfo{
		ObservedAt:  timestamppb.New(observedAt),
		Location:    location,
		Description: description,
	}
	if gofakeit.Bool() {
		info.Color = wrapperspb.String(gofakeit.Color())
	}

	if gofakeit.Bool() {
		info.Sound = wrapperspb.Bool(gofakeit.Bool())
	}

	if gofakeit.Bool() {
		info.DurationSeconds = wrapperspb.Int32(gofakeit.Int32())
	}

	resp, err := client.Create(ctx, &ufo_v1.CreateRequest{Info: info})
	if err != nil {
		return "", err
	}
	return resp.Uuid, nil
}

func DeleteSight(ctx context.Context, client ufo_v1.UFOServiceClient, uuid string) error {
	if uuid == "" {
		return fmt.Errorf("DeleteSight: uuid не может быть пустым")
	}

	_, err := client.Delete(context.Background(), &ufo_v1.DeleteRequest{Uuid: uuid})
	if err != nil {
		return fmt.Errorf("DeleteSight: ошибка при удалении наблюдения с UUID %s: %w", uuid, err)
	}
	return nil
}

func GetSgithing(ctx context.Context, client ufo_v1.UFOServiceClient, uuid string) (*ufo_v1.Sighting, error) {
	resp, err := client.Get(ctx, &ufo_v1.GetRequest{Uuid: uuid})
	if err != nil {
		return nil, err
	}
	return resp.Sighting, nil
}

func GetAllSights(ctx context.Context, client ufo_v1.UFOServiceClient) ([]*ufo_v1.Sighting, error) {
	resp, err := client.GetAll(ctx, &ufo_v1.GetAllRequest{})
	if err != nil {
		return nil, err
	}
	if resp == nil {
		return nil, fmt.Errorf("GetAllSights: empty response")
	}
	return resp.Sightings, nil
}

func main() {
	conn, err := grpc.NewClient(serverAddress,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		log.Printf("Failed to connect:%v,\n", err)
	}

	defer func() {
		if err := conn.Close(); err != nil {
			log.Printf("Failed to close connection: %v\n", err)
		}
	}()

	client := ufo_v1.NewUFOServiceClient(conn)
	log.Println("=== Тестирование API для работы с наблюдениями НЛО ===")
	log.Println()

	scanner := bufio.NewScanner(os.Stdin)

	for {
		fmt.Println("\nВыберите действие:")
		fmt.Println("1. Создать наблюдение")
		fmt.Println("2. Получить наблюдение по UUID")
		fmt.Println("3. Получить все наблюдения")
		fmt.Println("4. Удаление наблюдение")
		fmt.Println("5. Обновить наблюдение")
		fmt.Println("0. Выход")
		fmt.Print("Введите номер действия: ")

		if !scanner.Scan() {
			break
		}
		choice := scanner.Text()

		switch choice {
		case "1":
			uuid, err := CreateSight(context.Background(), client)
			if err != nil {
				log.Printf("Ошибка при создании: %v\n", err)
				continue
			}
			log.Printf("Создано новое наблюдение с UUID: %s\n", uuid)

		case "2":
			fmt.Print("Введите UUID: ")
			if !scanner.Scan() {
				break
			}
			uuid := scanner.Text()
			sighting, err := GetSgithing(context.Background(), client, uuid)
			if err != nil {
				log.Printf("Ошибка при получении: %v\n", err)
				continue
			}
			log.Printf("Наблюдение: %+v\n", sighting)

		case "3":
			sightings, err := GetAllSights(context.Background(), client)
			if err != nil {
				log.Printf("Ошибка при получении всех наблюдений: %v\n", err)
				continue
			}
			log.Printf("Всего наблюдений: %d\n", len(sightings))
			for i, s := range sightings {
				log.Printf("%d: %+v\n\n", i+1, s)
			}

		case "4":
			fmt.Print("Введите UUID: ")
			if !scanner.Scan() {
				break
			}
			uuid := scanner.Text()
			err := DeleteSight(context.Background(), client, uuid)
			if err != nil {
				log.Printf("Ошибка при удалении: %v\n", err)
				continue
			}
			log.Printf("Наблюдение успешно удалена: %+v\n", uuid)

		case "5":
			fmt.Print("Введите UUID для обновления: ")
			if !scanner.Scan() {
				break
			}
			uuid := scanner.Text()
			err := UpdateSight(context.Background(), client, uuid)
			if err != nil {
				log.Printf("Ошибка при обновлении: %v\n", err)
				continue
			}
			log.Printf("Наблюдение успешно обновлено: %s\n", uuid)
		case "0":
			log.Println("Выход из программы")
			return

		default:
			log.Println("Неверный выбор, попробуйте снова")
		}
	}
}
