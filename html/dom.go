package html

func Tag(tag string, className string, children ...Renderer) *Node {
	node := &Node{}
	node.Tag = tag
	if className != "" {
		node.Attr("class", className)
	}
	node.Children = children
	return node
}

func Fragment(children ...Renderer) *Node {
	return Tag("", "", children...)
}

func Head(children ...Renderer) *Node {
	return Tag("head", "", children...)
}

func Link(href string) *Node {
	return Tag("link", "").
		Attr("rel", "stylesheet").
		Attr("href", href)
}

func Title(value string) *Node {
	return Tag("title", "", Text{value})
}

func Meta(attr, value string) *Node {
	return Tag("meta", "").Attr(attr, value)
}

func Div(className string, children ...Renderer) *Node { return Tag("div", className, children...) }

func P(text string) *Node  { return Tag("p", "", Text{text}) }
func H1(text string) *Node { return Tag("h1", "", Text{text}) }
func H2(text string) *Node { return Tag("h2", "", Text{text}) }
func H3(text string) *Node { return Tag("h3", "", Text{text}) }

func Form() *Node { return Tag("form", "") }
