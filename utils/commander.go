package utils

import (
	"os"
	"os/exec"
	"strings"
)

import log "github.com/cihub/seelog"

/* called when command ends with status ok */
type OnExitFunc func(identifier string) ()

type Commander struct {
	name       string
	args       string
	exitAction OnExitFunc
	identifier string
	cmd        exec.Cmd
}

func NewCommander(name string, args string, exitAction OnExitFunc) *Commander {

	_, err := exec.LookPath(name)
	if err != nil {
		log.Error(name + " not exits")
		return nil
	}

	varArgs := strings.Split(args, " ")
	commander := new(Commander)
	commander.name = name
	commander.args = args
	commander.exitAction = exitAction
	commander.cmd = *exec.Command(name, varArgs...)
	commander.identifier = ""

	return commander
}

/* starts command execution - concurrently */
func (cmder *Commander) Execute(identifier string) {

	if nil != cmder {

		cmder.identifier = identifier

		go func() {

			logStr := cmder.name + " with args: " + cmder.args
			log.Debug("Started command: " + logStr)

			err := cmder.cmd.Run()
			if err != nil {
				log.Info("Stopped command: " + logStr + " with error: " + err.Error())
			} else {
				log.Info("Stopped command: " + logStr + " without error")

				if nil != cmder.exitAction {
					cmder.exitAction(cmder.identifier)
				}
			}
		}()
	}
}

/* sends KILL signal to executed external command */
func (cmder *Commander) Stop() {

	if nil != cmder {

		logStr := cmder.name + " with args: " + cmder.args
		log.Debug("Stopped command manually " + logStr)
		cmder.cmd.Process.Signal(os.Kill)
	}
}

