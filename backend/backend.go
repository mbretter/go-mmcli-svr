package backend

type Backend interface {
	Exec(arg ...string) ([]byte, error)
	ExecModem(modem string, arg ...string) ([]byte, error)
}
