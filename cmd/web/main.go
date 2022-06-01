package main

import (
	"ParsissCrm/internal/config"
	"ParsissCrm/internal/driver"
	"ParsissCrm/internal/handlers"
	"ParsissCrm/internal/helpers"
	"ParsissCrm/internal/render"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/alexedwards/scs/v2"
)

const portNumber = ":8080"

var app config.AppConfig
var session *scs.SessionManager
var infoLog *log.Logger
var errorLog *log.Logger

func main() {
	db, err := Run()
	if err != nil {
		log.Fatal(err)
	}

	defer db.SQL.Close()

	srv := &http.Server{
		Addr:    portNumber,
		Handler: routes(&app),
	}

	err = srv.ListenAndServe()
	log.Fatal(err)
}

func Run() (*driver.DB, error) {
	inProduction := flag.Bool("production", true, "Application is in production")
	useCache := flag.Bool("cache", true, "Use template cache")
	dbName := flag.String("dbname", "", "Database name")
	dbHost := flag.String("dbhost", "localhost", "Database host")
	dbUser := flag.String("dbuser", "", "Database user")
	dbPass := flag.String("dbpass", "", "Database password")
	dbPort := flag.String("dbport", "5432", "Database port")
	dbSSL := flag.String("dbssl", "disable", "Database ssl settings (disable, prefer, require)")

	flag.Parse()

	if *dbName == "" || *dbUser == "" {
		fmt.Println("Missing required flags")
		os.Exit(1)
	}

	app.InProduction = *inProduction
	app.UseCache = *useCache

	infoLog = log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	app.InfoLog = infoLog

	errorLog = log.New(os.Stdout, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)
	app.ErrorLog = errorLog

	session = scs.New()
	session.Lifetime = 24 * time.Hour
	session.Cookie.Persist = true
	session.Cookie.SameSite = http.SameSiteLaxMode
	session.Cookie.Secure = app.InProduction

	app.Session = session

	log.Print("Connecting to database...")
	connectionString := fmt.Sprintf("host=%s port=%s dbname=%s user=%s password=%s sslmode=%s", *dbHost, *dbPort, *dbName, *dbUser, *dbPass, *dbSSL)
	db, err := driver.ConnectSQL(connectionString)
	if err != nil {
		log.Fatal("Cannot connect to database! Dying...")
	}
	log.Print("Connected to database!")

	tc, err := render.CreateTemplateCache()
	if err != nil {
		log.Fatal("Can not create template cache")

		return nil, err
	}

	file := driver.OpenExcelFile("Report.xlsx")
	driver.ParseExcelFile(&file)
	pis := driver.GetAllPersonalInformation()
	sis := driver.GetAllSurgeriesInformation()
	app.TemplateCache = tc

	repo := handlers.NewRepo(&app, db)
	handlers.NewHandlers(repo)
	render.NewRenderer(&app)
	helpers.NewHelpers(&app)
	for index, pi := range pis {
		id, _ := handlers.Repo.DB.AddPersonalInformation(pi)
		handlers.Repo.DB.AddSurgeriesInformation(sis[index], id)
	}
	return db, nil
}
