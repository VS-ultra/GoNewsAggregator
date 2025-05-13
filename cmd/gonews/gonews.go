package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"GoNewsAggregator/pkg/api"
	"GoNewsAggregator/pkg/rss"
	"GoNewsAggregator/pkg/storage"
)

// конфигурация приложения
type config struct {
	URLS   []string `json:"rss"`
	Period int      `json:"request_period"`
}

func main() {
	// Небольшая начальная пауза, чтобы docker успел поднять сеть
	time.Sleep(5 * time.Second)

	// --- ждём, пока Postgres будет готов ---
	var db *storage.DB
	var err error
	for i := 0; i < 10; i++ {
		db, err = storage.New()
		if err == nil {
			break
		}
		log.Printf("ждём Postgres… попытка %d: %v\n", i+1, err)
		time.Sleep(5 * time.Second)
	}
	if err != nil {
		log.Fatal("не удалось подключиться к БД после 10 попыток:", err)
	}

	// инициализируем API
	api := api.New(db)

	// читаем config.json из корня /app (куда его положили в Dockerfile)
	b, err := ioutil.ReadFile("./config.json")
	if err != nil {
		log.Fatal("конфиг не найден:", err)
	}
	var cfg config
	if err := json.Unmarshal(b, &cfg); err != nil {
		log.Fatal("не удалось распарсить config.json:", err)
	}

	// каналы для парсинга
	chPosts := make(chan []storage.Post)
	chErrs := make(chan error)
	for _, url := range cfg.URLS {
		go parseURL(url, chPosts, chErrs, cfg.Period)
	}

	// сохраняем в БД
	go func() {
		for posts := range chPosts {
			if err := db.StoreNews(posts); err != nil {
				log.Println("ошибка записи в БД:", err)
			}
		}
	}()

	// логируем ошибки парсинга
	go func() {
		for err := range chErrs {
			log.Println("ошибка парсинга RSS:", err)
		}
	}()

	// запускаем HTTP-сервер на порту 80
	log.Println("Server is up and running.")
	if err := http.ListenAndServe(":80", api.Router()); err != nil {
		log.Fatal(err)
	}
}

func parseURL(url string, posts chan<- []storage.Post, errs chan<- error, period int) {
	for {
		news, err := rss.Parse(url)
		if err != nil {
			errs <- err
		} else {
			posts <- news
		}
		time.Sleep(time.Duration(period) * time.Minute)
	}
}
