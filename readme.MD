# Cara Menjalankan Proyek

## Prasyarat

Pastikan Anda telah menginstal perangkat lunak berikut:

- Go: https://go.dev/doc/install
- Docker: https://docs.docker.com/get-docker/
- Docker Compose: https://docs.docker.com/compose/install/

## Langkah-langkah

### 1. Clone repositori ini:

git clone <URL_REPOSITORI_ANDA>
cd <NAMA_PROYEK>

### 2. Konfigurasi Database

Buat file `.env` di root direktori proyek dan isi dengan konten berikut:

DB_USER=user
DB_PASSWORD=password
DB_HOST=localhost
DB_PORT=5432
DB_NAME=your_db_name
JWT_SECRET=supersecretkey

Pastikan nilai-nilai ini sesuai dengan konfigurasi yang ada di file `docker-compose.yml` Anda.

### 3. Jalankan Database dengan Docker Compose

Untuk memulai database menggunakan Docker Compose, jalankan perintah berikut:

docker-compose up -d

### 4. Install Dependensi Go

Setelah database siap, instal dependensi Go dengan perintah:

go mod tidy

### 5. Jalankan Aplikasi

Untuk menjalankan aplikasi, gunakan perintah:

go run main.go

Aplikasi akan berjalan di `http://localhost:3000`.

## Struktur Proyek

Berikut adalah struktur folder proyek:

.
├── main.go                # Entry point aplikasi
├── go.mod                 # File modul Go
├── go.sum                 # Checksum file modul
├── .env                   # Variabel lingkungan
├── .gitignore             # File yang diabaikan Git
├── docker-compose.yml     # Konfigurasi Docker Compose
└── internal/
    ├── models/            # Definisi model GORM
    ├── repository/        # Logika akses data (Repository)
    ├── service/           # Logika bisnis (Service)
    └── handler/           # Logika permintaan HTTP (Handler)
