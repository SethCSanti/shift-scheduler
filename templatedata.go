package main

type TemplateData struct {
	User *User // Reference to the user for personalized greetings
	ErrorMessage string // Error message to display on the page, if any
	SuccessMessage string // Success message to display on the page, if any
}