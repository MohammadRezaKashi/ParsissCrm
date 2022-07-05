package dbrepo

import (
	"ParsissCrm/internal/driver"
	"ParsissCrm/internal/models"
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/jackc/pgtype"
	"github.com/nleeper/goment"
)

func (m *postgresDBRepo) AddPersonalInformation(personalInfo models.PersonalInformation) (int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := `
	INSERT INTO public."PatientsInformation"(
		name, family, age, phone, national_id, address, email, place_of_birth)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8) RETURNING id`

	err := m.DB.QueryRowContext(ctx, query, personalInfo.Name, personalInfo.Family, personalInfo.Age,
		personalInfo.PhoneNumber, personalInfo.NationalID, personalInfo.Address, personalInfo.Email,
		personalInfo.PlaceOfBirth).Scan(&personalInfo.ID)
	if err != nil {
		log.Println(err)
		return 0, err
	}
	return personalInfo.ID, nil
}

func (m *postgresDBRepo) AddSurgeriesInformation(surgeriesInformation models.SurgeriesInformation, patientID int) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := `
	INSERT INTO public."SurgeriesInformation"(
		patient_id, surgery_date, surgery_day, surgery_type, surgery_area, surgery_description,
	                                          surgery_result, surgeon_first, surgeon_second, resident, hospital,
	                                          hospital_type, hospital_address, ct, mr, fmri, dti, operator_first, operator_second,
	                                          start_time, stop_time, enter_time, exit_time, patient_enter_time,
	                                          head_fix_type, cancellation_reason, file_number, date_of_hospital_admission, surgery_time)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18, $19, $20, $21, $22,
		        $23, $24, $25, $26, $27, $28, $29)`

	_, err := m.DB.ExecContext(ctx, query, patientID, surgeriesInformation.SurgeryDate, surgeriesInformation.SurgeryDay,
		surgeriesInformation.SurgeryType, surgeriesInformation.SurgeryArea, surgeriesInformation.SurgeryDescription,
		surgeriesInformation.SurgeryResult, surgeriesInformation.SurgeonFirst, surgeriesInformation.SurgeonSecond,
		surgeriesInformation.Resident, surgeriesInformation.Hospital, surgeriesInformation.HospitalType,
		surgeriesInformation.HospitalAddress, surgeriesInformation.CT, surgeriesInformation.MR, surgeriesInformation.FMRI,
		surgeriesInformation.DTI, surgeriesInformation.OperatorFirst, surgeriesInformation.OperatorSecond,
		surgeriesInformation.StartTime.Time, surgeriesInformation.StopTime.Time, surgeriesInformation.EnterTime.Time,
		surgeriesInformation.ExitTime.Time, surgeriesInformation.PatientEnterTime.Time, surgeriesInformation.HeadFixType,
		surgeriesInformation.CancellationReason, surgeriesInformation.FileNumber, surgeriesInformation.DateOfHospitalAdmission.Time,
		surgeriesInformation.SurgeryTime)
	if err != nil {
		log.Println(err)
		return err
	}
	return nil
}

func (m *postgresDBRepo) AddFinancialInformation(financialInfo models.FinancialInformation, patientID int) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := `
	INSERT INTO public."FinancialInformation"(
		patient_id, payment_status, payment_date, payment_amount, payment_notes, receipt_number,
	                                          first_contact, first_caller, last_four_card_number, bank,
	                                          discount_percentage, discount_reason, credit_amount, insurance_type,
	                                          financial_verifier, receipt_received_date, receipt_receiver)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17)`
	_, err := m.DB.ExecContext(ctx, query, patientID, financialInfo.PaymentStatus, financialInfo.DateOfPayment.Time,
		financialInfo.CashAmount, financialInfo.PaymentNote, financialInfo.ReceiptNumber, financialInfo.DateOfFirstContact.Time,
		financialInfo.FirstCaller, financialInfo.LastFourDigitsCard, financialInfo.Bank, financialInfo.DiscountPercent,
		financialInfo.ReasonForDiscount, financialInfo.CreditAmount, financialInfo.TypeOfInsurance,
		financialInfo.FinancialVerifier, financialInfo.ReceiptDate.Time, financialInfo.ReceiptReceiver)
	if err != nil {
		log.Println(err)
		return err
	}
	return nil
}

func (m *postgresDBRepo) GetAllPatients() ([]models.PersonalInformation, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := `SELECT * FROM public."PatientsInformation"
	ORDER BY id DESC `

	var patients []models.PersonalInformation

	rows, err := m.DB.QueryContext(ctx, query)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var patient models.PersonalInformation
		err := rows.Scan(&patient.ID, &patient.Name, &patient.Family, &patient.Age, &patient.PhoneNumber, &patient.NationalID, &patient.Address, &patient.Email, &patient.PlaceOfBirth, new(time.Time), new(time.Time))
		if err != nil {
			log.Println(err)
			return nil, err
		}
		patients = append(patients, patient)
	}

	return patients, nil
}

func (m *postgresDBRepo) GetPatientByID(id int) (models.PersonalInformation, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := `SELECT * FROM public."PatientsInformation"
	WHERE id = $1`

	var patient models.PersonalInformation

	err := m.DB.QueryRowContext(ctx, query, id).Scan(&patient.ID, &patient.Name, &patient.Family, &patient.Age, &patient.PhoneNumber, &patient.NationalID, &patient.Address, &patient.Email, &patient.PlaceOfBirth, new(time.Time), new(time.Time))
	if err != nil {
		log.Println(err)
		return models.PersonalInformation{}, err
	}
	return patient, nil
}

func (m *postgresDBRepo) GetSurgicalInformationByPatientID(id int) ([]models.SurgeriesInformation, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := `SELECT * FROM public."SurgeriesInformation"
	WHERE patient_id = $1
	ORDER BY id ASC`

	var surgeries []models.SurgeriesInformation

	rows, err := m.DB.QueryContext(ctx, query, id)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var startTime, stopTime, enterTime, exitTime, patientEnterTime pgtype.Time
		var surgery models.SurgeriesInformation

		err := rows.Scan(&surgery.ID, &surgery.PatientID, &surgery.SurgeryDate, &surgery.SurgeryDay, &surgery.SurgeryTime,
			&surgery.SurgeryType, &surgery.SurgeryArea, &surgery.SurgeryDescription, &surgery.SurgeryResult,
			&surgery.SurgeonFirst, &surgery.SurgeonSecond, &surgery.Resident, &surgery.Hospital,
			&surgery.HospitalType, &surgery.HospitalAddress, &surgery.CT, &surgery.MR, &surgery.FMRI, &surgery.DTI,
			&surgery.OperatorFirst, &surgery.OperatorSecond, &startTime, &stopTime, &enterTime, &exitTime,
			&patientEnterTime, &surgery.HeadFixType, &surgery.CancellationReason, &surgery.FileNumber,
			&surgery.DateOfHospitalAdmission, new(time.Time), new(time.Time))
		if err != nil {
			log.Println(err)
			return nil, err
		}
		st, _ := goment.New(startTime.Microseconds * 1000)
		surgery.StartTime.Time = st.ToTime().UTC()
		st, _ = goment.New(stopTime.Microseconds * 1000)
		surgery.StopTime.Time = st.ToTime().UTC()
		st, _ = goment.New(enterTime.Microseconds * 1000)
		surgery.EnterTime.Time = st.ToTime().UTC()
		st, _ = goment.New(exitTime.Microseconds * 1000)
		surgery.ExitTime.Time = st.ToTime().UTC()
		st, _ = goment.New(patientEnterTime.Microseconds * 1000)
		surgery.PatientEnterTime.Time = st.ToTime().UTC()

		surgeries = append(surgeries, surgery)
	}

	return surgeries, nil
}

func (m *postgresDBRepo) GetFinancialInformationByPatientID(id int) ([]models.FinancialInformation, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := `SELECT * FROM public."FinancialInformation"
	WHERE patient_id = $1
	ORDER BY id ASC`

	var financials []models.FinancialInformation

	rows, err := m.DB.QueryContext(ctx, query, id)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var financial models.FinancialInformation

		err := rows.Scan(&financial.ID, &financial.PatientID, &financial.PaymentStatus, &financial.DateOfPayment,
			&financial.CashAmount, &financial.PaymentNote, &financial.ReceiptNumber, &financial.DateOfFirstContact,
			&financial.FirstCaller, &financial.LastFourDigitsCard, &financial.Bank, &financial.DiscountPercent,
			&financial.ReasonForDiscount, &financial.CreditAmount, &financial.TypeOfInsurance, &financial.FinancialVerifier,
			&financial.ReceiptDate, &financial.ReceiptReceiver, new(time.Time), new(time.Time))
		if err != nil {
			log.Println(err)
			return nil, err
		}
		financials = append(financials, financial)
	}
	return financials, nil
}

func (m *postgresDBRepo) PutPersonalInformation(personalInfo models.PersonalInformation) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := `
	UPDATE public."PatientsInformation"
	SET name = $1, family = $2, age = $3, phone = $4, national_id = $5, address = $6, email = $7,
	    place_of_birth = $8, updated_at = $9
	WHERE id = $10`
	_, err := m.DB.ExecContext(ctx, query, personalInfo.Name, personalInfo.Family, personalInfo.Age,
		personalInfo.PhoneNumber, personalInfo.NationalID, personalInfo.Address, personalInfo.Email,
		personalInfo.PlaceOfBirth, time.Now(), personalInfo.ID)
	if err != nil {
		log.Println(err)
		return err
	}
	return nil
}

func (m *postgresDBRepo) PutSurgeriesInformation(surgeriesInfo models.SurgeriesInformation) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	if surgeriesInfo.DateOfHospitalAdmission.Status == pgtype.Undefined {
		surgeriesInfo.DateOfHospitalAdmission.Status = pgtype.Present
	}

	query := `
	UPDATE public."SurgeriesInformation"
	SET patient_id = $1, surgery_date = $2, surgery_day = $3, surgery_type = $4, surgery_area = $5,
	    surgery_description = $6, surgery_result = $7, surgeon_first = $8, surgeon_second = $9,
	    resident = $10, hospital = $11, hospital_type = $12, hospital_address = $13, ct = $14,
	    mr = $15, fmri = $16, dti = $17, operator_first = $18, operator_second = $19, start_time = $20,
	    stop_time = $21, enter_time = $22, exit_time = $23, patient_enter_time = $24, head_fix_type = $25,
	    cancellation_reason = $26, file_number = $27, date_of_hospital_admission = $28, updated_at = $29, surgery_time = $31
	WHERE id = $30`
	_, err := m.DB.ExecContext(ctx, query, surgeriesInfo.PatientID, surgeriesInfo.SurgeryDate,
		surgeriesInfo.SurgeryDay, surgeriesInfo.SurgeryType, surgeriesInfo.SurgeryArea,
		surgeriesInfo.SurgeryDescription, surgeriesInfo.SurgeryResult, surgeriesInfo.SurgeonFirst,
		surgeriesInfo.SurgeonSecond, surgeriesInfo.Resident, surgeriesInfo.Hospital,
		surgeriesInfo.HospitalType, surgeriesInfo.HospitalAddress, surgeriesInfo.CT,
		surgeriesInfo.MR, surgeriesInfo.FMRI, surgeriesInfo.DTI, surgeriesInfo.OperatorFirst,
		surgeriesInfo.OperatorSecond, surgeriesInfo.StartTime.Time, surgeriesInfo.StopTime.Time,
		surgeriesInfo.EnterTime.Time, surgeriesInfo.ExitTime.Time, surgeriesInfo.PatientEnterTime.Time,
		surgeriesInfo.HeadFixType, surgeriesInfo.CancellationReason, surgeriesInfo.FileNumber,
		surgeriesInfo.DateOfHospitalAdmission, time.Now(), surgeriesInfo.ID, surgeriesInfo.SurgeryTime)
	if err != nil {
		log.Println(err)
		return err
	}
	return nil
}

func (m *postgresDBRepo) PutFinancialInformation(financialInfo models.FinancialInformation) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := `
	UPDATE public."FinancialInformation"
	SET patient_id = $1, payment_status = $2, payment_date = $3, payment_amount = $4, payment_notes = $5,
	    receipt_number = $6, first_contact = $7, first_caller = $8, last_four_card_number = $9,
	    bank = $10, discount_percentage = $11, discount_reason = $12, credit_amount = $13, insurance_type = $14,
	    financial_verifier = $15, receipt_received_date = $16, receipt_receiver = $17, updated_at = $18
	WHERE id = $19`
	_, err := m.DB.ExecContext(ctx, query, financialInfo.PatientID, financialInfo.PaymentStatus,
		financialInfo.DateOfPayment, financialInfo.CashAmount, financialInfo.PaymentNote, financialInfo.ReceiptNumber,
		financialInfo.DateOfFirstContact, financialInfo.FirstCaller, financialInfo.LastFourDigitsCard,
		financialInfo.Bank, financialInfo.DiscountPercent, financialInfo.ReasonForDiscount, financialInfo.CreditAmount,
		financialInfo.TypeOfInsurance, financialInfo.FinancialVerifier, financialInfo.ReceiptDate,
		financialInfo.ReceiptReceiver, time.Now(), financialInfo.ID)
	if err != nil {
		log.Println(err)
		return err
	}
	return nil
}

func (m *postgresDBRepo) GetDistinctList(tableName string, columnName string) ([]interface{}, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := fmt.Sprintf("SELECT DISTINCT %s FROM public.\"%s\"", columnName, tableName)

	rows, err := m.DB.QueryContext(ctx, query)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	var list []interface{}
	for rows.Next() {
		var value interface{}
		err := rows.Scan(&value)
		if err != nil {
			log.Println(err)
			return nil, err
		}
		list = append(list, value)
	}
	return list, nil
}

func (m *postgresDBRepo) GetFilterData(filter interface{}) ([]models.PersonalInformation, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := `SELECT pa.* FROM public."PatientsInformation" pa
	FULL OUTER JOIN public."SurgeriesInformation" su ON pa.id = su.patient_id
	FULL OUTER JOIN public."FinancialInformation" fi ON pa.id = fi.patient_id
	`
	f := filter.(map[string]interface{})
	if len(f) > 0 {
		query += "WHERE "
	}
	var ands []string
	for key, value := range filter.(map[string]interface{}) {
		var ors []string
		var typeOfValue string
		for k, v := range value.(map[string]interface{}) {
			if k == "type" {
				typeOfValue = v.(string)
			} else {
				for i, v := range v.([]interface{}) {
					if typeOfValue == "checkbox" {
						ors = append(ors, fmt.Sprintf("%s = '%s'", key, v.(string)))
					} else if typeOfValue == "date" {
						d := driver.ConvertStringToDate(v.(string))
						switch i {
						case 0:
							ors = append(ors, fmt.Sprintf("%s >= '%s'", key, d.Time.Format("2006-01-02")))
						case 1:
							ors = append(ors, fmt.Sprintf("%s <= '%s'", key, d.Time.Format("2006-01-02")))
						}
					}
				}
			}
		}
		if typeOfValue == "checkbox" {
			ands = append(ands, fmt.Sprintf("(%s)", strings.Join(ors, " OR ")))
		} else if typeOfValue == "date" {
			ands = append(ands, fmt.Sprintf("(%s)", strings.Join(ors, " AND ")))
		}

	}
	query += strings.Join(ands, " AND ")
	rows, err := m.DB.QueryContext(ctx, query)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	defer rows.Close()

	var patients []models.PersonalInformation
	for rows.Next() {
		//var surgery models.SurgeriesInformation
		//
		//err := rows.Scan(&surgery.ID, &surgery.PatientID, &surgery.SurgeryDate, &surgery.SurgeryDay, &surgery.SurgeryTime, &surgery.SurgeryType,
		//	&surgery.SurgeryArea, &surgery.SurgeryDescription, &surgery.SurgeryResult, &surgery.SurgeonFirst, &surgery.SurgeonSecond,
		//	&surgery.Resident, &surgery.Hospital, &surgery.HospitalType, &surgery.HospitalAddress, &surgery.CT, &surgery.MR, &surgery.FMRI,
		//	&surgery.DTI, &surgery.OperatorFirst, &surgery.OperatorSecond, &surgery.StartTime.Time, &surgery.StopTime.Time, &surgery.EnterTime.Time,
		//	&surgery.ExitTime.Time, &surgery.PatientEnterTime.Time, &surgery.HeadFixType, &surgery.CancellationReason, &surgery.FileNumber,
		//	&surgery.DateOfHospitalAdmission, new(time.Time), new(time.Time))
		var patient models.PersonalInformation
		err := rows.Scan(&patient.ID, &patient.Name, &patient.Family, &patient.Age, &patient.PhoneNumber,
			&patient.NationalID, &patient.Address, &patient.Email, &patient.PlaceOfBirth, new(time.Time), new(time.Time))
		if err != nil {
			log.Println(err)
			return nil, err
		}
		patients = append(patients, patient)
	}
	return patients, nil
}
