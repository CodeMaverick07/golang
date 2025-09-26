package main

import (
	"fmt"
)

func main() {
	// variables
	var name string = "hemant"
	fmt.Printf("hello %s\n", name)
	age := 21
	fmt.Printf("my age is %d\n", age)

	var city string = "mumbai"

	fmt.Println(city)

	var pincode, area string = "413512", "Laxmi colony"
	fmt.Println(pincode, area)

	var (
		isEmployed bool   = true
		salary     int    = 10000
		role       string = "developer"
	)

	fmt.Printf("is employed %t, salary %d, role %s \n", isEmployed, salary, role)

	var defaultInt int
	var defaultFloat float32
	var defaultString string
	var defaultBoolean bool

	fmt.Printf("int: %d \nfloat: %f \nstring: %s \nstring: %t \n", defaultInt, defaultFloat, defaultString, defaultBoolean)

	const typedAge int = 21
	const untypedAge = 21

	fmt.Println(typedAge == untypedAge)
	const (
		jan = iota + 1
		feb
		mar
		apr
	)
	fmt.Println(jan, feb, mar, apr)
	addition := add(4, 5)
	fmt.Println(addition)
	addition, multiplication := addAndMultiply(5, 6)
	fmt.Println(addition, multiplication)

	// if else
	if addition > 10 {
		fmt.Println("addition is greater than 10")
	} else if addition < 12 {
		fmt.Println("good")
	} else {
		fmt.Println("fuck off")
	}

	// switches
	day := "monday"
	switch day {
	case "monday":
		fmt.Print("monday")
	case "tuesday":
		fmt.Print("tuesday")
	default:
		fmt.Println("its weekend")
	}
	//for loops
	for i := 0; i < 5; i++ {
		fmt.Println(i)
	}

	//while loop
	// golang dose not have while loop it uses for keyword for while loop
	counter := 0
	for counter < 5 {
		fmt.Println(counter)
		counter++
	}

	//arrays and slices
	numbers := [5]int{1, 2, 3, 4, 5}
	len := len(numbers)
	fmt.Printf("this is our array %v and length is %d\n", numbers, len)

	//slices
	//allNumbers := numbers[:]
	//firstThree := numbers[0:3]

	initSlice := []int{11, 12, 13}
	initSlice = append(initSlice, 14)

	//itrate on arrays or slice

	for index, value := range initSlice {
		fmt.Printf("index %d, value %d\n", index, value)
	}

	//maps
	capitals := map[string]string{
		"india": "Delhi",
		"usa":   "dc",
	}

	capital, exits := capitals["japan"]

	if exits {
		fmt.Println(capital)
	} else {
		fmt.Println()
	}

	delete(capitals, "india")
	person := Person{name: "hemant", age: 12}

	fmt.Printf("printing the struct %+v\n", person)

	employee := struct {
		salary int
		role   string
	}{salary: 1000, role: "developer"}
	fmt.Printf("printing the struct %+v\n", employee)

}

type Person struct {
	name string
	age  int
}

func add(a int, b int) int {
	return a + b
}

func addAndMultiply(a, b int) (int, int) {
	return a + b, a * b
}
