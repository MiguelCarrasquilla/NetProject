package models

import "database/sql"

type Fair struct {
	ID          int            `json:"id_feria"`
	Titulo      string         `json:"titulo"`
	Descripcion string         `json:"descripcion"`
	FechaInicio string         `json:"fecha_inicio"` // Usa `time.Time` si prefieres manejar fechas
	IdUsuario   int            `json:"id_usuario"`   // FK para relacionar el usuario creador
	FotoFeria   sql.NullString `json:"foto_feria"`   // Nueva propiedad para la foto de la feria
}
