package services

import (
	"dbconnection/models"
	"dbconnection/repositories"
)

type PreferenceService struct {
	PreferenceRepo *repositories.PreferenceRepository
}

func (service *PreferenceService) UpdatePreferences(pref *models.Preference, token string) (*models.Preference, error) {
	return service.PreferenceRepo.UpdatePreferences(pref, token)
}
func (service *PreferenceService) GetPreferencesByUserID(idUsuario int, token string) (*models.Preference, error) {
	pref, err := service.PreferenceRepo.GetPreferencesByUserID(idUsuario, token)
	if err != nil {
		if err.Error() == "La tabla 'preferencias_usuarios' no existe en la base de datos" {
			return nil, nil // No content, sin datos para devolver
		}
		return nil, err
	}
	return pref, nil
}

// CreatePreferences - Llama al repositorio para crear nuevas preferencias
func (service *PreferenceService) CreatePreferences(pref *models.Preference, token string) (*models.Preference, error) {
	return service.PreferenceRepo.CreatePreferences(pref, token)
}
