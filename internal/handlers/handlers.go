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
	"github.com/jinzhu/copier"
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

	surgeryDay, surgerytime, surgeryarea, surgeryresult, hospitaltype, paymentstatus, headfixtype, imagevalidity := GetAllSelectOptions()

	data["patient"] = models.PersonalInformation{}
	data["surgeryinfo"] = models.SurgeriesInformation{}
	data["surgeryday"] = surgeryDay
	data["surgerytime"] = surgerytime
	data["surgeryarea"] = surgeryarea
	data["surgeryresult"] = surgeryresult
	data["hospitaltype"] = hospitaltype
	data["paymentstatus"] = paymentstatus
	data["headfixtype"] = headfixtype
	data["imagevalidity"] = imagevalidity
	data["baseurl"] = "http://localhost:8080"

	render.Template(rw, r, "addNewReport.page.html", &models.TemplateData{
		Form: forms.New(nil),
		Data: data,
	})
}

func (m *Repository) PostAddNewReport(rw http.ResponseWriter, r *http.Request) {
	var patient models.PersonalInformation

	patient.Name = r.Form.Get("name")
	patient.Family = r.Form.Get("family")
	age, err := strconv.Atoi(r.Form.Get("age"))
	if err != nil {
		m.App.Session.Put(r.Context(), "error", "can't parse age!")
		http.Redirect(rw, r, "/report", http.StatusTemporaryRedirect)
		return
	}
	patient.Age = age
	patient.Address = r.Form.Get("address")
	patient.Email = r.Form.Get("email")
	patient.NationalID = r.Form.Get("national_id")
	patient.PhoneNumber = r.Form.Get("phone_number")
	patient.PlaceOfBirth = r.Form.Get("place_of_birth")

	id, err := m.DB.AddPersonalInformation(patient)
	if err != nil {
		m.App.Session.Put(r.Context(), "error", "can't add personal information!")
		http.Redirect(rw, r, "/report", http.StatusTemporaryRedirect)
		return
	}

	var surgeryinfo models.SurgeriesInformation

	err = m.DB.AddSurgeriesInformation(surgeryinfo, id)
	if err != nil {
		m.App.Session.Put(r.Context(), "error", "can't add surgeries information!")
		http.Redirect(rw, r, "/report", http.StatusTemporaryRedirect)
		return
	}

	http.Redirect(rw, r, "/report", http.StatusTemporaryRedirect)
}

func (m *Repository) PostUpdateReport(rw http.ResponseWriter, r *http.Request) {
	patient, ok := m.App.Session.Get(r.Context(), "patient").(models.PersonalInformation)
	if !ok {
		m.App.Session.Put(r.Context(), "error", "can't get personal information data!")
		http.Redirect(rw, r, "/report", http.StatusTemporaryRedirect)
	}

	err := r.ParseForm()
	if err != nil {
		m.App.Session.Put(r.Context(), "error", "can't parse from!")
		http.Redirect(rw, r, "/report", http.StatusTemporaryRedirect)
		return
	}

	patient.Name = r.Form.Get("name")
	patient.Family = r.Form.Get("family")
	age, err := strconv.Atoi(r.Form.Get("age"))
	if err != nil {
		m.App.Session.Put(r.Context(), "error", "can't parse age!")
		http.Redirect(rw, r, "/report", http.StatusTemporaryRedirect)
		return
	}
	patient.Age = age
	patient.Address = r.Form.Get("address")
	patient.Email = r.Form.Get("email")
	patient.NationalID = r.Form.Get("national_id")
	patient.PhoneNumber = r.Form.Get("phone_number")
	patient.PlaceOfBirth = r.Form.Get("place_of_birth")

	err = m.DB.PutPersonalInformation(patient)
	if err != nil {
		m.App.Session.Put(r.Context(), "error", "can't update personal information!")
		http.Redirect(rw, r, "/report", http.StatusTemporaryRedirect)
		return
	}

	surgery, ok := m.App.Session.Get(r.Context(), "surgeryinfo").(models.SurgeriesInformation)
	if !ok {
		m.App.Session.Put(r.Context(), "error", "can't get surgery information data!")
		http.Redirect(rw, r, "/report", http.StatusTemporaryRedirect)
	}

	surgery.FileNumber = r.Form.Get("file_number")
	surgery.DateOfHospitalAdmission = driver.ConvertStringToDate(r.Form.Get("date_of_hospital_admission"))
	surgery.SurgeryDate = driver.ConvertStringToDate(r.Form.Get("surgery_date"))
	surgery.SurgeryDay, _ = strconv.Atoi(r.Form.Get("surgery_day"))
	surgery.SurgeryTime, _ = strconv.Atoi(r.Form.Get("surgery_time"))
	surgery.SurgeryType = r.Form.Get("surgery_type")
	surgery.SurgeryArea, _ = strconv.Atoi(r.Form.Get("surgery_area"))
	surgery.SurgeryDescription = r.Form.Get("surgery_description")
	surgery.SurgeryResult, _ = strconv.Atoi(r.Form.Get("surgery_result"))
	surgery.SurgeonFirst = r.Form.Get("surgeon_first")
	surgery.SurgeonSecond = r.Form.Get("surgeon_second")
	surgery.Resident = r.Form.Get("resident")
	surgery.Hospital = r.Form.Get("hospital")
	surgery.HospitalType, _ = strconv.Atoi(r.Form.Get("hospital_type"))
	surgery.HospitalAddress = r.Form.Get("hospital_address")
	surgery.CT, _ = strconv.Atoi(r.Form.Get("ct"))
	surgery.MR, _ = strconv.Atoi(r.Form.Get("mr"))
	surgery.FMRI, _ = strconv.Atoi(r.Form.Get("fmri"))
	surgery.DTI, _ = strconv.Atoi(r.Form.Get("dti"))
	surgery.OperatorFirst = r.Form.Get("operator_first")
	surgery.OperatorSecond = r.Form.Get("operator_second")
	surgery.HeadFixType, _ = strconv.Atoi(r.Form.Get("head_fix_type"))
	surgery.OperatorSecond = r.Form.Get("cancelation_reason")

	err = m.DB.PutSurgeriesInformation(surgery)
	if err != nil {
		m.App.Session.Put(r.Context(), "error", "can't update surgery information!")
		http.Redirect(rw, r, "/report", http.StatusTemporaryRedirect)
		return
	}

	m.App.Session.Put(r.Context(), "flash", "Saved successfully")
	http.Redirect(rw, r, "/report/detail/"+strconv.Itoa(patient.ID)+"/show", http.StatusSeeOther)
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

	surgeryDay, surgerytime, surgeryarea, surgeryresult, hospitaltype, paymentstatus, headfixtype, imagevalidity := GetAllSelectOptions()

	for index, item := range surgeryDay {
		val, err := strconv.Atoi(item.Value)
		if err != nil {
			continue
		}

		if val == surgeryInfo[0].SurgeryDay {
			surgeryDay[index].Selected = "selected"
		}
	}

	for index, item := range surgerytime {
		val, err := strconv.Atoi(item.Value)
		if err != nil {
			continue
		}

		if val == surgeryInfo[0].SurgeryTime {
			surgerytime[index].Selected = "selected"
		}
	}

	for index, item := range surgeryarea {
		val, err := strconv.Atoi(item.Value)
		if err != nil {
			continue
		}

		if val == surgeryInfo[0].SurgeryArea {
			surgeryarea[index].Selected = "selected"
		}
	}

	for index, item := range surgeryresult {
		val, err := strconv.Atoi(item.Value)
		if err != nil {
			continue
		}

		if val == surgeryInfo[0].SurgeryResult {
			surgeryresult[index].Selected = "selected"
		}
	}

	for index, item := range hospitaltype {
		val, err := strconv.Atoi(item.Value)
		if err != nil {
			continue
		}

		if val == surgeryInfo[0].HospitalType {
			hospitaltype[index].Selected = "selected"
		}
	}

	for index, item := range headfixtype {
		val, err := strconv.Atoi(item.Value)
		if err != nil {
			continue
		}

		if val == surgeryInfo[0].HeadFixType {
			headfixtype[index].Selected = "selected"
		}
	}

	var ct []models.Option
	copier.Copy(&ct, &imagevalidity)

	for index, item := range ct {
		val, err := strconv.Atoi(item.Value)
		if err != nil {
			continue
		}

		if val == surgeryInfo[0].CT {
			ct[index].Selected = "selected"
		}
	}

	var mr []models.Option
	copier.Copy(&mr, &imagevalidity)
	for index, item := range mr {
		val, err := strconv.Atoi(item.Value)
		if err != nil {
			continue
		}

		if val == surgeryInfo[0].MR {
			mr[index].Selected = "selected"
		}
	}

	var fmri []models.Option
	copier.Copy(&fmri, &imagevalidity)
	for index, item := range fmri {
		val, err := strconv.Atoi(item.Value)
		if err != nil {
			continue
		}

		if val == surgeryInfo[0].FMRI {
			fmri[index].Selected = "selected"
		}
	}

	var dti []models.Option
	copier.Copy(&dti, &imagevalidity)
	for index, item := range dti {
		val, err := strconv.Atoi(item.Value)
		if err != nil {
			continue
		}

		if val == surgeryInfo[0].DTI {
			dti[index].Selected = "selected"
		}
	}

	data["patient"] = patient
	data["surgeryinfo"] = surgeryInfo[0]
	data["surgeryday"] = surgeryDay
	data["surgerytime"] = surgerytime
	data["surgeryarea"] = surgeryarea
	data["surgeryresult"] = surgeryresult
	data["hospitaltype"] = hospitaltype
	data["paymentstatus"] = paymentstatus
	data["headfixtype"] = headfixtype
	data["ct"] = ct
	data["mr"] = mr
	data["fmri"] = fmri
	data["dti"] = dti
	data["baseurl"] = "http://localhost:8080"

	m.App.Session.Put(r.Context(), "surgeryinfo", surgeryInfo[0])
	m.App.Session.Put(r.Context(), "patient", patient)

	render.Template(rw, r, "addNewReport.page.html", &models.TemplateData{
		Form: forms.New(nil),
		Data: data,
	})
}

func GetAllSelectOptions() ([]models.Option, []models.Option, []models.Option, []models.Option, []models.Option, []models.Option, []models.Option, []models.Option) {
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
		{Value: "3", Text: "ENT & Neurosurgery", Selected: ""},
		{Value: "4", Text: "CMF", Selected: ""},
		{Value: "5", Text: "Spine", Selected: ""},
		{Value: "6", Text: "Orthopedics", Selected: ""},
	}

	surgeryresult := []models.Option{
		{Value: "1", Text: "Success", Selected: ""},
		{Value: "2", Text: "Canceled", Selected: ""},
		{Value: "3", Text: "Fail", Selected: ""},
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
		{Value: "2", Text: "Mayfield", Selected: ""},
		{Value: "3", Text: "Other", Selected: ""},
	}

	imagevalidity := []models.Option{
		{Value: "1", Text: "Exist", Selected: ""},
		{Value: "2", Text: "Not Exist", Selected: ""},
		{Value: "3", Text: "Exist And Valid", Selected: ""},
		{Value: "4", Text: "Exist Not Valid", Selected: ""},
	}

	return surgeryDay, surgerytime, surgeryarea, surgeryresult, hospitaltype, paymentstatus, headfixtype, imagevalidity
}
