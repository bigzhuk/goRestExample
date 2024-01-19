package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/go-chi/chi"
	"net/http"
	"strconv"
)

const defaultLimit = 10

type Artist struct {
	ID    string   `json:"id"`    // id коллектива
	Name  string   `json:"name"`  // название группы
	Born  string   `json:"born"`  // год основания группы
	Genre string   `json:"genre"` // жанр
	Songs []string `json:"songs"` // популярные песни, это слайс строк, так как песен может быть несколько
}

type Filter struct {
	Genre  string `json:"genre"`
	Born   string `json:"born"`
	Limit  int    `json:"limit"`
	Offset int    `json:"offset"`
}

var artists = map[string]Artist{
	"1": {
		ID:    "1",
		Name:  "30 Seconds To Mars",
		Born:  "1998",
		Genre: "alternative",
		Songs: []string{
			"The Kill",
			"A Beautiful Lie",
			"Attack",
			"Live Like A Dream",
		},
	},
	"2": {
		ID:    "2",
		Name:  "Garbage",
		Born:  "1994",
		Genre: "alternative",
		Songs: []string{
			"Queer",
			"Shut Your Mouth",
			"Cup of Coffee",
			"Til the Day I Die",
		},
	},
	"3": {
		ID:    "3",
		Name:  "Queen",
		Born:  "1970",
		Genre: "rock",
		Songs: []string{
			"We Will Rock You",
			"I want to break free",
			"We are the champions",
		},
	},
}

func validateArtistListRequest(r *http.Request) (*Filter, error) {
	var err error
	var filter Filter

	limit := r.URL.Query().Get("limit")
	if limit != "" {
		filter.Limit, err = strconv.Atoi(r.URL.Query().Get("limit"))
		if err != nil {
			return nil, fmt.Errorf("limit shoul be an integer")
		}
	} else {
		filter.Limit = defaultLimit // тут можно было ругнуться, а можно сделать так, чтобы приложение была менее хрупким.
	}

	offset := r.URL.Query().Get("offset")
	if offset != "" {
		filter.Offset, err = strconv.Atoi(r.URL.Query().Get("offset"))
		if err != nil {
			return nil, fmt.Errorf("offset shoul be an integer")
		}
	}

	filter.Genre = r.URL.Query().Get("genre")
	filter.Born = r.URL.Query().Get("born")

	return &filter, err
}

func GetArtistList(w http.ResponseWriter, r *http.Request) {

	filter, err := validateArtistListRequest(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// сериализуем данные из слайса artists
	resp, err := json.Marshal(getArtistList(filter))
	// FYI: идеоматически в го принято называть методы без get, просто artistList, но мне не нравится метод без глагола.
	// маршалинг в json это формат данных. Если мы хотим чтобы наш сервис поддерживал разные форматы, целесообразно спрятать это зад интерфейс.
	// но в большинстве случае это избыточно.
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// в заголовок записываем тип контента, у нас это данные в формате JSON
	w.Header().Set("Content-Type", "application/json")
	// так как все успешно, то статус OK
	w.WriteHeader(http.StatusOK)
	// записываем сериализованные в JSON данные в тело ответа
	w.Write(resp)
}

func getArtistList(filter *Filter) map[string]Artist {
	// этот метод должен быть методом сервиса и под капотом обращаться за данным в storage.
	// Atrist - структура доменной модели (бизнес логики). При походе в БД и возврате ответа целесообразно сделать отедельные dto,
	// куда нужно смапить данные из Artist
	artistList := make([]Artist, 0)
	for _, artist := range artists {
		if (filter.Genre == artist.Genre || filter.Genre == "") && (filter.Born == artist.Born || filter.Born == "") {
			artistList = append(artistList, artist)
		}
	}

	return convertArtistsToMap(artistList)
}

func convertArtistsToMap(artistList []Artist) map[string]Artist {
	artistListMap := make(map[string]Artist)
	for _, artist := range artistList {
		artistListMap[artist.ID] = artist
	}

	return artistListMap
}

func SaveArtist(w http.ResponseWriter, r *http.Request) {
	var artist Artist
	var buf bytes.Buffer

	_, err := buf.ReadFrom(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err = json.Unmarshal(buf.Bytes(), &artist); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	artists[artist.ID] = artist

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
}

func GetArtist(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	artist, ok := artists[id]
	if !ok {
		http.Error(w, "Артист не найден", http.StatusNoContent)
		return
	}

	resp, err := json.Marshal(artist)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(resp)
}

func getArtistSongs(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	artist, ok := artists[id]
	if !ok {
		http.Error(w, "Артист не найден", http.StatusNoContent)
		return
	}

	resp, err := json.Marshal(artist.Songs)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(resp)
}

func main() {
	// создаём новый роутер
	r := chi.NewRouter()

	// делать одинаковые маршруты для разных методов (например POST и GET) можно, но я бы не советовал.
	// в данном случае сделал как принято в спецификации api у нас в компании сущность/дейтсвие
	r.Get("/artistList/get", GetArtistList)
	r.Post("/artist/save", SaveArtist)
	r.Get("/artist/{id}", GetArtist)

	// запускаем сервер
	if err := http.ListenAndServe(":8080", r); err != nil {
		fmt.Printf("Ошибка при запуске сервера: %s", err.Error())
		return
	}
}
