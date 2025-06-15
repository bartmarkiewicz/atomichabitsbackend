package integration

import (
	"bytes"
	"encoding/json"
	"fmt"
	"habitgobackend/cmd/api/resource/habit"
	"habitgobackend/test/util"
	"net/http"
	"strings"
	"sync"
	"testing"

	"github.com/google/uuid"
)

const (
	baseURL = "http://localhost:8080/v1"
)

func TestApiSmoke_GetHabits(testing *testing.T) {
	testing.Parallel()

	resp, err := http.Get(fmt.Sprintf("%s/habits", baseURL))

	if err != nil {
		testing.Fatalf("Failed to get habits: %s", err)
	}

	util.IsEqual(testing, resp.StatusCode, http.StatusOK)

	habits := &habit.Habits{}

	if err = json.NewDecoder(resp.Body).Decode(habits); err != nil {
		testing.Fatalf("Failed to decode habits: %s", err)
	}

	expected := habit.Habits{}

	util.HabitsEqual(testing, *habits, expected)
}

func TestApiSmoke_CreateHabit(testing *testing.T) {
	testing.Parallel()

	// Create a new habit
	newHabit := habit.JsonHabit{
		Description: "Test habit",
		ColourHex:   "#FF5733",
		IconBase64:  "data:image/png;base64,iVBORw0KGgoAAAANSUhEUgAAAAEAAAABCAQAAAC1HAwCAAAAC0lEQVR42mNk+A8AAQUBAScY42YAAAAASUVORK5CYII=",
		ModeType:    "daily",
	}

	habitJSON, err := json.Marshal(newHabit)
	if err != nil {
		testing.Fatalf("Failed to marshal habit: %s", err)
	}

	// Send POST request to create habit
	resp, err := http.Post(
		fmt.Sprintf("%s/habits", baseURL),
		"application/json",
		bytes.NewBuffer(habitJSON),
	)
	if err != nil {
		testing.Fatalf("Failed to create habit: %s", err)
	}
	defer resp.Body.Close()

	// Check response status code
	util.IsEqual(testing, resp.StatusCode, http.StatusCreated)

	// Decode response
	var createdHabit habit.JsonHabit
	if err = json.NewDecoder(resp.Body).Decode(&createdHabit); err != nil {
		testing.Fatalf("Failed to decode created habit: %s", err)
	}

	// Verify habit was created with correct data
	util.IsEqual(testing, createdHabit.Description, newHabit.Description)
	util.IsEqual(testing, createdHabit.ColourHex, newHabit.ColourHex)
	util.IsEqual(testing, createdHabit.IconBase64, newHabit.IconBase64)
	util.IsEqual(testing, createdHabit.ModeType, newHabit.ModeType)

	// Verify ID was generated
	_, err = uuid.Parse(createdHabit.ID)
	if err != nil {
		testing.Fatalf("Invalid UUID returned: %s", err)
	}
}

func TestApiSmoke_GetHabitById(testing *testing.T) {
	testing.Parallel()

	// First create a habit
	newHabit := habit.JsonHabit{
		Description: "Test habit for get by ID",
		ColourHex:   "#3366FF",
		IconBase64:  "data:image/png;base64,iVBORw0KGgoAAAANSUhEUgAAAAEAAAABCAQAAAC1HAwCAAAAC0lEQVR42mNk+A8AAQUBAScY42YAAAAASUVORK5CYII=",
		ModeType:    "weekly",
	}

	habitJSON, err := json.Marshal(newHabit)
	if err != nil {
		testing.Fatalf("Failed to marshal habit: %s", err)
	}

	// Send POST request to create habit
	resp, err := http.Post(
		fmt.Sprintf("%s/habits", baseURL),
		"application/json",
		bytes.NewBuffer(habitJSON),
	)
	if err != nil {
		testing.Fatalf("Failed to create habit: %s", err)
	}

	// Decode response to get the ID
	var createdHabit habit.JsonHabit
	if err = json.NewDecoder(resp.Body).Decode(&createdHabit); err != nil {
		testing.Fatalf("Failed to decode created habit: %s", err)
	}
	resp.Body.Close()

	// Now get the habit by ID
	getResp, err := http.Get(fmt.Sprintf("%s/habits/%s", baseURL, createdHabit.ID))
	if err != nil {
		testing.Fatalf("Failed to get habit by ID: %s", err)
	}
	defer getResp.Body.Close()

	// Check response status code
	util.IsEqual(testing, getResp.StatusCode, http.StatusOK)

	// Decode response
	var retrievedHabit habit.JsonHabit
	if err = json.NewDecoder(getResp.Body).Decode(&retrievedHabit); err != nil {
		testing.Fatalf("Failed to decode retrieved habit: %s", err)
	}

	// Verify habit data matches
	util.IsEqual(testing, retrievedHabit.ID, createdHabit.ID)
	util.IsEqual(testing, retrievedHabit.Description, newHabit.Description)
	util.IsEqual(testing, retrievedHabit.ColourHex, newHabit.ColourHex)
	util.IsEqual(testing, retrievedHabit.IconBase64, newHabit.IconBase64)
	util.IsEqual(testing, retrievedHabit.ModeType, newHabit.ModeType)
}

func TestApiSmoke_UpdateHabit(testing *testing.T) {
	testing.Parallel()

	// First create a habit
	newHabit := habit.JsonHabit{
		Description: "Test habit for update",
		ColourHex:   "#33CC33",
		IconBase64:  "data:image/png;base64,iVBORw0KGgoAAAANSUhEUgAAAAEAAAABCAQAAAC1HAwCAAAAC0lEQVR42mNk+A8AAQUBAScY42YAAAAASUVORK5CYII=",
		ModeType:    "monthly",
	}

	habitJSON, err := json.Marshal(newHabit)
	if err != nil {
		testing.Fatalf("Failed to marshal habit: %s", err)
	}

	// Send POST request to create habit
	resp, err := http.Post(
		fmt.Sprintf("%s/habits", baseURL),
		"application/json",
		bytes.NewBuffer(habitJSON),
	)
	if err != nil {
		testing.Fatalf("Failed to create habit: %s", err)
	}

	// Decode response to get the ID
	var createdHabit habit.JsonHabit
	if err = json.NewDecoder(resp.Body).Decode(&createdHabit); err != nil {
		testing.Fatalf("Failed to decode created habit: %s", err)
	}
	resp.Body.Close()

	// Update the habit
	updatedHabit := habit.JsonHabit{
		ID:          createdHabit.ID,
		Description: "Updated test habit",
		ColourHex:   "#9933FF",
		IconBase64:  createdHabit.IconBase64,
		ModeType:    "yearly",
	}

	updatedJSON, err := json.Marshal(updatedHabit)
	if err != nil {
		testing.Fatalf("Failed to marshal updated habit: %s", err)
	}

	// Create PUT request
	req, err := http.NewRequest(
		http.MethodPut,
		fmt.Sprintf("%s/habits/%s", baseURL, createdHabit.ID),
		bytes.NewBuffer(updatedJSON),
	)
	if err != nil {
		testing.Fatalf("Failed to create PUT request: %s", err)
	}
	req.Header.Set("Content-Type", "application/json")

	// Send PUT request
	client := &http.Client{}
	putResp, err := client.Do(req)
	if err != nil {
		testing.Fatalf("Failed to update habit: %s", err)
	}
	defer putResp.Body.Close()

	// Check response status code
	util.IsEqual(testing, putResp.StatusCode, http.StatusOK)

	// Decode response
	var returnedHabit habit.JsonHabit
	if err = json.NewDecoder(putResp.Body).Decode(&returnedHabit); err != nil {
		testing.Fatalf("Failed to decode updated habit: %s", err)
	}

	// Verify habit was updated with correct data
	util.IsEqual(testing, returnedHabit.ID, updatedHabit.ID)
	util.IsEqual(testing, returnedHabit.Description, updatedHabit.Description)
	util.IsEqual(testing, returnedHabit.ColourHex, updatedHabit.ColourHex)
	util.IsEqual(testing, returnedHabit.IconBase64, updatedHabit.IconBase64)
	util.IsEqual(testing, returnedHabit.ModeType, updatedHabit.ModeType)

	// Verify by getting the habit again
	getResp, err := http.Get(fmt.Sprintf("%s/habits/%s", baseURL, createdHabit.ID))
	if err != nil {
		testing.Fatalf("Failed to get updated habit: %s", err)
	}
	defer getResp.Body.Close()

	var retrievedHabit habit.JsonHabit
	if err = json.NewDecoder(getResp.Body).Decode(&retrievedHabit); err != nil {
		testing.Fatalf("Failed to decode retrieved habit: %s", err)
	}

	// Verify habit data matches the update
	util.IsEqual(testing, retrievedHabit.Description, updatedHabit.Description)
	util.IsEqual(testing, retrievedHabit.ColourHex, updatedHabit.ColourHex)
	util.IsEqual(testing, retrievedHabit.ModeType, updatedHabit.ModeType)
}

func TestApiSmoke_DeleteHabit(testing *testing.T) {
	testing.Parallel()

	// First create a habit
	newHabit := habit.JsonHabit{
		Description: "Test habit for deletion",
		ColourHex:   "#FF33CC",
		IconBase64:  "data:image/png;base64,iVBORw0KGgoAAAANSUhEUgAAAAEAAAABCAQAAAC1HAwCAAAAC0lEQVR42mNk+A8AAQUBAScY42YAAAAASUVORK5CYII=",
		ModeType:    "daily",
	}

	habitJSON, err := json.Marshal(newHabit)
	if err != nil {
		testing.Fatalf("Failed to marshal habit: %s", err)
	}

	// Send POST request to create habit
	resp, err := http.Post(
		fmt.Sprintf("%s/habits", baseURL),
		"application/json",
		bytes.NewBuffer(habitJSON),
	)
	if err != nil {
		testing.Fatalf("Failed to create habit: %s", err)
	}

	// Decode response to get the ID
	var createdHabit habit.JsonHabit
	if err = json.NewDecoder(resp.Body).Decode(&createdHabit); err != nil {
		testing.Fatalf("Failed to decode created habit: %s", err)
	}
	resp.Body.Close()

	// Create DELETE request
	req, err := http.NewRequest(
		http.MethodDelete,
		fmt.Sprintf("%s/habits/%s", baseURL, createdHabit.ID),
		nil,
	)
	if err != nil {
		testing.Fatalf("Failed to create DELETE request: %s", err)
	}

	// Send DELETE request
	client := &http.Client{}
	deleteResp, err := client.Do(req)
	if err != nil {
		testing.Fatalf("Failed to delete habit: %s", err)
	}
	defer deleteResp.Body.Close()

	// Check response status code
	util.IsEqual(testing, deleteResp.StatusCode, http.StatusOK)

	// Verify habit was deleted by trying to get it
	getResp, err := http.Get(fmt.Sprintf("%s/habits/%s", baseURL, createdHabit.ID))
	if err != nil {
		testing.Fatalf("Failed to get deleted habit: %s", err)
	}
	defer getResp.Body.Close()

	// Should get a 404 Not Found
	util.IsEqual(testing, getResp.StatusCode, http.StatusNotFound)
}

func TestApiSmoke_ValidationErrors(testing *testing.T) {
	testing.Parallel()

	// Test case 1: Missing required fields
	invalidHabit := habit.JsonHabit{
		Description: "", // Missing required field
		ColourHex:   "#FF5733",
		IconBase64:  "data:image/png;base64,iVBORw0KGgoAAAANSUhEUgAAAAEAAAABCAQAAAC1HAwCAAAAC0lEQVR42mNk+A8AAQUBAScY42YAAAAASUVORK5CYII=",
		ModeType:    "daily",
	}

	habitJSON, err := json.Marshal(invalidHabit)
	if err != nil {
		testing.Fatalf("Failed to marshal habit: %s", err)
	}

	// Send POST request with invalid data
	resp, err := http.Post(
		fmt.Sprintf("%s/habits", baseURL),
		"application/json",
		bytes.NewBuffer(habitJSON),
	)
	if err != nil {
		testing.Fatalf("Failed to send request: %s", err)
	}
	defer resp.Body.Close()

	// Should get a 400 Bad Request
	util.IsEqual(testing, resp.StatusCode, http.StatusBadRequest)

	// Decode error response
	var errorResp map[string]interface{}
	if err = json.NewDecoder(resp.Body).Decode(&errorResp); err != nil {
		testing.Fatalf("Failed to decode error response: %s", err)
	}

	// Verify error message contains validation information
	errorMsg, ok := errorResp["error"].(string)
	if !ok {
		testing.Fatalf("Error response does not contain 'error' field")
	}

	if !strings.Contains(strings.ToLower(errorMsg), "validation") {
		testing.Errorf("Error message does not mention validation: %s", errorMsg)
	}

	// Test case 2: Invalid UUID format
	getResp, err := http.Get(fmt.Sprintf("%s/habits/not-a-valid-uuid", baseURL))
	if err != nil {
		testing.Fatalf("Failed to get habit with invalid ID: %s", err)
	}
	defer getResp.Body.Close()

	// Should get a 400 Bad Request or 404 Not Found
	if getResp.StatusCode != http.StatusBadRequest && getResp.StatusCode != http.StatusNotFound {
		testing.Errorf("Expected status code 400 or 404, got %d", getResp.StatusCode)
	}
}

func TestApiSmoke_NotFoundErrors(testing *testing.T) {
	testing.Parallel()

	// Generate a random UUID that shouldn't exist in the database
	randomUUID := uuid.New().String()

	// Try to get a non-existent habit
	resp, err := http.Get(fmt.Sprintf("%s/habits/%s", baseURL, randomUUID))
	if err != nil {
		testing.Fatalf("Failed to get non-existent habit: %s", err)
	}
	defer resp.Body.Close()

	// Should get a 404 Not Found
	util.IsEqual(testing, resp.StatusCode, http.StatusNotFound)

	// Try to update a non-existent habit
	updateHabit := habit.JsonHabit{
		ID:          randomUUID,
		Description: "This habit doesn't exist",
		ColourHex:   "#FF5733",
		IconBase64:  "data:image/png;base64,iVBORw0KGgoAAAANSUhEUgAAAAEAAAABCAQAAAC1HAwCAAAAC0lEQVR42mNk+A8AAQUBAScY42YAAAAASUVORK5CYII=",
		ModeType:    "daily",
	}

	updateJSON, err := json.Marshal(updateHabit)
	if err != nil {
		testing.Fatalf("Failed to marshal habit: %s", err)
	}

	// Create PUT request
	req, err := http.NewRequest(
		http.MethodPut,
		fmt.Sprintf("%s/habits/%s", baseURL, randomUUID),
		bytes.NewBuffer(updateJSON),
	)
	if err != nil {
		testing.Fatalf("Failed to create PUT request: %s", err)
	}
	req.Header.Set("Content-Type", "application/json")

	// Send PUT request
	client := &http.Client{}
	putResp, err := client.Do(req)
	if err != nil {
		testing.Fatalf("Failed to update non-existent habit: %s", err)
	}
	defer putResp.Body.Close()

	// Should get a 404 Not Found
	util.IsEqual(testing, putResp.StatusCode, http.StatusNotFound)

	// Try to delete a non-existent habit
	deleteReq, err := http.NewRequest(
		http.MethodDelete,
		fmt.Sprintf("%s/habits/%s", baseURL, randomUUID),
		nil,
	)
	if err != nil {
		testing.Fatalf("Failed to create DELETE request: %s", err)
	}

	// Send DELETE request
	deleteResp, err := client.Do(deleteReq)
	if err != nil {
		testing.Fatalf("Failed to delete non-existent habit: %s", err)
	}
	defer deleteResp.Body.Close()

	// Should get a 404 Not Found
	util.IsEqual(testing, deleteResp.StatusCode, http.StatusNotFound)
}

func TestApiSmoke_ConcurrentOperations(testing *testing.T) {
	testing.Parallel()

	// Number of concurrent operations
	const numOperations = 5

	// Create a wait group to wait for all goroutines to finish
	var wg sync.WaitGroup
	wg.Add(numOperations)

	// Create a channel to collect errors
	errorChan := make(chan error, numOperations)

	// Create habits concurrently
	for i := 0; i < numOperations; i++ {
		go func(index int) {
			defer wg.Done()

			// Create a new habit
			newHabit := habit.JsonHabit{
				Description: fmt.Sprintf("Concurrent test habit %d", index),
				ColourHex:   "#FF5733",
				IconBase64:  "data:image/png;base64,iVBORw0KGgoAAAANSUhEUgAAAAEAAAABCAQAAAC1HAwCAAAAC0lEQVR42mNk+A8AAQUBAScY42YAAAAASUVORK5CYII=",
				ModeType:    "daily",
			}

			habitJSON, err := json.Marshal(newHabit)
			if err != nil {
				errorChan <- fmt.Errorf("Failed to marshal habit %d: %s", index, err)
				return
			}

			// Send POST request to create habit
			resp, err := http.Post(
				fmt.Sprintf("%s/habits", baseURL),
				"application/json",
				bytes.NewBuffer(habitJSON),
			)
			if err != nil {
				errorChan <- fmt.Errorf("Failed to create habit %d: %s", index, err)
				return
			}
			defer resp.Body.Close()

			// Check response status code
			if resp.StatusCode != http.StatusCreated {
				errorChan <- fmt.Errorf("Expected status code %d, got %d for habit %d", http.StatusCreated, resp.StatusCode, index)
				return
			}

			// Decode response
			var createdHabit habit.JsonHabit
			if err = json.NewDecoder(resp.Body).Decode(&createdHabit); err != nil {
				errorChan <- fmt.Errorf("Failed to decode created habit %d: %s", index, err)
				return
			}

			// Verify habit was created with correct data
			if createdHabit.Description != newHabit.Description {
				errorChan <- fmt.Errorf("Description mismatch for habit %d. Got %s, expected %s",
					index, createdHabit.Description, newHabit.Description)
				return
			}

			// Verify ID was generated
			if _, err = uuid.Parse(createdHabit.ID); err != nil {
				errorChan <- fmt.Errorf("Invalid UUID returned for habit %d: %s", index, err)
				return
			}
		}(i)
	}

	// Wait for all goroutines to finish
	wg.Wait()
	close(errorChan)

	// Check if there were any errors
	var errors []string
	for err := range errorChan {
		errors = append(errors, err.Error())
	}

	if len(errors) > 0 {
		testing.Fatalf("Concurrent operations failed with errors: %s", strings.Join(errors, ", "))
	}
}
