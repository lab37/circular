package main

import (
	"001/data"
	"fmt"
	"net/http"
	"strings"
	"unicode"
)

// GET /login
// Show the login page
func login(writer http.ResponseWriter, request *http.Request) {
	t := parseTemplateFiles("login.layout", "public.navbar", "login")
	t.Execute(writer, nil)
}

// GET /signup
// Show the signup page
func signup(writer http.ResponseWriter, request *http.Request) {
	sess, err := session(writer, request)
	if err != nil {
		http.Redirect(writer, request, "/login", 302)

	} else {
		user, err := sess.GetUser()
		if err != nil {
			http.Redirect(writer, request, "/login", 302)
		}
		if user.Email == "36ee@163.com" {
			generateHTML(writer, nil, "login.layout", "public.navbar", "signup")
		} else {
			fmt.Fprintf(writer, "只有周京成能添加帐号，请联系周京成")
		}
	}
}

// POST /signup
// Create the user account
func signupAccount(writer http.ResponseWriter, request *http.Request) {
	sess, err := session(writer, request)
	if err != nil {
		http.Redirect(writer, request, "/login", 302)

	} else {
		userNow, err := sess.GetUser()
		if err != nil {
			http.Redirect(writer, request, "/login", 302)

		}
		if userNow.Email == "36ee@163.com" {

			err := request.ParseForm()
			if err != nil {
				danger(err, "Cannot parse form")
			}
			user := data.User{
				Name:     request.PostFormValue("name"),
				Email:    request.PostFormValue("email"),
				Password: request.PostFormValue("password"),
			}
			if err := user.Create(); err != nil {
				danger(err, "Cannot create user")
			}
			http.Redirect(writer, request, "/login", 302)
		} else {
			fmt.Fprintf(writer, "注册失败！只有周京成能添加帐号，请联系周京成")
		}
	}
}

// POST /authenticate
// Authenticate the user given the email and password
func authenticate(writer http.ResponseWriter, request *http.Request) {
	err := request.ParseForm()
	user, err := data.GetUserByEmail(request.PostFormValue("email"))
	if err != nil {
		danger(err, "Cannot find user")
	}
	if user.Password == data.Encrypt(request.PostFormValue("password")) {
		session, err := user.CreateSession()
		if err != nil {
			danger(err, "Cannot create session")
		}
		cookie := http.Cookie{
			Name:     "_cookie",
			Value:    session.Uuid,
			HttpOnly: true,
		}
		http.SetCookie(writer, &cookie)
		http.Redirect(writer, request, "/jxc", 302)
	} else {
		http.Redirect(writer, request, "/login", 302)
	}

}

// GET /logout
// Logs the user out
func logout(writer http.ResponseWriter, request *http.Request) {
	cookie, err := request.Cookie("_cookie")
	if err != http.ErrNoCookie {
		warning(err, "Failed to get cookie")
		session := data.Session{Uuid: cookie.Value}
		session.DeleteByUUID()
	}
	http.Redirect(writer, request, "/", 302)
}

func updateAccount(writer http.ResponseWriter, request *http.Request) {

	sess, err := session(writer, request)
	if err != nil {
		http.Redirect(writer, request, "/login", 302)
	} else {
		userNew := data.User{}
		userNew.Name = strings.TrimFunc(request.PostFormValue("name"), unicode.IsSpace)
		userNew.Email = strings.TrimFunc(request.PostFormValue("email"), unicode.IsSpace)
		userOld, _ := sess.GetUser()
		userNew.Id = userOld.Id
		userNew.Update()
		http.Redirect(writer, request, "/login", 302)
	}
}

func updatePassword(writer http.ResponseWriter, request *http.Request) {

	sess, err := session(writer, request)
	if err != nil {
		http.Redirect(writer, request, "/login", 302)
	} else {
		userNew := data.User{}
		userNew.Password = strings.TrimFunc(request.PostFormValue("password"), unicode.IsSpace)
		userOld, _ := sess.GetUser()
		userNew.Id = userOld.Id
		userNew.UpdatePassword()
		http.Redirect(writer, request, "/login", 302)
	}
}
