package main

import (
	"bufio"
	"os"
)

func lireTexteDepuisFichier(fileName string) string {
	var f, _ = os.Open(fileName)
	sc := bufio.NewScanner(f)

	var lines string
	for sc.Scan() {
		line := sc.Text()
		if line != "" {
			lines += line + "\n"
		}
	}

	return lines
}

func lireBytesDepuisFichier(nomFichier string) ([]byte, error) {
	// Lire le fichier entier
	data, err := os.ReadFile(nomFichier)
	if err != nil {
		return nil, err
	}
	return data, nil
}

// Écrit une entrée quelconque dans un fichier
func ecrireDansFichier(nomFichier string, input []byte) error {
	// Ouvrir ou créer le fichier en mode écriture
	fichier, err := os.Create(nomFichier)
	if err != nil {
		return err
	}
	defer fichier.Close()

	// Convertir l'entrée en une chaîne de caractères
	_, err = fichier.Write(input)
	if err != nil {
		return err
	}

	return nil
}
