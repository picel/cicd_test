// main.go
package main

import (
	"context"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/google/uuid" // UUID 패키지 사용
	"github.com/segmentio/kafka-go"
)

var (
	kafkaBroker  = "kafka.event-queue.svc.cluster.local:9092"
	kafkaWriters = map[string]*kafka.Writer{
		"service-a": kafka.NewWriter(kafka.WriterConfig{
			Brokers: []string{kafkaBroker},
			Topic:   "service-a",
		}),
		"service-b": kafka.NewWriter(kafka.WriterConfig{
			Brokers: []string{kafkaBroker},
			Topic:   "service-b",
		}),
	}
)

func publishKafkaEventAsync(topic, message, correlationID string) {
	go func() {
		if err := publishKafkaEvent(topic, message, correlationID); err != nil {
			log.Printf("Failed to send Kafka event asynchronously: %v", err)
		}
	}()
}

func main() {
	// 서버 핸들러 등록
	http.HandleFunc("/api/v1/bff/test1", handleTest1)
	http.HandleFunc("/api/v1/bff/test2", handleTest2)

	// 프로그램 종료 시 Kafka Writer 닫기
	defer closeKafkaWriters()

	// 서버 시작
	log.Println("BFF Server started at :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

// Kafka Writer 닫기 함수
func closeKafkaWriters() {
	for _, writer := range kafkaWriters {
		writer.Close()
	}
}

func handleTest1(w http.ResponseWriter, r *http.Request) {
	correlationID := uuid.New().String() // Correlation ID 생성

	// Kafka Event 발행
	publishKafkaEventAsync("service-a", "Test1 Kafka Event", correlationID)

	// HTTP 요청 보내기 (context with timeout)
	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, "http://service-b.b-server.svc.cluster.local:8081/api/v1/service-b/test1", nil)
	if err != nil {
		log.Printf("Failed to create HTTP request: %v", err)
		http.Error(w, "Failed to create HTTP request", http.StatusInternalServerError)
		return
	}

	req.Header.Set("X-Correlation-ID", correlationID) // Correlation ID를 헤더에 추가

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Printf("Failed to send HTTP request: %v", err)
		http.Error(w, "Failed to send HTTP request", http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		http.Error(w, "Service B returned an error", resp.StatusCode)
		return
	}

	body, _ := io.ReadAll(resp.Body)
	w.Write(body)
}

func handleTest2(w http.ResponseWriter, r *http.Request) {
	correlationID := uuid.New().String() // Correlation ID 생성

	// Kafka Event 발행
	publishKafkaEventAsync("service-b", "Test2 Kafka Event", correlationID)

	// HTTP 요청 보내기 (context with timeout)
	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, "http://service-a.a-server.svc.cluster.local:8081/api/v1/service-a/test1", nil)
	if err != nil {
		log.Printf("Failed to create HTTP request: %v", err)
		http.Error(w, "Failed to create HTTP request", http.StatusInternalServerError)
		return
	}

	req.Header.Set("X-Correlation-ID", correlationID) // Correlation ID를 헤더에 추가

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Printf("Failed to send HTTP request: %v", err)
		http.Error(w, "Failed to send HTTP request", http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		http.Error(w, "Service A returned an error", resp.StatusCode)
		return
	}

	body, _ := io.ReadAll(resp.Body)
	w.Write(body)
}

func publishKafkaEvent(topic, message, correlationID string) error {
	// Writer 재사용
	writer, exists := kafkaWriters[topic]
	if !exists {
		log.Printf("Kafka writer for topic %s does not exist", topic)
		return nil
	}

	// 메시지에 Correlation ID 포함
	msg := kafka.Message{
		Key:   []byte(correlationID),
		Value: []byte(message),
	}

	// 메시지 전송
	err := writer.WriteMessages(context.Background(), msg)
	if err != nil {
		log.Printf("Error writing message to Kafka: %v", err)
	}
	return err
}
