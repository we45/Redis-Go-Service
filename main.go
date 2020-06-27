package main

import (
	"context"
	"encoding/json"
	"fmt"
	"math/rand"
	"os"
	"time"

	"github.com/go-redis/redis"
	"github.com/rs/xid"
	"syreclabs.com/go/faker"
)

var ctx = context.Background()

type CCN struct {
	Pan  string `json:"card_num"`
	Name string `json:"card_name"`
	Cvv  int    `json:"cvv"`
}

func main() {
	host := os.Args[1]
	port := os.Args[2]
	redisAddr := fmt.Sprintf("%s:%s", host, port)
	rdb := redis.NewClient(&redis.Options{
		Addr:     redisAddr,
		Password: "",
		DB:       0,
	})

	for {
		// generate card number JSON
		newPan := faker.Finance().CreditCard()
		newCcn := CCN{
			Pan:  newPan,
			Name: faker.Name().Name(),
			Cvv:  rand.Intn(999),
		}
		guid := xid.New()
		json, jsonErr := json.Marshal(newCcn)

		if jsonErr != nil {
			panic(jsonErr)
		}

		err := rdb.Set(guid.String(), json, 0).Err()
		if err != nil {
			panic(err)
		}
		fmt.Printf("Loaded New Card with UUID: %s and Card Number %s\n", guid.String(), newPan)
		time.Sleep(4000 * time.Millisecond)
	}

}
