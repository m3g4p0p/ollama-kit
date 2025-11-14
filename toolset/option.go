package toolset

type ToolOption func(t *Toolset)

func WithTool[T any](name, description string, handler func(T) string) ToolOption {
	return func(t *Toolset) {
		if err := AddTool(t, name, description, handler); err != nil {
			panic(err)
		}
	}
}
