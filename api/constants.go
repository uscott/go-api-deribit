package api

// Various constants
const (
	Big           = OneMillion * OneMillion
	Bit           = 100 * Satoshi
	Bp            = Pct * Pct // one basis point
	BTC           = "BTC"
	Edp           = "edp"
	ETH           = "ETH"
	Huge          = Big * Big
	OneMillion    = OneThousand * OneThousand
	OneThousand   = 1000
	Pct           = 0.01
	Satoshi       = 1.0 / (100 * OneMillion)
	SecondsInDay  = 3600 * 24
	SecondsInYear = 365 * SecondsInDay
	Small         = 1.0 / Big
	Swap          = "BTC-PERPETUAL"
	Tiny          = 1.0 / Huge
)

// Direction direction, `buy` or `sell`
const (
	DrctnBuy  = "buy"
	DrctnSell = "sell"
)

// OrderState order state, `"open"`, `"filled"`, `"rejected"`, `"cancelled"`, `"untriggered"`
const (
	OrdrStateOpen       = "open"
	OrdrStateFilled     = "filled"
	OrdrStateRejected   = "rejected"
	OrdrStateCancelled  = "cancelled"
	OrdrStateUntriggerd = "untriggered"
)

// OrderType order type, `"limit"`, `"market"`, `"stop_limit"`, `"stop_market"`
const (
	OrdrTypeLmt     = "limit"
	OrdrTypeMkt     = "market"
	OrdrTypeStopLmt = "stop_limit"
	OrdrTypeStopMkt = "stop_market"
)

// TriggerType trigger type, `"index_price"`, `"mark_price"`, `"last_price"`
const (
	TriggerTypeIdxPrc = "index_price"
	TriggerTypeMrkPrc = "mark_price"
	TriggerTypeLstPrc = "last_price"
)

const (
	TmInFrcIOC = "immediate_or_cancel"
)
