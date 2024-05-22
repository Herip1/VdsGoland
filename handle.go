package main

import (
	"app/config"
	"app/engine"
	"app/entity"
	"app/service"
	"encoding/json"
	"fmt"
	"gorm.io/gorm"
	"log"
	"net/http"
	"reflect"
	"strings"
	"time"

	"github.com/google/uuid"
)

var types map[string]bool
var hdl *service.Service
var apiMap map[string]map[string]reflect.Value

func init() {
	apiMap = make(map[string]map[string]reflect.Value)
	types = make(map[string]bool)
	types[".ico"] = true
	types[".html"] = true
	types[".js"] = true
	types[".svg"] = true
	types[".png"] = true
	hdl = &service.Service{}
	services := reflect.ValueOf(hdl)
	_struct := reflect.TypeOf(hdl)
	for methodNum := 0; methodNum < _struct.NumMethod(); methodNum++ {
		method := _struct.Method(methodNum)
		val, ok := config.Get().Api[method.Name]
		if !ok {
			continue
		}
		if _, ok := apiMap[val.Method]; !ok {
			apiMap[val.Method] = make(map[string]reflect.Value)
		}
		apiMap[val.Method][val.Url] = services.Method(methodNum)
	}
	http.HandleFunc("/logout", logoutHandler)
}
func handle(w http.ResponseWriter, r *http.Request) {
	ctx := engine.Context{
		Response: w,
		Request:  r,
	}
	url := r.URL
	path := url.Path[1:]
	pathArr := strings.Split(path, "/")

	// Проверка для маршрута /auth
	if pathArr[0] == "auth" && r.Method == http.MethodPost {
		AuthHandler(w, r)
		return
	}

	if pathArr[0] == "" || pathArr[0] == "index.html" {
		http.ServeFile(w, r, "./tpl/index.html")
		return
	}

	// Обработка запроса на получение таблицы пользователей
	if pathArr[0] == "getUsersTable" && r.Method == http.MethodGet {
		getUsersTableHandler(w, r)
		return
	}

	if pathArr[0] == "deleteUser" && r.Method == http.MethodPost {
		deleteUserHandler(w, r)
		return
	}

	if pathArr[0] == "addUser" && r.Method == http.MethodPost {
		addUserHandler(w, r)
		return
	}

	if pathArr[0] == "addVds" && r.Method == http.MethodPost {
		addVdsHandler(w, r)
		return
	}

	if pathArr[0] == "getVdsTable" && r.Method == http.MethodGet {
		getVdsTableHandler(w, r)
		return
	}

	if pathArr[0] == "deleteVds" && r.Method == http.MethodPost {
		deleteVdsHandler(w, r)
		return
	}

	// Обработка статических файлов
	last := pathArr[len(pathArr)-1]
	pos := strings.LastIndex(last, ".")
	if pos > 0 {
		str := last[pos:]
		if len(str) > 1 && types[str] {
			http.ServeFile(w, r, "./tpl/"+path)
			return
		}
	}

	// Проверка валидности токена перед отображением таблицы пользователей
	if isValidToken(r.Header.Get("Authorization")) {
		ctx.Response.Write([]byte("true"))
	} else {
		ctx.Response.Write([]byte("false"))
	}
}

func sendFile(url string, ctx engine.Context) {
	ctx.Response.Write([]byte{})
}

type SAuth struct {
	Login string `json:"login"`
	Pass  string `json:"pass"`
}

func AuthHandler(w http.ResponseWriter, r *http.Request) {
	// Получаем username и password из запроса
	d := json.NewDecoder(r.Body)
	var auth SAuth
	d.Decode(&auth)

	// Поиск пользователя в базе данных
	var user entity.User
	result := db.Where("login = ? AND pass = ?", auth.Login, auth.Pass).First(&user)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			http.Error(w, "User not found", http.StatusUnauthorized)
		} else {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		}
		return

	}

	// Генерируем токен
	tokenString := generateToken(user.Uid)

	// Отправляем ответ с токеном, uid, name и access
	response := struct {
		Token  string `json:"token"`
		Uid    uint64 `json:"uid"`
		Name   string `json:"name"`
		Access bool   `json:"access"`
	}{
		Token:  tokenString,
		Uid:    user.Uid,
		Name:   user.Name,
		Access: user.Access,
	}

	// Преобразуем структуру в JSON
	responseJSON, err := json.Marshal(response)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// Устанавливаем Content-Type и отправляем ответ
	w.Header().Set("Content-Type", "application/json")
	w.Write(responseJSON)
}

func generateToken(uid uint64) string {
	// Создаем UUID
	var token entity.Token
	db.Table("tokens").Where("uid = ?", uid).First(&token)
	token.Exp = time.Now().Add(time.Hour * 24)
	token.Uuid = uuid.NewString()
	token.Uid = uid
	db.Save(&token)
	return token.Uuid
}

// Ваша функция для проверки валидности токена в postgres
func isValidToken(token string) bool {
	var tokenEntity entity.Token
	result := db.Table("tokens").Where("uuid = ?", gorm.Expr("?::text", token)).First(&tokenEntity)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return false
		}
		return false
	}

	// Проверяем, не истек ли срок действия токена
	if time.Now().After(tokenEntity.Exp) {
		return false
	}

	return true
}

func getUsersTableHandler(w http.ResponseWriter, r *http.Request) {
	// Получение данных пользователей из базы данных
	users := []entity.User{}
	db.Find(&users)

	// Преобразование данных в формат JSON и отправка их
	responseJSON, err := json.Marshal(users)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(responseJSON)
}
func deleteUserHandler(w http.ResponseWriter, r *http.Request) {
	// Парсим тело запроса для получения userId
	var requestData struct {
		Uid uint64 `json:"uid"`
	}
	err := json.NewDecoder(r.Body).Decode(&requestData)
	if err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	// Удаление пользователя из базы данных
	result := db.Where("uid = ?", requestData.Uid).Delete(&entity.User{})

	if result.Error != nil {
		http.Error(w, "Error deleting user", http.StatusInternalServerError)
		return
	}
	log.Printf("Attempting to delete user with UID: %d", requestData.Uid)

	// Отправка ответа о успешном удалении
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{"success": true, "message": "User deleted successfully"})
}
func addUserHandler(w http.ResponseWriter, r *http.Request) {
	var user entity.User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}
	result := db.Create(&user)
	if result.Error != nil {
		http.Error(w, "Error adding user", http.StatusInternalServerError)
		return
	}

	// Отправка ответа о успешном добавлении
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{"success": true, "message": "User added successfully"})
}
func logoutHandler(w http.ResponseWriter, r *http.Request) {
	// Извлекаем токен из запроса
	tokenString := r.Header.Get("Authorization")

	// Делаем токен недействительным (пример: установка истечения на текущее время)
	var token entity.Token
	db.Model(&token).Where("uuid = ?", tokenString).Update("exp", time.Now())

	// Отправляем подтверждение об успешном выходе
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Вы успешно вышли из системы"))
}

func addVdsHandler(w http.ResponseWriter, r *http.Request) {
	var vds entity.Vds
	err := json.NewDecoder(r.Body).Decode(&vds)
	if err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}
	result := db.Create(&vds)
	if result.Error != nil {
		http.Error(w, "Error adding vds", http.StatusInternalServerError)
		return
	}

	// Отправка ответа о успешном добавлении
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{"success": true, "message": "Vds added successfully"})
}

func getVdsTableHandler(w http.ResponseWriter, r *http.Request) {
	// Получение данных пользователей из базы данных
	vds := []entity.Vds{}
	db.Find(&vds)

	// Преобразование данных в формат JSON и отправка их
	responseJSON, err := json.Marshal(vds)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		fmt.Println("Ошибочка тут")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(responseJSON)
}

func deleteVdsHandler(w http.ResponseWriter, r *http.Request) {
	// Парсим тело запроса для получения userId
	var requestData struct {
		Vid uint64 `json:"vid"`
	}
	err := json.NewDecoder(r.Body).Decode(&requestData)
	if err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	// Удаление пользователя из базы данных
	result := db.Where("vid = ?", requestData.Vid).Delete(&entity.Vds{})

	if result.Error != nil {
		http.Error(w, "Error deleting vds", http.StatusInternalServerError)
		return
	}
	log.Printf("Attempting to delete Vds with Vds: %d", requestData.Vid)

	// Отправка ответа о успешном удалении
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{"success": true, "message": "Vds deleted successfully"})
}
