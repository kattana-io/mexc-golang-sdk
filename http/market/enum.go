package mexchttpmarket

type Side string

const (
	SideBuy  Side = "BUY"
	SideSell Side = "SELL"
)

type Type string

const (
	TypeLimit             Type = "LIMIT"
	TypeMarket            Type = "MARKET"
	TypeLimitMarket       Type = "LIMIT_MARKET"
	TypeImmediateOrCancel Type = "IMMEDIATE_OR_CANCEL"
	TypeFillOrKill        Type = "FILL_OR_KILL"
)

type Status string

//nolint:misspell
const (
	StatusNew                Status = "NEW"
	StatusFilled             Status = "FILLED"
	StatusPartiallyFilled    Status = "PARTIALLY_FILLED"
	StatusCancelled          Status = "CANCELLED"
	StatusPartiallyCancelled Status = "PARTIALLY_CANCELLED"
)
