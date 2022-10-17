package ndau

import (
	"net/http"
	"time"
)

type AccountListReq struct {
	Limit int    `json:"limit"`
	After string `json:"after"`
}

// AccountListResp
type AccountListResp struct {
	NumAccounts int      `json:"NumAccounts"`
	FirstIndex  int      `json:"FirstIndex"`
	After       string   `json:"After"`
	NextAfter   string   `json:NextAfter"`
	Accounts    []string `json:Accounts"`
}

// Account ...
type Account struct {
	CurrencySeatDate time.Time `json:"CurrencySeatDate"`
	Id               string    `json:"id"`
	Balance          int       `json:"balance"`
}

// AccountResp
type AccountResp map[string]Account

// CurrentPriceResp
type CurrentPriceResp struct {
	MarketPrice   int `json:"marketPrice"`
	TargetPrice   int `json:"targetPrice"`
	FloorPrice    int `json:"floorPrice"`
	TotalReleased int `json:"totalReleased"`
	TotalIssued   int `json:"totalIssued"`
	TotalNdau     int `json:"totalNdau"`
	TotalBurned   int `json:"totalBurned"`
	SIB           int `json:"sib"`
}

//go:generate mockgen -destination=./mocks/mock_http_client.go -package=mocks github.com/ndau/go-ndau/ndau HttpClient
type HttpClient interface {
	Do(req *http.Request) (*http.Response, error)
}
