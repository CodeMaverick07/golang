package api

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"regexp"

	"github.com/codemaverick07/api/internals/store"
	"github.com/codemaverick07/api/internals/utils"
)

type registerUserRequest struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
	Bio      string `json:"bio"`
}

type UserHandler struct {
	UserStore store.UserStore
	logger    *log.Logger
}

func NewUserHandler(UserStore store.UserStore, logger *log.Logger) *UserHandler {
	return &UserHandler{
		UserStore: UserStore,
		logger:    logger,
	}
}

func (h *UserHandler) ValidateRegisterRequest(req *registerUserRequest) error {
	if req.Username == "" {
		return errors.New("name is too short")
	}
	if len(req.Username) > 50 {
		return errors.New("username cant be more than 50letters")
	}
	if req.Email == "" {
		return errors.New("email is not provided")
	}
	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	if !emailRegex.MatchString(req.Email) {
		return errors.New("email is not correct")
	}
	if req.Password == "" {
		return errors.New("password is required")
	}
	return nil
}

func (p *UserHandler) HandleRegisterUser(w http.ResponseWriter, r *http.Request) {
	var req registerUserRequest
	err := json.NewDecoder(r.Body).Decode(&req)

	if err != nil {
		p.logger.Printf("ERROR: error while decoding request: %v", err)
		utils.WriteJSON(w, http.StatusBadRequest, utils.Envelope{"error": err.Error()})
		return
	}
	err = p.ValidateRegisterRequest(&req)
	if err != nil {
		p.logger.Printf("ERROR: error while decoding request: %v", err)
		utils.WriteJSON(w, http.StatusBadRequest, utils.Envelope{"error": err.Error()})
		return
	}

	user := &store.User{
		UserName: req.Username,
		Email:    req.Email,
	}

	if req.Bio != "" {
		user.Bio = req.Bio
	}

	err = user.PasswordHash.Set(req.Password)
	if err != nil {
		p.logger.Printf("ERROR: error while hashing the password: %v", err)
		utils.WriteJSON(w, http.StatusInternalServerError, utils.Envelope{"error": "internal server error"})
		return
	}
	err = p.UserStore.CreateUser(user)
	if err != nil {
		p.logger.Printf("ERROR: error while registering user: %v", err)
		utils.WriteJSON(w, http.StatusInternalServerError, utils.Envelope{"error": "internal server error"})
		return
	}

	utils.WriteJSON(w, http.StatusOK, utils.Envelope{"data": user})

}
