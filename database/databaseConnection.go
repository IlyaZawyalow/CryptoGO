package database

import (
	"cryptoGo/depth"
	"database/sql"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/adshao/go-binance/v2"
	_ "github.com/lib/pq"
)

func AddNewRow(kline binance.WsKline, depth depth.DepthInfo, tableName string, db *sql.DB) {
	// Преобразование временной метки в строку
	timestamp := time.Unix(0, kline.StartTime*int64(time.Millisecond)).Format("2006-01-02 15:04:05")

	// Подготовленный запрос с параметрами
	klineQuery := fmt.Sprintf(
		"INSERT INTO %s (timestamp, Asset, Open, High, Low, Close, Volume, TradeNum, ActiveBuyVolume, Count, MaxDid, MinAsk, Imbalance, BuyWAP, SellWAP) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15)",
		tableName,
	)

	// Выполнение запроса
	_, err := db.Exec(klineQuery, timestamp, kline.Symbol, kline.Open, kline.High, kline.Low, kline.Close, kline.Volume, kline.TradeNum, kline.ActiveBuyVolume, depth.Count, depth.MaxDid, depth.MinAsk, depth.Imbalance, depth.BuyWAP, depth.SellWAP)
	if err != nil {
		fmt.Println(err)
	}

}

func ConnectToDatabase() (*sql.DB, error) {
	userName := os.Getenv("UserName")
	bdName := os.Getenv("DbName")
	password := os.Getenv("Password")
	sslmode := os.Getenv("Sslmode")
	port := os.Getenv("DbPort")

	// Добавляем порт в строку подключения
	connStr := fmt.Sprintf("user=%s password=%s dbname=%s sslmode=%s port=%s", userName, password, bdName, sslmode, port)

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}

	err = db.Ping()
	if err != nil {
		log.Fatal(err)
	}

	return db, nil
}
