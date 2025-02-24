package cmd

import (
	"fmt"
	"os"
	"sync"

	"github.com/spf13/cobra"
)

var (
	dbFile     string
	dbName     string
	rename     bool
	deleteDB   bool
	username   string
	password   string
	host       string
	destFolder string
)

var rootCmd = &cobra.Command{
	Use:   "dbhopper",
	Short: "A tool to dump MySQL databases",
	Long:  `A command-line tool to dump MySQL databases with options to rename the database in the dump file and delete the database after dumping.`,
	Run: func(cmd *cobra.Command, args []string) {
		if username == "" {
			username = os.Getenv("MYSQL_USERNAME")
		}
		if password == "" {
			password = os.Getenv("MYSQL_PASSWORD")
		}
		if host == "" {
			host = os.Getenv("MYSQL_HOST")
		}

		if username == "" || password == "" || host == "" {
			fmt.Println("Error: MySQL connection details (username, password, host) are required.")
			fmt.Println("Provide them via command-line arguments or environment variables.")
			os.Exit(1)
		}

		if dbFile == "" && dbName == "" {
			fmt.Println("Error: Either --file or --name must be provided.")
			os.Exit(1)
		}

		if destFolder != "" {
			if _, err := os.Stat(destFolder); os.IsNotExist(err) {
				err := os.MkdirAll(destFolder, os.ModePerm)
				if err != nil {
					fmt.Println("Error creating directory:", err)
					return
				}
				fmt.Println("Destination folder created:", destFolder)
			}

		}

		var databases []string
		if dbFile != "" {
			var err error
			databases, err = readDatabases(dbFile)
			if err != nil {
				fmt.Printf("Error reading databases from file: %v\n", err)
				os.Exit(1)
			}
		} else {
			databases = []string{dbName}
		}

		var wg sync.WaitGroup
		wg.Add(len(databases))

		for _, db := range databases {
			go func(dbName string) {
				defer wg.Done()
				processDatabase(destFolder, dbName)
			}(db)
		}

		wg.Wait()
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.PersistentFlags().BoolP("help", "", false, "help for dbhopper")

	rootCmd.Flags().StringVarP(&dbFile, "file", "f", "", "Path to .txt file containing list of databases (one per line)")
	rootCmd.Flags().StringVarP(&destFolder, "output", "o", "", "Path to destination folder of databases dumps")
	rootCmd.Flags().StringVarP(&dbName, "name", "n", "", "Name of a single database to dump")
	rootCmd.Flags().BoolVarP(&rename, "rename", "r", false, "Rename the database in the dump file with prefix 'R4_'")
	rootCmd.Flags().BoolVarP(&deleteDB, "delete", "d", false, "Drop the database after dumping")

	rootCmd.Flags().StringVarP(&username, "username", "u", "", "MySQL username (optional, falls back to MYSQL_USERNAME environment variable)")
	rootCmd.Flags().StringVarP(&password, "password", "p", "", "MySQL password (optional, falls back to MYSQL_PASSWORD environment variable)")
	rootCmd.Flags().StringVarP(&host, "host", "h", "", "MySQL host (optional, falls back to MYSQL_HOST environment variable)")
}
