package repositories

import (
	"database/sql"
	"dbconnection/models"
	"log"
)

type FairRepository struct {
	DB *sql.DB
}

// GetAllFairs obtiene todas las ferias de la base de datos
func (repo *FairRepository) GetAllFairs() ([]models.Fair, error) {
	var fairs []models.Fair
	query := "SELECT id_feria, titulo, descripcion, fecha_inicio, id_usuario, foto_feria FROM Feria"
	rows, err := repo.DB.Query(query)
	if err != nil {
		log.Printf("Error al obtener las ferias: %v", err)
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var fair models.Fair
		// Usar Scan para asignar el valor nulo a FotoFeria como sql.NullString
		if err := rows.Scan(&fair.ID, &fair.Titulo, &fair.Descripcion, &fair.FechaInicio, &fair.IdUsuario, &fair.FotoFeria); err != nil {
			log.Printf("Error al escanear la feria: %v", err)
			return nil, err
		}
		fairs = append(fairs, fair)
	}

	if err := rows.Err(); err != nil {
		log.Printf("Error al leer las filas: %v", err)
		return nil, err
	}

	return fairs, nil
}

// GetFairByID obtiene una feria por su ID
func (repo *FairRepository) GetFairByID(id int) (*models.Fair, error) {
	fair := &models.Fair{}
	query := "SELECT id_feria, titulo, descripcion, fecha_inicio, id_usuario, foto_feria FROM Feria WHERE id_feria = ?"
	err := repo.DB.QueryRow(query, id).Scan(&fair.ID, &fair.Titulo, &fair.Descripcion, &fair.FechaInicio, &fair.IdUsuario, &fair.FotoFeria)
	if err != nil {
		if err == sql.ErrNoRows {
			log.Printf("No se encontró una feria con ID %d", id)
		} else {
			log.Printf("Error al ejecutar SELECT en Feria: %v", err)
		}
		return nil, err
	}

	return fair, nil
}

// CreateFair inserta una nueva feria en la base de datos y la devuelve
func (repo *FairRepository) CreateFair(fair *models.Fair) (*models.Fair, error) {
	result, err := repo.DB.Exec("INSERT INTO Feria (titulo, descripcion, fecha_inicio, id_usuario, foto_feria) VALUES (?, ?, ?, ?, ?)",
		fair.Titulo, fair.Descripcion, fair.FechaInicio, fair.IdUsuario, fair.FotoFeria)
	if err != nil {
		log.Printf("Error al ejecutar INSERT en Feria: %v", err)
		return nil, err
	}

	// Obtener el ID de la feria recién creada
	fairID, err := result.LastInsertId()
	if err != nil {
		log.Printf("Error al obtener el ID de la feria recién creada: %v", err)
		return nil, err
	}

	// Recuperar la feria recién creada para devolverla completa
	newFair := &models.Fair{}
	query := "SELECT id_feria, titulo, descripcion, fecha_inicio, id_usuario, foto_feria FROM Feria WHERE id_feria = ?"
	err = repo.DB.QueryRow(query, fairID).Scan(&newFair.ID, &newFair.Titulo, &newFair.Descripcion, &newFair.FechaInicio, &newFair.IdUsuario, &newFair.FotoFeria)
	if err != nil {
		log.Printf("Error al ejecutar SELECT en Feria para recuperar la nueva feria: %v", err)
		return nil, err
	}

	return newFair, nil
}

func (repo *FairRepository) UpdateFair(id int, fair *models.Fair) (*models.Fair, error) {
	// Preparar la consulta de actualización
	query := `UPDATE Feria SET titulo = ?, descripcion = ?, fecha_inicio = ?, id_usuario = ?, foto_feria = ? WHERE id_feria = ?`

	// Aquí utilizamos .String si FotoFeria tiene valor, y "" si es nulo
	fotoFeriaValue := ""
	if fair.FotoFeria.Valid {
		fotoFeriaValue = fair.FotoFeria.String
	}

	_, err := repo.DB.Exec(query, fair.Titulo, fair.Descripcion, fair.FechaInicio, fair.IdUsuario, fotoFeriaValue, id)
	if err != nil {
		log.Printf("Error al ejecutar UPDATE en Feria: %v", err)
		return nil, err
	}

	// Recuperar la feria actualizada
	updatedFair := &models.Fair{}
	query = "SELECT id_feria, titulo, descripcion, fecha_inicio, id_usuario, foto_feria FROM Feria WHERE id_feria = ?"
	err = repo.DB.QueryRow(query, id).Scan(&updatedFair.ID, &updatedFair.Titulo, &updatedFair.Descripcion, &updatedFair.FechaInicio, &updatedFair.IdUsuario, &updatedFair.FotoFeria)
	if err != nil {
		log.Printf("Error al ejecutar SELECT en Feria para recuperar la feria actualizada: %v", err)
		return nil, err
	}

	return updatedFair, nil
}

// DeleteFair elimina una feria de la base de datos
func (repo *FairRepository) DeleteFair(id int) error {
	// Preparar la consulta para eliminar la feria por ID
	query := "DELETE FROM Feria WHERE id_feria = ?"
	_, err := repo.DB.Exec(query, id)
	if err != nil {
		log.Printf("Error al ejecutar DELETE en Feria: %v", err)
		return err
	}

	// Si se elimina correctamente, no devolvemos error
	return nil
}
