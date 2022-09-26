package core

import (
	"encoding/hex"
	"fmt"
	"io/ioutil"
	"math/rand"
	"net/http"
	"os"
	"os/exec"
	"time"
)

const (
	Template = `{
    "type": "Tx AlonzoEra",
    "description": "",
    "cborHex": "%s"
}
`
)

type Server struct {
}

func NewServer() *Server {
	return &Server{}
}

func (s *Server) submit(w http.ResponseWriter, req *http.Request) {
	fmt.Println("There is a new request.")
	body, err := ioutil.ReadAll(req.Body)
	if err != nil {
		http.Error(w, "Invalid request", 400)
		return
	}

	// TODO: Do cardano validation here.
	content := fmt.Sprintf(Template, hex.EncodeToString(body))
	fileName := GetRandomString(16)

	err = ioutil.WriteFile(fileName, []byte(content), 0644)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	defer os.Remove(fileName)

	cmd := exec.Command("cardano-cli", "transaction", "submit", "--tx-file", fileName, "--testnet-magic", "1")
	out, err := cmd.CombinedOutput()

	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	fmt.Println(string(out))
}

func (s *Server) Run() {
	rand.Seed(time.Now().UnixNano())

	http.HandleFunc("/submit", s.submit)

	fmt.Println("Listening on port 8090...")
	http.ListenAndServe(":8090", nil)
}
