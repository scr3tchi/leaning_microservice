package data

import (
	"encoding/json"
	"fmt"
	"io"
	"time"
)

type Product struct {
	ID          int     `json:"id"`
	Name        string  `json:"name"`        //Название
	Description string  `json:"description"` //Описание
	Price       float32 `json:"price"`       //Цена
	SKU         string  `json:"sku"`         //Артикл(Внутренний индентификатор продукта)
	CreateOn    string  `json:"_"`           // Дата создания
	UpdateOn    string  `json:"_"`           //Дата обновления. Если мы напишем `json"_'` то данное поле не буде показываться в выводе
	DeleteOn    string  `json:"_"`           //Дата удаления
}

// Методо обратный методу ToJSON
// io.Reader - Интерфейс который читаеть запрос
func (p *Product) FromJSON(r io.Reader) error {
	e := json.NewDecoder(r)
	return e.Decode(p)
}

type Products []*Product // Список всех продуктов().Позволит хранить большое количество данных, не копируя их каждый раз

// Метод который проверят что ответ для клиента будет в JSON формате
// w io.Writer -  интерфейс кторый говорит куда мы будем записывать данные
func (p *Products) ToJSON(w io.Writer) error {
	e := json.NewEncoder(w) //Создает объект для преобразования в JSON
	return e.Encode(p)      // Преобразует p в JSON и записывает в w
}

// В данном методе мы возвращаем список продуктов
func GetProducts() Products {
	return productList
}

func AddProducts(p *Product) {
	p.ID = getNextID()
	productList = append(productList, p)
}

func UpdateProduct(id int, p *Product) error {
	_, pos, err := findProduct(id)
	if err != nil {
		return err
	}
	p.ID = id
	//Формально заменяем текущий продкут на новый по позиции
	productList[pos] = p
	return nil
}

var ErrProductNotFounded = fmt.Errorf("Product not founded")

// Данный метод ищет продукт по ID
func findProduct(id int) (*Product, int, error) {
	for index, p := range productList {
		if p.ID == id {
			return p, index, nil
		}
	}
	return nil, -1, ErrProductNotFounded
}

func getNextID() int {
	lp := productList[len(productList)-1]
	return lp.ID + 1
}

// Создаем слайс типа структуры
var productList = []*Product{
	&Product{
		ID:          1,
		Name:        "Late",
		Description: "Frothy milky coffe",
		Price:       2.45,
		SKU:         "abs123",
		CreateOn:    time.Now().UTC().String(),
		UpdateOn:    time.Now().UTC().String(),
	},
	&Product{
		ID:          2,
		Name:        "Expresso",
		Description: "Short and strong coffe without milk",
		Price:       1.99,
		SKU:         "fgd34",
		CreateOn:    time.Now().UTC().String(),
		UpdateOn:    time.Now().UTC().String(),
	},
}
