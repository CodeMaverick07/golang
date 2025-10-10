package api

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/codemaverick07/api/internals/store"
	"github.com/codemaverick07/api/internals/tokens"
	"github.com/codemaverick07/api/internals/utils"
)

type TokenHandler struct {
	TokenStore store.TokenStore
	UserStore  store.UserStore
	logger     *log.Logger
}

type CreateTokenRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func NewTokenHandler(tokenStore store.TokenStore, userStore store.UserStore, logger *log.Logger) *TokenHandler {
	return &TokenHandler{
		TokenStore: tokenStore,
		UserStore:  userStore,
		logger:     logger,
	}
}

func (h *TokenHandler) HandleCreateToken(w http.ResponseWriter, r *http.Request) {
	var req CreateTokenRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		h.logger.Printf("ERROR: error while decoding request: %v", err)
		utils.WriteJSON(w, http.StatusBadRequest, utils.Envelope{"error": "error while decoding token"})
		return
	}

	user, err := h.UserStore.GetUserByUserName(req.Username)
	if err != nil {
		if err == sql.ErrNoRows {
			h.logger.Printf("ERROR: user not found: %s", req.Username)
			utils.WriteJSON(w, http.StatusUnauthorized, utils.Envelope{"error": "invalid credentials"})
			return
		}
		h.logger.Printf("ERROR: error while getting user by name: %v", err)
		utils.WriteJSON(w, http.StatusInternalServerError, utils.Envelope{"error": "internal server error"})
		return
	}
	passwordDoMatch, err := user.PasswordHash.Matches(req.Password)

	if err != nil {
		h.logger.Printf("ERROR: password dose not match in handleCreateToken: %v", err)
		utils.WriteJSON(w, http.StatusBadRequest, utils.Envelope{"error": "internal server error"})
		return
	}

	if !passwordDoMatch {
		utils.WriteJSON(w, http.StatusInternalServerError, utils.Envelope{"error": "password dose not match"})
		return
	}

	token, err := h.TokenStore.CreateNewToken(user.ID, 24*time.Hour, tokens.ScopeAuth)
	if err != nil {
		h.logger.Printf("ERROR: error while creating token: %v", err)
		utils.WriteJSON(w, http.StatusBadRequest, utils.Envelope{"error": "internal server error "})
		return
	}

	utils.WriteJSON(w, http.StatusAccepted, utils.Envelope{"data": token})

}
