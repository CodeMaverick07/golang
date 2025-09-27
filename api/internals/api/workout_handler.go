package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/codemaverick07/api/internals/store"
	"github.com/go-chi/chi/v5"
)

type WorkoutHandler struct {
	WorkoutStore store.WorkoutStore
}

func NewWorkoutHandler(workoutStore store.WorkoutStore) *WorkoutHandler {
	return &WorkoutHandler{
		WorkoutStore: workoutStore,
	}
}

func (wh *WorkoutHandler) HandleGetWorkByID(w http.ResponseWriter, r *http.Request) {
	paramsWorkoutID := chi.URLParam(r, "id")
	if paramsWorkoutID == "" {
		http.NotFound(w, r)
		return
	}
	workoutID, err := strconv.ParseInt(paramsWorkoutID, 10, 64)
	if err != nil {
		http.NotFound(w, r)
		return
	}
	fmt.Fprintf(w, "workoutId is %d\n", workoutID)
}

func (wh *WorkoutHandler) HandleCreateWorkout(w http.ResponseWriter, r *http.Request) {
	var workout store.Workout
	err := json.NewDecoder(r.Body).Decode(&workout)
	if err != nil {
		fmt.Println("Error decoding request body:", err)
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	createdWorkout, err := wh.WorkoutStore.CreateWorkout(&workout)
	if err != nil {
		fmt.Println("Error decoding request body:", err)
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(createdWorkout)
}
