package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os/exec"
	"time"
)

type Promotion struct {
	Name     string
	Course   string
	Year     int
	Students []Student
}

type Student struct {
	Name     string
	LastName string
	Age      int
	Gender   string
}

var Mentors = Promotion{
	Name:   "Mentor'ac",
	Course: "Informatique",
	Year:   5,
	Students: []Student{{Name: "Cyril", LastName: "RODRIGUES", Age: 22, Gender: "homme"},
		{Name: "Kheir-Eddine", LastName: "MEDERREG", Age: 22, Gender: "homme"},
		{Name: "Alan", LastName: "PHILIPIERT", Age: 28, Gender: "homme"}},
}

type PromoData struct {
	Title            string
	PromoTitle       string
	PromoNb          int
	StudentsTitle    string
	CurrentPromotion Promotion
}

type ChangeData struct {
	Title       string
	Message     string
	ChangeCount int
}

type User struct {
	LastName string
	Name     string
	Birthday string
	Gender   string
}

var lastUser User

type UserInitData struct {
	Title     string
	FormTitle string
}

type UserDisplayData struct {
	Title     string
	MainTitle string
	UserInfo  User
}

func main() {
	tmpl, err := template.ParseGlob("../templates/*.html")
	if err != nil {
		log.Fatal(err)
	}

	http.HandleFunc("/promo", func(w http.ResponseWriter, r *http.Request) {
		dataPage := PromoData{
			Title:            "Promotion",
			PromoTitle:       "Informations sur la promotion",
			PromoNb:          len(Mentors.Students),
			StudentsTitle:    "Liste des étudiants",
			CurrentPromotion: Mentors,
		}
		tmpl.ExecuteTemplate(w, "promo", dataPage)
	})

	var changeCount int

	http.HandleFunc("/change", func(w http.ResponseWriter, r *http.Request) {
		changeCount++
		var message string
		if changeCount%2 == 0 {
			message = "Le nombre de visite est pair !"
		} else {
			message = "Le nombre de visite est impair !"
		}
		dataPage := ChangeData{
			Title:       "Change",
			Message:     message,
			ChangeCount: changeCount,
		}
		tmpl.ExecuteTemplate(w, "change", dataPage)
	})

	http.HandleFunc("/user/init", func(w http.ResponseWriter, r *http.Request) {
		dataPage := UserInitData{
			Title:     "Inscription",
			FormTitle: "Renseignez vos données",
		}
		tmpl.ExecuteTemplate(w, "init", dataPage)
	})

	http.HandleFunc("/user/treatment", func(w http.ResponseWriter, r *http.Request) {
		var gender string
		switch r.FormValue("gender") {
		case "male":
			gender = "M"
		case "female":
			gender = "F"
		default:
			gender = "Inconnu"
		}
		lastUser = User{
			LastName: r.FormValue("lastName"),
			Name:     r.FormValue("name"),
			Birthday: r.FormValue("birthday"),
			Gender:   gender,
		}
		http.Redirect(w, r, "/user/display", 301)
	})

	http.HandleFunc("/user/display", func(w http.ResponseWriter, r *http.Request) {
		dataPage := UserDisplayData{
			Title:     "Données personnelles",
			MainTitle: "Vos données personnelles",
			UserInfo:  lastUser,
		}
		tmpl.ExecuteTemplate(w, "display", dataPage)
	})

	fs := http.FileServer(http.Dir("../css"))
	http.Handle("/css/", http.StripPrefix("/css/", fs))
	fs = http.FileServer(http.Dir("../assets"))
	http.Handle("/assets/", http.StripPrefix("/assets/", fs))
	port := "localhost:8080"
	url := "http://" + port + "/promo"
	cmd := exec.Command("explorer", url)
	go http.ListenAndServe(port, nil)
	fmt.Println("Server is running...")
	time.Sleep(time.Second * 5)
	cmd.Run()
	fmt.Println("If the navigator didn't open on its own, just go to " + url + " on your navigator.")
	isRunning := true
	for isRunning {
		fmt.Print("If you want to end the server, type 'stop' here : ")
		var command string
		fmt.Scanln(&command)
		if command == "stop" {
			isRunning = false
		}
	}
}
