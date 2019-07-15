package pkg

import (
	"fmt"
	"github.com/jroimartin/gocui"
	"strings"
	"time"
)

// Move view cursor to the bottom
func moveViewCursorDown(g *gocui.Gui, v *gocui.View, allowEmpty bool) error {
	cx, cy := v.Cursor()
	ox, oy := v.Origin()
	nextLine, err := getNextViewLine(g, v)
	if err != nil {
		return err
	}
	if !allowEmpty && nextLine == "" {
		return nil
	}
	if err := v.SetCursor(cx, cy+1); err != nil {
		if err := v.SetOrigin(ox, oy+1); err != nil {
			return err
		}
	}
	return nil
}

// Move view cursor to the top
func moveViewCursorUp(g *gocui.Gui, v *gocui.View, dY int) error {
	ox, oy := v.Origin()
	cx, cy := v.Cursor()
	if cy > dY {
		if err := v.SetCursor(cx, cy-1); err != nil && oy > 0 {
			if err := v.SetOrigin(ox, oy-1); err != nil {
				return err
			}
		}
	}
	return nil
}

// Get view line (relative to the cursor)
func getViewLine(g *gocui.Gui, v *gocui.View) (string, error) {
	var l string
	var err error

	_, cy := v.Cursor()
	if l, err = v.Line(cy); err != nil {
		l = ""
	}

	return l, err
}

// Get the next view line (relative to the cursor)
func getNextViewLine(g *gocui.Gui, v *gocui.View) (string, error) {
	var l string
	var err error

	_, cy := v.Cursor()
	if l, err = v.Line(cy + 1); err != nil {
		l = ""
	}

	return l, err
}

// Set view cursor to default line
func setViewCursorToLine(g *gocui.Gui, v *gocui.View, lines []string, selLine string) error {
	ox, _ := v.Origin()
	cx, _ := v.Cursor()
	for y, line := range lines {
		if line == selLine {
			if err := v.SetCursor(ox, y); err != nil {
				if err := v.SetOrigin(cx, y); err != nil {
					return err
				}
			}
		}
	}
	return nil
}

// Get pod name form line
func getPodNameFromLine(line string) string {
	if line == "" {
		return ""
	}

	i := strings.Index(line, " ")
	if i == -1 {
		return line
	}

	return line[0:i]
}

// Get selected pod
func getSelectedPod(g *gocui.Gui) (string, error) {
	v, err := g.View("pods")
	if err != nil {
		return "", err
	}
	l, err := getViewLine(g, v)
	if err != nil {
		return "", err
	}
	p := getPodNameFromLine(l)

	return p, nil
}

// Show views logs
func showViewPodsLogs(g *gocui.Gui) error {
	vn := "logs"

	switch LOG_MOD {
	case "pod":
		// Get current selected pod
		p, err := getSelectedPod(g)
		if err != nil {
			return err
		}

		// Display pod containers
		vLc, err := g.View(vn + "-containers")
		if err != nil {
			return err
		}
		vLc.Clear()
		for _, c := range getPodContainers(p) {
			fmt.Fprintln(vLc, c)
		}
		vLc.SetCursor(0, 0)

		// Display logs
		refreshPodsLogs(g)
	}

	debug(g, "Action: Show view logs")
	g.SetViewOnTop(vn)
	g.SetViewOnTop(vn + "-containers")
	g.SetCurrentView(vn)

	return nil
}

// Refresh pods logs
func refreshPodsLogs(g *gocui.Gui) error {
	vn := "logs"

	// Get current selected pod
	p, err := getSelectedPod(g)
	if err != nil {
		return err
	}

	vLc, err := g.View(vn + "-containers")
	if err != nil {
		return err
	}

	c, err := getViewLine(g, vLc)
	if err != nil {
		return err
	}

	vL, err := g.View(vn)
	if err != nil {
		return err
	}
	getPodContainerLogs(p, c, vL)

	return nil
}

// Display error
func displayError(g *gocui.Gui, e error) error {
	lMaxX, lMaxY := g.Size()
	minX := lMaxX / 6
	minY := lMaxY / 6
	maxX := 5 * (lMaxX / 6)
	maxY := 5 * (lMaxY / 6)

	if v, err := g.SetView("errors", minX, minY, maxX, maxY); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}

		// Settings
		v.Title = " ERROR "
		v.Frame = true
		v.Wrap = true
		v.Autoscroll = true
		v.BgColor = gocui.ColorRed
		v.FgColor = gocui.ColorWhite

		// Content
		v.Clear()
		fmt.Fprintln(v, e.Error())

		// Send to forground
		g.SetCurrentView(v.Name())
	}

	return nil
}

// Hide error box
func hideError(g *gocui.Gui) {
	g.DeleteView("errors")
}

// Display confirmation message
func displayConfirmation(g *gocui.Gui, m string) error {
	lMaxX, lMaxY := g.Size()

	if v, err := g.SetView("confirmation", -1, lMaxY-3, lMaxX, lMaxY-1); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}

		// Settings
		v.Frame = false

		// Content
		fmt.Fprintln(v, textPadCenter(m, lMaxX))

		// Auto-hide message
		hide := func() {
			hideConfirmation(g)
		}
		time.AfterFunc(time.Duration(2)*time.Second, hide)
	}

	return nil
}

// delete pod confirmation message
func deleteConfirmation(g *gocui.Gui, m string) error {
	lMaxX, lMaxY := g.Size()

	if v, err := g.SetView("delconfirmation", lMaxX/2-25, lMaxY/2-3, lMaxX/2+25, lMaxY/2+3); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		// Settings
		v.Frame = true
		v.Wrap = true

		// Content
		fmt.Fprintln(v, m+" pod will be deleted, press Enter to confirm, or 'q' to cancel. ")
		g.SetCurrentView(v.Name())
	}
	return nil
}

// Hide confirmation message
func hideConfirmation(g *gocui.Gui) {
	g.DeleteView("confirmation")
}