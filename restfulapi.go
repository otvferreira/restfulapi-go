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

func main() {
	http.HandleFunc("/pessoa", handlePessoas)
	http.HandleFunc("/pessoa/", handlePessoaByID)

	log.Fatal(http.ListenAndServe(":8080", nil))
}

func handlePessoas(w http.ResponseWriter, r *http.Request) {
	
	switch r.Method {
		case http.MethodGet:
			obterPessoas(w, r)
		case http.MethodPost:
			criarPessoa(w, r)
		default:
			http.Error(w, "Método não existe", http.StatusMethodNotAllowed)
	}

}

func handlePessoaByID(w. http.ResponseWriter, r *http.Request) {

	switch r.Method {
		case http.MethodGet:
			obterPessoaPorID(w, r)
		case http.MethodDelete:
			deletarPessoaPorID(w, r)
		default:
			http.Error(w, "Método não existe", http.StatusMethodNotAllowed)
	}

}

func obterPessoas(w http.ResponseWriter, r *http.Request) {
	jsonBytes, err := json.Marshal(pessoas)
	if err != nil {
		http.Error(w, "Erro ao listar todas as pessoas.", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonBytes)
}


func criarPessoa(w http.ResponseWriter, r *http.Request) {
	var pessoa Pessoa
	err := json.NewDecoder(r.Body).Decode(&pessoa)
	if err != nil {
		http.Error(w, "JSON vazio ou incorreto", http.StatusBadRequest)
		return
	}

	pessoa.ID = len(pessoas) + 1
	pessoas = append(pessoas, pessoa)

	w.WriteHeader(http.StatusCreated)
}

func obterPessoaPorID(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Path[len("/pessoa/"):]
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "ID inexistente", http.StatusBadRequest)
		return
	}

	for _, pessoa := range pessoas {
		if pessoa.ID == id {
			jsonBytes, _ := json.Marshal(pessoa)
			w.Header().Set("Content-Type", "application/json")
			w.Write(jsonBytes)
			return
		}
	}

	w.WriteHeader(http.StatusNotFound)
}

func deletarPessoaPorID(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Path[len("/pessoa/"):]
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "ID inexistente", http.StatusBadRequest)
		return
	}

	var novasPessoas []Pessoa
	removido := false
	for _, pessoa := range pessoas {
		if pessoa.ID != id {
			novasPessoas = append(novasPessoas, pessoa)
		} else {
			removido = true
		}
	}
	pessoas = novasPessoas

	if removido {
		w.WriteHeader(http.StatusNoContent)
	} else {
		w.WriteHeader(http.StatusNotFound)
	}
}