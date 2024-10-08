package ui

import (
	"io"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/alessio/shellescape"
	"github.com/chelnak/ysmrr/pkg/tput"
	"golang.org/x/term"

	builtinv1 "github.com/knita-io/knita/api/events/builtin/v1"
	"github.com/knita-io/knita/internal/event"
)

const targetFPS = 10

type Element interface {
	ID() string
	Update(fc int)
	Height() int
	Render(writer io.Writer, width int)
}

type Container interface {
	AddChildElement(ele Element)
	GetChildElement(id string) (Element, bool)
}

type Manager struct {
	stream              event.Stream
	done                func()
	writer              io.Writer
	container           *ElementContainer
	ticker              *time.Ticker
	exitC               chan struct{}
	doneC               chan struct{}
	started             bool
	mu                  sync.Mutex
	needsUpdate         bool
	runtimeIDToTenderID map[string]string
}

func NewManager(stream event.Stream) *Manager {
	ui := &Manager{
		stream:              stream,
		writer:              os.Stdout,
		needsUpdate:         true,
		doneC:               make(chan struct{}),
		exitC:               make(chan struct{}),
		runtimeIDToTenderID: make(map[string]string),
	}
	ui.container = NewElementContainer(ui)
	return ui
}

func (ui *Manager) Start() {
	ui.done = ui.stream.Subscribe(ui.onEventCallback)
	ui.ticker = time.NewTicker(time.Second / targetFPS)
	go ui.mainLoop()
	ui.started = true
}

func (ui *Manager) Stop() {
	if !ui.started {
		return
	}
	ui.started = false
	ui.done()
	ui.ticker.Stop()
	close(ui.exitC)
	<-ui.doneC
}

func (ui *Manager) AddChildElement(ele Element) {
	ui.container.AddChildElement(ele)
}

func (ui *Manager) GetChildElement(id string) (Element, bool) {
	return ui.container.GetChildElement(id)
}

func (ui *Manager) notifyUpdate() {
	ui.needsUpdate = true
}

func (ui *Manager) mainLoop() {
	tput.Civis(ui.writer)
	defer tput.Cnorm(ui.writer)

	var lastWidth, lastHeight, fc int
	loop := func() {
		ui.mu.Lock()
		defer ui.mu.Unlock()

		fc++
		// Capture terminal width
		width, _, _ := term.GetSize(int(os.Stdout.Fd())) // TODO correct FD?
		if width <= 0 || width > 140 {
			width = 140
		}
		// Update all elements
		ui.container.Update(fc)
		// Capture content height
		height := ui.container.Height()
		// Determine if an update is needed
		if !(ui.needsUpdate || lastWidth != width || lastHeight != height) {
			return
		}
		ui.needsUpdate = false
		// Rewind to the top of the view port
		if lastHeight > 0 {
			tput.Cuu(ui.writer, lastHeight)
		}
		// Do a render pass
		ui.container.Render(ui.writer, width)
		lastWidth = width
		lastHeight = height
		// TODO: A virtual screen buffer would help us only redraw what is necessary
		// NOTE: If we ever remove elements e.g. height<lastHeight then we need to clear the delta from the screen
	}
	for {
		select {
		case <-ui.exitC:
			loop()
			close(ui.doneC)
			return
		case <-ui.ticker.C:
			loop()
		}
	}
}

func (ui *Manager) onEventCallback(event *event.Event) {
	ui.mu.Lock()
	defer ui.mu.Unlock()

	switch p := event.Payload.(type) {
	case *builtinv1.RuntimeTenderStartEvent:
		ui.AddChildElement(NewRuntimeElement(ui, p.TenderId, p.Opts))
		withElement(ui, p.TenderId, func(ele *RuntimeElement) {
			ele.StartTendering()
		})
	case *builtinv1.RuntimeSettlementEndEvent:
		ui.runtimeIDToTenderID[p.RuntimeId] = p.TenderId
		withElement(ui, p.TenderId, func(ele *RuntimeElement) {
			ele.EndTendering()
			ele.StartOpening()
		})
	case *builtinv1.RuntimeOpenEndEvent:
		withElement(ui, ui.runtimeIDToTenderID[p.RuntimeId], func(ele *RuntimeElement) {
			ele.EndOpening()
		})
	case *builtinv1.RuntimeCloseEndEvent:
		withElement(ui, ui.runtimeIDToTenderID[p.RuntimeId], func(ele *RuntimeElement) {
			switch s := p.Status.(type) {
			case *builtinv1.RuntimeCloseEndEvent_Result:
				ele.Complete("")
			case *builtinv1.RuntimeCloseEndEvent_Error:
				ele.Complete(s.Error.Message)
			}
			delete(ui.runtimeIDToTenderID, p.RuntimeId)
		})
	case *builtinv1.ExecStartEvent:
		withElement(ui, ui.runtimeIDToTenderID[p.RuntimeId], func(ele *RuntimeElement) {
			ele.StartExec()
			ele.AddChildElement(NewExecElement(ui, p.ExecId, p.Opts))
		})
	case *builtinv1.ExecEndEvent:
		withElement(ui, ui.runtimeIDToTenderID[p.RuntimeId], func(ele *RuntimeElement) {
			ele.EndExec()
		})
		withElement(ui, p.ExecId, func(ele *ExecElement) {
			switch s := p.Status.(type) {
			case *builtinv1.ExecEndEvent_Result:
				ele.Complete(s.Result.ExitCode, "")
			case *builtinv1.ExecEndEvent_Error:
				ele.Complete(-1, s.Error.Message)
			}
		})
	case *builtinv1.ImportStartEvent:
		withElement(ui, ui.runtimeIDToTenderID[p.RuntimeId], func(ele *RuntimeElement) {
			ele.StartImport()
		})
	case *builtinv1.ImportEndEvent:
		withElement(ui, ui.runtimeIDToTenderID[p.RuntimeId], func(ele *RuntimeElement) {
			ele.EndImport()
		})
	case *builtinv1.ExportStartEvent:
		withElement(ui, ui.runtimeIDToTenderID[p.RuntimeId], func(ele *RuntimeElement) {
			ele.StartExport()
		})
	case *builtinv1.ExportEndEvent:
		withElement(ui, ui.runtimeIDToTenderID[p.RuntimeId], func(ele *RuntimeElement) {
			ele.EndExport()
		})
	case *builtinv1.StdoutEvent:
		switch s := p.Source.Source.(type) {
		case *builtinv1.LogEventSource_Exec:
			withElement(ui, s.Exec.ExecId, func(ele *ExecElement) {
				ele.SetMessage(string(p.Data))
			})
		}
	case *builtinv1.StderrEvent:
		switch s := p.Source.Source.(type) {
		case *builtinv1.LogEventSource_Exec:
			withElement(ui, s.Exec.ExecId, func(ele *ExecElement) {
				ele.SetMessage(string(p.Data))
			})
		}
	}
}

// withElement locates the element with id, and if it exists, and is of type K, invokes fn. Otherwise, nops.
func withElement[K Element](ui *Manager, id string, fn func(K)) {
	ele, ok := ui.GetChildElement(id)
	if ok {
		eleT, ok := ele.(K)
		if ok {
			fn(eleT)
		}
	}
}

func formatUntrustedText(text string) string {
	return shellescape.StripUnsafe(strings.Trim(text, "\r\n\t"))
}
