package main

import (
	"context"
	"database/sql"
	"testing"
	"time"

	"golang-assessment/golang-assessment/proto"

	_ "github.com/glebarez/sqlite"
)

type mockTaskServiceClient struct {
	proto.UnimplementedTaskServiceServer
}

func (m *mockTaskServiceClient) SendTask(ctx context.Context, req *proto.TaskRequest) (*proto.TaskResponse, error) {
	return &proto.TaskResponse{
		Status: "Task saved successfully",
	}, nil
}

func TestSendTask(t *testing.T) {
	mockClient := &mockTaskServiceClient{}

	req := &proto.TaskRequest{
		Type:  2,
		Value: 50,
	}

	resp, err := mockClient.SendTask(context.Background(), req)
	if err != nil {
		t.Errorf("Error in SendTask: %v", err)
	}

	if resp.Status != "Task saved successfully" {
		t.Errorf("Expected 'Task saved successfully', got '%s'", resp.Status)
	}
}

func TestSaveTask(t *testing.T) {
	db, err := sql.Open("sqlite", ":memory:")
	if err != nil {
		t.Fatalf("Failed to open test database: %v", err)
	}
	defer db.Close()

	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS tasks (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		type INTEGER NOT NULL,
		value INTEGER NOT NULL,
		state TEXT NOT NULL,
		created_at DATETIME NOT NULL,
		updated_at DATETIME NOT NULL
	)`)
	if err != nil {
		t.Fatalf("Failed to create tasks table: %v", err)
	}

	s := NewTaskServiceServer(db)

	task := &Task{
		Type:      2,
		Value:     50,
		State:     "received",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	err = s.SaveTask(task)
	if err != nil {
		t.Errorf("Error saving task: %v", err)
	}

	var taskCount int
	err = db.QueryRow("SELECT COUNT(*) FROM tasks WHERE type = ? AND value = ?", 2, 50).Scan(&taskCount)
	if err != nil {
		t.Errorf("Error querying task: %v", err)
	}

	if taskCount != 1 {
		t.Errorf("Expected 1 task, found %d", taskCount)
	}
}
