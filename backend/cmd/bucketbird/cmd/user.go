package cmd

import "github.com/spf13/cobra"

var userCmd = &cobra.Command{
	Use:   "user",
	Short: "User management commands",
	Long:  `Create, delete, and manage user accounts.`,
}

func init() {
	rootCmd.AddCommand(userCmd)
}
