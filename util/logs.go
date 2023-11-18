package util

import (
	"TogetherForever/gui"
	"fmt"
	"strconv"
	"strings"
)

func WriteLine(text string, format ...string) {
	if len(format) > 0 {
		for i, f := range format {
			text = strings.Replace(text, "{"+strconv.Itoa(i)+"}", f, -1)
		}
	}
	if gui.Enabled {
		gui.OutBox.SetText(gui.OutBox.Text + "\n" + text)
		gui.OutBox.CursorRow = len(strings.Split(gui.OutBox.Text, "\n"))
	} else {
		fmt.Println(text)
	}
}
