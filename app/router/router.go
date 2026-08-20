package router

import (
	"keuangan-api/app/handler"
	"keuangan-api/app/middleware"
	"keuangan-api/app/repository"
	"keuangan-api/app/service"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func Setup(db *gorm.DB) *gin.Engine {
	r := gin.Default()

	// --- Shared repositories ---
	userRepo   := &repository.UserRepository{DB: db}
	notifRepo  := &repository.NotificationRepository{DB: db}
	budgetRepo := &repository.BudgetRepository{DB: db}
	txRepo     := &repository.TransactionRepository{DB: db}

	// --- Auth ---
	authHandler := &handler.AuthHandler{
		Service: &service.AuthService{
			Repo: &repository.AuthRepository{DB: db},
		},
	}

	// --- Users / Contacts ---
	userHandler := &handler.UserHandler{
		Service: &service.UserService{
			Repo: userRepo,
		},
	}

	// --- Categories ---
	categoryHandler := &handler.CategoryHandler{
		Service: &service.CategoryService{
			Repo: &repository.CategoryRepository{DB: db},
		},
	}

	// --- Transactions ---
	transactionHandler := &handler.TransactionHandler{
		Service: &service.TransactionService{
			Repo:       txRepo,
			BudgetRepo: budgetRepo,
			NotifRepo:  notifRepo,
		},
	}

	// --- Budgets ---
	budgetHandler := &handler.BudgetHandler{
		Service: &service.BudgetService{
			Repo:            budgetRepo,
			TransactionRepo: txRepo,
		},
	}

	// --- Savings ---
	savingHandler := &handler.SavingHandler{
		Service: &service.SavingService{
			Repo:      &repository.SavingRepository{DB: db},
			UserRepo:  userRepo,
			NotifRepo: notifRepo,
		},
	}

	// --- Agendas ---
	agendaHandler := &handler.AgendaHandler{
		Service: &service.AgendaService{
			Repo:     &repository.AgendaRepository{DB: db},
			UserRepo: userRepo,
		},
	}

	// --- Notifications ---
	notifHandler := &handler.NotificationHandler{
		Service: &service.NotificationService{
			Repo: notifRepo,
		},
	}

	// --- Password Reset ---
	passwordResetHandler := &handler.PasswordResetHandler{
		Service: &service.PasswordResetService{
			Repo:     &repository.PasswordResetRepository{DB: db},
			AuthRepo: &repository.AuthRepository{DB: db},
		},
	}

	// --- API Docs ---
	apiDocHandler := &handler.APIDocHandler{
		Service: &service.APIDocService{
			Repo: &repository.APIDocRepository{DB: db},
		},
	}

	// Web UI Dokumentasi API
	r.GET("/docs", apiDocHandler.RenderDocsUI)

	// =========================================================
	// Public Routes
	// =========================================================
	api := r.Group("/api")
	{
		api.GET("/docs", apiDocHandler.GetAll)
		api.POST("/register", authHandler.Register)
		api.POST("/login", authHandler.Login)
		api.POST("/forgot-password", passwordResetHandler.ForgotPassword)
		api.POST("/reset-password", passwordResetHandler.ResetPassword)
	}

	// =========================================================
	// Protected Routes
	// =========================================================
	protected := r.Group("/api")
	protected.Use(middleware.AuthMiddleware())
	{
		// Account & Contacts
		protected.POST("/change-password", authHandler.ChangePassword)
		protected.GET("/contacts", userHandler.GetRecentContacts)

		// Categories
		protected.GET("/categories", categoryHandler.GetCategories)
		protected.POST("/categories", categoryHandler.CreateCategory)

		// Transactions CRUD + Reports & Export
		protected.POST("/transactions", transactionHandler.Create)
		protected.GET("/transactions", transactionHandler.GetAll)
		protected.GET("/transactions/report", transactionHandler.GetReport)
		protected.GET("/transactions/export", transactionHandler.ExportCSV)
		protected.PUT("/transactions/:id", transactionHandler.Update)
		protected.DELETE("/transactions/:id", transactionHandler.Delete)

		// Budgets CRUD
		protected.POST("/budgets", budgetHandler.Create)
		protected.GET("/budgets", budgetHandler.GetMyBudgets)
		protected.PUT("/budgets/:id", budgetHandler.Update)
		protected.DELETE("/budgets/:id", budgetHandler.Delete)

		// Savings CRUD + actions
		protected.POST("/savings", savingHandler.CreateGoal)
		protected.GET("/savings", savingHandler.GetMyGoals)
		protected.PUT("/savings/:id", savingHandler.UpdateGoal)
		protected.DELETE("/savings/:id", savingHandler.DeleteGoal)
		protected.POST("/savings/:id/members", savingHandler.AddMember)
		protected.POST("/savings/:id/contribute", savingHandler.Contribute)
		protected.POST("/savings/:id/withdraw", savingHandler.Withdraw)
		protected.GET("/savings/:id/contributions", savingHandler.GetContributionHistory)

		// Agendas CRUD + status + actions
		protected.POST("/agendas", agendaHandler.CreateAgenda)
		protected.GET("/agendas", agendaHandler.GetMyAgendas)
		protected.PUT("/agendas/:id", agendaHandler.UpdateAgenda)
		protected.PUT("/agendas/:id/status", agendaHandler.UpdateStatus)
		protected.DELETE("/agendas/:id", agendaHandler.DeleteAgenda)
		protected.POST("/agendas/:id/members", agendaHandler.AddMember)

		// Notifications
		protected.GET("/notifications", notifHandler.GetAll)
		protected.PUT("/notifications/:id/read", notifHandler.MarkAsRead)
	}

	return r
}
