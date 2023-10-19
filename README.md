# iso8583
Este paquete permite validar el formato de iso8583

# Instación
Ejecuta el siguiente comando

```bash
go get -u github.com/pruebasgolang/iso8385Utils

```

# Uso

```go
package main

import (
	"github.com/pruebasgolang/iso8385Utils"
	"iso8583/models"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

var (
	message  = models.ISO8583Message{}
)



func main() {

	http.HandleFunc("/iso8583", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		err := json.NewDecoder(r.Body).Decode(&message)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		defer responseIso8583(w, r)

		er := utils.ValidateISO8583Message(message)
		if er != nil {
			for i := 0; i < len(er); i++ {
				http.Error(w, er[i].Error(), http.StatusBadRequest)
				//json.NewEncoder(w).Encode(er[i])
			}
		}

	})

	log.Fatal(http.ListenAndServe(":8080", nil))

}

func responseIso8583(w http.ResponseWriter, r *http.Request) {
 	simulateResponseMessage()
	json, err := json.Marshal(message)
	if err != nil {
		fmt.Fprintf(w, "Error: %v", err)
		return
	}
	fmt.Fprintf(w, "%s", json)
}

func simulateResponseMessage() {
	message.MTI_MessageOrigin = 1
}


```
