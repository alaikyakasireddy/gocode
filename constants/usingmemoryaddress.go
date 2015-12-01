package main

import "fmt"

const gramtokilogram int =1/1000


func main() {
	var grams int
	fmt.Println("enter the grams")
	fmt.Scan(&grams)
	kilogram := grams * gramtokilogram
	fmt.Println(grams, "grams is", kilogram , "kilogram ")
}