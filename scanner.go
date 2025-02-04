package main

import (
	"crypto/sha256"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sync"
	"sync/atomic"
)

// FileStatus representa el estado de un archivo comparado
type FileStatus struct {
	Path   string
	Status string // "OK", "MODIFIED", "MISSING", "EXTRA"
}

// CompareFiles compara la instalación local con la versión oficial
func CompareFiles(moodlePath, extractedPath string) ([]FileStatus, error) {
	var results []FileStatus
	resultsChan := make(chan FileStatus, 100)
	var count int32
	var totalFiles int32
	var wg sync.WaitGroup
	semaphore := make(chan struct{}, 10)

	fmt.Println("📊 Contando archivos en ambas ubicaciones...")

	ignorePatterns := readIgnoreList()

	// Contar archivos en la instalación local
	err := filepath.Walk(moodlePath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return nil
		}

		relPath, _ := filepath.Rel(moodlePath, path)
		relPath = filepath.ToSlash(relPath)

		// Ignorar directorios completos
		if shouldIgnore(relPath, ignorePatterns) {
			fmt.Println("🚫 Ignorando:", relPath)
			if info.IsDir() {
				return filepath.SkipDir
			}
			return nil
		}

		atomic.AddInt32(&totalFiles, 1)
		return nil
	})
	if err != nil {
		return nil, err
	}

	// Contar archivos en el ZIP
	err = filepath.Walk(extractedPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return nil
		}

		relPath, _ := filepath.Rel(extractedPath, path)
		relPath = filepath.ToSlash(relPath)

		// Ignorar directorios completos
		if shouldIgnore(relPath, ignorePatterns) {
			fmt.Println("🚫 Ignorando:", relPath)
			if info.IsDir() {
				return filepath.SkipDir
			}
			return nil
		}

		atomic.AddInt32(&totalFiles, 1)
		return nil
	})
	if err != nil {
		return nil, err
	}

	fmt.Printf("📁 Total de archivos a comparar: %d\n", totalFiles)
	fmt.Println("📊 Iniciando comparación de archivos...")

	// Comparar archivos de la instalación local con el ZIP
	go func() {
		filepath.Walk(moodlePath, func(localFile string, info os.FileInfo, err error) error {
			if err != nil {
				return nil
			}

			relPath, _ := filepath.Rel(moodlePath, localFile)
			relPath = filepath.ToSlash(relPath)

			// Ignorar directorios completos
			if shouldIgnore(relPath, ignorePatterns) {
				fmt.Println("🚫 Ignorando:", relPath)
				if info.IsDir() {
					return filepath.SkipDir
				}
				return nil
			}

			originalFile := filepath.Join(extractedPath, relPath)

			wg.Add(1)
			semaphore <- struct{}{}
			go func(relPath, originalFile, localFile string, isDir bool) {
				defer wg.Done()
				defer func() { <-semaphore }()

				if isDir {
					if !fileExists(originalFile) {
						resultsChan <- FileStatus{Path: relPath, Status: "EXTRA"}
					}
				} else {
					if !fileExists(originalFile) {
						resultsChan <- FileStatus{Path: relPath, Status: "EXTRA"}
					} else if !compareFileHash(localFile, originalFile) {
						resultsChan <- FileStatus{Path: relPath, Status: "MODIFIED"}
					} else {
						resultsChan <- FileStatus{Path: relPath, Status: "OK"}
					}
				}
				atomic.AddInt32(&count, 1)
			}(relPath, originalFile, localFile, info.IsDir())

			return nil
		})

		wg.Wait()
		close(resultsChan)
	}()

	// Recoger los resultados del canal
	for result := range resultsChan {
		results = append(results, result)
	}

	fmt.Println("✅ Comparación de archivos completada.")
	return results, nil
}

// compareFileHash compara los hashes SHA256 de dos archivos
func compareFileHash(file1, file2 string) bool {
	hash1, err1 := fileHash(file1)
	hash2, err2 := fileHash(file2)

	if err1 != nil || err2 != nil {
		return false
	}
	return hash1 == hash2
}

// fileHash calcula el hash SHA256 de un archivo
func fileHash(filename string) (string, error) {
	f, err := os.Open(filename)
	if err != nil {
		return "", err
	}
	defer f.Close()

	h := sha256.New()
	if _, err := io.Copy(h, f); err != nil {
		return "", err
	}

	return fmt.Sprintf("%x", h.Sum(nil)), nil
}

// fileExists verifica si un archivo o carpeta existe
func fileExists(path string) bool {
	_, err := os.Stat(path)
	return !os.IsNotExist(err)
}

// getReportsDir devuelve la ruta del directorio de informes
func getReportsDir() string {
	reportsDir := "./reports"
	if _, err := os.Stat(reportsDir); os.IsNotExist(err) {
		os.Mkdir(reportsDir, os.ModePerm)
	}
	return reportsDir
}

// saveReport guarda los resultados en la carpeta ./reports/
func saveReport(results []FileStatus) error {
	reportsDir := getReportsDir()

	file, err := os.Create(filepath.Join(reportsDir, "security_report.txt"))
	if err != nil {
		return fmt.Errorf("error creando el informe: %w", err)
	}
	defer file.Close()

	alertFile, err := os.Create(filepath.Join(reportsDir, "security_report_alert.txt"))
	if err != nil {
		return fmt.Errorf("error creando el informe de alertas: %w", err)
	}
	defer alertFile.Close()

	for _, result := range results {
		line := fmt.Sprintf("[%s] %s\n", result.Status, result.Path)
		file.WriteString(line)

		// Guardar alertas
		if result.Status == "MODIFIED" || result.Status == "EXTRA" || result.Status == "MISSING" {
			alertFile.WriteString(line)
		}
	}

	fmt.Println("📄 Informe completo:", filepath.Join(reportsDir, "security_report.txt"))
	fmt.Println("⚠️ Informe de alertas:", filepath.Join(reportsDir, "security_report_alert.txt"))
	return nil
}
