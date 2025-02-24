package cmd

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

func readDatabases(filePath string) ([]string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var databases []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		dbName := strings.TrimSpace(scanner.Text())
		if dbName != "" {
			databases = append(databases, dbName)
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return databases, nil
}

func processDatabase(dbName string) {
	fullPath := fmt.Sprintf("%s.sql", dbName)
	if destFolder != "" {
		if _, err := os.Stat(destFolder); os.IsNotExist(err) {
			err := os.MkdirAll(destFolder, os.ModePerm)
			if err != nil {
				fmt.Println("Error creating directory:", err)
				return
			}
			fmt.Println("Destination folder created:", destFolder)
		}
		fullPath = filepath.Join(destFolder, fmt.Sprintf("%s.sql", dbName))
	}

	err := createDump(dbName, fullPath)
	if err != nil {
		fmt.Printf("Error creating dump for database %s: %v\n", dbName, err)
		return
	}
	fmt.Printf("Dump created successfully for database %s: %s\n", dbName, fullPath)

	if rename {
		err = renameDatabaseInDump(fullPath, dbName)
		if err != nil {
			fmt.Printf("Error renaming database in dump for database %s: %v\n", dbName, err)
			return
		}
		fmt.Printf("Database renamed in dump file with prefix 'R4_' for database %s: %s\n", dbName, fullPath)
	}

	if deleteDB {
		err = dropDatabase(dbName)
		if err != nil {
			fmt.Printf("Error dropping database %s: %v\n", dbName, err)
			return
		}
		fmt.Printf("Database dropped: %s\n", dbName)
	}
}

func createDump(dbName, dumpFile string) error {
	cmd := exec.Command("mysqldump", "-BR", "-h", host, "-u", username, "-p"+password, dbName)
	output, err := os.Create(dumpFile)
	if err != nil {
		return err
	}
	defer output.Close()
	cmd.Stdout = output
	return cmd.Run()
}

func renameDatabaseInDump(dumpFile, dbName string) error {
	data, err := os.ReadFile(dumpFile)
	if err != nil {
		return err
	}

	prefixedName := "R4_" + dbName
	newData := strings.ReplaceAll(string(data), fmt.Sprintf("`%s`", dbName), fmt.Sprintf("`%s`", prefixedName))

	err = os.WriteFile(dumpFile, []byte(newData), 0644)
	if err != nil {
		return err
	}
	return nil
}

func dropDatabase(dbName string) error {
	cmd := exec.Command("mysql", "-h", host, "-u", username, "-p"+password, "-e", fmt.Sprintf("DROP DATABASE `%s`;", dbName))
	return cmd.Run()
}
