package repositories

import (
	"database/sql"
	"dbconnection/models"
	"errors"
	"log"

	"github.com/go-sql-driver/mysql"
)

type PreferenceRepository struct {
	DB *sql.DB
}

func (repo *PreferenceRepository) GetPreferencesByUserID(idUsuario int, token string) (*models.Preference, error) {

	pref := &models.Preference{}
	query := "SELECT id_pref, linkedinlink, instagramlink, xlink FROM preferenciasusuarios WHERE id_usuario = ?"

	// Ejecutamos la consulta
	err := repo.DB.QueryRow(query, idUsuario).Scan(&pref.ID, &pref.Linkedin, &pref.Instagram, &pref.XLink)

	if err != nil {
		// Verificar si el error es el 1146 (Tabla no encontrada)
		if mysqlErr, ok := err.(*mysql.MySQLError); ok && mysqlErr.Number == 1146 {
			return nil, errors.New("La tabla 'preferenciasusuarios' no existe en la base de datos")
		}
		return nil, err
	}

	return pref, nil
}

// UpdatePreferences - Actualiza las preferencias de un usuario
func (repo *PreferenceRepository) UpdatePreferences(pref *models.Preference, token string) (*models.Preference, error) {
	// Verificar si el id_usuario existe
	var count int
	query := "SELECT COUNT(*) FROM usuario WHERE id_usuario = ?"
	err := repo.DB.QueryRow(query, pref.IdUsuario).Scan(&count)
	if err != nil {
		return nil, err
	}

	if count == 0 {
		return nil, errors.New("El usuario con el id proporcionado no existe")
	}

	// Actualizar las preferencias
	query = "UPDATE preferenciasusuarios SET linkedinlink = ?, instagramlink = ?, xlink = ? WHERE id_usuario = ?"
	_, err = repo.DB.Exec(query, pref.Linkedin, pref.Instagram, pref.XLink, pref.IdUsuario)
	if err != nil {
		log.Printf("Error al ejecutar UPDATE en preferenciasusuarios: %v", err)
		return nil, err
	}

	// Recuperar las preferencias actualizadas
	updatedPref := &models.Preference{}
	query = "SELECT id_pref, linkedinlink, instagramlink, xlink, id_usuario FROM preferenciasusuarios WHERE id_usuario = ?"
	err = repo.DB.QueryRow(query, pref.IdUsuario).Scan(&updatedPref.ID, &updatedPref.Linkedin, &updatedPref.Instagram, &updatedPref.XLink, &updatedPref.IdUsuario)
	if err != nil {
		log.Printf("Error al ejecutar SELECT en preferenciasusuarios para recuperar las preferencias actualizadas: %v", err)
		return nil, err
	}

	return updatedPref, nil
}

func (repo *PreferenceRepository) CreatePreferences(pref *models.Preference, token string) (*models.Preference, error) {
	// Verificar si el usuario existe antes de intentar insertar las preferencias
	query := "SELECT COUNT(*) FROM usuario WHERE id_usuario = ?"
	var count int
	err := repo.DB.QueryRow(query, pref.IdUsuario).Scan(&count)
	if err != nil {
		return nil, err
	}

	if count == 0 {
		return nil, errors.New("El usuario con el id proporcionado no existe")
	}

	// Inserta las preferencias en la base de datos
	query = "INSERT INTO preferenciasusuarios (id_usuario, linkedinlink, instagramlink, xlink) VALUES (?, ?, ?, ?)"
	result, err := repo.DB.Exec(query, pref.IdUsuario, pref.Linkedin, pref.Instagram, pref.XLink)
	if err != nil {
		log.Printf("Error al ejecutar INSERT en preferenciasusuarios: %v", err)
		return nil, err
	}

	// Obtener el ID del registro recién creado
	prefID, err := result.LastInsertId()
	if err != nil {
		log.Printf("Error al obtener el ID del nuevo registro de preferencias: %v", err)
		return nil, err
	}

	// Recuperar las preferencias recién creadas
	newPref := &models.Preference{}
	query = "SELECT id_pref, id_usuario, linkedinlink, instagramlink, xlink FROM preferenciasusuarios WHERE id_pref = ?"
	err = repo.DB.QueryRow(query, prefID).Scan(&newPref.ID, &newPref.IdUsuario, &newPref.Linkedin, &newPref.Instagram, &newPref.XLink)
	if err != nil {
		log.Printf("Error al recuperar las preferencias recién creadas: %v", err)
		return nil, err
	}

	return newPref, nil
}
