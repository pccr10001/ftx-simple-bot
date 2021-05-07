package ftx

import (
	"github.com/go-numb/go-ftx/auth"
	"github.com/go-numb/go-ftx/rest"
	"github.com/go-numb/go-ftx/rest/private/account"
	"github.com/go-numb/go-ftx/rest/private/orders"
	"github.com/go-numb/go-ftx/rest/public/futures"
	"github.com/go-numb/go-ftx/rest/public/markets"
	"github.com/pccr10001/jrb-bot/app"
	"log"
	"strconv"
)

var Client *rest.Client

var CoinMap map[string]string

const (
	SIDE_BUY  = "buy"
	SIDE_SELL = "sell"
)

func fillMap() {
	CoinMap = make(map[string]string)
	for _, m := range app.AppConfig.Market {
		CoinMap[m.Symbol] = m.Market
	}
}

//func GetBalance(symbol string) (*wallet.Balance, error) {
//	balances, err := Client.Balances(&wallet.RequestForBalances{})
//	if err != nil {
//		return nil, err
//	}
//	for _, b := range *balances {
//		if b.Coin == CoinMap[symbol]{
//
//		}
//	}
//}

func Init() {
	Client = rest.New(
		auth.New(
			app.AppConfig.Exchange.Ftx.APIKey,
			app.AppConfig.Exchange.Ftx.APISecret,
			auth.SubAccount{
				UUID:     1,
				Nickname: app.AppConfig.Exchange.Ftx.SubAccount,
			},
		))
	Client.Auth.UseSubAccountID(1)

	i, err := GetInformation()

	if err != nil {
		log.Fatalln(err)
	}

	fillMap()

	log.Printf("Account %s logged in\n", i.Username)
}

func GetInformation() (*account.ResponseForInformation, error) {
	i, err := Client.Information(&account.RequestForInformation{})
	if err != nil {
		return nil, err
	}
	return i, nil
}

func GetFuture(symbol string) (*futures.Future, error) {
	f, err := Client.Future(&futures.RequestForFuture{ProductCode: CoinMap[symbol]})
	if err != nil {
		return nil, err
	}
	var ff = futures.Future(*f)
	return &ff, nil
}

func GetPosition(symbol string) (*account.Position, error) {
	f, err := Client.Positions(&account.RequestForPositions{ShowAvgPrice: true})
	if err != nil {
		return nil, err
	}
	for _, p := range *f {
		if p.Future == CoinMap[symbol] {
			if p.Size == 0 {
				return nil, nil
			}
			return &p, nil
		}
	}
	return nil, nil
}

func GetMarket(symbol string) (float64, error) {
	m, err := Client.Markets(&markets.RequestForMarkets{ProductCode: CoinMap[symbol]})
	if err != nil {
		return 0, err
	}
	return (*m)[0].Price, nil
}

func PlaceOrder(symbol string, amount float64, side string) (*orders.ResponseForOrderStatus, error) {
	price, err := GetMarket(symbol)
	if err != nil {
		return nil, err
	}
	size := amount / price
	order, err := Client.PlaceOrder(&orders.RequestForPlaceOrder{
		Type:   "market",
		Market: CoinMap[symbol],
		Side:   side,
		Price:  0,
		Size:   size,
	})

	if err != nil {
		return nil, err
	}

	status, err := Client.OrderStatus(&orders.RequestForOrderStatus{
		ClientID: "",
		OrderID:  strconv.Itoa(order.ID),
	})

	if err != nil {
		return nil, err
	}

	return status, nil

}

func FillPosition(position account.Position) (*orders.ResponseForOrderStatus, error) {

	side := SIDE_BUY

	if position.Side == SIDE_BUY {
		side = SIDE_SELL
	}

	order, err := Client.PlaceOrder(&orders.RequestForPlaceOrder{
		Type:   "market",
		Market: position.Future,
		Side:   side,
		Price:  0,
		Size:   position.Size,
	})

	if err != nil {
		return nil, err
	}

	status, err := Client.OrderStatus(&orders.RequestForOrderStatus{
		ClientID: "",
		OrderID:  strconv.Itoa(order.ID),
	})

	if err != nil {
		return nil, err
	}

	return status, nil

}
