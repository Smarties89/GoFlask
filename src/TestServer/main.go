package main

import "shttp"

func peter() string {
	return "You are awesome!"
}

func martin() string {
	return "Welcome Martin"
}

func main() {
	for {
		s := shttp.CreateHTTPServer(":5000")
		s.AddRoute("/peter", peter)
		s.AddRoute("/martin", martin)
		s.Run()
	}
}
