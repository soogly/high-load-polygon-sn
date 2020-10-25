package models

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/jackc/pgx/pgxpool"
)

var dbpoolM, dbpoolS *pgxpool.Pool // Master, Slave connection pools
var connM, connS *pgxpool.Conn

func Init() {
	log.Println("initialization DB connection pools")
	var err error

	dbpoolM, err = pgxpool.Connect(context.Background(), os.Getenv("DB_MASTER_URL"))
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		os.Exit(1)
	}

	dbpoolS, err = pgxpool.Connect(context.Background(), os.Getenv("DB_SLAVE_URL"))

	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		os.Exit(1)
	}

	go checkPool(dbpoolM, "MASTER")
	go checkPool(dbpoolS, "SLAVE")
}

func checkPool(pool *pgxpool.Pool, db_role string) {
	log.Println("Acquiring connection from pool " + db_role)
	conn, err := pool.Acquire(context.Background())
	log.Println("Acquired")

	if err != nil {
		log.Fatal(err)
	}
	defer conn.Release()

	var greeting string

	err = conn.QueryRow(context.Background(), "select 'HELLO FROM Timing "+db_role+" db!'").Scan(&greeting)
	if err != nil {
		log.Println("err2")

		fmt.Fprintf(os.Stderr, "[Timing_mS] "+db_role+" DB :: QueryRow failed: %v\n", err)
		os.Exit(1)
	}

	fmt.Println(greeting)
}
