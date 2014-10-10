// Activation example used by the activation unit tests.
package main

import (
	"fmt"
	"net"
	"os"

	"github.com/coreos/locksmith/third_party/github.com/coreos/go-systemd/activation"
)

func fixListenPid() {
	if os.Getenv("FIX_LISTEN_PID") != "" {
		// HACK: real systemd would set LISTEN_PID before exec'ing but
		// this is too difficult in golang for the purpose of a test.
		// Do not do this in real code.
		os.Setenv("LISTEN_PID", fmt.Sprintf("%d", os.Getpid()))
	}
}

func main() {
	fixListenPid()

	pc, _ := activation.PacketConns(false)

	if len(pc) == 0 {
		panic("No packetConns")
	}

	if os.Getenv("LISTEN_PID") == "" || os.Getenv("LISTEN_FDS") == "" {
		panic("Should not unset envs")
	}

	pc, err := activation.PacketConns(true)
	if err != nil {
		panic(err)
	}

	if os.Getenv("LISTEN_PID") != "" || os.Getenv("LISTEN_FDS") != "" {
		panic("Can not unset envs")
	}

	udp1, ok := pc[0].(*net.UDPConn)
	if !ok {
		panic("packetConn 1 not UDP")
	}
	udp2, ok := pc[1].(*net.UDPConn)
	if !ok {
		panic("packetConn 2 not UDP")
	}

	_, addr1, err := udp1.ReadFromUDP(nil)
	if err != nil {
		panic(err)
	}
	_, addr2, err := udp2.ReadFromUDP(nil)
	if err != nil {
		panic(err)
	}

	// Write out the expected strings to the two pipes
	_, err = udp1.WriteToUDP([]byte("Hello world"), addr1)
	if err != nil {
		panic(err)
	}
	_, err = udp2.WriteToUDP([]byte("Goodbye world"), addr2)
	if err != nil {
		panic(err)
	}

	return
}
