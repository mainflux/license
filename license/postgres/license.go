package postgres

import (
	"context"
	"database/sql"
	"encoding/json"
	"time"

	"github.com/lib/pq"
	"github.com/mainflux/license/license"
)

var _ license.Repository = (*licenseRepository)(nil)

const (
	errDuplicate  = "unique_violation"
	errInvalid    = "invalid_text_representation"
	errTruncation = "string_data_right_truncation"
)

type licenseRepository struct {
	db Database
}

// New instantiates a PostgreSQL implementation of license
// repository.
func New(db Database) license.Repository {
	return &licenseRepository{
		db: db,
	}
}

func (repo licenseRepository) Save(ctx context.Context, l license.License) (string, error) {
	q := `INSERT INTO licenses (id, issuer, device_id, created_at, expires_at, updated_at, updated_by, services, plan)
          VALUES (:id, :issuer, :device_id, :created_at, :expires_at, :updated_at, :updated_by, :services, :plan)`

	dbl, err := toDBLicense(l)
	if err != nil {
		return "", err
	}

	if _, err := repo.db.NamedExecContext(ctx, q, dbl); err != nil {

		pqErr, ok := err.(*pq.Error)
		if ok {
			if pqErr.Code.Name() == errDuplicate {
				return "", license.ErrConflict
			}
		}

		return "", err
	}

	return dbl.ID, nil
}

func (repo licenseRepository) Retrieve(ctx context.Context, issuer, id string) (license.License, error) {
	q := `SELECT id, issuer, device_id, created_at, expires_at, updated_at, updated_by, services, plan FROM licenses WHERE issuer = $1 AND id = $2`
	dbl := dbLicense{
		ID:     id,
		Issuer: issuer,
	}
	if err := repo.db.QueryRowxContext(ctx, q, issuer, id).StructScan(&dbl); err != nil {
		pqErr, ok := err.(*pq.Error)
		if err == sql.ErrNoRows || ok && errInvalid == pqErr.Code.Name() {
			return license.License{}, license.ErrNotFound
		}

		return license.License{}, err
	}

	return toLicense(dbl)
}

func (repo licenseRepository) Update(ctx context.Context, l license.License) error {
	q := `UPDATE licenses SET plan = :plan, services = :services, updated_at = :updated_at, expires_at = :expires_at
		  WHERE issuer = :issuer AND id = :id;`
	dbl, err := toDBLicense(l)
	if err != nil {
		return err
	}

	res, err := repo.db.NamedExecContext(ctx, q, dbl)
	if err != nil {
		pqErr, ok := err.(*pq.Error)
		if ok {
			switch pqErr.Code.Name() {
			case errInvalid, errTruncation:
				return license.ErrMalformedEntity
			}
		}

		return err
	}

	cnt, err := res.RowsAffected()
	if err != nil {
		return err
	}

	if cnt == 0 {
		return license.ErrNotFound
	}

	return nil
}

func (repo licenseRepository) Remove(ctx context.Context, owner, id string) error {
	q := `DELETE FROM licenses WHERE issuer = $1 AND id = $2`

	if _, err := repo.db.ExecContext(ctx, q, owner, id); err != nil {
		return err
	}

	return nil
}

type dbLicense struct {
	ID        string         `db:"id"`
	Issuer    string         `db:"issuer"`
	DeviceID  string         `db:"device_id"`
	Active    bool           `db:"active"`
	CreatedAt time.Time      `db:"created_at"`
	ExpiresAt time.Time      `db:"expires_at"`
	UpdatedAt time.Time      `db:"updated_at"`
	UpdatedBy string         `db:"updated_by"`
	Services  pq.StringArray `db:"services"`
	Plan      []byte         `db:"plan"`
}

func toDBLicense(l license.License) (dbLicense, error) {
	data := []byte("{}")
	if len(l.Plan) > 0 {
		b, err := json.Marshal(l.Plan)
		if err != nil {
			return dbLicense{}, err
		}
		data = b
	}

	return dbLicense{
		ID:        l.ID,
		Issuer:    l.Issuer,
		DeviceID:  l.DeviceID,
		Active:    l.Active,
		CreatedAt: l.CreatedAt,
		ExpiresAt: l.ExpiresAt,
		UpdatedAt: l.UpdatedAt,
		UpdatedBy: l.UpdatedBy,
		Services:  l.Services,
		Plan:      data,
	}, nil
}

func toLicense(l dbLicense) (license.License, error) {
	var plan map[string]interface{}
	if err := json.Unmarshal([]byte(l.Plan), &plan); err != nil {
		return license.License{}, err
	}
	return license.License{
		ID:        l.ID,
		Issuer:    l.Issuer,
		DeviceID:  l.DeviceID,
		Active:    l.Active,
		CreatedAt: l.CreatedAt,
		ExpiresAt: l.ExpiresAt,
		UpdatedAt: l.UpdatedAt,
		UpdatedBy: l.UpdatedBy,
		Services:  l.Services,
		Plan:      plan,
	}, nil
}
