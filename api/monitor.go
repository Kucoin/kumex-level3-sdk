package api

import (
	"encoding/json"
	"time"
)

func (s *Server) GetChanLen(message *TokenMessage, reply *string) error {
	if err := s.checkToken(message.Token); err != "" {
		*reply = err
		return nil
	}

	ret := map[string]int{
		"level3Builder.Messages": len(s.level3Builder.Messages),
		"eventWatcher.Messages":  len(s.eventWatcher.Messages),
	}

	if data, err := json.Marshal(ret); err == nil {
		*reply = s.success(string(data))
		return nil
	}

	*reply = s.failure(ServerErrorCode, "json failed")
	return nil
}

func (s *Server) Time(message *TokenMessage, reply *string) error {
	if err := s.checkToken(message.Token); err != "" {
		*reply = err
		return nil
	}

	ret := map[string]int64{
		"time": time.Now().Unix(),
	}

	if data, err := json.Marshal(ret); err == nil {
		*reply = s.success(string(data))
		return nil
	}

	*reply = s.failure(ServerErrorCode, "json failed")
	return nil
}
