// Package mockscsynth provides an OSC server that can mock
// scsynth, the SuperCollider server component.
package mockscsynth

import (
	"net"
	"testing"

	"github.com/scgolang/osc"
)

// New creates a new mock scsynth server.
func New(t *testing.T, listenAddr string) osc.Conn {
	udpAddr, err := net.ResolveUDPAddr("udp", listenAddr)
	if err != nil {
		t.Fatal(err)
	}
	conn, err := osc.ListenUDP("udp", udpAddr)
	if err != nil {
		t.Fatal(err)
	}
	synthdefDoneMsg := osc.Message{
		Address: "/done", // TODO: reuse constants from sc
		Arguments: osc.Arguments{
			osc.String("/d_recv"),
		},
	}
	go func() {
		if err := conn.Serve(1, osc.Dispatcher{
			"/d_recv": osc.Method(func(m osc.Message) error {
				return conn.SendTo(m.Sender, synthdefDoneMsg)
			}),
			"/g_new": osc.Method(func(m osc.Message) error {
				return nil
			}),
		}); err != nil {
			t.Fatal(err)
		}
	}()
	return conn
}
