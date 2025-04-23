package main

import (
	"atomic/internal/db"
	"atomic/internal/ws"
	_ "atomic/internal/ws/command/collection"
	_ "atomic/internal/ws/command/ai" 
	"log"

	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("⚠️ Không tìm thấy file .env, dùng biến môi trường hệ thống")
	}
	db.Connect() // Kết nối đến MongoDB

	ws.Connect()

}
