package controllers

import (
	"dbconnection/models"
	"dbconnection/services"
	"encoding/json"
	"log"
	"net/http"
	"strconv"
)

type PreferenceController struct {
	PreferenceService *services.PreferenceService
}

// CreatePreferences - Endpoint para crear nuevas preferencias de un usuario
func (controller *PreferenceController) CreatePreferences(w http.ResponseWriter, r *http.Request) {
	token := r.Header.Get("Authorization")
	if token == "" {
		http.Error(w, "Token de autenticación no proporcionado", http.StatusUnauthorized)
		return
	}

	var pref models.Preference
	if err := json.NewDecoder(r.Body).Decode(&pref); err != nil {
		log.Printf("Error al decodificar el cuerpo de la solicitud: %v", err)
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	createdPref, err := controller.PreferenceService.CreatePreferences(&pref, token)
	if err != nil {
		log.Printf("Error al crear las preferencias: %v", err)
		http.Error(w, "Error creating preferences", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(createdPref)
}

func (controller *PreferenceController) GetPreferences(w http.ResponseWriter, r *http.Request) {
	idUsuarioStr := r.URL.Query().Get("id")
	if idUsuarioStr == "" {
		http.Error(w, "Missing user ID", http.StatusBadRequest)
		return
	}

	idUsuario, err := strconv.Atoi(idUsuarioStr)
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	token := r.Header.Get("Authorization")

	preferences, err := controller.PreferenceService.GetPreferencesByUserID(idUsuario, token)
	if err != nil {
		// Si se recibe un error de la tabla inexistente, devolver 204 No Content
		if err.Error() == "La tabla 'preferencias_usuarios' no existe en la base de datos" {
			w.WriteHeader(http.StatusNoContent) // 204 No Content
			return
		}
		http.Error(w, "Error fetching preferences", http.StatusInternalServerError)
		return
	}

	if preferences == nil {
		// Si no se encontraron preferencias y no hay error, también devolver No Content
		w.WriteHeader(http.StatusNoContent) // 204 No Content
		return
	}

	// Si todo está bien, devolver las preferencias del usuario
	json.NewEncoder(w).Encode(preferences)
}

// UpdatePreferences - Endpoint para actualizar las preferencias del usuario
func (c *PreferenceController) UpdatePreferences(w http.ResponseWriter, r *http.Request) {
	token := r.Header.Get("Authorization")
	if token == "" {
		http.Error(w, "Token de autenticación no proporcionado", http.StatusUnauthorized)
		return
	}

	// Obtener el id_usuario y id_pref
	var pref models.Preference
	if err := json.NewDecoder(r.Body).Decode(&pref); err != nil {
		log.Printf("Error al decodificar el cuerpo de la solicitud: %v", err)
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	// Verificar si ya existe una preferencia para este usuario
	existingPref, err := c.PreferenceService.GetPreferencesByUserID(pref.IdUsuario, token)
	if err != nil {
		http.Error(w, "Error fetching preferences", http.StatusInternalServerError)
		return
	}

	if existingPref == nil {
		// Si no existen preferencias, crear nuevas preferencias
		newPref, err := c.PreferenceService.CreatePreferences(&pref, token)
		if err != nil {
			http.Error(w, "Error creating preferences", http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(newPref)
		return
	}

	// Si ya existen preferencias, actualizar las preferencias existentes
	updatedPref, err := c.PreferenceService.UpdatePreferences(&pref, token)
	if err != nil {
		log.Printf("Error al actualizar las preferencias: %v", err)
		http.Error(w, "Error updating preferences", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(updatedPref)
}
