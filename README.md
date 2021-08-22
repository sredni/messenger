# Messenger
Simple library that allows to send multiple messages without exceeding receiver limits 
### Installation
```shell
go get github.com/sredni/messenger
```
### Example usage
```go
package main

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/sredni/messenger"
)

func main() {
	ctx := context.Background()
	cancelCtx, cancel := context.WithCancel(ctx)
	defer cancel()

	errorHandler := func(err messenger.DeliveryError) {
		log.Println(err)
	}

	w := messenger.NewHttpPostWorker(&http.Client{}, "https://example.com/test", 4 * time.Second)
	msgr := messenger.NewManager(w, 2)
	msgr.Run(cancelCtx, errorHandler)
	
	err := msgr.Send([]messenger.Message{
		{Content: []byte("1")},
		{Content: []byte("2")},
		{Content: []byte("3")},
		{Content: []byte("4")},
	})
	if err != nil {
		log.Println(err)
		return
	}

	msgr.Sync()
}
```
