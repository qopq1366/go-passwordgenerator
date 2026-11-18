package main

import (
	"errors"
	"flag"
	"fmt"
	"io"
	"os"

	"github.com/qopq1366/go-passwordgenerator/internal/password"
)

const version = "1.0.0"

func main() {
	os.Exit(run(os.Args[1:], os.Stdout, os.Stderr))
}

func run(args []string, stdout, stderr io.Writer) int {
	flags := flag.NewFlagSet("passwordgen", flag.ContinueOnError)
	flags.SetOutput(stderr)

	var (
		length    int
		count     int
		lowercase bool
		uppercase bool
		digits    bool
		symbols   bool
		showVer   bool
	)

	flags.IntVar(&length, "length", 16, "длина пароля")
	flags.IntVar(&length, "l", 16, "сокращение для --length")
	flags.IntVar(&count, "count", 1, "количество паролей")
	flags.IntVar(&count, "n", 1, "сокращение для --count")
	flags.BoolVar(&lowercase, "lowercase", true, "использовать строчные буквы")
	flags.BoolVar(&uppercase, "uppercase", true, "использовать прописные буквы")
	flags.BoolVar(&digits, "digits", true, "использовать цифры")
	flags.BoolVar(&symbols, "symbols", true, "использовать специальные символы")
	flags.BoolVar(&showVer, "version", false, "показать версию")

	flags.Usage = func() {
		fmt.Fprintf(stderr, "Генератор криптографически стойких паролей.\n\n")
		fmt.Fprintf(stderr, "Использование:\n  passwordgen [настройки]\n\nНастройки:\n")
		flags.PrintDefaults()
		fmt.Fprintf(stderr, "\nЧтобы отключить набор символов, укажите, например, --symbols=false.\n")
	}

	if err := flags.Parse(args); err != nil {
		if errors.Is(err, flag.ErrHelp) {
			return 0
		}
		return 2
	}
	if flags.NArg() != 0 {
		fmt.Fprintf(stderr, "ошибка: неожиданный аргумент: %s\n", flags.Arg(0))
		flags.Usage()
		return 2
	}
	if showVer {
		fmt.Fprintf(stdout, "passwordgen %s\n", version)
		return 0
	}
	if count <= 0 {
		fmt.Fprintln(stderr, "ошибка: количество паролей должно быть больше нуля")
		return 2
	}

	opts := password.Options{
		Length:    length,
		Lowercase: lowercase,
		Uppercase: uppercase,
		Digits:    digits,
		Symbols:   symbols,
	}
	generator := password.New()
	for i := 0; i < count; i++ {
		generated, err := generator.Generate(opts)
		if err != nil {
			fmt.Fprintf(stderr, "ошибка: %v\n", err)
			return 2
		}
		fmt.Fprintln(stdout, generated)
	}
	return 0
}
