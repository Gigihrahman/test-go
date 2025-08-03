package models

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	ID           uint       `gorm:"primaryKey;autoIncrement"`
	Nama         string     `gorm:"type:varchar(255)"`
	Email        string     `gorm:"type:varchar(255);unique"`
	KataSandi    string     `gorm:"type:varchar(255)"`
	NoTelp       string     `gorm:"type:varchar(255);unique"`
	TanggalLahir *time.Time `gorm:"type:date"`
	JenisKelamin string     `gorm:"type:varchar(255)"`
	Tentang      string     `gorm:"type:text"`
	Pekerjaan    string     `gorm:"type:varchar(255)"`
	IDProvinsi   uint
	IDKota       uint
	IsAdmin      bool `gorm:"type:boolean"`
	CreatedAt    time.Time
	UpdatedAt    time.Time

	Alamat []Alamat `gorm:"foreignKey:IDUser"`
	Trx    []Trx    `gorm:"foreignKey:IDUser"`
}

type Alamat struct {
	gorm.Model
	ID           uint `gorm:"primaryKey;autoIncrement"`
	IDUser       uint
	JudulAlamat  string `gorm:"type:varchar(255)"`
	NamaPenerima string `gorm:"type:varchar(255)"`
	NoTelp       string `gorm:"type:varchar(255)"`
	DetailAlamat string `gorm:"type:varchar(255)"`
	CreatedAt    time.Time
	UpdatedAt    time.Time

	User User `gorm:"foreignKey:IDUser"`
}

type Toko struct {
	gorm.Model
	ID          uint `gorm:"primaryKey;autoIncrement"`
	IDUser      uint
	NamaToko    string `gorm:"type:varchar(255)"`
	URLFotoToko string `gorm:"type:varchar(255)"`
	CreatedAt   time.Time
	UpdatedAt   time.Time

	User       User         `gorm:"foreignKey:IDUser"`
	Product    []Product    `gorm:"foreignKey:IDToko"`
	ProductLog []ProductLog `gorm:"foreignKey:IDToko"`
}

type Category struct {
	gorm.Model
	ID           uint   `gorm:"primaryKey;autoIncrement"`
	NamaCategory string `gorm:"type:varchar(255)"`
	CreatedAt    time.Time
	UpdatedAt    time.Time

	Product    []Product    `gorm:"foreignKey:IDCategory"`
	ProductLog []ProductLog `gorm:"foreignKey:IDCategory"`
}

type Product struct {
	gorm.Model
	ID            uint `gorm:"primaryKey;autoIncrement"`
	IDToko        uint
	IDCategory    uint
	NamaProduct   string `gorm:"type:varchar(255)"`
	Slug          string `gorm:"type:varchar(255)"`
	HargaReseller int
	HargaKonsumen int
	Stok          int
	Deskripsi     string `gorm:"type:text"`
	CreatedAt     time.Time
	UpdatedAt     time.Time

	Toko         Toko           `gorm:"foreignKey:IDToko"`
	Category     Category       `gorm:"foreignKey:IDCategory"`
	ProductPhoto []ProductPhoto `gorm:"foreignKey:ProductID"`
	DetailTrx    []DetailTrx    `gorm:"foreignKey:ProductID"`
}

type ProductPhoto struct {
	gorm.Model
	ID        uint `gorm:"primaryKey;autoIncrement"`
	ProductID uint
	URL       string `gorm:"type:varchar(255)"`
	CreatedAt time.Time
	UpdatedAt time.Time

	Product Product `gorm:"foreignKey:ProductID"`
}

type ProductLog struct {
	gorm.Model
	ID            uint `gorm:"primaryKey;autoIncrement"`
	ProductID     uint
	IDToko        uint
	IDCategory    uint
	NamaProduct   string `gorm:"type:varchar(255)"`
	Slug          string `gorm:"type:varchar(255)"`
	HargaReseller int
	HargaKonsumen int
	Deskripsi     string `gorm:"type:text"`
	CreatedAt     time.Time
	UpdatedAt     time.Time

	Product  Product  `gorm:"foreignKey:ProductID"`
	Toko     Toko     `gorm:"foreignKey:IDToko"`
	Category Category `gorm:"foreignKey:IDCategory"`
}

type Trx struct {
	gorm.Model
	ID               uint `gorm:"primaryKey;autoIncrement"`
	IDUser           uint
	AlamatPengiriman uint
	KodeInvoice      string `gorm:"type:varchar(255)"`
	MethodBayar      string `gorm:"type:varchar(255)"`
	HargaTotal       int
	CreatedAt        time.Time
	UpdatedAt        time.Time

	User      User        `gorm:"foreignKey:IDUser"`
	Alamat    Alamat      `gorm:"foreignKey:AlamatPengiriman"`
	DetailTrx []DetailTrx `gorm:"foreignKey:IDTrx"`
}

type DetailTrx struct {
	gorm.Model
	ID         uint `gorm:"primaryKey;autoIncrement"`
	IDTrx      uint
	ProductID  uint
	IDToko     uint
	Kuantitas  int
	HargaTotal int
	CreatedAt  time.Time
	UpdatedAt  time.Time

	Trx     Trx        `gorm:"foreignKey:IDTrx"`
	Product ProductLog `gorm:"foreignKey:ProductID"`
	Toko    Toko       `gorm:"foreignKey:IDToko"`
}

type DetailTrxPayload struct {
	ProductID uint `json:"product_id"`
	Kuantitas int  `json:"kuantitas"`
}

type TrxPayload struct {
	MethodBayar string             `json:"method_bayar"`
	AlamatKirim uint               `json:"alamat_kirim"`
	DetailTrx   []DetailTrxPayload `json:"detail_trx"`
}
