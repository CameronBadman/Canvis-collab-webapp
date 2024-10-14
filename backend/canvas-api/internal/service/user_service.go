// file: internal/service/user_service.go

package service

import (
	"canvas-api/internal/models"

	"github.com/gocql/gocql"
)

type UserService struct {
	session *gocql.Session
}

func NewUserService(session *gocql.Session) *UserService {
	return &UserService{session: session}
}

func (s *UserService) CreateUser(user *models.User) error {
	query := `
		INSERT INTO users (firebase_uid, created_at, id)
		VALUES (?, ?, ?)
	`
	return s.session.Query(query, user.FirebaseUID, user.CreatedAt, user.ID).Exec()
}

func (s *UserService) GetUserByFirebaseUID(firebaseUID string) (*models.User, error) {
	var user models.User
	query := `
		SELECT firebase_uid, created_at, id
		FROM users
		WHERE firebase_uid = ?
	`
	err := s.session.Query(query, firebaseUID).Scan(&user.FirebaseUID, &user.CreatedAt, &user.ID)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (s *UserService) DeleteUserByFirebaseUID(firebaseUID string) error {
	query := `DELETE FROM users WHERE firebase_uid = ?`
	return s.session.Query(query, firebaseUID).Exec()
}

func (s *UserService) UpdateUser(user *models.User) error {
	query := `
		UPDATE users
		SET created_at = ?
		WHERE firebase_uid = ?
	`
	return s.session.Query(query, user.CreatedAt, user.FirebaseUID).Exec()
}
