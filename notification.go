package main

import (
	"github.com/godbus/dbus/v5"
)

type Notifier struct {
	conn           *dbus.Conn
	timeoutSeconds int32 // Notification timeout in seconds (0 = server default, -1 = never expire)
}

func NewNotifier() *Notifier {
	conn, err := dbus.SessionBus()
	if err != nil {
		// Return a notifier that will fail gracefully
		return &Notifier{conn: nil, timeoutSeconds: 5}
	}
	return &Notifier{conn: conn, timeoutSeconds: 15} // Default: 5 seconds
}

// SetTimeout sets the notification timeout in seconds
// 0 = use server default, -1 = never expire
func (n *Notifier) SetTimeout(seconds int32) {
	n.timeoutSeconds = seconds
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
	// expire_timeout: milliseconds (0 = server default, -1 = never expire)
	expireTimeout := n.timeoutSeconds * 1000
	call := obj.Call("org.freedesktop.Notifications.Notify", 0,
		"Emotional Support",       // app_name
		uint32(0),                 // replaces_id
		iconPath,                  // app_icon
		title,                     // summary
		message,                   // body
		[]string{},                // actions
		map[string]dbus.Variant{}, // hints
		expireTimeout)             // expire_timeout in ms

	return call.Err
}
