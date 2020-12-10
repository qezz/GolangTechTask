package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/go-chi/chi"
)

type buffHandler struct {
	store BuffStore
}

func (b *buffHandler) GetBuff(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if id == "" {
		http.Error(w, "missing buff id in path", http.StatusBadRequest)
		log.Print("[Error] missing buff id in path")
		return
	}

	parsedid, err := strconv.ParseUint(id, 10, 64)
	if err != nil {
		http.Error(w, fmt.Sprintf("failed to parse id, %v", err), http.StatusBadRequest)
		log.Printf("[Error] failed to parse id, %v", err)
		return
	}

	buff, err := b.store.GetBuff(parsedid)
	if err != nil {
		http.Error(w, fmt.Sprintf("failed to get buff, %v", err), http.StatusInternalServerError)
		log.Printf("[Error] failed to get buff, %v", err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(buff); err != nil {
		http.Error(w, fmt.Sprintf("failed encode buff, %v", err), http.StatusInternalServerError)
		log.Printf("[Error] failed encode buff, %v", err)
		return
	}

	log.Printf("GetBuff request is ok: User-Agent - %v, IP-Address - %v", r.UserAgent(), IPAddress(r))
}

func (b *buffHandler) CreateBuff(w http.ResponseWriter, r *http.Request) {
	var cbr CreateBuff
	if err := json.NewDecoder(r.Body).Decode(&cbr); err != nil {
		http.Error(w, fmt.Sprintf("failed to unmarshal request, %v", err), http.StatusBadRequest)
		log.Printf("[Error] failed to unmarshal request, %v", err)
		return
	}

	buff := &Buff{
		ID:       0,
		Question: cbr.Question,
		Answers:  cbr.Answers,
	}

	id, err := b.store.SetBuff(buff)
	if err != nil {
		http.Error(w, fmt.Sprintf("failed create buff, %v", err), http.StatusInternalServerError)
		log.Printf("[Error] failed create buff, %v", err)
		return
	}

	buff, err = b.store.GetBuff(id)
	if err != nil {
		http.Error(w, fmt.Sprintf("failed to get buff, %v", err), http.StatusInternalServerError)
		log.Printf("[Error] failed to get buff, %v", err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(buff); err != nil {
		http.Error(w, fmt.Sprintf("failed encode buff, %v", err), http.StatusInternalServerError)
		log.Printf("[Error] failed encode buff, %v", err)
		return
	}

	log.Printf("CreateBuff request is ok: User-Agent - %v, IP-Address - %v", r.UserAgent(), IPAddress(r))
}
