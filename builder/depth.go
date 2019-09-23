package builder

import "github.com/Kucoin/kumex-level3-sdk/level3stream"

//[5]interface{}{"orderTime", "orderId", "price", "size", "ts"}
type DepthResponse struct {
	Sequence uint64           `json:"sequence"`
	Asks     [][5]interface{} `json:"asks"`
	Bids     [][5]interface{} `json:"bids"`
}

//[3]string{"orderId", "price", "size"}
type FullOrderBook struct {
	Sequence uint64      `json:"sequence"`
	Asks     [][3]string `json:"asks"`
	Bids     [][3]string `json:"bids"`
}

func (b *Builder) GetAtomicFullOrderBook() (*DepthResponse, error) {
	resp, err := b.apiService.AtomicFullOrderBook(b.symbol)
	if err != nil {
		return nil, err
	}

	var fullOrderBook DepthResponse
	if err := resp.ReadJson(&fullOrderBook); err != nil {
		return nil, err
	}

	return &fullOrderBook, nil
}

func (b *Builder) GetPullingMessages(start, end uint64) ([]*level3stream.StreamDataModel, error) {
	resp, err := b.apiService.AtomicFullOrderBook(b.symbol)
	if err != nil {
		return nil, err
	}

	var pullingMessages []*level3stream.StreamDataModel
	if err := resp.ReadJson(&pullingMessages); err != nil {
		return nil, err
	}

	return pullingMessages, nil
}
