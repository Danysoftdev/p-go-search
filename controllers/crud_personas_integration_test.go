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

	"github.com/danysoftdev/p-go-search/config"
	"github.com/danysoftdev/p-go-search/controllers"
	"github.com/danysoftdev/p-go-search/repositories"
	"github.com/danysoftdev/p-go-search/services"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

func TestEndpointsControllerIntegration(t *testing.T) {
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

	os.Setenv("MONGO_URI", "mongodb://"+endpoint)
	os.Setenv("MONGO_DB", "testdb")
	os.Setenv("COLLECTION_NAME", "personas_test")

	err = config.ConectarMongo()
	assert.NoError(t, err)
	defer config.CerrarMongo()

	repositories.SetCollection(config.Collection)
	services.SetPersonaRepository(repositories.RealPersonaRepository{})

	// Setup router
	router := mux.NewRouter()

	router.HandleFunc("/personas/{documento}", controllers.ObtenerPersonaPorDocumento).Methods("GET")

	// 3. Obtener por documento
	reqBuscar := httptest.NewRequest("GET", "/personas/999", nil)
	reqBuscar = mux.SetURLVars(reqBuscar, map[string]string{"documento": "999"})
	resBuscar := httptest.NewRecorder()
	router.ServeHTTP(resBuscar, reqBuscar)

	assert.Equal(t, http.StatusOK, resBuscar.Code)

}
