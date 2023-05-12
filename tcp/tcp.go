package tcp

func new(listener chan *Message, myListener chan *Message, port int, host string) *TCP {
	return &TCP{
		Port:       port,
		Host:       host,
		ListenerCh: listener,
		quitch:     make(chan struct{}),
		myListener: myListener,
	}
}
