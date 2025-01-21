package shell

import (
	"os/signal"
	"syscall"
)

func Init() {
	signal.Ignore(syscall.SIGINT, syscall.SIGQUIT)
}
