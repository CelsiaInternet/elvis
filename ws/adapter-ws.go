package ws

type AdapterWS struct {
	conn *Client
}

/**
* Subscribed
* @param channel string
**/
func (s *AdapterWS) Subscribed(channel string) {
	channel = clusterChannel(channel)
	s.conn.Subscribe(channel, func(msg Message) {
		if msg.tp == TpDirect {
			s.conn.SendMessage(msg.Id, msg)
		} else {
			s.conn.Publish(msg.Channel, msg)
		}
	})
}

/**
* UnSubscribed
* @param sub channel string
**/
func (s *AdapterWS) UnSubscribed(channel string) {
	channel = clusterChannel(channel)
	s.conn.Unsubscribe(channel)
}

/**
* Publish
* @param sub channel string
**/
func (s *AdapterWS) Publish(channel string, msg Message) {
	channel = clusterChannel(channel)
	s.conn.Publish(channel, msg)
}
