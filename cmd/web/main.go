package main

import (
	"database/sql"
	"flag"
	"log"
	"net/http"
	"os"

	"github.com/BradPreston/snippetbox/internal/models"
	_ "github.com/go-sql-driver/mysql"
)

type config struct {
    addr string
    staticDir string
    dsn string
}

type application struct {
    errorLog *log.Logger
    infoLog *log.Logger
    snippets *models.SnippetModel
}

func main() {
    var cfg config 
    flag.StringVar(&cfg.addr, "addr", ":4000", "HTTP network address")
    flag.StringVar(&cfg.staticDir, "static-dir", "./ui/static", "Path to static assets")
    flag.StringVar(&cfg.dsn, "dsn", "web:pass@/snippetbox?parseTime=true", "MySQL data source name")
    flag.Parse()

    infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
    errorLog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

    db, err := openDB(cfg.dsn)
    if err != nil {
        errorLog.Fatal(err)
    }
    defer db.Close()

    app := &application{
        errorLog: errorLog,
        infoLog: infoLog,
        snippets: &models.SnippetModel{DB: db},
    }

    srv := &http.Server{
        Addr: cfg.addr,
        ErrorLog: errorLog,
        Handler: app.routes(),
    }
	infoLog.Printf("Starting server on %s", cfg.addr)
    err = srv.ListenAndServe()
    errorLog.Fatal(err)
}

func openDB(dsn string) (*sql.DB, error) {
    db, err := sql.Open("mysql", dsn)
    if err != nil {
        return nil, err
    }
    if err = db.Ping(); err != nil {
        return nil, err
    }
    return db, nil
}
