# Week 4 — Structs and Methods (Go)

This folder contains runnable solutions split by exercise so separate `main()` functions do not conflict.

## Run the exercises
From this directory:

```sh
go run ./01_structs
go run ./02_composition
go run ./03_receivers
go run ./04_bank_account
go test -bench=. -benchmem ./05_benchmarks
go run ./06_shapes
go run ./07_stringers
go test ./...
```

The compiler-error examples in Task 3 are documented below rather than left uncommented in a `.go` file, because a program containing the expected errors cannot build.

## Task 1: Structs

A struct is a typed collection of named fields. Go does not have classes or class inheritance. Structs group data; methods attach behaviour; composition builds larger data types.

Four ways to create a `Student`:

```go
var s1 Student
s2 := Student{Name: "Aruzhan", ID: 1, GPA: 3.8}
s3 := Student{"Dias", 2, 3.2, []string{"Go", "Databases"}}
p := &Student{Name: "Timur"}
```

Prefer named fields because positional literals depend on field order. A zero-value struct is usable: strings become `""`, numeric fields become `0`, and a slice becomes `nil`. Exported fields start with uppercase letters; lowercase fields are package-private.

`%v` prints field values, `%+v` prints field names and values, `%#v` prints a Go-syntax representation (including the package-qualified type and explicit field values).

### Struct copying and slices

Assigning `t2 := t1` copies the struct fields. An array field is copied element by element. A slice field is a slice header (pointer, length, capacity), so the copy refers to the same backing array. Consequently, changing `t2.Name` and `t2.Scores[0]` does not change those fields in `t1`, while `t2.Members[0]` can change an element visible through `t1.Members`.

Passing `Student` by value copies it. `append` returns a slice header with the updated length. If the function receives only a copy of that header, the caller's length is not updated. Passing `*Student` lets the function assign the new slice header back into the original struct.

### Struct equality

Structs can use `==` only when all fields are comparable. A `Point{1,2}` equals another `Point{1,2}`. Two pointers to different Point objects are not equal, even if the pointed-to fields match. Dereferencing them and comparing `*c == *d` compares the struct values and is true. `Team` with a `[]string` field cannot be compared with `==`, because slices are not comparable.

`struct{}` has no fields and occupies zero bytes. `map[string]struct{}` is a common set representation: the key represents membership, and the value needs no storage for meaningful data.

## Task 2: Composition and embedding

Go uses composition rather than class inheritance. `Person` can be a named field (`Address Address`) or an embedded field (`Person`). Embedding promotes fields and methods, allowing `student.Name` or `student.Contact()`, while the actual field remains accessible as `student.Person`.

`Student` and `Teacher` each embed `Person`. `Teacher.Contact` shadows the promoted `Person.Contact`; call `teacher.Person.Contact()` to use the embedded method directly.

Embedding is not inheritance or dynamic method overriding. In the Animal/Dog example, `dog.Sound()` returns `Woof!`, but `dog.Describe()` invokes `Animal.Describe` with an `Animal` receiver. Inside that method, `a.Sound()` means `Animal.Sound()`, so it prints `Rex says ...`, not `Rex says Woof!`.

If two embedded fields both promote an `ID`, `record.ID` is ambiguous and does not compile. Explicit paths such as `record.Device.ID` and `record.User.ID` work. Embedding a pointer also promotes methods, but calling through a nil embedded pointer can panic.

## Task 3: Compile or not?

Given:

```go
type Counter struct{ n int }
func (c Counter) Get() int { return c.n }
func (c *Counter) Inc() { c.n++ }
func NewCounter() Counter { return Counter{} }
```

1. `c.Inc()` — compiles: local variable `c` is addressable; Go rewrites it as `(&c).Inc()`.
2. `(&c).Inc()` — compiles: explicitly passes a pointer.
3. `Counter{}.Get()` — compiles: `Get` has a value receiver, which can receive the temporary value.
4. `Counter{}.Inc()` — does not compile: the composite literal is not addressable, so Go cannot implicitly take its address for a pointer-receiver call.
5. `NewCounter().Inc()` — does not compile: a function result is not addressable for this implicit `&` operation.
6. `m["a"].Get()` where `m` is `map[string]Counter` — compiles: the value-receiver method uses a copy of the map element.
7. `m["a"].Inc()` — does not compile: map elements are not addressable, so Go cannot take the element's address.
8. `s[0].Inc()` where `s` is `[]Counter` — compiles: slice elements are addressable.
9. `pm["a"].Inc()` where `pm` is `map[string]*Counter` — compiles: the map element is a pointer value, so the method can be called through that pointer.

### Pointer and value receivers

A value receiver (`func (c Counter)`) works on a copy. A pointer receiver (`func (c *Counter)`) can modify the original object and avoids copying a large receiver. Go can automatically take the address of an addressable value for a pointer-method call, and dereference a pointer to call a value method. This convenience does not make map elements or function return values addressable.

The `range` form `for _, acc := range accounts` copies each element. Calling a pointer method on `acc` changes that copy. Use `for i := range accounts { accounts[i].Deposit(...) }` to mutate the actual slice elements.

Choose value receivers for small, cheap-to-copy types that behave like values and whose methods do not need to mutate the receiver. Choose pointer receivers when mutation is needed, copying is expensive, or consistency requires it. The assignment recommends using pointer receivers consistently once a type has methods that need them.

## Task 4: Bank account

`Owner` is exported; `balance` is unexported to protect the account state from direct access outside the package. `Deposit` changes the balance, `Withdraw` returns an error for invalid/insufficient funds, and `Balance` reads the balance. This solution uses pointer receivers consistently because the account has mutation methods. The interest loop uses indices so the original accounts are updated.

## Task 5: Value vs pointer receiver benchmarks

Run:

```sh
go test -bench=. -benchmem ./05_benchmarks
```

The `Report` array `[1024]float64` occupies 8192 bytes, excluding other fields/padding. `//go:noinline` asks the compiler not to inline that function, making it easier to observe receiver-call costs in this demonstration. A value receiver copies the receiver's value; a pointer receiver passes a pointer instead. `B/op` measures heap bytes allocated per operation, not all bytes copied. A stack copy can cost time without allocating heap memory, so `B/op` may be 0 for both. `ns/op` is the column that may reveal the time cost. Results depend on machine, Go version, and compiler; interpret your local run rather than copying a preset number.

## Task 6: Shapes and interfaces

An interface specifies method signatures. A type implements an interface implicitly if its method set contains all required methods. `Shape` requires `Area() float64` and `Perimeter() float64`. `Rectangle` uses value receivers, so both `Rectangle` and `*Rectangle` implement Shape. `Circle` uses pointer receivers, so only `*Circle` implements Shape; use `&Circle{Radius: 1}`, not `Circle{Radius: 1}`. `Scale` uses a pointer receiver to update the original rectangle, but is not required by the Shape interface.

## Task 7: Tour of Go — Stringers

`fmt.Stringer` is the interface `{ String() string }`. The `IPAddr` type is `[4]byte`; the `String` method formats it as a dotted IPv4 address, e.g. `1.2.3.4`. When formatted, `fmt` calls `String()` for values whose method set implements `fmt.Stringer`.

The additional Book demonstration uses `func (b *Book) String() string`. Because this method has a pointer receiver, `*Book` implements Stringer but `Book` does not. Therefore `fmt.Println(book)` uses the default struct format, while `fmt.Println(&book)` uses the custom string.

## Method sets / interface rule

- The method set of `T` contains methods declared with receiver `T`.
- The method set of `*T` contains methods declared with receiver `T` and `*T`.

A pointer variable can call a value-receiver method because Go dereferences it. A local addressable value can call a pointer-receiver method because Go can take its address. But an interface stores a value and cannot generally make a non-pointer value addressable. Thus `Dog{}` does not implement an interface if the required method exists only on `*Dog`; `&Dog{}` does.

For embedded methods, embedding `T` promotes `T` methods to the outer type and promotes pointer-receiver methods to the pointer-to-outer type; embedding `*T` promotes both value and pointer methods to both outer method sets. See the assignment's method-set table for the precise combinations.

### Other defence questions from the PDF

- **What is a struct vs a class?** A struct groups fields. Go has no classes, constructors, `extends`, or built-in inheritance; methods and composition provide related capabilities.
- **What happens on struct assignment?** Every field is copied. If a field is a slice, only its header is copied; its backing array may remain shared.
- **Why can a value call a pointer method but a map element cannot?** The value must be addressable for Go to insert `&`; variables and slice elements are addressable, map elements are not.
- **What is a method?** A function with a receiver parameter attached to a named type.
- **What is the bonus method-value output?** `12 40`: `f := r.Area` captures a copy of `r` when the method value is created; later `r.Width = 10` does not change the receiver saved inside `f`.
- **What does `fmt.Stringer` do?** It enables a type to define its custom string representation with `String() string`.

### Expected answers for the struct-formatting example

- `fmt.Printf("%v\\n", s1)` prints the field values without field names (the zero-value fields are empty/zero; a nil slice displays as `[]` in the default formatting).
- `fmt.Printf("%+v\\n", s2)` includes field names, for example `{Name:Aruzhan ID:1 GPA:3.8 Courses:[]}`.
- `fmt.Printf("%#v\\n", s3)` prints a Go-syntax representation, like `main.Student{Name:"Dias", ID:2, GPA:3.2, Courses:[]string{"Go", "Databases"}}` (the package qualifier can differ).
- `s1.Courses == nil` is `true`, because the zero value of a slice is `nil`.

### Method sets with embedded fields (question A–D)

Using `type Speaker interface { Speak() string }`, a `Base` with `func (Base) Hello() string` and `func (*Base) Speak() string`:

- A: `var _ Speaker = &ByValue{}` — compiles. `*ByValue` gets the promoted pointer-receiver method.
- B: `var _ Speaker = ByValue{}` — does not compile. `ByValue`'s method set does not include the promoted pointer-receiver `Speak` method.
- C: `var _ Speaker = ByPointer{Base: &Base{}}` — compiles. Embedding `*Base` promotes both `Base` and `*Base` methods to `ByPointer`.
- D: `var _ Speaker = &ByPointer{Base: &Base{}}` — compiles for the same reason.

Summary table:

| Outer type embeds | Method set of outer value `S` | Method set of pointer `*S` |
|---|---|---|
| `T` | methods with receiver `T` | methods with receiver `T` and `*T` |
| `*T` | methods with receiver `T` and `*T` | methods with receiver `T` and `*T` |

### Methods on named types

Methods can be defined on a named type declared in your own package, including a named basic type such as `type Celsius float64`. They cannot be added directly to built-in types such as `int`, or to types owned by another package. This is why `func (i int) Double() int` is invalid.

### Value receiver vs pointer receiver: short answer for defence

- Use `func (x T) Method()` for a small value-like type when a copy is fine and the method doesn't need to mutate the original.
- Use `func (x *T) Method()` when changing the original is required, the value is large to copy, or the type should use pointer receivers consistently.
- For the BankAccount task, all methods use `*BankAccount`: Deposit/Withdraw mutate it, and keeping Balance on a pointer receiver is consistent even though Balance only reads.
