package hud

import (
	"time"

	"github.com/gdamore/tcell"

	"github.com/astro-walker/tilt/internal/hud/view"
	"github.com/astro-walker/tilt/internal/rty"
	"github.com/astro-walker/tilt/pkg/model"
)

type buildStatus struct {
	edits      []string
	duration   time.Duration
	status     string
	deployTime time.Time
	reason     model.BuildReason
	muted      bool
}

func (bs buildStatus) defaultTextColor() tcell.Color {
	if bs.muted {
		return cLightText
	}
	return tcell.ColorDefault
}

func makeBuildStatus(res view.Resource, triggerMode model.TriggerMode) buildStatus {
	status := "Pending"
	duration := time.Duration(0)
	edits := []string{}
	deployTime := time.Time{}
	reason := model.BuildReason(0)

	if res.IsTiltfile {
		return buildStatus{
			status: "N/A",
		}
	}

	if !res.CurrentBuild.Empty() {
		status = "In prog."
		duration = time.Since(res.CurrentBuild.StartTime)
		edits = res.CurrentBuild.Edits
		reason = res.CurrentBuild.Reason
	} else if !res.PendingBuildSince.IsZero() {
		status = "Pending"
		if triggerMode.AutoOnChange() {
			duration = time.Since(res.PendingBuildSince)
		}
		edits = res.PendingBuildEdits
		reason = res.PendingBuildReason
	} else if !res.LastBuild().FinishTime.IsZero() {
		lastBuild := res.LastBuild()
		if lastBuild.Error != nil {
			status = "Error"
		} else {
			status = "OK"
		}
		duration = lastBuild.Duration()
		edits = lastBuild.Edits
		deployTime = res.LastDeployTime
		reason = lastBuild.Reason
	}

	return buildStatus{
		status:     status,
		duration:   duration,
		edits:      edits,
		deployTime: deployTime,
		reason:     reason,
	}
}

func buildStatusCell(bs buildStatus) rty.Component {
	textColor := bs.defaultTextColor()
	showingDuration := bs.duration != 0
	lhsWidth := BuildStatusCellMinWidth
	if !showingDuration {
		lhsWidth += BuildDurCellMinWidth
	}
	lhs := rty.NewMinLengthLayout(lhsWidth, rty.DirHor).
		Add(rty.ColoredString(bs.status, textColor))
	if !showingDuration {
		return lhs
	}

	sb := rty.NewStringBuilder()
	sb.Fg(cLightText).Text(" (")
	sb.Fg(textColor).Text(formatBuildDuration(bs.duration))
	sb.Fg(cLightText).Text(")")
	rhs := rty.NewMinLengthLayout(BuildDurCellMinWidth, rty.DirHor).
		SetAlign(rty.AlignEnd).
		Add(sb.Build())

	return rty.NewConcatLayout(rty.DirHor).
		Add(lhs).
		Add(rhs)
}
