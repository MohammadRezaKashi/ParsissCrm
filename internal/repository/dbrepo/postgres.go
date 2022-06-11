package dbrepo

import (
	"ParsissCrm/internal/models"
	"context"
	"log"
	"time"

	"github.com/jackc/pgtype"
	"github.com/nleeper/goment"
)

func (m *postgresDBRepo) AddPersonalInformation(personalInfo models.PersonalInformation) (int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := `
	INSERT INTO public."PatientsInformation"(
		name, phone, national_id, address, email)
		VALUES ($1, $2, $3, $4, $5) RETURNING id`
	err := m.DB.QueryRowContext(ctx, query, personalInfo.Name, personalInfo.PhoneNumber,
		personalInfo.NationalID, personalInfo.Address, "").Scan(&personalInfo.ID)
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
<<<<<<< HEAD
	                                          start_time, stop_time, enter_time, exit_time, patient_enter_time,
	                                          head_fix_type, cancellation_reason, file_number, date_of_hospital_admission, surgery_time)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18, $19, $20, $21, $22,
		        $23, $24, $25, $26, $27, $28, $29)`
=======
	                                          start_time, stop_time, enter_time, exit_time, enter_patient_time,
	                                          head_fix_type, cancellation_reason, file_number, date_of_hospital_admission)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18, $19, $20, $21, $22,
		        $23, $24, $25, $26, $27, $28)`
>>>>>>> 0156e4322382c6f487914cb5381720cd572d6914
	_, err := m.DB.ExecContext(ctx, query, patientID, surgeriesInformation.SurgeryDate, surgeriesInformation.SurgeryDay,
		surgeriesInformation.SurgeryType, surgeriesInformation.SurgeryArea, surgeriesInformation.SurgeryDescription,
		surgeriesInformation.SurgeryResult, surgeriesInformation.SurgeonFirst, surgeriesInformation.SurgeonSecond,
		surgeriesInformation.Resident, surgeriesInformation.Hospital, surgeriesInformation.HospitalType,
		surgeriesInformation.HospitalAddress, surgeriesInformation.CT, surgeriesInformation.MR, surgeriesInformation.FMRI,
		surgeriesInformation.DTI, surgeriesInformation.OperatorFirst, surgeriesInformation.OperatorSecond,
		surgeriesInformation.StartTime.Time, surgeriesInformation.StopTime.Time, surgeriesInformation.EnterTime.Time,
		surgeriesInformation.ExitTime.Time, surgeriesInformation.PatientEnterTime.Time, surgeriesInformation.HeadFixType,
<<<<<<< HEAD
		surgeriesInformation.CancellationReason, surgeriesInformation.FileNumber, surgeriesInformation.DateOfHospitalAdmission.Time,
		surgeriesInformation.SurgeryTime)
=======
		surgeriesInformation.CancellationReason, "", surgeriesInformation.DateOfHospitalAdmission.Time)
>>>>>>> 0156e4322382c6f487914cb5381720cd572d6914
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
		patient_id, payment_status, payment_date, payment_amount, payment_type, payment_notes, receipt_number,
	                                          first_contact, first_caller, last_four_card_number, bank,
	                                          discount_percentage, discount_reason, credit_amount, insurance_type,
	                                          financial_verifier, receipt_received_date, receipt_receiver)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18)`
	_, err := m.DB.ExecContext(ctx, query, patientID, financialInfo.PaymentStatus, financialInfo.DateOfPayment.Time,
		financialInfo.CashAmount, 0, "", financialInfo.ReceiptNumber, financialInfo.DateOfFirstContact.Time,
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
	ORDER BY id ASC `

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

func (m *postgresDBRepo) PutPersonalInformation(personalInfo models.PersonalInformation) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := `
	UPDATE public."PatientsInformation"
	SET name = $1, family = $2, age = $3, phone = $4, national_id = $5, address = $6, email = $7,
	    place_of_birthday = $8, updated_at = $9
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
	SET patient_id = $1, payment_status = $2, payment_date = $3, payment_amount = $4, payment_type = $5,
	    payment_notes = $6, receipt_number = $7, first_contact = $8, first_caller = $9, last_four_card_number = $10,
	    bank = $11, discount_percentage = $12, discount_reason = $13, credit_amount = $14, insurance_type = $15,
	    financial_verifier = $16, receipt_received_date = $17, receipt_receiver = $18, updated_at = $19
	WHERE id = $27`
	_, err := m.DB.ExecContext(ctx, query, financialInfo.PatientID, financialInfo.PaymentStatus,
		financialInfo.DateOfPayment, financialInfo.CashAmount, 0, "", financialInfo.ReceiptNumber,
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
