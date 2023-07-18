//go:build linux

package utils

import (
	"os/exec"

	"github.com/godbus/dbus/v5"
)

type linuxNotifier struct {
}

func NewNotifier() Notifier {
	return &linuxNotifier{}
}

func (n *linuxNotifier) push(c NotifyContent) {
	cmd := func() error {
		send, err := exec.LookPath("sw-notify-send")
		if err != nil {
			send, err = exec.LookPath("notify-send")
			if err != nil {
				return err
			}
		}

		c := exec.Command(send, c.Title, c.Text, "-i", c.Icon)
		return c.Run()
	}

	conn, err := dbus.SessionBus()
	if err != nil {
		cmd()
		return
	}
	obj := conn.Object("org.freedesktop.Notifications", dbus.ObjectPath("/org/freedesktop/Notifications"))

	call := obj.Call("org.freedesktop.Notifications.Notify", 0, "", uint32(0), c.Icon, c.Title,
		c.Text, []string{}, map[string]dbus.Variant{"string": dbus.MakeVariant(byte(c.urgency))}, int32(-1))
	if call.Err != nil {
		e := cmd()
		if e != nil {
			send, err := exec.LookPath("kdialog")
			if err != nil {
				return
			}
			c := exec.Command(send, "--title", c.Title, "--passivepopup", c.Text, "10", "--icon", c.Icon)
			c.Run()
		}
	}
}
