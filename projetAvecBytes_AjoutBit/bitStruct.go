package main

import "bytes"

// Bit writer
type BitWriter struct {
	buffer   byte   // octet en cours de construction
	bitCount uint8  // nombre de bits actuellement dans buffer
	output   []byte // tous les octets complets
}

func NewBitWriter() *BitWriter {
	return &BitWriter{
		buffer:   0,               // buffer de bits sur un octet
		bitCount: 0,               // indique combien de bits sont déjà dans le buffer
		output:   make([]byte, 0), // contient le résultat final.
	}
}

func (bw *BitWriter) WriteBits(value uint64, nbits int) {
	for i := nbits - 1; i >= 0; i-- {
		// prendre le i-ème bit du code à écrire
		bit := (value >> i) & 1
		bw.buffer <<= 1
		bw.buffer |= byte(bit)
		bw.bitCount++

		// Si on a un octet, on l'ajoute au résultat final et on commence un nouvel octet
		if bw.bitCount == 8 {
			bw.output = bytes.Clone(append(bw.output, bw.buffer))
			bw.buffer = 0
			bw.bitCount = 0
		}
	}
}

func (bw *BitWriter) Flush() {
	if bw.bitCount > 0 {
		// compléter avec des 0 à droite
		bw.buffer <<= (8 - bw.bitCount)
		bw.output = append(bw.output, bw.buffer)
	}
}

func (bw *BitWriter) Bytes() []byte {
	bw.Flush()
	return bw.output
}

// Bit reader
type BitReader struct {
	data     []byte
	position int // position en bits
}

func NewBitReader(data []byte) *BitReader {
	return &BitReader{data: data}
}

func (br *BitReader) ReadBits(nbits uint8) (uint16, bool) {
	var result uint16
	for i := uint8(0); i < nbits; i++ {
		bytePos := br.position / 8
		if bytePos >= len(br.data) {
			return 0, false
		}
		bitOffset := 7 - (br.position % 8)
		bit := (br.data[bytePos] >> bitOffset) & 1
		result = (result << 1) | uint16(bit)
		br.position++
	}
	return result, true
}
