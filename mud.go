package avp

import (
	"log"
	"net"
	"sync"
)

// MUD represents a MUD.
type MUD struct {
	sync.RWMutex
	clients  []*Client
	listener net.Listener
}

// BroadcastAll broadcasts a message to all connected clients.
func (m *MUD) BroadcastAll(msg string, mustLoggedIn bool) {
	m.RLock()
	defer m.RUnlock()
	for _, client := range m.clients {
		if mustLoggedIn {
			if !client.State.IsLoggedIn() {
				continue
			}
		}
		client.Write <- msg
	}
}

// CallbackErr handles errors encountered by goroutines spawned by the MUD.
func (m *MUD) CallbackErr(err error) {
	log.Println(err)
}

// ListenAndServe concurrently accepts and serves new Clients.
func (m *MUD) ListenAndServe(address string) (err error) {
	if m.listener, err = net.Listen("tcp", address); err != nil {
		return
	}
	defer m.listener.Close()

	var conn net.Conn
	for {
		conn, err = m.listener.Accept()
		if err != nil {
			return
		}
		client := new(Client).ListenAndServe(conn, m)
		m.Lock()
		m.clients = append(m.clients, client)
		m.Unlock()
	}
}
