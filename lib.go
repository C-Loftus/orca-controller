//go:generate go run ./generator
package main

import (
	"errors"

	"github.com/godbus/dbus/v5"
)

func NewOrcaClient() (*OrcaClient, error) {
	conn, err := dbus.ConnectSessionBus()
	if err != nil {
		return nil, err
	}

	return &OrcaClient{conn: conn}, nil
}

func (c *OrcaClient) Close() {
	c.conn.Close()
}

// Present a message to the user user Orca's speech and/or
// braille output
func (c *OrcaClient) PresentMessage(message string) error {
	obj := c.conn.Object(OrcaServiceName, OrcaObjectPath)
	var result bool
	err := obj.Call("org.gnome.Orca.Service.PresentMessage", 0, message).Store(&result)
	if err != nil {
		return err
	}
	if result {
		// this could be due to voice being disabled in Orca
		return errors.New("Message was not presented")
	}
	return nil
}
