package server_test

import (
	"encoding/json"
	"github.com/Zarux/BGServer/internal/app/server"
	"log"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"
)

func TestGetBoardGames(t *testing.T) {
	var (
		testResp []server.BoardGame
	)
	testQuery := "Agricola"
	testLimit := 10

	expectedId := uint64(25)
	expectedName := "Agricola (2007)"

	reqBadRequest, err := http.NewRequest("GET", "/boardGames", nil)
	if err != nil {
		t.Fatal(err)
	}
	log.Println(reqBadRequest.URL.String())
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(server.GetBoardGames)
	handler.ServeHTTP(rr, reqBadRequest)
	if status := rr.Code; status != http.StatusBadRequest {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusBadRequest)
	}

	reqGoodRequest, err := http.NewRequest("GET", "/boardGames", nil)
	query := reqGoodRequest.URL.Query()
	query.Add("search", testQuery)
	query.Add("limit", strconv.Itoa(testLimit))
	reqGoodRequest.URL.RawQuery = query.Encode()
	log.Println(reqGoodRequest.URL.String())
	if err != nil {
		t.Fatal(err)
	}

	rr = httptest.NewRecorder()
	handler = http.HandlerFunc(server.GetBoardGames)
	handler.ServeHTTP(rr, reqGoodRequest)
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	err = json.NewDecoder(rr.Body).Decode(&testResp)
	if err != nil {
		t.Fatal(err)
	}
	if len(testResp) > testLimit {
		t.Errorf("limit failed: expected <= %v, got %v", testLimit, len(testResp))
	}

	if testResp[0].Name != expectedName {
		t.Errorf("name check failed: expected %v, got %v", expectedName, testResp[0].Name)
	}

	if testResp[0].Id != expectedId {
		t.Errorf("id check failed: expected %v, got %v", expectedId, testResp[0].Id)
	}

}
