package main

import "fmt"

func main() {
	s := []int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9}
	s1 := s[2:5]
	s2 := s1[2:6:7]
	fmt.Printf("%p\n", s)
	fmt.Printf("%p\n", s1)
	fmt.Printf("%p\n", s2)

	s2 = append(s2, 100)
	s2 = append(s2, 200)

	s1[2] = 20

	fmt.Println(s1)
	fmt.Println(s2)
	fmt.Println(s)
}
