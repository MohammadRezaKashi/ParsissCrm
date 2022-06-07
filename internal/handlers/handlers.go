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

	data := make(map[string]interface{})

	surgeryDay, surgerytime, surgeryarea, surgeryresult, hospitaltype, paymentstatus, headfixtype := GetAllSelectOptions()

	data["patient"] = models.PersonalInformation{}
	data["surgeryinfo"] = models.SurgeriesInformation{}
	data["surgeryday"] = surgeryDay
	data["surgerytime"] = surgerytime
	data["surgeryarea"] = surgeryarea
	data["surgeryresult"] = surgeryresult
	data["hospitaltype"] = hospitaltype
	data["paymentstatus"] = paymentstatus
	data["headfixtype"] = headfixtype
	data["baseurl"] = "http://localhost:8080"

	render.Template(rw, r, "addNewReport.page.html", &models.TemplateData{
		Form: forms.New(nil),
		Data: data,
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

	surgeryDay, surgerytime, surgeryarea, surgeryresult, hospitaltype, paymentstatus, headfixtype := GetAllSelectOptions()

	data["patient"] = patient
	data["surgeryinfo"] = surgeryInfo[0]
	data["surgeryday"] = surgeryDay
	data["surgerytime"] = surgerytime
	data["surgeryarea"] = surgeryarea
	data["surgeryresult"] = surgeryresult
	data["hospitaltype"] = hospitaltype
	data["paymentstatus"] = paymentstatus
	data["headfixtype"] = headfixtype
	data["baseurl"] = "http://localhost:8080"

	render.Template(rw, r, "addNewReport.page.html", &models.TemplateData{
		Form: forms.New(nil),
		Data: data,
	})
}

func GetAllSelectOptions() ([]models.Option, []models.Option, []models.Option, []models.Option, []models.Option, []models.Option, []models.Option) {
	surgeryDay := []models.Option{
		{Value: "1", Text: "Saturday", Selected: ""},
		{Value: "2", Text: "Sunday", Selected: ""},
		{Value: "3", Text: "Monday", Selected: ""},
		{Value: "4", Text: "Tuesday", Selected: ""},
		{Value: "5", Text: "Wednesday", Selected: ""},
		{Value: "6", Text: "Thursday", Selected: ""},
		{Value: "7", Text: "Friday", Selected: ""},
	}

	surgerytime := []models.Option{
		{Value: "1", Text: "Morning", Selected: ""},
		{Value: "2", Text: "Afternoon", Selected: ""},
		{Value: "3", Text: "Evening", Selected: ""},
	}

	surgeryarea := []models.Option{
		{Value: "1", Text: "Neurosurgery", Selected: ""},
		{Value: "2", Text: "ENT", Selected: ""},
		{Value: "3", Text: "CMF", Selected: ""},
		{Value: "4", Text: "Spine", Selected: ""},
		{Value: "5", Text: "Orthopedics", Selected: ""},
	}

	surgeryresult := []models.Option{
		{Value: "1", Text: "Success", Selected: ""},
		{Value: "2", Text: "Fail", Selected: ""},
		{Value: "3", Text: "Canceled", Selected: ""},
	}

	hospitaltype := []models.Option{
		{Value: "1", Text: "Private", Selected: ""},
		{Value: "2", Text: "Govermental", Selected: ""},
		{Value: "3", Text: "Other", Selected: ""},
	}

	paymentstatus := []models.Option{
		{Value: "1", Text: "Paid", Selected: ""},
		{Value: "2", Text: "Unpaid", Selected: ""},
		{Value: "3", Text: "Free", Selected: ""},
		{Value: "4", Text: "Health Plan", Selected: ""},
		{Value: "5", Text: "Paid By Hospital", Selected: ""},
	}

	headfixtype := []models.Option{
		{Value: "1", Text: "Headband", Selected: ""},
		{Value: "1", Text: "Mayfield", Selected: ""},
		{Value: "1", Text: "Other", Selected: ""},
	}
	return surgeryDay, surgerytime, surgeryarea, surgeryresult, hospitaltype, paymentstatus, headfixtype
}
