package depth

import (
	"fmt"
	"strconv"

	"github.com/adshao/go-binance/v2/common"
)

type DepthInfo struct {
	bidsList  []common.PriceLevel
	asksList  []common.PriceLevel
	Count     int
	MaxDid    float64
	MinAsk    float64
	Imbalance float64
	BuyWAP    float64
	SellWAP   float64
}

func (d *DepthInfo) SetDepth(bidsList, asksList []common.PriceLevel) {
	d.bidsList = bidsList
	d.asksList = asksList
	d.SetCount()
	d.SetMaxDid()
	d.SetMinAsk()
	d.SetImbalance()
	d.SetWAP()
}

// Установка количества заявок
func (d *DepthInfo) SetCount() {
	d.Count = len(d.bidsList) + len(d.asksList)
}

// Установка максимальной цены покупки
func (d *DepthInfo) SetMaxDid() {
	if len(d.bidsList) == 0 {
		d.MaxDid = 0
		return
	}

	// Инициализация maxDid
	maxDid, _ := strconv.ParseFloat(d.bidsList[0].Price, 64)

	// Поиск максимальной цены
	for _, bid := range d.bidsList {
		price, _ := strconv.ParseFloat(bid.Price, 64)

		if price > maxDid {
			maxDid = price
		}
	}

	d.MaxDid = maxDid
}

// Установка минимальной цены продажи
func (d *DepthInfo) SetMinAsk() {
	if len(d.asksList) == 0 {
		d.MinAsk = 0
		return
	}

	// Инициализация minAsk
	minAsk, _ := strconv.ParseFloat(d.asksList[0].Price, 64)

	// Поиск минимальной цены
	for _, ask := range d.asksList {
		price, _ := strconv.ParseFloat(ask.Price, 64)

		if price < minAsk {
			minAsk = price
		}
	}

	d.MinAsk = minAsk
}

// Установка индекса дисбаланса
func (d *DepthInfo) SetImbalance() {
	var totalBuyVolume, totalSellVolume float64

	// Суммируем объемы заявок на продажу
	if len(d.asksList) > 0 {
		for _, ask := range d.asksList {
			quantity, err := strconv.ParseFloat(ask.Quantity, 64)
			if err != nil {
				fmt.Println("Error parsing ask quantity:", err)
				continue // Игнорируем ошибочные заявки
			}
			totalSellVolume += quantity
		}
	}

	// Суммируем объемы заявок на покупку
	if len(d.bidsList) > 0 {
		for _, bid := range d.bidsList {
			quantity, err := strconv.ParseFloat(bid.Quantity, 64)
			if err != nil {
				fmt.Println("Error parsing bid quantity:", err)
				continue // Игнорируем ошибочные заявки
			}
			totalBuyVolume += quantity
		}
	}

	// Вычисляем индекс дисбаланса
	if totalBuyVolume+totalSellVolume == 0 {
		d.Imbalance = 0.0
		return
	}

	d.Imbalance = (totalBuyVolume - totalSellVolume) / (totalBuyVolume + totalSellVolume)
}

// Установка WAP
func (d *DepthInfo) SetWAP() {
	var totalBuyVolume, totalSellVolume float64
	var totalBuyPrice, totalSellPrice float64

	// Суммируем объемы и цены заявок на покупку
	for _, bid := range d.bidsList {
		price, err := strconv.ParseFloat(bid.Price, 64)
		if err != nil {
			fmt.Println("Error parsing bid price:", err)
			continue
		}
		quantity, err := strconv.ParseFloat(bid.Quantity, 64)
		if err != nil {
			fmt.Println("Error parsing bid quantity:", err)
			continue
		}

		totalBuyVolume += quantity
		totalBuyPrice += price * quantity // Взвешиваем цену на объем
	}

	// Суммируем объемы и цены заявок на продажу
	for _, ask := range d.asksList {
		price, err := strconv.ParseFloat(ask.Price, 64)
		if err != nil {
			fmt.Println("Error parsing ask price:", err)
			continue
		}
		quantity, err := strconv.ParseFloat(ask.Quantity, 64)
		if err != nil {
			fmt.Println("Error parsing ask quantity:", err)
			continue
		}

		totalSellVolume += quantity
		totalSellPrice += price * quantity // Взвешиваем цену на объем
	}

	// Вычисляем BuyWAP
	if totalBuyVolume > 0 {
		d.BuyWAP = totalBuyPrice / totalBuyVolume
	} else {
		d.BuyWAP = 0
	}

	// Вычисляем SellWAP
	if totalSellVolume > 0 {
		d.SellWAP = totalSellPrice / totalSellVolume
	} else {
		d.SellWAP = 0
	}
}
