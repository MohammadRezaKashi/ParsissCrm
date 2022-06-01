package dbrepo

import (
	"ParsissCrm/internal/models"
	"context"
	"log"
	"time"
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
	                                          hospital_type, hospital_address, ct, mr, operator_first, operator_second,
	                                          start_time, stop_time, enter_time, exit_time, enter_patient_time,
	                                          head_fix_type, cancellation_reason)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18, $19, $20, $21, $22,
		        $23, $24)`
	_, err := m.DB.ExecContext(ctx, query, patientID, surgeriesInformation.SurgeryDate, surgeriesInformation.SurgeryDay,
		surgeriesInformation.SurgeryType, surgeriesInformation.SurgeryArea, surgeriesInformation.SurgeryDescription,
		surgeriesInformation.SurgeryResult, surgeriesInformation.SurgeonFirst, surgeriesInformation.SurgeonSecond,
		surgeriesInformation.Resident, surgeriesInformation.Hospital, surgeriesInformation.HospitalType,
		surgeriesInformation.HospitalAddress, surgeriesInformation.CT, surgeriesInformation.MR,
		surgeriesInformation.OperatorFirst, surgeriesInformation.OperatorSecond,
		surgeriesInformation.StartTime.Time, surgeriesInformation.StopTime.Time, surgeriesInformation.EnterTime.Time,
		surgeriesInformation.ExitTime.Time, surgeriesInformation.PatientEnterTime.Time, surgeriesInformation.HeadFixType,
		surgeriesInformation.CancelationReason)
	if err != nil {
		log.Println(err)
		return err
	}
	return nil
}
