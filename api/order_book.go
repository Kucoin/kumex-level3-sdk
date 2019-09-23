package api

func (s *Server) GetOrderBook(message *TokenMessage, reply *string) error {
	if err := s.checkToken(message.Token); err != "" {
		*reply = err
		return nil
	}

	data, err := s.level3Builder.SnapshotBytes()
	if err != nil {
		*reply = s.failure(TickerErrorCode, err.Error())
		return nil
	}

	*reply = s.success(string(data))
	return nil
}

type GetPartOrderBookMessage struct {
	Number int `json:"number"`
	TokenMessage
}

func (s *Server) GetPartOrderBook(message *GetPartOrderBookMessage, reply *string) error {
	if err := s.checkToken(message.Token); err != "" {
		*reply = err
		return nil
	}

	data, err := s.level3Builder.GetPartOrderBook(message.Number)
	if err != nil {
		*reply = s.failure(TickerErrorCode, err.Error())
		return nil
	}

	*reply = s.success(string(data))
	return nil
}
