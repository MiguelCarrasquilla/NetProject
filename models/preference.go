package models

type Preference struct {
	ID        int    `json:"id_pref"`
	IdUsuario int    `json:"id_usuario"` // FK para relacionar con el usuario
	Linkedin  string `json:"linkedinlink"`
	Instagram string `json:"instagramlink"`
	XLink     string `json:"xlink"`
}
