package cmd

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"time"

	"bucketbird/backend/internal/config"
	"bucketbird/backend/internal/logging"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/spf13/cobra"
)

var userListCmd = &cobra.Command{
	Use:   "list",
	Short: "List all user accounts",
	Long:  `List all user accounts in the system with their details.`,
	Run:   runUserList,
}

func init() {
	userCmd.AddCommand(userListCmd)
}

func runUserList(cmd *cobra.Command, args []string) {
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

	// Query all users
	rows, err := pool.Query(ctx, `
		SELECT id, email, first_name, last_name, created_at
		FROM users
		ORDER BY created_at DESC
	`)
	if err != nil {
		logger.Error("failed to query users", slog.Any("error", err))
		os.Exit(1)
	}
	defer rows.Close()

	fmt.Println("\nUsers:")
	fmt.Println("================================================================================")
	fmt.Printf("%-38s %-30s %-20s %s\n", "ID", "Email", "Name", "Created At")
	fmt.Println("--------------------------------------------------------------------------------")

	count := 0
	for rows.Next() {
		var (
			id        string
			email     string
			firstName string
			lastName  string
			createdAt time.Time
		)

		if err := rows.Scan(&id, &email, &firstName, &lastName, &createdAt); err != nil {
			logger.Error("failed to scan user row", slog.Any("error", err))
			continue
		}

		name := fmt.Sprintf("%s %s", firstName, lastName)
		fmt.Printf(
			"%-38s %-30s %-20s %s\n",
			id,
			email,
			name,
			createdAt.UTC().Format("2006-01-02 15:04:05"),
		)
		count++
	}

	if err := rows.Err(); err != nil {
		logger.Error("error iterating users", slog.Any("error", err))
		os.Exit(1)
	}

	fmt.Println("--------------------------------------------------------------------------------")
	fmt.Printf("Total: %d user(s)\n\n", count)
}
