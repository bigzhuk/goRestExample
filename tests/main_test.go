package main

import (
	"bytes"
	"github.com/stretchr/testify/require"
	"io"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

// тестируем API нашего приложения. Наш пример может быть запущен отдельно от приложения.
// Единственное что связывает тесты с кодом приложения dto запросов ответов (в нашем случае мы используем модель, что неправильно)

type TestCase struct {
	URL              string
	Request          string
	ExpectedResponse ExpectedResponse
}

type ExpectedResponse struct {
	StatusCode int
	Body       string
}

func getArtistTestCases() []TestCase {
	return []TestCase{
		{
			URL:     "http://localhost:8080/artist/1",
			Request: "",
			ExpectedResponse: ExpectedResponse{
				StatusCode: http.StatusOK,
				Body:       "{\"id\":\"1\",\"name\":\"30 Seconds To Mars\",\"born\":1998,\"genre\":\"alternative\",\"songs\":[\"The Kill\",\"A Beautiful Lie\",\"Attack\",\"Live Like A Dream\"]}",
			},
		},
		{
			URL:     "http://localhost:8080/artist/2",
			Request: "",
			ExpectedResponse: ExpectedResponse{
				StatusCode: http.StatusOK,
				Body:       "{\"id\":\"2\",\"name\":\"Garbage\",\"born\":1994,\"genre\":\"alternative\",\"songs\":[\"Queer\",\"Shut Your Mouth\",\"Cup of Coffee\",\"Til the Day I Die\"]}",
			},
		},
		{
			URL:     "http://localhost:8080/artist/10",
			Request: "",
			ExpectedResponse: ExpectedResponse{
				StatusCode: http.StatusNoContent,
				Body:       "",
			},
		},
		{
			URL:     "http://localhost:8080/artistNotFound/1",
			Request: "",
			ExpectedResponse: ExpectedResponse{
				StatusCode: http.StatusNotFound,
				Body:       "404 page not found\n",
			},
		},
	}
}

func TestGetArtist(t *testing.T) {
	testCaseList := getArtistTestCases()

	for _, testCase := range testCaseList {
		resp, err := http.Get(testCase.URL)

		require.NoError(t, err)

		body, err := io.ReadAll(resp.Body)
		assert.Equal(t, testCase.ExpectedResponse.StatusCode, resp.StatusCode)
		assert.Equal(t, testCase.ExpectedResponse.Body, string(body))
		if resp != nil && resp.Body != nil {
			resp.Body.Close()
		}
	}
	t.Run()
}

func saveArtistTestCases() []TestCase {
	return []TestCase{
		{
			URL:     "http://localhost:8080/artist/save",
			Request: "{\n    \"id\": \"4\",\n    \"name\": \"Руки вверх\",\n    \"born\": 1998,\n    \"genre\": \"pop\",\n    \"songs\": [\n        \"Крошка моя\"\n    ]\n}",
			ExpectedResponse: ExpectedResponse{
				StatusCode: http.StatusCreated,
				Body:       "\"Артист успешно сохранен\"",
			},
		},
		{
			URL:     "http://localhost:8080/artist/save",
			Request: "{\n    \"id\": \"4\",\n    \"name\": \"Руки вверх\",\n    \"born\": 1800,\n    \"genre\": \"pop\",\n    \"songs\": [\n        \"Крошка моя\"\n    ]\n}",
			ExpectedResponse: ExpectedResponse{
				StatusCode: http.StatusBadRequest,
				Body:       "ошибка валидации запроса: Key: 'Artist.Born' Error:Field validation for 'Born' failed on the 'gt' tag\n",
			},
		},
	}
}

func TestSaveArtist(t *testing.T) {
	testCaseList := saveArtistTestCases()
	for _, testCase := range testCaseList {
		resp, err := http.Post(testCase.URL, "application/json", bytes.NewBuffer([]byte(testCase.Request)))
		require.NoError(t, err)

		body, err := io.ReadAll(resp.Body)
		assert.Equal(t, testCase.ExpectedResponse.StatusCode, resp.StatusCode)
		assert.Equal(t, testCase.ExpectedResponse.Body, string(body))
		if resp != nil && resp.Body != nil {
			resp.Body.Close()
		}
	}
}

// Минусы - тест зависит от данных. Решение - тесты разворачиваются в отдельном контуре. При старте данные заливаются из фикстур.
// Тест должен тестировать логику, а не данные.
