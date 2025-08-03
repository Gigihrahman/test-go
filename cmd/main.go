package main

import (
	"fmt"
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"

	category_handler "test-rakamin/internal/handler/category"
	product_handler "test-rakamin/internal/handler/product"
	toko_handler "test-rakamin/internal/handler/toko"
	trx_handler "test-rakamin/internal/handler/trx"
	user_handler "test-rakamin/internal/handler/user"
	"test-rakamin/internal/models"
	category_repository "test-rakamin/internal/repository/category"
	product_repository "test-rakamin/internal/repository/product"
	product_photo_repository "test-rakamin/internal/repository/product_photo"
	toko_repository "test-rakamin/internal/repository/toko"
	trx_repository "test-rakamin/internal/repository/trx"
	user_repository "test-rakamin/internal/repository/user"
	category_service "test-rakamin/internal/service/category"
	product_service "test-rakamin/internal/service/product"
	toko_service "test-rakamin/internal/service/toko"
	trx_service "test-rakamin/internal/service/trx"
	user_service "test-rakamin/internal/service/user"
	"test-rakamin/pkg/internalsql"
)

func main() {

	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	dbName := os.Getenv("DB_NAME")

	dataSourceName := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", dbUser, dbPassword, dbHost, dbPort, dbName)

	db, err := internalsql.Connect(dataSourceName)
	if err != nil {
		log.Fatalf("Gagal terhubung ke database: %v", err)
	}

	err = db.AutoMigrate(
		&models.User{},
		&models.Alamat{},
		&models.Toko{},
		&models.Category{},
		&models.Product{},
		&models.ProductPhoto{},
		&models.ProductLog{},
		&models.Trx{},
		&models.DetailTrx{},
	)
	if err != nil {
		log.Fatalf("Gagal migrasi database: %v", err)
	}

	app := fiber.New()

	userRepo := user_repository.NewUserRepository(db)
	categoryRepo := category_repository.NewCategoryRepository(db)
	tokoRepo := toko_repository.NewTokoRepository(db)
	productRepo := product_repository.NewProductRepository(db)
	productPhotoRepo := product_photo_repository.NewProductPhotoRepository(db)
	trxRepo := trx_repository.NewTrxRepository(db)

	userService := user_service.NewUserService(userRepo)
	categoryService := category_service.NewCategoryService(categoryRepo)
	tokoService := toko_service.NewTokoService(tokoRepo)
	productService := product_service.NewProductService(productRepo, productPhotoRepo)
	trxService := trx_service.NewTrxService(trxRepo, productRepo)

	userHandler := user_handler.NewUserHandler(userService)
	categoryHandler := category_handler.NewCategoryHandler(categoryService)
	tokoHandler := toko_handler.NewTokoHandler(tokoService)
	productHandler := product_handler.NewProductHandler(productService)
	trxHandler := trx_handler.NewTrxHandler(trxService)

	userHandler.RegisterRoutes(app)
	categoryHandler.RegisterRoutes(app)
	tokoHandler.RegisterRoutes(app)
	productHandler.RegisterRoutes(app)
	trxHandler.RegisterRoutes(app)

	log.Println("Server berjalan di http://localhost:3000")
	log.Fatal(app.Listen(":3000"))
}
