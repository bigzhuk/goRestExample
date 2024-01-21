package main

import (
	"github.com/stretchr/testify/require"
	"io"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

type TestCase struct {
	Address          string
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
			Address: "http://localhost:8080/artist/1",
			Request: "",
			ExpectedResponse: ExpectedResponse{
				StatusCode: http.StatusOK,
				Body:       "{\"id\":\"1\",\"name\":\"30 Seconds To Mars\",\"born\":\"1998\",\"genre\":\"alternative\",\"songs\":[\"The Kill\",\"A Beautiful Lie\",\"Attack\",\"Live Like A Dream\"]}",
			},
		},
		{
			Address: "http://localhost:8080/artist/2",
			Request: "",
			ExpectedResponse: ExpectedResponse{
				StatusCode: http.StatusOK,
				Body:       "{\"id\":\"2\",\"name\":\"Garbage\",\"born\":\"1994\",\"genre\":\"alternative\",\"songs\":[\"Queer\",\"Shut Your Mouth\",\"Cup of Coffee\",\"Til the Day I Die\"]}",
			},
		},
	}
}

func TestGetArtist(t *testing.T) {
	testCaseList := getArtistTestCases()

	for _, testCase := range testCaseList {
		resp, err := http.Get(testCase.Address)

		require.NoError(t, err)

		body, err := io.ReadAll(resp.Body)
		assert.Equal(t, testCase.ExpectedResponse.StatusCode, resp.StatusCode)
		assert.Equal(t, testCase.ExpectedResponse.Body, string(body))
		if resp != nil && resp.Body != nil {
			resp.Body.Close()
		}
	}
}

func TestSaveArtist(t *testing.T) {
	testCaseList := getArtistTestCases()

	for _, testCase := range testCaseList {
		resp, err := http.Get(testCase.Address)

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
