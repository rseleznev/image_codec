package models

const (
	HeaderSize = 13

	MagicOffset = 0
	VersionOffset = 4
	WidthOffset = 5
	HeightOffset = 7
	DataSizeOffset = 9
)

type Pixel struct {
	R byte
	G byte
	B byte
}

type DeltaEncodedElement struct {
	R int16
	G int16
	B int16
}

type RLEEncodedElement struct {
	Count byte
	Value DeltaEncodedElement
}

type FileHeader struct {
	Magic [4]byte
	Version byte
	Width uint16
	Height uint16
	DataSize uint32
}

type HeapElement struct {
	Type string
	Value byte
	Freq int
	LeftChild *HeapElement
	RightChild *HeapElement
}

type HaffmanTreeUnit struct {
	TreeNode *HeapElement
	Code HaffmanCode
}

type HaffmanCode struct {
	BitCode uint32
	CodeLen uint32
}