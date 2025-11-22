package runner

import (
	"m3g4p0p/agents/runner/history"
	"m3g4p0p/agents/toolset"
	"m3g4p0p/agents/util"
)

func WithProcessor(processor history.Processor) util.Option[Runner] {
	return func(r *Runner) {
		r.processor = processor
	}
}

func WithToolset(toolset toolset.Toolset) util.Option[Runner] {
	return func(r *Runner) {
		r.toolset = toolset
	}
}
