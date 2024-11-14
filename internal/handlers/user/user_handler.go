package user

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
)

type UserHandler struct {
	store *UserStore
}

func NewHandler(store *UserStore) *UserHandler {
	return &UserHandler{store: store}
}

func (h *UserHandler) UserRoutes(router *mux.Router) {
	router.HandleFunc("/users/{id}/status", h.GetUserStatus).Methods("GET")
	router.HandleFunc("/users/leaderboard", h.GetTopUserBalance).Methods("GET")
	router.HandleFunc("/users/{id}/task/complete", h.CompleteTask).Methods("POST")
	router.HandleFunc("/users/{id}/referrer", h.SetReferrer).Methods("POST")
}

func (h *UserHandler) GetUserStatus(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userID := vars["id"]

	user, err := h.store.GetUserStatus(userID)
	if err != nil {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
}

func (h *UserHandler) GetTopUserBalance(w http.ResponseWriter, r *http.Request) {
	users, err := h.store.GetTopUsersByBalance(3)
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(users); err != nil {
		http.Error(w, "Error encoding response", http.StatusInternalServerError)
		return
	}
}

func (h *UserHandler) CompleteTask(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userID := vars["id"]

	var requestData struct {
		TaskID int `json:"task_id"`
	}

	if err := json.NewDecoder(r.Body).Decode(&requestData); err != nil {
		http.Error(w, "Error body request", http.StatusBadRequest)
		return
	}

	task, err := h.store.GetTaskByID(requestData.TaskID)
	if err != nil {
		http.Error(w, "Task is not a found", http.StatusNotFound)
		return
	}

	userTask, err := h.store.GetUserTask(userID, requestData.TaskID)
	if err == nil && userTask.Completed {
		http.Error(w, "Task has already been completed ", http.StatusBadRequest)
		return
	}

	if err != nil {
		if err := h.store.CreateUserTask(userID, requestData.TaskID); err != nil {
			http.Error(w, "Failed to create user task record", http.StatusInternalServerError)
			return
		}
	}

	if err := h.store.CompleteUserTask(userID, requestData.TaskID); err != nil {
		http.Error(w, "Failed to complete task", http.StatusInternalServerError)
		return
	}

	if err := h.store.UpdateUserBalance(userID, task.Price); err != nil {
		http.Error(w, "Failed to update user balance", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{
		"message": "Task completed successfully",
	})
}

func (h *UserHandler) SetReferrer(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userID := vars["id"]

	var requestData struct {
		ReferrerID int `json:"referrer_id"`
	}

	err := json.NewDecoder(r.Body).Decode(&requestData)
	if err != nil {
		http.Error(w, "Error body req", http.StatusBadRequest)
		return
	}

	user, err := h.store.GetUserByID(userID)
	if err != nil {
		http.Error(w, "User not", http.StatusNotFound)
		return
	}
	if user.ReferrerID != nil {
		http.Error(w, "Referrer id has been already", http.StatusBadRequest)
		return
	}

	err = h.store.SetReferrerID(userID, requestData.ReferrerID)
	if err != nil {
		http.Error(w, "Failed to set referral ID", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{
		"message": "Referral ID successfully set",
	})
}
