package utils

import (
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"time"
)

// SaveUploadedFile menyimpan file yang diunggah ke direktori public/uploads
// dan mengembalikan nama file baru yang dibuat.
func SaveUploadedFile(file *multipart.FileHeader) (string, error) {
	// Buat direktori "public/uploads" jika belum ada
	uploadDir := "./public/uploads"
	if err := os.MkdirAll(uploadDir, 0755); err != nil {
		return "", fmt.Errorf("gagal membuat direktori upload: %w", err)
	}

	// Buat nama file unik dengan timestamp untuk menghindari duplikasi
	filename := fmt.Sprintf("%d-%s", time.Now().Unix(), filepath.Base(file.Filename))
	filePath := filepath.Join(uploadDir, filename)

	// Buka file yang diunggah
	src, err := file.Open()
	if err != nil {
		return "", fmt.Errorf("gagal membuka file yang diunggah: %w", err)
	}
	defer src.Close()

	// Buat file baru di direktori tujuan
	dst, err := os.Create(filePath)
	if err != nil {
		return "", fmt.Errorf("gagal membuat file tujuan: %w", err)
	}
	defer dst.Close()

	// Salin isi file
	if _, err := io.Copy(dst, src); err != nil {
		return "", fmt.Errorf("gagal menyalin file: %w", err)
	}

	return filename, nil
}
