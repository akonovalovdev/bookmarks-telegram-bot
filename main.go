package main

import ( 
	"fmt"
	"net/http"
	"html/template"
)

type User struct {
	Name string
	Age uint16
	Money int16
	AvgGrades float64
	Happiness float64
	Hobbies []string
}

func (u User) getAllInfo() string {
	return fmt.Sprintf("User name is: %s. He is %d and he has money equal: %d",
	u.Name, u.Age, u.Money)
}
func (u *User) setNewName(newName string) {
	u.Name = newName
}

func homePage(w http.ResponseWriter, r *http.Request) {
	bob := User{"Bob", 25, -50, 4.2, 0.8, []string{"Football", "Skate", "Dance"}}
	//bob := User{Name: "Bob", Age: 25, Money: -50, AvgGrades: 4.2, Happiness: 0.8}
	// fmt.Fprintf(w, `<h1>Main Text</h1>
	// <b>Main Text</b>`)
	tmpl, err := template.ParseFiles("templates/home_page.html")
	if err != nil{
		panic(err)
	}
	tmpl.Execute(w, bob)
}

func contactsPage(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Contacts page")
}

func handleRequest() {
	http.HandleFunc("/", homePage)
	http.HandleFunc("/contacts/", contactsPage)
	http.ListenAndServe(":8080", nil)

}

func main() {
	// var bob User =
	// bob := User{Name: "Bob", Age: 25, Money: -50, AvgGrades: 4.2, Happiness: 0.8}  User("Bob", 25, -50, 4.2, 0.8)

	handleRequest()

}
