package handler

import (
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
)

// router.Get("/", h.Market)
// router.Get("/indicator", h.MarketIndicator)
// router.Post("/", h.ChangeMarketStatus)

func TestMarketHandler(t *testing.T) {

	app := fiber.New()

	readerMock := MaketRetrieverMock{}
	writerMock := MarketSaverMock{}

	f := MarketHandler{
		r: readerMock,
		w: writerMock,
	}
	f.InitRoute(app)
	go func() {
		app.Listen(":3000")
	}()

	t.Run("시장 단계 조회", func(t *testing.T) {
		t.Run("성공 테스트", func(t *testing.T) {
			param := MarketStatusParam{
				Date: "202-08-29",
			}
			err := sendReqeust(app, "/market", "GET", param)
			assert.NoError(t, err)
		})

	})

	t.Run("시장 단계 저장", func(t *testing.T) {
		t.Run("성공 테스트", func(t *testing.T) {
			param := SaveMarketStatusParam{
				Status: 1,
			}
			err := sendReqeust(app, "/market", "POST", param)
			assert.NoError(t, err)
		})

	})

	t.Run("시장 지표 조회", func(t *testing.T) {
		t.Run("성공 테스트", func(t *testing.T) {
			param := MarketStatusParam{
				Date: "202-08-29",
			}
			err := sendReqeust(app, "/market/indicator", "GET", param)
			assert.NoError(t, err)
		})

	})

	app.Shutdown()
}
