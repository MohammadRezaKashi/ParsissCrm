package models

import (
	"github.com/jackc/pgtype"
)

type PersonalInformation struct {
	ID           int
	Name         string
	Family       string
	Age          int
	PhoneNumber  string
	NationalID   string
	Address      string
	Email        string
	PlaceOfBirth string
}

type SurgeriesInformation struct {
	ID                      int
	PatientID               int
	SurgeryDate             pgtype.Date
	SurgeryDay              string
	SurgeryType             string
	SurgeryArea             int
	SurgeryDescription      string
	SurgeryResult           int
	SurgeonFirst            string
	SurgeonSecond           string
	Resident                string
	Hospital                string
	HospitalType            int
	HospitalAddress         string
	CT                      string
	MR                      string
	OperatorFirst           string
	OperatorSecond          string
	StartTime               pgtype.Timestamp
	StopTime                pgtype.Timestamp
	EnterTime               pgtype.Timestamp
	ExitTime                pgtype.Timestamp
	PatientEnterTime        pgtype.Timestamp
	HeadFixType             int
	CancelationReason       string
	FileNumber              string
	DateOfHospitalAdmission pgtype.Date
}

type FinancialInformation struct {
	ID                 int
	PatientID          int
	PaymentStatus      string
	DateOfFirstContact pgtype.Date
	FirstCaller        string
	DateOfPayment      pgtype.Date
	LastFourDigitsCard string
	CashAmount         string
	Bank               string
	DiscountPercent    float64
	ReasonForDiscount  string
	CreditAmount       int
	TypeOfInsurance    string
	FinancialVerifier  string
	ReceiptNumber      int
	ReceiptDate        pgtype.Date
	ReceiptReceiver    string
}
