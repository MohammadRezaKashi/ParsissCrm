package repository

import "ParsissCrm/internal/models"

type DatabaseRepo interface {
	AddPersonalInformation(information models.PersonalInformation) (int, error)
	AddSurgeriesInformation(information models.SurgeriesInformation, id int) error
}
