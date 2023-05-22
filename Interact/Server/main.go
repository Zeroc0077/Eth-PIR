// Package Server
/**
 * @author zeroc
 * @date 14:56 2023/5/22
 * @file main.go
 **/
package main

import (
	"Eth-PIR/Utils"
	pir2 "Eth-PIR/pir"
	"bufio"
	"log"
	"math/big"
	"net"
	"strconv"
	"strings"
)

const LOGQ = uint64(32)
const SEC_PARAM = uint64(1 << 10)
const privatekey = "582583d82e41eb96922349fd0d2fb4405f7b6282b3364beb84553d360204179d"

func boolToBinaryString(value bool) string {
	if value {
		return "1"
	} else {
		return "0"
	}
}

func process(conn net.Conn) {
	defer func(conn net.Conn) {
		err := conn.Close()
		if err != nil {
			log.Fatalf("Failed to close connection: %v", err)
		}
	}(conn)
	for {
		// 1337 is the chainID of the private chain
		chainID := big.NewInt(1337)
		// charge the server
		Utils.ChargeServer(privatekey, chainID)
		reader := bufio.NewReader(conn)
		var buf [128]byte
		n, err := reader.Read(buf[:])
		if err != nil {
			log.Fatalf("Failed to read data: %v", err)
		}
		recvData := string(buf[:n])
		num1, num2 := strings.Split(recvData, ",")[0], strings.Split(recvData, ",")[1]
		s1, _ := strconv.Atoi(num1)
		s2, _ := strconv.Atoi(num2)
		log.Printf("N: %v, d: %v", s1, s2)
		// Run the PIR protocol
		N := uint64(1 << s1)
		d := uint64(s2)
		pir := pir2.DoublePIR{}
		p := pir.PickParams(N, d, SEC_PARAM, LOGQ)

		DB := pir2.MakeRandomDB(N, d, &p)
		_, _, flag, val := pir2.RunPIR(&pir, DB, p, []uint64{0})

		sendData := boolToBinaryString(flag) + "," + strconv.FormatUint(val, 10)
		_, err = conn.Write([]byte(sendData))
		if err != nil {
			return
		}
	}
}

func main() {
	listen, err := net.Listen("tcp", "127.0.0.1:20000")
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}
	log.Println("Listening on the Port 20000")
	for {
		conn, err := listen.Accept()
		log.Printf("Accept a connection from %v", conn.RemoteAddr().String())
		if err != nil {
			log.Fatalf("Failed to accept: %v", err)
		}
		go process(conn)
	}
}
