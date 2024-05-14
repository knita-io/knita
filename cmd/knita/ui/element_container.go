package ui

import "io"

type ElementContainer struct {
	ui           *Manager
	children     []Element
	childrenByID map[string]Element
}

func NewElementContainer(ui *Manager) *ElementContainer {
	return &ElementContainer{ui: ui, childrenByID: map[string]Element{}}
}

func (e *ElementContainer) AddChildElement(ele Element) {
	e.children = append(e.children, ele)
	e.childrenByID[ele.ID()] = ele
	e.ui.notifyUpdate()
}

func (e *ElementContainer) GetChildElement(id string) (Element, bool) {
	ele, ok := e.childrenByID[id]
	if ok {
		return ele, ok
	}
	for _, ele := range e.children {
		container, ok := ele.(Container)
		if ok {
			ele, ok := container.GetChildElement(id)
			if ok {
				return ele, ok
			}
		}
	}
	return nil, false
}

func (e *ElementContainer) Update(fc int) {
	for _, ele := range e.children {
		ele.Update(fc)
	}
}

func (e *ElementContainer) Height() int {
	var totalHeight int
	for _, ele := range e.children {
		totalHeight += ele.Height()
	}
	return totalHeight
}

func (e *ElementContainer) Render(writer io.Writer, width int) {
	for _, ele := range e.children {
		ele.Render(writer, width)
	}
}
