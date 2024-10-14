// file: internal/service/vector_data_service.go

package service

import (
	"canvas-api/internal/models"

	"github.com/gocql/gocql"
)

type VectorDataService struct {
	session *gocql.Session
}

func NewVectorDataService(session *gocql.Session) *VectorDataService {
	return &VectorDataService{session: session}
}

func (s *VectorDataService) CreateVectorData(vectorData *models.VectorData) error {
	query := `
		INSERT INTO vector_data (canvas_id, vector_id, type, x, y)
		VALUES (?, ?, ?, ?, ?)
	`
	return s.session.Query(query,
		vectorData.CanvasID,
		vectorData.VectorID,
		vectorData.Type,
		vectorData.X,
		vectorData.Y).Exec()
}

func (s *VectorDataService) GetVectorDataByCanvasID(canvasID gocql.UUID) ([]*models.VectorData, error) {
	query := `
		SELECT canvas_id, vector_id, type, x, y
		FROM vector_data
		WHERE canvas_id = ?
	`
	iter := s.session.Query(query, canvasID).Iter()
	var vectorDataList []*models.VectorData
	var vd models.VectorData
	for iter.Scan(&vd.CanvasID, &vd.VectorID, &vd.Type, &vd.X, &vd.Y) {
		vectorDataList = append(vectorDataList, &vd)
	}
	if err := iter.Close(); err != nil {
		return nil, err
	}
	return vectorDataList, nil
}

func (s *VectorDataService) UpdateVectorData(vectorData *models.VectorData) error {
	query := `
		UPDATE vector_data
		SET type = ?, x = ?, y = ?
		WHERE canvas_id = ? AND vector_id = ?
	`
	return s.session.Query(query,
		vectorData.Type,
		vectorData.X,
		vectorData.Y,
		vectorData.CanvasID,
		vectorData.VectorID).Exec()
}

func (s *VectorDataService) DeleteVectorData(canvasID gocql.UUID, vectorID gocql.UUID) error {
	query := `
		DELETE FROM vector_data
		WHERE canvas_id = ? AND vector_id = ?
	`
	return s.session.Query(query, canvasID, vectorID).Exec()
}

func (s *VectorDataService) DeleteVectorDataByCanvasID(canvasID gocql.UUID) error {
	query := `DELETE FROM vector_data WHERE canvas_id = ?`
	return s.session.Query(query, canvasID).Exec()
}
