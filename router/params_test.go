package router

import (
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/julienschmidt/httprouter"
	"github.com/stretchr/testify/assert"
)

func TestRouterParams(t *testing.T) {

	router := httprouter.New()

	router.GET("/products/:id", func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		id := p.ByName("id")
		text := "Product " + id
		fmt.Fprint(w, text)
	})

	request := httptest.NewRequest("GET", "http://localhost:3000/products/1", nil)
	recorder := httptest.NewRecorder()

	router.ServeHTTP(recorder, request)

	response := recorder.Result()

	body, _ := io.ReadAll(response.Body)

	assert.Equal(t, "Product 1", string(body))

}

func TestParamsComments(t *testing.T) {
	router := httprouter.New()

	router.GET("/comments/:id", func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		text := "Comment " + p.ByName("id")
		fmt.Fprint(w, text)
	})

	request := httptest.NewRequest("GET", "http://localhost:3000/comments/1", nil)
	recorder := httptest.NewRecorder()

	router.ServeHTTP(recorder, request)

	response := recorder.Result()

	body, _ := io.ReadAll(response.Body)

	assert.Equal(t, "Comment 1", string(body))

}

func TestParamsNamed(t *testing.T) {
	router := httprouter.New()

	router.GET("/products/:id/items/:itemid", func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		text := "Product " + p.ByName("id") + " Item " + p.ByName("itemid")
		fmt.Fprint(w, text)
	})

	request := httptest.NewRequest("GET", "http://localhost:3000/products/1/items/1", nil)
	recorder := httptest.NewRecorder()

	router.ServeHTTP(recorder, request)

	response := recorder.Result()

	body, _ := io.ReadAll(response.Body)

	assert.Equal(t, "Product 1 Item 1", string(body))
}

func TestParamsCatchAll(t *testing.T) {
	router := httprouter.New()

	router.GET("/images/*image", func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		text := "image : " + p.ByName("image")
		fmt.Fprint(w, text)
	})

	request := httptest.NewRequest("GET", "http://localhost:3000/images/small/profile.png", nil)
	recorder := httptest.NewRecorder()

	router.ServeHTTP(recorder, request)

	response := recorder.Result()

	body, _ := io.ReadAll(response.Body)

	assert.Equal(t, "image : /small/profile.png", string(body))
}

func TestPanicHandler(t *testing.T) {
	router := httprouter.New()

	router.PanicHandler = func(w http.ResponseWriter, r *http.Request, i interface{}) {
		fmt.Fprint(w, "Panic : ", i)
	}

	router.GET("/", func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		panic("Ups")
	})

	request := httptest.NewRequest("GET", "http://localhost:3000/", nil)
	recorder := httptest.NewRecorder()

	router.ServeHTTP(recorder, request)

	response := recorder.Result()

	body, _ := io.ReadAll(response.Body)

	assert.Equal(t, "Panic : Ups", string(body))
}

func TestNotFoundHandler(t *testing.T) {
	router := httprouter.New()

	router.NotFound = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, "Ga Ketemu")
	})

	// router.GET("/", func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	// 	fmt.Fprint(w, "Hello World")
	// })

	request := httptest.NewRequest("GET", "http://localhost:3000/", nil)
	recorder := httptest.NewRecorder()

	router.ServeHTTP(recorder, request)

	response := recorder.Result()

	body, _ := io.ReadAll(response.Body)

	assert.Equal(t, "Ga Ketemu", string(body))
}

func TestMethodNotAllowed(t *testing.T) {
	router := httprouter.New()

	router.MethodNotAllowed = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, "Ga Boleh")
	})

	router.POST("/", func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		fmt.Fprint(w, "Hello")
	})

	request := httptest.NewRequest("POST", "http://localhost:3000/", nil)
	recorder := httptest.NewRecorder()

	router.ServeHTTP(recorder, request)

	response := recorder.Result()

	body, _ := io.ReadAll(response.Body)

	assert.Equal(t, "Ga Boleh", string(body))
}

func TestAllErrorHandler(t *testing.T) {

	router := httprouter.New()

	router.PanicHandler = func(w http.ResponseWriter, r *http.Request, i interface{}) {
		fmt.Fprint(w, "Panic :", i)
	}

	router.NotFound = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, "Ga Ketemu")
	})

	router.MethodNotAllowed = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, "Ga Boleh")
	})

}

type LogMiddleware struct {
	http.Handler
}

func (middleware *LogMiddleware) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	fmt.Println("Receive Request")
	middleware.Handler.ServeHTTP(writer, request)
}

func TestMiddleware(t *testing.T) {
	router := httprouter.New()

	router.GET("/", func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		fmt.Fprint(w, "Hello")
	})

	request := httptest.NewRequest("GET", "http://localhost:3000/", nil)
	recorder := httptest.NewRecorder()

	middleware := LogMiddleware{Handler: router}

	middleware.ServeHTTP(recorder, request)

	response := recorder.Result()

	body, _ := io.ReadAll(response.Body)

	assert.Equal(t, "Hello", string(body))

}
