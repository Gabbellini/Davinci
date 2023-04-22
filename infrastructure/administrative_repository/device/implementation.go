package device

import (
	"base/domain/entities"
	"context"
	"database/sql"
	"errors"
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

func (r repository) Create(ctx context.Context, device entities.Device, userId int64) error {
	query := `
	INSERT INTO device (name, id_resolution, id_orientation, id_user) 
	VALUES (?,?,?,?)
	`

	stmt, err := r.db.PrepareContext(ctx, query)
	if err != nil {
		log.Println("[Update] Error PrepareContext", err)
		return err
	}
	defer stmt.Close()

	_, err = stmt.ExecContext(
		ctx,
		device.Name,
		device.Resolution.Id,
		device.Orientation,
		userId,
	)
	if err != nil {
		log.Println("[Create] Error ExecContext", err)
		return err
	}

	return nil
}

func (r repository) Update(ctx context.Context, device entities.Device, userId int64) error {
	command := `
	UPDATE device
	SET name = ?, 
	    id_resolution = ?, 
	    id_orientation = ?
	WHERE id = ? AND 
	      id_user = ?
	`

	_, err := r.db.ExecContext(
		ctx,
		command,
		device.Name,
		device.Resolution.Id,
		device.Orientation,
		device.Id,
		userId,
	)
	if err != nil {
		log.Println("[Update] Error ExecContext", err)
		return err
	}

	return nil
}

func (r repository) Delete(ctx context.Context, deviceId int64, userId int64) error {
	//language=sql
	command := `
	UPDATE device 
	SET status_code = ?
	WHERE id = ? 
	  AND id_user = ?
	`

	_, err := r.db.ExecContext(ctx, command, entities.StatusDeleted, deviceId, userId)
	if err != nil {
		log.Println("[Delete] Error ExecContext", err)
		return err
	}

	return nil
}

func (r repository) GetAll(ctx context.Context, userId int64) ([]entities.Device, error) {
	query := `
	SELECT d.id,
	       d.name,
	       d.id_orientation,
	       d.status_code,
		   r.width,
		   r.height
	FROM device d
		INNER JOIN resolution r on d.id_resolution = r.id
	WHERE id_user = ? AND 
	      d.status_code = ?
	`
	rows, err := r.db.QueryContext(ctx, query, userId, entities.StatusOk)
	if err != nil {
		log.Println("[GetAll] Error QueryContext", err)
		return nil, err
	}
	defer rows.Close()

	var devices []entities.Device
	for rows.Next() {
		var device entities.Device

		err = rows.Scan(
			&device.Id,
			&device.Name,
			&device.Orientation,
			&device.StatusCode,
			&device.Resolution.Width,
			&device.Resolution.Height,
		)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				return nil, nil
			}

			log.Println("[GetAll] Error Scan", err)
			return nil, err
		}

		devices = append(devices, device)
	}

	return devices, nil
}

func (r repository) GetById(ctx context.Context, deviceId int64, userId int64) (*entities.Device, error) {
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
	WHERE d.id = ? AND 
	      d.id_user = ? AND
	      d.status_code = ?
	`

	var device entities.Device
	err := r.db.QueryRowContext(
		ctx,
		query,
		deviceId,
		userId,
		entities.StatusOk,
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
		log.Println("[GetById] Error Scan", err)
		return nil, err
	}

	return &device, nil
}
