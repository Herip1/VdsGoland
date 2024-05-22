package storage

import (
	"app/entity"
	"testing"
)

func TestUserStorage(t *testing.T) {
	// Тестирование добавления пользователя
	user := entity.User{Name: "John Doe"} // Удалите поле Age, если оно не определено в entity.User
	createdUser := CreateUser(user)
	if createdUser.Uid == 0 {
		t.Errorf("Failed to create user: user ID is zero")
	}

	// Тестирование получения пользователя по ID
	retrievedUser, err := GetUserByID(createdUser.Uid)
	if err != nil || retrievedUser.Uid != createdUser.Uid {
		t.Errorf("Failed to retrieve user: %v", err)
	}

	// Тестирование обновления пользователя
	updatedUser := *retrievedUser
	updatedUser.Name = "Jane Doe"
	err = UpdateUser(updatedUser)
	if err != nil {
		t.Errorf("Failed to update user: %v", err)
	}

	// Проверка обновленных данных
	updatedUserCheck, _ := GetUserByID(updatedUser.Uid)
	if updatedUserCheck.Name != "Jane Doe" {
		t.Errorf("User update failed: expected name %s, got %s", "Jane Doe", updatedUserCheck.Name)
	}

	// Тестирование удаления пользователя
	err = DeleteUser(createdUser.Uid)
	if err != nil {
		t.Errorf("Failed to delete user: %v", err)
	}

	// Проверка наличия удаленного пользователя
	_, err = GetUserByID(createdUser.Uid)
	if err == nil {
		t.Errorf("User was not deleted")
	}

	// Тестирование получения всех пользователей
	users, err := GetAllUsers()
	if err != nil || len(users) != 0 {
		t.Errorf("Failed to get all users")
	}
}
