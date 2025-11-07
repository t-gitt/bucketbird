package cmd

import (
	"context"
	"fmt"
	"log/slog"
	"os"

	"bucketbird/backend/internal/config"
	"bucketbird/backend/internal/logging"
	"bucketbird/backend/internal/repository"
	"bucketbird/backend/pkg/crypto"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/spf13/cobra"
)

var (
	resetPasswordEmail    string
	resetPasswordPassword string
)

var userResetPasswordCmd = &cobra.Command{
	Use:   "reset-password",
	Short: "Reset a user's password",
	Long:  `Reset a user's password by email address.`,
	Run:   runUserResetPassword,
}

func init() {
	userCmd.AddCommand(userResetPasswordCmd)

	userResetPasswordCmd.Flags().StringVarP(&resetPasswordEmail, "email", "e", "", "User email address (required)")
	userResetPasswordCmd.Flags().StringVarP(&resetPasswordPassword, "password", "p", "", "New password (required, min 8 characters)")

	userResetPasswordCmd.MarkFlagRequired("email")
	userResetPasswordCmd.MarkFlagRequired("password")
}

func runUserResetPassword(cmd *cobra.Command, args []string) {
	if resetPasswordEmail == "" || resetPasswordPassword == "" {
		fmt.Fprintln(os.Stderr, "Error: Both --email and --password flags are required")
		os.Exit(1)
	}

	if len(resetPasswordPassword) < 8 {
		fmt.Fprintln(os.Stderr, "Error: Password must be at least 8 characters")
		os.Exit(1)
	}

	// Load configuration
	cfg := config.Load()
	logger := logging.NewLogger(cfg.AppName, cfg.Env)

	// Connect to database
	ctx := context.Background()

	pool, err := pgxpool.New(ctx, cfg.DBDSN)
	if err != nil {
		logger.Error("failed to connect to database", slog.Any("error", err))
		os.Exit(1)
	}
	defer pool.Close()

	// Initialize repositories
	repos := repository.NewRepositories(pool)

	// Look up user by email
	user, err := repos.Users.GetByEmail(ctx, resetPasswordEmail)
	if err != nil {
		if err == repository.ErrNotFound {
			fmt.Fprintf(os.Stderr, "User not found: %s\n", resetPasswordEmail)
			os.Exit(1)
		}
		fmt.Fprintf(os.Stderr, "Failed to find user: %v\n", err)
		os.Exit(1)
	}

	// Hash the new password
	passwordHash, err := crypto.HashPassword(resetPasswordPassword)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to hash password: %v\n", err)
		os.Exit(1)
	}

	// Update the password in database
	if err := repos.Users.UpdatePassword(ctx, user.ID, passwordHash); err != nil {
		fmt.Fprintf(os.Stderr, "Failed to update password: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("✓ Password successfully reset for user: %s (%s %s)\n", user.Email, user.FirstName, user.LastName)
}
