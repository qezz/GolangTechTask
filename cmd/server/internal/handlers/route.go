package handlers

import (
	"fmt"
	"github.com/go-chi/chi"
	"net/http"
)

func defaultPage(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Works!")
	return
}

func Routes(store Store) *chi.Mux {
	r := chi.NewRouter()

	bh := buffHandler{
		store: store,
	}

	r.Route("/", func(r chi.Router) {
		r.Get("/", defaultPage)
	})

	r.Route("/buff", func(r chi.Router) {
		r.Get("/{id}", bh.GetBuff)
		r.Post("/", bh.CreateBuff)
	})

	sh := streamHandler{
		store: store,
	}

	r.Route("/stream", func(r chi.Router) {
		r.Post("/", sh.CreateStream)
		r.Get("/", defaultPage)
		r.Get("/page={page:[0-9]+}", sh.ListStreams)

		r.Route("/{id}", func(r chi.Router) {
			r.Get("/", sh.GetStream)

		})
	})

	return r
}
