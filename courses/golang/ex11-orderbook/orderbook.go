package orderbook

import "container/list"

type Orderbook struct {
	Bids *list.List // sorted descending
	Asks *list.List // sorted ascending
}

func New() *Orderbook {
	return &Orderbook{
		Bids: list.New(),
		Asks: list.New(),
	}
}

func (orderbook *Orderbook) Match(order *Order) ([]*Trade, *Order) {
	switch order.Side {
	case SideBid:
		return orderbook.handleBid(order)
	case SideAsk:
		return orderbook.handleAsk(order)
	}

	return nil, nil
}

func (orderbook *Orderbook) addBid(new *Order) {
	var bid *Order
	for curr := orderbook.Bids.Front(); curr != nil; curr = curr.Next() {
		bid = curr.Value.(*Order)
		// sorted descending
		if new.Price > bid.Price {
			orderbook.Bids.InsertBefore(new, curr)
			return
		}
	}
	orderbook.Bids.PushBack(new)
}

func (orderbook *Orderbook) addAsk(new *Order) {
	var ask *Order
	for curr := orderbook.Asks.Front(); curr != nil; curr = curr.Next() {
		ask = curr.Value.(*Order)
		// sorted ascending
		if new.Price < ask.Price {
			orderbook.Asks.InsertBefore(new, curr)
			return
		}
	}
	orderbook.Asks.PushBack(new)
}

func (orderbook *Orderbook) handleBid(order *Order) ([]*Trade, *Order) {
	var trades []*Trade

	curr := orderbook.Asks.Front()
	var next *list.Element

	for ; curr != nil; curr = next {
		next = curr.Next()
		ask := curr.Value.(*Order)

		if order.Kind == KindMarket || order.Price >= ask.Price {
			trade := &Trade{
				Price: ask.Price,
				Bid:   order,
				Ask:   ask,
			}

			// processing trade volumes
			if ask.Volume > order.Volume {
				trade.Volume = order.Volume
				ask.Volume -= order.Volume
				order.Volume = 0
			} else {
				trade.Volume = ask.Volume
				order.Volume -= ask.Volume
				ask.Volume = 0
				orderbook.Asks.Remove(curr)
			}

			trades = append(trades, trade)

			// is order satisfied?
			if order.Volume == 0 {
				break
			}
		} else {
			break
		}
	}

	// if order is not satisfied
	if order.Volume > 1 {
		// if it was Market order, cancel it
		if order.Kind == KindMarket {
			return trades, order
		}
		// otherwise, make it resting
		orderbook.addBid(order)
	}

	return trades, nil
}

func (orderbook *Orderbook) handleAsk(order *Order) ([]*Trade, *Order) {
	var trades []*Trade

	curr := orderbook.Bids.Front()
	var next *list.Element

	for ; curr != nil; curr = next {
		next = curr.Next()
		bid := curr.Value.(*Order)

		if order.Kind == KindMarket || bid.Price >= order.Price {
			trade := &Trade{
				Price: bid.Price,
				Bid:   bid,
				Ask:   order,
			}

			// processing trade volumes
			if bid.Volume > order.Volume {
				trade.Volume = order.Volume
				bid.Volume -= order.Volume
				order.Volume = 0
			} else {
				trade.Volume = bid.Volume
				order.Volume -= bid.Volume
				bid.Volume = 0
				orderbook.Bids.Remove(curr)
			}

			trades = append(trades, trade)

			// is order satisfied?
			if order.Volume == 0 {
				break
			}
		} else {
			break
		}
	}

	// if order is not satisfied
	if order.Volume > 1 {
		// if it was Market order, cancel it
		if order.Kind == KindMarket {
			return trades, order
		}
		// otherwise, make it resting
		orderbook.addAsk(order)
	}

	return trades, nil
}
