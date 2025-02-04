package main

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"strings"
)

// DBConfig almacena la configuración de la base de datos
type DBConfig struct {
	Type   string
	Host   string
	Name   string
	User   string
	Pass   string
	Prefix string
}

// Config almacena toda la configuración de Moodle
type Config struct {
	DB        DBConfig
	Variables map[string]string // Guarda TODAS las variables de config.php
}

// LoadConfig lee config.php y extrae la información
func LoadConfig(moodlePath string) (*Config, error) {
	configPath := moodlePath + "/config.php"
	file, err := os.Open(configPath)
	if err != nil {
		return nil, fmt.Errorf("no se pudo abrir el archivo config.php: %w", err)
	}
	defer file.Close()

	// Expresiones regulares para extraer información
	reVar := regexp.MustCompile(`\$CFG->(\w+)\s*=\s*['"](.+?)['"];`)

	config := &Config{
		Variables: make(map[string]string),
	}

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())

		if matches := reVar.FindStringSubmatch(line); matches != nil {
			key := matches[1]
			value := matches[2]
			config.Variables[key] = value

			// Guardamos valores importantes en la estructura DBConfig
			switch key {
			case "dbtype":
				config.DB.Type = value
			case "dbhost":
				config.DB.Host = value
			case "dbname":
				config.DB.Name = value
			case "dbuser":
				config.DB.User = value
			case "dbpass":
				config.DB.Pass = value
			case "prefix":
				config.DB.Prefix = value
			}
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("error al leer config.php: %w", err)
	}

	return config, nil
}
