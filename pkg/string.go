package pkg

import (
	"fmt"
	"github.com/jroimartin/gocui"
	"strings"
	"time"
)

// Repeat string
func timesString(str string, n int) (out string) {
	for i := 0; i < n; i++ {
		out += str
	}
	return
}

// Text pad center
func textPadCenter(s string, l int) string {
	pc := " "
	p := timesString(pc, (l/2)-len(s)/2)
	o := p + s + p
	if (len(s) < l) && (l < len(o)) {
		o = o[0:l]
	}
	return o
}

// StringFormatBoth fg and bg colors
// Thanks https://github.com/mephux/komanda-cli/blob/master/komanda/color/color.go
func stringFormatBoth(fg, bg int, str string, args []string) string {
	return fmt.Sprintf("\x1b[48;5;%dm\x1b[38;5;%d;%sm%s\x1b[0m", bg, fg, strings.Join(args, ";"), str)
}

// Frame text with colors
func frameText(text string) string {
	return stringFormatBoth(15, 0, text, []string{"1"})
}

// Useful to debug Pody (display with CTRL+D)
func debug(g *gocui.Gui, output interface{}) {
	v, err := g.View("debug")
	if err == nil {
		t := time.Now()
		tf := t.Format("2006-01-02 15:04:05")
		output = tf + " > " + output.(string)
		fmt.Fprintln(v, output)
	}
}