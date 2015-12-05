package main
import "fmt"

func main () {
	a := fmt.Fprint("this", 1, "is a demonstration", true, "of Fprint")
	fmt.Print(a)
}