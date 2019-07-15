package pkg

import "github.com/jroimartin/gocui"

// Configure globale keys
var keys []Key = []Key{
	Key{"", gocui.KeyCtrlC, actionGlobalQuit},
	Key{"", gocui.KeyCtrlD, actionGlobalToggleViewDebug},
	Key{"pods", gocui.KeyCtrlN, actionGlobalToggleViewNamespaces},
	Key{"pods", gocui.KeyArrowUp, actionViewPodsUp},
	Key{"pods", gocui.KeyArrowDown, actionViewPodsDown},
	Key{"pods", 'd', actionViewPodsDelete},
	Key{"pods", 'l', actionViewPodsLogs},
	Key{"logs", 'l', actionViewPodsLogsHide},
	Key{"logs", gocui.KeyArrowUp, actionViewPodsLogsUp},
	Key{"logs", gocui.KeyArrowDown, actionViewPodsLogsDown},
	Key{"namespaces", gocui.KeyArrowUp, actionViewNamespacesUp},
	Key{"namespaces", gocui.KeyArrowDown, actionViewNamespacesDown},
	Key{"namespaces", gocui.KeyEnter, actionViewNamespacesSelect},
	Key{"delconfirmation", gocui.KeyEnter, actionDelPodConfirm},
	Key{"delconfirmation", 'q', actionBackToPod},
}

type Key struct {
	viewname string
	key      interface{}
	handler  func(*gocui.Gui, *gocui.View) error
}

// Define UI key bindings
func UiKey(g *gocui.Gui) error {
	for _, key := range keys {
		if err := g.SetKeybinding(key.viewname, key.key, gocui.ModNone, key.handler); err != nil {
			return err
		}
	}

	return nil
}
