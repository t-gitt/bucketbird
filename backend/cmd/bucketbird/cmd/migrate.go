package cmd

import (
	"fmt"
	"os"
	"os/exec"

	"bucketbird/backend/internal/config"

	"github.com/spf13/cobra"
)

var migrateCmd = &cobra.Command{
	Use:   "migrate",
	Short: "Database migration commands",
	Long:  `Run database migrations using golang-migrate.`,
}

var migrateUpCmd = &cobra.Command{
	Use:   "up",
	Short: "Apply all pending migrations",
	Run:   runMigrateUp,
}

var migrateDownCmd = &cobra.Command{
	Use:   "down",
	Short: "Rollback the last migration",
	Run:   runMigrateDown,
}

var migrateVersionCmd = &cobra.Command{
	Use:   "version",
	Short: "Show current migration version",
	Run:   runMigrateVersion,
}

func init() {
	rootCmd.AddCommand(migrateCmd)
	migrateCmd.AddCommand(migrateUpCmd)
	migrateCmd.AddCommand(migrateDownCmd)
	migrateCmd.AddCommand(migrateVersionCmd)
}

func runMigrate(args ...string) error {
	cfg := config.Load()

	cmdArgs := []string{
		"-path", "/app/migrations",
		"-database", cfg.DBDSN,
	}
	cmdArgs = append(cmdArgs, args...)

	cmd := exec.Command("migrate", cmdArgs...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	return cmd.Run()
}

func runMigrateUp(cmd *cobra.Command, args []string) {
	if err := runMigrate("up"); err != nil {
		fmt.Fprintf(os.Stderr, "migration failed: %v\n", err)
		os.Exit(1)
	}
}

func runMigrateDown(cmd *cobra.Command, args []string) {
	if err := runMigrate("down", "1"); err != nil {
		fmt.Fprintf(os.Stderr, "migration failed: %v\n", err)
		os.Exit(1)
	}
}

func runMigrateVersion(cmd *cobra.Command, args []string) {
	if err := runMigrate("version"); err != nil {
		fmt.Fprintf(os.Stderr, "migration failed: %v\n", err)
		os.Exit(1)
	}
}
