package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
)

type Pessoa struct {
	ID int 'json:"id"'
	Nome string 'json:"nome"'
}

var pessoa []Pessoa

func main(){
	http.HandleFunc("/pessoa", handlePessoas)
	http.HandleFunc("/pessoa/", handlePessoaByID)

	log.Fatal(http.ListenAndServe(":8080", nil))
}

