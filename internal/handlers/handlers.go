package handlers

import (
	"ParsissCrm/internal/config"
	"ParsissCrm/internal/driver"
	"ParsissCrm/internal/forms"
	"ParsissCrm/internal/models"
	"ParsissCrm/internal/render"
	"ParsissCrm/internal/repository"
	"ParsissCrm/internal/repository/dbrepo"
	"net/http"
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
