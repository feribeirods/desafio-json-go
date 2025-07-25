package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"sort"
	"time"

	"github.com/gin-gonic/gin"
)

type Projects struct {
	Name      string `json:"name"`
	Completed bool   `json:"completed"`
}

type Team struct {
	Name     string     `json:"name"`
	Leader   bool       `json:"leader"`
	Projects []Projects `json:"projects"`
}

type Logs struct {
	Date   string `json:"date"`
	Action string `json:"action"`
}

type Users struct {
	Id      string `json:"id"`
	Name    string `json:"name"`
	Age     int    `json:"age"`
	Score   int    `json:"score"`
	Active  bool   `json:"active"`
	Country string `json:"country"`
	Team    Team   `json:"team"`
	Logs    []Logs `json:"logs"`
}

type EndpointResult struct {
	Path            string  `json:"path"`
	ResponseSuccess bool    `json:"response_success"`
	DurationMs      float64 `json:"duration_ms"`
	ValidJSON       bool    `json:"valid_json"`
}

var users []Users
var superusers []Users

func setupRouter() *gin.Engine {
	r := gin.Default()

	r.POST("/users", handlePostUsers)
	r.GET("/superusers", handleSuperusers)
	r.GET("/top-countries", handleTopCountries)
	r.GET("/team-insights", handleTeamInsights)
	r.GET("/active-users-per-day", handleActiveUsers)
	r.GET("/evaluation", handleEvaluation)

	return r
}

func handlePostUsers(c *gin.Context) {
	reqInit := time.Now()

	// Pega o arquivo enviado no form
	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"erro": "Nenhum arquivo foi enviado",
		})
		return
	}

	// Abre o arquivo para leitura
	openedFile, err := file.Open()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"erro": "Erro ao abrir o arquivo",
		})
		return
	}
	defer openedFile.Close()

	// Lê o conteúdo do arquivo
	fileContent, err := io.ReadAll(openedFile)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"erro": "Erro ao ler o arquivo",
		})
		return
	}

	err = json.Unmarshal(fileContent, &users)

	reqEnd := time.Now()
	reqDuration := reqEnd.Sub(reqInit)

	if err == nil {
		c.JSON(http.StatusOK, gin.H{
			"response": "Arquivo recebido com sucesso!",
			"duration": reqDuration.Milliseconds(),
		})
	}

}

func handleSuperusers(c *gin.Context) {
	reqInit := time.Now()

	helperSuperusers()

	reqEnd := time.Now()
	reqDuration := reqEnd.Sub(reqInit)

	c.JSON(http.StatusOK, gin.H{
		"response":       superusers,
		"duration":       reqDuration.Milliseconds(),
		"superusers_qty": len(superusers),
	})
}

func handleTopCountries(c *gin.Context) {
	reqInit := time.Now()

	helperSuperusers()

	type CountryInfo struct {
		Country string `json:"country"`
		Qty     int    `json:"qty"`
	}

	var topCountries []CountryInfo

	for _, u := range superusers {

		existePais := false

		for i, t := range topCountries {

			if t.Country == u.Country {
				topCountries[i].Qty += 1
				existePais = true
			}
		}

		if !existePais {
			topCountries = append(topCountries, CountryInfo{u.Country, 1})
		}
	}

	sort.Slice(topCountries, func(i, j int) bool {
		return topCountries[i].Qty > topCountries[j].Qty
	})

	reqEnd := time.Now()
	reqDuration := reqEnd.Sub(reqInit)

	c.JSON(http.StatusOK, gin.H{
		"response": topCountries[0:5],
		"duration": reqDuration.Milliseconds(),
	})
}

func handleTeamInsights(c *gin.Context) {
	reqInit := time.Now()

	type TeamInfo struct {
		Name             string   `json:"name"`
		Qty              int      `json:"qty"`
		Leaders          []string `json:"leaders"`
		FinishedProjects []string `json:"finished_projects"`
		ActiveMembers    int      `json:"active_members"`
		PercentActive    float32  `json:"percent_active"`
	}

	var teamInfos []TeamInfo

	for _, u := range users {

		existTeam := false

		for i, t := range teamInfos {

			if t.Name == u.Team.Name {
				teamInfos[i].Qty += 1
				existTeam = true

				if u.Active {
					teamInfos[i].ActiveMembers += 1
				}

				if u.Team.Leader {
					leaderName := u.Name
					teamInfos[i].Leaders = append(teamInfos[i].Leaders, leaderName)
				}

				for _, p := range u.Team.Projects {
					existProject := false

					for _, fp := range teamInfos[i].FinishedProjects {
						if p.Name == fp {
							existProject = true
						}
					}

					if p.Completed && !existProject {
						teamInfos[i].FinishedProjects = append(teamInfos[i].FinishedProjects, p.Name)
					}
				}
			}
		}

		if !existTeam {
			teamInfos = append(teamInfos, TeamInfo{u.Team.Name, 1, []string{}, []string{}, 0, 0})

		}
	}

	for i, t := range teamInfos {
		teamInfos[i].PercentActive = float32(t.ActiveMembers) / float32(t.Qty) * 100
	}

	reqEnd := time.Now()
	reqDuration := reqEnd.Sub(reqInit)

	c.JSON(http.StatusOK, gin.H{
		"response": teamInfos,
		"duration": reqDuration.Milliseconds(),
	})
}

func handleActiveUsers(c *gin.Context) {
	reqInit := time.Now()

	type UsersLogin struct {
		Date string `json:"date"`
		Qty  int    `json:"qty"`
	}

	var usersLogin []UsersLogin

	for _, u := range users {
		for _, l := range u.Logs {
			existLogin := false
			for i, ul := range usersLogin {
				if ul.Date == l.Date {
					existLogin = true
					usersLogin[i].Qty += 1
				}
			}

			if l.Action == "login" && !existLogin {
				usersLogin = append(usersLogin, UsersLogin{l.Date, 0})
			}
		}
	}

	reqEnd := time.Now()
	reqDuration := reqEnd.Sub(reqInit)

	c.JSON(http.StatusOK, gin.H{
		"response": usersLogin,
		"duration": reqDuration.Milliseconds(),
	})
}

func handleEvaluation(c *gin.Context) {

	var results []EndpointResult
	var result map[string]interface{}

	// Requisição POST
	filePath := "usuarios.json"
	file, err := os.Open(filePath)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	// Prepara o corpo multipart
	var body bytes.Buffer
	writer := multipart.NewWriter(&body)

	// Adiciona o arquivo ao form
	part, err := writer.CreateFormFile("arquivo", filepath.Base(file.Name()))
	if err != nil {
		panic(err)
	}

	_, err = io.Copy(part, file)
	if err != nil {
		panic(err)
	}

	writer.Close()

	// Faz a requisição POST
	resp, err := http.Post("http://localhost:8081/users", writer.FormDataContentType(), &body)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	// Lê e imprime a resposta
	respBody, _ := io.ReadAll(resp.Body)
	json.Unmarshal(respBody, &result)

	fmt.Printf("Resposta JSON: %+v\n", result)

	results = append(results, EndpointResult{"http://localhost:8081/users", resp.StatusCode == 200, result["duration"].(float64), true})

	// Requisições GET
	endpointsGet := []string{
		"http://localhost:8081/superusers",
		"http://localhost:8081/top-countries",
		"http://localhost:8081/team-insights",
		"http://localhost:8081/active-users-per-day",
	}

	for _, url := range endpointsGet {
		resp, _ := http.Get(url)
		body, _ := io.ReadAll(resp.Body)

		// Valida se o status da resposta é 200
		responseSuccess := resp.StatusCode == 200

		// Valida se o retorno é um JSON válido
		if err := json.Unmarshal(body, &result); err != nil {
			fmt.Println("Erro ao parsear JSON:", err)
			return
		}

		results = append(results, EndpointResult{url, responseSuccess, result["duration"].(float64), true})

	}

	c.JSON(http.StatusOK, gin.H{
		"evaluation": results,
	})
}

func main() {
	// Cria uma instância do Gin
	router := setupRouter()
	println("Servidor rodando em http://localhost:8081")
	router.Run(":8081")
}

func helperSuperusers() {
	for _, u := range users {
		if u.Score > 900 && u.Active {
			superusers = append(superusers, u)
		}
	}
}
