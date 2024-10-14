// file: internal/service/canvas_service.go

package service

import (
	"canvas-api/internal/models"

	"github.com/gocql/gocql"
)

type CanvasService struct {
	session *gocql.Session
}

func NewCanvasService(session *gocql.Session) *CanvasService {
	return &CanvasService{session: session}
}

func (s *CanvasService) CreateCanvas(canvas *models.Canvas) error {
	query := `
		INSERT INTO canvases (firebase_uid, id, created_at, name)
		VALUES (?, ?, ?, ?)
	`
	return s.session.Query(query, canvas.FirebaseUID, canvas.ID, canvas.CreatedAt, canvas.Name).Exec()
}

func (s *CanvasService) GetCanvas(firebaseUID string, id gocql.UUID) (*models.Canvas, error) {
	var canvas models.Canvas
	query := `
		SELECT firebase_uid, id, created_at, name
		FROM canvases
		WHERE firebase_uid = ? AND id = ?
	`
	err := s.session.Query(query, firebaseUID, id).Scan(&canvas.FirebaseUID, &canvas.ID, &canvas.CreatedAt, &canvas.Name)
	if err != nil {
		return nil, err
	}
	return &canvas, nil
}

func (s *CanvasService) GetUserCanvases(firebaseUID string) ([]*models.Canvas, error) {
	query := `
		SELECT firebase_uid, id, created_at, name
		FROM canvases
		WHERE firebase_uid = ?
	`
	iter := s.session.Query(query, firebaseUID).Iter()
	var canvases []*models.Canvas
	var canvas models.Canvas
	for iter.Scan(&canvas.FirebaseUID, &canvas.ID, &canvas.CreatedAt, &canvas.Name) {
		canvases = append(canvases, &canvas)
	}
	if err := iter.Close(); err != nil {
		return nil, err
	}
	return canvases, nil
}

func (s *CanvasService) UpdateCanvas(canvas *models.Canvas) error {
	query := `
		UPDATE canvases
		SET name = ?
		WHERE firebase_uid = ? AND id = ?
	`
	return s.session.Query(query, canvas.Name, canvas.FirebaseUID, canvas.ID).Exec()
}

func (s *CanvasService) DeleteCanvas(firebaseUID string, id gocql.UUID) error {
	query := `DELETE FROM canvases WHERE firebase_uid = ? AND id = ?`
	return s.session.Query(query, firebaseUID, id).Exec()
}
