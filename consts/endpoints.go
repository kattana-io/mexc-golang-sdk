package consts

// HTTP
const (
	// Market
	EndpointExchangeInfo           = "/api/v3/exchangeInfo"
	EndpointOrder                  = "/api/v3/order"
	EndpointOrderBook              = "/api/v3/depth"
	EndpointPing                   = "/api/v3/ping"
	EndpointTime                   = "/api/v3/time"
	EndpointTradeFee               = "/api/v3/tradeFee"
	EndpointInternalTransfer       = "/api/v3/capital/transfer/internal"
	EndpointWithdraw               = "/api/v3/capital/withdraw"
	EndpointWithdrawHistory        = "/api/v3/capital/withdraw/history"
	EndpointGetCurrencyInformation = "/api/v3/capital/config/getall"

	// Stream
	EndpointStream = "/api/v3/userDataStream"
)
