package database

import (
	"sync"

	"backend/prisma/db"
)

var (
	client *db.PrismaClient
	once   sync.Once
)

func Client() *db.PrismaClient {
	once.Do(func() {
		client = db.NewClient()
		if err := client.Prisma.Connect(); err != nil {
			panic(err)
		}
	})
	return client
}

func Close() {
	if client != nil {
		if err := client.Prisma.Disconnect(); err != nil {
			panic(err)
		}
	}
}
