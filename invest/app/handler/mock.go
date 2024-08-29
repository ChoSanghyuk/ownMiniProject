package handler

import (
	"fmt"
	m "invest/model"
	"time"

	"gorm.io/datatypes"
)

/***************************** Asset ***********************************/
type AssetRetrieverMock struct {
	err error
}

func (mock AssetRetrieverMock) RetrieveAssetList() ([]map[string]interface{}, error) {
	fmt.Println("RetrieveAssetList Called")

	if mock.err != nil {
		return nil, mock.err
	}
	return []map[string]interface{}{
		{
			"id":   1,
			"name": "비트코인",
		},
		{
			"id":   2,
			"name": "TigerS&P500",
		},
	}, nil
}
func (mock AssetRetrieverMock) RetrieveAsset(id uint) (*m.Asset, error) {
	fmt.Println("RetrieveAsset Called")

	if mock.err != nil {
		return nil, mock.err
	}
	return &m.Asset{
		ID:        1,
		Name:      "비트코인",
		Category:  6,
		Currency:  "USD",
		Top:       9800,
		Bottom:    6800,
		SellPrice: 8800,
		BuyPrice:  7800,
		Path:      "",
	}, nil
}

func (mock AssetRetrieverMock) RetrieveAssetHist(id uint) ([]m.Invest, error) {
	fmt.Println("RetrieveAssetHist Called")

	if mock.err != nil {
		return nil, mock.err
	}
	return []m.Invest{
		{
			ID:      1,
			FundID:  3,
			AssetID: 1,
			Price:   7800,
			Count:   5,
		},
	}, nil
}

type AssetInfoSaverMock struct {
	err error
}

func (mock AssetInfoSaverMock) SaveAssetInfo(name string, category uint, currency string, top float64, bottom float64, selPrice float64, buyPrice float64, path string) error {
	fmt.Println("SaveAssetInfo Called")

	if mock.err != nil {
		return mock.err
	}
	return nil
}
func (mock AssetInfoSaverMock) UpdateAssetInfo(name string, category uint, currency string, top float64, bottom float64, selPrice float64, buyPrice float64, path string) error {
	fmt.Println("UpdateAssetInfo Called")

	if mock.err != nil {
		return mock.err
	}
	return nil
}
func (mock AssetInfoSaverMock) DeleteAssetInfo(id uint) error {
	fmt.Println("DeleteAssetInfo Called")

	if mock.err != nil {
		return mock.err
	}
	return nil
}

/***************************** Fund ***********************************/
type FundRetrieverMock struct {
	err error
}

func (mock FundRetrieverMock) RetreiveFundsSummary() ([]m.InvestSummary, error) {
	fmt.Println("RetreiveFundsSummary Called")

	if mock.err != nil {
		return nil, mock.err
	}
	return []m.InvestSummary{
		{
			ID:      1,
			FundID:  1,
			AssetID: 1,
			Sum:     568210,
		},
	}, nil
}
func (mock FundRetrieverMock) RetreiveFundSummaryById(id uint) ([]m.InvestSummary, error) {
	fmt.Println("RetreiveFundSummaryById Called")

	if mock.err != nil {
		return nil, mock.err
	}
	return []m.InvestSummary{
		{
			ID:      1,
			FundID:  1,
			AssetID: 1,
			Sum:     568210,
		},
	}, nil
}
func (mock FundRetrieverMock) RetreiveAFundInvestsById(id uint) ([]m.Invest, error) {
	fmt.Println("RetreiveAFundInvestsById Called")

	if mock.err != nil {
		return nil, mock.err
	}
	return []m.Invest{
		{
			ID:      id,
			FundID:  3,
			AssetID: 1,
			Price:   7800,
			Count:   5,
		},
	}, nil
}

type FundWriterMock struct {
	err error
}

func (mock FundWriterMock) SaveFund(name string) error {
	fmt.Println("SaveFund Called")

	if mock.err != nil {
		return mock.err
	}
	return nil
}

type ExchageRateGetterMock struct {
}

func (mock ExchageRateGetterMock) GetRealtimeExchageRate() float64 {
	fmt.Println("SaveFund Called")

	return 1334.3
}

/***************************** Market ***********************************/
type MaketRetrieverMock struct {
	err error
}

func (mock MaketRetrieverMock) RetrieveMarketStatus(date string) (*m.Market, error) {
	fmt.Println("RetrieveMarketStatus Called")

	if mock.err != nil {
		return nil, mock.err
	}
	return &m.Market{
		CreatedAt: datatypes.Date(time.Now()),
		Status:    3,
	}, nil
}

func (mock MaketRetrieverMock) RetrieveMarketIndicator(date string) (*m.DailyIndex, *m.CliIndex, error) {
	fmt.Println("RetrieveMarketIndicator Called")

	if mock.err != nil {
		return nil, nil, mock.err
	}
	return &m.DailyIndex{
			CreatedAt:      datatypes.Date(time.Now()),
			FearGreedIndex: 23,
			NasDaq:         17556.03,
		}, &m.CliIndex{
			CreatedAt: datatypes.Date(time.Now()),
			Index:     102,
		}, nil
}

type MarketSaverMock struct {
	err error
}

func (mock MarketSaverMock) SaveMarketStatus(status uint) error {
	fmt.Println("SaveMarketStatus Called")

	if mock.err != nil {
		return mock.err
	}
	return nil
}

/***************************** Invest ***********************************/
type InvestRetrieverMock struct {
	err error
}

func (mock InvestRetrieverMock) RetrieveInvestHist(fundId uint, assetId uint, start string, end string) ([]m.Invest, error) {
	fmt.Println("RetrieveInvestHist Called")

	if mock.err != nil {
		return nil, mock.err
	}
	return []m.Invest{
		{
			ID:      1,
			FundID:  fundId,
			AssetID: assetId,
			Price:   7800,
			Count:   5,
		},
	}, nil
}

type InvestSaverMock struct {
	err error
}

func (mock InvestSaverMock) SaveInvest(fundId uint, assetId uint, price float64, count int) error {
	fmt.Println("SaveInvest Called")

	if mock.err != nil {
		return mock.err
	}
	return nil
}
