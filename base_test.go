package main

import (
   "net/http"
   "net/http/httptest"
   "testing"
   "github.com/stretchr/testify/assert"
   _ "fmt"
)

func performRequest(r http.Handler, method, path string) *httptest.ResponseRecorder {
   req, _ := http.NewRequest(method, path, nil)
   w := httptest.NewRecorder()
   r.ServeHTTP(w, req)
   return w
}

func Test_GetProductFromCart(t *testing.T) {
   router := SetupRouter()

   w := performRequest(router, "GET", "/cart/getProducts")

   assert.Equal(t, http.StatusOK, w.Code)
}
