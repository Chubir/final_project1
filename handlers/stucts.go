package handlers

import repo "myApp/internal/repository"

type createRequest struct {
	Date    string `json:"date"`
	Title   string `json:"title"`
	Comment string `json:"comment"`
	Repeat  string `json:"repeat"`
}

type createResponseId struct {
	Id string `json:"id"`
}

type task struct {
	Id      string `json:"id"`
	Date    string `json:"date"`
	Title   string `json:"title"`
	Comment string `json:"comment"`
	Repeat  string `json:"repeat"`
}

type listResponse struct {
	Tasks []repo.Task `json:"tasks"`
}

type ResponseError struct {
	Error string `json:"error"`
}
