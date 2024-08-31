package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	_ "github.com/go-sql-driver/mysql"
)

type Pessoa struct {
	ID   int    `json:"id"`
	Nome string `json:"nome"`
}

var db *sql.DB

func initDB() {
	var err error
	db, err = sql.Open("mysql", "aluno:iftm@tcp(127.0.0.1:3306)/mydatabase")
	if err != nil {
		panic(err)
	}

	_, err = db.Exec("CREATE TABLE IF NOT EXISTS pessoas (id INT AUTO_INCREMENT PRIMARY KEY, nome VARCHAR(255) NOT NULL)")
	if err != nil {
		panic(err)
	}
}

func getListPessoas(w http.ResponseWriter, r *http.Request) {
	rows, err := db.Query("SELECT id, nome FROM pessoas")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var pessoas []Pessoa
	for rows.Next() {
		var pessoa Pessoa
		if err := rows.Scan(&pessoa.ID, &pessoa.Nome); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		pessoas = append(pessoas, pessoa)
	}
	json.NewEncoder(w).Encode(pessoas)
}

func getPessoa(w http.ResponseWriter, r *http.Request) {
	nome := r.URL.Query().Get("nome")
	idStr := r.URL.Query().Get("id")

	var pessoa Pessoa
	var err error
	if idStr != "" {
		id, _ := strconv.Atoi(idStr)
		err = db.QueryRow("SELECT id, nome FROM pessoas WHERE id = ?", id).Scan(&pessoa.ID, &pessoa.Nome)
	} else if nome != "" {
		err = db.QueryRow("SELECT id, nome FROM pessoas WHERE nome = ?", nome).Scan(&pessoa.ID, &pessoa.Nome)
	}

	if err != nil {
		if err == sql.ErrNoRows {
			json.NewEncoder(w).Encode(nil)
		} else {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}
	json.NewEncoder(w).Encode(pessoa)
}

func postPessoa(w http.ResponseWriter, r *http.Request) {
	var pessoa Pessoa
	_ = json.NewDecoder(r.Body).Decode(&pessoa)
	result, err := db.Exec("INSERT INTO pessoas (nome) VALUES (?)", pessoa.Nome)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	id, err := result.LastInsertId()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	pessoa.ID = int(id)
	json.NewEncoder(w).Encode(pessoa)
}

func deletePessoa(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.Atoi(r.URL.Query().Get("id"))

	_, err := db.Exec("DELETE FROM pessoas WHERE id = ?", id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	fmt.Fprintf(w, "Pessoa com ID %d foi deletada com sucesso.", id)
}

func main() {
	initDB()

	http.HandleFunc("/getListPessoas", getListPessoas)
	//curl localhost:3333/getListPessoas
	http.HandleFunc("/getPessoa", getPessoa)
	//curl "localhost:3333/getPessoa?nome=Henrique"
	//curl "localhost:3333/getPessoa?id=1"
	http.HandleFunc("/postPessoa", postPessoa)
	//curl -X POST -H "Content-Type: application/json" -d '{"nome":"Henrique"}' localhost:3333/postPessoa
	http.HandleFunc("/deletePessoa", deletePessoa)
	//curl -X DELETE "localhost:3333/deletePessoa?id=1"

	_ = http.ListenAndServe(":3333", nil)
}