package handlers

import (
	"context"
	"example/data"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

type Products struct {
	l *log.Logger
}

func NewProducts(l *log.Logger) *Products {
	return &Products{l}
}

// Функция которая реализует вывод всех продкутов
func (p *Products) GetProducts(w http.ResponseWriter, r *http.Request) {
	p.l.Println("Handle GET Product")
	lp := data.GetProducts()
	err := lp.ToJSON(w)
	if err != nil {
		http.Error(w, "Unable to marshal json", http.StatusInternalServerError)
	}
}

func (p *Products) AddProduct(w http.ResponseWriter, r *http.Request) {
	p.l.Println("Handle POST Product")
	prod := r.Context().Value(KeyProduct{}).(*data.Product)
	data.AddProducts(prod)
}

func (p *Products) UpdateProducts(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)                 // Возрващает карту map[string]string( например клиент ввел запрос .../123, то карта будет иметь вид - map[string]string{"id": "123"})
	id, err := strconv.Atoi(vars["id"]) //Присваиваем переменно id значение из карты
	if err != nil {
		http.Error(w, "Unable to convert id", http.StatusBadRequest)
	}

	p.l.Println("Handle PUT Product", id)

	prod := r.Context().Value(KeyProduct{}).(*data.Product)
	err = data.UpdateProduct(id, prod)

	if err == data.ErrProductNotFounded {
		http.Error(w, "Product not found", http.StatusNotFound)
		return
	}
	if err != nil {
		http.Error(w, "Product not found", http.StatusNotFound)
		return
	}
}

type KeyProduct struct{} //Уникальный ключ

// Middleware - Функция которая которая выполняться перед основным обработчиком запросов
func (p *Products) MiddlewareProductValidation(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		prod := &data.Product{}

		err := prod.FromJSON(r.Body)
		if err != nil {
			http.Error(w, "Unable unmarshal json", http.StatusBadRequest)
			return
		}
		ctx := context.WithValue(r.Context(), KeyProduct{}, prod) // Создвем контекст и записываем туда информацию о продукте.KyeProduct{} - это формально его уникальное имя(ярлык), чтобы другие программы могли с ним работать
		req := r.WithContext(ctx)                                 // Создаем новый запрос с новый контекстом

		next.ServeHTTP(w, req) // Передаем новый запрос с обновленным контекстом дальше
	})
}
