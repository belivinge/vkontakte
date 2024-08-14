package sqlite

import (
	"database/sql"
	"testing"
	"time"
)

func setupTestDB(t *testing.T) (*sql.DB, func()) {
	db, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		t.Fatalf("failed to open database: %v", err)
	}

	_, err = db.Exec(`CREATE TABLE categories (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		name TEXT NOT NULL,
		created DATETIME NOT NULL
	)`)
	if err != nil {
		t.Fatalf("failed to create table: %v", err)
	}

	_, err = db.Exec(`INSERT INTO categories (name, created) VALUES
		('Tech', ?),
		('Lifestyle', ?)`,
		time.Now().Format(time.RFC3339),
		time.Now().Format(time.RFC3339),
	)
	if err != nil {
		t.Fatalf("failed to seed table: %v", err)
	}

	teardown := func() {
		db.Close()
	}

	return db, teardown
}

func TestCategoryModel_Latest(t *testing.T) {
	db, teardown := setupTestDB(t)
	defer teardown()

	model := &CategoryModel{DB: db}

	categories, err := model.Latest()
	if err != nil {
		t.Fatalf("Latest() returned an error: %v", err)
	}

	expectedNumCategories := 2
	if len(categories) != expectedNumCategories {
		t.Errorf("Expected %d categories, got %d", expectedNumCategories, len(categories))
	}

	for _, cat := range categories {
		if cat.ID == 0 {
			t.Error("Category ID should not be 0")
		}
		if cat.Name == "" {
			t.Error("Category Name should not be empty")
		}
		if cat.Created.IsZero() {
			t.Error("Category Created should not be zero")
		}
	}
}
