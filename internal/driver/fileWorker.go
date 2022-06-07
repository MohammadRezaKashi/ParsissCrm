package driver

import (
	"ParsissCrm/internal/models"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/jackc/pgtype"
	"github.com/xuri/excelize/v2"
	ptime "github.com/yaa110/go-persian-calendar"
	"golang.org/x/exp/slices"
)

func CreateColumnName() []string {
	col := []string{}
	for _, itr := range []string{"", "A", "B"} {
		for r := 'A'; r <= 'Z'; r++ {
			col = append(col, itr+string(r))
			if itr+string(r) == "BD" {
				break
			}
		}
	}
	return col
}

func CreateMonths() []string {
	months := []string{}
	for r := 1; r <= 12; r++ {
		months = append(months, ptime.Month(r).String())
	}
	return months
}

func ConvertMonthStringToInt(month string) int {
	switch month {
	case "فروردین":
		return 1
		break
	case "اردیبهشت":
		return 2
		break
	case "خرداد":
		return 3
		break
	case "تیر":
		return 4
		break
	case "مرداد":
		return 5
		break
	case "شهریور":
		return 6
		break
	case "مهر":
		return 7
		break
	case "آبان":
		return 8
		break
	case "آذر":
		return 9
		break
	case "دی":
		return 10
		break
	case "بهمن":
		return 11
		break
	case "اسفند":
		return 12
		break
	}
	return 0
}

func OpenExcelFile(path string) excelize.File {
	file, err := excelize.OpenFile(path)
	if err != nil {
		log.Fatal(err)
	}
	return *file
}

var excelCell = [][]string{}

func ParseExcelFile(file *excelize.File) {
	columnName := CreateColumnName()
	months := CreateMonths()
	for _, sheet := range file.GetSheetList() {
		if slices.Contains(months, strings.Split(sheet, " ")[0]) {
			i := 8
			for {
				s, err := file.GetCellValue(sheet, columnName[1]+strconv.Itoa(i))
				if err != nil {
					log.Fatal(err)
				}
				if s == "" {
					break
				}
				var rowCell = []string{}
				rowCell = append(rowCell, strings.Split(sheet, " ")[1])
				rowCell = append(rowCell, strconv.Itoa(ConvertMonthStringToInt(strings.Split(sheet, " ")[0])))
				for index, col := range columnName[5:] {
					s, err := file.GetCellValue(sheet, col+strconv.Itoa(i))
					if index == 4 {
						if s == "" {
							rowCell = nil
							break
						}
					}
					if err != nil {
						log.Fatal(err)
					}
					rowCell = append(rowCell, s)
				}
				if rowCell != nil {
					excelCell = append(excelCell, rowCell)
				}
				i++
			}
		}
	}
}

func GetAllPersonalInformation() []models.PersonalInformation {
	var personalInfo = []models.PersonalInformation{}
	if excelCell == nil {
		return personalInfo
	}
	for _, row := range excelCell {
		personalInfo = append(personalInfo, models.PersonalInformation{
			Name:        row[6],
			PhoneNumber: row[14],
			NationalID:  row[23],
			Address:     row[24],
		})
	}
	return personalInfo
}

func GetAllSurgeriesInformation() []models.SurgeriesInformation {
	var surgeryInfo = []models.SurgeriesInformation{}
	if excelCell == nil {
		return surgeryInfo
	}
	for _, row := range excelCell {
		surgeryInfo = append(surgeryInfo, models.SurgeriesInformation{
			SurgeryDate:       ConvertStringToDate(row[0] + "-" + row[1] + "-" + row[2]),
			SurgeryDay:        row[3],
			SurgeonFirst:      row[7],
			SurgeonSecond:     row[8],
			Resident:          row[9],
			Hospital:          row[10],
			HospitalType:      ConvertHospitalTypeToInt(row[11]),
			SurgeryResult:     ConvertSurgeryResultToInt(row[13]),
			CT:                row[15],
			MR:                row[16],
			SurgeryType:       row[17],
			SurgeryArea:       ConvertSurgeryAreaToInt(row[18]),
			OperatorFirst:     row[19],
			OperatorSecond:    row[20],
			StartTime:         ConvertStringToTimestamp(row[34]),
			StopTime:          ConvertStringToTimestamp(row[35]),
			EnterTime:         ConvertStringToTimestamp(row[36]),
			ExitTime:          ConvertStringToTimestamp(row[37]),
			PatientEnterTime:  ConvertStringToTimestamp(row[38]),
			HeadFixType:       ConvertHeadFixTypeToInt(row[39]),
			CancelationReason: row[40],
		})
	}
	return surgeryInfo
}

func GetAllFinancialInformation() []models.FinancialInformation {
	var financialInfo []models.FinancialInformation
	if excelCell == nil {
		return financialInfo
	}
	for _, row := range excelCell {
		receiptNumber, _ := strconv.Atoi(row[49])
		financialInfo = append(financialInfo, models.FinancialInformation{
			PaymentStatus:      row[5],
			DateOfFirstContact: ConvertStringToDate(row[21]),
			FirstCaller:        row[22],
			DateOfPayment:      ConvertStringToDate(row[25]),
			LastFourDigitsCard: row[26],
			CashAmount:         row[27],
			Bank:               row[28],
			DiscountPercent:    ConvertDiscountPercentToFloat64(row[29]),
			ReasonForDiscount:  row[30],
			CreditAmount:       ConvertCreditAmountToInt(row[31]),
			TypeOfInsurance:    row[32],
			FinancialVerifier:  row[33],
			ReceiptNumber:      receiptNumber,
			ReceiptDate:        ConvertStringToDate(row[50]),
			ReceiptReceiver:    row[51],
		})
	}
	return financialInfo
}

func ConvertCreditAmountToInt(s string) int {
	if s == "**" {
		return 0
	}
	i, _ := strconv.Atoi(s)
	return i
}

func ConvertDiscountPercentToFloat64(s string) float64 {
	if s == "**" {
		return 0
	}
	f, _ := strconv.ParseFloat(s, 32)
	return f
}

func ConvertSurgeryAreaToInt(s string) int {
	switch s {
	case "N":
		return 1
		break
	case "E":
		return 2
		break
	case "E+N":
		return 3
		break
	case "C":
		return 4
		break
	case "S":
		return 5
		break
	case "O":
		return 6
		break
	}
	return 0
}

func ConvertStringToDate(date string) pgtype.Date {
	d := strings.FieldsFunc(date, func(r rune) bool {
		return r == '-' || r == '/'
	})
	if len(d) == 3 {
		t := ptime.Time{}
		year, _ := strconv.Atoi(d[0])
		month, _ := strconv.Atoi(d[1])
		day, _ := strconv.Atoi(d[2])
		t.Set(year, ptime.Month(month), day, 0, 0, 0, 0, ptime.Iran())
		gt := t.Time()
		var pgDate = pgtype.Date{
			Time:             gt,
			Status:           2,
			InfinityModifier: 0,
		}
		return pgDate
	} else {
		return pgtype.Date{}
	}
}

func ConvertStringToTimestamp(ts string) pgtype.Timestamp {
	t, _ := time.Parse("2006-01-02 15:04", "2006-01-02 "+ts)
	var pgTime = pgtype.Timestamp{
		Time:             t,
		Status:           2,
		InfinityModifier: 0,
	}
	return pgTime
}

func ConvertHospitalTypeToInt(s string) int {
	switch s {
	case "خصوصی":
		return 1
		break
	case "دولتی":
		return 2
		break
	}
	return 0
}

func ConvertSurgeryResultToInt(s string) int {
	switch s {
	case "برگزار شد":
		return 1
		break
	case "کنسل شد":
		return 2
		break
	}
	return 0
}

func ConvertHeadFixTypeToInt(s string) int {
	switch s {
	case "میفیلد":
		return 1
		break
	case "هدبند":
		return 2
		break
	}
	return 0
}
