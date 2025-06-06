package integration

import (
	"encoding/json"
	"fmt"
	"habitgobackend/cmd/api/resource/habit"
	"habitgobackend/test/util"
	"net/http"
	"testing"
)

func TestApiSmoke_GetHabits(testing *testing.T) {
	testing.Parallel()

	resp, err := http.Get(fmt.Sprintf("http://localhost:%v/v1/habits", 8080))

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
