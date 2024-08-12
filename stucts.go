package main

type createRequest struct {
	Date    string `json:"date"`
	Title   string `json:"title"`
	Comment string `json:"comment"`
	Repeat  string `json:"repeat"`
}

type createResponseId struct {
	Id string `json:"id"`
}
type createResponseError struct {
	Error string `json:"error"`
}

type task struct {
	Id      string `json:"id"`
	Date    string `json:"date"`
	Title   string `json:"title"`
	Comment string `json:"comment"`
	Repeat  string `json:"repeat"`
}

type listResponse struct {
	Tasks []task `json:"tasks"`
}

type listErrResponse struct {
	Error string `json:"error"`
}
