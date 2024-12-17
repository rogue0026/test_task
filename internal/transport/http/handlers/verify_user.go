package handlers

import (
	"encoding/json"
	"errors"
	"github.com/rogue0026/test_/internal/storage"
	"github.com/rogue0026/test_/internal/storage/users/postgres"
	"log/slog"
	"net/http"
	"time"
)

func VerifyUser(logger *slog.Logger, usersRepo postgres.UsersRepository) http.Handler {
	type Request struct {
		UserID           string `json:"user_id"`
		VerificationCode string `json:"verification_code"`
	}
	h := func(w http.ResponseWriter, r *http.Request) {
		in := Request{}
		err := json.NewDecoder(r.Body).Decode(&in)
		if err != nil {
			js, _ := json.Marshal(map[string]interface{}{"error": "invalid input data format"})
			w.WriteHeader(http.StatusBadRequest)
			_, _ = w.Write(js)
			return
		}
		usr, err := usersRepo.GetUnregisteredUserByID(r.Context(), in.UserID)
		if err != nil {
			if errors.Is(err, storage.ErrUserNotFound) {
				w.WriteHeader(http.StatusNotFound)
				return
			}
			logger.Error(err.Error())
			w.WriteHeader(http.StatusInternalServerError)
			js, _ := json.Marshal(map[string]interface{}{"error": "internal server error occurred"})
			_, _ = w.Write(js)
			return
		}
		if usr.VerificationCode != in.VerificationCode {
			w.WriteHeader(http.StatusUnauthorized)
			js, _ := json.Marshal(map[string]interface{}{"error": "invalid verification code"})
			_, _ = w.Write(js)
			return
		}
		t, _ := time.Parse("02.01.2006 15:04:05", usr.VerificationCodeExpires)
		if t.Before(time.Now()) {
			w.WriteHeader(http.StatusUnauthorized)
			js, _ := json.Marshal(map[string]interface{}{"error": "verification code is expired"})
			_, _ = w.Write(js)
			return
		}
		err = usersRepo.SaveRegisteredUser(r.Context(), usr)
		if err != nil {
			logger.Error(err.Error())
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
		js, _ := json.Marshal(map[string]interface{}{"message": "user successfully registered"})
		_, _ = w.Write(js)
	}
	return http.HandlerFunc(h)
}
