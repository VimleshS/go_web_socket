package main

import (
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
)

type page struct {
	WsEndPoint string
}

func main() {
	fs := http.FileServer(http.Dir("static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {

		// tmpl := template.Must(template.New("layout").ParseFiles("templates/index.html"))

		tmpl, err := template.ParseFiles("templates/index.html", "templates/input.html")
		if err != nil {
			log.Println(err.Error())
			http.Error(w, http.StatusText(500), 500)
			return
		}

		p := page{WsEndPoint: r.Host}
		if err := tmpl.ExecuteTemplate(w, "layout", p); err != nil {
			log.Println(err.Error())
			http.Error(w, http.StatusText(500), 500)
		}

		// p := page{WsEndPoint: r.Host}
		// tmpl.Execute(w, p)
	})

	http.HandleFunc("/ws", ConnWs)
	err := http.ListenAndServe(":9090", nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}

}

func ConnWs(w http.ResponseWriter, r *http.Request) {
	ws, err := websocket.Upgrade(w, r, nil, 1024, 1024)
	if _, ok := err.(websocket.HandshakeError); ok {
		http.Error(w, "Not a websocket handshake", 400)
		return
	} else if err != nil {
		log.Println(err)
		return
	}

	var img64 []byte

	res := map[string]interface{}{}
	for {
		// messageType, p, err := ws.ReadMessage()
		// fmt.Println(messageType)
		// fmt.Println(string(p))
		// fmt.Println(err)

		if err = ws.ReadJSON(&res); err != nil {
			fmt.Printf("%v \n", res)

			if err.Error() == "EOF" {
				return
			}
			// ErrShortWrite means a write accepted fewer bytes than requested then failed to return an explicit error.
			if err.Error() == "unexpected EOF" {
				return
			}
			fmt.Println("Read : " + err.Error())
			return
		}

		res["a"] = "a"
		log.Println(res)

		// for {
		files, _ := ioutil.ReadDir("./images")
		for _, f := range files {
			img64, _ = ioutil.ReadFile("./images/" + f.Name())

			/*Below code works*/
			// str := base64.StdEncoding.EncodeToString(img64)
			// res["img64"] = str

			// if err = ws.WriteJSON(&res); err != nil {
			// 	fmt.Println("watch dir - Write : " + err.Error())
			// }

			/*Efficient only for small data*/
			/*TextMessage*/
			/*Remember ruby multipart upload, always encoded into base64 */

			// _dst := make([]byte, base64.StdEncoding.EncodedLen(len(img64)))
			// base64.StdEncoding.Encode(_dst, img64)

			/*For larger data use NewEncoder */
			// _dst := &bytes.Buffer{}
			// encoder := base64.NewEncoder(base64.StdEncoding, _dst)
			// encoder.Write(img64)
			// encoder.Close()

			// ws.WriteMessage(websocket.TextMessage, _dst.Bytes())

			/*BinaryMessage*/
			/*utf-8 encoding*/
			ws.WriteMessage(websocket.BinaryMessage, img64)

			time.Sleep(50 * time.Millisecond)
		}
		time.Sleep(50 * time.Millisecond)
		// }
	}
}
