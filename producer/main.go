package main

import (
	"context"
	"fmt"
	"golang-assessment/golang-assessment/proto"
	"golang-assessment/shared"
	"log"
	"math/rand"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"google.golang.org/grpc"
)

var version = "1.0.0"

var taskCounter = prometheus.NewCounterVec(
	prometheus.CounterOpts{
		Name: "tasks_produced_total",
		Help: "Total number of tasks produced",
	},
	[]string{"type"},
)

func init() {
	prometheus.MustRegister(taskCounter)
}

// produceTask generates a random task type and task value
func produceTask() (int, int) {
	r := rand.New(rand.NewSource(time.Now().UnixNano())) // Create a new random number generator
	taskType := r.Intn(10)                               // Random task type between 0 and 9
	taskValue := r.Intn(100)                             // Random task value between 0 and 99
	return taskType, taskValue
}

func main() {
	// Check for -version flag
	if len(os.Args) > 1 && os.Args[1] == "-version" {
		fmt.Println("Version:", version)
		return
	}

	// Load configuration
	config, err := shared.LoadConfig()
	if err != nil {
		log.Fatalf("Error loading config: %v", err)
	}

	// Initialize logger
	logger := shared.InitLogger(config.LogLevel)
	logger.Info("Producer service started")

	// Start Prometheus metrics server
	go func() {
		http.Handle("/metrics", promhttp.Handler())
		log.Fatal(http.ListenAndServe("0.0.0.0:9091", nil)) // Producer metrics on port 9091
	}()

	// Connect to gRPC consumer
	conn, err := grpc.Dial("consumer:50051", grpc.WithInsecure())
	if err != nil {
		logger.Fatalf("Failed to connect to consumer: %v", err)
	}
	defer conn.Close()

	// Create a new TaskService client
	taskServiceClient := proto.NewTaskServiceClient(conn)

	// Produce tasks
	for i := 0; i < config.MaxBacklog; i++ {
		taskType, taskValue := produceTask()

		logger.Infof("Produced task type: %d, value: %d", taskType, taskValue)

		// Increment the counter for the produced task
		taskCounter.With(prometheus.Labels{"type": strconv.Itoa(taskType)}).Inc()

		// Create a TaskRequest
		taskRequest := &proto.TaskRequest{
			Type:  int32(taskType),
			Value: int32(taskValue),
		}

		// Send the task to the consumer via gRPC
		_, err := taskServiceClient.SendTask(context.Background(), taskRequest)
		if err != nil {
			logger.Errorf("Failed to send task: %v", err)
		} else {
			logger.Infof("Task sent successfully: type=%d, value=%d", taskType, taskValue)
		}

		time.Sleep(time.Duration(100) * time.Millisecond) // Simulate delay
	}

	// Keep the service running
	select {}
}
