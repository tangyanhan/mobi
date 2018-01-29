/*
Package mobi is used to generate mobi files
*/
package main

import (
	"encoding/binary"
	"fmt"
	"os"
	"flag"
)

// Compression options
const (
	CompressionNone      = 1
	CompressionPalmDOC   = 2
	CompressionHuffCDIC  = 17480
)

// PalmDOCHeader is the very beginning
type PalmDOCHeader struct {
	Compression int16
	Unused      int16  // Always zero
	TextLength  uint32  // Uncompressed length of the entire text of the book
	RecordCount int16  // Number of PDB records used for the text of the book.
	RecordSize  int16  // Maximum size of each record containing text, always 4096
	CurrentPosition uint32  // Current reading position, as an offset into the uncompressed text
}

// MobiType constants
const (
	MobiTypeMobiPocketBook = 2
	MobiTypePalmDocBook    = 3
	MobiTypeMobiGenByKindle = 232
	MobiTypeKF8ByKindle     = 248
	MobiTypeNews            = 257
	MobiTypePics            = 513 + iota
	MobiTypeWord
	MobiTypeXLS
	MobiTypePPT
	MobiTypeText
	MobiTypeHTML
)

// TextEncoding constants
const (
	TextEncodingCP1252 = 1252
	TextEncodingUTF8   = 65001
)

// IndexUnavailable for index unavailable in MobiHeader
const IndexUnavailable = 0xFFFFFFFF

// Constants for locale
const (
	LocaleCodeUSEnglish = 1033
	LocaleCodeUKEnglish = 2057
)
// Header is in record 0 following PalmDOCHeader
type Header struct {
	Identifier [4]byte  // characters MOBI
	HeaderLength uint32   // Length of MobiHeader
	MobiType  uint32
	TextEncoding uint32
	UniqueID    uint32   // some kind of unique ID number
	FileVersion uint32

	// Fields below are IndexUnavailable when unavailable
	OrtographicIndex uint32 // Section number of orthographic meta index
	InflectionIndex uint32
	IndexNames uint32
	IndexKeys uint32
	ExtraIndex [6]uint32

	FirstNonBookIndex uint32
	FullNameOffset    uint32
	FullNameLength    uint32
	LocaleCode        uint32
	InputLanguage     uint32
	OutputLanguage    uint32
	MinVersion        uint32
	FirstImageIndex   uint32  // First record number that contains image

	HuffmanRecordOffset uint32
	HuffmanRecordCount  uint32
	HuffmanTableOffset  uint32
	HuffmanTableLength  uint32
	EXTHFlags           uint32  // Bitfield,if bit 6 is set, then there's an EXTH record
	UnknownBytes        [36]byte  // Last 4 bytes as 0xFFFFFFFF
	DRMOffset uint32
	DRMCount  uint32
	DRMSize   uint32
	DRMFlags  uint32
}

var fileName = flag.String("filename", "", "File name of a mobi")
func main() {
	flag.Parse()

	die := func(msg string, err error) {
		fmt.Printf("%s %s  Error:%v\n", *fileName, msg, err)
		os.Exit(1)
	}

	if *fileName == "" {
		die( "No mobi file given, exit", nil)
	}

	fh, err := os.Open(*fileName)
	if err != nil {
		die("open", err)
	}

	defer func() {
		if err := fh.Close(); err != nil {
			die("closing", err)
		}
	}()

	palmHeader := PalmDOCHeader{}
	err = binary.Read(fh, binary.LittleEndian, &palmHeader)
	if err != nil {
		die("Reading PalmDOCHeader", err)
	}

	fmt.Println("PALMHeader:", palmHeader)
	mobiHeader := Header{}
	err = binary.Read(fh, binary.LittleEndian, &mobiHeader)
	if err != nil {
		die("Reading MOBI Header", err)
	}
	fmt.Println("MobiHeader:", mobiHeader, "Identifier:", string(mobiHeader.Identifier[:]))
}
