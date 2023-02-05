package gormsqlite

import "os"

const tmpDb = "./test-auth.db"

func setup() {
	if _, err := os.Stat(tmpDb); os.IsNotExist(err) {
		_, err := os.Create(tmpDb)
		if err != nil {
			panic(err)
		}
	}
}

func teardown() {
	if _, err := os.Stat(tmpDb); err == nil {
		err = os.Remove(tmpDb)
		if err != nil {
			panic(err)
		}
	}
}
