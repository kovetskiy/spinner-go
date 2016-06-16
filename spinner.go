package spinner

import (
	"fmt"
	"os"
	"strings"
	"sync"
	"syscall"
	"time"
	"unsafe"
)

var (
	defaultZeroFrame = " "
	defaultFrames    = []string{"—", "\\", "|", "/", "-", "\\", "|", "/"}
	defaultInterval  = time.Millisecond * 300
)

type Spinner struct {
	status    string
	frames    []string
	interval  time.Duration
	process   *sync.WaitGroup
	active    bool
	zeroframe string
	iteration int
}

func New() *Spinner {
	spinner := &Spinner{
		frames:    defaultFrames,
		interval:  defaultInterval,
		process:   &sync.WaitGroup{},
		zeroframe: defaultZeroFrame,
	}

	return spinner
}

func (spinner *Spinner) Schedule() {
	spinner.Update()
	if spinner.active {
		spinner.NewLine()
		spinner.process.Done()
	}
}

func (spinner *Spinner) Update() {
	frame := spinner.zeroframe
	if spinner.active {
		frame = spinner.frames[spinner.iteration]
	}

	fmt.Fprintf(
		os.Stdout,
		"\r"+spinner.status+frame+getSpinnerSuffix(len(spinner.status)),
	)
}

func (spinner *Spinner) NewLine() {
	fmt.Fprint(os.Stdout, "\n")
}

func (spinner *Spinner) Start() {
	spinner.active = true
	spinner.process.Add(1)

	go func() {
		spinner.iteration = 0
		for spinner.active {
			spinner.Schedule()
			spinner.iteration = (spinner.iteration + 1) % len(spinner.frames)
			time.Sleep(spinner.interval)
		}

		spinner.Schedule()
	}()
}

func (spinner *Spinner) Stop() {
	spinner.active = false
	spinner.process.Wait()
}

func (spinner *Spinner) Call(status string, methods ...func() error) error {
	spinner.Start()

	for _, method := range methods {
		err := method()
		if err != nil {
			return err
		}
	}

	spinner.Stop()

	return nil
}

func getTerminalWidth() int {
	term, err := os.Open("/dev/tty")
	if err != nil {
		term = os.Stdin
	}

	window := struct {
		Rows    uint16
		Columns uint16
		X       uint16
		Y       uint16
	}{}

	result, _, err := syscall.Syscall(
		syscall.SYS_IOCTL,
		term.Fd(),
		uintptr(syscall.TIOCGWINSZ),
		uintptr(unsafe.Pointer(&window)),
	)
	if int(result) == -1 || err != nil {
		return 0
	}

	return int(window.Columns)
}

func getSpinnerSuffix(length int) string {
	suffix := getTerminalWidth() - length
	if suffix > 0 {
		return strings.Repeat(" ", suffix)
	}

	return ""
}
