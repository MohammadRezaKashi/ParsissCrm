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
	                                          start_time, stop_time, enter_time, exit_time, enter_patient_time,
	                                          head_fix_type, cancellation_reason)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18, $19, $20, $21, $22,
		        $23, $24, $25, $26)`
	_, err := m.DB.ExecContext(ctx, query, patientID, surgeriesInformation.SurgeryDate, surgeriesInformation.SurgeryDay,
		surgeriesInformation.SurgeryType, surgeriesInformation.SurgeryArea, surgeriesInformation.SurgeryDescription,
		surgeriesInformation.SurgeryResult, surgeriesInformation.SurgeonFirst, surgeriesInformation.SurgeonSecond,
		surgeriesInformation.Resident, surgeriesInformation.Hospital, surgeriesInformation.HospitalType,
		surgeriesInformation.HospitalAddress, surgeriesInformation.CT, surgeriesInformation.MR, surgeriesInformation.FMRI,
		surgeriesInformation.FMRI, surgeriesInformation.OperatorFirst, surgeriesInformation.OperatorSecond,
		surgeriesInformation.StartTime.Time, surgeriesInformation.StopTime.Time, surgeriesInformation.EnterTime.Time,
		surgeriesInformation.ExitTime.Time, surgeriesInformation.PatientEnterTime.Time, surgeriesInformation.HeadFixType,
		surgeriesInformation.CancelationReason)
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
	                                          payment_receipt_date, first_contact, first_caller,
	                                          last_four_card_number, bank, discount_percentage, discount_reason,
	                                          credit_amount, insurance_type, financial_verifier, receipt_received_date,
	                                          receipt_receiver)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18, $19)`
	_, err := m.DB.ExecContext(ctx, query, patientID, financialInfo.PaymentStatus, financialInfo.DateOfPayment.Time,
		financialInfo.CashAmount, 0, "", financialInfo.ReceiptNumber, financialInfo.ReceiptDate.Time,
		financialInfo.DateOfFirstContact.Time, financialInfo.FirstCaller, financialInfo.LastFourDigitsCard,
		financialInfo.Bank, financialInfo.DiscountPercent, financialInfo.ReasonForDiscount, financialInfo.CreditAmount,
		financialInfo.TypeOfInsurance, financialInfo.FinancialVerifier, financialInfo.ReceiptDate.Time,
		financialInfo.ReceiptReceiver)
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
		err := rows.Scan(&patient.ID, &patient.Name, &patient.PhoneNumber, &patient.NationalID, &patient.Address, &patient.Email, new(time.Time), new(time.Time))
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

	err := m.DB.QueryRowContext(ctx, query, id).Scan(&patient.ID, &patient.Name, &patient.PhoneNumber, &patient.NationalID, &patient.Address, &patient.Email, new(time.Time), new(time.Time))
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
		err := rows.Scan(&surgery.ID, &surgery.PatientID, &surgery.SurgeryDate, &surgery.SurgeryDay,
			&surgery.SurgeryType, &surgery.SurgeryArea, &surgery.SurgeryDescription, &surgery.SurgeryResult,
			&surgery.SurgeonFirst, &surgery.SurgeonSecond, &surgery.Resident, &surgery.Hospital,
			&surgery.HospitalType, &surgery.HospitalAddress, &surgery.CT, &surgery.MR, &surgery.FMRI, &surgery.DTI,
			&surgery.OperatorFirst, &surgery.OperatorSecond, &startTime, &stopTime, &enterTime, &exitTime,
			&patientEnterTime, &surgery.HeadFixType, &surgery.CancelationReason, new(time.Time), new(time.Time))
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
