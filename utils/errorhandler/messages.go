package errorhandler

import (
	"errors"
	"fmt"

	"github.com/enescakir/emoji"
)

var (
	ErrInterrupt = errors.New("^C")
	ErrExist     = errors.New("already exist")
)

var (
	ExitMessage     = fmt.Errorf("\n%s%s\n", "Program Exited", emoji.Parse(":exclamation:"))
	DoneMessage     = fmt.Errorf("\n%s%s\n", "Done", emoji.Parse(":sparkles:"))
	CompleteMessage = fmt.Errorf("%s%s\n", "Completed", emoji.Parse(":sparkles:"))
	WaveMessage     = emoji.WavingHand
	Exclamation     = emoji.ExclamationMark
)
