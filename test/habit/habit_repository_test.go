package habit

import (
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/google/uuid"
	"habitgobackend/cmd/api/resource/habit"
	"habitgobackend/test/util"
	"testing"
)

func TestRepository_GetHabits(testing *testing.T) {
	testing.Parallel()

	database, mock, err := util.NewMockDatabase()
	util.NoError(testing, err)

	repository := habit.NewRepository(database)

	expectedHabits := habit.Habits{
		{
			ID:          uuid.New(),
			Description: "Random habit",
			ColourHex:   "#000000",
			IconBase64:  "data:image/png;base64,iVBORw0KGgoAAAANSUhEUgAAAAEAAAABCAYAAAAfFcSJAAAADUlEQVR42mNk+P+ErkJggg==",
			ModeType:    "daily",
		},
		{
			ID:          uuid.New(),
			Description: "Drink water",
			ColourHex:   "#ffffff",
			IconBase64:  "data:image/png;base64,iVBORw0KGgoAAk+P+ErkJggg==",
			ModeType:    "monthly",
		},
		{
			ID:          uuid.New(),
			Description: "Write some code",
			ColourHex:   "#bbbbbb",
			IconBase64:  "data:image/png;base64,iVBORwsdsdasgoAAk+P+ErkJggg==",
			ModeType:    "daily",
		},
	}

	rows := sqlmock.NewRows([]string{"ID", "Description", "ColourHex", "IconBase64", "ModeType"}).
		AddRow(expectedHabits[0].ID, expectedHabits[0].Description, expectedHabits[0].ColourHex,
			expectedHabits[0].IconBase64, expectedHabits[0].ModeType).
		AddRow(expectedHabits[1].ID, expectedHabits[1].Description, expectedHabits[1].ColourHex,
			expectedHabits[1].IconBase64, expectedHabits[1].ModeType).
		AddRow(expectedHabits[2].ID, expectedHabits[2].Description, expectedHabits[2].ColourHex,
			expectedHabits[2].IconBase64, expectedHabits[2].ModeType)

	mock.ExpectQuery("SELECT (.+) FROM \"habits\"").WillReturnRows(rows)

	habits, err := repository.GetHabits()

	util.NoError(testing, err)

	util.IsEqual(testing, len(habits), 3)
	util.HabitsEqual(testing, habits, expectedHabits)
}

func TestRepository_CreateHabit(testing *testing.T) {
	testing.Parallel()

	database, mock, err := util.NewMockDatabase()
	util.NoError(testing, err)

	repository := habit.NewRepository(database)

	id := uuid.New()

	mock.ExpectBegin()
	mock.ExpectExec("^INSERT INTO \"habits\" ").
		WithArgs(id, "Description", "#000000", "data:image/png;base64,iVBORw0KGgoAAAANSUhE+ErkJggg==", "daily").
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	newHabit := &habit.Habit{ID: id, Description: "Description",
		IconBase64: "data:image/png;base64,iVBORw0KGgoAAAANSUhE+ErkJggg==",
		ColourHex:  "#000000", ModeType: "daily"}

	result, err := repository.CreateHabit(newHabit)
	util.NoError(testing, err)
	util.HabitsEqual(testing, habit.Habits{newHabit}, habit.Habits{result})
}

func TestRepository_UpdateHabit(testing *testing.T) {
	testing.Parallel()

	database, mock, err := util.NewMockDatabase()
	util.NoError(testing, err)

	repository := habit.NewRepository(database)

	id := uuid.New()
	_ = sqlmock.NewRows([]string{"id", "description", "colour_hex", "icon_base_64", "mode_type"}).
		AddRow(id, "Description", "Hex", "Icon", "Mode").
		AddRow(uuid.New(), "Random Description", "Colour Hex", "Icon base 64", "Mode Type")

	mock.ExpectBegin()
	mock.ExpectExec("^UPDATE \"habits\" SET").
		WithArgs("Updated Description", "Updated Hex", "Updated Icon", "Updated Mode", id).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	newHabit := &habit.Habit{ID: id, Description: "Updated Description",
		IconBase64: "Updated Icon",
		ColourHex:  "Updated Hex", ModeType: "Updated Mode"}

	result, err := repository.UpdateHabit(newHabit)
	util.NoError(testing, err)
	util.IsEqual(testing, 1, result)
}

func TestRepository_GetHabit(testing *testing.T) {
	testing.Parallel()

	database, mock, err := util.NewMockDatabase()
	util.NoError(testing, err)

	repository := habit.NewRepository(database)

	expectedHabit := habit.Habit{
		ID:          uuid.New(),
		Description: "Random habit",
		ColourHex:   "#000000",
		IconBase64:  "data:image/png;base64,iVBORw0KGgoAAAANSUhEUgAAAAEAAAABCAYAAAAfFcSJAAAADUlEQVR42mNk+P+ErkJggg==",
		ModeType:    "daily",
	}

	rows := sqlmock.NewRows([]string{"ID", "Description", "ColourHex", "IconBase64", "ModeType"}).
		AddRow(expectedHabit.ID, expectedHabit.Description, expectedHabit.ColourHex, expectedHabit.IconBase64,
			expectedHabit.ModeType)

	mock.ExpectQuery("SELECT (.+) FROM \"habits\" WHERE (.+)").
		WithArgs(expectedHabit.ID, 1).
		WillReturnRows(rows)

	result, err := repository.GetHabit(expectedHabit.ID)

	util.NoError(testing, err)

	util.HabitsEqual(testing, habit.Habits{result}, habit.Habits{&expectedHabit})
}

func TestRepository_DeleteHabit(testing *testing.T) {
	testing.Parallel()

	database, mock, err := util.NewMockDatabase()
	util.NoError(testing, err)

	repository := habit.NewRepository(database)

	expectedHabit := habit.Habit{
		ID:          uuid.New(),
		Description: "Random habit",
		ColourHex:   "#000000",
		IconBase64:  "data:image/png;base64,iVBORw0KGgoAAAANSUhEUgAAAAEAAAABCAYAAAAfFcSJAAAADUlEQVR42mNk+P+ErkJggg==",
		ModeType:    "daily",
	}

	mock.ExpectBegin()
	mock.ExpectExec("DELETE * FROM \"habits\" WHERE (.+)").
		WithArgs(expectedHabit.ID).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	result, err := repository.DeleteHabit(expectedHabit.ID)

	util.NoError(testing, err)
	util.IsEqual(testing, result, 1)
}
