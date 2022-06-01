package models

import (
	"github.com/jackc/pgtype"
)

type PersonalInformation struct {
	ID          int
	Name        string
	PhoneNumber string
	NationalID  string
	Address     string
}

type SurgeriesInformation struct {
	ID                 int
	PatientID          int
	SurgeryDate        pgtype.Date
	SurgeryDay         string
	SurgeryType        string
	SurgeryArea        int
	SurgeryDescription string
	SurgeryResult      int
	SurgeonFirst       string
	SurgeonSecond      string
	Resident           string
	Hospital           string
	HospitalType       int
	HospitalAddress    string
	CT                 string
	MR                 string
	OperatorFirst      string
	OperatorSecond     string
	StartTime          pgtype.Timestamp
	StopTime           pgtype.Timestamp
	EnterTime          pgtype.Timestamp
	ExitTime           pgtype.Timestamp
	PatientEnterTime   pgtype.Timestamp
	HeadFixType        int
	CancelationReason  string
}
