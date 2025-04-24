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

type WithdrawStatus int32

const (
	WithdrawStatusApply         WithdrawStatus = 1
	WithdrawStatusAuditing      WithdrawStatus = 2
	WithdrawStatusWait          WithdrawStatus = 3
	WithdrawStatusProcessing    WithdrawStatus = 4
	WithdrawStatusWaitPackaging WithdrawStatus = 5
	WithdrawStatusWaitConfirm   WithdrawStatus = 6
	WithdrawStatusSuccess       WithdrawStatus = 7
	WithdrawStatusFailed        WithdrawStatus = 8
	WithdrawStatusCancel        WithdrawStatus = 9
	WithdrawStatusManual        WithdrawStatus = 10
)

type TransferType int32

const (
	TransferTypeOutside TransferType = 0 // вывод на внешний адрес
	TransferTypeInside  TransferType = 1 // внутренний перевод (возможно, между аккаунтами)
)
