package main

import (
	"bytes"
	"fmt"
	"math"
	"math/bits"
	"os"
	"path/filepath"
	"strings"
	"time"
)

func lenInt(n int) int {
	if n == 0 {
		return 1
	}
	bitLength := bits.Len(uint(n)) // !
	chunks := (bitLength + 7) / 8
	return chunks
}

// splitInt permet de séparer un entier n en k octets
func splitInt(n int, s int) []int {
	if n == 0 {
		result := make([]int, s)
		return result
	}

	listeSplited := []int{}
	for n > 0 {
		part := int(n & 0xFF)
		listeSplited = append(listeSplited, part)
		n >>= 8
	}

	for len(listeSplited) < s {
		listeSplited = append(listeSplited, 0)
	}

	return listeSplited
}

// packInt permet de reformer une liste d'octet en un entier
func packInt(chunks []byte) int {
	result := 0
	for i := len(chunks) - 1; i >= 0; i-- {
		result = (result << 8) | int(chunks[i])
	}
	return result
}

// LZWCompresser fait la compression en suivant LZW
func LZWCompresser(data []byte, dico map[string]int) []byte {
	var w []byte
	var result []byte
	k := 1 // Pour gérer les puissances de 2 pour savoir sur combien de bits on encode/décode

	// TANT QUE (il reste des caractères à lire dans Texte) FAIRE
	for _, c := range data { // c ← Lire(Texte)
		p := append(w, c) // p ← Concaténer(w, c)

		// SI Existe(p, dictionnaire) ALORS
		if _, exists := dico[string(p)]; exists {
			w = p // w ← p
		} else {
			code := dico[string(w)]
			if int(math.Pow(2, float64(8*k))) <= len(dico) {
				k++
			}

			bytesCode := splitInt(code, k)
			for _, b := range bytesCode {
				result = append(result, byte(b)) // Écrire dictionnaire[w]
			}

			dico[string(p)] = len(dico) // Ajouter(p, dictionnaire)

			w = []byte{c}
		}
	}

	if len(w) > 0 {
		code := dico[string(w)]
		if int(math.Pow(2, float64(8*k))) <= len(dico) {
			k++
		}

		bytesCode := splitInt(code, k)
		for _, b := range bytesCode {
			result = append(result, byte(b))
		}
	}

	return result
}

// LZWDecompresser fait la compression en suivant LZW
func LZWDecompresser(data []byte, dico map[int][]byte) []byte {
	var result []byte
	k := 1
	i := 0

	if int(math.Pow(2, float64(8*k))) <= len(dico) {
		k++
	}

	if len(data) < k {
		return result
	}

	v := packInt(data[i : i+k])
	i += k

	w := dico[v]
	// append chaque élément de w à result ([]byte)
	result = append(result, w...) // Écrire dictionnaire[v]

	for i < len(data) {
		if (1 << (8 * k)) <= len(dico)+1 {
			k++
		}

		v := packInt(data[i : i+k]) // v ← Lire(Code)
		i += k

		var entry []byte
		// SI Existe(v, dictionnaire) ALORS
		if val, exists := dico[v]; exists {
			entry = val // entrée ← dictionnaire[v]
		} else {
			entry = append(w, w[0]) // entrée ← Concaténer(w, w[0])
		}

		result = append(result, entry...) // Écrire entrée
		// clone fait une copie pour pas qu'on écrive au-delà de ce qu'il y a dans la slice
		dico[len(dico)] = bytes.Clone(append(w, entry[0])) // Ajouter(Concaténer(w,entrée[0]), dictionnaire)
		w = entry                                          // w ← entrée

	}

	return result
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: go run main.go <nom_du_fichier>")
		return
	}

	testFile := os.Args[1]
	fmt.Println("### Lempel / Ziv / Welch compression\n")

	// ----------- COMPRESSION -----------
	fmt.Printf("## compressing file '%s'\n", testFile)

	// -- Reading file --
	start := time.Now()
	data, err := os.ReadFile(testFile)
	if err != nil {
		fmt.Println("Erreur de lecture du fichier :", err)
		return
	}
	readDuration := time.Since(start)
	fmt.Printf(">>> reading file:   (%6dms)  => %d bytes\n", readDuration.Milliseconds(), len(data))

	// -- Compression --
	dicoC := make(map[string]int)
	for i := 0; i < 256; i++ {
		dicoC[string([]byte{byte(i)})] = i
	}

	start = time.Now()
	compressed := LZWCompresser(data, dicoC)
	compDuration := time.Since(start)

	ratio := 100 - int(float64(len(compressed))*100.0/float64(len(data)))
	fmt.Printf(">>> compressing:    (%6dms)  => %d bytes, compression %d%%\n", compDuration.Milliseconds(), len(compressed), ratio)

	// -- Saving file --
	// Nom compressé
	ext := filepath.Ext(testFile)
	base := strings.TrimSuffix(filepath.Base(testFile), ext)
	compressedFile := base + "_compresse" + ext

	start = time.Now()
	err = ecrireDansFichier(compressedFile, compressed)
	saveDuration := time.Since(start)
	if err != nil {
		fmt.Println("Erreur d'écriture (compression) :", err)
		return
	}
	fmt.Printf(">>> saving:         (%6dms)\n", saveDuration.Milliseconds())
	fmt.Printf("## compressed file '%s' saved\n\n", compressedFile)

	// ----------- DECOMPRESSION -----------
	fmt.Printf("## decompressing file '%s'\n", compressedFile)

	// -- Read file --
	start = time.Now()
	toDecod, err := lireBytesDepuisFichier(compressedFile)
	if err != nil {
		fmt.Println("Erreur de lecture (compression) :", err)
		return
	}
	readDecDuration := time.Since(start)
	fmt.Printf("<<< read file:      (%6dms)  => %d bytes\n", readDecDuration.Milliseconds(), len(toDecod))

	// -- Decompression --
	dicoD := make(map[int][]byte)
	for i := 0; i < 256; i++ {
		dicoD[i] = []byte{byte(i)}
	}

	start = time.Now()
	decompressed := LZWDecompresser(toDecod, dicoD)
	decompDuration := time.Since(start)
	fmt.Printf("<<< decompressing:  (%6dms)  => %d bytes\n", decompDuration.Milliseconds(), len(decompressed))

	// -- Saving file --
	decompressedFile := base + "_decompresse" + ext
	start = time.Now()
	err = ecrireDansFichier(decompressedFile, decompressed)
	writeDecDuration := time.Since(start)

	if err != nil {
		fmt.Println("Erreur d'écriture (décompression) :", err)
		return
	}

	fmt.Printf("<<< saving:         (%6dms)\n", writeDecDuration.Milliseconds())
	fmt.Printf("## decompressed file '%s' saved\n\n", decompressedFile)

	// ----------- COMPARAISON -----------
	start = time.Now()
	same := sameFiles(testFile, decompressedFile)
	compareDuration := time.Since(start)

	status := "FAIL"
	if same {
		status = "OK"
	}

	fmt.Printf("## comparing: %s    (%6dms)\n", status, compareDuration.Milliseconds())
}

// pour ajouter bit à bit
// buffer qui va contenir un octet mais icomplet =>
// 		// struct avec tab d'octets, un buffer de type byte et une méthode qui permet de faire append d'un certain nb de bit

// regarder si ça vaut le coup de remettre les dictionnaires à 0
