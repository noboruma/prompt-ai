package main

import (
	"fmt"
	"strings"

	"github.com/gdamore/tcell/v2"
	"github.com/noboruma/prompt-ai/ais"
	"github.com/noboruma/prompt-ai/clipboards"
	"github.com/noboruma/prompt-ai/prompt"
	str "github.com/noboruma/prompt-ai/strings"
	"github.com/rivo/tview"
)

func prepareChatViewSection() *tview.TextView {
	chatView := tview.NewTextView().
		SetTextAlign(tview.AlignLeft).
		SetText("...").
		SetDynamicColors(true).
		SetRegions(true).
		SetWordWrap(true)
	chatView.SetBorder(true).SetTitle("Chat")

	return chatView
}

func prepareInputSection(app *tview.Application, chatView *tview.TextView, copy_clipboards *clipboards.CopyClipboards) *tview.Flex {
	chat_history := ""

	errorText := tview.NewTextView()
	errorText.SetTextColor(tcell.ColorRed)
	statusText := tview.NewTextView().
		SetScrollable(false).
		SetToggleHighlights(false).
		SetDynamicColors(true).
		SetRegions(true).
		SetWordWrap(true)
	errorLog := func(msg string) {
		errorText.SetText(msg)
	}

	updateQuota := func() {
		usage, err := ais.GetOpenAIQuotaUsage()
		if err != nil {
			errorLog(err.Error())
		}
		fmt.Fprintf(statusText, "[yellow]Usage: $%v", usage.Used)
	}
	updateQuota()

	inputTextArea := tview.NewInputField().
		SetLabel(prompt.DefaultPrompt).
		SetPlaceholder("E.g. why is 42 the answer")

	inputTextArea.SetDoneFunc(func(key tcell.Key) {
		if key == tcell.KeyEnter {
			inputTextArea.SetDisabled(true)
			spinctl := prompt.StartSpinner(app, inputTextArea)
			go func() {
				defer inputTextArea.SetDisabled(false)
				defer func() { spinctl <- struct{}{} }()
				responses, err := ais.SendPrompt(inputTextArea.GetText(), 100)
				if err != nil {
					chatView.SetText(err.Error())
					chatView.SetBackgroundColor(tcell.ColorRed)
					return
				}
				tmp := strings.Join(responses, "\n")
				chat_history += "[yellow]>>" + inputTextArea.GetText() + "\n[white]"

				sections := str.ExtractMarkdownSections(tmp)
				for i := range sections {
					if sections[i].Markdown {
						id := copy_clipboards.Append(sections[i].Content)
						chat_history += "[green]" + fmt.Sprintf("[F%d]", id+1) + "\n---[blue]"
						chat_history += sections[i].Content + "---[white]"
					} else {
						chat_history += sections[i].Content
					}
				}
				chat_history += "\n"
				chatView.SetText(chat_history)
				inputTextArea.SetText("")
			}()
		}
	})

	inputField := tview.NewFlex().SetDirection(tview.FlexRow).
		AddItem(inputTextArea, 0, 1, true).
		AddItem(tview.NewTextView().
			SetDynamicColors(true).
			SetText("[green]F1-F8[white]: copy to system clipboard"), 0, 1, false).
		AddItem(tview.NewTextView().
			SetDynamicColors(true).
			SetText("[green]TAB[white]: switch panes"), 0, 1, false).
		AddItem(statusText, 0, 1, true)

	inputField.SetBorder(true).SetTitle("Input")
	return inputField
}

func main() {

	copy_clipboards := clipboards.NewCopyClipboards()
	app := tview.NewApplication()
	chatView := prepareChatViewSection()

	inputField := prepareInputSection(app, chatView, &copy_clipboards)
	inputField.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		key := event.Key()
		if key >= tcell.KeyF1 && key <= tcell.KeyF8 {
			copy_clipboards.Fetch(int(key - tcell.KeyF1))
		} else if key == tcell.KeyTab {
			app.SetFocus(chatView)
		}
		return event
	})

	chatView.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		key := event.Key()
		if key == tcell.KeyTab {
			app.SetFocus(inputField)
		}
		return event
	})

	mainLayout := tview.NewFlex().
		AddItem(tview.NewFlex().SetDirection(tview.FlexRow).
			AddItem(chatView, 0, 2, false).
			AddItem(inputField, 6, 1, true), 0, 1, true)

	if err := app.SetRoot(mainLayout, true).EnableMouse(true).Run(); err != nil {
		panic(err)
	}
}
