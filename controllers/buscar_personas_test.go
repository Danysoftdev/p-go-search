package controllers_test

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/danysoftdev/p-go-search/controllers"
	"github.com/danysoftdev/p-go-search/models"
	"github.com/danysoftdev/p-go-search/services"
	"github.com/danysoftdev/p-go-search/tests/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/gorilla/mux"

)

func TestObtenerPersonaPorDocumento_Success(t *testing.T) {
	mockRepo := new(mocks.MockPersonaRepo)
	services.SetPersonaRepository(mockRepo)

	personaEsperada := models.Persona{
		Documento: "12345",
		Nombre:    "Ana",
		Apellido:  "Díaz",
		Edad:      30,
		Correo:    "ana@correo.com",
		Telefono:  "3001234567",
		Direccion: "Calle Falsa",
	}

	mockRepo.On("ObtenerPersonaPorDocumento", "12345").Return(personaEsperada, nil)

	req := httptest.NewRequest("GET", "/personas/12345", nil)
	req = mux.SetURLVars(req, map[string]string{"documento": "12345"})

	rec := httptest.NewRecorder()
	controllers.ObtenerPersonaPorDocumento(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)

	var resp models.Persona
	json.NewDecoder(rec.Body).Decode(&resp)
	assert.Equal(t, personaEsperada.Nombre, resp.Nombre)
	mockRepo.AssertExpectations(t)
}

func TestObtenerPersonaPorDocumento_NotFound(t *testing.T) {
	mockRepo := new(mocks.MockPersonaRepo)
	services.SetPersonaRepository(mockRepo)

	mockRepo.On("ObtenerPersonaPorDocumento", "99999").Return(models.Persona{}, errors.New("persona no encontrada"))

	req := httptest.NewRequest("GET", "/personas/99999", nil)
	req = mux.SetURLVars(req, map[string]string{"documento": "99999"})

	rec := httptest.NewRecorder()
	controllers.ObtenerPersonaPorDocumento(rec, req)

	assert.Equal(t, http.StatusNotFound, rec.Code)
	assert.Contains(t, rec.Body.String(), "persona no encontrada")
	mockRepo.AssertExpectations(t)
}
