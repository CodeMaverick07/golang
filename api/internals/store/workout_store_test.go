package store

import (
	"database/sql"
	"testing"

	_ "github.com/jackc/pgx/v4/stdlib"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func SetupTestDB(t *testing.T) *sql.DB {
	db, err := sql.Open("pgx", "host=localhost user=postgres password=postgres dbname=postgres port=5433")
	if err != nil {
	t.Fatalf("opening test db %v", err)
	}
	err = Migrate(db, "../../migrations/")

	if err != nil {
	t.Fatalf("error in migration %v", err)
	}

	_, err = db.Exec(`TRUNCATE workouts,workout_entries CASCADE`)
	if err != nil {
	t.Fatalf("error in truncate db %v", err)
	}
	return db
}

func TestCreateWorkout(t *testing.T) {
	db := SetupTestDB(t)
	defer db.Close()
	store := NewPostgresWorkoutStore(db)

	tests := []struct {
		name    string
		Workout *Workout
		wantErr bool
	}{
		{
			name: "valid workout",
			Workout: &Workout{
				Title:           "push",
				Description:     "dddddddd",
				DurationMinutes: 60,
				CaloriesBurned:  200,
				Entries: []WorkoutEntry{
					{
						ExerciseName: "bench press",
						Sets:         6,
						Reps:         IntPtr(10),
						Weight:       *FloatPtr(77.8),
						Notes:        "nice notes",
						OrderIndex:   1,
					},
				},
			},
			wantErr: false,
		},
		{
			name: "workout with invalid entries",
			Workout: &Workout{
				Title:           "full body",
				Description:     "complete workout ",
				DurationMinutes: 60,
				CaloriesBurned:  200,
				Entries: []WorkoutEntry{
					{
						ExerciseName: "bench press",
						Sets:         6,
						Reps:         IntPtr(10),
						Notes:        "nice notes",
						OrderIndex:   1,
					},
					{
						ExerciseName:    "error exercise",
						Sets:            6,
						Reps:            IntPtr(10),
						DurationSeconds: IntPtr(60),
						Notes:           "nice notes",
						OrderIndex:      1,
					},
				},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			createdWorkout, err := store.CreateWorkout(tt.Workout)
			if tt.wantErr {
				assert.Error(t, err)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, tt.Workout.Title, createdWorkout.Title)
			assert.Equal(t, tt.Workout.Description, createdWorkout.Description)
			assert.Equal(t, tt.Workout.CaloriesBurned, createdWorkout.CaloriesBurned)
			retrieved, err := store.GetWorkoutByID(int64(createdWorkout.ID))
			require.NoError(t, err)
			assert.Equal(t, createdWorkout.ID, retrieved.ID)
			assert.Equal(t, len(retrieved.Entries), len(createdWorkout.Entries))
		})
	}
}

func IntPtr(i int) *int {
	return &i
}

func FloatPtr(i float64) *float64 {
	return &i
}
