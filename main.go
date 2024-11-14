package main

import (
	"dbconnection/config"
	"dbconnection/controllers"
	"dbconnection/db"
	"dbconnection/repositories"
	"dbconnection/services"
	"log"
	"net/http"

	"github.com/cloudinary/cloudinary-go/v2" // Asegúrate de que esta importación esté presente
	"github.com/gorilla/mux"
	"github.com/rs/cors"
)

func main() {
	// Cargar configuración y conectar a la base de datos
	cfg := config.LoadConfig()
	database, err := db.Connect(cfg)
	if err != nil {
		log.Fatalf("Error de conexión a la base de datos: %v", err)
	}
	defer database.Close()

	// Inicializar repositorios, servicios y controladores
	userRepo := &repositories.UserRepository{DB: database}
	userService := &services.UserService{UserRepo: userRepo}
	userController := &controllers.UserController{UserService: userService}

	fairRepo := &repositories.FairRepository{DB: database}
	fairService := &services.FairService{FairRepo: fairRepo}
	fairController := &controllers.FairController{FairService: fairService}

	preferenceRepo := &repositories.PreferenceRepository{DB: database}
	preferenceService := &services.PreferenceService{PreferenceRepo: preferenceRepo}
	preferenceController := &controllers.PreferenceController{PreferenceService: preferenceService}

	// Configurar Cloudinary
	cld, err := cloudinary.NewFromParams("drlf5ytmk", "241212669924127", "kJQMb-K02kyMVnSjGSWPMd_Vpgs")
	if err != nil {
		log.Fatalf("Error de configuración de Cloudinary: %v", err)
	}
	log.Println("Cloudinary configurado correctamente")

	// Pasar la instancia de Cloudinary al controlador
	userController.Cloudinary = cld
	fairController.Cloudinary = cld

	// Configurar las rutas de la API
	mux := mux.NewRouter()
	mux.HandleFunc("/api/login", userController.Login)
	mux.HandleFunc("/api/users", userController.CreateUser)
	mux.HandleFunc("/api/users/get", userController.GetUser)
	mux.HandleFunc("/api/users/update/{id}", userController.UpdateUserProfile)
	mux.HandleFunc("/api/fairs", fairController.CreateFair)
	mux.HandleFunc("/api/fairs/get", fairController.GetFair)
	mux.HandleFunc("/api/fairs/getAll", fairController.GetAllFairs)
	mux.HandleFunc("/api/fairs/update/{id}", fairController.UpdateFair)
	mux.HandleFunc("/api/fairs/delete/{id}", fairController.DeleteFair)
	mux.HandleFunc("/api/preferences", preferenceController.GetPreferences)
	mux.HandleFunc("/api/preferences/update", preferenceController.UpdatePreferences)
	mux.HandleFunc("/api/preferences/create", preferenceController.CreatePreferences)

	// Configurar el middleware CORS
	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"https://net-project-nextjs.vercel.app"}, // Permitir solicitudes desde localhost:3000 (frontend)
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Content-Type", "Authorization"},
		AllowCredentials: true,
	})

	// Envolver el servidor mux con CORS
	handler := c.Handler(mux)

	// Iniciar el servidor
	log.Println("Servidor escuchando en http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", handler))
}
