package habitica

import (
	"fmt"
	"github.com/rivo/tview"
	"github.com/wtfutil/wtf/modules/habitica/api"
	"github.com/wtfutil/wtf/view"
	"strings"
)

// Widget is the container for your module's data
type Widget struct {
	view.TextWidget
	settings *Settings

	client *api.Client

	todos []*api.Task

	Err error
}

// NewWidget creates and returns an instance of Widget
func NewWidget(tviewApp *tview.Application, redrawChan chan bool, pages *tview.Pages, settings *Settings) *Widget {
	widget := Widget{
		TextWidget: view.NewTextWidget(tviewApp, redrawChan, pages, settings.common),
		settings:   settings,

		client: api.NewClient(settings.userId, settings.apiToken),
	}

	return &widget
}

/* -------------------- Exported Functions -------------------- */

// Refresh updates the onscreen contents of the widget
func (widget *Widget) Refresh() {
	widget.refresh()
	widget.TextWidget.Redraw(widget.display)
}

/* -------------------- Unexported Functions -------------------- */

func (widget *Widget) content() string {
	builder := new(strings.Builder)

	if widget.Err != nil {
		builder.WriteString("\n [red]Error[white]\n")
		_, _ = fmt.Fprintf(builder, "%+v", widget.Err)

		return builder.String()
	}

	builder.WriteString("\n [red]TODOs[white]\n")

	for _, todo := range widget.todos {
		builder.WriteString(todo.Text)
		builder.WriteByte('\n')
	}

	return builder.String()
}

func (widget *Widget) display() (string, string, bool) {
	return widget.CommonSettings().Title, widget.content(), widget.Err != nil
}

func (widget *Widget) refresh() {
	widget.todos, widget.Err = widget.client.ListTodos()
	if widget.Err != nil {
		return
	}
}
