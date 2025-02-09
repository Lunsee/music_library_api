package endpoints

import (
	"encoding/json"
	"fmt"
	"log"
	"music_library_api/internal/database"
	"music_library_api/internal/models"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"

	_ "github.com/lib/pq"
)

// @Summary Получить список песен
// @Description Возвращает список песен с возможностью фильтрации по группе и названию, а также пагинацией
// @Tags Songs
// @Produce json
// @Param group query string false "Фильтр по названию группы"
// @Param song query string false "Фильтр по названию песни"
// @Param page query int false "Номер страницы (по умолчанию 1)"
// @Success 200 {array} models.Song "Список найденных песен"
// @Failure 400 {string} string "Некорректный запрос (например, если страница вне диапазона)"
// @Failure 500 {string} string "Ошибка сервера при получении данных из БД"
// @Router /api/getSongs [get]
func GetSongs(w http.ResponseWriter, r *http.Request) {
	log.Println("Info: <GetSongs> endpoint...")
	group := r.URL.Query().Get("group")
	song := r.URL.Query().Get("song")
	pageStr := r.URL.Query().Get("page")

	//default
	page := 1
	const limit = 10 // max limit items on page

	//get page number
	if pageStr != "" {
		var err error
		page, err = strconv.Atoi(pageStr)
		if err != nil { //check errors
			log.Printf("Error: Invalid page parameter: %v", err)
			http.Error(w, "Invalid page parameter", http.StatusBadRequest)
			return
		}
		log.Printf("Debug: Page parameter is set to %d", page)
	}

	var Songs_from_db []models.Song
	var filteredSongs []models.Song
	db := database.GetDB()
	result_err := db.Find(&Songs_from_db)
	if result_err.Error != nil { //check errors
		log.Printf("Error: Failed to find songs from DB: %v", result_err.Error)
		http.Error(w, "Error getting data from DB: "+result_err.Error.Error(), http.StatusInternalServerError)
		return
	}
	log.Printf("Info: Find %d songs from the database", len(Songs_from_db))

	for _, s := range Songs_from_db {
		// group sort
		if group != "" && !strings.Contains(strings.ToLower(s.Group), strings.ToLower(group)) {
			continue
		}
		// song sort
		if song != "" && !strings.Contains(strings.ToLower(s.Song), strings.ToLower(song)) {
			continue
		}

		filteredSongs = append(filteredSongs, s)
	}

	log.Printf("Debug: Filtered %d songs based on query parameters", len(filteredSongs))
	index_song_start := (page - 1) * limit
	index_song_end := index_song_start + limit

	if index_song_start > len(filteredSongs) {
		log.Printf("Error: Page out of range, page: %d ", page)
		http.Error(w, "Page out of range", http.StatusBadRequest)
		return
	}

	if index_song_end > len(filteredSongs) {
		index_song_end = len(filteredSongs)
	}

	log.Printf("Debug: Returning songs from index %d to %d", index_song_start, index_song_end)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(filteredSongs[index_song_start:index_song_end])
}

// @Summary Добавить песню
// @Description Добавляет новую песню в базу данных, предварительно получая дополнительную информацию с внешнего API
// @Tags Songs
// @Accept json
// @Produce json
// @Param song body models.Song true "Данные о песне"
// @Success 201 {object} models.Song "Песня успешно добавлена"
// @Failure 400 {string} string "Некорректный JSON-запрос"
// @Failure 404 {string} string "Информация о песне не найдена во внешнем API"
// @Failure 405 {string} string "Неверный метод запроса (требуется POST)"
// @Failure 500 {string} string "Ошибка сервера при обращении к API или сохранении в БД"
// @Router /api/AddSong [post]
func AddSongs(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodPost { // check POST
		log.Printf("Error: Invalid request method, expected POST but got %s", r.Method)
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	log.Println("Info: <AddSongs> endpoint...")

	var newSong models.Song
	err := json.NewDecoder(r.Body).Decode(&newSong)
	if err != nil {
		log.Printf("Error: Failed to decode JSON body: %v", err)
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}
	log.Printf("Debug: Incoming endpoint request body %+v", newSong)

	apiURL := fmt.Sprintf("%s/info?group=%s&song=%s", os.Getenv("API_URL"),
		url.QueryEscape(newSong.Group), url.QueryEscape(newSong.Song)) // url format

	resp, err := http.Get(apiURL) // REQUEST to your SWAGGER API

	if err != nil {
		log.Printf("Error: Failed to fetch song info from swagger api: %v", err)
		http.Error(w, "Failed to fetch song info from swagger api", http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close() // close after beginning request

	if resp.StatusCode != http.StatusOK {
		log.Printf("Error: Song info not found, received status code %d", resp.StatusCode)
		http.Error(w, "Song info not found", http.StatusNotFound)
		return
	}

	//struct for answer swagger api
	var songInfo struct {
		ReleaseDate string `json:"releaseDate"`
		Text        string `json:"text"`
		Link        string `json:"link"`
	}

	err = json.NewDecoder(resp.Body).Decode(&songInfo) // get response params
	if err != nil {
		log.Printf("Error: Failed to decode API response: %v", err)
		http.Error(w, "Invalid API response", http.StatusInternalServerError)
		return
	}

	log.Printf("Debug: Incoming request body from swagger api (/Info) %+v", songInfo)

	parsedDate, err := time.Parse("02.01.2006", songInfo.ReleaseDate) // Output date format from api swagger
	if err != nil {
		log.Printf("Error: Invalid date format from API: %v", err)
		http.Error(w, "Invalid date format from API", http.StatusInternalServerError)
		return
	}

	newSong.ReleaseDate = parsedDate
	newSong.Text = songInfo.Text
	newSong.Link = songInfo.Link

	db := database.GetDB()

	log.Printf("Debug: This will be added: %+v", songInfo)
	result_err := db.Create(&newSong) //added to db

	if result_err.Error != nil { //check errors
		log.Printf("Error: Failed to save song to database: %v", result_err.Error)
		http.Error(w, "Error saving song to database: "+result_err.Error.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated) // code 201
	json.NewEncoder(w).Encode(newSong)
}

// @Summary Удалить песню
// @Description Удаляет песню из базы данных по её ID
// @Tags Songs
// @Produce json
// @Param songId query int true "ID песни для удаления"
// @Success 200 {object} models.Song "Песня успешно удалена"
// @Failure 400 {string} string "Некорректный ID песни или отсутствует параметр songId"
// @Failure 404 {string} string "Песня не найдена"
// @Failure 405 {string} string "Неверный метод запроса (требуется DELETE)"
// @Failure 500 {string} string "Ошибка сервера при удалении песни"
// @Router /api/deleteSong [delete]
func DeleteSong(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodDelete {
		log.Printf("Error: Invalid request method. Expected DELETE, got %s", r.Method)
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	log.Println("Info: <DeleteSong> endpoint...")

	songId := r.URL.Query().Get("songId")
	db := database.GetDB()

	songIntId := 0

	if songId == "" {
		log.Println("Error: Missing songId parameter")
		http.Error(w, "Missing songId parameter", http.StatusBadRequest)
		return
	}

	if songId != "" {
		var err error
		songIntId, err = strconv.Atoi(songId)
		if err != nil { //check errors
			log.Printf("Error: Invalid id parameter: %v", err)
			http.Error(w, "Invalid id parameter", http.StatusBadRequest)
			return
		}
		log.Printf("Debug: songId converted to int: %d", songIntId)
	}

	var song models.Song
	if err := db.First(&song, songIntId).Error; err != nil {
		log.Printf("Error: Song with id %d not found", songIntId)
		http.Error(w, "Song not found", http.StatusNotFound)
		return
	}

	log.Printf("Info: Found song to delete: %+v", song)
	if err := db.Delete(&song).Error; err != nil {
		log.Printf("Error: Failed to delete song with id %d: error: %v", songIntId, err)
		http.Error(w, "Failed to delete song", http.StatusInternalServerError)
		return
	}

	log.Printf("Info: Song with id %d deleted successfully", songIntId)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(song)
}

// @Summary Редактировать песню
// @Description Обновляет указанные параметры песни в базе данных
// @Tags Songs
// @Accept json
// @Produce json
// @Param songId query int true "ID песни для редактирования"
// @Param paramsToEdit query string true "Список параметров для изменения (через запятую), например: group,song,releaseDate"
// @Param body body object true "JSON с новыми значениями полей"
// @Success 200 {object} models.Song "Песня успешно обновлена"
// @Failure 400 {string} string "Некорректный ID или параметры запроса"
// @Failure 404 {string} string "Песня не найдена"
// @Failure 405 {string} string "Неверный метод запроса (требуется PUT)"
// @Failure 500 {string} string "Ошибка сервера при обновлении"
// @Router /api/EditSong [put]
func EditSong(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodPut {
		log.Printf("Error: Invalid request method. Expected PUT, got %s", r.Method)
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}
	log.Println("Info: <EditSong> endpoint...")

	songId := r.URL.Query().Get("songId")
	songIntId := 0

	if songId == "" {
		log.Println("Error: Missing songId parameter")
		http.Error(w, "Missing songId parameter", http.StatusBadRequest)
		return
	}

	if songId != "" {
		var err error
		songIntId, err = strconv.Atoi(songId)
		if err != nil {
			log.Printf("Error: Invalid id parameter: %v", err)
			http.Error(w, "Invalid id parameter", http.StatusBadRequest)
			return
		}

	}
	log.Printf("Debug: songId converted to int: %d", songIntId)

	db := database.GetDB()
	var song models.Song
	if err := db.First(&song, songIntId).Error; err != nil {
		log.Printf("Error: Song with id %d not found", songIntId)
		http.Error(w, "Song not found", http.StatusNotFound)
		return
	}

	log.Printf("Info: Found song to edit: %+v", song)

	paramsToEdit := r.URL.Query().Get("paramsToEdit")
	if paramsToEdit == "" {
		log.Println("Error: Missing paramsToEdit parameter")
		http.Error(w, "Missing paramsToEdit parameter", http.StatusBadRequest)
		return
	}
	log.Printf("Debug: paramsToEdit parameter: %s", paramsToEdit)

	params := strings.Split(paramsToEdit, ",")

	// json decode params
	var requestData map[string]string
	if err := json.NewDecoder(r.Body).Decode(&requestData); err != nil {
		log.Printf("Error: Invalid JSON body: %v", err)
		http.Error(w, "Invalid JSON body", http.StatusBadRequest)
		return
	}

	log.Printf("Debug: Incoming request body: %+v", requestData)

	// update
	for _, field := range params {
		switch field {
		case "group":
			song.Group = requestData["group"]
		case "song":
			song.Song = requestData["song"]
		case "releaseDate":
			releaseDate, err := time.Parse("2006-01-02", requestData["releaseDate"])
			if err != nil {
				log.Printf("Error: Invalid releaseDate format. Expected YYYY-MM-DD: %v", err)
				http.Error(w, "Invalid releaseDate format. Use YYYY-MM-DD", http.StatusBadRequest)
				return
			}
			song.ReleaseDate = releaseDate
		case "text":
			song.Text = requestData["text"]
		case "link":
			song.Link = requestData["link"]
		}
	}

	// new update time
	song.UpdatedAt = time.Now()

	// save to db
	if err := db.Save(&song).Error; err != nil {
		log.Printf("Error: Failed to update song with id %d: %v", songIntId, err)
		http.Error(w, "Failed to update song", http.StatusInternalServerError)
		return
	}

	log.Printf("Info: Song with id %d updated successfully", songIntId)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(song)
}

// @Summary Получить текст песни постранично
// @Description Возвращает текст песни, разбитый на страницы по указанному лимиту строк
// @Tags Songs
// @Accept json
// @Produce json
// @Param songId query int true "ID песни"
// @Param page query int false "Номер страницы (по умолчанию 1)"
// @Param limit query int false "Количество строф на странице (по умолчанию 2)"
// @Success 200 {object} map[string]interface{} "Текст песни постранично"
// @Failure 400 {string} string "Некорректные параметры запроса"
// @Failure 404 {string} string "Песня не найдена или страница вне диапазона"
// @Failure 405 {string} string "Неверный метод запроса (требуется GET)"
// @Router /api/GetSongText [get]
func GetSongText(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		log.Printf("Error: Invalid request method. Expected GET, got %s", r.Method)
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}
	log.Println("Info: <GetSongText> endpoint...")

	// songId
	songId := r.URL.Query().Get("songId")
	if songId == "" {
		log.Println("Error: Missing songId parameter")
		http.Error(w, "Missing songId parameter", http.StatusBadRequest)
		return
	}
	log.Printf("Debug: songId: %s", songId)

	// id to int
	songIntId, err := strconv.Atoi(songId)
	if err != nil {
		log.Printf("Error: Invalid songId parameter: %v", err)
		http.Error(w, "Invalid songId parameter", http.StatusBadRequest)
		return
	}

	page, err := strconv.Atoi(r.URL.Query().Get("page"))
	if err != nil || page <= 0 {
		page = 1
	}
	limit, err := strconv.Atoi(r.URL.Query().Get("limit"))
	if err != nil || limit <= 0 {
		limit = 2
	}

	log.Printf("Debug: Page: %d, Limit: %d", page, limit)

	if page <= 0 {
		page = 1
	}
	if limit <= 0 {
		limit = 2
	}
	db := database.GetDB()
	var song models.Song

	if err := db.First(&song, songIntId).Error; err != nil {
		log.Printf("Error: Song with id %d not found", songIntId)
		http.Error(w, "Song not found", http.StatusNotFound)
		return
	}

	log.Printf("Info: Found song with id %d: %+v", songIntId, song)

	//split text
	verses := strings.Split(song.Text, "\n\n")

	start := (page - 1) * limit
	end := start + limit

	if start >= len(verses) {
		log.Println("Error: Page out of range")
		http.Error(w, "Page out of range", http.StatusNotFound)
		return
	}
	if end > len(verses) {
		end = len(verses)
	}

	// return verses
	response := map[string]interface{}{
		"song":    song.Song,
		"group":   song.Group,
		"verses":  verses[start:end],
		"page":    page,
		"perPage": limit,
		"total":   len(verses),
	}
	log.Printf("Info: Returning verses: %v from page %d", response["verses"], page)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)

}
