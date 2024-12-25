package handlers

import (
	"fmt"
	"io"
	"log"
	"net/http"
)

// Создали структуру в которой хранится наш лог. Нуже для того(пока) чтобы, если мы подключали бд мы могли использовать свой лог( и его модернизировать) а не общий
type Hello struct {
	l *log.Logger
}

// Конструктор. МЫ поставили & т.к с указателями программа быстрее работает чем с оригеналом
// &Hello{l} - & поставили для того чтобы мы могли вносить изменения, и могли передавать указатель а не копии
func NewHello(l *log.Logger) *Hello {
	return &Hello{l}
}

// h *Hello - это простыми словами объект структуры Hello. Поставили * чтобы мы могли работать с объектами напрямую
func (h *Hello) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	h.l.Println("Hello world")
	// Мы написали 2 значения т.к если метод не смодежет прочитать тело то он запищет !nil в err. Иначе запишите nil
	// Мы можем опутсить err, и просто никуда не записывать
	/*Вот пример
	d, _ := io.ReadAll(r.Body)
	И обработчик записывать ненадо
	*/
	d, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Ooops", http.StatusBadRequest)
		return
	}
	// Записали %s т.к мы хотим вренуть строку
	fmt.Fprintf(w, "Hello %s", d)
}
