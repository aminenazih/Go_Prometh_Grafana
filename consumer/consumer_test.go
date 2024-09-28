package main

import (
	"testing"
	"time"
)

func TestSaveTask(t *testing.T) {
	// Simulating saving a task
	task := Task{
		Type:      1,
		Value:     100,
		State:     "done",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	if task.Type < 0 || task.Value < 0 {
		t.Errorf("Invalid task: %+v", task)
	}
}
