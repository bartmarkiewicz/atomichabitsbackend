package util

import (
	"habitgobackend/cmd/api/resource/habit"
	"testing"
)

func NoError(testing *testing.T, err error) {
	if err != nil {
		testing.Fatalf("err: %e", err)
	}
}

func IsEqual[T comparable](testing *testing.T, obj, other T) {
	if obj != other {
		testing.Fatalf("Is not equal: %v, other: %v", obj, other)
	}
}

func HabitsEqual(t *testing.T, actual habit.Habits, expected habit.Habits) {
	if len(actual) != len(expected) {
		t.Fatalf("Habit lists have different lengths. Got %d, expected %d", len(actual), len(expected))
	}

	for i, actualHabit := range actual {
		expectedHabit := expected[i]

		if actualHabit.ID != expectedHabit.ID {
			t.Errorf("Habit %d: ID mismatch. Got %v, expected %v", i, actualHabit.ID, expectedHabit.ID)
		}
		if actualHabit.Description != expectedHabit.Description {
			t.Errorf("Habit %d: Description mismatch. Got %v, expected %v", i, actualHabit.Description, expectedHabit.Description)
		}
		if actualHabit.ColourHex != expectedHabit.ColourHex {
			t.Errorf("Habit %d: ColourHex mismatch. Got %v, expected %v", i, actualHabit.ColourHex, expectedHabit.ColourHex)
		}
		if actualHabit.IconBase64 != expectedHabit.IconBase64 {
			t.Errorf("Habit %d: IconBase64 mismatch. Got %v, expected %v", i, actualHabit.IconBase64, expectedHabit.IconBase64)
		}
		if actualHabit.ModeType != expectedHabit.ModeType {
			t.Errorf("Habit %d: ModeType mismatch. Got %v, expected %v", i, actualHabit.ModeType, expectedHabit.ModeType)
		}
	}
}
