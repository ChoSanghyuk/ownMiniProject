package model

import (
	"time"

	"gorm.io/datatypes"
	"gorm.io/gorm"
)

type Fund struct {
	ID   uint
	Name string
}

type Asset struct {
	ID        uint
	Name      string
	Category  uint
	Currency  string
	Top       float64
	Bottom    float64
	SellPrice float64
	BuyPrice  float64
	Path      string
}

type Invest struct {
	ID      uint
	FundID  uint
	Fund    Fund
	AssetID uint
	Asset   Asset
	Price   float64
	Count   int
	gorm.Model
}

type InvestSummary struct {
	ID      uint
	FundID  uint
	Fund    Fund
	AssetID uint
	Asset   Asset
	Sum     float64
}

type Market struct {
	CreatedAt datatypes.Date `gorm:"primaryKey"`
	Status    uint
}

type DailyIndex struct {
	CreatedAt      datatypes.Date `gorm:"primaryKey"`
	FearGreedIndex uint
	NasDaq         float64
}

type CliIndex struct {
	CreatedAt datatypes.Date `gorm:"primaryKey"`
	Index     float64
}

type Sample struct {
	ID   uint `gorm:"primaryKey"`
	Date datatypes.Date
	Time time.Time
	gorm.Model
}
