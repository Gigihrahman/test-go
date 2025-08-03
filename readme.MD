Cara Menjalankan Proyek

### Prasyarat

Pastikan Anda telah menginstal:

- [Go](https://go.dev/doc/install)

- [Docker](https://docs.docker.com/get-docker/)

- [Docker Compose](https://docs.docker.com/compose/install/)

### Langkah-langkah

1. **Clone repositori ini:**

git clone <URL_REPOSITORI_ANDA>
cd <NAMA_PROYEK>

2. **Konfigurasi Database:**
   Buat file `.env` di root direktori proyek dengan konten berikut:

DB_USER=user
DB_PASSWORD=password
DB_HOST=localhost
DB_PORT=5432
DB_NAME=your_db_name
JWT_SECRET=supersecretkey

Pastikan nilai-nilai ini cocok dengan konfigurasi di `docker-compose.yml` Anda.

3. **Jalankan Database dengan Docker Compose:**

docker-compose up -d

4. **Install dependensi Go:**

go mod tidy

5. **Jalankan Aplikasi:**

go run main.go

Aplikasi akan berjalan di `http://localhost:3000`.

## Struktur Proyek

.
├── main.go # Entry point aplikasi
├── go.mod # File modul Go
├── go.sum # Checksum file modul
├── .env # Variabel lingkungan
├── .gitignore # File yang diabaikan Git
├── docker-compose.yml # Konfigurasi Docker Compose
└── internal/
├── models/ # Definisi model GORM
├── repository/ # Logika akses data (Repository)
├── service/ # Logika bisnis (Service)
└── handler/ # Logika permintaan HTTP (Handler)
