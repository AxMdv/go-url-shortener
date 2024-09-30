package compressor_test

import (
	"bytes"
	"compress/gzip"
)

func Example() {
	var b bytes.Buffer
	w := gzip.NewWriter(&b)
	reqBody := `{"url": "https://yandex.ru"} `

	_, err := w.Write([]byte(reqBody))
	if err != nil {
		panic(err)
	}
	// compressReader, err := compressor.NewCompressReader(w)

	// w := compressor.NewCompressWriter()
	// request, err := http.NewRequest(http.MethodPost, "/", writer)
	// request.Header
	// if err != nil {
	// 	panic(err)
	// }
	// compressReader, err := compressor.NewCompressReader(request.Body)
	// if err != nil {
	// 	panic(err)
	// }
	// request.Body = compressReader
	// bytes, err := io.ReadAll(request.Body)

	// if err != nil {
	// 	panic(err)
	// }
	// fmt.Println(bytes)

	// Output:
	// {"url": "https://yandex.ru"}
}
