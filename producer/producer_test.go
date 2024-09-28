package main

import "testing"

func TestTaskProduction(t *testing.T) {
	taskType, taskValue := produceTask()
	if taskType < 0 || taskType > 9 {
		t.Errorf("Invalid task type: %d", taskType)
	}
	if taskValue < 0 || taskValue > 99 {
		t.Errorf("Invalid task value: %d", taskValue)
	}
}
