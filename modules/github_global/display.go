package github_global

import (
	"fmt"

	ghb "github.com/google/go-github/v32/github"
)

func (widget *Widget) display() {
	widget.TextWidget.Redraw(widget.content)
}

func (widget *Widget) content() (string, string, bool) {
	// Choses the correct place to scroll to when changing sources
	if len(widget.View.GetHighlights()) > 0 {
		widget.View.ScrollToHighlight()
	} else {
		widget.View.ScrollToBeginning()
	}

	// initial maxItems count
	widget.Items = make([]int, 0)
	widget.SetItemCount(0)

	title := fmt.Sprintf("%s", widget.CommonSettings().Title)

	var str string
	if widget.settings.showOpenReviewRequests {
		str += fmt.Sprintf("\n [%s]Open Review Requests[white]\n", widget.settings.Colors.Subheading)
		str += widget.displayMyReviewRequests()
	}
	if widget.settings.showMyPullRequests {
		str += fmt.Sprintf("\n [%s]My Pull Requests[white]\n", widget.settings.Colors.Subheading)
		str += widget.displayMyPullRequests()
	}
	if widget.settings.showAssignedIssues {
		str += fmt.Sprintf("\n [%s]My Assigned Issues[white]\n", widget.settings.Colors.Subheading)
		str += widget.displayAssignedIssues()
	}
	for _, customQuery := range widget.settings.customQueries {
		str += fmt.Sprintf("\n [%s]%s[white]\n", widget.settings.Colors.Subheading, customQuery.title)
		str += widget.displayCustomQuery(customQuery.filter, customQuery.perPage)
	}

	return title, str, false
}

func (widget *Widget) displayMyPullRequests() string {
	prs := widget.data.myPullRequests()

	prLength := len(prs)

	if prLength == 0 {
		return " [grey]none[white]\n"
	}

	maxItems := widget.GetItemCount()

	str := ""
	for idx, pr := range prs {
		str += fmt.Sprintf(` %s[green]["%d"]%4d[""][white] %s`, widget.mergeString(pr), maxItems+idx, *pr.Number, *pr.Title)
		str += "\n"
		widget.Items = append(widget.Items, *pr.Number)
	}

	widget.SetItemCount(maxItems + prLength)

	return str
}

func (widget *Widget) displayMyReviewRequests() string {
	prs := widget.data.myReviewRequests()

	if len(prs) == 0 {
		return " [grey]none[white]\n"
	}

	maxItems := widget.GetItemCount()

	str := ""
	for idx, pr := range prs {
		str += fmt.Sprintf(` %s[green]["%d"]%4d[""][white] %s`, widget.mergeString(pr), maxItems+idx, *pr.Number, *pr.Title)
		str += "\n"
		widget.Items = append(widget.Items, *pr.Number)
	}

	widget.SetItemCount(maxItems + len(prs))

	return str
}

func (widget *Widget) displayAssignedIssues() string {
	issues := widget.data.MyAssignedIssues

	prLength := len(issues)

	if prLength == 0 {
		return " [grey]none[white]\n"
	}

	maxItems := widget.GetItemCount()

	str := ""
	for idx, pr := range issues {
		str += fmt.Sprintf(` [green]["%d"]%4d[""][white] %s`, maxItems+idx, *pr.Number, *pr.Title)
		str += "\n"
		widget.Items = append(widget.Items, *pr.Number)
	}

	widget.SetItemCount(maxItems + prLength)

	return str
}

func (widget *Widget) displayCustomQuery(filter string, perPage int) string {
	res := widget.data.customIssueQuery(filter, perPage)

	if res == nil {
		return " [grey]Invalid Query[white]\n"
	}

	issuesLength := len(res.Issues)

	if issuesLength == 0 {
		return " [grey]none[white]\n"
	}

	maxItems := widget.GetItemCount()

	str := ""
	for idx, issue := range res.Issues {
		str += fmt.Sprintf(` [green]["%d"]%4d[""][white] %s`, maxItems+idx, *issue.Number, *issue.Title)
		str += "\n"
		widget.Items = append(widget.Items, *issue.Number)
	}

	widget.SetItemCount(maxItems + issuesLength)

	return str
}

var mergeIcons = map[string]string{
	"dirty":    "[red]\u0021[white] ",
	"clean":    "[green]\u2713[white] ",
	"unstable": "[red]\u2717[white] ",
	"blocked":  "[red]\u2717[white] ",
}

func (widget *Widget) mergeString(pr *ghb.PullRequest) string {
	if !widget.settings.enableStatus {
		return ""
	}
	if str, ok := mergeIcons[pr.GetMergeableState()]; ok {
		return str
	}
	return "? "
}
