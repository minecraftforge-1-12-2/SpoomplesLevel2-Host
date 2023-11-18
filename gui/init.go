package gui

import (
	_ "embed"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
	"image/color"
)

var (
	//go:embed icon.png
	Icon    []byte
	Enabled = false
	App     fyne.App
	Window  fyne.Window

	OutBox   *widget.Entry
	InputBox *widget.Entry
)

func InitGui() {
	Enabled = true
	App = app.New()
	Window = App.NewWindow("Spoomples Level 2 Host")
	Window.SetFixedSize(true)
	Window.Resize(fyne.NewSize(600, 405))
	ico := fyne.NewStaticResource("icon.png", Icon)
	Window.SetIcon(ico)

	OutBox = widget.NewMultiLineEntry()
	InputBox = widget.NewEntry()

	OutBox.Disable()
	OutBox.SetMinRowsVisible(15)
	OutBox.TextStyle.Monospace = true
	OutBox.Wrapping = fyne.TextWrapWord

	InputBox.SetPlaceHolder("NO")
	InputBox.Disable()

	img := canvas.NewImageFromResource(ico)
	img.SetMinSize(fyne.NewSize(50, 50))
	c := container.New(
		layout.NewVBoxLayout(),
		container.New(
			layout.NewHBoxLayout(),
			layout.NewSpacer(),
			img,
			MakeGuiText(" Spoomples Level 2 Host", 30, color.White),
			layout.NewSpacer(),
		),
		layout.NewSpacer(),
		OutBox,
		InputBox,
	)

	Window.SetContent(c)
}

func StartGui() {
	OutBox.SetText("")
	Window.ShowAndRun()
}

func MakeGuiText(text string, size float32, color color.Color) *canvas.Text {
	t := canvas.NewText(text, color)
	t.TextSize = size
	t.TextStyle.Bold = true
	return t
}
