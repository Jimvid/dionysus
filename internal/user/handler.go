package user

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/jimvid/dionysus/internal/jwt"
	"github.com/jimvid/dionysus/internal/model"
)

type UserHandler struct {
	userService *UserService
}

func NewUserHandler(userService *UserService) *UserHandler {
	return &UserHandler{
		userService: userService,
	}
}

func (h *UserHandler) RegisterUserHandler(w http.ResponseWriter, r *http.Request) {
	var userRequest model.UserRequest
	body, err := io.ReadAll(r.Body)
	defer r.Body.Close()

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = json.Unmarshal([]byte(body), &userRequest)

	if userRequest.Username == "" || userRequest.Password == "" {
		http.Error(w, "invalid request, the fields cannot be empty", http.StatusBadRequest)
		return
	}

	if userRequest.Password != userRequest.ConfirmPassword {
		http.Error(w, "Passwords does not match", http.StatusBadRequest)
		return
	}

	userExist, err := h.userService.DoesUserExist(userRequest.Username)
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	if userExist {
		http.Error(w, "User already exists", http.StatusConflict)
		return
	}

	newUser, err := h.userService.NewUserWithHashedPassword(userRequest)
	if err != nil {
		http.Error(w, "Could not generate hashed password", http.StatusInternalServerError)
		return
	}

	// we know that this user does not exist
	err = h.userService.InsertUser(newUser)
	if err != nil {
		http.Error(w, "Error inserting user", http.StatusInternalServerError)
		return
	}

	successMsg := fmt.Sprintf(`registered with: %s"}`, userRequest)
	_, _ = w.Write([]byte(successMsg))
}

func (h *UserHandler) LoginUserHandler(w http.ResponseWriter, r *http.Request) {
	var loginRequest model.UserLoginRequest

	body, err := io.ReadAll(r.Body)
	defer r.Body.Close()
	err = json.Unmarshal([]byte(body), &loginRequest)

	if err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	newUser, err := h.userService.GetUser(loginRequest.Username)
	if err != nil {
		http.Error(w, error.Error(err), http.StatusInternalServerError)
		return
	}

	if !h.userService.ValidatePassword(newUser.PasswordHash, loginRequest.Password) {
		http.Error(w, "Invalid user credentials", http.StatusUnauthorized)
		return
	}

	accessToken, err := jwt.CreateToken(newUser)
	if err != nil {
		http.Error(w, "Could not create token", http.StatusInternalServerError)
		return
	}

	successMsg := fmt.Sprintf(`{"access_token": "%s"}`, accessToken)

	_, _ = w.Write([]byte(successMsg))
}
