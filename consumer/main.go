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

	_ "net/http/pprof"

	_ "github.com/glebarez/sqlite"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/sirupsen/logrus"
	"golang.org/x/time/rate"
	"google.golang.org/grpc"
)

var version = "1.0.0"

type Task struct {
	ID        int
	Type      int
	Value     int
	State     string
	CreatedAt time.Time
	UpdatedAt time.Time
}

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

var limiter = rate.NewLimiter(1, 5)

func init() {
	prometheus.MustRegister(tasksProcessed)
	prometheus.MustRegister(taskState)
}

func NewTaskServiceServer(db *sql.DB) *TaskServiceServer {
	return &TaskServiceServer{db: db}
}

func (s *TaskServiceServer) SaveTask(task *Task) error {
	_, err := s.db.Exec("INSERT INTO tasks (type, value, state, created_at, updated_at) VALUES (?, ?, ?, ?, ?)",
		task.Type, task.Value, task.State, task.CreatedAt, task.UpdatedAt)
	return err
}

func (s *TaskServiceServer) SendTask(ctx context.Context, req *proto.TaskRequest) (*proto.TaskResponse, error) {

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

	taskState.With(prometheus.Labels{"state": task.State}).Inc()

	time.Sleep(time.Duration(task.Value) * time.Millisecond)
	task.State = "done"

	taskState.With(prometheus.Labels{"state": task.State}).Inc()

	err := s.SaveTask(&task)
	if err != nil {
		logrus.Error("Failed to save task: ", err)
		return nil, fmt.Errorf("failed to save task: %v", err)
	}

	tasksProcessed.With(prometheus.Labels{"type": strconv.Itoa(task.Type)}).Inc()

	logrus.Infof("Task saved: %+v", task)

	return &proto.TaskResponse{
		Status: "Task saved successfully",
	}, nil
}

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
	if len(os.Args) > 1 && os.Args[1] == "-version" {
		fmt.Println("Version:", version)
		return
	}

	logrus.SetFormatter(&logrus.JSONFormatter{})
	logrus.SetLevel(logrus.InfoLevel)

	go func() {
		log.Println(http.ListenAndServe("localhost:6062", nil))
	}()

	dbPath := "/app/tasks.db"
	db, err := sql.Open("sqlite", dbPath)
	if err != nil {
		log.Fatalf("Failed to connect to the database: %v", err)
	}
	defer db.Close()

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

	runMigrations(db)

	go func() {
		http.Handle("/metrics", promhttp.Handler())
		log.Fatal(http.ListenAndServe("0.0.0.0:9092", nil))
	}()

	listener, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatal("Failed to listen on port 50051: ", err)
	}

	grpcServer := grpc.NewServer()
	taskServiceServer := NewTaskServiceServer(db)

	proto.RegisterTaskServiceServer(grpcServer, taskServiceServer)

	logrus.Info("gRPC server is running on port 50051")

	if err := grpcServer.Serve(listener); err != nil {
		log.Fatalf("Failed to serve gRPC server: %v", err)
	}

	select {}
}
