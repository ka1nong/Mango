package common

type Error string

func (e Error) Error() string {
	return string(e)
}

const (
	STOCKINFO_DATE                = "date"
	STOCKINFO_OPEN                = "open"
	STOCKINFO_HIGHT               = "hight"
	STOCKINFO_LOW                 = "low"
	STOCKINFO_CLOSE               = "close"
	STOCKINFO_CHANGE              = "hange"
	STOCKINFO_VOLUME              = "volume"
	STOCKINFO_MONEY               = "money"
	STOCKINFO_TRADED_MARKET_VALUE = "traded_market_value"
	STOCKINFO_MARKET_VALUE        = "market_value"
	STOCKINFO_TURNOVER            = "turnover"
	STOCKINFO_ADJUST_PRICE        = "adjust_price"
	STOCKINFO_REPORT_TYPE         = "report_type"
	STOCKINFO_REPORT_DATE         = "report_date"
	STOCKINFO_PE_TTM              = "PE_TTM"
	STOCKINFO_PS_TTM              = "PS_TTM"
	STOCKINFO_PC_TTM              = "PC_TTM"
	STOCKINFO_PB                  = "PB"
)

const (
	STOCK_CODE    = "code"
	STOCK_NAME    = "name"
	STOCK_ADDRESS = "address"
)
