package main

import (
	"database/sql"
	"errors"
	"fmt"
	"html/template"
	"io"
	"net/http"
	"net/url"
	"os"
	"path/filepath"

	_ "github.com/jackc/pgx/v4/stdlib"
	"github.com/joho/godotenv"
	migrate "github.com/rubenv/sql-migrate"
	"github.com/yayuyokitano/livefetcher/i18nloader"
	"github.com/yayuyokitano/livefetcher/lib/api/endpoints"
	"github.com/yayuyokitano/livefetcher/lib/api/router"
	runner "github.com/yayuyokitano/livefetcher/lib/core"
	"github.com/yayuyokitano/livefetcher/lib/core/logging"
	"github.com/yayuyokitano/livefetcher/lib/core/queries"
	"github.com/yayuyokitano/livefetcher/lib/services"
)

func main() {
	// load env in non-containerized execution
	if os.Getenv("CONTAINERIZED") != "true" {
		err := godotenv.Load()
		if err != nil {
			panic("Error loading .env file")
		}
	}

	switch os.Args[len(os.Args)-1] {
	case "migrate":
		fmt.Println("Performing migration...")
		performMigration(true)
		fmt.Println("Migration complete!")
		return
	case "test":
		err := runner.RunConnector(os.Args[2])
		fmt.Println(err)
		return
	case "start":
		fmt.Println("Starting server...")
	default:
		fmt.Println("Invalid command")
		return
	}
	err := services.Start()
	defer services.Stop()
	if err != nil {
		panic(err)
	}
	err = i18nloader.Init()
	if err != nil {
		panic(err)
	}

	fmt.Println("Connected to Postgres!")

	go logging.ServeLogs()
	startServer()
}

func performMigration(firstTime bool) {
	migrations := &migrate.FileMigrationSource{
		Dir: "./migrations",
	}

	var domain string
	if os.Getenv("CONTAINERIZED") == "true" {
		domain = "db"
	} else {
		domain = "localhost"
	}

	db, err := sql.Open("pgx", fmt.Sprintf("postgresql://%s:%s@%s:5432/%s", os.Getenv("POSTGRES_USER"), url.QueryEscape(os.Getenv("POSTGRES_PASSWORD")), domain, os.Getenv("POSTGRES_DB")))
	if err != nil {
		panic(err)
	}
	defer db.Close()

	_, err = db.Exec(fmt.Sprintf("CREATE DATABASE %s", os.Getenv("POSTGRES_DB")))
	if err != nil {
		fmt.Println("Failed to create database, probably already exists.")
	}

	_, err = db.Exec(fmt.Sprintf("CREATE USER %s WITH PASSWORD '%s'", os.Getenv("POSTGRES_GRAFANA_USER"), os.Getenv("POSTGRES_GRAFANA_PASSWORD")))
	if err != nil {
		fmt.Println("Failed to create user, probably already exists.")
	}

	_, err = db.Exec(fmt.Sprintf("GRANT pg_read_all_data TO %s", os.Getenv("POSTGRES_GRAFANA_USER")))
	if err != nil {
		fmt.Println("Failed to grant user permissions, probably already exists.")
	}

	n, err := migrate.Exec(db, "postgres", migrations, migrate.Up)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Printf("Applied %d migrations!\n", n)

}

func startServer() {
	fs := http.FileServer(http.Dir("./static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))
	/*for livehouse := range coreconnectors.Connectors {
		err := runner.RunConnector(livehouse)
		fmt.Println(err)
	}*/
	/*err := runner.RunConnector("ShimokitazawaMosaic")
	fmt.Println(err)
	err = runner.RunConnector("ShimokitazawaShelter")
	fmt.Println(err)
	err = runner.RunConnector("ShinjukuLoft")
	fmt.Println(err)
	err = runner.RunConnector("ShinsaibashiBronze")
	fmt.Println(err)*/

	router.Handle("/api/lives", router.Methods{
		GET: endpoints.GetLives,
	})
	router.Handle("/", router.Methods{
		GET: serveTemplate,
	})
	fmt.Println("Listening on port 9999")
	http.ListenAndServe(":9999", nil)
}

func serveTemplate(w io.Writer, r *http.Request) *logging.StatusError {
	lp := filepath.Join("templates", "layout.html")
	fp := filepath.Join("templates", filepath.Clean(r.URL.Path), "index.html")

	// Return a 404 if the template doesn't exist
	info, err := os.Stat(fp)
	if err != nil {
		if os.IsNotExist(err) {
			return logging.SE(http.StatusNotFound, errors.New("404 page not found"))
		}
	}

	// Return a 404 if the request is for a directory
	if info.IsDir() {
		return logging.SE(http.StatusNotFound, errors.New("404 page not found"))
	}

	tmpl, err := template.New("layout").Funcs(template.FuncMap{
		"T":    i18nloader.GetLocalizer(w, r).Localize,
		"Lang": func() string { return i18nloader.GetMainLanguage(w, r) },
	}).ParseFiles(lp, fp)
	if err != nil {
		return logging.SE(http.StatusInternalServerError, err)
	}

	areas, err := queries.GetAllAreas(r.Context())
	if err != nil {
		return logging.SE(http.StatusInternalServerError, err)
	}

	err = tmpl.ExecuteTemplate(w, "layout", areas)
	if err != nil {
		return logging.SE(http.StatusInternalServerError, err)
	}
	return nil
}
