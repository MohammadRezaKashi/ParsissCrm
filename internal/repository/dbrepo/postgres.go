package dbrepo

import (
	"ParsissCrm/internal/models"
	"context"
	"log"
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

func (m *postgresDBRepo) AddReport(report models.Report) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := `
	INSERT INTO public.reports(
		id, date, access_level, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5);
	`

	_, err := m.DB.ExecContext(ctx, query, report.ID, report.Date.AddDate(0, 0, 1), report.Acccess_level, time.Now(), time.Now())
	if err != nil {
		log.Println(err)
		return err
	}

	return nil
}
