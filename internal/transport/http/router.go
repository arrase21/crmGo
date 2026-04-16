package http

import (
	"github.com/arrase21/crm-users/internal/service"
	"github.com/arrase21/crm-users/internal/transport/middleware"
	"github.com/gin-gonic/gin"
)

func NewRouter(
	userSvc *service.UserService,
	roleSvc *service.RoleService,
	permissionSvc *service.PermissionService,
	employeeSvc *service.EmployeeService,
	payrollConceptSvc *service.PayrollConceptService,
	payrollCalculatorSvc *service.PayrollCalculatorService,
	payrollSvc *service.PayrollService,
	payrollStateSvc *service.PayrollStateService,
	batchPayrollSvc *service.PayrollBatchService,
) *gin.Engine {
	r := gin.Default()

	//CORS
	r.Use(func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Content-Type, Authorization, X-Tenant-ID")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		c.Next()
	})
	// Health Check
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status":  "ok",
			"service": "crm-api",
		})
	})

	// API v1
	v1 := r.Group("/api/v1")
	v1.Use(middleware.TenantMiddleware())
	// routes
	users := v1.Group("/users")
	{
		userHandler := NewUserHandler(userSvc)
		users.POST("", userHandler.Create)
		users.GET("", userHandler.List)
		users.GET("/search", userHandler.GetByDni)
		users.GET("/:id", userHandler.GetByID)
		users.PUT("/:id", userHandler.Update)
		users.DELETE("/:id", userHandler.Delete)
	}
	roles := v1.Group("/roles")
	{
		roleHandler := NewRoleHandler(roleSvc, permissionSvc)

		// CRUD roles
		roles.POST("", roleHandler.Create)
		roles.GET("", roleHandler.List)
		roles.GET("/:id", roleHandler.GetByID)
		roles.PUT("/:id", roleHandler.Update)
		roles.DELETE("/:id", roleHandler.Delete)

		// Permissions
		roles.POST("/:id/permissions", roleHandler.AssignPermission)
		roles.DELETE("/:id/permissions/:actionId", roleHandler.RevokePermission)

		// User ↔ Role
		roles.POST("/assign", roleHandler.AssignRoleToUser)
		roles.POST("/revoke", roleHandler.RevokeRoleFromUser)
	}

	// Employees
	employees := v1.Group("/employees")
	{
		employeeHandler := NewEmployeeHandler(employeeSvc)
		employees.POST("", employeeHandler.Create)
		employees.GET("", employeeHandler.List)
		employees.GET("/search", employeeHandler.GetByUserID)
		employees.GET("/:id", employeeHandler.GetByID)
		employees.PUT("/:id", employeeHandler.Update)
		employees.DELETE("/:id", employeeHandler.Delete)
	}

	// Payroll Concepts
	payrollConcepts := v1.Group("/payroll-concepts")
	{
		conceptHandler := NewPayrollConceptHandler(payrollConceptSvc)
		payrollConcepts.POST("", conceptHandler.Create)
		payrollConcepts.GET("", conceptHandler.List)
		payrollConcepts.GET("/search", conceptHandler.GetByCode)
		payrollConcepts.GET("/active", conceptHandler.GetActiveConcepts)
		payrollConcepts.POST("/seed", conceptHandler.SeedDefaultConcepts)
		payrollConcepts.GET("/:id", conceptHandler.GetByID)
		payrollConcepts.PUT("/:id", conceptHandler.Update)
		payrollConcepts.DELETE("/:id", conceptHandler.Delete)
	}

	// Payroll (Nómina)
	payroll := v1.Group("/payroll")
	{
		payrollHandler := NewPayrollHandler(payrollCalculatorSvc, payrollSvc)
		payroll.POST("/calculate", payrollHandler.Calculate)
		payroll.POST("/calculate-and-save", payrollHandler.CalculateAndSave)
		payroll.GET("/employee/:employeeId", payrollHandler.ListByEmployee)
		payroll.GET("/:id", payrollHandler.GetByID)
		payroll.DELETE("/:id", payrollHandler.Delete)

		// Nuevos endpoints de estado y batch
		stateHandler := NewPayrollStateHandler(payrollStateSvc, batchPayrollSvc, payrollSvc)
		payroll.POST("/:id/mark-paid", stateHandler.MarkAsPaid)
		payroll.POST("/:id/revert-to-draft", stateHandler.RevertToDraft)
		payroll.GET("/:id/payment", stateHandler.GetPaymentInfo)
		payroll.POST("/batch", stateHandler.ProcessBatch)
		payroll.GET("/summary", stateHandler.GetPayrollSummary)
	}

	return r
}
