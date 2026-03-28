package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/gofiber/swagger"

	_ "github.com/heracle/pt.heracle.fit.go/docs"
	"github.com/heracle/pt.heracle.fit.go/internal/ai"
	"github.com/heracle/pt.heracle.fit.go/internal/config"
	"github.com/heracle/pt.heracle.fit.go/internal/database"
	"github.com/heracle/pt.heracle.fit.go/internal/handler"
	mw "github.com/heracle/pt.heracle.fit.go/internal/middleware"
	"github.com/heracle/pt.heracle.fit.go/internal/repository"
	"github.com/heracle/pt.heracle.fit.go/internal/service"
)

// @title           Heracle Fitness API
// @version         1.0
// @description     Backend API for the Heracle personal training & fitness platform.
// @host            localhost:8080
// @BasePath        /
// @securityDefinitions.apikey BearerAuth
// @in              header
// @name            Authorization
// @description     JWT Bearer token. Format: "Bearer {token}"

func main() {
	cfg := config.Load()

	if cfg.DatabaseURL == "" {
		log.Fatal("DATABASE_URL is required")
	}

	// Database
	pool := database.NewPool(cfg.DatabaseURL)
	defer pool.Close()

	// Repositories
	userRepo := repository.NewUserRepo(pool)
	profileRepo := repository.NewUserProfileRepo(pool)
	mealRepo := repository.NewMealRepo(pool)
	exerciseRepo := repository.NewExerciseRepo(pool)
	sessionRepo := repository.NewSessionRepo(pool)
	workoutLogRepo := repository.NewWorkoutLogRepo(pool)
	sleepRepo := repository.NewSleepRepo(pool)
	trainerRepo := repository.NewTrainerRepo(pool)
	suggestionRepo := repository.NewDietSuggestionRepo(pool)
	foodItemRepo := repository.NewFoodItemRepo(pool)
	splitRepo := repository.NewSplitRepo(pool)

	// AI
	aiRouter := ai.NewAIRouter(cfg)

	// Services
	authSvc := service.NewAuthService(cfg, userRepo, trainerRepo)
	userSvc := service.NewUserService(userRepo, profileRepo, trainerRepo)
	dietSvc := service.NewDietService(mealRepo, profileRepo, suggestionRepo, foodItemRepo, trainerRepo, aiRouter)
	workoutSvc := service.NewWorkoutService(profileRepo, exerciseRepo, sessionRepo, workoutLogRepo, trainerRepo)
	sleepSvc := service.NewSleepService(sleepRepo, aiRouter)
	trainerSvc := service.NewTrainerService(trainerRepo, userRepo, profileRepo, mealRepo)
	splitSvc := service.NewSplitService(splitRepo, trainerRepo)

	// Handlers
	authH := handler.NewAuthHandler(authSvc)
	userH := handler.NewUserHandler(userSvc)
	dietH := handler.NewDietHandler(dietSvc)
	workoutH := handler.NewWorkoutHandler(workoutSvc)
	sleepH := handler.NewSleepHandler(sleepSvc)
	trainerH := handler.NewTrainerHandler(trainerSvc)
	splitH := handler.NewSplitHandler(splitSvc)

	// Fiber
	app := fiber.New(fiber.Config{
		BodyLimit: 50 * 1024 * 1024, // 50MB for image uploads
		ErrorHandler: func(c *fiber.Ctx, err error) error {
			code := fiber.StatusInternalServerError
			if e, ok := err.(*fiber.Error); ok {
				code = e.Code
			}
			return c.Status(code).JSON(fiber.Map{
				"statusCode": code,
				"message":    err.Error(),
			})
		},
	})

	// Middleware
	app.Use(recover.New())
	app.Use(logger.New(logger.Config{
		Format: "${time} ${status} ${method} ${path} ${latency}\n",
	}))
	app.Use(cors.New(cors.Config{
		AllowOrigins: "*",
		AllowMethods: strings.Join([]string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"}, ","),
		AllowHeaders: "Origin, Content-Type, Accept, Authorization",
	}))

	// ── Swagger ─────────────────────────────────────────────────────────────────
	app.Get("/swagger/*", swagger.HandlerDefault)

	// ── Public routes ───────────────────────────────────────────────────────────
	auth := app.Group("/auth")
	auth.Post("/google/token", authH.AuthenticateGoogleToken)
	auth.Post("/admin/login", authH.AdminLogin)
	auth.Get("/dev/token", authH.GetDevToken)

	// Health check
	app.Get("/health", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{"status": "ok"})
	})

	// ── Protected routes (JWT required) ─────────────────────────────────────────
	api := app.Group("", mw.AuthMiddleware(cfg.JWTSecret))

	// User
	user := api.Group("/user")
	user.Get("/profile", userH.GetProfile)
	user.Get("/body-metrics", userH.GetBodyMetrics)
	user.Get("/onboarding-status", userH.GetOnboardingStatus)
	user.Get("/caltest", userH.TestCalendar)
	user.Post("/body-metrics", userH.SaveBodyMetrics)
	user.Post("/targets", userH.SaveTargets)
	user.Patch("/trainer/body-metrics/:clientId", mw.TrainerGuard(), userH.TrainerSaveBodyMetrics)

	// Diet
	diet := api.Group("/diet")
	diet.Get("/status", dietH.GetStatus)
	diet.Get("/today", dietH.GetTodayDiet)
	diet.Get("/food/search", dietH.SearchFood)
	diet.Get("/meals", dietH.GetMealsByDate)
	diet.Get("/preferences", dietH.GetDietPreferences)
	diet.Post("/preferences", dietH.SaveDietPreferences)
	diet.Post("/meal", dietH.LogMeal)
	diet.Post("/ai/food", dietH.AnalyseFood)
	diet.Patch("/trainer/targets/:clientId", mw.TrainerGuard(), dietH.TrainerUpdateTargets)
	diet.Get("/trainer/meals/:clientId", mw.TrainerGuard(), dietH.TrainerGetMealsByDate)
	diet.Get("/trainer/status/:clientId", mw.TrainerGuard(), dietH.TrainerGetStatus)

	// Workout
	workout := api.Group("/workout")
	workout.Get("/today", workoutH.GetTodayWorkout)
	workout.Get("/sessions", workoutH.GetStaticSessions)
	workout.Get("/exercises", workoutH.GetExercises)
	workout.Get("/preferences", workoutH.GetWorkoutPreferences)
	workout.Post("/preferences", workoutH.SaveWorkoutPreferences)
	workout.Post("/session", workoutH.CreateSession)
	workout.Get("/session/:id", workoutH.GetSession)
	workout.Get("/my-sessions", workoutH.GetUserSessions)
	workout.Patch("/session/:id", workoutH.UpdateSession)
	workout.Delete("/session/:id", workoutH.DeleteSession)
	workout.Post("/log", workoutH.CreateWorkoutLog)
	workout.Get("/log/:id", workoutH.GetWorkoutLog)
	workout.Get("/logs", workoutH.GetWorkoutLogs)
	workout.Patch("/log/:id", workoutH.UpdateWorkoutLog)
	workout.Delete("/log/:id", workoutH.DeleteWorkoutLog)
	workout.Patch("/trainer/session/:clientId/:sessionId", mw.TrainerGuard(), workoutH.TrainerUpdateSession)
	workout.Post("/trainer/session/:clientId", mw.TrainerGuard(), workoutH.TrainerCreateSession)
	workout.Get("/trainer/sessions/:clientId", mw.TrainerGuard(), workoutH.TrainerGetSessions)
	workout.Delete("/trainer/session/:clientId/:sessionId", mw.TrainerGuard(), workoutH.TrainerDeleteSession)
	workout.Patch("/trainer/log-review/:logId", mw.TrainerGuard(), workoutH.TrainerAddLogReview)

	// Sleep
	sleep := api.Group("/sleep")
	sleep.Post("/", sleepH.AddSleepData)
	sleep.Get("/", sleepH.GetSleepData)
	sleep.Get("/insight", sleepH.GetAIInsight)

	// Trainer
	trainer := api.Group("/trainer", mw.TrainerGuard())
	trainer.Get("/clients", trainerH.GetClients)
	trainer.Get("/client/:clientId", trainerH.GetClientDetails)
	trainer.Post("/clients/add", trainerH.AddClient)
	trainer.Delete("/clients/remove/:clientId", trainerH.RemoveClient)

	// Admin-only route for adding trainers
	api.Post("/trainer/admin/add", mw.AdminGuard(), trainerH.AdminAddTrainer)

	// Split
	split := api.Group("/split")
	split.Get("/", splitH.GetMySplit)
	split.Get("/trainer/:clientId", mw.TrainerGuard(), splitH.TrainerGetClientSplit)
	split.Put("/trainer/:clientId", mw.TrainerGuard(), splitH.UpsertClientSplit)

	// ── Start ───────────────────────────────────────────────────────────────────
	port := cfg.Port
	if !strings.HasPrefix(port, ":") {
		port = ":" + port
	}

	// Graceful shutdown
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		fmt.Println("\n🛑 Shutting down...")
		_ = app.Shutdown()
	}()

	fmt.Printf("🚀 Heracle Go server starting on %s\n", port)
	fmt.Printf("📖 Swagger docs available at http://localhost%s/swagger/index.html\n", port)
	if err := app.Listen(port); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}
