# Moodle Security 🔒

**Moodle Security** es una herramienta en **Go** diseñada para analizar instalaciones de **Moodle** y detectar archivos sospechosos, modificados o faltantes.

### 🚀 Características:
✅ Analiza la instalación de Moodle y genera un informe de seguridad.  
✅ Compara la instalación con una versión oficial descargada desde **GitHub** o un archivo **ZIP local**.  
✅ Permite omitir archivos y carpetas específicas mediante `ignore.conf`.  
✅ Elimina archivos temporales y reportes con `cleaner`.  
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

### 🗑️ **Limpiar archivos temporales y reportes**
```sh
moodle-security clean
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
│── github.go              # Descarga de Moodle desde GitHub
│── config.go              # Manejo de configuración
│── README.md              # Este archivo 📖
```