package core

import (
	"trojan/util"
)

func TrojanGFWStop() error {
	if err := util.ExecCommand("docker stop trojan-panel-trojanGFW"); err != nil {
		return err
	} else {
		return nil
	}
}

func TrojanGFWRestart() error {
	if err := util.ExecCommand("docker restart trojan-panel-trojanGFW"); err != nil {
		return err
	} else {
		return nil
	}
}

func TrojanGFWStatus() string {
	return util.ExecCommandWithResult("docker ps | grep trojan-panel-trojanGFW")
}

func TrojanGOStop() error {
	if err := util.ExecCommand("docker stop trojan-panel-trojanGO"); err != nil {
		return err
	} else {
		return nil
	}
}

func TrojanGORestart() error {
	if err := util.ExecCommand("docker restart trojan-panel-trojanGO"); err != nil {
		return err
	} else {
		return nil
	}
}

func TrojanGOStatus() string {
	return util.ExecCommandWithResult("docker ps | grep trojan-panel-trojanGO")
}
