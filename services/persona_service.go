package services

import (
	"errors"
	"strings"

	"github.com/danysoftdev/p-go-search/models"
	"github.com/danysoftdev/p-go-search/repositories"
	"go.mongodb.org/mongo-driver/mongo"
)

var Repo repositories.PersonaRepository

func SetPersonaRepository(r repositories.PersonaRepository) {
	Repo = r
}

func BuscarPersonaPorDocumento(doc string) (models.Persona, error) {
	if strings.TrimSpace(doc) == "" {
		return models.Persona{}, errors.New("el documento no puede estar vacío")
	}

	persona, err := Repo.ObtenerPersonaPorDocumento(doc)
	if err == mongo.ErrNoDocuments {
		return models.Persona{}, errors.New("persona no encontrada")
	}

	return persona, err
}
