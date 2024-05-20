package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"requester/colors"
	"strconv"
)

type Configuration struct {
	Ip   string
	Port string

	Method   string
	Path     string
	Protocol string
	Url      string
	Headers  []string
	Body     string
}

var config Configuration

func loadConfig() {
	colors.Cyan()
	fmt.Printf("    Reading from config.json\n")
	colors.Reset()
	// Open config.json
	jsonFile, err := os.Open("config.json")
	if err != nil {
		colors.RedBold()
		fmt.Printf("\n    Error opening config.json.\n")
		colors.Reset()
		fmt.Println(err)
	}
	defer jsonFile.Close()
	// Read from config.json file.
	jsonFileData, err := io.ReadAll(jsonFile)
	if err != nil {
		colors.RedBold()
		fmt.Printf("\n    Error reading from jsonFileData.\n")
		colors.Reset()
		fmt.Println(err)
	}
	// Copy config.json data to &config.
	err = json.Unmarshal(jsonFileData, &config)
	if err != nil {
		colors.RedBold()
		fmt.Println(err)
		colors.Reset()
	}
}

func urlInfo(url string) {
	colors.Cyan()
	fmt.Printf("    Gathering info for url: %v\n\n", url)
	colors.Reset()
	// Find http port for URL.
	port, err := net.LookupPort("tcp", "http")
	if err != nil {
		colors.RedBold()
		fmt.Println(err)
		colors.Reset()
	}
	fmt.Printf("Port: %v\n", port)
	config.Port = strconv.Itoa(port)
	// Returns ip of URL.
	// Builds URL:port.
	ip, err := net.ResolveTCPAddr("tcp", url+":"+config.Port)
	if err != nil {
		colors.RedBold()
		fmt.Println(err)
		colors.Reset()
	}
	fmt.Printf("IP: %v\n", ip)
	config.Ip = ip.IP.String()
	// Lookup other names for given ip.
	otherAddr, err := net.LookupAddr(config.Ip)
	if err != nil {
		colors.RedBold()
		fmt.Println(err)
		colors.Reset()
	}
	fmt.Printf("Other names for ip addr: %v\n\n", otherAddr)
}

func dial() {
	// Set ip.
	ip, err := net.ResolveTCPAddr("tcp", config.Ip+":"+config.Port)
	if err != nil {
		panic(err)
	}
	// Connect to TCP socket.
	conn, err := net.DialTCP("tcp", nil, ip)
	if err != nil {
		panic(err)
	}
	defer conn.Close()
	// Format request properly.
	// Redundant whitespace can break this.
	request := config.Method + " " + config.Path + " " + config.Protocol + "\r\n"
	request += "Host: " + config.Url + "\r\n"
	for _, header := range config.Headers {
		request += header + "\r\n"
	}
	request += "" + "\r\n"
	// Write request to opened connection.
	conn.Write([]byte(request))
	colors.GreenBold()
	fmt.Printf("Your request:\n")
	colors.Reset()
	log.Printf("sent request:\n\n%s", request)
	// Read response.
	colors.GreenBold()
	fmt.Printf("Response:\n")
	colors.Reset()
	for scanner := bufio.NewScanner(conn); scanner.Scan(); {
		line := scanner.Bytes()
		if _, err := fmt.Fprintf(os.Stdout, "%s\n", line); err != nil {
			log.Printf("error writing to connection: %s", err)
		}
		if scanner.Err() != nil {
			log.Printf("error reading from connection: %s", err)
			return
		}
	}
}

func main() {
	colors.ClearTerminal()
	loadConfig()
	urlInfo(config.Url)
	dial()
}
