package main

import (
	"database/sql"
	"fmt"

	_ "github.com/go-sql-driver/mysql" // Driver para MySQL
	_ "github.com/lib/pq"              // Driver para PostgreSQL
)

// GetMoodleVersion obtiene la versión de Moodle desde la base de datos
func GetMoodleVersion(dbConfig DBConfig) (string, error) {
	// Seleccionamos el driver y construimos el DSN (Data Source Name)
	var dsn string
	var driver string

	switch dbConfig.Type {
	case "mysqli", "mariadb":
		driver = "mysql"
		dsn = fmt.Sprintf("%s:%s@tcp(%s)/%s?parseTime=true",
			dbConfig.User, dbConfig.Pass, dbConfig.Host, dbConfig.Name)
	case "pgsql":
		driver = "postgres"
		dsn = fmt.Sprintf("postgres://%s:%s@%s/%s?sslmode=disable",
			dbConfig.User, dbConfig.Pass, dbConfig.Host, dbConfig.Name)
	default:
		return "", fmt.Errorf("tipo de base de datos no soportado: %s", dbConfig.Type)
	}

	// Conectar a la base de datos
	db, err := sql.Open(driver, dsn)
	if err != nil {
		return "", fmt.Errorf("error al conectar con la base de datos: %w", err)
	}
	defer db.Close()

	// Realizar la consulta para obtener la versión de Moodle
	var version string
	query := fmt.Sprintf("SELECT value FROM %sconfig WHERE name = 'release';", dbConfig.Prefix)
	err = db.QueryRow(query).Scan(&version)
	if err != nil {
		return "", fmt.Errorf("error al obtener la versión de Moodle: %w", err)
	}

	return version, nil
}
