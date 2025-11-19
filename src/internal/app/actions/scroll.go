package actions

import (
	"github.com/jroimartin/gocui"
)

func ScrollDown(v *gocui.View) error {
	if v == nil {
		return nil
	}

	ox, oy := v.Origin()
	_, height := v.Size()
	if height <= 0 {
		return nil
	}

	lines := len(v.BufferLines())
	maxY := lines - height
	if maxY < 0 {
		maxY = 0
	}
	if oy >= maxY {
		return nil
	}

	return v.SetOrigin(ox, oy+1)
}

func ScrollUp(v *gocui.View) error {
	if v == nil {
		return nil
	}

	ox, oy := v.Origin()
	if oy > 1 {
		return v.SetOrigin(ox, oy-1)
	}
	return nil
}

func JumpPageUp (v *gocui.View) error {
	if v == nil {
		return  nil
	}
	ox, _ := v.Origin()
	return v.SetOrigin(ox, 0)
} 
