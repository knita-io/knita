package ui

import (
	"fmt"
	"io"
	"strings"
	"unicode/utf8"

	"github.com/chelnak/ysmrr/pkg/tput"

	executorv1 "github.com/knita-io/knita/api/executor/v1"
)

type RuntimeElement struct {
	*ElementContainer
	ui             *Manager
	tenderID       string
	opts           *executorv1.Opts
	height         int
	message        string
	complete       bool
	err            string
	tendering      bool
	opening        bool
	currentExecs   int
	currentImports int
	currentExports int
}

func NewRuntimeElement(ui *Manager, tenderID string, opts *executorv1.Opts) *RuntimeElement {
	return &RuntimeElement{ui: ui, tenderID: tenderID, opts: opts, height: 1, ElementContainer: NewElementContainer(ui)}
}

func (e *RuntimeElement) ID() string {
	return e.tenderID
}

func (e *RuntimeElement) Update(fc int) {
	e.ElementContainer.Update(fc)
}

func (e *RuntimeElement) Height() int {
	return e.height + e.ElementContainer.Height()
}

func (e *RuntimeElement) Render(writer io.Writer, width int) {
	displayName := e.tenderID
	if e.opts.Tags != nil {
		name, ok := e.opts.Tags["name"]
		if ok {
			displayName = formatUntrustedText(name)
		}
	}
	var text string
	if e.complete {
		if e.err == "" {
			text = fmt.Sprintf("%s: finished\r\n", displayName)
		} else {
			text = fmt.Sprintf("%s: failed: %s\r\n", displayName, e.err)
		}
	} else {
		var states []string
		if e.tendering {
			states = append(states, "locating")
		}
		if e.opening {
			states = append(states, "provisioning")
		}
		if e.currentExecs > 0 {
			states = append(states, "executing")
		}
		if e.currentImports > 0 {
			states = append(states, "importing")
		}
		if e.currentExports > 0 {
			states = append(states, "exporting")
		}
		if len(states) == 0 {
			states = append(states, "idle")
		}
		text = fmt.Sprintf("%s: %s\r\n", displayName, strings.Join(states, ", "))
	}
	if utf8.RuneCountInString(text) > width {
		text = text[:width]
	}
	tput.ClearLine(writer)
	fmt.Fprint(writer, text)
	e.ElementContainer.Render(writer, width)
}

func (e *RuntimeElement) StartTendering() {
	e.tendering = true
	e.ui.notifyUpdate()
}

func (e *RuntimeElement) EndTendering() {
	e.tendering = false
	e.ui.notifyUpdate()
}

func (e *RuntimeElement) StartOpening() {
	e.opening = true
	e.ui.notifyUpdate()
}

func (e *RuntimeElement) EndOpening() {
	e.opening = false
	e.ui.notifyUpdate()
}

func (e *RuntimeElement) StartExec() {
	e.currentExecs++
	e.ui.notifyUpdate()
}

func (e *RuntimeElement) EndExec() {
	e.currentExecs--
	e.ui.notifyUpdate()
}

func (e *RuntimeElement) StartImport() {
	e.currentImports++
	e.ui.notifyUpdate()
}

func (e *RuntimeElement) EndImport() {
	e.currentImports--
	e.ui.notifyUpdate()
}

func (e *RuntimeElement) StartExport() {
	e.currentExports++
	e.ui.notifyUpdate()
}

func (e *RuntimeElement) EndExport() {
	e.currentExports--
	e.ui.notifyUpdate()
}

func (e *RuntimeElement) SetMessage(message string) {
	e.message = message
	e.ui.notifyUpdate()
}

func (e *RuntimeElement) Complete(err string) {
	e.err = err
	e.complete = true
	e.ui.notifyUpdate()
}
