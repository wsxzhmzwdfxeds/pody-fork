package pkg

import (
	"fmt"
	"github.com/jroimartin/gocui"
)

var DEBUG_DISPLAYED bool = false
var NAMESPACES_DISPLAYED bool = false
var Current_ns string
var LOG_MOD string = "pod"

// Global action: Quit
func actionGlobalQuit(g *gocui.Gui, v *gocui.View) error {
	return gocui.ErrQuit
}

// Global action: Toggle debug
func actionGlobalToggleViewDebug(g *gocui.Gui, v *gocui.View) error {
	vn := "debug"

	if !DEBUG_DISPLAYED {
		debug(g, "Action: Display debug popup")
		g.SetViewOnTop(vn)
		g.SetCurrentView(vn)
	} else {
		debug(g, "Action: Hide debug popup")
		g.SetViewOnBottom(vn)
		g.SetCurrentView("pods")
	}

	DEBUG_DISPLAYED = !DEBUG_DISPLAYED

	return nil
}

// View namespaces: Toggle display
func actionGlobalToggleViewNamespaces(g *gocui.Gui, v *gocui.View) error {
	vn := "namespaces"

	if !NAMESPACES_DISPLAYED {
		debug(g, "Action: Display namespaces popup")
		g.SetViewOnTop(vn)
		g.SetCurrentView(vn)
		changeStatusContext(g, "SE")
	} else {
		debug(g, "Action: Hide namespaces popup")
		g.SetViewOnBottom(vn)
		g.SetCurrentView("pods")
		changeStatusContext(g, "D")
	}

	NAMESPACES_DISPLAYED = !NAMESPACES_DISPLAYED

	return nil
}

// View pods: Up
func actionViewPodsUp(g *gocui.Gui, v *gocui.View) error {
	moveViewCursorUp(g, v, 2)
	debug(g, "Select up in pods view")
	return nil
}

// View pods: Down
func actionViewPodsDown(g *gocui.Gui, v *gocui.View) error {
	moveViewCursorDown(g, v, false)
	debug(g, "Select down in pods view")
	return nil
}

// View pods: Delete
//func actionViewPodsDelete(g *gocui.Gui, v *gocui.View) error {
//	p, err := getSelectedPod(g)
//	if err != nil {
//		return err
//	}
//
//	if err := deletePod(p); err != nil {
//		return err
//	}
//
//	debug(g, "Delete pod: "+p)
//	displayConfirmation(g, p+" pod deleted")
//
//	go viewPodsRefreshList(g)
//
//	return nil
//}

func actionViewPodsDelete(g *gocui.Gui, v *gocui.View) error {
	p, err := getSelectedPod(g)
	if err != nil {
		return err
	}
	deleteConfirmation(g, p)

	return nil
}

func actionDelPodConfirm(g *gocui.Gui, v *gocui.View) error {
	p, err := getSelectedPod(g)
	if err != nil {
		return err
	}
	if err := deletePod(p); err != nil {
		return err
	}
	debug(g, "Delete pod: "+p)
	go viewPodsRefreshList(g)
	g.DeleteView("delconfirmation")
	g.SetCurrentView("pods")
	return nil
}

func actionBackToPod(g *gocui.Gui, v *gocui.View) error {
	g.DeleteView("delconfirmation")
	g.SetCurrentView("pods")
	return nil
}


// View pods: Logs
func actionViewPodsLogs(g *gocui.Gui, v *gocui.View) error {
	LOG_MOD = "pod"
	err := showViewPodsLogs(g)

	changeStatusContext(g, "SL")

	return err
}

// View pod logs: Up
func actionViewPodsLogsUp(g *gocui.Gui, v *gocui.View) error {
	vLc, err := g.View("logs-containers")
	if err != nil {
		return err
	}
	moveViewCursorUp(g, vLc, 0)
	refreshPodsLogs(g)
	debug(g, "Select up in logs view")
	return nil
}

// View pod logs: Down
func actionViewPodsLogsDown(g *gocui.Gui, v *gocui.View) error {
	vLc, err := g.View("logs-containers")
	if err != nil {
		return err
	}
	moveViewCursorDown(g, vLc, false)
	refreshPodsLogs(g)
	debug(g, "Select down in logs view")
	return nil
}

// View logs: Hide
func actionViewPodsLogsHide(g *gocui.Gui, v *gocui.View) error {
	g.SetViewOnBottom("logs")
	g.SetViewOnBottom("logs-containers")
	g.SetCurrentView("pods")

	v.Clear()

	changeStatusContext(g, "D")

	debug(g, "Action: Hide view logs)")

	return nil
}

// View namespaces: Up
func actionViewNamespacesUp(g *gocui.Gui, v *gocui.View) error {
	moveViewCursorUp(g, v, 0)
	debug(g, "Select up in namespaces view")
	return nil
}

// View namespaces: Down
func actionViewNamespacesDown(g *gocui.Gui, v *gocui.View) error {
	moveViewCursorDown(g, v, false)
	debug(g, "Select down in namespaces view")
	return nil
}

// Namespace: Choose
func actionViewNamespacesSelect(g *gocui.Gui, v *gocui.View) error {
	line, err := getViewLine(g, v)
	debug(g, "Select namespace: "+line)
	NAMESPACE = line
	go viewPodsRefreshList(g)
	actionGlobalToggleViewNamespaces(g, v)
	displayConfirmation(g, line+" namespace selected")
	lMaxX,_ := v.Size()
	tv,_ := g.View("title")
	tv.Clear()
	fmt.Fprintln(tv, versionTitle(lMaxX))
	return err
}
