package cmd

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"strings"

	"bucketbird/backend/internal/config"
	"bucketbird/backend/internal/logging"
	"bucketbird/backend/internal/repository"
	"bucketbird/backend/internal/service"
	"bucketbird/backend/pkg/jwt"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/spf13/cobra"
)

var (
	createUserEmail     string
	createUserPassword  string
	createUserFirstName string
	createUserLastName  string
)

var userCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a new user account",
	Long:  `Create a new user account with email and password.`,
	Run:   runUserCreate,
}

func init() {
	userCmd.AddCommand(userCreateCmd)

	userCreateCmd.Flags().StringVarP(&createUserEmail, "email", "e", "", "Email address (required)")
	userCreateCmd.Flags().StringVarP(&createUserPassword, "password", "p", "", "Password (required)")
	userCreateCmd.Flags().StringVar(&createUserFirstName, "first-name", "", "First name")
	userCreateCmd.Flags().StringVar(&createUserLastName, "last-name", "", "Last name")

	userCreateCmd.MarkFlagRequired("email")
	userCreateCmd.MarkFlagRequired("password")
}

func runUserCreate(cmd *cobra.Command, args []string) {
	if strings.TrimSpace(createUserEmail) == "" || strings.TrimSpace(createUserPassword) == "" {
		fmt.Fprintln(os.Stderr, "error: --email and --password are required")
		os.Exit(2)
	}

	cfg := config.Load()
	logger := logging.NewLogger(cfg.AppName, cfg.Env)

	ctx := context.Background()

	// Connect to database
	pool, err := pgxpool.New(ctx, cfg.DBDSN)
	if err != nil {
		logger.Error("failed to connect to database", slog.Any("error", err))
		os.Exit(1)
	}
	defer pool.Close()

	// Initialize repositories
	repos := repository.NewRepositories(pool)

	// Initialize services
	tokenManager := jwt.NewTokenManager(cfg.JWTSecret, cfg.AccessTokenTTL)
	authService := service.NewAuthService(
		repos.Users,
		repos.Sessions,
		tokenManager,
		cfg.RefreshTokenTTL,
		logger,
	)

	result, err := authService.Register(ctx, service.RegisterInput{
		Email:     strings.TrimSpace(createUserEmail),
		Password:  strings.TrimSpace(createUserPassword),
		FirstName: strings.TrimSpace(createUserFirstName),
		LastName:  strings.TrimSpace(createUserLastName),
	})
	if err != nil {
		logger.Error("failed to create user", slog.Any("error", err))
		os.Exit(1)
	}

	fmt.Printf("Created user %s (ID: %s)\n", result.User.Email, result.User.ID)
}
