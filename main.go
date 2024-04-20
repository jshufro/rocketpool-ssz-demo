package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"os"

	"github.com/rocket-pool/smartnode/shared/services/rewards/ssz_types"
)

func printCid(filename string, data []byte) bool {
	if filename == "" {
		return false
	}

	cid, err := singleFileDirIPFSCid(data, filename)
	if err != nil {
		panic(err)
	}

	fmt.Println(cid)
	return true
}

func main() {
	input, output := os.Stdin, os.Stdout
	flag.Func("file", "path to file to encode/decode (stdin if not passed)", func(p string) (err error) {
		input, err = os.Open(p)
		return err
	})
	flag.Func("output", "path to write to (stdout if not passed)", func(p string) (err error) {
		output, err = os.OpenFile(p, os.O_CREATE|os.O_WRONLY, 0644)
		return err
	})
	encode := flag.Bool("encode", false, "If passed, input is treated as json and encoded to ssz instead of decoded to json")
	getCID := flag.String("cid", "", "If passed, input is treated as ssz and the IPFS CID corresponding to the ssz and filename is printed to stdout")
	defer input.Close()
	defer output.Close()
	flag.Parse()

	data, err := io.ReadAll(input)
	if err != nil {
		panic(err)
	}

	sszFile := ssz_types.NewSSZFile_v1()
	if *encode {
		// Refuse to write to a terminal
		if o, _ := output.Stat(); o.Mode()&os.ModeCharDevice != 0 {
			panic("Refusing to write SSZ to terminal")
		}

		err := json.Unmarshal(data, sszFile)
		if err != nil {
			panic(err)
		}
		data, err = sszFile.MarshalSSZ()
		if err != nil {
			panic(err)
		}
		if printCid(*getCID, data) {
			return
		}
	} else {
		if printCid(*getCID, data) {
			return
		}
		err = sszFile.UnmarshalSSZ(data)
		if err != nil {
			panic(err)
		}
		data, err = json.Marshal(sszFile)
		if err != nil {
			panic(err)
		}
	}

	_, err = output.Write(data)
	if err != nil {
		panic(err)
	}
	return
}
