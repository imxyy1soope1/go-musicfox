//go:build darwin

package utils

import (
	"os/exec"
)

type osxNotifier struct {
}

func NewNotifier() Notifier {
	return &osxNotifier{}
}

func (_ *osxNotifier) getNotifierCmd() string {
	localDir := GetLocalDataDir()
	notifierPath := path.Join(localDir, "musicfox-notifier.app")
	if _, err := os.Stat(notifierPath); os.IsNotExist(err) {
		err = CopyDirFromEmbed("embed/musicfox-notifier.app", notifierPath)
		if err != nil {
			log.Printf("copy musicfox-notifier.app failed, err: %+v", errors.WithStack(err))
		}
	} else if err != nil {
		log.Printf("musicfox-notifier.app status err: %+v", errors.WithStack(err))
	}

	return notifierPath + "/Contents/MacOS/musicfox-notifier"
}

func (n *osxNotifier) push(c NotifyContent) {
	cmd := func() exec.Cmd {
		cmdPath := n.getNotifierCmd()
		if FileOrDirectoryExists(cmdPath) {
			var args = []string{"-title", n.appName, "-message", text, "-subtitle", title, "-contentImage", iconPath}
			if redirectUrl != "" {
				args = append(args, "-open", redirectUrl)
			}
			if groupId != "" {
				args = append(args, "-group", groupId)
			}
			return exec.Command(cmdPath, args...)
		} else if notificator.CheckMacOSVersion() {
			title = strings.Replace(title, `"`, `\"`, -1)
			text = strings.Replace(text, `"`, `\"`, -1)

			notification := fmt.Sprintf("display notification \"%s\" with title \"%s\" subtitle \"%s\"", text, n.appName, title)
			return exec.Command("osascript", "-e", notification)
		} else {
			return exec.Command("growlnotify", "-n", n.appName, "--image", iconPath, "-m", title, "--url", redirectUrl)
		}
	}
	cmd().Run()
}
