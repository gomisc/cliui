package editor

import (
	"bytes"
	"encoding/json"

	"git.eth4.dev/golibs/errors"
	"github.com/gdamore/tcell/v2"
	"github.com/pgavlin/femto"
	"github.com/pgavlin/femto/runtime"
	"github.com/rivo/tview"
	"gopkg.in/yaml.v2"
)

const ErrUnsupportedObjectFormat = errors.Const("unsupported object format")

func EditObject(obj any, fileType string) error {
	buf := &bytes.Buffer{}

	switch fileType {
	case "yaml":
		encoder := yaml.NewEncoder(buf)
		if err := encoder.Encode(obj); err != nil {
			return errors.Wrap(err, "encode data to yaml")
		}
	case "json":
		encoder := json.NewEncoder(buf)
		encoder.SetIndent("", "  ")

		if err := encoder.Encode(&obj); err != nil {
			return errors.Wrap(err, "encode data to json")
		}
	default:
		return ErrUnsupportedObjectFormat
	}

	data, err := Edit(buf.Bytes(), fileType)
	if err != nil {
		return errors.Wrap(err, "edit object data")
	}

	buf = bytes.NewBuffer(data)

	switch fileType {
	case "yaml":
		decoder := yaml.NewDecoder(buf)
		if err = decoder.Decode(obj); err != nil {
			return errors.Wrapf(err, "unmarshal yaml")
		}
	case "json":
		decoder := json.NewDecoder(buf)
		if err = decoder.Decode(&obj); err != nil {
			return errors.Wrapf(err, "unmarshal json")
		}
	}

	return nil
}

func Edit(content []byte, fileType string) (result []byte, err error) {
	var colorscheme femto.Colorscheme

	if monokai := runtime.Files.FindFile(femto.RTColorscheme, "monokai"); monokai != nil {
		var data []byte

		if data, err = monokai.Data(); err == nil {
			colorscheme = femto.ParseColorscheme(string(data))
		}
	}

	app := tview.NewApplication()
	buffer := makeBuffer(content, fileType)
	root := femto.NewView(buffer)
	root.SetRuntimeFiles(runtime.Files)
	root.SetColorscheme(colorscheme)
	root.SetBorderColor(tcell.ColorDarkOrange)
	root.SetBorder(true)
	root.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Key() {
		case tcell.KeyCtrlW:
			result = []byte(buffer.String())
			app.Stop()
			return nil
		case tcell.KeyCtrlQ:
			result = content
			app.Stop()
			return nil
		}
		return event
	})
	app.SetRoot(root, true)

	if err = app.Run(); err != nil {
		return nil, errors.Wrap(err, "run editor")
	}

	return result, nil
}

func makeBuffer(content []byte, filetype string) *femto.Buffer {
	buff := femto.NewBufferFromString(string(content), "")
	buff.Settings["filetype"] = filetype
	buff.Settings["keepautoindent"] = true
	buff.Settings["statusline"] = false
	buff.Settings["softwrap"] = true
	buff.Settings["scrollbar"] = true
	buff.Settings["smartpaste"] = true

	return buff
}
