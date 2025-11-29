package main

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/brianvoe/gofakeit"
	ufoV1 "github.com/mbakhodurov/examples/week_1/grpc_with_interceptor/pkg/proto/ufo/v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/protobuf/types/known/timestamppb"
	"google.golang.org/protobuf/types/known/wrapperspb"
)

const serverAddress = "localhost:50051"

func CreateSighting(ctx context.Context, client ufoV1.UFOServiceClient) (string, error) {
	observedAt := gofakeit.DateRange(
		time.Now().AddDate(-3, 0, 0), // за последние 3 года
		time.Now(),
	)
	location := gofakeit.City() + ", " + gofakeit.StreetName()
	description := gofakeit.Sentence(gofakeit.Number(5, 15))

	info := &ufoV1.SightingInfo{
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

	resp, err := client.Create(ctx, &ufoV1.CreateRequest{Info: info})
	if err != nil {
		return "", err
	}
	return resp.Uuid, nil
}

func GetSight(ctx context.Context, client ufoV1.UFOServiceClient, uuid string) (*ufoV1.Sighting, error) {
	resp, err := client.Get(ctx, &ufoV1.GetRequest{Uuid: uuid})
	if err != nil {
		return nil, err
	}
	return resp.Sighting, nil
}

func GetAllSights(ctx context.Context, client ufoV1.UFOServiceClient) ([]*ufoV1.Sighting, error) {
	resp, err := client.GetAll(ctx, &ufoV1.GetAllRequest{})
	if err != nil {
		return nil, fmt.Errorf("GetAllSights: %w", err)
	}
	if resp == nil {
		return nil, fmt.Errorf("GetAllSights: empty response")
	}

	return resp.GetSightings(), nil
}

func DeleteSight(ctx context.Context, client ufoV1.UFOServiceClient, uuid string) error {
	if uuid == "" {
		return fmt.Errorf("DeleteSight: uuid не может быть пустым")
	}

	_, err := client.Delete(ctx, &ufoV1.DeleteRequest{Uuid: uuid})
	if err != nil {
		return fmt.Errorf("DeleteSight: ошибка при удалении наблюдения с UUID %s: %w", uuid, err)
	}
	return nil
}

func UpdateSight(ctx context.Context, client ufoV1.UFOServiceClient, uuid string) error {
	if uuid == "" {
		return fmt.Errorf("UpdateSight: uuid не может быть пустым")
	}
	updateInfo := &ufoV1.SightingUpdateInfo{
		Description: wrapperspb.String("Обновленное описание наблюдения"),
		Color:       wrapperspb.String("Red"),
	}
	_, err := client.Update(ctx, &ufoV1.UpdateRequest{
		Uuid:       uuid,
		UpdateInfo: updateInfo,
	})
	if err != nil {
		return fmt.Errorf("UpdateSight: ошибка при обновлении наблюдения с UUID %s: %w", uuid, err)
	}
	log.Printf("Наблюдение с UUID %s успешно обновлено", uuid)
	return nil
}

func main() {
	conn, err := grpc.NewClient(serverAddress,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		log.Printf("failed to connect: %v\n", err)
		return
	}
	defer func() {
		if cerr := conn.Close(); cerr != nil {
			log.Printf("failed to close connect: %v", cerr)
		}
	}()

	client := ufoV1.NewUFOServiceClient(conn)
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
			uuid, err := CreateSighting(context.Background(), client)
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
			sighting, err := GetSight(context.Background(), client, uuid)
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
