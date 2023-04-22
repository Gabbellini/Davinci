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

func NewRepository(db *sql.DB) Repository {
	return &repository{
		db: db,
	}
}

func (r repository) GetByName(ctx context.Context, deviceName string, userId int64) (*entities.Device, error) {
	//language=sql
	query := `
	SELECT d.id,
	       d.name,
	       d.id_orientation,
	       d.status_code, 
	       d.id_resolution,
	       r.width,
	       r.height
	FROM device d
		INNER JOIN resolution r on d.id_resolution = r.id
	WHERE d.name = ? AND d.id_user = ?
	`
	var device entities.Device

	err := r.db.QueryRowContext(
		ctx,
		query,
		deviceName,
		userId,
	).Scan(
		&device.Id,
		&device.Name,
		&device.Orientation,
		&device.StatusCode,
		&device.Resolution.Id,
		&device.Resolution.Width,
		&device.Resolution.Height,
	)

	if err != nil {
		log.Println("[GetByName] Error Scan", err)
		return nil, err
	}

	return &device, nil
}
