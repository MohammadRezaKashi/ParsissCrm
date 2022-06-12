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
	SurgeryDay              int
	SurgeryTime             int
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
	CT                      int
	MR                      int
	FMRI                    int
	DTI                     int
	OperatorFirst           string
	OperatorSecond          string
	StartTime               pgtype.Timestamp
	StopTime                pgtype.Timestamp
	EnterTime               pgtype.Timestamp
	ExitTime                pgtype.Timestamp
	PatientEnterTime        pgtype.Timestamp
	HeadFixType             int
	CancellationReason      string
	FileNumber              string
	DateOfHospitalAdmission pgtype.Date
}

type FinancialInformation struct {
	ID                 int
	PatientID          int
	PaymentStatus      int
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

type SoftwareInformation struct {
	FRE                     float64
	RegistrationTime        int
	RegistrationTryNumber   int
	SurfaceRegistrationTime int
	RegistrationToolName    string
}

type Option struct {
	Value, Text string
	Selected    string
}
