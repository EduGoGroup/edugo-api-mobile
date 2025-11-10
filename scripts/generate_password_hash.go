package main

import (
	"fmt"
	"log"

	"github.com/EduGoGroup/edugo-shared/auth"
)

func main() {
	password := "password123"

	fmt.Println("Generando hash bcrypt para:", password)
	fmt.Println("==================================================")

	hash, err := auth.HashPassword(password)
	if err != nil {
		log.Fatalf("Error generando hash: %v", err)
	}

	fmt.Println("Hash generado:")
	fmt.Println(hash)
	fmt.Println("")

	// Verificar que el hash funciona
	err = auth.VerifyPassword(hash, password)
	if err != nil {
		log.Fatalf("Error verificando hash: %v", err)
	}

	fmt.Println("✅ Hash verificado exitosamente")
	fmt.Println("")
	fmt.Println("Puedes usar este hash en tu archivo 02_seed_data.sql")
}
