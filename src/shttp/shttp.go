package shttp

import (
	"container/list"
	"log"
	"net"
	"strconv"
	"strings"
)

type SHTTPServer struct {
	l net.Listener
	routes *list.List
}


func extractrequest(inputstring string) (string, string) {
	splitted := strings.Split(inputstring, " ")
	return splitted[0],  splitted[1]
}

func findroutefun(s SHTTPServer, route string) Route {
	for r := s.routes.Front(); r != nil; r = r.Next() {

		if route == r.Value.(RouteElement).route  {
			return r.Value.(RouteElement).fun

		}
	}
	return nil
}

func handleconnection(s SHTTPServer, c net.Conn) {
	var input []byte = make([]byte, 1024)
	i, err := c.Read(input)
	_ = err

	inputstring := string(input[:i])
	method, loc := extractrequest(inputstring)
	_ = method
	routefun := findroutefun(s, loc)
	if routefun != nil {
		c.Write(createheader("200", routefun() ))
	} else {
		c.Write(createheader("404", "not found!"))
	}

	c.Close()
}

func (s SHTTPServer)Run() {
	for {
		conn, err := s.l.Accept()
		if err != nil {
			log.Fatal(err)
		} else {
			go handleconnection(s, conn)
		}
	}

}

func createheader(statuscode string, content string) []byte {
	head := "HTTP/1.1 200 ok\n"
	head += "Content-Type: text/html; charset=utf-8\n"
	head += "Content-Length: " + strconv.Itoa(len(content)) + "\n"
	head += "\n"
	head += content
	return []byte(head)
}

func CreateHTTPServer(port string) SHTTPServer  {
	l, err := net.Listen("tcp", port)

	if err != nil {
		log.Fatal(err)
	}
	return SHTTPServer{l, list.New() }
}

type RouteElement struct {
	route string
	fun Route
}

type Route func() string

func (s SHTTPServer)AddRoute(url string, f Route) {
	s.routes.PushBack(RouteElement{url, f})
}
