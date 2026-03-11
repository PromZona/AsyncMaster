package runtime

type MockContext struct {
}

func (m *MockContext) ChatID() int64 {
	return 0
}

func (m *MockContext) FirstName() string {
	return "-"
}

func (m *MockContext) Text() string {
	return "24"
}
func (m *MockContext) Callback() string {
	return "32"
}
func (m *MockContext) Send(msg string) error {
	return nil
}

type MockRuntime struct {
	middlewares []Middleware
}

func (mr *MockRuntime) HandleText(h Handler) {
	mr.apply(h)
}

func (mr *MockRuntime) HandleCallback(h Handler) {
	mr.apply(h)

}

func (mr *MockRuntime) Start() error {
	return nil
}

func (t *MockRuntime) Use(m Middleware) {
	t.middlewares = append(t.middlewares, m)
}

func (t *MockRuntime) apply(h Handler) Handler {
	for i := len(t.middlewares) - 1; i >= 0; i-- {
		h = t.middlewares[i](h)
	}
	return h
}
