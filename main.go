// Package Eth_PIR
/**
 * @author zeroc
 * @date 1:15 2023/5/21
 * @file main.go
 **/
package main

import pir2 "Eth-PIR/pir"

const LOGQ = uint64(32)
const SEC_PARAM = uint64(1 << 10)

func main() {
	N := uint64(1 << 28)
	d := uint64(3)
	pir := pir2.DoublePIR{}
	p := pir.PickParams(N, d, SEC_PARAM, LOGQ)

	DB := pir2.MakeRandomDB(N, d, &p)
	pir2.RunPIR(&pir, DB, p, []uint64{0})
}
