package main

import (
	"cryptoGo/database"
	"cryptoGo/parser"
	"fmt"
	"log"
	"time"

	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	db, err := database.ConnectToDatabase()

	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	fmt.Print("Успешное подключение к базе данных!\n")

	p := parser.InitParser(db)

	p.StartParser()
	time.Sleep(time.Second * 400)
	p.StopParser()

}
