package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/arrase21/crm-users/internal/config"
	"github.com/arrase21/crm-users/internal/database"
	"github.com/arrase21/crm-users/internal/repository"
	"github.com/arrase21/crm-users/internal/service"
	transportHttp "github.com/arrase21/crm-users/internal/transport/http"
	"github.com/joho/godotenv"
)

func main() {

	// Cargar .env solo en local
	if err := godotenv.Load(); err != nil {
		log.Println("ℹ️ .env no encontrado, usando variables del sistema")
	}

	log.Println("1️⃣ cargando configuración")
	pgCfg := config.LoadPostgres()

	log.Println("2️⃣ conectando a la base de datos")

	db, err := database.Connect(pgCfg)
	if err != nil {
		log.Fatalf("❌ Failed to connect to database: %v", err)
	}
	defer database.Close(db)

	if err := database.Automigrate(db); err != nil {
		log.Fatalf("❌ Failed to migrate database: %v", err)
	}

	userRepo := repository.NewGormUserRepository(db)
	userService := service.NewUserService(userRepo)
	roleRepo := repository.NewGormRoleRepository(db)
	roleService := service.NewRoleService(roleRepo)
	userRoleRepo := repository.NewGormUserRoleRepository(db)
	permissionService := service.NewPermissionService(userRoleRepo, roleRepo)

	// Employee
	employeeRepo := repository.NewGormEmployeeRepository(db)
	employeeService := service.NewEmployeeService(employeeRepo)

	// Payroll Concept
	payrollConceptRepo := repository.NewGormPayrollConceptRepository(db)
	payrollConceptService := service.NewPayrollConceptService(payrollConceptRepo)

	// Employee Contract
	contractRepo := repository.NewGormEmployeeContractRepository(db)

	// Payroll
	payrollRepo := repository.NewGormPayrollRepository(db)
	payrollService := service.NewPayrollService(payrollRepo)
	payrollItemRepo := repository.NewGormPayrollItemRepository(db)

	// Payment
	paymentRepo := repository.NewGormPaymentRepository(db)

	// Payroll Calculator
	payrollCalculatorService := service.NewPayrollCalculatorService(
		payrollRepo,
		payrollItemRepo,
		employeeRepo,
		contractRepo,
		payrollConceptRepo,
	)

	// Payroll State Service (transiciones de estado)
	payrollStateService := service.NewPayrollStateService(
		payrollRepo,
		paymentRepo,
		employeeRepo,
	)

	// Batch Payroll Service
	batchPayrollService := service.NewPayrollBatchService(
		payrollRepo,
		payrollItemRepo,
		employeeRepo,
		contractRepo,
		payrollConceptRepo,
		payrollStateService,
	)

	router := transportHttp.NewRouter(
		userService,
		roleService,
		permissionService,
		employeeService,
		payrollConceptService,
		payrollCalculatorService,
		payrollService,
		payrollStateService,
		batchPayrollService,
	)
	port := getEnv("PORT", "8080")
	srv := &http.Server{
		Addr:         ":" + port,
		Handler:      router,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}
	go func() {
		log.Printf("🚀 Server starting on port %s", port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("❌ Failed to start server: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("🛑 Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("❌ Server forced to shutdown: %v", err)
	}

	log.Println("✅ Server exited gracefully")

}
func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}
