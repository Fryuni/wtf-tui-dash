package github_global

import (
	"github.com/rivo/tview"
	"github.com/wtfutil/wtf/view"
	"strconv"
)

// Widget define wtf widget to register widget later
type Widget struct {
	view.TextWidget

	data *data

	settings *Settings
	Selected int
	maxItems int
	Items    []int
}

// NewWidget creates a new instance of the widget
func NewWidget(tviewApp *tview.Application, redrawChan chan bool, pages *tview.Pages, settings *Settings) *Widget {
	widget := Widget{
		TextWidget: view.NewTextWidget(tviewApp, redrawChan, pages, settings.Common),
		settings:   settings,
		data:       newData(settings),
	}

	widget.initializeKeyboardControls()

	widget.View.SetRegions(true)

	widget.Unselect()

	return &widget
}

/* -------------------- Exported Functions -------------------- */

// SetItemCount sets the amount of PRs RRs and other PRs throughout the widgets display creation
func (widget *Widget) SetItemCount(items int) {
	widget.maxItems = items
}

// GetItemCount returns the amount of PRs RRs and other PRs calculated so far as an int
func (widget *Widget) GetItemCount() int {
	return widget.maxItems
}

// GetSelected returns the index of the currently highlighted item as an int
func (widget *Widget) GetSelected() int {
	if widget.Selected < 0 {
		return 0
	}
	return widget.Selected
}

// Next cycles the currently highlighted text down
func (widget *Widget) Next() {
	widget.Selected++
	if widget.Selected >= widget.maxItems {
		widget.Selected = 0
	}
	widget.View.Highlight(strconv.Itoa(widget.Selected))
	widget.View.ScrollToHighlight()
}

// Prev cycles the currently highlighted text up
func (widget *Widget) Prev() {
	widget.Selected--
	if widget.Selected < 0 {
		widget.Selected = widget.maxItems - 1
	}
	widget.View.Highlight(strconv.Itoa(widget.Selected))
	widget.View.ScrollToHighlight()
}

// Unselect stops highlighting the text and jumps the scroll position to the top
func (widget *Widget) Unselect() {
	widget.Selected = -1
	widget.View.Highlight()
	widget.View.ScrollToBeginning()
}

// Refresh reloads the github data via the Github API and reruns the display
func (widget *Widget) Refresh() {
	widget.data.Refresh()
	widget.display()
}

func (widget *Widget) openPr() {
	currentSelection := widget.View.GetHighlights()
	if widget.Selected >= 0 && len(widget.Items) > 0 && currentSelection[0] != "" {
		//url := (*widget.currentGithubRepo().RemoteRepo.HTMLURL + "/pull/" + strconv.Itoa(widget.Items[widget.Selected]))
		//utils.OpenFile(url)
	}
}

func (widget *Widget) openPulls() {
	widget.data.OpenPulls()
}

func (widget *Widget) openIssues() {
	widget.data.OpenIssues()
}
