package main

import (
	"github.com/godbus/dbus/v5"
)

type Notifier struct {
	conn *dbus.Conn
}

func NewNotifier() *Notifier {
	conn, err := dbus.SessionBus()
	if err != nil {
		// Return a notifier that will fail gracefully
		return &Notifier{conn: nil}
	}
	return &Notifier{conn: conn}
}

func (n *Notifier) Send(title, message, iconPath string) error {
	if n.conn == nil {
		// Try to reconnect
		conn, err := dbus.SessionBus()
		if err != nil {
			return err
		}
		n.conn = conn
	}

	obj := n.conn.Object("org.freedesktop.Notifications", "/org/freedesktop/Notifications")
	call := obj.Call("org.freedesktop.Notifications.Notify", 0,
		"Emotional Support",       // app_name
		uint32(0),                 // replaces_id
		iconPath,                  // app_icon
		title,                     // summary
		message,                   // body
		[]string{},                // actions
		map[string]dbus.Variant{}, // hints
		int32(5000))               // expire_timeout in ms

	return call.Err
}
