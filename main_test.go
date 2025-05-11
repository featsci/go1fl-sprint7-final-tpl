package main

import (
	"net/http"
	"net/http/httptest"
	"slices"
	"strconv"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCafeNegative(t *testing.T) {
	handler := http.HandlerFunc(mainHandle)

	requests := []struct {
		request string
		status  int
		message string
	}{
		{"/cafe", http.StatusBadRequest, "unknown city"},
		{"/cafe?city=omsk", http.StatusBadRequest, "unknown city"},
		{"/cafe?city=tula&count=na", http.StatusBadRequest, "incorrect count"},
	}
	for _, v := range requests {
		response := httptest.NewRecorder()
		req := httptest.NewRequest("GET", v.request, nil)
		handler.ServeHTTP(response, req)

		assert.Equal(t, v.status, response.Code)
		assert.Equal(t, v.message, strings.TrimSpace(response.Body.String()))
	}
}

func TestCafeWhenOk(t *testing.T) {
	handler := http.HandlerFunc(mainHandle)

	requests := []string{
		"/cafe?count=2&city=moscow",
		"/cafe?city=tula",
		"/cafe?city=moscow&search=ложка",
	}
	for _, v := range requests {
		response := httptest.NewRecorder()
		req := httptest.NewRequest("GET", v, nil)

		handler.ServeHTTP(response, req)

		assert.Equal(t, http.StatusOK, response.Code)
	}
}

/*
Требования к каждой тестовой функции даны ниже. При реализации используйте функции пакетов assert и require. С помощью require.Equal() следует проверить, что запрос успешно обработан, а функции assert использовать для проверки корректности результатов.
*/

// проверяет работу сервера при разных значениях параметра count;
func TestCafeCount(t *testing.T) {
	countTest := 0
	handler := http.HandlerFunc(mainHandle)

	requests := []struct {
		count int // передаваемое значение count
		want  int // ожидаемое количество кафе в ответе
	}{
		{count: 0, want: 0}, // len(cafeList["moscow"])},
		{count: 1, want: 1}, // len(cafeList["tula"])},
		{count: 2, want: 2},
		{count: 100, want: len(cafeList["moscow"])},
	}

	for _, v := range requests {
		response := httptest.NewRecorder()
		req := httptest.NewRequest("GET", "/cafe?city=moscow&count="+strconv.Itoa(v.count), nil)
		handler.ServeHTTP(response, req)

		resBody := strings.Split(strings.TrimSpace(response.Body.String()), ",")
		if !slices.Contains(resBody, "") {
			countTest = len(resBody)
		}

		require.Equal(t, http.StatusOK, response.Code)
		assert.Equal(t, v.want, countTest)
		require.Equal(t, v.want, countTest)
	}

}

// проверяет результат поиска кафе по указанной подстроке в параметре search.
func TestCafeSearch(t *testing.T) {
	sWant := 0
	handler := http.HandlerFunc(mainHandle)

	requests := []struct {
		search    string // передаваемое значение search
		wantCount int    // ожидаемое количество кафе в ответе
	}{
		{search: "фасоль", wantCount: 0},
		{search: "кофе", wantCount: 2},
		{search: "вилка", wantCount: 1},
	}

	// fmt.Println(requests)
	for _, v := range requests {
		response := httptest.NewRecorder()
		req := httptest.NewRequest("GET", "/cafe?city=moscow&search="+v.search, nil)
		handler.ServeHTTP(response, req)

		if !strings.Contains(req.RequestURI, strings.ToLower(v.search)) {
			t.Errorf("not contain %s, in %s", strings.ToLower(v.search), req.RequestURI)
		}

		resBody := strings.Split(strings.TrimSpace(response.Body.String()), ",")
		if !slices.Contains(resBody, "") {
			sWant = len(resBody)
		}

		require.Equal(t, http.StatusOK, response.Code)
		assert.Equal(t, v.wantCount, sWant)
		require.Equal(t, v.wantCount, sWant)
	}

}
