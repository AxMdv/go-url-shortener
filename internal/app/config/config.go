package config

import "flag"

var Flags struct {
	RunAddr            *string
	ResponseResultAddr *string
}

func init() {
	Flags.RunAddr = flag.String("a", "localhost:8080", "address and port to run server")
	Flags.ResponseResultAddr = flag.String("b", "http://localhost:8080/", "resut basic response address (before shortened URL)")
}

// func ParseFlags() {
// 	// регистрируем переменную flagRunAddr
// 	// как аргумент -a со значением :8080 по умолчанию
// 	flag.StringVar(&Flags.RunAddr, "a", "localhost:8080", "address and port to run server")
// 	// парсим переданные серверу аргументы в зарегистрированные переменные
// 	flag.Parse()
// }
