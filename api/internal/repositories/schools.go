package repositories

import (
	"context"

	"github.com/XDoubleU/essentia/pkg/database"
	"github.com/XDoubleU/essentia/pkg/database/postgres"

	"check-in/api/internal/dtos"
	"check-in/api/internal/models"
)

type SchoolRepository struct {
	db postgres.DB
}

func (repo SchoolRepository) GetTotalCount(ctx context.Context) (*int64, error) {
	query := `
		SELECT COUNT(*)
		FROM schools
	`

	var total *int64

	err := repo.db.QueryRow(ctx, query).Scan(&total)
	if err != nil {
		return nil, postgres.PgxErrorToHTTPError(err)
	}

	return total, nil
}

func (repo SchoolRepository) GetAll(ctx context.Context) ([]*models.School, error) {
	query := `
		SELECT id, name
		FROM schools
		ORDER BY name ASC
	`

	rows, err := repo.db.Query(ctx, query)
	if err != nil {
		return nil, postgres.PgxErrorToHTTPError(err)
	}

	schools := []*models.School{}

	for rows.Next() {
		var school models.School

		err = rows.Scan(
			&school.ID,
			&school.Name,
		)

		if err != nil {
			return nil, postgres.PgxErrorToHTTPError(err)
		}

		schools = append(schools, &school)
	}

	if err = rows.Err(); err != nil {
		return nil, postgres.PgxErrorToHTTPError(err)
	}

	return schools, nil
}

func (repo SchoolRepository) GetAllSortedByLocation(
	ctx context.Context,
	locationID string,
) ([]*models.School, error) {
	query := `
		SELECT id, name
		FROM schools
		ORDER BY
			CASE
				WHEN read_only = true THEN -1
				ELSE (
					SELECT COUNT(*)
					FROM check_ins
					WHERE check_ins.location_id = $1
					AND check_ins.school_id = schools.id
				)
			END
		DESC, name ASC
	`

	rows, err := repo.db.Query(ctx, query, locationID)
	if err != nil {
		return nil, postgres.PgxErrorToHTTPError(err)
	}

	schools := []*models.School{}

	for rows.Next() {
		var school models.School

		err = rows.Scan(
			&school.ID,
			&school.Name,
		)

		if err != nil {
			return nil, postgres.PgxErrorToHTTPError(err)
		}

		schools = append(schools, &school)
	}

	if err = rows.Err(); err != nil {
		return nil, postgres.PgxErrorToHTTPError(err)
	}

	return schools, nil
}

func (repo SchoolRepository) GetAllPaginated(
	ctx context.Context,
	limit int64,
	offset int64,
) ([]*models.School, error) {
	query := `
		SELECT id, name, read_only
		FROM schools
		ORDER BY name ASC
		LIMIT $1 OFFSET $2
	`

	rows, err := repo.db.Query(ctx, query, limit, offset)
	if err != nil {
		return nil, postgres.PgxErrorToHTTPError(err)
	}

	schools := []*models.School{}

	for rows.Next() {
		var school models.School

		err = rows.Scan(
			&school.ID,
			&school.Name,
			&school.ReadOnly,
		)

		if err != nil {
			return nil, postgres.PgxErrorToHTTPError(err)
		}

		schools = append(schools, &school)
	}

	if err = rows.Err(); err != nil {
		return nil, postgres.PgxErrorToHTTPError(err)
	}

	return schools, nil
}

func (repo SchoolRepository) GetByID(
	ctx context.Context,
	id int64,
) (*models.School, error) {
	query := `
		SELECT name, read_only
		FROM schools
		WHERE id = $1
	`

	//nolint:exhaustruct //other fields are optional
	school := models.School{
		ID: id,
	}

	err := repo.db.QueryRow(
		ctx,
		query,
		id).Scan(&school.Name, &school.ReadOnly)

	if err != nil {
		return nil, postgres.PgxErrorToHTTPError(err)
	}

	return &school, nil
}

func (repo SchoolRepository) GetByName(
	ctx context.Context,
	name string,
) (*models.School, error) {
	query := `
		SELECT id, read_only
		FROM schools
		WHERE name = $1
	`

	//nolint:exhaustruct //other fields are optional
	school := models.School{
		Name: name,
	}

	err := repo.db.QueryRow(
		ctx,
		query,
		name).Scan(&school.ID, &school.ReadOnly)

	if err != nil {
		return nil, postgres.PgxErrorToHTTPError(err)
	}

	return &school, nil
}

func (repo SchoolRepository) GetByIDWithoutReadOnly(
	ctx context.Context,
	id int64,
) (*models.School, error) {
	query := `
		SELECT name
		FROM schools
		WHERE id = $1 AND read_only = false
	`

	//nolint:exhaustruct //other fields are optional
	school := models.School{
		ID:       id,
		ReadOnly: false,
	}

	err := repo.db.QueryRow(
		ctx,
		query,
		id).Scan(&school.Name)

	if err != nil {
		return nil, postgres.PgxErrorToHTTPError(err)
	}

	return &school, nil
}

func (repo SchoolRepository) Create(
	ctx context.Context,
	name string,
) (*models.School, error) {
	query := `
		INSERT INTO schools (name)
		VALUES ($1)
		RETURNING id
	`

	//nolint:exhaustruct //other fields are optional
	school := models.School{
		Name: name,
	}

	err := repo.db.QueryRow(ctx, query, name).Scan(&school.ID)

	if err != nil {
		return nil, postgres.PgxErrorToHTTPError(err)
	}

	return &school, nil
}

func (repo SchoolRepository) Update(
	ctx context.Context,
	school models.School,
	schoolDto dtos.SchoolDto,
) (*models.School, error) {
	school.Name = schoolDto.Name

	query := `
		UPDATE schools
		SET name = $2
		WHERE id = $1 AND read_only = false
	`

	result, err := repo.db.Exec(ctx, query, school.ID, school.Name)
	if err != nil {
		return nil, postgres.PgxErrorToHTTPError(err)
	}

	rowsAffected := result.RowsAffected()
	if rowsAffected == 0 {
		return nil, database.ErrResourceNotFound
	}

	return &school, nil
}

func (repo SchoolRepository) Delete(ctx context.Context, id int64) error {
	query := `
		DELETE FROM schools
		WHERE id = $1 AND read_only = false
	`

	result, err := repo.db.Exec(ctx, query, id)
	if err != nil {
		return postgres.PgxErrorToHTTPError(err)
	}

	rowsAffected := result.RowsAffected()
	if rowsAffected == 0 {
		return database.ErrResourceNotFound
	}

	return nil
}
