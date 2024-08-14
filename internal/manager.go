package internal

import (
	"crypto/tls"
	"database/sql"
	"flag"
	"forum/models"
	sqlite "forum/sqlite"
	"html/template"
	"log"
	"net/http"
	"net/url"
	"os"
	"sync"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

type contextKey string

var contextKeyUser = contextKey("user")

type Config struct {
	Addr      string
	StaticDir string
}

type application struct {
	errorLog     *log.Logger
	infoLog      *log.Logger
	sessionStore map[string]map[string]interface{}
	sessionMutex sync.Mutex
	activeSessions map[int]string
	posts        interface {
		Insert(string, string, string, string, string) (int, error)
		Get(int) (*models.Post, error)
		Latest() ([]*models.Post, error)
		Filter(url.Values, func(int, string, string) (bool, error), func(int, []string, int) (bool, error)) ([]*models.Post, error)
		UpdateReactions(int, func(int) (int, error), func(int) (int, error)) error
		Paginate([]*models.Post, int, int) ([]*models.Post, int, error)
	}
	templatecache map[string]*template.Template
	users         interface {
		Insert(string, string, string, string) error
		Authenticate(string, string) (int, error)
		Get(int) (*models.User, error)
	}
	comments interface {
		Insert(string, string, string) error
		Latest(int) ([]*models.Comment, error)
		UpdateReactions(int, func(int) (int, error), func(int) (int, error)) error
	}
	categories interface {
		Latest() ([]*models.Category, error)
	}
	post_category interface {
		Insert(string, []string) error
		Get(int) ([]string, error)
		FilterByCategories(int, []string, int) (bool, error)
	}
	post_reactions interface {
		Insert(string, string, string) (int, error)
		Get(string, string) (string, error)
		Delete(string, string) error
		Likes(int) (int, error)
		Dislikes(int) (int, error)
		FilterByLiked(int, string, string) (bool, error)
	}
	comment_reactions interface {
		Insert(string, string, string) (int, error)
		Get(string, string) (string, error)
		Delete(string, string) error
		Likes(int) (int, error)
		Dislikes(int) (int, error)
	}
}

func CreateServer(infoLog *log.Logger, errorLog *log.Logger, cfg *Config) *http.Server {
	dsn := flag.String("dsn", "db/forum.db?parseTime=true", "MySQL database")
	f, err := os.OpenFile("/tmp/info.log", os.O_RDWR|os.O_CREATE, 0o666)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()
	db, err := openDB(*dsn)
	if err != nil {
		errorLog.Fatal(err)
	}
	templateCache, err := newTemplateCache("./ui/html/")
	if err != nil {
		errorLog.Fatal(err)
	}
	app := &application{
		errorLog:          errorLog,
		infoLog:           infoLog,
		sessionStore:      make(map[string]map[string]interface{}),
		posts:             &sqlite.PostModel{DB: db},
		templatecache:     templateCache,
		users:             &sqlite.UserModel{DB: db},
		comments:          &sqlite.CommentModel{DB: db},
		categories:        &sqlite.CategoryModel{DB: db},
		post_category:     &sqlite.PostCategoryModel{DB: db},
		post_reactions:    &sqlite.PostReactionModel{DB: db},
		comment_reactions: &sqlite.CommentReactionModel{DB: db},
	}
	tlsConfig := &tls.Config{
		PreferServerCipherSuites: true,
		CurvePreferences: []tls.CurveID{
			tls.X25519,
			tls.CurveP256,
		},
	}
	srv := &http.Server{
		Addr:         cfg.Addr,
		ErrorLog:     errorLog,
		Handler:      app.routes(),
		TLSConfig:    tlsConfig,
		IdleTimeout:  time.Minute,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
	}
	return srv
}

func openDB(dsn string) (*sql.DB, error) {
	db, err := sql.Open("sqlite3", dsn)
	if err != nil {
		return nil, err
	}
	if err = db.Ping(); err != nil {
		return nil, err
	}
	return db, nil
}
