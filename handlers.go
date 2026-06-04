package main

import (
	// native Go packages
	"html/template"
	"net/http"

	// internal packages
	// 3rd party packages
	"golang.org/x/crypto/bcrypt"
)

func home(w http.ResponseWriter, r *http.Request) {
	username, err := getSessionUser(r)
	data := TemplateData{}
	if err == nil {
		data.User = app.Users[username]
	} else {
		data.User = nil
	}
	tmpl, err := template.ParseFiles("templates/base.html", "templates/home.html")
	if err != nil {
		http.Error(w, "template error", http.StatusInternalServerError)
		return
	}
	tmpl.ExecuteTemplate(w, "base", data)
}

// phase 1 handlers

func registerView(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("templates/base.html", "templates/register.html")
	if err != nil {
		http.Error(w, "template error", http.StatusInternalServerError)
		return
	}
	tmpl.ExecuteTemplate(w, "base", TemplateData{})
}

func registerSubmitView(w http.ResponseWriter, r *http.Request) {
    // Parse form data
	username := r.FormValue("username")
	password := r.FormValue("password")
	email := r.FormValue("email")
	confirmPassword := r.FormValue("confirmPassword")

	// Basic validation
	if password != confirmPassword {
		data := TemplateData{
			ErrorMessage: "Passwords do not match",
		}
		tmpl, err := template.ParseFiles("templates/base.html", "templates/register.html")
		if err != nil {
			http.Error(w, "template error", http.StatusInternalServerError)
			return
		}
		tmpl.ExecuteTemplate(w, "base", data)
		return
	}

	if _, exists := app.Users[username]; exists {
		data := TemplateData{
			ErrorMessage: "Username already exists",
		}
		tmpl, err := template.ParseFiles("templates/base.html", "templates/register.html")
		if err != nil {
			http.Error(w, "template error", http.StatusInternalServerError)
			return
		}
		tmpl.ExecuteTemplate(w, "base", data)
		return
	}

	// Hash the password
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		http.Error(w, "error hashing password", http.StatusInternalServerError)
		return
	}

    // Create new user and store in app state
    role := "employee" // Default role for new users
    if len(app.Users) == 0 {
        role = "admin" // First user becomes admin
    }
    app.Users[username] = &User{
        UserName: username,
        PasswordHash: string(hash),
        Role: role,
        Email: email,
    }

    // Redirect to login page after successful registration
    saveData()
    http.Redirect(w, r, "/login", http.StatusSeeOther)
}

func loginView(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("templates/base.html", "templates/login.html")
	if err != nil {
		http.Error(w, "template error", http.StatusInternalServerError)
		return
	}
	tmpl.ExecuteTemplate(w, "base", TemplateData{})
}

func loginSubmitView(w http.ResponseWriter, r *http.Request) {
    // Parse form data
    username := r.FormValue("username")
    password := r.FormValue("password")

    // Validate user credentials
    user, exists := app.Users[username]
    if !exists || bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)) != nil {
        data := TemplateData{
            ErrorMessage: "Invalid username or password",
        }
        tmpl, err := template.ParseFiles("templates/base.html", "templates/login.html")
        if err != nil {
            http.Error(w, "template error", http.StatusInternalServerError)
            return
        }
        tmpl.ExecuteTemplate(w, "base", data)
        return
    }

    // Create a new session for the user
    setSessionUser(w, username)
    // Redirect to the home page
    saveData()
    http.Redirect(w, r, "/", http.StatusSeeOther)
}

func logoutView(w http.ResponseWriter, r *http.Request) {
    // Clear the session
    clearSessionUser(w)
    // Redirect to the home page
    http.Redirect(w, r, "/", http.StatusSeeOther)
}

// phase 2 handlers

func scheduleView(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Schedule view"))
}

func approvalView(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Approval view"))
}

func submitView(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Submit view"))
    saveData()
}

func decideView(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Decide view"))
    saveData()
}

// phase 3 handlers

func profileView(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Profile view"))
}

func pictureView(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Profile picture update view"))
    saveData()
}
