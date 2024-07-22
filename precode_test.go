package main

import (
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var cafeList = map[string][]string{
	"moscow": []string{"Мир кофе", "Сладкоежка", "Кофе и завтраки", "Сытый студент"},
}

func mainHandle(w http.ResponseWriter, req *http.Request) {
	countStr := req.URL.Query().Get("count")
	if countStr == "" {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("count missing"))
		return
	}

	count, err := strconv.Atoi(countStr)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("wrong count value"))
		return
	}

	city := req.URL.Query().Get("city")

	cafe, ok := cafeList[city]

	if !ok {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("wrong city value"))
		return
	}

	if count > len(cafe) {
		count = len(cafe)
	}

	answer := strings.Join(cafe[:count], ",")

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(answer))
}

// Проверяем ответ сервера если запрос сформирован корректно
func TestMainHandlerCorrectResponse(t *testing.T) {
	req, err := http.NewRequest("GET", "/?count=4&city=moscow", nil) // Запрос к серверу с количеством больше, чем в списке

	responseRecorder := httptest.NewRecorder()
	handler := http.HandlerFunc(mainHandle)
	handler.ServeHTTP(responseRecorder, req)

	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, responseRecorder.Code)
	assert.NotEmpty(t, responseRecorder.Body)
}

// В запросе передается город которого нет в списке
func TestMainHandlerWhenCityNotExist(t *testing.T) {
	//Запрос города которого нет в списке
	req, err := http.NewRequest("GET", "/?count=4&city=kazan", nil) // Запрос к серверу с неправиильнм городом

	responseRecorder := httptest.NewRecorder()
	handler := http.HandlerFunc(mainHandle)
	handler.ServeHTTP(responseRecorder, req)

	require.NoError(t, err)
	assert.Equal(t, http.StatusBadRequest, responseRecorder.Code)
	assert.NotEqual(t, responseRecorder.Body, "wrong city value")
}

// В запросе передается неправильное количество кафешек
func TestMainHandlerWhenCountMoreThanTotal(t *testing.T) {
	totalCount := 4

	// Запрос с кол-вом кофе больше чем в списке
	req, err := http.NewRequest("GET", "/?count="+strconv.Itoa(totalCount+1)+"&city=moscow", nil) // Запрос к серверу с количеством больше, чем в списке

	responseRecorder := httptest.NewRecorder()
	handler := http.HandlerFunc(mainHandle)
	handler.ServeHTTP(responseRecorder, req)

	require.NoError(t, err)
	assert.Equal(t, len(strings.Split(responseRecorder.Body.String(), ",")), totalCount)

}
