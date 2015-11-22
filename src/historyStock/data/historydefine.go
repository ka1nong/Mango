package historyStock

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
	STOCKINFO_CHANGE              = "hange" //什么鬼，语法不能插入change字段，会报错。我操。只能把change改成hange
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
	USER     = "root"
	PASSWORD = "199212"
	DATABASE = "historyStock"
)
