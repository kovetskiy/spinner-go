package spinner

import (
	"io"
	"time"
)

var (
	spinner = New()
)

func SetActive(active bool) *Spinner {
	spinner.active = active
	return spinner
}

func SetOutput(output io.Writer) *Spinner {
	spinner.output = output
	return spinner
}

func SetStatus(status string) *Spinner {
	spinner.status = status
	return spinner
}

func SetFrames(frame ...string) *Spinner {
	spinner.frames = frame
	return spinner
}

func SetInterval(interval time.Duration) *Spinner {
	spinner.interval = interval
	return spinner
}

func SetEmptyFrame(frame string) *Spinner {
	spinner.emptyFrame = frame
	return spinner
}

func IsActive() bool {
	return spinner.active
}

func Spin() {
	spinner.Spin()
}

func Start() {
	spinner.Start()
}

func Stop() {
	spinner.Stop()
}

func Call(methods ...func() error) error {
	return spinner.Call(methods...)
}
