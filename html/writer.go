package html

import (
	"bytes"
	"fmt"
	"html"
)

type Writer struct {
	buf    bytes.Buffer
	stack  []string
	inattr bool
	invoid bool
}

func NewWriter() *Writer {
	return &Writer{}
}

func (w *Writer) String() string {
	if len(w.stack) != 0 {
		panic("writing not finished")
	}
	return w.buf.String()
}

func (w *Writer) Bytes() []byte {
	if len(w.stack) != 0 {
		panic("writing not finished")
	}
	return w.buf.Bytes()
}

func (w *Writer) writeByte(b byte) {
	w.buf.WriteByte(b)
}

func (w *Writer) writeString(s string) {
	w.buf.WriteString(s)
}

func (w *Writer) Open(tag string) {
	w.CloseAttributes()

	w.stack = append(w.stack, tag)
	w.invoid = voidElements[tag]
	w.inattr = true

	w.writeByte('<')
	w.writeString(tag)
}

func (w *Writer) CloseAttributes() {
	if !w.inattr {
		return
	}

	w.inattr = false
	w.writeByte('>')
}

func (w *Writer) Attr(name, value string) {
	if !w.inattr {
		panic("not in attributes section")
	}

	w.writeByte(' ')
	w.writeString(name)
	w.writeByte('=')
	w.writeByte('"')
	if name == "href" || name == "src" {
		w.writeString(NormalizeURL(value))
	} else {
		//w.writeString(EscapeAttribute(value))
		w.writeString(value)
	}
	w.writeByte('"')
}

func (w *Writer) Text(text string) {
	w.CloseAttributes()
	w.writeString(html.EscapeString(text))
}

func (w *Writer) UnsafeWrite(text string) {
	w.writeString(text)
}

func (w *Writer) UnsafeContent(text string) {
	w.CloseAttributes()
	w.writeString(text)
}

func (w *Writer) Close(tag string) {
	w.CloseAttributes()

	if len(w.stack) == 0 {
		panic("no unclosed tags")
	}

	var current string
	n := len(w.stack) - 1
	current, w.stack = w.stack[n], w.stack[:n]
	if current != tag {
		panic(fmt.Sprintf("closing tag %q expected %q", tag, current))
	}

	w.invoid = (len(w.stack) > 0) && voidElements[w.stack[len(w.stack)-1]]
	// void elements have only a single tag
	if voidElements[tag] {
		return
	}

	w.writeString("</")
	w.writeString(tag)
	w.writeByte('>')
}

// Section 12.1.2, "Elements", gives this list of void elements. Void elements
// are those that can't have any contents.
var voidElements = map[string]bool{
	"area":    true,
	"base":    true,
	"br":      true,
	"col":     true,
	"command": true,
	"embed":   true,
	"hr":      true,
	"img":     true,
	"input":   true,
	"keygen":  true,
	"link":    true,
	"meta":    true,
	"param":   true,
	"source":  true,
	"track":   true,
	"wbr":     true,
}
