package main
import "fmt"

func main () {
	mySentence := "3 4"
	var one, two int
	fmt.Sscan(mySentence, &one, &two)
	fmt.Println(one, two)


}