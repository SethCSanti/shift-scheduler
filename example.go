// models.go
type Note struct {
    Title   string
    Content string
}

type Author struct {
    Username     string
    PasswordHash string
}

type AppState struct {
    Authors map[string]*Author
    Notes   map[string]*Note
}

type TemplateData struct {
    User    *Author
    Error   string
    Success string
}


// handlers.go
var app AppState

// GET handler — just renders a template with data
func noteView(w http.ResponseWriter, r *http.Request) {
    // 1. get the session user
    username, err := getSessionUser(r)
    if err != nil {
        http.Redirect(w, r, "/login", http.StatusSeeOther)
        return
    }

    // 2. build the data to pass to the template
    data := TemplateData{
        User: app.Authors[username],
    }

    // 3. parse the templates (base + page specific)
    tmpl, err := template.ParseFiles("templates/base.html", "templates/note.html")
    if err != nil {
        http.Error(w, "template error", http.StatusInternalServerError)
        return
    }

    // 4. execute the template with data
    tmpl.ExecuteTemplate(w, "base", data)
}

// POST handler — processes form data then redirects
func createNoteView(w http.ResponseWriter, r *http.Request) {
    // 1. get the session user
    username, err := getSessionUser(r)
    if err != nil {
        http.Redirect(w, r, "/login", http.StatusSeeOther)
        return
    }

    // 2. read form values
    title := r.FormValue("title")
    content := r.FormValue("content")

    // 3. validate — return early if something is wrong
    if title == "" {
        data := TemplateData{
            User:  app.Authors[username],
            Error: "Title cannot be empty",
        }
        tmpl, _ := template.ParseFiles("templates/base.html", "templates/create_note.html")
        tmpl.ExecuteTemplate(w, "base", data)
        return
    }

    // 4. do the actual work
    app.Notes[title] = &Note{
        Title:   title,
        Content: content,
    }

    // 5. redirect on success
    http.Redirect(w, r, "/notes", http.StatusSeeOther)
}

// POST handler — involves password hashing
func registerAuthorView(w http.ResponseWriter, r *http.Request) {
    // 1. read form values
    username := r.FormValue("username")
    password := r.FormValue("password")
    confirmPassword := r.FormValue("confirmPassword")

    // 2. validate — check passwords match
    if password != confirmPassword {
        data := TemplateData{Error: "Passwords do not match"}
        tmpl, _ := template.ParseFiles("templates/base.html", "templates/register.html")
        tmpl.ExecuteTemplate(w, "base", data)
        return
    }

    // 3. check if user already exists
    if _, exists := app.Authors[username]; exists {
        data := TemplateData{Error: "Username already taken"}
        tmpl, _ := template.ParseFiles("templates/base.html", "templates/register.html")
        tmpl.ExecuteTemplate(w, "base", data)
        return
    }

    // 4. hash the password
    hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
    if err != nil {
        http.Error(w, "server error", http.StatusInternalServerError)
        return
    }

    // 5. assign role based on whether any authors exist
    role := "viewer"
    if len(app.Authors) == 0 {
        role = "admin"
    }
    _ = role // use role however you need

    // 6. store the new author
    app.Authors[username] = &Author{
        Username:     username,
        PasswordHash: string(hash),
    }

    // 7. redirect to login
    http.Redirect(w, r, "/login", http.StatusSeeOther)
}

// POST handler — involves bcrypt comparison + setting a cookie
func loginAuthorView(w http.ResponseWriter, r *http.Request) {
    // 1. read form values
    username := r.FormValue("username")
    password := r.FormValue("password")

    // 2. look up the user
    author, exists := app.Authors[username]
    if !exists {
        data := TemplateData{Error: "Invalid username or password"}
        tmpl, _ := template.ParseFiles("templates/base.html", "templates/login.html")
        tmpl.ExecuteTemplate(w, "base", data)
        return
    }

    // 3. compare password against hash
    err := bcrypt.CompareHashAndPassword([]byte(author.PasswordHash), []byte(password))
    if err != nil {
        data := TemplateData{Error: "Invalid username or password"}
        tmpl, _ := template.ParseFiles("templates/base.html", "templates/login.html")
        tmpl.ExecuteTemplate(w, "base", data)
        return
    }

    // 4. set the session cookie
    setSessionUser(w, username)

    // 5. redirect to home
    http.Redirect(w, r, "/", http.StatusSeeOther)
}

// POST handler — clears session and redirects
func logoutAuthorView(w http.ResponseWriter, r *http.Request) {
    clearSessionUser(w)
    http.Redirect(w, r, "/login", http.StatusSeeOther)
}