package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/go-chi/chi"
)

type streamHandler struct {
	store Store
}

func (s *streamHandler) GetStream(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	if id == "" {
		http.Error(w, fmt.Sprintf("missing stream id in path"), http.StatusBadRequest)
		log.Print("[Error] missing stream id in path")
		return
	}

	parsedid, err := strconv.ParseUint(id, 10, 64)
	if err != nil {
		http.Error(w, fmt.Sprintf("failed to parse id, %v", err), http.StatusBadRequest)
		log.Printf("[Error] failed to parse id, %v", err)
		return
	}

	stream, err := s.store.GetStream(parsedid)
	if err != nil {
		http.Error(w, fmt.Sprintf("failed to get stream, %v", err), http.StatusInternalServerError)
		log.Printf("[Error] failed to get stream, %v", err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(stream); err != nil {
		http.Error(w, fmt.Sprintf("failed encode stream, %v", err), http.StatusInternalServerError)
		log.Printf("[Error] failed encode stream, %v", err)
		return
	}

	log.Printf("GetStream request is ok: User-Agent - %v, IP-Address - %v", r.UserAgent(), IPAddress(r))

}

func (s *streamHandler) CreateStream(w http.ResponseWriter, r *http.Request) {
	var csr CreateStream
	if err := json.NewDecoder(r.Body).Decode(&csr); err != nil {
		http.Error(w, fmt.Sprintf("failed to unmarshal request, %v", err), http.StatusBadRequest)
		log.Printf("[Error] failed to unmarshal request, %v", err)
		return
	}

	// User could send any data, for instance non-existing buff id,
	// we are checking the id and if any of those doesn;t exist
	// return error.

	for _, b := range csr.Buffs {

		_, err := s.store.GetBuff(b)
		if err != nil {
			http.Error(w, fmt.Sprintf("Wrong request: failed create stream, %v", err), http.StatusInternalServerError)
			log.Printf("[Error] Wrong request: failed create stream, %v", err)
			return

		}

	}

	stream := &Stream{
		ID:    0,
		Name:  csr.Name,
		Buffs: csr.Buffs,
	}

	id, err := s.store.SetStream(stream)
	if err != nil {
		http.Error(w, fmt.Sprintf("failed create stream, %v", err), http.StatusInternalServerError)
		log.Printf("[Error] failed create stream, %v", err)
		return
	}

	stream, err = s.store.GetStream(id)
	if err != nil {
		http.Error(w, fmt.Sprintf("failed to get stream, %v", err), http.StatusInternalServerError)
		log.Printf("[Error] failed to get stream, %v", err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(stream); err != nil {
		http.Error(w, fmt.Sprintf("failed encode stream, %v", err), http.StatusInternalServerError)
		log.Printf("[Error] failed encode stream, %v", err)
		return
	}

	log.Printf("CreateStream request is ok: User-Agent - %v, IP-Address - %v", r.UserAgent(), IPAddress(r))

}

func (s *streamHandler) ListStreams(w http.ResponseWriter, r *http.Request) {

	total, _ := s.store.Count()

	if total == 0 {
		http.Error(w, fmt.Sprintf("No Data"), http.StatusInternalServerError)
		log.Printf("[Error] No Data, please fill a Memory Storage")
		return

	}

	page := chi.URLParam(r, "page")
	pageSize := r.URL.Query().Get("pagesize")

	if page == "" {
		page = "1"
	}

	if pageSize == "" {
		pageSize = "25"
	}

	parsedPage, err := strconv.ParseUint(page, 10, 64)
	if err != nil {
		http.Error(w, fmt.Sprintf("failed to parse page, %v", err), http.StatusBadRequest)
		log.Printf("[Error] failed to parse page, %v", err)
		return
	}

	parsedPageSize, err := strconv.ParseUint(pageSize, 10, 64)
	if err != nil {
		http.Error(w, fmt.Sprintf("failed to parse pageSize, %v", err), http.StatusBadRequest)
		log.Printf("[Error] failed to parse pageSize, %v", err)
		return
	}

	if total+parsedPageSize < parsedPage*parsedPageSize {
		http.Error(w, fmt.Sprintf("No Data"), http.StatusBadRequest)
		log.Printf("[Warning] No Data, wrong Parameters")
		return

	}

	index, _ := s.store.Streams()

	keys := make([]uint64, 0, len(index))
	values := make([]*Stream, 0, len(index))
	for k, v := range index {
		keys = append(keys, k)
		values = append(values, v)
	}

	var stream *Stream

	startPageEntry := parsedPage*parsedPageSize - parsedPageSize
	endPageEntry := parsedPage * parsedPageSize

	if endPageEntry > total {
		endPageEntry = total

	}

	result := make([]map[string]interface{}, endPageEntry-startPageEntry)

	if startPageEntry > endPageEntry {
		log.Printf("[Error] Wrong Parameters")
		return
	}

	for j, i := 0, startPageEntry; i < endPageEntry; i, j = i+1, j+1 {

		stream, _ = s.store.GetStream(i + 1)

		if stream != nil {

			result[j] = map[string]interface{}{"name": stream.Name, "id": stream.ID}
		}
	}

	pretty, err := json.MarshalIndent(result, "", "    ")
	if err != nil {
		http.Error(w, fmt.Sprintf("failed encode stream, %v", err), http.StatusInternalServerError)
		log.Printf("[Error] failed encode stream, %v", err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	fmt.Fprintf(w, string(pretty))

	log.Printf("ListStreams request is ok: User-Agent - %v, IP-Address - %v", r.UserAgent(), IPAddress(r))

}
