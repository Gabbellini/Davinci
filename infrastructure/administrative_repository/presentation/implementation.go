package presentation

import (
	"context"
	"database/sql"
	"davinci/domain/entities"
	"davinci/settings"
	"encoding/json"
	"log"
)

type repository struct {
	db       *sql.DB
	settings settings.Settings
}

func NewPresentationRepository(settings settings.Settings, db *sql.DB) Repository {
	return &repository{
		db:       db,
		settings: settings,
	}
}

func (r repository) Create(ctx context.Context, presentation entities.Presentation, userId int64) (int64, error) {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		log.Println("[Create] Error BeginTx", err)
		return 0, err
	}

	presentationId, err := r.createPresentation(ctx, tx, presentation, userId)
	if err != nil {
		_ = tx.Rollback()
		log.Println("[Create] Error createPresentation", err)
		return 0, err
	}

	for _, page := range presentation.Pages {
		_, err = r.createPage(ctx, tx, page, presentationId)
		if err != nil {
			_ = tx.Rollback()
			log.Println("[Create] Error createPage", err)
			return 0, err
		}
	}

	err = tx.Commit()
	if err != nil {
		log.Println("[Create] Error Commit", err)
		return 0, err
	}

	return presentationId, nil
}

func (r repository) createPresentation(ctx context.Context, tx *sql.Tx, presentation entities.Presentation, userId int64) (int64, error) {
	command := `
	INSERT INTO presentation (name, id_user, id_resolution)
	VALUES (?,?,?)
	`

	var result sql.Result
	var err error

	if tx != nil {
		result, err = tx.ExecContext(ctx, command, presentation.Name, userId, presentation.ResolutionId)
	} else {
		result, err = r.db.ExecContext(ctx, command, presentation.Name, userId, presentation.ResolutionId)
	}

	if err != nil {
		log.Println("[createPresentation] Error ExecContext", err)
		return 0, err
	}

	presentationId, err := result.LastInsertId()
	if err != nil {
		log.Println("[createPresentation] Error LastInsertId", err)
		return 0, err
	}

	return presentationId, nil
}

func (r repository) createPage(ctx context.Context, tx *sql.Tx, page entities.Page, presentationId int64) (int64, error) {
	command := `
	INSERT INTO page (id_presentation, component, duration)
	VALUES (?,?,?)
	`

	var result sql.Result
	var err error

	b, err := json.Marshal(page.Component)
	if err != nil {
		log.Println("[createPage] Error Marshal", err)
		return 0, err
	}
	componentString := string(b)

	if tx != nil {
		result, err = r.db.ExecContext(ctx, command, presentationId, componentString, page.Duration)
	} else {
		result, err = r.db.ExecContext(ctx, command, presentationId, componentString, page.Duration)
	}
	if err != nil {
		log.Println("[createPage] Error ExecContext", err)
		return 0, err
	}

	pageId, err := result.LastInsertId()
	if err != nil {
		log.Println("[createPage] Error LastInsertId", err)
		return 0, err
	}

	return pageId, nil
}

func (r repository) Update(ctx context.Context, presentationId int64, presentation entities.Presentation, userId int64) error {
	tx, err := r.db.Begin()
	if err != nil {
		log.Println("[Update] error in Begin tx", err)
		return err
	}

	err = r.updatePresentation(ctx, tx, presentationId, presentation, userId)
	if err != nil {
		_ = tx.Rollback()
		log.Println("[Update] Error updatePresentation", err)
		return err
	}

	for _, page := range presentation.Pages {
		err = r.updatePage(ctx, tx, presentationId, page)
		if err != nil {
			_ = tx.Rollback()
			log.Println("[Update] Error updatePage", err)
			return err
		}
	}

	_ = tx.Commit()
	if err != nil {
		log.Println("[Update] Error Commit", err)
		return err
	}

	return nil
}

func (r repository) updatePresentation(ctx context.Context, tx *sql.Tx, presentationId int64, presentation entities.Presentation, userId int64) error {
	command := `
	UPDATE presentation 
	SET name = ?, id_resolution = ?
	WHERE id = ? AND id_user = ?`

	var err error
	if tx != nil {
		_, err = tx.ExecContext(ctx, command, presentation.Name, presentation.ResolutionId, presentationId, userId)
	} else {
		_, err = r.db.ExecContext(ctx, command, presentation.Name, presentation.ResolutionId, presentationId, userId)
	}

	if err != nil {
		log.Println("[updatePresentation] Error ExecContext", err)
		return err
	}

	return nil
}

func (r repository) updatePage(ctx context.Context, tx *sql.Tx, presentationId int64, page entities.Page) error {
	command := `
	UPDATE page 
	SET id_presentation = ? , component = ?, duration = ?
	WHERE id = ?`

	var err error

	b, err := json.Marshal(page.Component)
	if err != nil {
		log.Println("[updatePage] Error Marshal", err)
		return err
	}
	componentString := string(b)

	if tx != nil {
		_, err = r.db.ExecContext(ctx, command, presentationId, componentString, page.Duration, presentationId)
	} else {
		_, err = r.db.ExecContext(ctx, command, presentationId, componentString, page.Duration, presentationId)
	}
	if err != nil {
		log.Println("[updatePage] Error ExecContext", err)
		return err
	}

	return nil
}

func (r repository) Delete(ctx context.Context, presentationId int64, userId int64) error {
	tx, err := r.db.Begin()
	if err != nil {
		log.Println("[Delete] Error Begin", err)
		return err
	}

	err = r.deletePresentation(ctx, tx, presentationId, userId)
	if err != nil {
		log.Println("[Delete] Error deletePresentation", err)
		return err
	}

	err = r.deletePresentationPages(ctx, tx, presentationId)
	if err != nil {
		log.Println("[Delete] Error deletePresentationPages", err)
		return err
	}

	err = tx.Commit()
	if err != nil {
		log.Println("[Delete] Error Commit", err)
		return err
	}

	return nil
}

func (r repository) deletePresentation(ctx context.Context, tx *sql.Tx, presentationId int64, userId int64) error {
	command := `
	UPDATE presentation 
	SET status_code = ?
	WHERE id = ? AND id_user = ?`

	var err error
	if tx != nil {
		_, err = tx.ExecContext(ctx, command, entities.StatusDeleted, presentationId, userId)
	} else {
		_, err = r.db.ExecContext(ctx, command, entities.StatusDeleted, presentationId, userId)
	}
	if err != nil {
		log.Println("[deletePresentation] Error ExecContext", err)
		return err
	}

	return nil
}

func (r repository) deletePresentationPages(ctx context.Context, tx *sql.Tx, presentationId int64) error {
	command := `
	UPDATE page 
	SET status_code = ?
	WHERE id_presentation = ?`

	var err error
	if tx != nil {
		_, err = r.db.ExecContext(ctx, command, presentationId)
	} else {
		_, err = r.db.ExecContext(ctx, command, presentationId)
	}
	if err != nil {
		log.Println("[deletePresentationPages] Error ExecContext", err)
		return err
	}

	return nil
}

func (r repository) GetAll(ctx context.Context, userId int64) ([]entities.Presentation, error) {
	query := `
	SELECT p.id,
	       p.name,
	       p.status_code, 
	       p.created_at, 
	       p.modified_at
	FROM presentation as p
	WHERE id_user = ?
	`

	rows, err := r.db.QueryContext(ctx, query, userId)
	if err != nil {
		log.Println("[GetAll] Error QueryContext", err)
		return nil, err
	}
	defer rows.Close()

	var presentations []entities.Presentation
	for rows.Next() {
		var presentation entities.Presentation
		err = rows.Scan(
			&presentation.Id,
			&presentation.Name,
			&presentation.StatusCode,
			&presentation.CreatedAt,
			&presentation.ModifiedAt,
		)
		if err != nil {
			log.Println("[GetAll] Error Scan", err)
			return nil, err
		}

		presentations = append(presentations, presentation)
	}

	return presentations, nil
}

func (r repository) GetById(ctx context.Context, id int64, userId int64) (*entities.Presentation, error) {
	query := `
	SELECT p.id,
	       p.name,
	       p.status_code, 
	       p.created_at, 
	       p.modified_at,
	       p.id_resolution
	FROM presentation as p
	WHERE p.id = ? AND p.id_user = ?
	`

	queryPages := `
	SELECT p.id, 
	       p.id_presentation, 
	       p.duration, 
	       p.component, 
	       p.status_code, 
	       p.created_at, 
	       p.modified_at
	FROM page p
	WHERE id_presentation = ?
	`

	var presentation entities.Presentation
	err := r.db.QueryRowContext(
		ctx,
		query,
		id,
		userId,
	).Scan(
		&presentation.Id,
		&presentation.Name,
		&presentation.StatusCode,
		&presentation.CreatedAt,
		&presentation.ModifiedAt,
		&presentation.ResolutionId,
	)
	if err != nil {
		log.Println("[GetById] error in QueryRowContext", err)
		return nil, err
	}

	rows, err := r.db.QueryContext(ctx, queryPages, id)
	if err != nil {
		log.Println("[GetById] error in QueryContext", err)
		return nil, err
	}
	defer rows.Close()

	var pages []entities.Page
	for rows.Next() {
		var page entities.Page
		err = rows.Scan(
			&page.Id,
			&page.PresentationId,
			&page.Duration,
			&page.Component,
			&page.StatusCode,
			&page.CreatedAt,
			&page.ModifiedAt,
		)
		if err != nil {
			log.Println("[GetById] error in Scan", err)
			return nil, err
		}

		pages = append(pages, page)
	}
	presentation.Pages = pages

	return &presentation, nil
}

func (r repository) getPresentation(ctx context.Context, presentationId int64) (*entities.Presentation, error) {
	//language=sql
	query := `
	SELECT id, 
	       id_resolution, 
	       name, 
	       status_code, 
	       created_at, 
	       modified_at
	FROM presentation
	WHERE id = ?`

	var presentation entities.Presentation
	err := r.db.QueryRowContext(ctx, query, presentationId).Scan(
		&presentation.Id,
		&presentation.ResolutionId,
		&presentation.Name,
		&presentation.StatusCode,
		&presentation.CreatedAt,
		&presentation.ModifiedAt,
	)
	if err != nil {
		log.Println("[getPresentation] Error Scan", err)
		return nil, err
	}

	return &presentation, nil
}

func (r repository) getPresentationPages(ctx context.Context, presentationId int64) ([]entities.Page, error) {
	//language=sql
	query := `
	SELECT id, 
	       id_presentation, 
	       duration, 
	       component, 
	       status_code, 
	       created_at, 
	       modified_at
	FROM page
	WHERE id_presentation = ?
	`

	rows, err := r.db.QueryContext(ctx, query, presentationId)
	if err != nil {
		log.Println("[getPresentationPages] Error QueryContext", err)
		return nil, err
	}
	defer rows.Close()

	var pages []entities.Page
	for rows.Next() {
		var page entities.Page
		var componentString string
		err = rows.Scan(
			&page.Id,
			&page.PresentationId,
			&page.Duration,
			&componentString,
			&page.StatusCode,
			&page.CreatedAt,
			&page.ModifiedAt,
		)
		if err != nil {
			log.Println("[getPresentationPages] Error Scan", err)
			return nil, err
		}

		err = json.Unmarshal([]byte(componentString), &page.Component)
		if err != nil {
			log.Println("[getPresentationPages] Error Unmarshal", err)
			return nil, err
		}

		pages = append(pages, page)
	}

	return pages, err
}
