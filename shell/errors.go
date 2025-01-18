package shell

import "errors"

// ErrExit signals that the shell should terminate.
var ErrExit = errors.New("exit")
