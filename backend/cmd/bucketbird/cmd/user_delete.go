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

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/spf13/cobra"
)

var (
	deleteUserEmail string
	deleteUserID    string
)

var userDeleteCmd = &cobra.Command{
	Use:   "delete",
	Short: "Delete a user account",
	Long:  `Delete a user account by email or ID. This will cascade delete all user data.`,
	Run:   runUserDelete,
}

func init() {
	userCmd.AddCommand(userDeleteCmd)

	userDeleteCmd.Flags().StringVarP(&deleteUserEmail, "email", "e", "", "Email address of user to delete")
	userDeleteCmd.Flags().StringVar(&deleteUserID, "id", "", "User ID (UUID) to delete")
}

func runUserDelete(cmd *cobra.Command, args []string) {
	// Must provide either email or ID
	if strings.TrimSpace(deleteUserEmail) == "" && strings.TrimSpace(deleteUserID) == "" {
		fmt.Fprintln(os.Stderr, "error: either --email or --id is required")
		os.Exit(2)
	}

	// Cannot provide both
	if strings.TrimSpace(deleteUserEmail) != "" && strings.TrimSpace(deleteUserID) != "" {
		fmt.Fprintln(os.Stderr, "error: provide either --email or --id, not both")
		os.Exit(2)
	}

	cfg := config.Load()
	logger := logging.NewLogger(cfg.AppName, cfg.Env)

	ctx := context.Background()

	// Connect directly to database
	pool, err := pgxpool.New(ctx, cfg.DBDSN)
	if err != nil {
		logger.Error("failed to connect to database", slog.Any("error", err))
		os.Exit(1)
	}
	defer pool.Close()

	repos := repository.NewRepositories(pool)

	var user *repository.User

	// Look up user by email or ID
	if strings.TrimSpace(deleteUserEmail) != "" {
		user, err = repos.Users.GetByEmail(ctx, strings.TrimSpace(deleteUserEmail))
		if err != nil {
			if err == repository.ErrNotFound {
				fmt.Fprintf(os.Stderr, "error: user with email %s not found\n", deleteUserEmail)
				os.Exit(1)
			}
			logger.Error("failed to get user by email", slog.Any("error", err))
			os.Exit(1)
		}
	} else {
		id, err := uuid.Parse(strings.TrimSpace(deleteUserID))
		if err != nil {
			fmt.Fprintf(os.Stderr, "error: invalid user ID format: %v\n", err)
			os.Exit(2)
		}
		user, err = repos.Users.GetByID(ctx, id)
		if err != nil {
			if err == repository.ErrNotFound {
				fmt.Fprintf(os.Stderr, "error: user with ID %s not found\n", id)
				os.Exit(1)
			}
			logger.Error("failed to get user by ID", slog.Any("error", err))
			os.Exit(1)
		}
	}

	// Delete the user
	err = repos.Users.Delete(ctx, user.ID)
	if err != nil {
		logger.Error("failed to delete user", slog.Any("error", err), slog.String("user_id", user.ID.String()))
		os.Exit(1)
	}

	fmt.Printf("Successfully deleted user %s (ID: %s)\n", user.Email, user.ID)
	fmt.Println("Note: S3 files were not deleted and may need manual cleanup")
}
