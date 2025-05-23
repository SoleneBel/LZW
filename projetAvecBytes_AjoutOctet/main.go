package main

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func LZWCompresser(data []byte, dico map[string]int) []byte {
	var w []byte
	var result []byte
	//k := 1 // Pour gérer les puissances de 2 pour savoir sur combien de bits on encode/décode

	// TANT QUE (il reste des caractères à lire dans Texte) FAIRE
	for _, c := range data { // c ← Lire(Texte)
		p := append(w, c) // p ← Concaténer(w, c)

		// SI Existe(p, dictionnaire) ALORS
		if _, exists := dico[string(p)]; exists {
			w = p // w ← p
		} else {
			dico[string(p)] = len(dico) // Ajouter(p, dictionnaire)

			// code est un entier (int) représentant l'index de w dans le dictionnaire
			// Pour stocker ce code en 2 octets, il est divisé en deux parties :
			// code >> 8 : récupère les 8 bits de poids fort
			// code & 0xFF : récupère les 8 bits de poids faible
			// Ensuite, on convertit chaque partie en byte pour qu'ils tiennent sur 8 bits, et on les ajoute à la slice result
			code := dico[string(w)]
			result = append(result, byte(code>>8), byte(code&0xFF)) // Écrire dictionnaire[w]
			w = []byte{c}                                           // w ← c
		}
	}

	if len(w) > 0 {
		code := dico[string(w)]
		result = append(result, byte(code>>8), byte(code&0xFF))
	}

	return result
}

func LZWDecompresser(data []byte, dico map[int][]byte) []byte {
	var result []byte

	if len(data) < 2 {
		return result
	}

	// récupère sous forme de int le v en prenant ses bits de poids fort et bits de poids faible
	readCode := func(i int) int {
		return int(data[i])<<8 | int(data[i+1])
	}

	v := readCode(0) // v ← Lire(Code)
	w := dico[v]
	// append chaque élément de w à result ([]byte)
	result = append(result, w...) // Écrire dictionnaire[v]

	for i := 2; i < len(data); i += 2 {
		v = readCode(i) // v ← Lire(Code)

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

func sameFiles(filename1, filename2 string) bool {
	file1, _ := os.ReadFile(filename1)
	file2, _ := os.ReadFile(filename2)

	return bytes.Equal(file1, file2)
}

func main() {
	/*testFile := "input.txt"
	data, _ := os.ReadFile(testFile)*/
	if len(os.Args) < 2 {
		fmt.Println("Usage: go run main.go <nom_du_fichier>")
		return
	}

	testFile := os.Args[1] // On récupère le nom du fichier passé en argument
	data, err := os.ReadFile(testFile)
	if err != nil {
		fmt.Println("Erreur de lecture du fichier :", err)
		return
	}

	// ENCODAGE
	dicoC := make(map[string]int)
	for i := 0; i < 256; i++ {
		dicoC[string([]byte{byte(i)})] = i
	}
	compressed := LZWCompresser(data, dicoC)

	// Générer le nom du fichier compressé : nom + "_compresse" + extension
	ext := filepath.Ext(testFile)                            // extension
	base := strings.TrimSuffix(filepath.Base(testFile), ext) // nom
	nomFichierCompresse := base + "_compresse" + ext

	err1 := ecrireDansFichier(nomFichierCompresse, compressed)
	if err1 != nil {
		return
	}

	// DECODAGE
	toDecod, _ := lireBytesDepuisFichier(nomFichierCompresse)

	dicoD := make(map[int][]byte)
	for i := 0; i < 256; i++ {
		dicoD[i] = []byte{byte(i)}
	}
	decompressed := LZWDecompresser(toDecod, dicoD)

	nomFichierDecompresse := base + "_decompresse" + ext
	err2 := ecrireDansFichier(nomFichierDecompresse, decompressed)
	if err2 != nil {
		return
	}

	same := sameFiles(testFile, nomFichierDecompresse)
	fmt.Println(same)

	// VERSION SANS LECTURE/ECRITURE DANS FICHIERS
	/*data, _ := os.ReadFile("input.txt")

	dicoC := make(map[string]int)
	for i := 0; i < 256; i++ {
		dicoC[string([]byte{byte(i)})] = i
	}
	compressed := LZWCompresser(data, dicoC)

	dicoD := make(map[int][]byte)
	for i := 0; i < 256; i++ {
		dicoD[i] = []byte{byte(i)}
	}
	decompressed := LZWDecompresser(compressed, dicoD)

	fmt.Println("Décompressé :")
	fmt.Println(string(decompressed))*/

}

// pour ajouter bit à bit
// buffer qui va contenir un octet mais icomplet =>
// 		// struct avec tab d'octets, un buffer de type byte et une méthode qui permet de faire append d'un certain nb de bit

// regarder si ça vaut le coup de remettre les dictionnaires à 0
