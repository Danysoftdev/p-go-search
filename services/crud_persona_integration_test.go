//go:build integration
// +build integration

package services_test

import (
	"context"
	"testing"
	"time"

	"github.com/danysoftdev/p-go-search/config"
	"github.com/danysoftdev/p-go-search/models"
	"github.com/danysoftdev/p-go-search/repositories"
	"github.com/danysoftdev/p-go-search/services"

	"github.com/stretchr/testify/assert"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func TestBuscarPersonaPorDocumento(t *testing.T) {
	ctx := context.Background()

	req := testcontainers.ContainerRequest{
		Image:        "mongo:6.0",
		ExposedPorts: []string{"27017/tcp"},
		WaitingFor:   wait.ForListeningPort("27017/tcp").WithStartupTimeout(20 * time.Second),
	}
	mongoC, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
	assert.NoError(t, err)
	defer mongoC.Terminate(ctx)

	endpoint, err := mongoC.Endpoint(ctx, "")
	assert.NoError(t, err)

	t.Setenv("MONGO_URI", "mongodb://"+endpoint)
	t.Setenv("MONGO_DB", "testdb")
	t.Setenv("COLLECTION_NAME", "personas_test")

	err = config.ConectarMongo()
	assert.NoError(t, err)
	defer config.CerrarMongo()

	repositories.SetCollection(config.Collection)
	services.SetPersonaRepository(repositories.RealPersonaRepository{})

	// Limpiar colección
	_, err = config.Collection.DeleteMany(ctx, bson.M{})
	assert.NoError(t, err)

	// Insertar persona directamente
	persona := models.Persona{
		Documento: "12345",
		Nombre:    "Persona",
		Apellido:  "Prueba",
		Edad:      28,
		Correo:    "persona@prueba.com",
		Telefono:  "3001234567",
		Direccion: "Calle Falsa 123",
	}
	_, err = config.Collection.InsertOne(ctx, persona)
	assert.NoError(t, err)

	// Buscar desde el servicio
	encontrada, err := services.BuscarPersonaPorDocumento(persona.Documento)
	assert.NoError(t, err)
	assert.Equal(t, "Persona", encontrada.Nombre)
}
