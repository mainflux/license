package postgres

import (
	"context"
	"database/sql"
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
	q := `INSERT INTO license (id, owner, active, created, duration, expires, metadata, plan)
          VALUES (:id, :owner, :active, :created, :duration, :expires, :metadata, :plan)`

	dbl := toDBLicense(l)
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

func (repo licenseRepository) Retrieve(ctx context.Context, owner, id string) (license.License, error) {
	q := `SELECT id, active, created, duration, expires, metadata, plan FROM keys WHERE owner = $1 AND id = $2`
	dbl := dbLicense{
		ID:    id,
		Owner: owner,
	}
	if err := repo.db.QueryRowxContext(ctx, q, owner, id).StructScan(&dbl); err != nil {
		pqErr, ok := err.(*pq.Error)
		if err == sql.ErrNoRows || ok && errInvalid == pqErr.Code.Name() {
			return license.License{}, license.ErrNotFound
		}

		return license.License{}, err
	}

	return toLicense(dbl), nil
}

func (repo licenseRepository) Update(ctx context.Context, l license.License) error {
	q := `UPDATE license SET plan = :plan, metadata = :metadata WHERE owner = :owner AND id = :id;`
	dbl := toDBLicense(l)

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
	q := `DELETE FROM license WHERE owner = $1 AND id = $2`

	if _, err := repo.db.ExecContext(ctx, q, owner, id); err != nil {
		return err
	}

	return nil
}

type dbLicense struct {
	ID       string                 `db:"id"`
	Owner    string                 `db:"owner"`
	Active   bool                   `db:"active"`
	Created  time.Time              `db:"created"`
	Duration *uint                  `db:"duration"`
	Expires  *time.Time             `db:"expires"`
	Metadata map[string]interface{} `db:"metadata"`
	Plan     map[string]interface{} `db:"plan"`
}

func toDBLicense(l license.License) dbLicense {
	return dbLicense{
		ID:       l.ID,
		Owner:    l.Owner,
		Active:   l.Active,
		Created:  l.Created,
		Duration: l.Duration,
		Expires:  l.Expires,
		Metadata: l.Metadata,
		Plan:     l.Plan,
	}
}

func toLicense(l dbLicense) license.License {
	return license.License{
		ID:       l.ID,
		Owner:    l.Owner,
		Active:   l.Active,
		Created:  l.Created,
		Duration: l.Duration,
		Expires:  l.Expires,
		Metadata: l.Metadata,
		Plan:     l.Plan,
	}
}
