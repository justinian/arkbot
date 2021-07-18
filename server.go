package main

import (
	"errors"
	"fmt"
	"log"
	"net"
	"syscall"
	"time"
)

// This is linux specific, unfortunately, but I don't see a great way of
// checking this in Go without string parsing
var connRefused = syscall.Errno(111)

type Server struct {
	Name       string
	Map        string
	Folder     string
	Game       string
	Players    int
	MaxPlayers int
}

func (s Server) String() string {
	return fmt.Sprintf("%s : %d/%d players", s.Name, s.Players, s.MaxPlayers)
}

var query = []byte("\xff\xff\xff\xffTSource Engine Query\x00")

func popString(data []byte, skip int) ([]byte, string) {
	var i int
	for i = skip; data[i] != 0; i++ {
	}
	return data[i+1:], string(data[skip:i])
}

func checkServer(address string, verbose bool) (*Server, error) {
	if verbose {
		log.Printf("Checking %s", address)
	}

	addr, err := net.ResolveUDPAddr("udp", address)
	if err != nil {
		return nil, fmt.Errorf("Resolving UDP address: %w", err)
	}

	if verbose {
		log.Printf("  resolved to %s", addr)
	}

	conn, err := net.DialUDP("udp", nil, addr)
	if err != nil {
		return nil, fmt.Errorf("Creating UDP socket: %w", err)
	}

	_, err = conn.Write(query)
	if err != nil {
		return nil, fmt.Errorf("Sending UDP request: %w", err)
	}

	err = conn.SetReadDeadline(time.Now().Add(10 * time.Second))
	if err != nil {
		return nil, fmt.Errorf("Setting UDP deadline: %w", err)
	}

	buffer := make([]byte, 1024)
	n, _, err := conn.ReadFrom(buffer)
	if err != nil {
		if errors.Is(err, connRefused) {
			if verbose {
				log.Print("  connection refused")
			}
			return nil, nil
		}
		if neterr, ok := err.(net.Error); ok {
			if neterr.Timeout() || neterr.Temporary() {
				if verbose {
					log.Print("  timeout")
				}
				return nil, nil
			}
		}
		return nil, fmt.Errorf("Reading UDP request: %w", err)
	}

	data := buffer[:n]

	var s Server
	data, s.Name = popString(data, 6)
	data, s.Map = popString(data, 0)
	data, s.Folder = popString(data, 0)
	data, s.Game = popString(data, 0)
	s.Players = int(data[2])
	s.MaxPlayers = int(data[3])

	return &s, nil
}
