package core

import (
	"fmt"
)

type CustomError struct {
	Message string
}

func (e *CustomError) Error() string {
	return fmt.Sprintf("CustomError: %s", e.Message)
}

func NewErrBannerDoesNotExist() *CustomError {
	return &CustomError{Message: "Баннер для пользователя не найден"}
}

func NewErrInternalServerError() *CustomError {
	return &CustomError{Message: "Внутренняя ошибка сервера"}
}

func NewErrIncorrectInput() *CustomError {
	return &CustomError{Message: "Некорректные данные"}
}

func NewErrAccessDenied() *CustomError {
	return &CustomError{Message: "У пользователя нет доступа"}
}
