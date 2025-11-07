package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "bucketbird",
	Short: "BucketBird - S3 bucket management platform",
	Long: `BucketBird is a platform for managing and interacting with S3-compatible storage.

It provides a unified interface for multiple S3 providers, bucket management,
and user administration.`,
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
