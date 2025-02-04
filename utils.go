package main

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

// cleanTempAndReports elimina las carpetas ./temp/ y ./reports/
func cleanTempAndReports() {
	tempDir := getTempDir()
	reportsDir := getReportsDir()

	fmt.Println("🧹 Limpiando archivos temporales y reportes...")

	// Eliminar la carpeta temp
	if fileExists(tempDir) {
		os.RemoveAll(tempDir)
		fmt.Println("🗑️ Eliminada carpeta temporal:", tempDir)
	}

	// Eliminar la carpeta reports
	if fileExists(reportsDir) {
		os.RemoveAll(reportsDir)
		fmt.Println("🗑️ Eliminada carpeta de informes:", reportsDir)
	}

	fmt.Println("✅ Limpieza completada.")
}

// readIgnoreList carga patrones de archivos y carpetas a ignorar desde ignore.conf
func readIgnoreList() []*regexp.Regexp {
	var ignorePatterns []*regexp.Regexp

	file, err := os.Open("ignore.conf")
	if err != nil {
		fmt.Println("⚠️ No se encontró ignore.conf, no se ignorarán archivos.")
		return ignorePatterns
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		pattern := strings.TrimSpace(scanner.Text())
		if pattern == "" || strings.HasPrefix(pattern, "#") {
			continue
		}

		// Convertir ".git/" en "**/.git/" para ignorar todos los directorios .git en cualquier parte
		// if pattern == ".git/" {
		// 	pattern = "**/.git/"
		// }

		// Convertir "*" en ".*" para hacer la coincidencia más flexible
		pattern = strings.ReplaceAll(pattern, ".", `\.`) // Escapar "."
		pattern = strings.ReplaceAll(pattern, "*", ".*") // Convertir "*" en ".*"

		// Asegurar que los directorios se ignoran recursivamente
		if strings.HasSuffix(pattern, "/") {
			pattern += ".*"
		}

		re, err := regexp.Compile("^" + pattern + "$")
		if err != nil {
			fmt.Println("❌ Error en patrón de ignore.conf:", pattern)
			continue
		}
		ignorePatterns = append(ignorePatterns, re)
	}

	if err := scanner.Err(); err != nil {
		fmt.Println("❌ Error al leer ignore.conf:", err)
	}

	return ignorePatterns
}

// shouldIgnore verifica si un archivo/carpeta debe ser ignorado
func shouldIgnore(path string, ignorePatterns []*regexp.Regexp) bool {
	path = filepath.ToSlash(path) // Normalizar rutas

	for _, re := range ignorePatterns {
		if re.MatchString(path) {
			return true
		}

		// Si un directorio padre está en ignore.conf, también ignoramos su contenido
		for strings.HasSuffix(path, "/") {
			path = filepath.Dir(path)
			if re.MatchString(path) {
				return true
			}
		}
	}

	return false
}
