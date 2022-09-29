package editor

import (
	"git.corout.in/golibs/errors"
	"github.com/gdamore/tcell/v2"
	"github.com/pgavlin/femto"
	"github.com/pgavlin/femto/runtime"
	"github.com/rivo/tview"
)

func Edit(text, fileType string) (result string, err error) {
	var colorscheme femto.Colorscheme

	if monokai := runtime.Files.FindFile(femto.RTColorscheme, "monokai"); monokai != nil {
		var data []byte

		if data, err = monokai.Data(); err == nil {
			colorscheme = femto.ParseColorscheme(string(data))
		}
	}

	app := tview.NewApplication()
	buffer := makeBufferFromString(text, fileType)
	root := femto.NewView(buffer)
	root.SetRuntimeFiles(runtime.Files)
	root.SetColorscheme(colorscheme)
	root.SetBorderColor(tcell.ColorDarkOrange)
	root.SetBorder(true)
	root.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Key() {
		case tcell.KeyCtrlW:
			result = buffer.String()
			app.Stop()
			return nil
		case tcell.KeyCtrlQ:
			app.Stop()
			return nil
		}
		return event
	})
	app.SetRoot(root, true)

	if err = app.Run(); err != nil {
		return "", errors.Wrap(err, "run editor")
	}

	return result, nil
}

func makeBufferFromString(content, filetype string) *femto.Buffer {
	buff := femto.NewBufferFromString(content, "")
	buff.Settings["filetype"] = filetype
	buff.Settings["keepautoindent"] = true
	buff.Settings["statusline"] = false
	buff.Settings["softwrap"] = true
	buff.Settings["scrollbar"] = true
	buff.Settings["smartpaste"] = true

	return buff
}
