package html

type Text struct{ Content string }

func (t Text) Render(w *Writer) { w.Text(t.Content) }

type Node struct {
	Tag      string // "" - fragment
	Attrs    []Attr
	Children []Renderer
}

type Attr struct {
	Name  string
	Value string
}

func (n *Node) Attr(name, value string) *Node {
	n.Attrs = append(n.Attrs, Attr{name, value})
	return n
}

func (n *Node) Class(value string) *Node {
	n.Attr("class", value)
	return n
}
func (n *Node) Text(text ...string) *Node {
	for _, t := range text {
		n.Children = append(n.Children, Text{t})
	}
	return n
}

func (n *Node) Child(rs ...Renderer) *Node {
	n.Children = append(n.Children, rs...)
	return n
}

func (n *Node) Render(w *Writer) {
	n.RenderOpen(w)
	n.RenderClose(w)
}

func (n *Node) RenderOpen(w *Writer) {
	if n.Tag != "" {
		w.Open(n.Tag)

		for _, attr := range n.Attrs {
			w.Attr(attr.Name, attr.Value)
		}
		w.CloseAttributes()
	}

	w.Render(n.Children...)
}

func (n *Node) RenderClose(w *Writer) {
	if n.Tag != "" {
		w.Close(n.Tag)
	}
}
