package main

import (
	"log"

	"github.com/rahadianir/dealls/internal/app"
	"github.com/rahadianir/dealls/migrations"
	"github.com/spf13/cobra"
)

func main() {
	rootCmd := &cobra.Command{Use: "hrapp"}

	migrateCmd := &cobra.Command{
		Use:   "migrate",
		Short: "Run DB migrations",
		Run: func(cmd *cobra.Command, args []string) {
			migrations.SetupData()
		},
	}

	serveHTTPCmd := &cobra.Command{
		Use:   "http",
		Short: "Start the HTTP server",
		Run: func(cmd *cobra.Command, args []string) {
			app.StartServer()
		},
	}

	rootCmd.AddCommand(migrateCmd)
	rootCmd.AddCommand(serveHTTPCmd)

	if err := rootCmd.Execute(); err != nil {
		log.Fatal(err)
	}
}
