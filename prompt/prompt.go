package prompt

import (
	"time"

	"github.com/rivo/tview"
)

const (
	DefaultPrompt = "Prompt>> "
	spinner_time  = 100 * time.Millisecond
)

var (
	frames = [...]string{"Prompt - ", "Prompt \\ ", "Prompt | ", "Prompt / "}
)

func StartSpinner(app *tview.Application, spinner *tview.InputField) chan struct{} {
	res := make(chan struct{})
	go func() {
		defer func() {
			app.QueueUpdateDraw(func() {
				spinner.SetLabel(DefaultPrompt)
			})
		}()
	main:
		for {
			for _, frame := range frames {
				// Set the text of the text view to the current frame.
				app.QueueUpdateDraw(func() {
					spinner.SetLabel(frame)
				})
				select {
				case <-res:
					break main
				default:
					time.Sleep(spinner_time)
				}
			}
		}
	}()
	return res
}
