package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"strconv"
	"time"

	"golang-assessment/golang-assessment/proto"

	_ "net/http/pprof" // for profiling

	_ "github.com/glebarez/sqlite"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/sirupsen/logrus"
	"golang.org/x/time/rate"
	"google.golang.org/grpc"
)

var version = "1.0.0"

// Task struct to hold task data
type Task struct {
	ID        int
	Type      int
	Value     int
	State     string
	CreatedAt time.Time
	UpdatedAt time.Time
}

// TaskServiceServer struct implementing the gRPC service
type TaskServiceServer struct {
	db *sql.DB
	proto.UnimplementedTaskServiceServer
}

var tasksProcessed = prometheus.NewCounterVec(
	prometheus.CounterOpts{
		Name: "tasks_processed_total",
		Help: "Total number of tasks processed",
	},
	[]string{"type"},
)

var taskState = prometheus.NewCounterVec(
	prometheus.CounterOpts{
		Name: "tasks_state_count",
		Help: "Number of tasks in each state",
	},
	[]string{"state"},
)

var limiter = rate.NewLimiter(1, 5) // 1 task per second with a burst of 5

func init() {
	prometheus.MustRegister(tasksProcessed)
	prometheus.MustRegister(taskState)
}

// NewTaskServiceServer initializes the TaskServiceServer with a DB connection
func NewTaskServiceServer(db *sql.DB) *TaskServiceServer {
	return &TaskServiceServer{db: db}
}

// SaveTask saves the task to the database
func (s *TaskServiceServer) SaveTask(task *Task) error {
	_, err := s.db.Exec("INSERT INTO tasks (type, value, state, created_at, updated_at) VALUES (?, ?, ?, ?, ?)",
		task.Type, task.Value, task.State, task.CreatedAt, task.UpdatedAt)
	return err
}

// SendTask implements the gRPC method to receive tasks from the producer
func (s *TaskServiceServer) SendTask(ctx context.Context, req *proto.TaskRequest) (*proto.TaskResponse, error) {
	// Apply rate limiting
	if err := limiter.Wait(ctx); err != nil {
		return nil, err
	}

	task := Task{
		Type:      int(req.Type),
		Value:     int(req.Value),
		State:     "received", // Initial state
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	// Increment the state count for received tasks
	logrus.Infof("Tracking state: %s", task.State)
	taskState.With(prometheus.Labels{"state": task.State}).Inc()

	// Simulate task processing
	time.Sleep(time.Duration(task.Value) * time.Millisecond)
	task.State = "done"

	// Increment the state count for done tasks
	logrus.Infof("Tracking state: %s", task.State)
	taskState.With(prometheus.Labels{"state": task.State}).Inc()

	err := s.SaveTask(&task)
	if err != nil {
		logrus.Error("Failed to save task: ", err)
		return nil, fmt.Errorf("failed to save task: %v", err)
	}

	// Track the processed task in Prometheus
	tasksProcessed.With(prometheus.Labels{"type": strconv.Itoa(task.Type)}).Inc()

	logrus.Infof("Task saved: %+v", task)

	return &proto.TaskResponse{
		Status: "Task saved successfully",
	}, nil
}

// runMigrations runs database migrations using golang-migrate
func runMigrations(db *sql.DB) {
	_, err := db.Exec(`CREATE TABLE IF NOT EXISTS tasks (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		type INTEGER NOT NULL,
		value INTEGER NOT NULL,
		state TEXT NOT NULL,
		created_at DATETIME NOT NULL,
		updated_at DATETIME NOT NULL
	)`)
	if err != nil {
		log.Fatalf("Failed to create tasks table: %v", err)
	}
	log.Println("Tasks table created or already exists.")
}

func main() {
	// Check for -version flag
	if len(os.Args) > 1 && os.Args[1] == "-version" {
		fmt.Println("Version:", version)
		return
	}

	// Set up logging
	logrus.SetFormatter(&logrus.JSONFormatter{})
	logrus.SetLevel(logrus.InfoLevel)

	// Start profiling
	go func() {
		log.Println("Starting pprof at localhost:6062")
		log.Println(http.ListenAndServe("localhost:6062", nil))
	}()

	// Set up database connection
	dbPath := "/app/tasks.db"
	db, err := sql.Open("sqlite", dbPath)
	if err != nil {
		log.Fatalf("Failed to connect to the database: %v", err)
	}
	defer db.Close()

	// Ensure the database file exists
	_, err = os.Stat(dbPath)
	if os.IsNotExist(err) {
		file, err := os.Create(dbPath)
		if err != nil {
			log.Fatalf("Failed to create the database file: %v", err)
		}
		file.Close()
		log.Printf("Database file created: %s\n", dbPath)
	} else if err != nil {
		log.Fatalf("Error checking the database file: %v", err)
	}

	// Run database migrations
	runMigrations(db)

	// Set up Prometheus metrics server
	go func() {
		http.Handle("/metrics", promhttp.Handler())
		log.Fatal(http.ListenAndServe("0.0.0.0:9092", nil)) // Exposing Prometheus metrics at port 9092
	}()

	// Set up gRPC server
	listener, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatal("Failed to listen on port 50051: ", err)
	}

	grpcServer := grpc.NewServer()
	taskServiceServer := NewTaskServiceServer(db)

	// Register the TaskServiceServer with the gRPC server
	proto.RegisterTaskServiceServer(grpcServer, taskServiceServer)

	logrus.Info("gRPC server is running on port 50051")

	// Start the gRPC server
	if err := grpcServer.Serve(listener); err != nil {
		log.Fatalf("Failed to serve gRPC server: %v", err)
	}

	// Keep the service running
	select {}
}
