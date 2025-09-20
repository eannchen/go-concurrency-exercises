package main

// Loader implements KeyStoreCacheLoader.
type Loader struct {
	DB *MockDB
}

// Load gets the data from the database.
func (l *Loader) Load(key string) string {
	val, err := l.DB.Get(key)
	if err != nil {
		panic(err)
	}
	return val
}
