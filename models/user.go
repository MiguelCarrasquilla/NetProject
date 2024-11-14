package models

type User struct {
	ID         int    `json:"id_usuario"`
	Nombre     string `json:"nombre"`
	Ocupacion  string `json:"ocupacion"`
	Password   string `json:"contraseña"`
	Email      string `json:"email"`
	FotoPerfil string `json:"foto_perfil"`
}
