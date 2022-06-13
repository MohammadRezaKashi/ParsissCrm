package models

import (
	"time"

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
	FailReason              string
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

func (s *SurgeriesInformation) FillDefaults() {
	s.SurgeryDate = pgtype.Date{Time: time.Time{}, Status: pgtype.Present}
	s.StartTime = pgtype.Timestamp{Time: time.Time{}, Status: pgtype.Present}
	s.StopTime = pgtype.Timestamp{Time: time.Time{}, Status: pgtype.Present}
	s.EnterTime = pgtype.Timestamp{Time: time.Time{}, Status: pgtype.Present}
	s.ExitTime = pgtype.Timestamp{Time: time.Time{}, Status: pgtype.Present}
	s.PatientEnterTime = pgtype.Timestamp{Time: time.Time{}, Status: pgtype.Present}
	s.DateOfHospitalAdmission = pgtype.Date{Time: time.Time{}, Status: pgtype.Present}
}

type FinancialInformation struct {
	ID                 int
	PatientID          int
	PaymentStatus      int
	DateOfFirstContact pgtype.Date
	PaymentNote        string
	FirstCaller        string
	DateOfPayment      pgtype.Date
	LastFourDigitsCard string
	CashAmount         string
	Bank               string
	DiscountPercent    float64
	ReasonForDiscount  string
	CreditAmount       string
	TypeOfInsurance    string
	FinancialVerifier  string
	ReceiptNumber      int
	ReceiptDate        pgtype.Date
	ReceiptReceiver    string
}

func (f *FinancialInformation) FillDefaults() {
	f.DateOfFirstContact = pgtype.Date{Time: time.Time{}, Status: pgtype.Present}
	f.DateOfPayment = pgtype.Date{Time: time.Time{}, Status: pgtype.Present}
	f.ReceiptDate = pgtype.Date{Time: time.Time{}, Status: pgtype.Present}
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
