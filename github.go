package main

import (
	"archive/zip"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"runtime"
	"time"
)

// cleanVersion limpia la versión eliminando "+" y "(Build: ...)"
func cleanVersion(version string) string {
	re := regexp.MustCompile(`\+.*|\(.*\)`) // Elimina "+..." y "(...)"
	return re.ReplaceAllString(version, "")
}

// getTempDir devuelve la ruta del directorio temporal
func getTempDir() string {
	tempDir := "./temp"
	if _, err := os.Stat(tempDir); os.IsNotExist(err) {
		os.Mkdir(tempDir, os.ModePerm)
	}
	return tempDir
}

// DownloadMoodle descarga Moodle o usa un ZIP local en ./temp/
func DownloadMoodle(version, localZip string) (string, error) {
	tempDir := getTempDir()
	var zipPath string

	if localZip != "" {
		fmt.Println("📂 Usando archivo ZIP local:", localZip)
		zipPath = localZip
	} else {
		cleanVer := cleanVersion(version)
		url := fmt.Sprintf("https://github.com/moodle/moodle/archive/refs/tags/v%s.zip", cleanVer)
		zipPath = filepath.Join(tempDir, fmt.Sprintf("moodle-%s.zip", cleanVer))

		// Eliminar ZIP si existe
		if fileExists(zipPath) {
			fmt.Println("🗑️ Eliminando ZIP antiguo:", zipPath)
			os.Remove(zipPath)
		}

		fmt.Println("⬇️ Descargando Moodle desde:", url)
		if err := downloadFile(url, zipPath); err != nil {
			return "", fmt.Errorf("error al descargar Moodle: %w", err)
		}
	}

	extractPath := filepath.Join(tempDir, filepath.Base(zipPath[:len(zipPath)-4]))

	// Eliminar carpeta extraída si ya existía
	if fileExists(extractPath) {
		fmt.Println("🗑️ Eliminando carpeta antigua:", extractPath)
		os.RemoveAll(extractPath)
	}

	// Extraer Moodle
	fmt.Println("📦 Extrayendo Moodle en:", extractPath)
	if err := Unzip(zipPath, extractPath); err != nil {
		return "", fmt.Errorf("error al extraer Moodle: %w", err)
	}

	// Si usamos un ZIP local, no lo eliminamos
	if localZip == "" {
		if err := os.Remove(zipPath); err == nil {
			fmt.Println("🗑️ Archivo ZIP eliminado:", zipPath)
		}
	}

	// Buscar la carpeta real dentro del ZIP extraído
	realExtractedPath, err := findMoodleRoot(extractPath)
	if err != nil {
		return "", fmt.Errorf("error al encontrar la carpeta de Moodle: %w", err)
	}

	fmt.Println("📂 Moodle listo en:", realExtractedPath)
	return realExtractedPath, nil
}

// downloadFile descarga un archivo desde una URL y lo guarda en un destino local
func downloadFile(url, dest string) error {
	resp, err := http.Get(url)
	if err != nil {
		return fmt.Errorf("error al descargar %s: %w", url, err)
	}
	defer resp.Body.Close()

	// Crear el archivo destino
	out, err := os.Create(dest)
	if err != nil {
		return fmt.Errorf("error al crear el archivo %s: %w", dest, err)
	}
	defer out.Close()

	// Copiar el contenido descargado en el archivo
	_, err = io.Copy(out, resp.Body)
	if err != nil {
		return fmt.Errorf("error al escribir en el archivo %s: %w", dest, err)
	}

	return nil
}

// findMoodleRoot encuentra la carpeta real dentro del directorio extraído
func findMoodleRoot(basePath string) (string, error) {
	files, err := os.ReadDir(basePath)
	if err != nil {
		return "", fmt.Errorf("error al leer %s: %w", basePath, err)
	}

	// Si la carpeta extraída tiene una única subcarpeta, usamos esa
	if len(files) == 1 && files[0].IsDir() {
		realPath := filepath.Join(basePath, files[0].Name())
		fmt.Println("📂 Moodle detectado en:", realPath)
		return realPath, nil
	}

	// Si hay múltiples archivos/directorios, asumimos que Moodle está en `basePath`
	fmt.Println("📂 Moodle detectado en:", basePath)
	return basePath, nil
}

// Unzip extrae un archivo ZIP a un directorio
func Unzip(src, dest string) error {
	r, err := zip.OpenReader(src)
	if err != nil {
		return err
	}
	defer r.Close() // Cerramos antes de salir de la función

	for _, f := range r.File {
		fpath := filepath.Join(dest, f.Name)

		if f.FileInfo().IsDir() {
			os.MkdirAll(fpath, os.ModePerm)
			continue
		}

		if err := os.MkdirAll(filepath.Dir(fpath), os.ModePerm); err != nil {
			return err
		}

		outFile, err := os.OpenFile(fpath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
		if err != nil {
			return err
		}
		defer outFile.Close()

		rc, err := f.Open()
		if err != nil {
			return err
		}

		_, err = io.Copy(outFile, rc)
		rc.Close() // Cerrar el archivo inmediatamente después de copiarlo

		if err != nil {
			return err
		}
	}

	return nil
}

// removeFileWithRetry intenta eliminar un archivo varias veces en Windows
func removeFileWithRetry(path string) error {
	maxRetries := 5
	for i := 0; i < maxRetries; i++ {
		err := os.Remove(path)
		if err == nil {
			return nil
		}

		// En Windows, puede que el archivo aún esté en uso, intentamos de nuevo
		if runtime.GOOS == "windows" {
			fmt.Printf("⌛ Reintentando eliminar el ZIP (%d/%d)...\n", i+1, maxRetries)
			time.Sleep(500 * time.Millisecond)
		} else {
			return err // En otros sistemas, no reintentamos
		}
	}
	return fmt.Errorf("no se pudo eliminar el archivo ZIP tras %d intentos", maxRetries)
}
