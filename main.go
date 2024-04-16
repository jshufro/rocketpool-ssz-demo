package main

import (
	"encoding/json"
	"flag"
	"io"
	"os"

	"github.com/rocket-pool/smartnode/shared/services/rewards/ssz_types"
)

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
	} else {

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
