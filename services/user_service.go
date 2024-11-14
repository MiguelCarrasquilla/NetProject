package services

import (
	"dbconnection/models"
	"dbconnection/repositories"
	"errors"
	"time"

	"github.com/dgrijalva/jwt-go"
)

type UserService struct {
	UserRepo *repositories.UserRepository
}

var jwtKey = []byte("ASDFGHLJKQKWJDAKSDQWPWEASDL")

// Estructura de respuesta para el login (token + datos del usuario)
type LoginResponse struct {
	Token string       `json:"token"`
	User  *models.User `json:"user"`
}

func (service *UserService) Login(email, password string) (*LoginResponse, error) {
	// Buscar el usuario por email
	user, err := service.UserRepo.GetUserByEmail(email)
	if err != nil {
		if err.Error() == "usuario no encontrado" {
			return nil, errors.New("No se encontró un usuario con el email proporcionado.")
		}
		return nil, errors.New("Error al buscar usuario en la base de datos.")
	}

	// Verificar si la contraseña coincide usando una comparación de texto plano
	if user.Password != password {
		// Si la comparación falla, significa que la contraseña es incorrecta
		return nil, errors.New("La contraseña proporcionada es incorrecta.")
	}

	// Crear el token JWT
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"userId": user.ID,
		"exp":    time.Now().Add(time.Hour * 24).Unix(), // El token expira en 24 horas
	})

	// Firmar el token con la clave secreta
	tokenString, err := token.SignedString(jwtKey)
	if err != nil {
		return nil, errors.New("Error al generar el token de autenticación.")
	}

	// Limpiar la contraseña antes de devolver los datos
	user.Password = "" // Dejar la contraseña vacía para que no se incluya en la respuesta

	// Retornar el token y los datos del usuario
	return &LoginResponse{
		Token: tokenString,
		User:  user,
	}, nil
}

func (service *UserService) RegisterUser(user *models.User) (*models.User, error) {
	return service.UserRepo.CreateUser(user)
}

func (service *UserService) GetUserProfile(id int) (*models.User, error) {
	return service.UserRepo.GetUserByID(id)
}

// Función para actualizar el perfil del usuario
func (service *UserService) UpdateUserProfile(id int, user *models.User) (*models.User, error) {
	// Llamamos al repositorio para actualizar el perfil
	updatedUser, err := service.UserRepo.UpdateUserProfile(id, user)
	if err != nil {
		return nil, err
	}

	return updatedUser, nil
}
