package device

import (
	"base/domain/entities"
	"context"
	"database/sql"
	"log"
)

type repository struct {
	db *sql.DB
}

func (r repository) Create(ctx context.Context, device entities.Device) error {
	//TODO implement me
	panic("implement me")
}

func (r repository) Update(ctx context.Context, device entities.Device) error {
	//TODO implement me
	panic("implement me")
}

func (r repository) Delete(ctx context.Context, device entities.Device) error {
	//TODO implement me
	panic("implement me")
}

func (r repository) GetAll(ctx context.Context) ([]entities.Device, error) {
	devices := make([]entities.Device, 0)

	query := `
	SELECT d.id,
	       d.name,
	       d.id_orientation,
	       d.status_code, 
	       d.created_at, 
	       d.modified_at
	FROM device as d
	`

	stmt, err := r.db.PrepareContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	result, err := stmt.QueryContext(ctx)
	if err != nil {
		log.Printf("Error in [QueryContext]: %v", err)
		return nil, err
	}
	defer result.Close()

	for result.Next() {
		var dev entities.Device

		err = result.Scan(
			&dev.Id,
			&dev.Name,
			&dev.Orientation,
			&dev.StatusCode,
			&dev.CreatedAt,
			&dev.ModifiedAt,
		)
		if err != nil {
			log.Printf("Error in [Scan]: %v", err)
			return nil, err
		}

		devices = append(devices, dev)
	}

	return devices, nil
}

func (r repository) GetById(ctx context.Context, id int64) (entities.Device, error) {
	query := `
	SELECT d.id,
	       d.name,
	       d.id_orientation,
	       d.status_code, 
	       d.created_at, 
	       d.modified_at,
	       d.id_resolution,
	       r.width,
	       r.height
	FROM device as d
		INNER JOIN resolution r on d.id_resolution = r.id
	WHERE d.id = ?
	`
	var dev entities.Device
	var res entities.Resolution

	err := r.db.QueryRowContext(
		ctx,
		query,
		id,
	).Scan(
		&dev.Id,
		&dev.Name,
		&dev.Orientation,
		&dev.StatusCode,
		&dev.CreatedAt,
		&dev.ModifiedAt,
		&res.Id,
		&res.Width,
		&res.Height,
	)

	if err != nil {
		return dev, err
	}

	dev.Resolution = &res

	return dev, nil
}

func NewDeviceRepository(db *sql.DB) Repository {
	return &repository{
		db: db,
	}
}
