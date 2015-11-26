package onlineStock

type Error string

func (e Error) Error() string {
	return string(e)
}

const (
	STOCKINFO_OPEN            = "open"
	STOCKINFO_YESTERDAYCLOASE = "yesterdayClose"
	STOCKINFO_CURRENTPRICE    = "currentPrice"
	STOCKINFO_TODAYHIGHT      = "todayHight"
	STOCKINFO_TODAYLOW        = "toayLow"
	STOCKINFO_BUYONE          = "buyOne"
	STOCKINFO_SELLONE         = "sellOne"
	STOCKINFO_DEALCOUNT       = "dealCount"
	STOCKINFO_DEALPRICE       = "dealPrice"
	STOCKINFO_BUYONECOUNT     = "buyOneCount"
	STOCKINFO_BUYONEPRICE     = "buyOnePrice"
	STOCKINFO_BUYTWOCOUNT     = "buyTwoCount"
	STOCKINFO_BUYTWOPRICE     = "buyTwoPrice"
	STOCKINFO_BUYTHREECOUNT   = "buyThreeCount"
	STOCKINFO_BUYTHREEPRICE   = "buyThreePrice"
	STOCKINFO_BUYFOURCOUNT    = "buyFourCount"
	STOCKINFO_BUYFOURPRICE    = "buyFourPrice"
	STOCKINFO_BUYFIVECOUNT    = "buyFiveCount"
	STOCKINFO_BUYFIVEPRICE    = "buyFivePrice"
	STOCKINFO_SELLONECOUNT    = "sellOneCount"
	STOCKINFO_SELLONEPRICE    = "sellOnePrice"
	STOCKINFO_SELLTWOCOUNT    = "sellTwoCount"
	STOCKINFO_SELLTWOPRICE    = "sellTwoPrice"
	STOCKINFO_SELLTHREECOUNT  = "sellThreeCount"
	STOCKINFO_SELLTHREEPRICE  = "sellThreePrice"
	STOCKINFO_SELLFOURCOUNT   = "sellFourCount"
	STOCKINFO_SELLFOURPRICE   = "sellFourPrice"
	STOCKINFO_SELLFIVECOUNT   = "sellFiveCount"
	STOCKINFO_SELLFIVEPRICE   = "sellFivePrice"
	STOCKINFO_CURRENTDATA     = "currentData"
	STOCKINFO_CURRETNTIME     = "currentTime"
)

const (
	USER       = "root"
	PASSWORD   = "199212"
	DATABASE   = "onlineStock"
	MAIN_TABLE = "stockList"
)
