package mmcli

import (
	"errors"
	"os/exec"
	"strings"
)

type Mmcli struct {
	exec ExecCommand
}

type ExecCommand func(name string, arg ...string) ExecCommandOutput

type ExecCommandOutput interface {
	Output() ([]byte, error)
}

var execCommandFunc = func(name string, arg ...string) ExecCommandOutput {
	return exec.Command(name, arg...)
}

func Provide() *Mmcli {
	return &Mmcli{
		exec: execCommandFunc,
	}
}

func (b *Mmcli) Exec(args ...string) ([]byte, error) {
	cmd := b.exec("mmcli", append([]string{"-J"}, args...)...)
	buf, err := cmd.Output()
	if err != nil {
		var errInfo *exec.ExitError
		ok := errors.As(err, &errInfo)
		if ok {
			return nil, errors.New(strings.Trim(string(errInfo.Stderr), "\n "))
		}
		return nil, err
	}

	return buf, nil
}

func (b *Mmcli) ExecModem(modem string, args ...string) ([]byte, error) {
	if len(modem) == 0 {
		buf, _ := b.exec("bash", "-c", "mmcli -L -J | jq -r '.\"modem-list\"[0]'").Output()
		// contains "null" as string if no modem was found
		modem = strings.Trim(string(buf), "\n ")
	}

	return b.Exec(append([]string{"-m", modem}, args...)...)
}
