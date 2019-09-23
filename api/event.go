package api

type AddEventOrderIdsMessage struct {
	Data map[string][]string `json:"data"`
	TokenMessage
}

type AddEventClientOidsMessage struct {
	Data map[string][]string `json:"data"`
	TokenMessage
}

func (s *Server) AddEventOrderIdsToChannels(message *AddEventOrderIdsMessage, reply *string) error {
	if err := s.checkToken(message.Token); err != "" {
		*reply = err
		return nil
	}

	if len(message.Data) == 0 {
		*reply = s.failure(ServerErrorCode, "empty event data")
		return nil
	}

	s.eventWatcher.AddEventOrderIdsToChannels(message.Data)

	*reply = s.success("")
	return nil
}

// You must subscribe in advance according to the ClientOids subscription,
// or you will miss the receive message because of without the mapping relationship between message and orderId~
func (s *Server) AddEventClientOidsToChannels(message *AddEventClientOidsMessage, reply *string) error {
	if err := s.checkToken(message.Token); err != "" {
		*reply = err
		return nil
	}

	if len(message.Data) == 0 {
		*reply = s.failure(ServerErrorCode, "empty event data")
		return nil
	}

	s.eventWatcher.AddEventClientOidsToChannels(message.Data)

	*reply = s.success("")
	return nil
}
