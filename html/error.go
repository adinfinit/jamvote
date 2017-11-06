package html

type Error struct{ Err error }

func (err Error) Render(w *Writer) {
	w.Open("div")
	w.Attr("class", "h-error")
	w.Text("Error: " + err.Err.Error())
	w.Close("div")
}
