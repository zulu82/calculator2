package main

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"
)

func main() {
	reader := bufio.NewReader(os.Stdin)

	fmt.Print("Введите выражение (например, '1 + 2' или 'III / II'): ")
	input, _ := reader.ReadString('\n')
	input = strings.TrimSpace(input) // Убираем пробелы и переносы строки

	result, err := calculate(input)
	if err != nil {
		fmt.Println("Паника:", err)
		return
	}
	fmt.Println("Результат:", result)
}

func calculate(input string) (string, error) {
	re := regexp.MustCompile(`^(\d+|[I|V|X]+)\s*([\+\-\*\/])\s*(\d+|[I|V|X]+)$`)
	matches := re.FindStringSubmatch(input)
	if len(matches) != 4 {
		return "", fmt.Errorf("строка не является математической операцией")
	}

	firstOperand := matches[1]
	operator := matches[2]
	secondOperand := matches[3]

	isArabic := isArabicNumber(firstOperand) && isArabicNumber(secondOperand)
	isRoman := isRomanNumber(firstOperand) && isRomanNumber(secondOperand)

	if (!isArabic && !isRoman) || (isArabic && isRoman) {
		return "", fmt.Errorf("используются одновременно разные системы счисления")
	}

	var a, b int
	var err error

	if isArabic {
		a, err = strconv.Atoi(firstOperand)
		b, err = strconv.Atoi(secondOperand)

		if err != nil || a < 1 || a > 10 || b < 1 || b > 10 {
			return "", fmt.Errorf("вводимые числа должны быть от 1 до 10 включительно")
		}
	} else {
		if !isValidRoman(firstOperand) || !isValidRoman(secondOperand) {
			return "", fmt.Errorf("неправильные римские цифры")
		}
		a = romanToArabic(firstOperand)
		b = romanToArabic(secondOperand)

		if a < 1 || a > 10 || b < 1 || b > 10 {
			return "", fmt.Errorf("вводимые числа должны быть от 1 до 10 включительно")
		}
	}

	var result int
	switch operator {
	case "+":
		result = a + b
	case "-":
		result = a - b
		if !isArabic && result < 1 {
			return "", fmt.Errorf("в римской системе нет отрицательных чисел")
		}
	case "*":
		result = a * b
	case "/":
		result = a / b
	default:
		return "", fmt.Errorf("неизвестный оператор")
	}

	if isArabic {
		return strconv.Itoa(result), nil
	} else {
		if result < 1 {
			return "", fmt.Errorf("результат должен быть больше нуля в римской системе")
		}
		return arabicToRoman(result), nil
	}
}

func isArabicNumber(s string) bool {
	_, err := strconv.Atoi(s)
	return err == nil
}

func isRomanNumber(s string) bool {
	return regexp.MustCompile(`^[IVXLCDM]+$`).MatchString(s)
}

func isValidRoman(s string) bool {
	re := regexp.MustCompile(`^(M{0,3})(CM|CD|D?C{0,3})(XC|XL|L?X{0,3})(IX|IV|V?I{0,3})$`)
	return re.MatchString(s)
}

func romanToArabic(roman string) int {
	romanNumerals := map[rune]int{'I': 1, 'V': 5, 'X': 10}
	result := 0
	prevValue := 0

	for _, char := range roman {
		currentValue := romanNumerals[char]
		if currentValue > prevValue {
			result += currentValue - 2*prevValue
		} else {
			result += currentValue
		}
		prevValue = currentValue
	}
	return result
}

func arabicToRoman(num int) string {
	roman := []struct {
		value  int
		symbol string
	}{
		{1000, "M"},
		{900, "CM"},
		{500, "D"},
		{400, "CD"},
		{100, "C"},
		{90, "XC"},
		{50, "L"},
		{40, "XL"},
		{10, "X"},
		{9, "IX"},
		{5, "V"},
		{4, "IV"},
		{1, "I"},
	}

	result := ""
	for _, r := range roman {
		for num >= r.value {
			result += r.symbol
			num -= r.value
		}
	}
	return result
}
