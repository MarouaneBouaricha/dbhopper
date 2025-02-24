# dbhopper
Utility for migration of databases.

# Requirements
- Install MySQL client tools (e.g., mysqldump)

# Installation
Download the pre-released binaries from https://github.com/MarouaneBouaricha/dbhopper/releases

# How to

```shell
A command-line tool to dump MySQL databases with options to rename the database in the dump file and delete the database after dumping.

Usage:
  dbhopper [flags]

Flags:
  -d, --delete            Drop the database after dumping
  -f, --file string       Path to .txt file containing list of databases (one per line)
      --help              help for dbhopper
  -h, --host string       MySQL host (optional, falls back to MYSQL_HOST environment variable)
  -n, --name string       Name of a single database to dump
  -o, --output string     Path to destination folder of databases dumps
  -p, --password string   MySQL password (optional, falls back to MYSQL_PASSWORD environment variable)
  -r, --rename            Rename the database in the dump file with prefix 'R4_'
  -u, --username string   MySQL username (optional, falls back to MYSQL_USERNAME environment variable)
```

# Connect to MYSQL
## Using cli
```shell
dbhopper -h localhost -u root -p mysecretpassword
```
## Using environment variables
```shell
export MYSQL_USERNAME=root
export MYSQL_PASSWORD=secret
export MYSQL_HOST=localhost
```

# Examples
## one database
This will dump the database named "mydb"
```shell
dbhopper -h 172.17.0.1 -u root -p p_ssW0rd -n mydb
```

## multiple databases using a file
This will dump all the databases in side the file **databases.txt** and put the generated .sql files inside **data** folder
```shell
dbhopper -o data -f databases.txt
```
> [!NOTE]  
> "-o" option is used for destination folder

## rename dumped databases
```shell
dbhopper -r -o data -f list.txt
```

## delete databases after dump
```shell
dbhopper -d -o data -f list.txt
```

> [!WARNING]  
> The "-d" option will drop the databases after dump