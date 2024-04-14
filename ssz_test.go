package main

import (
	"bytes"
	_ "embed"
	"encoding/json"
	"testing"

	"github.com/go-test/deep"
	"github.com/klauspost/compress/zstd"
	"github.com/rocket-pool/smartnode/shared/services/rewards/ssz_types"
	"github.com/stretchr/testify/require"
)

//go:embed rp-rewards-holesky-155.json
var testFile string

func TestRoundTrip(t *testing.T) {
	// ssz format should be able to unmarshal identical json to the old format
	sszFile := ssz_types.NewSSZFile_v1()

	err := json.Unmarshal([]byte(testFile), &sszFile)
	if err != nil {
		t.Fatal(err)
	}

	// Check the header for equality
	if !bytes.Equal(sszFile.Magic[:], ssz_types.Magic[:]) {
		t.Fatalf("invalid header %v", sszFile.Magic)
	}

	// encode the result to ssz
	sszData, err := sszFile.MarshalSSZ()
	if err != nil {
		t.Fatal(err)
	}

	// We should be able to unmarshal it back into a new file
	sszFile2 := ssz_types.NewSSZFile_v1()
	err = sszFile2.UnmarshalSSZ(sszData)
	if err != nil {
		t.Fatal(err)
	}

	// generate sszFile2 proofs
	_, err = sszFile2.Proofs()
	if err != nil {
		t.Fatal(err)
	}

	// The two files should be deep-equal
	diffs := deep.Equal(sszFile2, sszFile)
	if len(diffs) > 0 {
		for _, d := range diffs {
			t.Log(d)
		}
		t.Fatal("Round trip parsed ssz differs from serialized ssz")
	}

	// Finally, serialize json and compare to the original
	jsonData, err := json.Marshal(sszFile)
	if err != nil {
		t.Fatal(err)
	}

	// The original file has a minipoolPerformanceFileCid string of --- that we must manually add
	var tmp map[string]any
	err = json.Unmarshal(jsonData, &tmp)
	if err != nil {
		t.Fatal(err)
	}
	tmp["minipoolPerformanceFileCid"] = "---"
	jsonData, err = json.Marshal(tmp)
	if err != nil {
		t.Fatal(err)
	}

	require.JSONEq(t, testFile, string(jsonData))
}

func TestInfo(t *testing.T) {
	sszFile := ssz_types.NewSSZFile_v1()
	err := json.Unmarshal([]byte(testFile), &sszFile)
	if err != nil {
		t.Fatal(err)
	}

	t.Logf("json-uncompressed:\t%d bytes\n", len([]byte(testFile)))

	sszData, err := sszFile.MarshalSSZ()
	if err != nil {
		t.Fatal(err)
	}

	t.Logf("ssz-uncompressed:\t%d bytes\n", len([]byte(sszData)))

	//compress json and ssz
	encoder, _ := zstd.NewWriter(nil, zstd.WithEncoderLevel(zstd.SpeedBestCompression))
	jsonCompressed := encoder.EncodeAll([]byte(testFile), make([]byte, 0, len(testFile)))
	sszCompressed := encoder.EncodeAll(sszData, make([]byte, 0, len(sszData)))
	t.Logf("json-compressed:\t%d bytes\n", len(jsonCompressed))
	t.Logf("ssz-compressed:\t%d bytes\n", len(sszCompressed))

}

func BenchmarkJSONParsing(b *testing.B) {
	for i := 0; i < b.N; i++ {
		sszFile := ssz_types.NewSSZFile_v1()

		err := json.Unmarshal([]byte(testFile), &sszFile)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkSSZParsing(b *testing.B) {
	sszFile := ssz_types.NewSSZFile_v1()

	err := json.Unmarshal([]byte(testFile), &sszFile)
	if err != nil {
		b.Fatal(err)
	}
	sszData, err := sszFile.MarshalSSZ()
	if err != nil {
		b.Fatal(err)
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		sszFile2 := ssz_types.NewSSZFile_v1()
		err = sszFile2.UnmarshalSSZ(sszData)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkSSZParsingAndProofGeneration(b *testing.B) {
	sszFile := ssz_types.NewSSZFile_v1()

	err := json.Unmarshal([]byte(testFile), &sszFile)
	if err != nil {
		b.Fatal(err)
	}
	sszData, err := sszFile.MarshalSSZ()
	if err != nil {
		b.Fatal(err)
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		sszFile2 := ssz_types.NewSSZFile_v1()
		err = sszFile2.UnmarshalSSZ(sszData)
		if err != nil {
			b.Fatal(err)
		}
		sszFile2.Proofs()
	}
}

func BenchmarkJSONSerialization(b *testing.B) {
	sszFile := ssz_types.NewSSZFile_v1()

	err := json.Unmarshal([]byte(testFile), &sszFile)
	if err != nil {
		b.Fatal(err)
	}
	sszFile.Proofs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err = json.Marshal(sszFile)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkSSZSerialization(b *testing.B) {
	sszFile := ssz_types.NewSSZFile_v1()

	err := json.Unmarshal([]byte(testFile), &sszFile)
	if err != nil {
		b.Fatal(err)
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err = sszFile.MarshalSSZ()
		if err != nil {
			b.Fatal(err)
		}
	}
}
