// Package Client
/**
 * @author zeroc
 * @date 14:56 2023/5/22
 * @file main.go
 **/
package main

import (
	"Eth-PIR/Utils"
	"bufio"
	"fmt"
	"log"
	"math/big"
	"net"
	"os"
	"strconv"
	"strings"
)

const privatekey = "fdce6cf6e724e00a5bed9c6ef3be624a307d6bb3d502493e2f043522c1791cf2"

func main() {
	// using net.Dial to connect to the server
	conn, err := net.Dial("tcp", "127.0.0.1:20000")
	if err != nil {
		log.Fatalf("Failed to dial: %v", err)
	}
	log.Println("Connected successfully")
	// close the connection when the client exits
	defer func(conn net.Conn) {
		err := conn.Close()
		if err != nil {
			log.Fatalf("Failed to close connection: %v", err)
		}
	}(conn)
	inputReader := bufio.NewReader(os.Stdin)
	for {
		// 1337 is the chainID of the private chain
		chainID := big.NewInt(1337)
		// start the service process
		Utils.StartProcess(privatekey, chainID)
		fmt.Print("Please input the data you want to retrive from the server: ")
		input, _ := inputReader.ReadString('\n')
		inputInfo := strings.Trim(input, "\r\n")
		if strings.ToUpper(inputInfo) == "Q" {
			log.Println("Client exit")
			_, err := conn.Write([]byte(inputInfo))
			if err != nil {
				log.Fatalf("Failed to write: %v", err)
				return
			}
			return
		}
		_, err := conn.Write([]byte(inputInfo))
		if err != nil {
			log.Fatalf("Failed to write: %v", err)
			return
		}
		buf := [128]byte{}
		n, err := conn.Read(buf[:])
		if err != nil {
			log.Fatalf("Failed to read: %v", err)
			return
		}
		recvData := string(buf[:n])
		num1 := strings.Split(recvData, ",")[0]
		flag, _ := strconv.Atoi(num1)
		if flag == 1 {
			Utils.ClientConfirm(privatekey, chainID, true)
		} else {
			Utils.ClientConfirm(privatekey, chainID, false)
		}
	}
}
