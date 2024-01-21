package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/go-chi/chi"
	"github.com/go-playground/validator"
	"gitlab.goodsteam.tech/goRestExample/model"
	"net/http"
	"strconv"
)

const defaultLimit = 10

var artists = map[string]model.Artist{
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

func validateArtistListRequest(r *http.Request) (*model.Filter, error) {
	var err error
	var filter model.Filter

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

func getArtistList(filter *model.Filter) map[string]model.Artist {
	// этот метод должен быть методом сервиса и под капотом обращаться за данным в storage.
	// Artist - структура доменной модели (бизнес логики). При походе в БД и возврате ответа целесообразно сделать отедельные dto,
	// куда нужно смапить данные из Artist
	artistList := make([]model.Artist, 0)
	for _, artist := range artists {
		if (filter.Genre == artist.Genre || filter.Genre == "") && (filter.Born == artist.Born || filter.Born == "") {
			artistList = append(artistList, artist)
		}
	}

	return convertArtistsToMap(artistList)
}

func convertArtistsToMap(artistList []model.Artist) map[string]model.Artist {
	artistListMap := make(map[string]model.Artist)
	for _, artist := range artistList {
		artistListMap[artist.ID] = artist
	}

	return artistListMap
}

func SaveArtist(w http.ResponseWriter, r *http.Request) {
	var artist model.Artist
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

	fmt.Println(artist)

	validate := validator.New()
	if err := validate.Struct(artist); err != nil {
		validationErrors := err.(validator.ValidationErrors)
		fmt.Println("validationErrors", validationErrors)
		http.Error(w, fmt.Errorf("ошибка валидации запроса: %w", validationErrors).Error(), http.StatusBadRequest)
		return
	}

	artists[artist.ID] = artist

	respMsg := "Артист успешно сохранен"
	resp, err := json.Marshal(respMsg)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	w.Write(resp)
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
	// домашнее задание: сделать методы для добавления и удаления списка персен кокретного артиста.
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
	// в данном случае сделал как принято в спецификации api у нас в компании: хвостик api -  сущность/дейтсвие
	r.Get("/artist/list", GetArtistList)
	r.Post("/artist/save", SaveArtist)
	r.Get("/artist/{id}", GetArtist)
	r.Get("/artist/{id}/song/list", getArtistSongs)

	// запускаем сервер
	if err := http.ListenAndServe(":8080", r); err != nil {
		fmt.Printf("Ошибка при запуске сервера: %s", err.Error())
		return
	}
}
