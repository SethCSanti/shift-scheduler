package main

import (
	// native Go packages
	"html/template"
	"net/http"
	"encoding/json"
	"fmt"
	"time"
	"io"
	"os"
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
		UserName:     username,
		PasswordHash: string(hash),
		Role:         role,
		Email:        email,
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
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func logoutView(w http.ResponseWriter, r *http.Request) {
	// Clear the session
	clearSessionUser(w)
	// Redirect to the home page
	http.Redirect(w, r, "/login", http.StatusSeeOther)
}

// phase 2 handlers

func scheduleView(w http.ResponseWriter, r *http.Request) {
	username, err := getSessionUser(r)
	if err != nil {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}
	data := TemplateData{
		Schedule: app.Schedule[username],
		User:     app.Users[username],
	}
	data.Days  = []string{"Mon", "Tue", "Wed", "Thu", "Fri"}
    data.Hours = []string{"08:00", "09:00", "10:00", "11:00", "12:00", "13:00", "14:00", "15:00", "16:00", "17:00"}
	tmpl, err := template.ParseFiles("templates/base.html", "templates/schedule.html")
	if err != nil {
		http.Error(w, "template error", http.StatusInternalServerError)
		return
	}
	tmpl.ExecuteTemplate(w, "base", data)
}

func approvalView(w http.ResponseWriter, r *http.Request) {
    username, err := getSessionUser(r)
    if err != nil {
        http.Redirect(w, r, "/login", http.StatusSeeOther)
        return
    }
    if app.Users[username].Role != "admin" {
        http.Error(w, "forbidden", http.StatusForbidden)
        return
    }
    data := TemplateData{
        Schedules: app.Schedule,
        User:      app.Users[username],
    }
    tmpl, err := template.ParseFiles("templates/base.html", "templates/approval.html")
    if err != nil {
        http.Error(w, "template error", http.StatusInternalServerError)
        return
    }
    tmpl.ExecuteTemplate(w, "base", data)
}

func submitView(w http.ResponseWriter, r *http.Request) {
    // get session user
    username, err := getSessionUser(r)
    if err != nil {
        http.Error(w, "unauthorized", http.StatusUnauthorized)
        return
    }

    // decode JSON body into a slice of TimeBlocks
    var blocks []TimeBlock
    if err := json.NewDecoder(r.Body).Decode(&blocks); err != nil {
        http.Error(w, "invalid request body", http.StatusBadRequest)
        return
    }

    // calculate daily totals and validate each block
    dailyTotals := make(map[string]float64)
    for _, block := range blocks {
        start, err := time.Parse("15:04", block.StartTime)
        if err != nil {
            http.Error(w, "invalid start time", http.StatusBadRequest)
            return
        }
        end, err := time.Parse("15:04", block.EndTime)
        if err != nil {
            http.Error(w, "invalid end time", http.StatusBadRequest)
            return
        }

        hours := end.Sub(start).Hours()

        // minimum shift is 3 hours
        if hours < 3 {
            w.WriteHeader(http.StatusBadRequest)
            w.Write([]byte(`<div id="status-banner" class="banner error">Each shift must be at least 3 hours.</div>`))
            return
        }

        dailyTotals[block.Day] += hours
    }

    // maximum 9 hours per day
    for day, total := range dailyTotals {
        if total > 9 {
            w.WriteHeader(http.StatusBadRequest)
            fmt.Fprintf(w, `<div id="status-banner" class="banner error">%s exceeds the 9 hour daily maximum.</div>`, day)
            return
        }
    }

    // calculate weekly total
    weeklyTotal := 0.0
    for _, total := range dailyTotals {
        weeklyTotal += total
    }

    // weekly total must be between 20 and 40 hours
    if weeklyTotal < 20 {
        w.WriteHeader(http.StatusBadRequest)
        w.Write([]byte(`<div id="status-banner" class="banner error">Weekly total must be at least 20 hours.</div>`))
        return
    }
    if weeklyTotal > 40 {
        w.WriteHeader(http.StatusBadRequest)
        w.Write([]byte(`<div id="status-banner" class="banner error">Weekly total cannot exceed 40 hours.</div>`))
        return
    }

    // save the schedule
    app.Schedule[username] = &Schedule{
        Blocks:      blocks,
        DailyTotal:  dailyTotals,
        WeeklyTotal: weeklyTotal,
        Status:      "pending",
        UserName:    username,
    }

    saveData()

    // return success banner
    w.Write([]byte(`<div id="status-banner" class="banner success">Schedule submitted successfully and is pending approval.</div>`))
}

func decideView(w http.ResponseWriter, r *http.Request) {
	username, err := getSessionUser(r)
    if err != nil {
        http.Redirect(w, r, "/login", http.StatusSeeOther)
        return
    }
    if app.Users[username].Role != "admin" {
        http.Error(w, "forbidden", http.StatusForbidden)
        return
    }
	targetUsername := r.FormValue("targetUsername")
	decision := r.FormValue("decision")
	adminComment := r.FormValue("adminComment")

	schedule, exists := app.Schedule[targetUsername]
	if !exists {
		http.Error(w, "schedule not found", http.StatusNotFound)
		return
	}

	if decision == "approve" {
		schedule.Status = "approved"
	} else if decision == "reject" {
		schedule.Status = "rejected"
		schedule.AdminComment = adminComment
	}
	saveData()
	fmt.Fprintf(w, `<div id="status-banner" class="banner success">Schedule %sd for %s.</div>`, decision, targetUsername)
	http.Redirect(w, r, "/approval", http.StatusSeeOther)
}

// phase 3 handlers

func profileView(w http.ResponseWriter, r *http.Request) {
	username, err := getSessionUser(r)
    if err != nil {
        http.Redirect(w, r, "/login", http.StatusSeeOther)
        return
    }
	data := TemplateData{
		User: app.Users[username],
	}
	tmpl, err := template.ParseFiles("templates/base.html", "templates/profile.html")
	if err != nil {
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}
	err = tmpl.ExecuteTemplate(w, "base", data)
	if err != nil {
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}
}

func pictureView(w http.ResponseWriter, r *http.Request) {
	// get session user
	username, err := getSessionUser(r)
	if err != nil {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	// parse the multipart form (32MB max memory)
	err = r.ParseMultipartForm(32 << 20)
	if err != nil {
		http.Error(w, "error parsing form", http.StatusBadRequest)
		return
	}

	// get the uploaded file
	file, header, err := r.FormFile("picture")
	if err != nil {
		http.Error(w, "error reading file", http.StatusBadRequest)
		return
	}
	defer file.Close()

	// create destination file in static/uploads/
	dst, err := os.Create("static/uploads/" + header.Filename)
	if err != nil {
		http.Error(w, "error saving file", http.StatusInternalServerError)
		return
	}
	defer dst.Close()

	// copy uploaded file to destination
	_, err = io.Copy(dst, file)
	if err != nil {
		http.Error(w, "error saving file", http.StatusInternalServerError)
		return
	}

	// update user's profile picture
	app.Users[username].ProfilePicture = header.Filename

	saveData()

	w.Write([]byte(`<div id="status-banner" class="banner success">Profile picture updated successfully.</div>`))
}
