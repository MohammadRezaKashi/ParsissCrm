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
	"time"

	"github.com/jackc/pgtype"

	"github.com/go-chi/chi"
	"github.com/jinzhu/copier"
)

var Repo *Repository

type Repository struct {
	App *config.AppConfig
	DB  repository.DatabaseRepo
}

var baseUrl = "http://localhost:8080"

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

	var filters = models.Filters{}

	_, _, _, _, hospitaltype, _, _ := GetAllSelectOptionsSurgery()

	filters.HospitalTypeOptions = hospitaltype

	data := make(map[string]interface{})
	data["patients"] = patients
	data["filters"] = filters
	data["baseurl"] = baseUrl
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

	surgeryDay, surgerytime, surgeryarea, surgeryresult, hospitaltype, headfixtype, imagevalidity := GetAllSelectOptionsSurgery()
	paymentstatus := GetAllSelectOptionsFinancial()

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
	data["baseurl"] = baseUrl

	render.Template(rw, r, "addNewReport.page.html", &models.TemplateData{
		Form: forms.New(nil),
		Data: data,
	})
}

func (m *Repository) PostAddNewReport(rw http.ResponseWriter, r *http.Request) {
	var patient models.PersonalInformation
	GetPersonalInformationFromForm(patient, r, m, rw)

	id, err := m.DB.AddPersonalInformation(patient)
	if err != nil {
		m.App.Session.Put(r.Context(), "error", "can't add personal information!")
		http.Redirect(rw, r, "/report", http.StatusTemporaryRedirect)
		return
	}

	surgeryinfo := models.SurgeriesInformation{}
	surgeryinfo.FillDefaults()

	GetSurgeryInformationFromForm(surgeryinfo, r, m)

	err = m.DB.AddSurgeriesInformation(surgeryinfo, id)
	if err != nil {
		m.App.Session.Put(r.Context(), "error", "can't add surgeries information!")
		http.Redirect(rw, r, "/report", http.StatusTemporaryRedirect)
		return
	}
	financial := models.FinancialInformation{}
	financial.FillDefaults()

	GetFinancialInformationFromForm(financial, r)

	err = m.DB.AddFinancialInformation(financial, id)
	if err != nil {
		m.App.Session.Put(r.Context(), "error", "can't add financial information!")
		http.Redirect(rw, r, "/report", http.StatusTemporaryRedirect)
		return
	}

	http.Redirect(rw, r, "/report", http.StatusSeeOther)
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

	shouldReturn := GetPersonalInformationFromForm(patient, r, m, rw)
	if shouldReturn {
		return
	}

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

	GetSurgeryInformationFromForm(surgery, r, m)

	err = m.DB.PutSurgeriesInformation(surgery)
	if err != nil {
		m.App.Session.Put(r.Context(), "error", "can't update surgery information!")
		http.Redirect(rw, r, "/report", http.StatusTemporaryRedirect)
		return
	}

	financial, ok := m.App.Session.Get(r.Context(), "financialinfo").(models.FinancialInformation)
	if !ok {
		m.App.Session.Put(r.Context(), "error", "can't get financial information data!")
		http.Redirect(rw, r, "/report", http.StatusTemporaryRedirect)
	}

	GetFinancialInformationFromForm(financial, r)

	err = m.DB.PutFinancialInformation(financial)
	if err != nil {
		m.App.Session.Put(r.Context(), "error", "can't update surgery information!")
		http.Redirect(rw, r, "/report", http.StatusTemporaryRedirect)
		return
	}

	m.App.Session.Put(r.Context(), "flash", "Saved successfully")
	http.Redirect(rw, r, "/report/detail/"+strconv.Itoa(patient.ID)+"/show", http.StatusSeeOther)
}

func GetFinancialInformationFromForm(financial models.FinancialInformation, r *http.Request) {
	financial.PaymentStatus, _ = strconv.Atoi(r.Form.Get("payment_status"))
	date := driver.ConvertStringToDate(r.Form.Get("first_contact"))
	if date.Status == 2 {
		financial.DateOfFirstContact = date
	}
	date = driver.ConvertStringToDate(r.Form.Get("payment_date"))
	if date.Status == 2 {
		financial.DateOfPayment = date
	}
	financial.FirstCaller = r.Form.Get("first_caller")
	financial.LastFourDigitsCard = r.Form.Get("payment_card_number")
	financial.CashAmount = r.Form.Get("payment_receipt_amount")
	financial.Bank = r.Form.Get("bank")
	financial.ReasonForDiscount = r.Form.Get("discount_reason")
	financial.TypeOfInsurance = r.Form.Get("insurance_type")
	financial.FinancialVerifier = r.Form.Get("financial_verifier")
	financial.ReceiptNumber, _ = strconv.Atoi(r.Form.Get("receipt_number"))
	financial.DiscountPercent, _ = strconv.ParseFloat(r.Form.Get("discount_percentage"), 32)
	date = driver.ConvertStringToDate(r.Form.Get("receipt_received_date"))
	if date.Status == 2 {
		financial.ReceiptDate = date
	}
	financial.ReceiptReceiver = r.Form.Get("receipt_receiver")
}

func GetSurgeryInformationFromForm(surgery models.SurgeriesInformation, r *http.Request, m *Repository) {
	surgery.FileNumber = r.Form.Get("file_number")
	date := driver.ConvertStringToDate(r.Form.Get("date_of_hospital_admission"))
	if date.Status == 2 {
		surgery.DateOfHospitalAdmission = date
	}
	date = driver.ConvertStringToDate(r.Form.Get("surgery_date"))
	if date.Status == 2 {
		surgery.SurgeryDate = date
	}
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
	surgery.CancellationReason = r.Form.Get("cancelation_reason")

	layout := "15:04"
	t, err := time.Parse(layout, r.Form.Get("start_time"))
	if err != nil {
		m.App.Session.Put(r.Context(), "error", "can't parse start time!")
	}
	surgery.StartTime = pgtype.Timestamp{Time: t, Status: pgtype.Present}

	t, err = time.Parse(layout, r.Form.Get("stop_time"))
	if err != nil {
		m.App.Session.Put(r.Context(), "error", "can't parse stop time!")
	}
	surgery.StopTime = pgtype.Timestamp{Time: t, Status: pgtype.Present}

	t, err = time.Parse(layout, r.Form.Get("enter_time"))
	if err != nil {
		m.App.Session.Put(r.Context(), "error", "can't parse enter time!")
	}
	surgery.EnterTime = pgtype.Timestamp{Time: t, Status: pgtype.Present}

	t, err = time.Parse(layout, r.Form.Get("exit_time"))
	if err != nil {
		m.App.Session.Put(r.Context(), "error", "can't parse exit time!")
	}
	surgery.ExitTime = pgtype.Timestamp{Time: t, Status: pgtype.Present}

	t, err = time.Parse(layout, r.Form.Get("patient_enter_time"))
	if err != nil {
		m.App.Session.Put(r.Context(), "error", "can't parse patient enter time!")
	}
	surgery.PatientEnterTime = pgtype.Timestamp{Time: t, Status: pgtype.Present}
}

func GetPersonalInformationFromForm(patient models.PersonalInformation, r *http.Request, m *Repository, rw http.ResponseWriter) bool {
	patient.Name = r.Form.Get("name")
	patient.Family = r.Form.Get("family")
	age, err := strconv.Atoi(r.Form.Get("age"))
	if err != nil {
		m.App.Session.Put(r.Context(), "error", "can't parse age!")
		http.Redirect(rw, r, "/report", http.StatusTemporaryRedirect)
		return true
	}
	patient.Age = age
	patient.Address = r.Form.Get("address")
	patient.Email = r.Form.Get("email")
	patient.NationalID = r.Form.Get("national_id")
	patient.PhoneNumber = r.Form.Get("phone_number")
	patient.PlaceOfBirth = r.Form.Get("place_of_birth")
	return false
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

	financialInfo, err := m.DB.GetFinancialInformationByPatientID(id)
	if err != nil {
		helpers.ServerError(rw, err)
		return
	}

	data := make(map[string]interface{})

	surgeryDay, surgerytime, surgeryarea, surgeryresult, hospitaltype, headfixtype, ct, mr, fmri, dti := InitSurgeryInfoComboBoxes(surgeryInfo)

	paymentstatus := InitFinancialInfoComboBoxes(financialInfo)

	data["patient"] = patient
	data["surgeryinfo"] = surgeryInfo[0]
	data["financialinfo"] = financialInfo[0]
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
	data["baseurl"] = baseUrl

	m.App.Session.Put(r.Context(), "surgeryinfo", surgeryInfo[0])
	m.App.Session.Put(r.Context(), "patient", patient)
	m.App.Session.Put(r.Context(), "financialinfo", financialInfo[0])

	render.Template(rw, r, "addNewReport.page.html", &models.TemplateData{
		Form: forms.New(nil),
		Data: data,
	})
}

func InitFinancialInfoComboBoxes(financialInfo []models.FinancialInformation) []models.Option {
	paymentstatus := GetAllSelectOptionsFinancial()
	for index, item := range paymentstatus {
		val, err := strconv.Atoi(item.Value)
		if err != nil {
			continue
		}
		if val == financialInfo[0].PaymentStatus {
			paymentstatus[index].Selected = "selected"
		}
	}
	return paymentstatus
}

func InitSurgeryInfoComboBoxes(surgeryInfo []models.SurgeriesInformation) ([]models.Option, []models.Option, []models.Option, []models.Option, []models.Option, []models.Option, []models.Option, []models.Option, []models.Option, []models.Option) {
	surgeryDay, surgerytime, surgeryarea, surgeryresult, hospitaltype, headfixtype, imagevalidity := GetAllSelectOptionsSurgery()

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
	return surgeryDay, surgerytime, surgeryarea, surgeryresult, hospitaltype, headfixtype, ct, mr, fmri, dti
}

func GetAllSelectOptionsFinancial() []models.Option {
	paymentstatus := []models.Option{
		{Value: "1", Text: "Paid", Selected: ""},
		{Value: "2", Text: "Unpaid", Selected: ""},
		{Value: "3", Text: "Free", Selected: ""},
		{Value: "4", Text: "Health Plan", Selected: ""},
		{Value: "5", Text: "Paid By Hospital", Selected: ""},
	}

	return paymentstatus
}

func GetAllSelectOptionsSurgery() ([]models.Option, []models.Option, []models.Option, []models.Option, []models.Option, []models.Option, []models.Option) {
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

	return surgeryDay, surgerytime, surgeryarea, surgeryresult, hospitaltype, headfixtype, imagevalidity
}

func (m *Repository) ShowFilters(w http.ResponseWriter, r *http.Request) {
}
