package habitica

import (
	"fmt"
	"github.com/wtfutil/wtf/modules/habitica/api"
	"strings"
)

type (
	renderComponent interface {
		display(state *renderState)
	}
	regionHandler interface {
		toggle()
		open()
	}
	renderState struct {
		widget         *Widget
		builder        *strings.Builder
		regionHandlers map[string]regionHandler
	}
)

func (s *renderState) addSection(title string, components []renderComponent) {

}

type Task api.Task

func (t *Task) display(state *renderState) {
	fmt.Fprintf(state.builder, ` - [%q][ ][""] %s`, t.Id, t.Text)

	for _, listItem := range t.Checklist {
		marker := ' '
		if listItem.Completed {
			marker = 'X'
		}
		fmt.Fprintf(state.builder, `   - [%q][%s][""] %s`, listItem.Id, marker, listItem.Text)
	}
}
