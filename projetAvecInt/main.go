package main

import (
	"fmt"
)

/* VERSION AVEC DES INT */

// LZWCompresser compresse du texte en utilisant LZW
func LZWCompresser(texte string, dictionnaire map[string]int) []int {
	var w string
	var result []int

	// TANT QUE (il reste des caractères à lire dans Texte) FAIRE
	for i := 0; i < len(texte); i++ {
		c := string(texte[i]) // c ← Lire(Texte)
		p := w + c            // p ← Concaténer(w, c)

		// SI Existe(p, dictionnaire) ALORS
		if _, exists := dictionnaire[p]; exists {
			w = p // w ← p
		} else {
			dictionnaire[p] = len(dictionnaire)      // Ajouter(p, dictionnaire)
			result = append(result, dictionnaire[w]) // Écrire dictionnaire[w]
			w = c                                    // w ← c
		}
	}

	// Écrire dictionnaire[w]
	if w != "" {
		result = append(result, dictionnaire[w])
	}

	return result
}

// LZWDecompresser décompresse un texte en utilisant LZW
func LZWDecompresser(code []int, dictionnaire map[int]string) string {
	var result string

	n := len(code)            // n ← |Code|
	v := code[0]              // v ← Lire(Code)
	result += dictionnaire[v] // Écrire dictionnaire[v]
	w := dictionnaire[v]      // w ← chr(v)

	// POUR i ALLANT DE 2 à n FAIRE
	for i := 1; i < n; i++ {
		var entree string
		v = code[i] // v ← Lire(Code)

		// SI Existe(v, dictionnaire) ALORS
		if val, exists := dictionnaire[v]; exists {
			entree = val // entrée ← dictionnaire[v]
		} else {
			entree = w + string(w[0]) // entrée ← Concaténer(w, w[0])
		}

		// Écrire entrée
		result += entree

		// Ajouter(Concaténer(w,entrée[0]), dictionnaire)
		dictionnaire[len(dictionnaire)] = w + string(entree[0])

		// w ← entrée
		w = entree
	}

	return result
}

func main() {
	dicoC := make(map[string]int)
	for i := 0; i < 256; i++ {
		dicoC[string(byte(i))] = i
	}
	texte := lireTexteDepuisFichier("input.txt")
	code := LZWCompresser(texte, dicoC)

	dicoD := make(map[int]string)
	for i := 0; i < 256; i++ {
		dicoD[i] = string(byte(i))
	}
	texteDecompresse := LZWDecompresser(code, dicoD)
	fmt.Println(texteDecompresse)
}
