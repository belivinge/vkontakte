package mock

import (
	"forum/models"
	"time"
)

var mockPost = &models.Post{
	ID:      1,
	Title:   "Fake title",
	Content: "Fake Content",
	Created: time.Now(),
}

type PostModel struct{}

func (m *PostModel) Insert(title, content, expires string) (int, error) {
	return 2, nil
}

func (m *PostModel) Get(id int) (*models.Post, error) {
	switch id {
	case 1:
		return mockPost, nil
	default:
		return nil, models.ErrNoRecord
	}
}

func (m *PostModel) Latest() ([]*models.Post, error) {
	return []*models.Post{mockPost}, nil
}
