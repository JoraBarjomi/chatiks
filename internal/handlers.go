package internal

import (
	"html/template"
	"log"
	"net/http"
)

type Base struct {
	// User       *utils.User
	IsLoggedIn bool
	Title      string
	AvatarSrc  string
}

// func WsHelperHandler(writer http.ResponseWriter, request *http.Request) {
// 	// session, err := gothic.Store.Get(request, "app_session")
// 	// utils.Check(err)
// 	// userID, ok := session.Values["user_id"]

// 	// var user *utils.User
// 	isLoggedIn := false

// 	// if ok {
// 	// 	user = utils.GetUserFromDB(userID.(int))
// 	// 	if user != nil {
// 	// 		isLoggedIn = true
// 	// 		log.Printf("User %s in the chat!", user.Email)
// 	// 	}
// 	// }

// 	data := Base{
// 		// User:       user,
// 		IsLoggedIn: isLoggedIn,
// 		Title:      "Chat - Site",
// 	}

// 	tmpl := template.Must(template.ParseFiles("ui/templates/base.html", "ui/templates/chat.html"))
// 	err := tmpl.ExecuteTemplate(writer, "base.html", data)
// 	if err != nil {

// 	}
// 	// utils.Check(err)
// }

func ServeChat(w http.ResponseWriter, r *http.Request) {
	log.Println(r.URL)
	// if r.URL.Path != "/" {
	// 	http.Error(w, "Not found", http.StatusNotFound)
	// 	return
	// }
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	data := Base{
		// IsLoggedIn: true,
		Title: "Chatiks",
	}

	tmpl := template.Must(template.ParseFiles("ui/templates/chatBase.html", "ui/templates/chat.html"))
	err := tmpl.ExecuteTemplate(w, "chatBase.html", data)
	if err != nil {
		log.Println(err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
	}
}
