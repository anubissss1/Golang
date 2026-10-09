package main

import "fmt"

type Person struct {
	Name  string
	Email string
}

func (p Person) Contact() string {
	return p.Name + " <" + p.Email + ">"
}

type Student struct {
	Person
	GPA float64
}

type Teacher struct {
	Person
	Department string
}

func (t Teacher) Contact() string {
	return t.Person.Contact() + " | Department: " + t.Department
}

type Animal struct{ Name string }

func (a Animal) Sound() string    { return "..." }
func (a Animal) Describe() string { return a.Name + " says " + a.Sound() }

type Dog struct{ Animal }

func (d Dog) Sound() string { return "Woof!" }

type Device struct{ ID int }
type User struct{ ID int }
type Record struct {
	Device
	User
}

type Logger struct{ prefix string }

func (l *Logger) Log(msg string) {
	fmt.Println(l.prefix + msg)
}

type Service struct {
	*Logger
	Name string
}

func main() {
	fmt.Println("===== TASK 2: UNIVERSITY COMPOSITION =====")
	student := Student{
		Person: Person{Name: "Ali", Email: "ali@example.com"},
		GPA:    3.7,
	}
	teacher := Teacher{
		Person:     Person{Name: "Aigerim", Email: "aigerim@example.com"},
		Department: "Computer Science",
	}

	fmt.Println("Student:", student.Contact())
	fmt.Println("Teacher:", teacher.Contact())
	fmt.Println("Teacher's embedded Person:", teacher.Person.Contact())

	fmt.Println("\n===== EMBEDDING IS NOT INHERITANCE =====")
	dog := Dog{Animal{Name: "Rex"}}
	fmt.Println("dog.Sound():", dog.Sound())
	fmt.Println("dog.Animal.Sound():", dog.Animal.Sound())
	fmt.Println("dog.Describe():", dog.Describe())
	// Describe belongs to Animal. Inside it, a.Sound() calls Animal.Sound;
	// it does not dynamically dispatch to Dog.Sound like a class override.

	fmt.Println("\n===== SAME-NAMED EMBEDDED FIELDS =====")
	record := Record{
		Device: Device{ID: 101},
		User:   User{ID: 202},
	}
	fmt.Println("Device ID:", record.Device.ID)
	fmt.Println("User ID:", record.User.ID)
	// record.ID would not compile: both embedded fields promote an ID, making
	// the selector ambiguous. Use the qualified paths above.

	fmt.Println("\n===== EMBEDDED POINTER =====")
	ok := Service{Logger: &Logger{prefix: "[svc] "}, Name: "users"}
	ok.Log("started")
	// Service{Name: "orders"}.Log(...) would panic because Logger is nil.
}
