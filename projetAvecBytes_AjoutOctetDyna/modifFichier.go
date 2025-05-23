package main

import (
	"bufio"
	"bytes"
	"os"
)

// lireTexteDepuisFichier permet de lire le fichier sous forme de string
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

// lireTexteDepuisFichier permet de lire le fichier sous forme de tableau de bytes
func lireBytesDepuisFichier(nomFichier string) ([]byte, error) {
	// Lire le fichier entier
	data, err := os.ReadFile(nomFichier)
	if err != nil {
		return nil, err
	}
	return data, nil
}

// ecrireDansFichier écrit une entrée sous forme de tableau de bytes dans un fichier
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

// sameFiles vérifie que 2 fichiers sont semblables
func sameFiles(filename1, filename2 string) bool {
	file1, _ := os.ReadFile(filename1)
	file2, _ := os.ReadFile(filename2)

	return bytes.Equal(file1, file2)
}
