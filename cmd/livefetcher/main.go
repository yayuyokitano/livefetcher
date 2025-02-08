package main

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"path/filepath"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/joho/godotenv"
	migrate "github.com/rubenv/sql-migrate"
	"github.com/yayuyokitano/livefetcher/internal/api/endpoints"
	"github.com/yayuyokitano/livefetcher/internal/api/router"
	runner "github.com/yayuyokitano/livefetcher/internal/core"
	"github.com/yayuyokitano/livefetcher/internal/core/logging"
	"github.com/yayuyokitano/livefetcher/internal/core/queries"
	"github.com/yayuyokitano/livefetcher/internal/core/util/datastructures"
	"github.com/yayuyokitano/livefetcher/internal/core/util/templatebuilder"
	i18nloader "github.com/yayuyokitano/livefetcher/internal/i18n"
	"github.com/yayuyokitano/livefetcher/internal/services"
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
		performMigration()
		fmt.Println("Migration complete!")
		return
	case "generatekeys":
		fmt.Println("Generating keys...")
		services.GenerateKeys()
		fmt.Println("Finished generating!")
		return
	case "test":
		err := runner.RunConnectorTest(os.Getenv("CONNECTOR_ID"))
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

func performMigration() {
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
	fs := http.FileServer(http.Dir("./web/static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))
	/*ctx := context.Background()
	for livehouse := range coreconnectors.Connectors {
		fmt.Println("running " + livehouse)
		err := runner.RunConnector(ctx, livehouse)
		fmt.Println(err)
	}*/
	/*
		err := runner.RunConnector("ShinsaibashiKnave")
		fmt.Println(err)

			err = runner.RunConnector("ShimokitazawaShelter")
			fmt.Println(err)
				err = runner.RunConnector("ShinjukuLoft")
				fmt.Println(err)
				err = runner.RunConnector("ShinsaibashiBronze")
				fmt.Println(err)*/

	router.Handle("/login", router.Methods{
		GET: endpoints.ShowLogin,
	})
	router.Handle("/user/{username}", router.Methods{
		GET: endpoints.ShowUser,
	})
	router.Handle("/api/user", router.Methods{
		PATCH: endpoints.PatchUser,
	})

	router.Handle("/api/dailylives/{year}/{month}/{day}", router.Methods{
		GET: endpoints.GetDailyLivesJSON,
	})
	router.Handle("/api/changepassword", router.Methods{
		POST: endpoints.ChangePassword,
	})
	router.Handle("/api/savedsearch", router.Methods{
		POST: endpoints.PostSavedSearch,
	})
	router.Handle("/list/{id}", router.Methods{
		GET: endpoints.ShowLiveList,
	})
	router.Handle("/livelistlive/{id}", router.Methods{
		DELETE: endpoints.DeleteLiveListLive,
	})
	router.Handle("/user/notifications", router.Methods{
		GET: endpoints.ListUserNotifications,
	})
	router.Handle("/notification/{id}", router.Methods{
		GET: endpoints.ShowNotification,
	})
	router.Handle("/calendarEvent/{id}", router.Methods{
		POST:   endpoints.AddToCalendar,
		DELETE: endpoints.RemoveFromCalendar,
	})
	router.Handle("/api/lives", router.Methods{
		GET: endpoints.GetLives,
	})
	router.Handle("/api/addToList", router.Methods{
		POST: endpoints.AddToLiveList,
	})
	router.Handle("/search", router.Methods{
		GET: endpoints.GetLives,
	})
	router.Handle("/favorites", router.Methods{
		GET: endpoints.GetFavoriteLives,
	})
	router.Handle("/modal/livelist", router.Methods{
		GET: endpoints.GetLiveLiveListModal,
	})
	router.Handle("/api/login", router.Methods{
		POST: endpoints.ExecuteLogin,
	})
	router.Handle("/api/register", router.Methods{
		POST: endpoints.Register,
	})
	router.Handle("/api/logout", router.Methods{
		POST: endpoints.Logout,
	})
	router.Handle("/api/favorite", router.Methods{
		POST: endpoints.Favorite,
	})
	router.Handle("/api/unfavorite", router.Methods{
		POST: endpoints.Unfavorite,
	})
	router.Handle("/settings", router.Methods{
		GET: endpoints.ShowSettings,
	})
	router.Handle("/authorizeGoogleCalendar", router.Methods{
		GET: endpoints.AuthorizeGoogleCalendar,
	})
	router.Handle("/api/updateTestLive", router.Methods{
		GET: func(au datastructures.AuthUser, w1 io.Writer, r *http.Request, w2 http.ResponseWriter) *logging.StatusError {
			err := runner.RunConnector(context.Background(), "ShimokitazawaTest")
			if err != nil {
				return logging.SE(http.StatusInternalServerError, err)
			}
			w1.Write([]byte("Successfully updated test live"))
			return nil
		},
	})
	router.Handle("/", router.Methods{
		GET: serveTemplate,
	})
	fmt.Println("Listening on port 9999")
	http.ListenAndServe(":9999", nil)
}

func serveTemplate(user datastructures.AuthUser, w io.Writer, r *http.Request, _ http.ResponseWriter) *logging.StatusError {
	lp := filepath.Join("web", "template", "layout.gohtml")
	fp := filepath.Join("web", "template", filepath.Clean(r.URL.Path), "index.gohtml")

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

	tmpl, err := templatebuilder.Build(w, r, user, nil, lp, fp)
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
