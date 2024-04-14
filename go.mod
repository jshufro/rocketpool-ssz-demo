module github.com/jshufro/ssz-demo

go 1.21.8

require (
	github.com/go-test/deep v1.1.0
	github.com/klauspost/compress v1.17.6
	github.com/rocket-pool/smartnode v1.12.1
	github.com/stretchr/testify v1.8.1
)

require (
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/ferranbt/fastssz v0.1.3 // indirect
	github.com/holiman/uint256 v1.2.4 // indirect
	github.com/klauspost/cpuid/v2 v2.2.7 // indirect
	github.com/minio/sha256-simd v1.0.1 // indirect
	github.com/mitchellh/mapstructure v1.5.0 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	github.com/wealdtech/go-merkletree v1.0.1-0.20190605192610-2bb163c2ea2a // indirect
	golang.org/x/crypto v0.19.0 // indirect
	golang.org/x/sys v0.17.0 // indirect
	gopkg.in/yaml.v2 v2.4.0 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)

replace github.com/rocket-pool/smartnode => github.com/jshufro/smartnode v1.1.1-0.20240414065551-a92a1ccdfe36

replace github.com/wealdtech/go-merkletree v1.0.1-0.20190605192610-2bb163c2ea2a => github.com/rocket-pool/go-merkletree v1.0.1-0.20220406020931-c262d9b976dd
