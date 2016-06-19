package spinner

import (
	"fmt"
	"io"
	"os"
	"strings"
	"sync"
	"syscall"
	"time"
	"unsafe"
)

var (
	// DefaultOutput represents io.Writer which will be used for new spinners.
	DefaultOutput = os.Stderr

	// DefaultEmptyFrame represents string representation of frame which will
	// be showed when spinner is finished.
	DefaultEmptyFrame = " "

	// DefaultFrames represents set of frames which will be iterated in cycle.
	DefaultFrames = []string{"—", "\\", "|", "/", "-", "\\", "|", "/"}

	// DefaultInterval represents time.Duration of spin iteration.
	DefaultInterval = time.Millisecond * 300
)

type Spinner struct {
	output     io.Writer
	status     string
	frames     []string
	interval   time.Duration
	process    *sync.WaitGroup
	active     bool
	emptyFrame string
	iteration  int
}

func New() *Spinner {
	spinner := &Spinner{
		output:     DefaultOutput,
		frames:     DefaultFrames,
		interval:   DefaultInterval,
		process:    &sync.WaitGroup{},
		emptyFrame: DefaultEmptyFrame,
	}

	return spinner
}

func (spinner *Spinner) SetOutput(output io.Writer) *Spinner {
	spinner.output = output
	return spinner
}

func (spinner *Spinner) SetStatus(status string) *Spinner {
	spinner.status = status
	return spinner
}

func (spinner *Spinner) SetFrames(frame ...string) *Spinner {
	spinner.frames = frame
	return spinner
}

func (spinner *Spinner) SetInterval(interval time.Duration) *Spinner {
	spinner.interval = interval
	return spinner
}

func (spinner *Spinner) SetEmptyFrame(frame string) *Spinner {
	spinner.emptyFrame = frame
	return spinner
}

func (spinner *Spinner) IsActive() bool {
	return spinner.active
}

func (spinner *Spinner) Spin() {
	spinner.spin()
	if !spinner.active {
		fmt.Fprint(spinner.output, "\n")
		spinner.process.Done()
	}
}

func (spinner *Spinner) SetActive(active bool) *Spinner {
	spinner.active = active
	return spinner
}

func (spinner *Spinner) spin() {
	frame := spinner.emptyFrame
	if spinner.active {
		frame = spinner.frames[spinner.iteration]
	}

	fmt.Fprintf(
		spinner.output,
		"\r"+spinner.status+frame+getSpinnerSuffix(len(spinner.status)),
	)
}

func (spinner *Spinner) Start() {
	spinner.active = true
	spinner.process.Add(1)

	go func() {
		spinner.iteration = 0
		for spinner.active {
			spinner.Spin()
			spinner.iteration = (spinner.iteration + 1) % len(spinner.frames)
			time.Sleep(spinner.interval)
		}

		spinner.Spin()
	}()
}

func (spinner *Spinner) Stop() {
	spinner.active = false
	spinner.process.Wait()
}

func (spinner *Spinner) Call(methods ...func() error) error {
	spinner.Start()
	defer spinner.Stop()

	for _, method := range methods {
		err := method()
		if err != nil {
			return err
		}
	}

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
