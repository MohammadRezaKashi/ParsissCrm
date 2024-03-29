package repository

import "ParsissCrm/internal/models"

type DatabaseRepo interface {
	AddPersonalInformation(information models.PersonalInformation) (int, error)
	AddSurgeriesInformation(information models.SurgeriesInformation, id int) error
	AddFinancialInformation(information models.FinancialInformation, id int) error
	GetAllPatients() ([]models.PersonalInformation, error)
	GetPatientByID(id int) (models.PersonalInformation, error)
	GetSurgicalInformationByPatientID(id int) ([]models.SurgeriesInformation, error)
	GetFinancialInformationByPatientID(id int) ([]models.FinancialInformation, error)
	PutPersonalInformation(personalInfo models.PersonalInformation) error
	PutSurgeriesInformation(surgeriesInfo models.SurgeriesInformation) error
	PutFinancialInformation(surgeriesInfo models.FinancialInformation) error
	GetDistinctList(tableName string, columnName string) ([]interface{}, error)
	GetFilterData(p interface{}) ([]models.PersonalInformation, error)
}
