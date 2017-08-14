package proxymap

func ExampleApp() {
	rpm, err := ToReverseProxyMap(map[string]string{
		"example.com": "http://getdown.usersys.redhat.com:8888/",
		"farts.com":   "http://localhost:8888/",
	})
	if err != nil {
		logrus.Fatal(err)
	}
	rph := NewReverseProxiesHandler(rpm)
	s := &http.Server{
		Addr:           ":http",
		Handler:        rph,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}
	logrus.Printf("listening on %q", s.Addr)
	logrus.Fatal(s.ListenAndServe())
}
