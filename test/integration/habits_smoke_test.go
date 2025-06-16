package integration

import (
	"bytes"
	"encoding/json"
	"fmt"
	headers "habitgobackend/cmd/api/resource/common/helpers"
	"habitgobackend/cmd/api/resource/habit"
	"habitgobackend/test/util"
	"io"
	"net/http"
	"strings"
	"testing"

	"github.com/google/uuid"
)

const (
	baseURL = "http://localhost:8080/v1"
)

func ClearDb(testing *testing.T) {
	resp, err := http.Get(fmt.Sprintf("%s/habits", baseURL))
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	habits := &habit.Habits{}

	if err = json.NewDecoder(resp.Body).Decode(habits); err != nil {
		testing.Fatalf("Failed to decode habits: %s", err)
	}

	for item := range *habits {
		request, err := http.NewRequest("DELETE",
			fmt.Sprintf("%s/habits/%v", baseURL, (*habits)[item].ID.String()),
			bytes.NewBuffer([]byte{}),
		)
		if err != nil {
			testing.Fatalf("Failed to create DELETE request: %s", err)
			return
		}

		client := &http.Client{}
		delRequest, err := client.Do(request)
		if err != nil || (delRequest.StatusCode != http.StatusOK && delRequest.StatusCode != http.StatusNotFound) {
			testing.Fatalf("Failed to create DELETE request: %v status: %v", err, delRequest.StatusCode)
			return
		}

	}
}

func TestSmoke_GetHabits(testing *testing.T) {
	ClearDb(testing)
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

func TestSmoke_CreateHabit(testing *testing.T) {
	ClearDb(testing)
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

	var createdHabitId = resp.Header.Get(headers.CREATED_ID)

	createdHabit := getHabit(testing, createdHabitId)

	util.IsEqual(testing, createdHabit.Description, newHabit.Description)
	util.IsEqual(testing, createdHabit.ColourHex, newHabit.ColourHex)
	util.IsEqual(testing, createdHabit.IconBase64, newHabit.IconBase64)
	util.IsEqual(testing, createdHabit.ModeType, newHabit.ModeType)

	_, err = uuid.Parse(createdHabit.ID)
	if err != nil {
		testing.Fatalf("Invalid UUID returned: %s", err)
	}
}

func TestSmoke_GetHabitById(testing *testing.T) {
	ClearDb(testing)
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

	createdHabitId := resp.Header.Get(headers.CREATED_ID)

	retrievedHabit := getHabit(testing, createdHabitId)

	util.IsEqual(testing, retrievedHabit.ID, createdHabitId)
	util.IsEqual(testing, retrievedHabit.Description, newHabit.Description)
	util.IsEqual(testing, retrievedHabit.ColourHex, newHabit.ColourHex)
	util.IsEqual(testing, retrievedHabit.IconBase64, newHabit.IconBase64)
	util.IsEqual(testing, retrievedHabit.ModeType, newHabit.ModeType)
}

func TestSmoke_UpdateHabit(testing *testing.T) {
	ClearDb(testing)
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

	createdHabitId := resp.Header.Get(headers.CREATED_ID)

	updatedHabit := habit.JsonHabit{
		ID:          createdHabitId,
		Description: "Updated test habit",
		ColourHex:   "#9933FF",
		IconBase64:  "data:image/png;base64,FISSHH",
		ModeType:    "yearly",
	}

	updatedJSON, err := json.Marshal(updatedHabit)
	if err != nil {
		testing.Fatalf("Failed to marshal updated habit: %s", err)
	}

	req, err := http.NewRequest(
		http.MethodPut,
		fmt.Sprintf("%s/habits/%s", baseURL, createdHabitId),
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

	returnedHabit := getHabit(testing, createdHabitId)

	util.IsEqual(testing, returnedHabit.ID, updatedHabit.ID)
	util.IsEqual(testing, returnedHabit.Description, updatedHabit.Description)
	util.IsEqual(testing, returnedHabit.ColourHex, updatedHabit.ColourHex)
	util.IsEqual(testing, returnedHabit.IconBase64, updatedHabit.IconBase64)
	util.IsEqual(testing, returnedHabit.ModeType, updatedHabit.ModeType)
}

func TestSmoke_InvalidScenarios(testing *testing.T) {
	ClearDb(testing)
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

	util.IsEqual(testing, resp.StatusCode, http.StatusUnprocessableEntity)

	var errorResp map[string]interface{}
	if err = json.NewDecoder(resp.Body).Decode(&errorResp); err != nil {
		testing.Fatalf("Failed to decode error response: %s", err)
	}

	errorMsg, ok := errorResp["error"].(string)
	if !ok {
		testing.Fatalf("Error response does not contain 'error' field")
	}

	if !strings.Contains(errorMsg, "Could not create entity") {
		testing.Errorf("Error message expected: Could not create entity, actual: %s", errorMsg)
	}

	getResp, err := http.Get(fmt.Sprintf("%s/habits/not-a-valid-uuid", baseURL))
	if err != nil {
		testing.Fatalf("Failed to get habit with invalid ID: %s", err)
	}
	defer getResp.Body.Close()

	if getResp.StatusCode != http.StatusBadRequest {
		testing.Errorf("Expected status code 400, got %d", getResp.StatusCode)
	}
}

func TestSmoke_NotFoundErrors(testing *testing.T) {
	ClearDb(testing)
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

func getHabit(testing *testing.T, habitId string) habit.JsonHabit {
	getResp, err := http.Get(fmt.Sprintf("%s/habits/%s", baseURL, habitId))
	if err != nil {
		testing.Fatalf("Failed to get habit by ID: %s", err)
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {

		}
	}(getResp.Body)

	util.IsEqual(testing, getResp.StatusCode, http.StatusOK)

	var retrievedHabit habit.JsonHabit
	if err = json.NewDecoder(getResp.Body).Decode(&retrievedHabit); err != nil {
		testing.Fatalf("Failed to decode retrieved habit: %s", err)
	}
	return retrievedHabit
}
