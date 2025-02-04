package main

import (
	"flag"
	"fmt"
	"os"
)

// Valor por defecto para la instalación de Moodle
const defaultMoodlePath = "/var/www"

func main() {
	// Definir comandos
	analyzeCmd := flag.NewFlagSet("analyze", flag.ExitOnError)
	fixCmd := flag.NewFlagSet("fix", flag.ExitOnError)

	// Definir flags para cada comando
	moodlePath := analyzeCmd.String("path", defaultMoodlePath, "Ruta de la instalación de Moodle")
	analyzeZip := analyzeCmd.String("zip", "", "Ruta a un ZIP local de Moodle (opcional)")

	fixPath := fixCmd.String("path", defaultMoodlePath, "Ruta de la instalación de Moodle")
	fixZip := fixCmd.String("zip", "", "Ruta a un ZIP local de Moodle (opcional)")

	// Verificar que el usuario ha pasado un comando
	if len(os.Args) < 2 {
		fmt.Println("Uso: moodle-security <comando> [opciones]")
		fmt.Println("Comandos disponibles:")
		fmt.Println("  analyze  - Analiza la instalación de Moodle")
		fmt.Println("  fix      - Repara archivos modificados")
		os.Exit(1)
	}

	// Procesar el comando
	switch os.Args[1] {
	case "analyze":
		analyzeCmd.Parse(os.Args[2:])
		fmt.Printf("Analizando Moodle en: %s\n", *moodlePath)

		config, err := LoadConfig(*moodlePath)
		if err != nil {
			fmt.Println("❌ Error al leer config.php:", err)
			os.Exit(1)
		}

		fmt.Println("✔️ Configuración de Moodle obtenida")

		version, err := GetMoodleVersion(config.DB)
		if err != nil {
			fmt.Println("❌ Error al obtener la versión de Moodle:", err)
			os.Exit(1)
		}
		fmt.Println("✔️ Versión de Moodle:", version)

		extractedPath, err := DownloadMoodle(version, *analyzeZip) // ✅ Ahora usa --zip correctamente
		if err != nil {
			fmt.Println("❌ Error al obtener Moodle:", err)
			os.Exit(1)
		}

		fmt.Println("🔍 Comparando archivos...")
		results, err := CompareFiles(*moodlePath, extractedPath)
		if err != nil {
			fmt.Println("❌ Error al comparar archivos:", err)
			os.Exit(1)
		}

		// Guardar el informe
		if err := saveReport(results); err != nil {
			fmt.Println("❌ Error al guardar el informe:", err)
			os.Exit(1)
		}

	case "fix":
		fixCmd.Parse(os.Args[2:])
		fmt.Printf("Reparando Moodle en: %s\n", *fixPath)

		// ✅ `fixZip` ya está disponible para su uso en `fix`
		fmt.Printf("Usando ZIP para reparar: %s\n", *fixZip)

	case "clean":
		fmt.Println("🧹 Ejecutando limpieza...")
		cleanTempAndReports()

	default:
		fmt.Println("Comando no reconocido.")
		os.Exit(1)
	}
}
