package toolset

type ToolOption func(t *Toolset)

func WithTool[T any](name, description string, handler func(T) (string, error)) ToolOption {
	return func(t *Toolset) {
		if err := AddTool(t, name, description, handler); err != nil {
			panic(err)
		}
	}
}

func WithFinalTool[T any](name, description string) ToolOption {
	return func(t *Toolset) {
		if err := AddFinalTool[T](t, name, description); err != nil {
			panic(err)
		}
	}
}
