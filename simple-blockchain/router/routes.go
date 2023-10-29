package router

import (
	"crypto/md5"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"simple-blockchain/models"
)

var BlockChain *models.Blockchain

func newBook(w http.ResponseWriter, r *http.Request) {
	var book models.Book

	if err := json.NewDecoder(r.Body).Decode(&book); err != nil{
		WriteJSON(w, "could not create new book", err)
		return
	}

	h := md5.New()
	io.WriteString(h, book.ISBN + book.PublicationDate)
	book.ID = fmt.Sprintf("%x", h.Sum(nil))
	resp, err := json.MarshalIndent(book, "", " ")
	if err != nil{
		WriteJSON(w, "could not save book data", err)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(resp)
}

func writeBlock(w http.ResponseWriter, r *http.Request) {
	var checkoutItem models.BookCheckout
	if err := json.NewDecoder(r.Body).Decode(&checkoutItem); err != nil {
		WriteJSON(w, "could not write block", err)
		return
	}

	BlockChain.AddBlock(checkoutItem)
	resp, err := json.MarshalIndent(checkoutItem, "", " ")
	if err != nil {
		WriteJSON(w, "could not marshal payload", err)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write(resp)
}

func getBlockchain(w http.ResponseWriter, r *http.Request){
	jBytes, err := json.MarshalIndent(BlockChain.Blocks, "", " ")
	if err != nil{
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(err)
		return
	}
	io.WriteString(w, string(jBytes))
}

func WriteJSON(w http.ResponseWriter, msg string, err error){
	w.WriteHeader(http.StatusInternalServerError)
	log.Printf(msg, " %v", err)
	w.Write([]byte(msg))
}