package psqlutils

import "fmt"

func BuildPostgresURL(
	DbHost string,
	DbPort int,
	DbUsername string,
	DbPassword string,
	DbName string,
	DbSslmode bool,
) string {
	sslmode := "disable"
	if DbSslmode {
		sslmode = "require"
	}

	// Building the URL
	return fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=%s",
		DbUsername, DbPassword, DbHost, DbPort, DbName, sslmode)
}

// BuildDSN builds the data source name (DSN) string based on the config struct
func BuildDSN(
	DbHost string,
	DbPort int,
	DbUsername string,
	DbPassword string,
	DbName string,
	DbSslmode bool,
) string {
	dsn := ""

	// Only include non-empty fields in the DSN
	if DbUsername != "" {
		dsn += fmt.Sprintf("user=%s ", DbUsername)
	}
	if DbPassword != "" {
		dsn += fmt.Sprintf("password=%s ", DbPassword)
	}
	if DbHost != "" {
		dsn += fmt.Sprintf("host=%s ", DbHost)
	}
	if DbPort != 0 {
		dsn += fmt.Sprintf("port=%d ", DbPort)
	}
	if DbName != "" {
		dsn += fmt.Sprintf("dbname=%s ", DbName)
	}

	// Append sslmode based on the boolean flag
	if DbSslmode {
		dsn += "sslmode=require "
	} else {
		dsn += "sslmode=disable "
	}

	return dsn
}
