package main

type createRequest struct {
	Date    string `json:"date"`
	Title   string `json:"title"`
	Comment string `json:"comment"`
	Repeat  string `json:"repeat"`
}

type createResponse struct {
	Id    int    `json:"id"`
	Error string `json:"error"`
}

type task struct {
	Id      int    `json:"id"`
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
