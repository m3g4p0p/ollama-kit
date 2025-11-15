package agentlib

import "m3g4p0p/agents/toolset"

// Stolen from https://openai.github.io/openai-agents-python/ref/extensions/handoff_prompt/#agents.extensions.handoff_prompt.RECOMMENDED_PROMPT_PREFIX
const HandoffPrefix = "# System context\nYou are part of a multi-agent system, designed to make agent coordination and execution easy. Agents uses two primary abstraction: **Agents** and **Handoffs**. An agent encompasses instructions and tools and can hand off a conversation to another agent when appropriate. Handoffs are achieved by calling a handoff function, generally named `transfer_to_<agent_name>`. Transfers between agents are handled seamlessly in the background; do not mention or draw attention to these transfers in your conversation with the user.\n"

type HandoffParams struct {
	Prompt string
}

func WithHandoff(agent Agent, name, description string) Option {
	return WithOptions(
		WithInstructions(HandoffPrefix),
		WithTool("transfer_to_"+name, description, func(args HandoffParams) (toolset.ToolResult, error) {
			return toolset.ToolResult{
				Content: args.Prompt,
				Args:    agent,
				Final:   true,
			}, nil
		}),
	)
}
