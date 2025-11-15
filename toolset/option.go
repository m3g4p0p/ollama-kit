package toolset

type ToolOption func(t *Toolset)

func WithTool[T any](name, description string, handler Handler[T]) ToolOption {
	return func(t *Toolset) {
		if err := AddTool(t, name, description, handler); err != nil {
			panic(err)
		}
	}
}

func WithSimpleTool[T any](name, description string, handler SimpleHandler[T]) ToolOption {
	return func(t *Toolset) {
		if err := AddSimpleTool(t, name, description, handler); err != nil {
			panic(err)
		}
	}
}

func WithStructuredOutput[T any](name, description string) ToolOption {
	return func(t *Toolset) {
		if err := AddStructuredOutput[T](t, name, description); err != nil {
			panic(err)
		}
	}
}
