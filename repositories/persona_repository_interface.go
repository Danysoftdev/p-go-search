package repositories

import "github.com/danysoftdev/p-go-search/models"

type PersonaRepository interface {
	ObtenerPersonaPorDocumento(documento string) (models.Persona, error)
}
