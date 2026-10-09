package main

import "fmt"

type Student struct {
	Name    string
	ID      int
	GPA     float64
	Courses []string
}

// addCourse receives a copy of the Student struct.
func addCourse(s Student, course string) {
	s.Courses = append(s.Courses, course)
}

// addCoursePtr receives a pointer, so it can update the original slice header.
func addCoursePtr(s *Student, course string) {
	s.Courses = append(s.Courses, course)
}

type Team struct {
	Name    string
	Scores  [3]int
	Members []string
}

type Point struct{ X, Y int }

func main() {
	fmt.Println("===== TASK 1: STUDENT RECORDS =====")

	// Four ways to create Student values.
	var s1 Student
	s2 := Student{Name: "Aruzhan", ID: 1, GPA: 3.8, Courses: []string{"Go"}}
	s3 := Student{"Dias", 2, 3.2, []string{"Go", "Databases"}}
	p := &Student{Name: "Timur", ID: 3, GPA: 3.5, Courses: []string{"Algorithms"}}

	fmt.Printf("Zero value: %+v\n", s1)
	fmt.Printf("Named fields: %+v\n", s2)
	fmt.Printf("Positional: %+v\n", s3)
	fmt.Printf("Pointer: %+v\n", p)
	fmt.Println("s1.Courses is nil:", s1.Courses == nil)

	fmt.Println("\n===== VALUE PARAMETER VS POINTER PARAMETER =====")
	student := Student{Name: "Ali", ID: 4, GPA: 3.6, Courses: []string{"Go"}}
	addCourse(student, "Databases")
	fmt.Println("After addCourse(Student):", student.Courses)
	// The original slice length does not change: append updates only the copied
	// slice header in addCourse. The backing array may be shared, but its length
	// in the caller stays the same.

	addCoursePtr(&student, "Databases")
	fmt.Println("After addCoursePtr(*Student):", student.Courses)

	fmt.Println("\n===== STRUCT COPY: ARRAY VS SLICE FIELD =====")
	t1 := Team{Name: "Gophers", Scores: [3]int{1, 2, 3}, Members: []string{"Ali", "Bota"}}
	t2 := t1
	t2.Name = "Rustaceans"
	t2.Scores[0] = 100
	t2.Members[0] = "Nurlan"
	fmt.Println("t1:", t1.Name, t1.Scores, t1.Members)
	fmt.Println("t2:", t2.Name, t2.Scores, t2.Members)
	// t1 == t2 does not compile because Team contains a []string field.

	fmt.Println("\n===== STRUCT COMPARISON =====")
	a, b := Point{1, 2}, Point{1, 2}
	c, d := &Point{1, 2}, &Point{1, 2}
	fmt.Println("a == b:", a == b)
	fmt.Println("c == d:", c == d)
	fmt.Println("*c == *d:", *c == *d)

	fmt.Println("\n===== ANONYMOUS STRUCT AND SET =====")
	point := struct{ X, Y int }{1, 2}
	fmt.Println("Anonymous point:", point)

	seen := map[string]struct{}{}
	seen["go"] = struct{}{}
	_, ok := seen["go"]
	fmt.Println("Set contains go:", ok)
}
