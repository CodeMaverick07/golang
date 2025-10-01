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
	workout, err := wh.WorkoutStore.GetWorkoutByID(workoutID)
	if err != nil {
		http.Error(w, "Workout not found", http.StatusNotFound)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(workout)
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

func (wh *WorkoutHandler) HandleUpdateWorkoutById(w http.ResponseWriter, r *http.Request) {
	paramsWorkoutId := chi.URLParam(r, "id")
	if paramsWorkoutId == "" {
		http.NotFound(w, r)
		return
	}
	workoutId, err := strconv.ParseInt(paramsWorkoutId, 10, 64)
	if err != nil {
		http.NotFound(w, r)
		return
	}
	existingWorkout, err := wh.WorkoutStore.GetWorkoutByID(workoutId)
	if err != nil {
		http.Error(w, "failed to fetch existing workouts", http.StatusInternalServerError)
		return
	}
	if existingWorkout == nil {
		http.NotFound(w, r)
		return
	}
	var updateWorkoutRequest struct {
		Title           string               `json:"title"`
		Description     string               `json:"description"`
		DurationMinutes int                  `json:"duration_minutes"`
		CaloriesBurned  int                  `json:"calories_burned"`
		Entries         []store.WorkoutEntry `json:"entries"`
	}
	err = json.NewDecoder(r.Body).Decode(&updateWorkoutRequest)
	if err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}
	if updateWorkoutRequest.Title != "" {
		existingWorkout.Title = updateWorkoutRequest.Title
	}
	if updateWorkoutRequest.Description != "" {
		existingWorkout.Description = updateWorkoutRequest.Description
	}
	if updateWorkoutRequest.DurationMinutes != 0 {
		existingWorkout.DurationMinutes = updateWorkoutRequest.DurationMinutes
	}
	if updateWorkoutRequest.CaloriesBurned != 0 {
		existingWorkout.CaloriesBurned = updateWorkoutRequest.CaloriesBurned
	}
	if updateWorkoutRequest.Entries != nil {
		existingWorkout.Entries = updateWorkoutRequest.Entries
	}
	err = wh.WorkoutStore.UpdateWorkout(existingWorkout)
	if err != nil {
		http.Error(w, "failed to update workout", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(existingWorkout)

}
