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
	var col []string
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
	var months []string
	for r := 1; r <= 12; r++ {
		months = append(months, ptime.Month(r).String())
	}
	return months
}

func ConvertMonthStringToInt(month string) int {
	switch month {
	case "فروردین":
		return 1
	case "اردیبهشت":
		return 2
	case "خرداد":
		return 3
	case "تیر":
		return 4
	case "مرداد":
		return 5
	case "شهریور":
		return 6
	case "مهر":
		return 7
	case "آبان":
		return 8
	case "آذر":
		return 9
	case "دی":
		return 10
	case "بهمن":
		return 11
	case "اسفند":
		return 12
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

var excelCell [][]string

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
				var rowCell []string
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
	var personalInfo []models.PersonalInformation
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
	var surgeryInfo []models.SurgeriesInformation
	if excelCell == nil {
		return surgeryInfo
	}
	for _, row := range excelCell {
		surgeryInfo = append(surgeryInfo, models.SurgeriesInformation{
			SurgeryDate:        ConvertStringToDate(row[0] + "-" + row[1] + "-" + row[2]),
			SurgeryDay:         ConvertSurgeryDayToInt(row[3]),
			SurgeonFirst:       row[7],
			SurgeonSecond:      row[8],
			Resident:           row[9],
			Hospital:           row[10],
			HospitalType:       ConvertHospitalTypeToInt(row[11]),
			SurgeryTime:        ConvertSurgeryTimeToInt(row[12]),
			SurgeryResult:      ConvertSurgeryResultToInt(row[13]),
			CT:                 ConvertImageValidityToInt(row[15]),
			MR:                 ConvertImageValidityToInt(row[16]),
			SurgeryType:        row[17],
			SurgeryArea:        ConvertSurgeryAreaToInt(row[18]),
			OperatorFirst:      row[19],
			OperatorSecond:     row[20],
			StartTime:          ConvertStringToTimestamp(row[34]),
			StopTime:           ConvertStringToTimestamp(row[35]),
			EnterTime:          ConvertStringToTimestamp(row[36]),
			ExitTime:           ConvertStringToTimestamp(row[37]),
			PatientEnterTime:   ConvertStringToTimestamp(row[38]),
			HeadFixType:        ConvertHeadFixTypeToInt(row[39]),
			CancellationReason: row[48],
		})
	}
	return surgeryInfo
}

func ConvertSurgeryTimeToInt(s string) int {
	switch s {
	case "صبح ", "صبح":
		return 1
	case "ظهر ", "ظهر":
		return 2
	case "عصر ", "عصر":
		return 3
	}
	return 0
}

func ConvertSurgeryDayToInt(s string) int {
	switch s {
	case "شنبه":
		return 1
	case "یک شنبه":
		return 2
	case "دو شنبه":
		return 3
	case "سه شنبه":
		return 4
	case "چهار شنبه":
		return 5
	case "پنج شنبه":
		return 6
	case "جمعه":
		return 7
	}
	return 0
}

func GetAllFinancialInformation() []models.FinancialInformation {
	var financialInfo []models.FinancialInformation
	if excelCell == nil {
		return financialInfo
	}
	for _, row := range excelCell {
		receiptNumber, _ := strconv.Atoi(row[49])
		financialInfo = append(financialInfo, models.FinancialInformation{
			PaymentStatus:      ConvertPaymentStatusToInt(row[5]),
			DateOfFirstContact: ConvertStringToDate(row[21]),
			FirstCaller:        row[22],
			DateOfPayment:      ConvertStringToDate(row[25]),
			LastFourDigitsCard: row[26],
			CashAmount:         row[27],
			Bank:               row[28],
			DiscountPercent:    ConvertDiscountPercentToFloat64(row[29]),
			ReasonForDiscount:  row[30],
			CreditAmount:       row[31],
			TypeOfInsurance:    row[32],
			FinancialVerifier:  row[33],
			ReceiptNumber:      receiptNumber,
			ReceiptDate:        ConvertStringToDate(row[50]),
			ReceiptReceiver:    row[51],
		})
	}
	return financialInfo
}

func ConvertPaymentStatusToInt(s string) int {
	switch s {
	case "پرداخت شد", "پرداخت شد ":
		return 1
	case "پرداخت نشد", "پرداخت نشد ":
		return 2
	case "رایگان", "رایگان ":
		return 3
	case "طرح سلامت", "طرح سلامت ":
		return 4
	case "توسط بیمارستان", "توسط بیمارستان ":
		return 5
	}
	return 0
}

func ConvertDiscountPercentToFloat64(s string) float64 {
	if s == "**" {
		return 0
	}
	f, _ := strconv.ParseFloat(s, 32)
	return f
}

func ConvertImageValidityToInt(s string) int {
	switch s {
	case "چک نشده ", "چک نشده":
		return 1
	case "ندارد", "ندارد ":
		return 2
	case "چک شد / تحویل بیمار ", "چک شد / تحویل بیمار":
		return 3
	case "نامعتبر", "نامعتبر ":
		return 4
	}
	return 0
}

func ConvertSurgeryAreaToInt(s string) int {
	switch s {
	case "N":
		return 1
	case "E":
		return 2
	case "E+N":
		return 3
	case "C":
		return 4
	case "S":
		return 5
	case "O":
		return 6
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
		if year >= 0 && year <= 80 {
			year += 1400
		} else if year >= 81 && year <= 99 {
			year += 1300
		}
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
	case "خصوصی", "خصوصی ":
		return 0
	case "دولتی", "دولتی ":
		return 1
	}
	return 2
}

func ConvertSurgeryResultToInt(s string) int {
	switch s {
	case "برگزار شد ", "برگزار شد":
		return 1
	case "کنسل شد ", "کنسل شد":
		return 2
	}
	return 0
}

func ConvertHeadFixTypeToInt(s string) int {
	switch s {
	case "هدبند", "هدبند ":
		return 1
	case "میفیلد", "میفیلد ":
		return 2
	case "هیچکدام", "هیچکدام ":
		return 3
	}
	return 0
}
