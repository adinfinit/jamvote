package html

type Renderer interface {
	Render(w *Writer)
}

type RenderFunc func(w *Writer)

func (fn RenderFunc) Render(w *Writer) { fn(w) }

type RendererComposite interface {
	Renderer
	RenderOpen(w *Writer)
	RenderClose(w *Writer)
}

func (w *Writer) Render(rs ...Renderer) {
	for _, r := range rs {
		r.Render(w)
	}
}
