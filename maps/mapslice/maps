package main
import "fmt"

func main () {
	myMap := map[string]string {
		"Alaikya":"Anusha",
		"allu":"sruthi",
		"vineela":"bhargav",
	}

	myMap["Navmit"] = "Danjeer"
	myMap["Navmit"] = "Something Else"

	delete(myMap, "Navmit")

	for key, val := range myMap {
		fmt.Println (key, " - ", val)
	}

	fmt.Println(len(myMap))

	if val, ok := myMap["allu"]; ok {
		fmt.Println(val)
	}
}
