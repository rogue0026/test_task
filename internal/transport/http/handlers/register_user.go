package handlers

import (
	"crypto/rand"
	"encoding/json"
	"fmt"
	"github.com/rogue0026/test_/pkg/email"
	"os"
	"time"

	"github.com/rogue0026/test_/internal/models"
	"github.com/rogue0026/test_/internal/storage/users/postgres"
	"log/slog"
	"net/http"
)

const (
	smtpHost string = "smtp.yandex.ru"
	smtpPort string = "465"
)

func RegisterUser(logger *slog.Logger, userRepo postgres.UsersRepository) http.Handler {
	h := func(w http.ResponseWriter, r *http.Request) {
		user := models.User{}
		err := json.NewDecoder(r.Body).Decode(&user)
		if err != nil {
			logger.Debug("error", err.Error())
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		// проверяем что пользователя с указанным логином нет в базе
		_, err = userRepo.GetRegisteredUserByLogin(r.Context(), user.Login)
		if err == nil { // значит что пользователь в базе есть. Значит отвечаем что логин не подходит
			w.WriteHeader(http.StatusBadRequest)
			js, _ := json.Marshal(map[string]interface{}{"error": fmt.Sprintf("login %s already in user", user.Login)})
			_, _ = w.Write(js)
			return
		}
		// проверяем что пользователя с указанным email нет в базе
		_, err = userRepo.GetRegisteredUserByEmail(r.Context(), user.Email)
		if err == nil { // значит что пользователь в базе есть. Отвечаем, что email не подходит
			w.WriteHeader(http.StatusBadRequest)
			js, _ := json.Marshal(map[string]interface{}{"error": fmt.Sprintf("login %s already in user", user.Email)})
			_, _ = w.Write(js)
			return
		}
		// если пользователя с указанным логином и почтой нет в базе, то генерим код активации, отправляем его на указанную почту
		// добавляем данные во временную таблицу и ждем когда клиент пришлет код активации аккаунта
		code, err := rand.Prime(rand.Reader, 17)
		if err != nil {
			logger.Error(err.Error())
			w.WriteHeader(http.StatusInternalServerError)
			responseData, _ := json.Marshal(map[string]interface{}{"error": "error while generating activation code"})
			_, _ = w.Write(responseData)
			return
		}
		codeExpirationTime := time.Now().Add(time.Minute * 10)
		login := os.Getenv("SMTP_LOGIN")
		password := os.Getenv("SMTP_PASSWORD")
		box := email.NewMailBox("smtp.yandex.ru", 465, login, password)
		err = box.SendEmail("p.nozdra4ev@yandex.ru", user.Email, "Account verification", fmt.Sprintf("Verification code for your account activation: %s", code.String()))
		if err != nil {
			logger.Error(err.Error())
			w.WriteHeader(http.StatusInternalServerError)
			js, _ := json.Marshal(map[string]interface{}{"error": "error while sending activation code to email address"})
			_, _ = w.Write(js)
			return
		}
		user.VerificationCode = code.String()
		user.VerificationCodeExpires = codeExpirationTime.Format("02.01.2006 15:04:05")
		userID, err := userRepo.SaveUnregisteredUser(r.Context(), user)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			responseData, _ := json.Marshal(map[string]interface{}{"error": err.Error()})
			_, _ = w.Write(responseData)
			return
		}
		responseData, _ := json.Marshal(map[string]interface{}{"user_id": userID})
		w.WriteHeader(http.StatusCreated)
		_, _ = w.Write(responseData)
	}
	return http.HandlerFunc(h)
}
