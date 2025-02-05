# Moodle Security 🔒

**Moodle Security** es una herramienta en **Go** diseñada para analizar instalaciones de **Moodle** y detectar archivos sospechosos, modificados o faltantes.

### 🚀 Características:
✅ Analiza la instalación de Moodle y genera un informe de seguridad.  
✅ Compara la instalación con una versión oficial descargada desde **GitHub** o un archivo **ZIP local**.  
✅ Permite omitir archivos y carpetas específicas mediante `ignore.conf`.  
✅ Elimina archivos temporales y reportes con `clean`.  
✅ **(En desarrollo)** Función `fix` para reparar archivos dañados automáticamente.

---

## 📥 Instalación

### 🔧 **Requisitos**
- **Go 1.23+**  
- **Acceso a internet** *(si deseas descargar Moodle desde GitHub)*  

### 📦 **Descargar el código**
```sh
git clone https://github.com/jhernandezj/moodle-security.git
cd moodle-security
go build -o moodle-security
```

---

## 🚀 **Uso**

### 🔍 **Analizar Moodle**
Compara los archivos de Moodle con una versión oficial y genera un informe.  
```sh
moodle-security analyze --path /ruta/de/moodle
```
Ejemplo en Windows:
```sh
moodle-security.exe analyze --path \"C:\localServer\moodle\"
```
Ejemplo con un **ZIP local**:
```sh
moodle-security analyze --path /ruta/de/moodle --zip /ruta/al/moodle.zip
```

### 🛠️ **Reparar archivos modificados o faltantes**
*(En desarrollo)*
```sh
moodle-security fix --path /ruta/de/moodle
```

### 🗑️ **Limpiar archivos temporales y reportes**
```sh
moodle-security cleaner
```

### ⚙️ **Opciones disponibles**
```sh
moodle-security --help
```

---

## 📂 **Estructura del Proyecto**
```
moodle-security/
│── temp/                  # Descargas temporales (Moodle ZIP y extraído)
│── reports/               # Reportes de análisis
│── ignore.conf            # Archivos y carpetas a ignorar
│── main.go                # Código principal
│── scanner.go             # Lógica de análisis
│── fixer.go               # Implementación de la reparación
│── database.go            # Funciones para extraer información desde la DB
│── github.go              # Descarga de Moodle desde GitHub
│── config.go              # Manejo de configuración
│── utils                  # Utilidades para manejar archivos y carpetas
│── README.md              # Este archivo 📖
```

---

## ⚠️ **Informe de Seguridad**
Los resultados se guardan en la carpeta `reports/`:

📄 **Reporte completo**  
```
reports/security_report.txt
```
⚠️ **Archivos sospechosos/modificados**  
```
reports/security_report_alert.txt
```

Ejemplo de contenido en `security_report_alert.txt`:
```
[MISSING] theme/boost/config.php
[MODIFIED] lib/setup.php
[EXTRA] local/backdoor.php
```

---

## 🔥 **Configuración de `ignore.conf`**
Para excluir archivos/carpeta del análisis, usa `ignore.conf`.  

Ejemplo:
```
# Ignorar cualquier carpeta .git en cualquier nivel
**/.git/

# Ignorar carpetas específicas y su contenido
user_privacy_utilities/
backup/
cache/*

# Ignorar archivos individuales
config.php
```

---

## 🏗 **Próximos pasos**
🔜 **Implementar `fix` para reparar archivos dañados automáticamente.**  
🔜 **Agregar logs más detallados en JSON.**  

---

## 👨‍💻 **Contribuir**
Si deseas mejorar Moodle Security, ¡envía un **PR** o abre un **issue** en GitHub!  
📌 **Repositorio:** [GitHub](https://github.com/jhernandezj/moodle-security)  

---

## 📜 **Licencia**
MIT License © 2025 