package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"strings"
)

type Champion struct {
	ID         int    `json:"id"`
	Anio       int    `json:"anio"`
	Equipo     string `json:"equipo"`
	Pais       string `json:"pais"`
	Entrenador string `json:"entrenador"`
	Estadio    string `json:"estadio"`
	Resultado  string `json:"resultado"`
}

var champions []Champion

func jsonError(w http.ResponseWriter, msg string, code int) {

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)

	json.NewEncoder(w).Encode(map[string]string{
		"error": msg,
	})

}

func cargarDatos() {

	file, err := os.ReadFile("data/champions.json")

	if err != nil {
		fmt.Println("Error leyendo archivo")
		return
	}

	json.Unmarshal(file, &champions)

}

func guardarDatos() {

	data, _ := json.MarshalIndent(champions, "", " ")

	os.WriteFile("data/champions.json", data, 0644)

}

func getAll(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json")

	pais := r.URL.Query().Get("pais")
	anio := r.URL.Query().Get("anio")

	var resultado []Champion

	for _, c := range champions {

		if pais != "" && !strings.EqualFold(c.Pais, pais) {
			continue
		}

		if anio != "" {

			a, _ := strconv.Atoi(anio)

			if c.Anio != a {
				continue
			}

		}

		resultado = append(resultado, c)

	}

	json.NewEncoder(w).Encode(resultado)

}

func getByID(w http.ResponseWriter, id int) {

	for _, c := range champions {

		if c.ID == id {

			json.NewEncoder(w).Encode(c)
			return

		}

	}

	jsonError(w, "Campeon no encontrado", 404)

}

func createChampion(w http.ResponseWriter, r *http.Request) {

	var nuevo Champion

	err := json.NewDecoder(r.Body).Decode(&nuevo)

	if err != nil {

		jsonError(w, "JSON invalido", 400)
		return

	}

	if nuevo.Equipo == "" || nuevo.Pais == "" {

		jsonError(w, "Campos requeridos faltantes", 400)
		return

	}

	champions = append(champions, nuevo)

	guardarDatos()

	w.WriteHeader(http.StatusCreated)

	json.NewEncoder(w).Encode(nuevo)

}

func updateChampion(w http.ResponseWriter, r *http.Request, id int) {

	var actualizado Champion

	json.NewDecoder(r.Body).Decode(&actualizado)

	for i, c := range champions {

		if c.ID == id {

			champions[i] = actualizado
			guardarDatos()

			json.NewEncoder(w).Encode(actualizado)
			return

		}

	}

	jsonError(w, "No encontrado", 404)

}

func deleteChampion(w http.ResponseWriter, id int) {

	for i, c := range champions {

		if c.ID == id {

			champions = append(champions[:i], champions[i+1:]...)

			guardarDatos()

			json.NewEncoder(w).Encode(map[string]string{
				"mensaje": "Eliminado correctamente",
			})

			return

		}

	}

	jsonError(w, "No encontrado", 404)

}

func handler(w http.ResponseWriter, r *http.Request) {

	path := strings.TrimPrefix(r.URL.Path, "/api/items")

	if path == "" || path == "/" {

		switch r.Method {

		case "GET":
			getAll(w, r)

		case "POST":
			createChampion(w, r)

		default:
			jsonError(w, "Metodo no permitido", 405)

		}

		return
	}

	idStr := strings.TrimPrefix(path, "/")

	id, err := strconv.Atoi(idStr)

	if err != nil {

		jsonError(w, "ID invalido", 400)
		return

	}

	switch r.Method {

	case "GET":
		getByID(w, id)

	case "PUT":
		updateChampion(w, r, id)

	case "DELETE":
		deleteChampion(w, id)

	default:
		jsonError(w, "Metodo no permitido", 405)

	}

}

func main() {

	cargarDatos()

	http.HandleFunc("/api/items", handler)
	http.HandleFunc("/api/items/", handler)

	fmt.Println("Servidor corriendo en puerto 24841")

	http.ListenAndServe(":24841", nil)

}
