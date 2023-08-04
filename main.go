package main

import (
	"fmt"
	"net/http"

	"github.com/danakin/web-dev-with-go-2-code_along/controllers"
	"github.com/danakin/web-dev-with-go-2-code_along/migrations"
	"github.com/danakin/web-dev-with-go-2-code_along/models"
	"github.com/danakin/web-dev-with-go-2-code_along/templates"
	"github.com/danakin/web-dev-with-go-2-code_along/views"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/gorilla/csrf"
)

// added modd for dynamic reloading
// go install github.com/cortesi/modd/cmd/modd@latest

func main() {
	// Set Up DB Connections
	cfg := models.DefaultPostgresConfig()
	db, err := models.Open(cfg)
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
	userService := models.UserService{
		DB: db,
	}
	sessionService := models.SessionService{
		DB:            db,
		BytesPerToken: 32,
	}

	// Set Up Middleware
	userMiddleware := controllers.UserMiddleware{
		SessionService: &sessionService,
	}

	csrfKey := "gFvi45R4fy5xNBlnEeZtQbfAVCYEIAUX"
	csrfMiddleware := csrf.Protect(
		[]byte(csrfKey),
		// TODO: set to true for HTTPS production ready code
		csrf.Secure(false),
	)

	// Set Up Controller
	userController := controllers.User{
		UserService:    &userService,
		SessionService: &sessionService,
	}
	userController.Templates.New = views.Must(views.ParseFS(templates.FS, "signup.gohtml", "tailwind.gohtml"))
	userController.Templates.Login = views.Must(views.ParseFS(templates.FS, "login.gohtml", "tailwind.gohtml"))

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
	r.Route("/users/me", func(r chi.Router) {
		r.Use(userMiddleware.RequireUser)
		r.Get("/", userController.CurrentUser)
		r.Get("/test", func(w http.ResponseWriter, r *http.Request) {
			fmt.Fprint(w, "Success")
		})
	})

	r.NotFound(func(w http.ResponseWriter, r *http.Request) { // not needed but nice to have
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
	})

	// Start the Server
	fmt.Println("Starting the server on :3000 ...")
	http.ListenAndServe(":3000", r)
}
