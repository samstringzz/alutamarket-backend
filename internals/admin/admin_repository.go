package admin

import (
	"context"
	"log"
	"os"
	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type repository struct {
	db *gorm.DB
}
func NewRepository() Repository {
	if err := godotenv.Load(); err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}
	dbURI := os.Getenv("DB_URI")

	// Initialize the database connection
	db, err := gorm.Open(postgres.Open(dbURI), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	return &repository{
		db: db,
	}
}

func (r *repository) Login(ctx context.Context, input *LoginAdminReq) (*LoginAdminRes,error){
	
	panic("")
}

func (r *repository) CreateAdmin(ctx context.Context, input *Admin) (*LoginAdminRes,error){
	
	panic("")
}

func (r *repository) GetAdmin(ctx context.Context, id uint8) (*Admin,error){
	
	panic("")
}

func (r *repository) GetAdmins(ctx context.Context) ([]*Admin,error){
	
	panic("")
}

func (r *repository) ApproveProduct(){
	

}
func (r *repository) MessageUser(){
	

}
