package repository

import "ParsissCrm/internal/models"

type DatabaseRepo interface {
	GetAllReports() ([]models.Report, error)
}
