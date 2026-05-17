package tracing

import (
	"context"
	"fmt"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/trace"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

// Константы для заголовков трассировки
const (
	// TraceIDHeader - заголовок для передачи trace ID
	TraceIDHeader = "x-trace-id"
)

// UnaryServerInterceptor создает gRPC unary interceptor для трассировки входящих запросов.
// Interceptor извлекает контекст трассировки из входящего запроса и создает новый спан для каждого запроса.
//
// Трейсер (Tracer) - это компонент, который создает и управляет спанами. Он отвечает за:
// - Создание новых спанов (начальных или дочерних)
// - Установку ID трейса и связей между спанами
// - Сбор информации о выполнении операций
//
// Пропагатор (Propagator) - это компонент, который отвечает за передачу контекста трассировки
// между сервисами. Он извлекает и внедряет данные трассировки в заголовки запросов,
// что позволяет поддерживать непрерывную трассировку через границы сервисов.
func UnaryServerInterceptor(serviceName string) grpc.UnaryServerInterceptor {
	// Получаем текущий трейсер и пропагатор
	tracer := otel.GetTracerProvider().Tracer(serviceName)
	propagator := otel.GetTextMapPropagator()

	// Возвращаем функцию-interceptor
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		fmt.Printf("\n=== [interceptor] метод: %s ===\n", info.FullMethod)

		// ШАГ 1: входящие метаданные от клиента
		md, ok := metadata.FromIncomingContext(ctx)
		if !ok {
			md = metadata.New(nil)
		}
		fmt.Printf("[шаг 1] входящие метаданные: %v\n", map[string][]string(md))

		// ШАГ 2: извлекаем trace-контекст из метаданных
		// если клиент передал traceparent — пропагатор положит его в ctx
		ctx = propagator.Extract(ctx, metadataCarrier(md))
		scBefore := trace.SpanContextFromContext(ctx)
		fmt.Printf("[шаг 2] после Extract — traceID: %s, valid: %v\n",
			scBefore.TraceID(), scBefore.IsValid())

		// ШАГ 3: создаём span для этого запроса
		ctx, span := tracer.Start(
			ctx,
			info.FullMethod,
			trace.WithSpanKind(trace.SpanKindServer),
		)
		defer span.End()
		sc := span.SpanContext()
		fmt.Printf("[шаг 3] span создан — traceID: %s, spanID: %s\n", sc.TraceID(), sc.SpanID())

		// ШАГ 4: пишем trace ID в ответные метаданные
		ctx = AddTraceIDToResponse(ctx)
		fmt.Printf("[шаг 4] x-trace-id добавлен в ответ: %s\n", sc.TraceID())

		// ШАГ 5: вызываем реальный хендлер
		resp, err := handler(ctx, req)
		if err != nil {
			span.RecordError(err)
			fmt.Printf("[шаг 5] ошибка (записана в span): %v\n", err)
		} else {
			fmt.Printf("[шаг 5] хендлер выполнен успешно\n")
		}

		return resp, err
	}
}

// UnaryClientInterceptor создает gRPC unary interceptor для трассировки исходящих запросов.
// Interceptor добавляет контекст трассировки в исходящий запрос.
//
// Outgoing context - это контекст, который отправляется в запросе к другому сервису.
// В него добавляются данные о текущем трейсе, чтобы следующий сервис мог
// продолжить цепочку трассировки.
func UnaryClientInterceptor(serviceName string) grpc.UnaryClientInterceptor {
	// Получаем текущий трейсер и пропагатор
	tracer := otel.GetTracerProvider().Tracer(serviceName)
	propagator := otel.GetTextMapPropagator()

	// Возвращаем функцию-interceptor
	return func(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
		fmt.Printf("\n=== [client interceptor] метод: %s ===\n", method)

		// Определяем имя спана в зависимости от наличия контекста трассировки
		spanName := formatSpanName(ctx, method)
		fmt.Printf("[шаг 1] имя span'а: %s\n", spanName)

		// Создаем новый спан с подготовленным именем
		ctx, span := tracer.Start(
			ctx,
			spanName,
			trace.WithSpanKind(trace.SpanKindClient),
		)
		defer span.End()
		sc := span.SpanContext()
		fmt.Printf("[шаг 2] span создан — traceID: %s, spanID: %s\n", sc.TraceID(), sc.SpanID())

		// Создаем переносчик метаданных для пропагации трассировки
		carrier := metadataCarrier(extractOutgoingMetadata(ctx))
		fmt.Printf("[шаг 3] carrier до Inject: %v\n", map[string][]string(metadata.MD(carrier)))

		// Внедряем контекст трассировки в метаданные
		// Пропагатор добавляет информацию о текущем трейсе в метаданные запроса
		propagator.Inject(ctx, carrier)
		fmt.Printf("[шаг 4] carrier после Inject: %v\n", map[string][]string(metadata.MD(carrier)))

		// Обновляем метаданные в контексте
		ctx = metadata.NewOutgoingContext(ctx, metadata.MD(carrier))
		fmt.Printf("[шаг 5] метаданные записаны в контекст, отправляем запрос\n")

		// Вызываем следующий обработчик с обогащенным контекстом
		err := invoker(ctx, method, req, reply, cc, opts...)
		// Если произошла ошибка, записываем её в спан
		if err != nil {
			trace.SpanFromContext(ctx).RecordError(err)
			fmt.Printf("[шаг 6] ошибка (записана в span): %v\n", err)
		} else {
			fmt.Printf("[шаг 6] запрос выполнен успешно\n")
		}

		return err
	}
}

// formatSpanName формирует имя спана в зависимости от наличия контекста трассировки
func formatSpanName(ctx context.Context, method string) string {
	if !trace.SpanContextFromContext(ctx).IsValid() {
		return "client." + method
	}

	return method
}

// extractOutgoingMetadata извлекает исходящие метаданные из контекста и создает их копию
func extractOutgoingMetadata(ctx context.Context) metadata.MD {
	md, ok := metadata.FromOutgoingContext(ctx)
	if !ok {
		return metadata.New(nil)
	}

	return md.Copy()
}

// GetTraceIDFromContext извлекает trace ID из контекста.
// Полезно для логирования и возврата trace ID клиенту.
func GetTraceIDFromContext(ctx context.Context) string {
	span := trace.SpanFromContext(ctx)
	if !span.SpanContext().IsValid() {
		return ""
	}

	return span.SpanContext().TraceID().String()
}

// AddTraceIDToResponse добавляет trace ID в исходящие метаданные gRPC ответа.
// Это позволяет клиенту получить trace ID для последующего поиска в системе трассировки.
func AddTraceIDToResponse(ctx context.Context) context.Context {
	traceID := GetTraceIDFromContext(ctx)
	if traceID == "" {
		return ctx
	}
	return metadata.AppendToOutgoingContext(ctx, TraceIDHeader, traceID)

	// Используем вспомогательную функцию для извлечения метаданных
	// md := extractOutgoingMetadata(ctx)

	// md.Set(TraceIDHeader, traceID)
	// return metadata.NewOutgoingContext(ctx, md)
}
