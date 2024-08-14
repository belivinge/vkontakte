package models

import (
	"errors"
	"testing"
	"time"
)

func TestConstants(t *testing.T) {
	tests := []struct {
		name     string
		expected error
		actual   error
	}{
		{"ErrNoRecord", ErrNoRecord, errors.New("models: no matching record found")},
		{"ErrInvalidCredentials", ErrInvalidCredentials, errors.New("models: invalid credentials")},
		{"ErrDuplicateName", ErrDuplicateName, errors.New("UNIQUE constraint failed: users.name")},
		{"ErrDuplicateEmail", ErrDuplicateEmail, errors.New("UNIQUE constraint failed: users.email")},
	}

	for _, test := range tests {
		if test.expected.Error() != test.actual.Error() {
			t.Errorf("Expected %s to be %v, got %v", test.name, test.expected, test.actual)
		}
	}
}

func TestModelDefinitions(t *testing.T) {
	// Create instances of each model with sample data
	post := Post{
		ID:       1,
		UserID:   "user123",
		Username: "john_doe",
		Title:    "Sample Post",
		Content:  "This is a sample post content",
		Created:  time.Now(),
		Likes:    10,
		Dislikes: 2,
	}

	user := User{
		ID:             1,
		Name:           "John Doe",
		Phone:          "123-456-7890",
		   Email:          "john.doe@example.com",
		HashedPassword: []byte("hashedpassword"),
		Created:        time.Now(),
	}

	comment := Comment{
		ID:       1,
		UserID:   "user123",
		PostID:   1,
		Username: "john_doe",
		Content:  "This is a sample comment",
		Created:  time.Now(),
		Likes:    5,
		Dislikes: 1,
	}

	category := Category{
		ID:      1,
		Name:    "Technology",
		Created: time.Now(),
	}

	postCategory := PostCategory{
		PostID:       1,
		CategoryName: "Technology",
		Created:      time.Now(),
	}

	postReaction := PostReaction{
		ID:      1,
		PostID:  1,
		UserID:  1,
		IsLike:  1,
		Created: time.Now(),
	}

	commentReaction := CommentReaction{
		ID:        1,
		CommentID: 1,
		UserID:    1,
		IsLike:    1,
		Created:   time.Now(),
	}

	// Validate that these instances can be created and have correct data
	models := []interface{}{
		post,
		user,
		comment,
		category,
		postCategory,
		postReaction,
		commentReaction,
	}

	for _, model := range models {
		switch m := model.(type) {
		case Post:
			if m.ID == 0 || m.UserID == "" || m.Username == "" || m.Title == "" || m.Content == "" {
				t.Errorf("Post model has invalid data: %+v", m)
			}
		case User:
			if m.ID == 0 || m.Name == "" || m.Phone == "" || m.Email == "" {
				t.Errorf("User model has invalid data: %+v", m)
			}
		case Comment:
			if m.ID == 0 || m.UserID == "" || m.PostID == 0 || m.Username == "" || m.Content == "" {
				t.Errorf("Comment model has invalid data: %+v", m)
			}
		case Category:
			if m.ID == 0 || m.Name == "" {
				t.Errorf("Category model has invalid data: %+v", m)
			}
		case PostCategory:
			if m.PostID == 0 || m.CategoryName == "" {
				t.Errorf("PostCategory model has invalid data: %+v", m)
			}
		case PostReaction:
			if m.ID == 0 || m.PostID == 0 || m.UserID == 0 {
				t.Errorf("PostReaction model has invalid data: %+v", m)
			}
		case CommentReaction:
			if m.ID == 0 || m.CommentID == 0 || m.UserID == 0 {
				t.Errorf("CommentReaction model has invalid data: %+v", m)
			}
		default:
			t.Errorf("Unknown model type: %T", m)
		}
	}
}
