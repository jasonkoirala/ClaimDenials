package config

// const DB_USERNAME = "root"
// const DB_PASSWORD = "jason123"
// const DB_NAME = "dev"
// const DB_HOST = "127.0.0.1"
// const DB_PORT = "3306"

// var db *gorm.DB

// func GetDatabase() *gorm.DB {
// 	db = ConnectToDatabase()
// 	return db
// }

// func ConnectToDatabase() *gorm.DB {
// 	var err error
// 	dsn := DB_USERNAME + ":" + DB_PASSWORD + "@tcp" + "(" + DB_HOST + ":" + DB_PORT + ")/" + DB_NAME + "?" + "parseTime=true&loc=Local"
// 	fmt.Println("dsn : ", dsn)
// 	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})

// 	if err != nil {
// 		fmt.Println("Error connecting to database : error: %v", err)
// 		return nil
// 	}

// 	return db
// }
