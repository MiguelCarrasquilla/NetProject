package controllers

import (
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

type UserController struct {
	UserService *services.UserService
	Cloudinary  *cloudinary.Cloudinary
}

func (controller *UserController) Login(w http.ResponseWriter, r *http.Request) {
	// Parsear los datos del request
	var loginData struct {
		Email    string `json:"email"`
		Password string `json:"contraseña"`
	}

	// Log para ver los datos que se están recibiendo en la solicitud
	log.Printf("Received login request: email=%s", loginData.Email)

	if err := json.NewDecoder(r.Body).Decode(&loginData); err != nil {
		log.Printf("Error decoding request body: %v", err)
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Log para verificar que los datos fueron correctamente decodificados
	log.Printf("Decoded login data: email=%s, password=%s", loginData.Email, loginData.Password)

	// Llamar al servicio para realizar el login
	loginResponse, err := controller.UserService.Login(loginData.Email, loginData.Password)
	if err != nil {
		// Log el error que ocurre en el servicio de Login
		log.Printf("Error during login: %v", err)
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	// Log para verificar la respuesta del login
	log.Printf("Login successful, user ID: %d", loginResponse.User.ID)

	// Retornar el token y la información del usuario
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(loginResponse)
}

func (controller *UserController) CreateUser(w http.ResponseWriter, r *http.Request) {
	var user models.User
	err := json.NewDecoder(r.Body).Decode(&user) // Decodificar JSON
	if err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	// Crear el usuario con los datos del formulario
	createdUser, err := controller.UserService.RegisterUser(&user)
	if err != nil {
		http.Error(w, "Error creating user", http.StatusInternalServerError)
		return
	}

	// Devolver el usuario creado
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(createdUser)
}

func (c *UserController) GetUser(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Query().Get("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	user, err := c.UserService.GetUserProfile(id)
	if err != nil {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	json.NewEncoder(w).Encode(user)
}

func (controller *UserController) UpdateUserProfile(w http.ResponseWriter, r *http.Request) {
	// Parsear la solicitud como multipart/form-data
	err := r.ParseMultipartForm(10 << 20) // Limitar el tamaño del archivo a 10 MB
	if err != nil {
		http.Error(w, "Error parsing form", http.StatusBadRequest)
		return
	}

	// Obtener los datos del formulario
	user := &models.User{
		Nombre:    r.FormValue("nombre"),
		Ocupacion: r.FormValue("ocupacion"),
		Email:     r.FormValue("email"),
	}

	// Obtener la foto de perfil si está presente
	file, _, err := r.FormFile("foto_perfil")
	if err != nil && err != http.ErrMissingFile {
		http.Error(w, "Error uploading file", http.StatusInternalServerError)
		return
	}

	// Subir la foto a Cloudinary
	if file != nil {
		// Usar UploadParams de Cloudinary
		uploadParams := uploader.UploadParams{
			Folder:    "profile_pictures", // Puedes definir una carpeta en Cloudinary
			PublicID:  "profile_picture",  // Definir un public ID para la foto o dejar que lo haga Cloudinary automáticamente
			Overwrite: boolPtr(true),      // Si quieres que sobrescriba si ya existe
		}

		// Subir la imagen a Cloudinary
		uploadResult, err := controller.Cloudinary.Upload.Upload(r.Context(), file, uploadParams)
		if err != nil {
			http.Error(w, "Error uploading image to Cloudinary", http.StatusInternalServerError)
			return
		}

		// Asignar la URL de la foto de perfil desde Cloudinary
		user.FotoPerfil = uploadResult.URL
	}

	// Obtener el ID del usuario desde la URL
	idStr := mux.Vars(r)["id"]
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	// Actualizar los datos del usuario (incluyendo la URL de la foto de perfil)
	updatedUser, err := controller.UserService.UpdateUserProfile(id, user)
	if err != nil {
		http.Error(w, "Error updating user profile", http.StatusInternalServerError)
		return
	}

	// Retornar el usuario actualizado
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(updatedUser)
}

// Función auxiliar para convertir un valor bool en un puntero a bool
func boolPtr(b bool) *bool {
	return &b
}
