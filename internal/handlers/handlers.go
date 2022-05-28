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
	reports, err := m.DB.GetAllReports()
	if err != nil {
		helpers.ServerError(rw, err)
		return
	}

	data := make(map[string]interface{})
	data["reports"] = reports

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

	err := r.ParseForm()
	if err != nil {
		m.App.Session.Put(r.Context(), "error", "can't parse form!")
		http.Redirect(rw, r, "/", http.StatusTemporaryRedirect)
		return
	}

	report := models.Report{}

	report.ID, err = strconv.Atoi(r.Form.Get("id"))
	if err != nil {
		m.App.Session.Put(r.Context(), "error", "can't parse id!")
		http.Redirect(rw, r, "/", http.StatusTemporaryRedirect)
		return
	}

	form := forms.New(r.PostForm)

	if !form.Valid() {
		data := make(map[string]interface{})
		data["report"] = report

		render.Template(rw, r, "addNewReport.page.html", &models.TemplateData{
			Form: form,
			Data: data,
		})
		return
	}

	err = m.DB.AddReport(report)
	if err != nil {
		m.App.Session.Put(r.Context(), "error", "can't insert report to database!")
		http.Redirect(rw, r, "/", http.StatusTemporaryRedirect)
		return
	}

	m.App.Session.Put(r.Context(), "flash", "Inserted successfully")
	http.Redirect(rw, r, "/report/add-new-report", http.StatusSeeOther)
}
