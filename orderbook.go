package orderbook

import (
	"fmt"
	"slices"
	"sort"
)

type OrderType string

const (
	Buy  OrderType = "buy"
	Sell OrderType = "sell"
)

type Order struct {
	ID     string
	Type   OrderType
	Price  float64 // Supports fractional amounts
	Amount float64 // Supports fractional amounts
}

type Trade struct {
	BuyOrderID  string
	SellOrderID string
	Price       float64
	Amount      float64
}

type orderBook interface {
	// Add a buy or sell order to the order book.
	Add(order Order) error
	// Cancel an order by ID.
	Cancel(orderID string) error
	// Retrieve the top buy and sell orders from the order book.
	Get() (buys []Order, sells []Order)
	// Match the top buy and sell orders from the order book.
	Match() ([]Trade, error)
}

var _ orderBook = (*OrderBook)(nil)

type OrderBook struct {
	buyOrders  []Order
	sellOrders []Order
}

func New() *OrderBook {
	return &OrderBook{}
}

func (ob *OrderBook) Add(o Order) error {
	switch o.Type {
	case Buy:
		ob.buyOrders = append(ob.buyOrders, o)
		// Highest price first.
		sort.Slice(ob.buyOrders, func(i, j int) bool {
			return ob.buyOrders[i].Price > ob.buyOrders[j].Price
		})
	case Sell:
		ob.sellOrders = append(ob.sellOrders, o)
		// Lowest price first.
		sort.Slice(ob.sellOrders, func(i, j int) bool {
			return ob.sellOrders[i].Price < ob.sellOrders[j].Price
		})
	default:
		return &OrderError{
			ID:     o.ID,
			Reason: fmt.Sprintf("invalid order type: %s", o.Type),
		}
	}

	return nil
}

func (ob *OrderBook) Cancel(orderID string) error {
	for i, o := range ob.sellOrders {
		if o.ID == orderID {
			ob.sellOrders = slices.Delete(ob.sellOrders, i, i+1)
			return nil
		}
	}

	for i, o := range ob.buyOrders {
		if o.ID == orderID {
			ob.buyOrders = slices.Delete(ob.buyOrders, i, i+1)
			return nil
		}
	}

	return &OrderError{
		ID:     orderID,
		Reason: "not found",
	}
}

func (ob *OrderBook) Match() (trades []Trade, err error) {
	for len(ob.buyOrders) > 0 && len(ob.sellOrders) > 0 {
		buy := ob.buyOrders[0]
		sell := ob.sellOrders[0]

		if buy.Price < sell.Price {
			break
		}

		amount := min(buy.Amount, sell.Amount)
		trades = append(trades, Trade{
			BuyOrderID:  buy.ID,
			SellOrderID: sell.ID,
			Price:       sell.Price,
			Amount:      amount,
		})

		ob.buyOrders[0].Amount -= amount
		ob.sellOrders[0].Amount -= amount

		if ob.buyOrders[0].Amount == 0 {
			ob.buyOrders = ob.buyOrders[1:]
		}
		if ob.sellOrders[0].Amount == 0 {
			ob.sellOrders = ob.sellOrders[1:]
		}
	}

	return
}

func (ob *OrderBook) Get() (buys []Order, sells []Order) {
	return ob.buyOrders, ob.sellOrders
}

type OrderError struct {
	ID     string
	Reason string
}

func (o *OrderError) Error() string {
	return fmt.Sprintf("Order(id=%s) %s", o.ID, o.Reason)
}
