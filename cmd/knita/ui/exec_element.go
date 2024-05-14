package ui

import (
	"fmt"
	"io"
	"strings"
	"time"
	"unicode/utf8"

	"github.com/chelnak/ysmrr/pkg/tput"

	executorv1 "github.com/knita-io/knita/api/executor/v1"
)

type ExecElement struct {
	ui           *Manager
	execID       string
	opts         *executorv1.ExecOpts
	height       int
	start        time.Time
	runTime      time.Duration
	complete     bool
	message      string
	exitCode     int32
	err          string
	spinnerChars []string
	spinnerFrame int
}

func NewExecElement(ui *Manager, execID string, opts *executorv1.ExecOpts) *ExecElement {
	return &ExecElement{
		ui:           ui,
		execID:       execID,
		opts:         opts,
		height:       1,
		spinnerChars: []string{"⠋", "⠙", "⠹", "⠸", "⠼", "⠴", "⠦", "⠧", "⠇", "⠏"},
		start:        time.Now(),
	}
}

func (e *ExecElement) ID() string {
	return e.execID
}

func (e *ExecElement) Update(fc int) {
	if !e.complete {
		frame := (fc % targetFPS) / (targetFPS / len(e.spinnerChars))
		if e.spinnerFrame != frame {
			e.spinnerFrame = frame
			e.ui.notifyUpdate()
		}
		runTime := time.Now().Sub(e.start)
		if e.runTime.Seconds() != runTime.Seconds() {
			e.runTime = runTime
			e.ui.notifyUpdate()
		}
	}
}

func (e *ExecElement) Height() int {
	return e.height
}

func (e *ExecElement) Render(writer io.Writer, width int) {
	displayName := e.execID
	if e.opts.Tags != nil {
		name, ok := e.opts.Tags["name"]
		if ok {
			displayName = formatUntrustedText(name)
		}
	}
	runTime := fmt.Sprintf("%s", e.runTime.Round(time.Millisecond*100))

	var text string
	if e.complete {
		if e.exitCode == 0 && e.err == "" {
			text = fmt.Sprintf(" ✓ %s (%s)\r\n", displayName, runTime)
		} else {
			text = fmt.Sprintf(" ✗ %s: %d - %s\r\n", displayName, e.exitCode, e.err)
		}
	} else {
		maxMessageWidth := width - len(e.spinnerChars[e.spinnerFrame]) - len(displayName) - len(runTime) - 9
		message := formatUntrustedText(e.message)
		if len(message) > maxMessageWidth {
			message = message[:maxMessageWidth]
		}
		padding := maxMessageWidth - len(message)
		text = fmt.Sprintf(" %s %s: %s%s     %s\r\n", e.spinnerChars[e.spinnerFrame], displayName, message, strings.Repeat(" ", padding), runTime)
	}

	if utf8.RuneCountInString(text) > width { // TODO use this for all strlen
		text = text[:width]
	}
	tput.ClearLine(writer)
	fmt.Fprint(writer, text)
}

func (e *ExecElement) SetMessage(message string) {
	e.message = message
	e.ui.notifyUpdate()
}

func (e *ExecElement) Complete(exitCode int32, err string) {
	e.exitCode = exitCode
	e.err = err
	e.complete = true
	e.ui.notifyUpdate()
}
