package repositories

import (
	"database/sql"
	"dbconnection/models"
	"errors"
	"fmt"
	"log"
)

type UserRepository struct {
	DB *sql.DB
}

func (repo *UserRepository) GetUserByEmail(email string) (*models.User, error) {
	user := &models.User{}
	query := "SELECT id_usuario, nombre, ocupacion, email, contraseña, foto_perfil FROM Usuario WHERE email = ?"
	err := repo.DB.QueryRow(query, email).Scan(&user.ID, &user.Nombre, &user.Ocupacion, &user.Email, &user.Password, &user.FotoPerfil)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("usuario no encontrado")
		}
		return nil, fmt.Errorf("error al ejecutar la consulta: %v", err) // Mejorar el error para obtener detalles
	}
	return user, nil
}

func (repo *UserRepository) CreateUser(user *models.User) (*models.User, error) {
	result, err := repo.DB.Exec("INSERT INTO Usuario (nombre, ocupacion, contraseña, email) VALUES (?, ?, ?, ?)",
		user.Nombre, user.Ocupacion, user.Password, user.Email)

	if err != nil {
		log.Printf("Error al ejecutar INSERT en Usuario: %v", err)
		return nil, err
	}

	// Obtener el ID del usuario recién creado
	userID, err := result.LastInsertId()
	if err != nil {
		log.Printf("Error al obtener el ID del usuario recién creado: %v", err)
		return nil, err
	}

	// Recuperar el usuario recién creado para devolverlo completo
	newUser := &models.User{}
	query := "SELECT id_usuario, nombre, ocupacion, email FROM Usuario WHERE id_usuario = ?"
	err = repo.DB.QueryRow(query, userID).Scan(&newUser.ID, &newUser.Nombre, &newUser.Ocupacion, &newUser.Email)
	if err != nil {
		log.Printf("Error al ejecutar SELECT en Usuario para recuperar el nuevo usuario: %v", err)
		return nil, err
	}

	return newUser, nil
}

func (repo *UserRepository) GetUserByID(id int) (*models.User, error) {
	user := &models.User{}
	query := "SELECT id_usuario, nombre, ocupacion, email, foto_perfil FROM Usuario WHERE id_usuario = ?"
	err := repo.DB.QueryRow(query, id).Scan(&user.ID, &user.Nombre, &user.Ocupacion, &user.Email, &user.FotoPerfil)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (repo *UserRepository) UpdateUserProfile(id int, user *models.User) (*models.User, error) {
	// Preparamos la consulta SQL para actualizar el perfil del usuario
	// Foto_Perfil puede ser NULL si no se proporciona una imagen
	var query string
	var args []interface{}

	// Verificar si se proporcionó una foto de perfil
	if user.FotoPerfil != "" {
		// Si se proporciona una foto de perfil, actualizamos el campo con la URL
		query = `UPDATE Usuario SET nombre = ?, ocupacion = ?, email = ?, foto_perfil = ? WHERE id_usuario = ?`
		args = append(args, user.Nombre, user.Ocupacion, user.Email, user.FotoPerfil, id)
	} else {
		// Si no se proporciona foto de perfil, no la actualizamos (dejar NULL)
		query = `UPDATE Usuario SET nombre = ?, ocupacion = ?, email = ? WHERE id_usuario = ?`
		args = append(args, user.Nombre, user.Ocupacion, user.Email, id)
	}

	// Ejecutamos la consulta
	_, err := repo.DB.Exec(query, args...)
	if err != nil {
		log.Printf("Error al actualizar el perfil del usuario: %v", err)
		return nil, err
	}

	// Recuperar el usuario actualizado
	updatedUser := &models.User{}
	query = "SELECT id_usuario, nombre, ocupacion, email, foto_perfil FROM Usuario WHERE id_usuario = ?"
	err = repo.DB.QueryRow(query, id).Scan(&updatedUser.ID, &updatedUser.Nombre, &updatedUser.Ocupacion, &updatedUser.Email, &updatedUser.FotoPerfil)
	if err != nil {
		log.Printf("Error al recuperar el usuario actualizado: %v", err)
		return nil, err
	}

	return updatedUser, nil
}
