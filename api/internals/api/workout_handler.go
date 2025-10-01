package api

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/codemaverick07/api/internals/store"
	"github.com/codemaverick07/api/internals/utils"
)

type WorkoutHandler struct {
	WorkoutStore store.WorkoutStore
	logger       *log.Logger
}

func NewWorkoutHandler(workoutStore store.WorkoutStore, logger *log.Logger) *WorkoutHandler {
	return &WorkoutHandler{
		WorkoutStore: workoutStore,
		logger:       logger,
	}
}

func (wh *WorkoutHandler) HandleGetWorkByID(w http.ResponseWriter, r *http.Request) {

	workoutID, err := utils.ReadParamId(r)
	if err != nil {
		wh.logger.Printf("ERROR: readId param %v", err)
		utils.WriteJSON(w, http.StatusBadRequest, utils.Envelope{"error": err})
		return
	}
	workout, err := wh.WorkoutStore.GetWorkoutByID(workoutID)
	if err != nil {
		wh.logger.Printf("ERROR: readId param %v", err)
		utils.WriteJSON(w, http.StatusBadRequest, utils.Envelope{"error": err})
		return
	}
	utils.WriteJSON(w, http.StatusOK, utils.Envelope{"data": workout})
}

func (wh *WorkoutHandler) HandleCreateWorkout(w http.ResponseWriter, r *http.Request) {
	var workout store.Workout
	err := json.NewDecoder(r.Body).Decode(&workout)
	if err != nil {
		wh.logger.Printf("ERROR: bad data createWorkout %v", err)
		utils.WriteJSON(w, http.StatusBadRequest, utils.Envelope{"error": err})
		return
	}
	createdWorkout, err := wh.WorkoutStore.CreateWorkout(&workout)
	if err != nil {
		wh.logger.Printf("ERROR: creating new workout %v", err)
		utils.WriteJSON(w, http.StatusInternalServerError, utils.Envelope{"error": err})
		return
	}
	utils.WriteJSON(w, http.StatusOK, utils.Envelope{"data": createdWorkout})
}

func (wh *WorkoutHandler) HandleUpdateWorkoutById(w http.ResponseWriter, r *http.Request) {

	workoutId, err := utils.ReadParamId(r)
	if err != nil {
		wh.logger.Printf("ERROR: wrong id %v \n", err)
		utils.WriteJSON(w, http.StatusNotFound, utils.Envelope{"error": "wrong id provided"})
		return
	}
	existingWorkout, err := wh.WorkoutStore.GetWorkoutByID(workoutId)
	if err != nil {
		wh.logger.Printf("ERROR: error while getting workout %v \n", err)
		utils.WriteJSON(w, http.StatusBadRequest, utils.Envelope{"error": "error while getting workoutById"})
		return
	}
	if existingWorkout == nil {
		wh.logger.Printf("ERROR: no workout found %v \n", err)
		utils.WriteJSON(w, http.StatusNotFound, utils.Envelope{"error": "no workout found"})
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
		wh.logger.Printf("ERROR: bad data %v \n", err)
		utils.WriteJSON(w, http.StatusBadRequest, utils.Envelope{"error": "bad data"})
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
		wh.logger.Printf("ERROR: failed to update error %v \n", err)
		utils.WriteJSON(w, http.StatusInternalServerError, utils.Envelope{"error": "failed to update workout"})
		return
	}
	utils.WriteJSON(w, http.StatusOK, utils.Envelope{"data": existingWorkout})

}

func (wh *WorkoutHandler) HandleDeleteWorkoutById(w http.ResponseWriter, r *http.Request) {

	workoutId, err := utils.ReadParamId(r)
	if err != nil {
		wh.logger.Printf("ERROR: wrong id %v \n", err)
		utils.WriteJSON(w, http.StatusNotFound, utils.Envelope{"error": "wrong id provided"})
		return
	}

	err = wh.WorkoutStore.DeleteWorkout(workoutId)
	if err != nil {
		wh.logger.Printf("ERROR: nothing found to delete %v \n", err)
		utils.WriteJSON(w, http.StatusNotFound, utils.Envelope{"error": "nothing found to delete"})
		return
	}

	utils.WriteJSON(w, http.StatusOK, utils.Envelope{"data": workoutId})

}
