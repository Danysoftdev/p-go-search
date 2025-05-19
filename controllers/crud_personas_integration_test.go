//go:build integration
// +build integration

package controllers_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"

	"github.com/danysoftdev/p-go-search/config"
	"github.com/danysoftdev/p-go-search/controllers"
	"github.com/danysoftdev/p-go-search/models"
	"github.com/danysoftdev/p-go-search/repositories"
	"github.com/danysoftdev/p-go-search/services"
	"go.mongodb.org/mongo-driver/bson"
)

func TestObtenerPersonaPorDocumento_Integration(t *testing.T) {
	ctx := context.Background()

	// 1. Levantar contenedor de MongoDB
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

	os.Setenv("MONGO_URI", "mongodb://"+endpoint)
	os.Setenv("MONGO_DB", "testdb")
	os.Setenv("COLLECTION_NAME", "personas_test")

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
		Documento: "999",
		Nombre:    "Lucía",
		Apellido:  "Pérez",
		Edad:      30,
		Correo:    "lucia@example.com",
		Telefono:  "3111234567",
		Direccion: "Calle 123",
	}
	_, err = config.Collection.InsertOne(ctx, persona)
	assert.NoError(t, err)

	// 2. Configurar router y endpoint
	router := mux.NewRouter()
	router.HandleFunc("/personas/{documento}", controllers.ObtenerPersonaPorDocumento).Methods("GET")

	// 3. Simular request al endpoint
	reqBuscar := httptest.NewRequest("GET", "/personas/999", nil)
	reqBuscar = mux.SetURLVars(reqBuscar, map[string]string{"documento": "999"})
	resBuscar := httptest.NewRecorder()
	router.ServeHTTP(resBuscar, reqBuscar)

	// 4. Verificar respuesta esperada
	assert.Equal(t, http.StatusOK, resBuscar.Code)
	assert.Contains(t, resBuscar.Body.String(), "Lucía")
}
