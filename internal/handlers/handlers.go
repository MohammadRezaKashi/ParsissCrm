package handlers

import (
	"ParsissCrm/internal/config"
	"ParsissCrm/internal/driver"
	"ParsissCrm/internal/forms"
	"ParsissCrm/internal/helpers"
	"ParsissCrm/internal/models"
	"ParsissCrm/internal/render"
	"ParsissCrm/internal/repository"
	"ParsissCrm/internal/repository/dbrepo"
	"net/http"
	"strconv"

	"github.com/go-chi/chi"
)

var Repo *Repository

type Repository struct {
	App *config.AppConfig
	DB  repository.DatabaseRepo
}

func NewRepo(a *config.AppConfig, db *driver.DB) *Repository {
	return &Repository{
		App: a,
		DB:  dbrepo.NewPostgresRepo(db.SQL, a),
	}
}

func NewHandlers(r *Repository) {
	Repo = r
}

func (m *Repository) Home(rw http.ResponseWriter, r *http.Request) {
	render.Template(rw, r, "home.page.html", &models.TemplateData{})
}

func (m *Repository) Report(rw http.ResponseWriter, r *http.Request) {
	patients, err := m.DB.GetAllPatients()
	if err != nil {
		helpers.ServerError(rw, err)
		return
	}
	data := make(map[string]interface{})
	data["patients"] = patients
	render.Template(rw, r, "report.page.html", &models.TemplateData{
		Data: data,
	})
}

func (m *Repository) About(rw http.ResponseWriter, r *http.Request) {
	render.Template(rw, r, "about.page.html", &models.TemplateData{})
}

func (m *Repository) Contact(rw http.ResponseWriter, r *http.Request) {
	render.Template(rw, r, "contact.page.html", &models.TemplateData{})
}

func (m *Repository) AddNewReport(rw http.ResponseWriter, r *http.Request) {
	render.Template(rw, r, "addNewReport.page.html", &models.TemplateData{
		Form: forms.New(nil),
	})
}

func (m *Repository) PostAddNewReport(rw http.ResponseWriter, r *http.Request) {

}

func (m *Repository) ShowDetail(rw http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		m.App.Session.Put(r.Context(), "error", "invalid data!")
		http.Redirect(rw, r, "/report", http.StatusTemporaryRedirect)
	}

	patient, err := m.DB.GetPatientByID(id)
	if err != nil {
		helpers.ServerError(rw, err)
		return
	}

	surgeryInfo, err := m.DB.GetSurgicalInformationByPatientID(id)
	if err != nil {
		helpers.ServerError(rw, err)
		return
	}

	data := make(map[string]interface{})

	data["patient"] = patient
	data["surgeryinfo"] = surgeryInfo
	data["baseurl"] = "http://localhost:8080"

	render.Template(rw, r, "addNewReport.page.html", &models.TemplateData{
		Form: forms.New(nil),
		Data: data,
	})
}
