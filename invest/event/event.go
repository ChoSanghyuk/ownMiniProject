package event

import (
	"cmp"
	"fmt"
	m "invest/model"
	"log"
	"slices"
	"strings"
	"time"
)

type Event struct {
	stg Storage
	rt  RtPoller
	dp  DailyPoller
}

func NewEvent(stg Storage, rtPoller RtPoller, dailyPoller DailyPoller) *Event {
	return &Event{
		stg: stg,
		rt:  rtPoller,
		dp:  dailyPoller,
	}
}

/*
작업 1. 자산의 현재가와 자산의 매도/매수 기준 비교하여 알림 전송
  - 보유 자산 list
  - 자산 정보
  - 현재가

작업 2. 자금별/종목별 현재 총액 갱신 + 최저가/최고가 갱신
  - investSummary list
  - 현재가
  - 환율
  - 자산 정보

직업 3. 현재 시장 단계에 맞는 변동 자산을 가지고 있는지 확인하여 알림 전송. 대상 시, 우선처분 대상 및 보유 자산 현환 전송
  - 시장 단계
  - 갱신된 investSummary list
*/

func (e Event) AssetEvent(c chan<- string) {

	// 등록 자산 목록 조회
	assetList, err := e.stg.RetrieveAssetList()
	if err != nil {
		c <- fmt.Sprintf("[AssetEvent] RetrieveAssetList 시, 에러 발생. %s", err)
		return
	}
	priceMap := make(map[uint]float64) // assetId => price

	// 등록 자산 매수/매도 기준 충족 시, 채널로 메시지 전달
	for _, a := range assetList {
		msg, err := e.buySellMsg(a.ID, priceMap)
		if err != nil {
			c <- fmt.Sprintf("[AssetEvent] buySellMsg시, 에러 발생. %s", err)
			return
		}
		if msg != "" {
			c <- msg
		}
	}

	// 자금별 종목 투자 내역 조회
	ivsmLi, err := e.stg.RetreiveFundsSummaryOrderByFundId()
	if err != nil {
		c <- fmt.Sprintf("[AssetEvent] RetreiveFundsSummaryOrderByFundId 시, 에러 발생. %s", err)
		return
	}
	if len(ivsmLi) == 0 {
		return
	}

	// 자금별/종목별 현재 총액 갱신
	err = e.updateFundSummarys(ivsmLi, priceMap)
	if err != nil {
		c <- fmt.Sprintf("[AssetEvent] updateFundSummary 시, 에러 발생. %s", err)
		return
	}

}

func (e Event) PortfolioEvent(c chan<- string) {

	// 자금별 종목 투자 내역 조회
	ivsmLi, err := e.stg.RetreiveFundsSummaryOrderByFundId()
	if err != nil {
		c <- fmt.Sprintf("[AssetEvent] RetreiveFundsSummaryOrderByFundId 시, 에러 발생. %s", err)
		return
	}
	if len(ivsmLi) == 0 {
		return
	}

	// 현재 시장 단계 이하로 변동 자산을 가지고 있는지 확인. (알림 전송)
	msg, err := e.portfolioMsg(ivsmLi)
	if err != nil {
		c <- fmt.Sprintf("[AssetEvent] portfolioMsg시, 에러 발생. %s", err)
	}
	if msg != "" {
		c <- msg
	}
}

func (e Event) RealEstateEvent(c chan<- string) {

	rtn, err := e.rt.RealEstateStatus()
	if err != nil {
		c <- fmt.Sprintf("크롤링 시 오류 발생. %s", err.Error())
		return
	}

	if rtn != "예정지구 지정" {
		c <- fmt.Sprintf("연신내 재개발 변동 사항 존재. 예정지구 지정 => %s", rtn)
	} else {
		log.Printf("연신내 변동 사항 없음. 현재 단계: %s", rtn)
	}
}

func (e Event) IndexEvent(c chan<- string) {

	// 1. 공포 탐욕 지수
	fgi, err := e.dp.FearGreedIndex()
	if err != nil {
		c <- fmt.Sprintf("공포 탐욕 지수 조회 시 오류 발생. %s", err.Error())
		return
	}
	// 2. Nasdaq 지수 조회
	nasdaq, err := e.dp.Nasdaq()
	if err != nil {
		c <- fmt.Sprintf("Nasdaq Index 조회 시 오류 발생. %s", err.Error())
		return
	}

	// 오늘분 저장
	err = e.stg.SaveDailyMarketIndicator(fgi, nasdaq)
	if err != nil {
		c <- fmt.Sprintf("Nasdaq Index 저장 시 오류 발생. %s", err.Error())
	}

	// 어제꺼 조회
	yesterday := time.Now().AddDate(0, 0, -1).Format("2006-01-02")
	di, _, err := e.stg.RetrieveMarketIndicator(yesterday)
	if err != nil {
		c <- fmt.Sprintf("어제자 Nasdaq Index 저장 시 오류 발생. %s", err.Error())
		c <- fmt.Sprintf("금일 공포 탐욕 지수 : %d\n금일 Nasdaq : %f", fgi, nasdaq)
	} else {
		c <- fmt.Sprintf("금일 공포 탐욕 지수 : %d (전일 : %d)\n금일 Nasdaq : %f (전일 : %f)", fgi, di.FearGreedIndex, nasdaq, di.NasDaq)
	}

}

/**********************************************************************************************************************
*********************************************Inner Function************************************************************
**********************************************************************************************************************/

func (e Event) buySellMsg(assetId uint, pm map[uint]float64) (msg string, err error) {

	// 자산 정보 조회
	a, err := e.stg.RetrieveAsset(assetId)
	if err != nil {
		return "", fmt.Errorf("[AssetEvent] RetrieveAsset 시, 에러 발생. %w", err)
	}

	// 자산별 현재 가격 조회
	category, err := m.ToCategory(a.Category)
	if err != nil {
		return "", fmt.Errorf("[AssetEvent] ToCategory시, 에러 발생. %w", err)
	}
	cp, err := e.rt.CurrentPrice(category, a.Code)
	if err != nil {
		return "", fmt.Errorf("[AssetEvent] CurrentPrice 시, 에러 발생. %w", err)
	}

	log.Printf("%s 현재 가격 %.3f", a.Name, cp)

	pm[assetId] = cp

	// 자산 매도/매수 기준 비교 및 알림 여부 판단. (알림 전송)
	if a.BuyPrice >= cp && !hasMsgCache(a.ID, false, a.BuyPrice) {
		msg = fmt.Sprintf("BUY %s. ID : %d. LOWER BOUND : %f. CURRENT PRICE :%f", a.Name, a.ID, a.BuyPrice, cp)
		setMsgCache(a.ID, false, a.BuyPrice)
	} else if a.SellPrice != 0 && a.SellPrice <= cp && !hasMsgCache(a.ID, true, a.SellPrice) {
		msg = fmt.Sprintf("SELL %s. ID : %d. UPPER BOUND : %f. CURRENT PRICE :%f", a.Name, a.ID, a.SellPrice, cp)
		setMsgCache(a.ID, true, a.SellPrice)
	}

	return
}

func (e Event) updateFundSummarys(list []m.InvestSummary, pm map[uint]float64) (err error) {
	for i := range len(list) {
		is := &list[i]
		is.Sum = pm[is.AssetID] * float64(is.Count)

		err = e.stg.UpdateInvestSummarySum(is.FundID, is.AssetID, is.Sum)
		if err != nil {
			return
		}
	}
	return nil
}

func (e Event) portfolioMsg(ivsmLi []m.InvestSummary) (msg string, err error) {
	// 현재 시장 단계 조회
	market, err := e.stg.RetrieveMarketStatus("")
	if err != nil {
		msg = fmt.Sprintf("[AssetEvent] RetrieveMarketStatus 시, 에러 발생. %s", err)
		return
	}
	marketLevel := m.MarketLevel(market.Status)

	// 환율까지 계산하여 원화로 변환
	ex := e.dp.ExchageRate()
	if ex == 0 {
		msg = "[AssetEvent] ExchageRate 시 환율 값 0 반환"
		return
	}

	keySet := make(map[uint]bool)
	stable := make(map[uint]float64)
	volatile := make(map[uint]float64)

	for i := range len(ivsmLi) {

		ivsm := &ivsmLi[i]

		keySet[ivsm.FundID] = true

		// 원화 가치로 환산
		var v float64
		if ivsm.Asset.Currency == m.USD.String() {
			v = ivsm.Sum * ex
		} else {
			v = ivsm.Sum
		}

		category, err := m.ToCategory(ivsm.Asset.Category)
		if err != nil {
			msg = fmt.Sprintf("[AssetEvent] investSummary loop내 ToCategory시, 에러 발생. %s", err)
			return msg, err
		}

		// 자금 종류별 안전 자산 가치, 변동 자산 가치 총합 계산
		if category.IsStable() {
			stable[ivsm.FundID] = stable[ivsm.FundID] + v
		} else {
			volatile[ivsm.FundID] = volatile[ivsm.FundID] + v
		}
	}

	var sb strings.Builder
	type priority struct {
		asset *m.Asset
		ap    float64
		cp    float64
		hp    float64
		score float64
	}
	for k := range keySet {
		if volatile[k]+stable[k] == 0 {
			continue
		}

		r := volatile[k] / (volatile[k] + stable[k])
		if r > marketLevel.MaxVolatileAssetRate() || r < marketLevel.MinVolatileAssetRate() { // 매도해야 함
			sb.WriteString(strings.Repeat("=", 20))
			sb.WriteString("\n")

			os := make([]priority, 0) // ordered slice
			for _, ivsm := range ivsmLi {
				if ivsm.FundID == k {

					a := &ivsm.Asset
					category, err := m.ToCategory(a.Category)
					if err != nil {
						return "", fmt.Errorf("[AssetEvent] ToCategory시, 에러 발생. %w", err)
					}

					cp, ap, hp, _, err := e.rt.AssetPriceInfo(category, a.Code)
					if err != nil {
						return "", fmt.Errorf("[AssetEvent] AssetPriceInfo, 에러 발생. %w", err)
					}

					os = append(os, priority{
						asset: a,
						ap:    ap,
						cp:    cp,
						hp:    hp,
						score: 0.6*((cp-ap)/cp) + 0.4*((cp-hp)/cp),
					})
				}
			}

			if r > marketLevel.MaxVolatileAssetRate() {
				sb.WriteString(fmt.Sprintf("자금 %d 변동 자산 비중 초과. 변동 자산 비율 : %.2f. 현재 시장 단계 : %s(%.1f)\n\n", k, r, marketLevel.String(), marketLevel.MaxVolatileAssetRate()))
				slices.SortFunc(os, func(a, b priority) int {
					// todo. 안전 자산일 때 매도 대상 후순위 로직 필요. + category 변환... 너무 번거로움. 방법 필요함
					return cmp.Compare(b.score, a.score) // 큰 게 앞으로
				})
			} else {
				sb.WriteString(fmt.Sprintf("자금 %d 변동 자산 비중 부족. 변동 자산 비율 : %.2f. 현재 시장 단계 : %s(%.1f)\n\n", k, r, marketLevel.String(), marketLevel.MinVolatileAssetRate()))
				slices.SortFunc(os, func(a, b priority) int {
					return cmp.Compare(a.score, b.score)
				})
			}

			for _, p := range os {
				sb.WriteString(fmt.Sprintf("AssetId : %d, AssetName : %s, CurrentPrice : %f, WeighedAveragePrice : %f, HighestPrice : %f\n", p.asset.ID, p.asset.Name, p.cp, p.ap, p.hp))
			}
		}

	}

	msg = sb.String()
	return
}

/*
[판단]
현재가가 고점 및 이평가보다 낮을수록 저평가(조정) => 매수
현재가가 고점 및 이평가보다 높을수록 고평가 => 매도

[수식]
cp - 현재가
ap - 평균가
hp - 최고가

매도매수지수 = 0.6*((cp-ap)/cp) + 0.4*((cp-hp))/cp)
매도매수지수 클수록 매도 우선 순위
매도매수지수 낮을수록 매수 우선순위
*/
