package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"runtime"
	"strings"
)

var clear map[string]func() //create a map for storing clear funcs
var execution bool = true

func init() {
	clear = make(map[string]func()) //Initialize it
	clear["linux"] = func() {
		cmd := exec.Command("clear") //Linux example, its tested
		cmd.Stdout = os.Stdout
		cmd.Run()
	}
	clear["windows"] = func() {
		cmd := exec.Command("cmd", "/c", "cls") //Windows example, its tested
		cmd.Stdout = os.Stdout
		cmd.Run()
	}
}

func CallClear() {
	value, ok := clear[runtime.GOOS] //runtime.GOOS -> linux, windows, darwin etc.
	if ok {                          //if we defined a clear func for that platform:
		value() //we execute it
	} else { //unsupported platform
		panic("Your platform is unsupported! I can't clear terminal screen :(")
	}
}

func main() {

	fmt.Println("Welcome to a new Web Browser in your terminal")
	for execution {
		var input string
		fmt.Scanln(&input)
		CallClear()
		rawbody, err := navegar(input)
		if err != nil {
			fmt.Println(err)
		}
		if err == nil {
			fmt.Println(rawbody)
		}
	}
}

func navegar(url string) (string, error) {
	if url == "exit" {
		exit()
		return "Se abandona el navegador", nil
	}
	resp, err := http.Get(url)
	if err != nil {
		return "", err
	}
	rawBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	body, err := getBodyTagContent(string(rawBody))
	if err != nil {
		return "", err
	}
	return body, nil
}

func exit() {
	execution = false
}

func getBodyTagContent(rawBody string) (string, error) {
	_, after, found := strings.Cut(rawBody, "<body")
	if !found {
		return "", fmt.Errorf("No se encontro el tag body")
	}
	_, afterClose, found := strings.Cut(after, ">")
	if !found {
		return "", fmt.Errorf("No se encontro el tag body")
	}
	bodyContent, _, found := strings.Cut(afterClose, "</body")

	return bodyContent, nil
}
