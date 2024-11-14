package services

import (
	"dbconnection/models"
	"dbconnection/repositories"
	"log"
	"net/http"

	"github.com/cloudinary/cloudinary-go"
)

type FairService struct {
	FairRepo   *repositories.FairRepository
	Cloudinary *cloudinary.Cloudinary
}

// DeleteFair elimina una feria usando el repositorio
func (service *FairService) DeleteFair(id int) error {
	// Llamar al repositorio para eliminar la feria de la base de datos
	err := service.FairRepo.DeleteFair(id)
	if err != nil {
		log.Printf("Error al eliminar la feria en el repositorio: %v", err)
		return err
	}

	// Si no hay error, la eliminación fue exitosa
	return nil
}

func (service *FairService) UpdateFair(id int, fair *models.Fair) (*models.Fair, error) {
	// Llamar al repositorio para actualizar la feria en la base de datos
	updatedFair, err := service.FairRepo.UpdateFair(id, fair)
	if err != nil {
		log.Printf("Error al actualizar la feria en el repositorio: %v", err)
		return nil, err
	}

	// Retornar la feria actualizada
	return updatedFair, nil
}

// CreateFair - Servicio para crear una feria y devolver el objeto creado
func (service *FairService) CreateFair(fair *models.Fair, r *http.Request) (*models.Fair, error) {
	// Llamar al repositorio para crear la feria en la base de datos
	createdFair, err := service.FairRepo.CreateFair(fair)
	if err != nil {
		log.Printf("Error al crear la feria en el repositorio: %v", err)
		return nil, err
	}

	// Retornar la feria creada
	return createdFair, nil
}

// GetAllFairs obtiene todas las ferias del repositorio
func (service *FairService) GetAllFairs() ([]models.Fair, error) {
	return service.FairRepo.GetAllFairs()
}

func (service *FairService) GetFairDetails(id int) (*models.Fair, error) {
	return service.FairRepo.GetFairByID(id)
}
