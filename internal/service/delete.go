package service

import (
	"context"
	"sync"
	"time"

	"github.com/AxMdv/go-url-shortener/internal/storage"
)

// DeleteTask contain data about user id and url to delete.
type DeleteTask struct {
	ShortenedURL string
	UUID         string
}

// DeleteURLBatch deletes batch of urls.
func (s *shortenerService) DeleteURLBatch(deleteBatch storage.DeleteBatch) error {

	// сигнальный канал для завершения горутин
	doneCh := make(chan struct{})
	// закрываем его при завершении программы
	defer close(doneCh)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// fanIn realizaton
	inputCh := generator(doneCh, deleteBatch)
	// получаем слайс каналов из 10 рабочих former
	channels := fanOut(doneCh, inputCh)
	// а теперь объединяем десять каналов в один
	formResultCh := fanIn(doneCh, channels...)

	var formedToDelete []storage.FormedURL
	for form := range formResultCh {
		formedToDelete = append(formedToDelete, form)
	}

	err := s.urlRepository.DeleteURLBatch(ctx, formedToDelete)
	return err
}

func generator(doneCh chan struct{}, input storage.DeleteBatch) chan DeleteTask {
	inputCh := make(chan DeleteTask)

	go func() {
		defer close(inputCh)

		for _, data := range input.ShortenedURL {
			task := DeleteTask{
				UUID:         input.UUID,
				ShortenedURL: data,
			}
			select {
			case <-doneCh:
				return
			case inputCh <- task:
			}

		}
	}()

	return inputCh
}

// form функция из предыдущего примера, делает то же, что и делала
func form(doneCh chan struct{}, inputCh chan DeleteTask) chan storage.FormedURL {
	formRes := make(chan storage.FormedURL)

	go func() {
		defer close(formRes)

		for data := range inputCh {

			formed := storage.FormedURL{
				UUID:         data.UUID,
				ShortenedURL: data.ShortenedURL,
			}

			select {
			case <-doneCh:
				return
			case formRes <- formed:
			}
		}
	}()
	return formRes
}

// fanOut принимает канал данных, порождает 10 горутин
func fanOut(doneCh chan struct{}, inputCh chan DeleteTask) []chan storage.FormedURL {
	// количество горутин add
	numWorkers := 10
	// каналы, в которые отправляются результаты
	channels := make([]chan storage.FormedURL, numWorkers)

	for i := 0; i < numWorkers; i++ {
		// получаем канал из горутины add
		formResultCh := form(doneCh, inputCh)
		// отправляем его в слайс каналов
		channels[i] = formResultCh
	}

	// возвращаем слайс каналов
	return channels
}

// fanIn объединяет несколько каналов resultChs в один.
func fanIn(doneCh chan struct{}, resultChs ...chan storage.FormedURL) chan storage.FormedURL {
	// конечный выходной канал в который отправляем данные из всех каналов из слайса, назовём его результирующим
	finalCh := make(chan storage.FormedURL)

	// понадобится для ожидания всех горутин
	var wg sync.WaitGroup

	// перебираем все входящие каналы
	for _, ch := range resultChs {
		// в горутину передавать переменную цикла нельзя, поэтому делаем так
		chClosure := ch

		// инкрементируем счётчик горутин, которые нужно подождать
		wg.Add(1)

		go func() {
			// откладываем сообщение о том, что горутина завершилась
			defer wg.Done()

			// получаем данные из канала
			for data := range chClosure {
				select {
				// выходим из горутины, если канал закрылся
				case <-doneCh:
					return
				// если не закрылся, отправляем данные в конечный выходной канал
				case finalCh <- data:
				}
			}
		}()
	}

	go func() {
		// ждём завершения всех горутин
		wg.Wait()
		// когда все горутины завершились, закрываем результирующий канал
		close(finalCh)
	}()

	// возвращаем результирующий канал
	return finalCh
}
