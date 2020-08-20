package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"time"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
)

const dbName = "Fanap.db"
const dbDriver = "sqlite3"
const serverPort = "4567"

type Model struct {
	ID        uint `gorm:"primary_key"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time
}

type RectangleModel struct {
	gorm.Model
	X      int
	Y      int
	Width  int
	Height int
}

type Rectangle struct {
	X      int `json:"x"`
	Y      int `json:"y"`
	Width  int `json:"width"`
	Height int `json:"height"`
}

type ResponseData struct {
	Main  Rectangle   `json:"main"`
	Input []Rectangle `json:"input"`
}

func (res *ResponseData) UnmarshalJSON(buf []byte) {
	json.Unmarshal(buf, &res)
}

// type ResponseData struct {
// 	main  Rectangle `json:"main"`
// 	input []interface{} `json:"input"`
// }

func main() {
	createTable()
	fmt.Println("Server Port:", serverPort)
	http.HandleFunc("/", requestHandler)
	http.ListenAndServe(":"+serverPort, nil)

}

func requestHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.Error(w, "404 not found.", http.StatusNotFound)
		return
	}

	switch r.Method {
	case "GET":
		db, err := gorm.Open(dbDriver, dbName)
		if err != nil {
			fmt.Println("dataHandler failed to connect database: ", err)
			os.Exit(1)
		}
		defer db.Close()
		var rectangles []RectangleModel
		_ = db.Find(&rectangles)
		res, _ := json.Marshal(rectangles)
		fmt.Println("res :", string(res))
		fmt.Fprintf(w, string(res))
		// http.ServeFile(w, r, "form.html")
	case "POST":
		if err := r.ParseForm(); err != nil {
			fmt.Fprintf(w, "ParseForm() err: %v", err)
			return
		}
		// fmt.Println("Post request")
		// fmt.Println(r)
		// fmt.Println("main:", r.FormValue("main"))
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			fmt.Println("getCoinInfo failed to ioutil.ReadAll(resp.Body):", err)
			os.Exit(1)
		}
		var data ResponseData
		data.UnmarshalJSON(body)
		// var inputs []Rectangle
		// json.Unmarshal(body, &data)
		// json.Unmarshal([]byte(data.input, &inputs)
		dataHandler(data)
		// fmt.Println("body:", string(body))
		// fmt.Println("main.width:", data.Main.Width)
		// for i, v := range data.Input {
		// 	fmt.Println(i, " data.input:", v)
		// }
		// fmt.Println("data.input:", data.Input)
		// fmt.Println("data.input length:", len(data.Input))
		fmt.Fprintf(w, "")
		// // main := r.FormValue("main")
		// ‫‪input‬‬Array := r.FormValue("input")
		// fmt.Fprintf(w, "main = %s\n", main)
		// fmt.Fprintf(w, "‫‪input‬‬ = %s\n", ‫‪input‬‬)
		// http.Response(w, "", http.StatusOK )
	default:
		http.Error(w, "Sorry, only GET and POST methods are supported.", http.StatusNotFound)
		fmt.Fprintf(w, "Sorry, only GET and POST methods are supported.")
	}
}

func createTable() {
	db, err := gorm.Open(dbDriver, dbName)
	if err != nil {
		fmt.Println("failed to connect database: ", err)
		os.Exit(1)
	}
	defer db.Close()

	db.AutoMigrate(&RectangleModel{})
}

func dataHandler(data ResponseData) {
	db, err := gorm.Open(dbDriver, dbName)
	if err != nil {
		fmt.Println("dataHandler failed to connect database: ", err)
		os.Exit(1)
	}
	defer db.Close()

	for _, rec := range data.Input {
		if checkIsCommon(data.Main, rec) {
			fmt.Println("Ok:", rec.X, ", ", rec.Y, ", ", rec.Width, ", ", rec.Height)
			db.Create(&RectangleModel{X: rec.X, Y: rec.Y, Width: rec.Width, Height: rec.Height})
		}
	}
}

func checkIsCommon(main Rectangle, r Rectangle) bool {
	if main.X <= r.X+r.Width && r.X <= main.X+main.Width && main.Y <= r.Y+r.Height && r.Y <= main.Y+main.Height {
		return true
	}
	return false
}


// {
// 	"main": {"x": 0, "y": 0, "width": 10, "height": 20},
// 	"input": [
// 		   {"x": 2, "y": 18, "width": 5, "height": 4},
// 		   {"x": 12, "y": 18, "width": 5, "height": 4},
// 		   {"x": -1, "y": -1, "width": 5, "height": 4}
// 	 ]
// 	}
