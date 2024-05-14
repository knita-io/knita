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

	executorv1 "github.com/knita-io/knita/api/executor/v1"
	"github.com/knita-io/knita/internal/event"
)

const targetFPS = 60

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
	stream      event.Stream
	done        func()
	writer      io.Writer
	container   *ElementContainer
	ticker      *time.Ticker
	exitC       chan struct{}
	doneC       chan struct{}
	started     bool
	mu          sync.Mutex
	needsUpdate bool
}

func NewManager(stream event.Stream) *Manager {
	ui := &Manager{
		stream:      stream,
		writer:      os.Stdout,
		needsUpdate: true,
		doneC:       make(chan struct{}),
		exitC:       make(chan struct{}),
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

func (ui *Manager) onEventCallback(event *executorv1.Event) {
	ui.mu.Lock()
	defer ui.mu.Unlock()

	switch p := event.Payload.(type) {
	case *executorv1.Event_RuntimeOpened:
		ui.AddChildElement(NewRuntimeElement(ui, p.RuntimeOpened.RuntimeId, p.RuntimeOpened.Opts))
	case *executorv1.Event_RuntimeClosed:
		withElement(ui, p.RuntimeClosed.RuntimeId, func(ele *RuntimeElement) {
			ele.Complete("") // TODO would be nice to get the error from a failed runtime
		})
	case *executorv1.Event_ExecStart:
		withElement(ui, p.ExecStart.RuntimeId, func(ele *RuntimeElement) {
			ele.StartExec()
			ele.AddChildElement(NewExecElement(ui, p.ExecStart.ExecId, p.ExecStart.Opts))
		})
	case *executorv1.Event_ExecEnd:
		withElement(ui, p.ExecEnd.RuntimeId, func(ele *RuntimeElement) {
			ele.EndExec()
		})
		withElement(ui, p.ExecEnd.ExecId, func(ele *ExecElement) {
			ele.Complete(p.ExecEnd.ExitCode, p.ExecEnd.Error)
		})
	case *executorv1.Event_ImportStart:
		withElement(ui, p.ImportStart.RuntimeId, func(ele *RuntimeElement) {
			ele.StartImport()
		})
	case *executorv1.Event_ImportEnd:
		withElement(ui, p.ImportEnd.RuntimeId, func(ele *RuntimeElement) {
			ele.EndImport()
		})
	case *executorv1.Event_ExportStart:
		withElement(ui, p.ExportStart.RuntimeId, func(ele *RuntimeElement) {
			ele.StartExport()
		})
	case *executorv1.Event_ExportEnd:
		withElement(ui, p.ExportEnd.RuntimeId, func(ele *RuntimeElement) {
			ele.EndExport()
		})
	case *executorv1.Event_Stdout:
		switch s := p.Stdout.Source.Source.(type) {
		case *executorv1.LogOutEventSource_Exec:
			withElement(ui, s.Exec.ExecId, func(ele *ExecElement) {
				ele.SetMessage(string(p.Stdout.Data))
			})
		}
	case *executorv1.Event_Stderr:
		switch s := p.Stderr.Source.Source.(type) {
		case *executorv1.LogOutEventSource_Exec:
			withElement(ui, s.Exec.ExecId, func(ele *ExecElement) {
				ele.SetMessage(string(p.Stderr.Data))
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
