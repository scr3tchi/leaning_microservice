package main

import (
	"context"
	"example/handlers"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/gorilla/mux"
)

func main() {
	//os.Stdout - стандартный вывод в Go
	//product-api - то есть все отправленные собщения будут начинаться с этого текста(запусти программы и сам посомтри)
	//log.LstdFlags - флаги которые определяют какой формат будет у лог-сообщения(стандартный формат котоырй включает в логах дату и время сообщения)
	l := log.New(os.Stdout, "product-api", log.LstdFlags)
	ph := handlers.NewProducts(l)

	sm := mux.NewRouter()

	getRouter := sm.Methods(http.MethodGet).Subrouter() //Создали маршрутизатор только для GET - запросов
	getRouter.HandleFunc("/", ph.GetProducts)

	putRouter := sm.Methods(http.MethodPut).Subrouter() //Создали маршрутизатор только для запросов PUT
	putRouter.HandleFunc("/{id:[-0-9]+}", ph.UpdateProducts)
	putRouter.Use(ph.MiddlewareProductValidation) //Принимает промежуточное по

	postRouter := sm.Methods(http.MethodPost).Subrouter()
	postRouter.HandleFunc("/", ph.AddProduct)
	postRouter.Use(ph.MiddlewareProductValidation)
	//sm.Handle("/products", ph)

	s := http.Server{
		Addr:    ":9090",
		Handler: sm,
		//Таймайты защищают сервер от медленный и некорректных пользователей который занимают ресурсы сервера слишком долго
		IdleTimeout:  120 * time.Second, //Время ожидания при простое соединения(120 сек)
		ReadTimeout:  1 * time.Second,   //Макс время для чтения запроса клиента
		WriteTimeout: 1 * time.Second,   //Макс время для записи ответа для клинета
	}

	//Созданне горутины.Чтобы основной поток программы мог работать
	go func() {
		//Программа будет ждать запросы до тех пор, пока сервер не получит сигнал остановки
		err := s.ListenAndServe() //Запускает сервер на указанном адресе, начинает обработку входящих HTTP запросов
		if err != nil {
			l.Fatal(err) //Елси возникла ошибка выводим ее
		}
	}()

	sigChan := make(chan os.Signal, 1) //Создается канал для получение сигналов от оперцаионной системе

	//Notify сообщает Go, что канад sigChan должен получить уведомление о сигнале прерывания
	//Interrupt(обычнл CTRL+C) или завершение процесаа Kill(Отправлен извне(например командой kill))
	signal.Notify(sigChan, os.Interrupt)
	signal.Notify(sigChan, os.Kill)

	sig := <-sigChan
	l.Println("Recived teminate, grateful shutdown", sig)

	tc, _ := context.WithTimeout(context.Background(), 30*time.Second) //Ставит впромежуток времени. Елси операция идет дольше 30 секунд то сревер завешает работу
	s.Shutdown(tc)                                                     //Безопасное завершение сервера
}
