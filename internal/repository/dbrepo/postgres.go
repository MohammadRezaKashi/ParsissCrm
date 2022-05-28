package dbrepo

import (
	"ParsissCrm/internal/models"
	"context"
	"time"
)

func (m *postgresDBRepo) GetAllReports() ([]models.Report, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var reports []models.Report

	query := `
	SELECT * FROM public.reports
	ORDER BY id ASC
	`

	rows, err := m.DB.QueryContext(ctx, query)
	if err != nil {
		return reports, err
	}
	defer rows.Close()

	for rows.Next() {
		var i models.Report
		err := rows.Scan(
			&i.ID,
			&i.Date,
			&i.Acccess_level,
			&i.Created_at,
			&i.Updated_at,
		)

		if err != nil {
			return reports, err
		}
		reports = append(reports, i)
	}

	if err = rows.Err(); err != nil {
		return reports, err
	}

	return reports, nil
}
