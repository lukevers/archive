package main

var MB *Mailbox

func main() {
	MB = NewMailbox(&Mailbox{
		Dirs: []string{"./emails"},
	})

	route()
}
