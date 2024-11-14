package controllers

import (
	"database/sql"
	"dbconnection/models"
	"dbconnection/services"
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/cloudinary/cloudinary-go/v2"
	"github.com/cloudinary/cloudinary-go/v2/api/uploader"
	"github.com/gorilla/mux"
)

type FairController struct {
	FairService *services.FairService
	Cloudinary  *cloudinary.Cloudinary
}

// DeleteFair - Endpoint para eliminar una feria por ID
func (c *FairController) DeleteFair(w http.ResponseWriter, r *http.Request) {
	// Obtener el ID de la feria desde la URL
	idStr := mux.Vars(r)["id"]
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid fair ID", http.StatusBadRequest)
		return
	}

	// Llamar al servicio para eliminar la feria
	err = c.FairService.DeleteFair(id)
	if err != nil {
		log.Printf("Error al eliminar la feria: %v", err)
		http.Error(w, "Error deleting fair", http.StatusInternalServerError)
		return
	}

	// Retornar una respuesta exitosa
	w.WriteHeader(http.StatusNoContent) // 204 No Content
}

func (c *FairController) UpdateFair(w http.ResponseWriter, r *http.Request) {
	// Parsear la solicitud como multipart/form-data
	err := r.ParseMultipartForm(10 << 20) // Limitar el tamaño del archivo a 10 MB
	if err != nil {
		http.Error(w, "Error parsing form", http.StatusBadRequest)
		return
	}

	// Obtener los datos del formulario
	fair := &models.Fair{
		Titulo:      r.FormValue("titulo"),
		Descripcion: r.FormValue("descripcion"),
		FechaInicio: r.FormValue("fecha_inicio"),
	}

	// Convertir el id_usuario (string) a int
	idUsuarioStr := r.FormValue("id_usuario")
	idUsuario, err := strconv.Atoi(idUsuarioStr)
	if err != nil {
		log.Printf("Error al convertir el id_usuario: %v", err)
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}
	fair.IdUsuario = idUsuario

	// Obtener el ID de la feria de la URL
	idStr := mux.Vars(r)["id"]
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid fair ID", http.StatusBadRequest)
		return
	}

	// Obtener la foto de la feria si está presente
	file, _, err := r.FormFile("foto_feria")
	if err != nil && err != http.ErrMissingFile {
		http.Error(w, "Error uploading file", http.StatusInternalServerError)
		return
	}

	// Subir la foto a Cloudinary
	if file != nil {
		// Crear un nombre único para la foto usando el ID de la feria
		publicID := "fair_picture_" + strconv.Itoa(id)

		// Usar UploadParams de Cloudinary con nombre dinámico
		uploadParams := uploader.UploadParams{
			Folder:    "fair_pictures", // Puedes definir una carpeta en Cloudinary
			PublicID:  publicID,        // Nombre único para la foto de la feria
			Overwrite: boolPtr(true),   // Sobrescribir si ya existe una imagen con el mismo nombre
		}

		// Subir la imagen a Cloudinary
		uploadResult, err := c.Cloudinary.Upload.Upload(r.Context(), file, uploadParams)
		if err != nil {
			http.Error(w, "Error uploading image to Cloudinary", http.StatusInternalServerError)
			return
		}

		// Asignar la URL de la foto de la feria desde Cloudinary
		// Convertir la URL de la foto a sql.NullString
		fair.FotoFeria = sql.NullString{
			String: uploadResult.URL,
			Valid:  true,
		}
	}

	// Llamar al servicio para actualizar la feria
	updatedFair, err := c.FairService.UpdateFair(id, fair) // Llamamos al servicio de actualización
	if err != nil {
		log.Printf("Error al actualizar la feria: %v", err)
		http.Error(w, "Error updating fair", http.StatusInternalServerError)
		return
	}

	// Retornar la feria actualizada
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(updatedFair)
}

func (c *FairController) CreateFair(w http.ResponseWriter, r *http.Request) {
	// Parsear la solicitud como multipart/form-data
	err := r.ParseMultipartForm(10 << 20) // Limitar el tamaño del archivo a 10 MB
	if err != nil {
		http.Error(w, "Error parsing form", http.StatusBadRequest)
		return
	}

	// Obtener los datos del formulario
	fair := &models.Fair{
		Titulo:      r.FormValue("titulo"),
		Descripcion: r.FormValue("descripcion"),
		FechaInicio: r.FormValue("fecha_inicio"),
	}

	// Convertir el id_usuario (string) a int
	idUsuarioStr := r.FormValue("id_usuario")
	idUsuario, err := strconv.Atoi(idUsuarioStr)
	if err != nil {
		log.Printf("Error al convertir el id_usuario: %v", err)
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}
	fair.IdUsuario = idUsuario

	// Obtener la foto de la feria si está presente
	file, _, err := r.FormFile("foto_feria")
	if err != nil && err != http.ErrMissingFile {
		http.Error(w, "Error uploading file", http.StatusInternalServerError)
		return
	}

	// Subir la foto a Cloudinary
	if file != nil {
		// Crear un nombre único para la foto usando el ID de la feria
		publicID := "fair_picture_" + strconv.Itoa(fair.IdUsuario)

		// Usar UploadParams de Cloudinary con nombre dinámico
		uploadParams := uploader.UploadParams{
			Folder:    "fair_pictures", // Puedes definir una carpeta en Cloudinary
			PublicID:  publicID,        // Nombre único para la foto de la feria
			Overwrite: boolPtr(true),   // Sobrescribir si ya existe una imagen con el mismo nombre
		}

		// Subir la imagen a Cloudinary
		uploadResult, err := c.Cloudinary.Upload.Upload(r.Context(), file, uploadParams)
		if err != nil {
			http.Error(w, "Error uploading image to Cloudinary", http.StatusInternalServerError)
			return
		}

		// Asignar la URL de la foto de la feria desde Cloudinary
		// Convertir la URL de la foto a sql.NullString
		fair.FotoFeria = sql.NullString{
			String: uploadResult.URL,
			Valid:  true,
		}
	}

	// Crear la feria con los datos del formulario y la URL de la foto de la feria
	createdFair, err := c.FairService.CreateFair(fair, r) // Pasar el objeto fair y el request r
	if err != nil {
		log.Printf("Error al crear la feria: %v", err)
		http.Error(w, "Error creating fair", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(createdFair)
}

func (c *FairController) GetFair(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Query().Get("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		log.Printf("Error al convertir el ID de la feria: %v", err)
		http.Error(w, "Invalid fair ID", http.StatusBadRequest)
		return
	}

	fair, err := c.FairService.GetFairDetails(id)
	if err != nil {
		log.Printf("Error al obtener la feria con ID %d: %v", id, err)
		http.Error(w, "Fair not found", http.StatusNotFound)
		return
	}

	// Convertir la foto_feria a una URL válida si existe
	if fair.FotoFeria.Valid {
		fair.FotoFeria.String = fair.FotoFeria.String // Foto_faire no nula, mostrar URL
	} else {
		fair.FotoFeria.String = "" // Si la foto es nula, asignamos vacío
	}

	json.NewEncoder(w).Encode(fair)
}

func (c *FairController) GetAllFairs(w http.ResponseWriter, r *http.Request) {
	fairs, err := c.FairService.GetAllFairs()
	if err != nil {
		http.Error(w, "Error al obtener las ferias: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Asignar el valor de FotoFeria en el resultado
	for i, fair := range fairs {
		if fair.FotoFeria.Valid {
			fairs[i].FotoFeria.String = fair.FotoFeria.String // Mostrar la URL de la foto
		} else {
			fairs[i].FotoFeria.String = "" // Si es nula, mostrar vacío
		}
	}

	json.NewEncoder(w).Encode(fairs)
}
