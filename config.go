package mmbot

import "mmbot/mmhook"

type Config struct {
	*mmhook.Config
	UserName         string
	OverrideUserName string
	IconURL          string
	DisableServer    bool
}
