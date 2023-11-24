package main

import (
	"fmt"
	"net/http"
	"os"
	"strconv"

	"github.com/danakin/web-dev-with-go-2-code_along/controllers"
	"github.com/danakin/web-dev-with-go-2-code_along/migrations"
	"github.com/danakin/web-dev-with-go-2-code_along/models"
	"github.com/danakin/web-dev-with-go-2-code_along/templates"
	"github.com/danakin/web-dev-with-go-2-code_along/views"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/gorilla/csrf"
	"github.com/joho/godotenv"
)

// added modd for dynamic reloading
// go install github.com/cortesi/modd/cmd/modd@latest

func main() {
	cfg, err := loadEnvConfig()
	if err != nil {
		panic(err)
	}

	// Set Up DB Connections
	db, err := models.Open(cfg.PSQL)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	err = models.MigrateFS(db, migrations.FS, ".")
	if err != nil {
		panic(err)
	}

	err = db.Ping()
	if err != nil {
		panic(err)
	}
	fmt.Println("Connected!")

	// Set Up Services
	userService := &models.UserService{
		DB: db,
	}

	sessionService := &models.SessionService{
		DB:            db,
		BytesPerToken: 32,
	}

	passwordResetService := &models.PasswordResetService{
		DB: db,
	}

	emailService := models.NewEmailService(cfg.SMTP)

	galleryService := models.GalleryService{
		DB: db,
	}

	// Set Up Middleware
	userMiddleware := controllers.UserMiddleware{
		SessionService: sessionService,
	}

	csrfKey := cfg.CSRF.Key
	csrfMiddleware := csrf.Protect(
		[]byte(csrfKey),
		csrf.Secure(cfg.CSRF.Secure),
		csrf.Path("/"),
	)

	// Set Up Controller
	userController := controllers.User{
		UserService:          userService,
		SessionService:       sessionService,
		PasswordResetService: passwordResetService,
		EmailService:         emailService,
	}
	userController.Templates.New = views.Must(views.ParseFS(templates.FS, "signup.gohtml", "tailwind.gohtml"))
	userController.Templates.Login = views.Must(views.ParseFS(templates.FS, "login.gohtml", "tailwind.gohtml"))
	userController.Templates.ForgotPassword = views.Must(views.ParseFS(templates.FS, "forgot-pw.gohtml", "tailwind.gohtml"))
	userController.Templates.CheckYourEmail = views.Must(views.ParseFS(templates.FS, "check-your-email.gohtml", "tailwind.gohtml"))
	userController.Templates.ResetPassword = views.Must(views.ParseFS(templates.FS, "reset-pw.gohtml", "tailwind.gohtml"))

	galleryController := controllers.Gallery{
		GalleryService: &galleryService,
	}
	galleryController.Templates.New = views.Must(views.ParseFS(templates.FS, "galleries/new.gohtml", "tailwind.gohtml"))
	galleryController.Templates.Edit = views.Must(views.ParseFS(templates.FS, "galleries/edit.gohtml", "tailwind.gohtml"))
	galleryController.Templates.Index = views.Must(views.ParseFS(templates.FS, "galleries/index.gohtml", "tailwind.gohtml"))
	galleryController.Templates.Show = views.Must(views.ParseFS(templates.FS, "galleries/show.gohtml", "tailwind.gohtml"))

	// Set Up Routing
	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(csrfMiddleware)
	r.Use(userMiddleware.SetUser)

	r.Get("/", controllers.StaticHandler(views.Must(views.ParseFS(templates.FS, "home.gohtml", "tailwind.gohtml"))))
	r.Get("/contact", controllers.StaticHandler(views.Must(views.ParseFS(templates.FS, "contact.gohtml", "tailwind.gohtml"))))
	r.Get("/faq", controllers.FAQ(views.Must(views.ParseFS(templates.FS, "faq.gohtml", "tailwind.gohtml"))))

	r.Get("/signup", userController.New)
	r.Post("/signup", userController.Store)

	r.Get("/login", userController.Login)
	r.Post("/login", userController.SignIn)

	r.Post("/logout", userController.SignOut)

	r.Get("/forgot-pw", userController.ForgotPassword)
	r.Post("/forgot-pw", userController.ProcessForgotPassword)

	r.Get("/reset-pw", userController.ResetPassword)
	r.Post("/reset-pw", userController.ProcessResetPassword)

	r.Route("/users/me", func(r chi.Router) {
		r.Use(userMiddleware.RequireUser)
		r.Get("/", userController.CurrentUser)
		r.Get("/test", func(w http.ResponseWriter, r *http.Request) {
			fmt.Fprint(w, "Success")
		})
	})

	r.Route("/galleries", func(r chi.Router) {
		r.Get("/{id}", galleryController.Show)

		r.Group(func(r chi.Router) {
			r.Use(userMiddleware.RequireUser)
			r.Get("/", galleryController.Index)
			r.Get("/new", galleryController.Create)
			r.Post("/", galleryController.Store)
			r.Get("/{id}/edit", galleryController.Edit)
			r.Post("/{id}", galleryController.Update) // todo: spoof put method?
		})
	})

	r.NotFound(func(w http.ResponseWriter, r *http.Request) { // not needed but nice to have
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
	})

	// Start the Server
	fmt.Printf("Starting the server on %s ..\n", cfg.Server.Address)
	err = http.ListenAndServe(cfg.Server.Address, r)
	if err != nil {
		panic(err)
	}
}

type config struct {
	PSQL models.PostgresConfig
	SMTP models.SMTPConfig
	CSRF struct {
		Key    string
		Secure bool
	}
	Server struct {
		Address string
	}
}

func loadEnvConfig() (config, error) {
	var cfg config

	err := godotenv.Load()
	if err != nil {
		return cfg, err
	}

	// PSQL Setup // TODO: get from .env
	cfg.PSQL = models.DefaultPostgresConfig()

	// SMTP Setup
	cfg.SMTP.Host = os.Getenv("MAIL_HOST")
	cfg.SMTP.Port, err = strconv.Atoi(os.Getenv("MAIL_PORT"))
	if err != nil {
		panic(err)
	}
	cfg.SMTP.Username = os.Getenv("MAIL_USERNAME")
	cfg.SMTP.Password = os.Getenv("MAIL_PASSWORD")

	// CSRF Setup // TODO: get from .env
	cfg.CSRF.Key = "gFvi45R4fy5xNBlnEeZtQbfAVCYEIAUX"
	cfg.CSRF.Secure = false

	// Server Setup // TODO: get from .env
	cfg.Server.Address = ":3000"

	return cfg, nil
}
