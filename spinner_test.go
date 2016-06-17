package spinner

import (
	"bytes"
	"errors"
	"io/ioutil"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestSpinner_Setters_ReturnsSpinner(t *testing.T) {
	test := assert.New(t)

	spinner := New()
	{
		returns := spinner.SetStatus("x")
		test.Equal(returns, spinner)
	}
	{
		returns := spinner.SetFrames("x")
		test.Equal(returns, spinner)
	}
	{
		returns := spinner.SetEmptyFrame("x")
		test.Equal(returns, spinner)
	}
	{
		returns := spinner.SetInterval(time.Millisecond)
		test.Equal(returns, spinner)
	}
}
func TestSpinner_SetStatus_SetsSpinnerStatus(t *testing.T) {
	test := assert.New(t)

	spinner := New()
	spinner.SetStatus("x")
	test.Equal(spinner.status, "x")
}

func TestSpinner_SetFrames_SetsSpinnerFrames(t *testing.T) {
	test := assert.New(t)

	spinner := New()
	spinner.SetFrames("x", "y")
	test.Equal(spinner.frames, []string{"x", "y"})
}

func TestSpinner_SetInterval_SetsSpinnerInterval(t *testing.T) {
	test := assert.New(t)

	spinner := New()
	spinner.SetInterval(time.Second)
	test.Equal(spinner.interval, time.Second)
}

func TestSpinner_SetEmptyFrame_SetsSpinnerEmptyFrame(t *testing.T) {
	test := assert.New(t)

	spinner := New()
	spinner.SetEmptyFrame("x")
	test.Equal(spinner.emptyFrame, "x")
}

func TestSpinner_SetOutput_SetsSpinnerOutput(t *testing.T) {
	test := assert.New(t)

	spinner := New()
	spinner.SetOutput(ioutil.Discard)
	test.Equal(spinner.output, ioutil.Discard)
}

func TestSpinner_StartStop_After400ms(t *testing.T) {
	test := assert.New(t)

	spinner, output := getTestableSpinner()

	spinner.Start()
	time.Sleep(time.Millisecond * 400)
	spinner.Stop()

	test.Equal(
		"\rx...1\rx...2\rx...0\n",
		output.String(),
	)
}

func TestSpinner_StartStop_Immidiately(t *testing.T) {
	test := assert.New(t)

	spinner, output := getTestableSpinner()

	spinner.Start()
	spinner.Stop()

	test.Equal(
		"\rx...0\n",
		output.String(),
	)
}

func TestSpinner_Call_RunsFuncAndDoSpin(t *testing.T) {
	test := assert.New(t)

	spinner, output := getTestableSpinner()

	err := spinner.Call(func() error {
		time.Sleep(time.Millisecond * 400)
		return nil
	})

	test.NoError(err)

	test.Equal(
		"\rx...1\rx...2\rx...0\n",
		output.String(),
	)
}

func TestSpinner_Call_ReturnsError(t *testing.T) {
	test := assert.New(t)

	spinner, output := getTestableSpinner()

	err := spinner.Call(func() error {
		time.Sleep(time.Millisecond * 400)
		return errors.New("blah")
	})

	test.EqualError(err, "blah")

	test.Equal(
		"\rx...1\rx...2\rx...0\n",
		output.String(),
	)
}

func getTestableSpinner() (*Spinner, *bytes.Buffer) {
	output := bytes.NewBuffer(nil)

	spinner := New()
	spinner.SetOutput(output)
	spinner.SetFrames("1", "2", "3", "4")
	spinner.SetEmptyFrame("0")
	spinner.SetStatus("x...")

	return spinner, output
}
