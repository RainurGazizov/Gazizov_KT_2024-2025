package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	"time"
	"unicode"
)

type Rotor struct {
	wiring   string
	position int
}

type Reflector struct {
	wiring string
}

type Enigma struct {
	rotors    []Rotor
	reflector Reflector
}

func NewRotor(wiring string, position int) Rotor {
	return Rotor{wiring: wiring, position: position}
}

func NewReflector(wiring string) Reflector {
	return Reflector{wiring: wiring}
}

func NewEnigma(rotors []Rotor, reflector Reflector) Enigma {
	return Enigma{rotors: rotors, reflector: reflector}
}

// rotateRotors вращает роторы после каждого символа
func (e *Enigma) rotateRotors() {
	// Вращаем первый ротор (самый правый)
	e.rotors[0].position = (e.rotors[0].position + 1) % 26

	// Проверяем, завершил ли первый ротор полный оборот
	if e.rotors[0].position == 0 {
		// Вращаем второй ротор
		e.rotors[1].position = (e.rotors[1].position + 1) % 26

		// Проверяем, завершил ли второй ротор полный оборот
		if e.rotors[1].position == 0 {
			// Вращаем третий ротор
			e.rotors[2].position = (e.rotors[2].position + 1) % 26
		}
	}
}

// encryptLetter шифрует одну букву
func (e *Enigma) encryptLetter(letter rune) rune {
	// Преобразуем букву в индекс (A=0, B=1, ..., Z=25)
	// Если буква строчная, преобразуем ее в заглавную
	index := int(unicode.ToUpper(letter) - 'A')

	// Проход через роторы вперед
	for i := range e.rotors {
		index = (index + e.rotors[i].position) % 26
		index = int(e.rotors[i].wiring[index] - 'A')
		index = (index - e.rotors[i].position + 26) % 26
	}

	// Проход через отражатель
	index = int(e.reflector.wiring[index] - 'A')

	// Проход через роторы в обратном порядке
	for i := len(e.rotors) - 1; i >= 0; i-- {
		index = (index + e.rotors[i].position) % 26
		index = strings.IndexRune(e.rotors[i].wiring, rune(index+'A'))
		index = (index - e.rotors[i].position + 26) % 26
	}

	// Вращаем роторы после шифрования буквы
	e.rotateRotors()

	// Преобразуем индекс обратно в букву
	encryptedLetter := rune(index + 'A')

	// Если исходная буква была строчной, возвращаем строчную букву
	if unicode.IsLower(letter) {
		return unicode.ToLower(encryptedLetter)
	}
	// Возвращаем заглавную букву
	return encryptedLetter
}

// Encrypt шифрует текст
func (e *Enigma) Encrypt(text string) string {
	var result strings.Builder
	for _, letter := range text {
		// Если символ - буква, шифруем ее
		if unicode.IsLetter(letter) {
			encryptedLetter := e.encryptLetter(letter)
			result.WriteRune(encryptedLetter)
		} else {
			// Если символ не является буквой (знак препинания, пробел и т.д.), добавляем его без изменений
			result.WriteRune(letter)
		}
	}
	return result.String() // Возвращаем зашифрованный текст
}

// writeToFile записывает текст в файл
func writeToFile(filename, text string) error {
	file, err := os.Create(filename) // Создаем файл
	if err != nil {
		return err
	}
	defer file.Close() // Закрываем файл после завершения

	_, err = file.WriteString(text) // Записываем текст в файл
	if err != nil {
		return err
	}

	return nil
}

func main() {
	startTime := time.Now()

	// Определяем роторы и отражатель
	rotor1 := NewRotor("EKMFLGDQVZNTOWYHXUSPAIBRCJ", 0)
	rotor2 := NewRotor("AJDKSIRUXBLHWTMCQGZNPYFVOE", 0)
	rotor3 := NewRotor("BDFHJLCPRTXVZNYEIWGAKMUSQO", 0)
	reflector := NewReflector("YRUHQSLDPXNGOKMIEBFZCWVJAT")

	// Создаем машину Enigma
	enigma := NewEnigma([]Rotor{rotor1, rotor2, rotor3}, reflector)

	// Чтение сообщения из файла
	filePath := "message.txt"
	messageBytes, err := ioutil.ReadFile(filePath)
	if err != nil {
		fmt.Println("Ошибка при чтении файла:", err)
		return
	}
	message := string(messageBytes) // Преобразуем байты в строку

	// Шифруем сообщение
	encryptedMessage := enigma.Encrypt(message)
	fmt.Println("Зашифрованное сообщение:", encryptedMessage)

	// Записываем зашифрованное сообщение в файл
	err = writeToFile("encrypted_message.txt", encryptedMessage)
	if err != nil {
		fmt.Println("Ошибка при записи зашифрованного сообщения в файл:", err)
		return
	}

	// Сбрасываем позиции роторов в исходное состояние
	enigma = NewEnigma([]Rotor{rotor1, rotor2, rotor3}, reflector)

	// Расшифровка зашифрованного сообщения
	decryptedMessage := enigma.Encrypt(encryptedMessage)
	fmt.Println("Расшифрованное сообщение:", decryptedMessage)

	// Записываем расшифрованное сообщение в файл
	err = writeToFile("decrypted_message.txt", decryptedMessage)
	if err != nil {
		fmt.Println("Ошибка при записи расшифрованного сообщения в файл:", err)
		return
	}

	duration := time.Since(startTime)
	fmt.Printf("Время выполнения кода: %v\n", duration)
}
