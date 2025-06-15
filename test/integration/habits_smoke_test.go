package integration

import (
	"bytes"
	"encoding/json"
	"fmt"
	"habitgobackend/cmd/api/resource/habit"
	"habitgobackend/test/util"
	"net/http"
	"strings"
	"testing"

	"github.com/google/uuid"
)

const (
	baseURL = "http://localhost:8080/v1"
)

func TestApiSmoke_GetHabits(testing *testing.T) {
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

	resp, err := http.Post(
		fmt.Sprintf("%s/habits", baseURL),
		"application/json",
		bytes.NewBuffer(habitJSON),
	)
	if err != nil {
		testing.Fatalf("Failed to create habit: %s", err)
	}
	defer resp.Body.Close()

	util.IsEqual(testing, resp.StatusCode, http.StatusCreated)

	var createdHabit habit.JsonHabit
	if err = json.NewDecoder(resp.Body).Decode(&createdHabit); err != nil {
		testing.Fatalf("Failed to decode created habit: %s", err)
	}

	util.IsEqual(testing, createdHabit.Description, newHabit.Description)
	util.IsEqual(testing, createdHabit.ColourHex, newHabit.ColourHex)
	util.IsEqual(testing, createdHabit.IconBase64, newHabit.IconBase64)
	util.IsEqual(testing, createdHabit.ModeType, newHabit.ModeType)

	_, err = uuid.Parse(createdHabit.ID)
	if err != nil {
		testing.Fatalf("Invalid UUID returned: %s", err)
	}
}

func TestApiSmoke_GetHabitById(testing *testing.T) {
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

	resp, err := http.Post(
		fmt.Sprintf("%s/habits", baseURL),
		"application/json",
		bytes.NewBuffer(habitJSON),
	)
	if err != nil {
		testing.Fatalf("Failed to create habit: %s", err)
	}

	var createdHabit habit.JsonHabit
	if err = json.NewDecoder(resp.Body).Decode(&createdHabit); err != nil {
		testing.Fatalf("Failed to decode created habit: %s", err)
	}
	resp.Body.Close()

	getResp, err := http.Get(fmt.Sprintf("%s/habits/%s", baseURL, createdHabit.ID))
	if err != nil {
		testing.Fatalf("Failed to get habit by ID: %s", err)
	}
	defer getResp.Body.Close()

	util.IsEqual(testing, getResp.StatusCode, http.StatusOK)

	var retrievedHabit habit.JsonHabit
	if err = json.NewDecoder(getResp.Body).Decode(&retrievedHabit); err != nil {
		testing.Fatalf("Failed to decode retrieved habit: %s", err)
	}

	util.IsEqual(testing, retrievedHabit.ID, createdHabit.ID)
	util.IsEqual(testing, retrievedHabit.Description, newHabit.Description)
	util.IsEqual(testing, retrievedHabit.ColourHex, newHabit.ColourHex)
	util.IsEqual(testing, retrievedHabit.IconBase64, newHabit.IconBase64)
	util.IsEqual(testing, retrievedHabit.ModeType, newHabit.ModeType)
}

func TestApiSmoke_UpdateHabit(testing *testing.T) {
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

	resp, err := http.Post(
		fmt.Sprintf("%s/habits", baseURL),
		"application/json",
		bytes.NewBuffer(habitJSON),
	)
	if err != nil {
		testing.Fatalf("Failed to create habit: %s", err)
	}

	var createdHabit habit.JsonHabit
	if err = json.NewDecoder(resp.Body).Decode(&createdHabit); err != nil {
		testing.Fatalf("Failed to decode created habit: %s", err)
	}
	resp.Body.Close()

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

	req, err := http.NewRequest(
		http.MethodPut,
		fmt.Sprintf("%s/habits/%s", baseURL, createdHabit.ID),
		bytes.NewBuffer(updatedJSON),
	)
	if err != nil {
		testing.Fatalf("Failed to create PUT request: %s", err)
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	putResp, err := client.Do(req)
	if err != nil {
		testing.Fatalf("Failed to update habit: %s", err)
	}
	defer putResp.Body.Close()

	util.IsEqual(testing, putResp.StatusCode, http.StatusOK)

	var returnedHabit habit.JsonHabit
	if err = json.NewDecoder(putResp.Body).Decode(&returnedHabit); err != nil {
		testing.Fatalf("Failed to decode updated habit: %s", err)
	}

	util.IsEqual(testing, returnedHabit.ID, updatedHabit.ID)
	util.IsEqual(testing, returnedHabit.Description, updatedHabit.Description)
	util.IsEqual(testing, returnedHabit.ColourHex, updatedHabit.ColourHex)
	util.IsEqual(testing, returnedHabit.IconBase64, updatedHabit.IconBase64)
	util.IsEqual(testing, returnedHabit.ModeType, updatedHabit.ModeType)

	getResp, err := http.Get(fmt.Sprintf("%s/habits/%s", baseURL, createdHabit.ID))
	if err != nil {
		testing.Fatalf("Failed to get updated habit: %s", err)
	}
	defer getResp.Body.Close()

	var retrievedHabit habit.JsonHabit
	if err = json.NewDecoder(getResp.Body).Decode(&retrievedHabit); err != nil {
		testing.Fatalf("Failed to decode retrieved habit: %s", err)
	}

	util.IsEqual(testing, retrievedHabit.Description, updatedHabit.Description)
	util.IsEqual(testing, retrievedHabit.ColourHex, updatedHabit.ColourHex)
	util.IsEqual(testing, retrievedHabit.ModeType, updatedHabit.ModeType)
}

func TestApiSmoke_DeleteHabit(testing *testing.T) {
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

	resp, err := http.Post(
		fmt.Sprintf("%s/habits", baseURL),
		"application/json",
		bytes.NewBuffer(habitJSON),
	)
	if err != nil {
		testing.Fatalf("Failed to create habit: %s", err)
	}

	var createdHabit habit.JsonHabit
	if err = json.NewDecoder(resp.Body).Decode(&createdHabit); err != nil {
		testing.Fatalf("Failed to decode created habit: %s", err)
	}
	resp.Body.Close()

	req, err := http.NewRequest(
		http.MethodDelete,
		fmt.Sprintf("%s/habits/%s", baseURL, createdHabit.ID),
		nil,
	)
	if err != nil {
		testing.Fatalf("Failed to create DELETE request: %s", err)
	}

	client := &http.Client{}
	deleteResp, err := client.Do(req)
	if err != nil {
		testing.Fatalf("Failed to delete habit: %s", err)
	}
	defer deleteResp.Body.Close()

	util.IsEqual(testing, deleteResp.StatusCode, http.StatusOK)

	getResp, err := http.Get(fmt.Sprintf("%s/habits/%s", baseURL, createdHabit.ID))
	if err != nil {
		testing.Fatalf("Failed to get deleted habit: %s", err)
	}
	defer getResp.Body.Close()

	util.IsEqual(testing, getResp.StatusCode, http.StatusNotFound)
}

func TestApiSmoke_ValidationErrors(testing *testing.T) {
	invalidHabit := habit.JsonHabit{
		Description: "",
		ColourHex:   "#FF5733",
		IconBase64:  "data:image/png;base64,iVBORw0KGgoAAAANSUhEUgAAAAEAAAABCAQAAAC1HAwCAAAAC0lEQVR42mNk+A8AAQUBAScY42YAAAAASUVORK5CYII=",
		ModeType:    "daily",
	}

	habitJSON, err := json.Marshal(invalidHabit)
	if err != nil {
		testing.Fatalf("Failed to marshal habit: %s", err)
	}

	resp, err := http.Post(
		fmt.Sprintf("%s/habits", baseURL),
		"application/json",
		bytes.NewBuffer(habitJSON),
	)
	if err != nil {
		testing.Fatalf("Failed to send request: %s", err)
	}
	defer resp.Body.Close()

	util.IsEqual(testing, resp.StatusCode, http.StatusBadRequest)

	var errorResp map[string]interface{}
	if err = json.NewDecoder(resp.Body).Decode(&errorResp); err != nil {
		testing.Fatalf("Failed to decode error response: %s", err)
	}

	errorMsg, ok := errorResp["error"].(string)
	if !ok {
		testing.Fatalf("Error response does not contain 'error' field")
	}

	if !strings.Contains(strings.ToLower(errorMsg), "validation") {
		testing.Errorf("Error message does not mention validation: %s", errorMsg)
	}

	getResp, err := http.Get(fmt.Sprintf("%s/habits/not-a-valid-uuid", baseURL))
	if err != nil {
		testing.Fatalf("Failed to get habit with invalid ID: %s", err)
	}
	defer getResp.Body.Close()

	if getResp.StatusCode != http.StatusBadRequest && getResp.StatusCode != http.StatusNotFound {
		testing.Errorf("Expected status code 400 or 404, got %d", getResp.StatusCode)
	}
}

func TestApiSmoke_NotFoundErrors(testing *testing.T) {
	randomUUID := uuid.New().String()

	resp, err := http.Get(fmt.Sprintf("%s/habits/%s", baseURL, randomUUID))
	if err != nil {
		testing.Fatalf("Failed to get non-existent habit: %s", err)
	}
	defer resp.Body.Close()

	util.IsEqual(testing, resp.StatusCode, http.StatusNotFound)

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

	req, err := http.NewRequest(
		http.MethodPut,
		fmt.Sprintf("%s/habits/%s", baseURL, randomUUID),
		bytes.NewBuffer(updateJSON),
	)
	if err != nil {
		testing.Fatalf("Failed to create PUT request: %s", err)
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	putResp, err := client.Do(req)
	if err != nil {
		testing.Fatalf("Failed to update non-existent habit: %s", err)
	}
	defer putResp.Body.Close()

	util.IsEqual(testing, putResp.StatusCode, http.StatusNotFound)

	deleteReq, err := http.NewRequest(
		http.MethodDelete,
		fmt.Sprintf("%s/habits/%s", baseURL, randomUUID),
		nil,
	)
	if err != nil {
		testing.Fatalf("Failed to create DELETE request: %s", err)
	}

	deleteResp, err := client.Do(deleteReq)
	if err != nil {
		testing.Fatalf("Failed to delete non-existent habit: %s", err)
	}
	defer deleteResp.Body.Close()

	util.IsEqual(testing, deleteResp.StatusCode, http.StatusNotFound)
}
