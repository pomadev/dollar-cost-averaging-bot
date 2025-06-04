package client

const (
	// Order types
	OrderTypeMarket = "market"
	OrderSideBuy    = "buy"
	
	// Precision constants for amount calculation
	BitbankPrecision  = 10000
	BitflyerBTCPrecision = 1000
	BitflyerETHPrecision = 100
	
	// HTTP Status
	HTTPStatusOK = 200
	
	// Success codes
	SuccessCode = 1
	
	// Format settings
	AmountFormatPrecision = 4
	AmountFormatVerb = 'f'
)