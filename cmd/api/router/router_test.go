package router

/*
func TestCORSMiddleware(t *testing.T) {
	router := NewRouter(":8080")
	go router.Setup()

	origin := "*"
	server := httptest.NewServer(router.GetRouter())
	defer server.Close()

	client := &http.Client{}

	// TESTING CORS WITH OPTIONS METHOD
	req, _ := http.NewRequest(
		"OPTIONS",
		"http://localhost:8080/ping",
		nil,
	)
	req.Header.Add("Origin", origin)

	get, err := client.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	o := get.Header.Get("Access-Control-Allow-Origin")
	if o != origin {
		t.Errorf("Got '%s' ; expecting origin '%s'", o, origin)
	}

	// TESTING CORS WITH GET METHOD
	req, _ = http.NewRequest(
		"GET",
		"http://localhost:8080/ping",
		nil,
	)
	req.Header.Add("Origin", origin)

	get, err = client.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	o = get.Header.Get("Access-Control-Allow-Origin")
	if o != origin {
		t.Errorf("Got '%s' ; expecting origin '%s'", o, origin)
	}
}
*/
