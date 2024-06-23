package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"

	//"github.com/360EntSecGroup-Skylar/excelize"
	"github.com/xuri/excelize/v2"
)

// Read Sheet
type readSong struct {
	y  string //int
	m  string //int
	d  string //int
	dn string //bool // 0: day, 1: night

	s_bNm string // book name
	s_sNo string //int
	s_sNm string
	s_sPg string //int

	s_mainUrl string // main
	s_allUrl  string // all
	s_sUrl    string // soprano
	s_aUrl    string // alto
	s_tUrl    string // tenor
	s_bUrl    string // base
}

// Set Sheet
type setSong struct {
	week  string
	dn    string
	s_bNm string
	s_sNm string
	s_sPg string
}

func checkWeek() int {
	var w int
	fmt.Print("Weeks: ")
	_, err := fmt.Scan(&w)
	if err != nil {
		fmt.Println("Input Error: ", err)
	}
	return w
}

/*
func setSongList() []readSong {

	var songs [ ]readSong
	var pt bool
	var yy, mm, dd, no, pg int
	var bNm, nm, mainUrl, allUrl, sUrl, aUrl, tUrl, bUrl string
	var w int

	w = checkWeek()

	reader := bufio.NewReader(os.Stdin)
	i := 1

	//songList()
	for i <= w {

		fmt.Printf("%-10s: ", "yy") // %-10s left , %10s right
		fmt.Scan(&yy)
		//_, err := fmt.Scan(&yy)
		//if err != nil {
		//	fmt.Println("Enter only number")
		//}
		//fmt.Printf("%-10s: ", "yy")
		//input, _ := reader.ReadString('\n')
		//y, _ = strconv.Atoi(strings.TrimSpace(input))
		//yy = fmt.Sprintf("%02d", y)

		fmt.Printf("%-10s: ", "mm")
		fmt.Scan(&mm)

		fmt.Printf("%-10s: ", "dd")
		fmt.Scan(&dd)

		fmt.Printf("%-25s: ", "0: morning , 1: evening")
		fmt.Scan(&pt)

		reader.ReadString('\n')

		fmt.Printf("%-10s: ", "book name")
		bNm, _ = reader.ReadString('\n')
		bNm = strings.TrimSpace(bNm)

		fmt.Printf("%-10s: ", "number")
		fmt.Scan(&no)

		//reader.ReadString('\n')

		fmt.Printf("%-10s: ", "name")
		nm, _ = reader.ReadString('\n')
		nm = strings.TrimSpace(nm)

		//reader.ReadString('\n')

		fmt.Printf("%-10s: ", "page")
		fmt.Scan(&pg)

		fmt.Printf("%-10s: ", "Main Url")
		mainUrl, _ = reader.ReadString('\n')
		mainUrl = strings.TrimSpace(mainUrl)

		fmt.Printf("%-10s: ", "All Url")
		allUrl, _ = reader.ReadString('\n')
		allUrl = strings.TrimSpace(allUrl)

		fmt.Printf("%-10s: ", "soprano")
		sUrl, _ = reader.ReadString('\n')
		sUrl = strings.TrimSpace(sUrl)

		fmt.Printf("%-10s: ", "alto")
		aUrl, _ = reader.ReadString('\n')
		aUrl = strings.TrimSpace(aUrl)

		fmt.Printf("%-10s: ", "tenor")
		tUrl, _ = reader.ReadString('\n')
		tUrl = strings.TrimSpace(tUrl)

		fmt.Printf("%-10s: ", "base")
		bUrl, _ = reader.ReadString('\n')
		bUrl = strings.TrimSpace(bUrl)

		asong := song{y: yy, m: mm, d: dd, part: pt, s_bNm: bNm, s_sNo: no, s_sNm: nm, s_sPg: pg, s_mainUrl: mainUrl, s_allUrl: allUrl, s_sUrl: sUrl, s_aUrl: aUrl, s_tUrl: tUrl, s_bUrl: bUrl}
		songs = append(songs, asong)

		i++
	}

	//allSongList(songs)
	printSongList(songs)

	return songs
}
*/

func readSongsOfExcel(name string) ([]readSong, error) {
	var songs []readSong

	// Open the Excel file
	f, err := excelize.OpenFile(name)
	if err != nil {
		return nil, fmt.Errorf("Excel 파일을 열 수 없습니다: %v", err)
	}
	defer func() {
		if err := f.Close(); err != nil {
			log.Fatalf("Excel 파일을 닫을 수 없습니다: %v", err)
		}
	}()

	// 시트 이름 지정
	sheetName := "Read"

	// 시트의 모든 행 읽기
	rows, err := f.GetRows(sheetName)
	if err != nil {
		return nil, fmt.Errorf("시트를 읽을 수 없습니다: %v", err)
	}

	// 첫 번째 행은 헤더이므로 건너뜁니다
	for i, row := range rows {
		if i == 0 {
			continue
		}

		if len(row) < 12 {
			log.Fatalf("행에 필요한 열이 부족합니다: %v", row)
		}

		date, err := time.Parse("06/1/2", row[0])
		if err != nil {
			log.Fatalf("날짜 형식이 잘못되었습니다: %v", err)
		}

		//yyyy := fmt.Sprintf("%04d", date.Year())
		//mm :=   fmt.Sprintf("%02d", date.Month())
		//dd :=   fmt.Sprintf("%02d", date.Day())

		s := readSong{
			y:         fmt.Sprintf("%04d", date.Year()),
			m:         fmt.Sprintf("%02d", date.Month()),
			d:         fmt.Sprintf("%02d", date.Day()),
			dn:        row[1],
			s_bNm:     row[2],
			s_sNo:     row[3],
			s_sNm:     row[4],
			s_sPg:     row[5],
			s_mainUrl: row[6],
			s_allUrl:  row[7],
			s_sUrl:    row[8],
			s_aUrl:    row[9],
			s_tUrl:    row[10],
			s_bUrl:    row[11],
		}

		songs = append(songs, s)
	}

	//printSongList(songs)

	return songs, nil
}

// YouTube API 관련 설정
const apiKey = "AIzaSyDGyhP_a7hscYkmB-8b0zZFqjv9ZZownDQ"
const searchURL = "https://www.googleapis.com/youtube/v3/search"

func searchYouTube(query string) (string, error) {
	url := fmt.Sprintf("%s?part=snippet&q=%s&key=%s", searchURL, query, apiKey)

	resp, err := http.Get(url)
	if err != nil {
		return "", fmt.Errorf("YouTube 검색 요청 실패: %v", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("응답 읽기 실패: %v", err)
	}

	var result map[string]interface{}
	if err := json.Unmarshal(body, &result); err != nil {
		return "", fmt.Errorf("JSON 언마샬 실패: %v", err)
	}

	if items, ok := result["items"].([]interface{}); ok {
		for _, item := range items {
			if id, ok := item.(map[string]interface{})["id"].(map[string]interface{}); ok {
				if videoId, ok := id["videoId"].(string); ok {
					return fmt.Sprintf("https://www.youtube.com/watch?v=%s", videoId), nil
				}
			}
		}
	}

	return "", fmt.Errorf("적합한 동영상 URL을 찾을 수 없습니다.")
}

func updateExcelFile(name string, sheetName string, songs []readSong) error {
	f, err := excelize.OpenFile(name)
	if err != nil {
		return fmt.Errorf("Excel 파일을 열 수 없습니다: %v", err)
	}
	defer func() {
		if err := f.Close(); err != nil {
			log.Fatalf("Excel 파일을 닫을 수 없습니다: %v", err)
		}
	}()

	for i, s := range songs {
		row := i + 2 // 헤더를 스킵하고 2번째 행부터 시작
		f.SetCellValue(sheetName, fmt.Sprintf("G%d", row), s.s_mainUrl)
		f.SetCellValue(sheetName, fmt.Sprintf("H%d", row), s.s_allUrl)
		f.SetCellValue(sheetName, fmt.Sprintf("I%d", row), s.s_sUrl)
		f.SetCellValue(sheetName, fmt.Sprintf("J%d", row), s.s_aUrl)
		f.SetCellValue(sheetName, fmt.Sprintf("K%d", row), s.s_tUrl)
		f.SetCellValue(sheetName, fmt.Sprintf("L%d", row), s.s_bUrl)
	}

	if err := f.Save(); err != nil {
		return fmt.Errorf("Excel 파일 저장 실패: %v", err)
	}

	return nil
}

func readExcelSheet(name string, sheetName string) (interface{}, error) {
	f, err := excelize.OpenFile(name)
	if err != nil {
		return nil, fmt.Errorf("Excel 파일을 열 수 없습니다: %v", err)
	}
	defer func() {
		if err := f.Close(); err != nil {
			log.Fatalf("Excel 파일을 닫을 수 없습니다: %v", err)
		}
	}()

	rows, err := f.GetRows(sheetName)
	if err != nil {
		return nil, fmt.Errorf("시트를 읽을 수 없습니다: %v", err)
	}

	if sheetName == "Set" {
		var songs []setSong
		for i, row := range rows {
			if i == 0 {
				continue // 헤더 스킵
			}
			if len(row) < 5 {
				continue // 부족한 열 스킵
			}
			s := setSong{
				week:  row[0],
				dn:    row[1],
				s_bNm: row[2],
				s_sNm: row[3],
				s_sPg: row[4],
			}
			songs = append(songs, s)
		}
		return songs, nil
	} else if sheetName == "Read" {
		var songs []readSong
		for i, row := range rows {
			if i == 0 {
				continue // 헤더 스킵
			}
			if len(row) < 12 {
				continue // 부족한 열 스킵
			}

			date, err := time.Parse("06/1/2", row[0])
			if err != nil {
				log.Fatalf("날짜 형식이 잘못되었습니다: %v", err)
			}

			s := readSong{
				y:         fmt.Sprintf("%04d", date.Year()),
				m:         fmt.Sprintf("%02d", date.Month()),
				d:         fmt.Sprintf("%02d", date.Day()),
				dn:        row[1],
				s_bNm:     row[2],
				s_sNo:     row[3],
				s_sNm:     row[4],
				s_sPg:     row[5],
				s_mainUrl: row[6],
				s_allUrl:  row[7],
				s_sUrl:    row[8],
				s_aUrl:    row[9],
				s_tUrl:    row[10],
				s_bUrl:    row[11],
			}
			songs = append(songs, s)
		}
		return songs, nil
	} else {
		return nil, fmt.Errorf("알 수 없는 시트 이름: %s", sheetName)
	}
}

func newSong() {
	setSheet := "Set"
	readSheet := "Read"
	filePath := "./cellDir/sharonSongList.xlsx"

	// Set 시트 읽기
	setSongsInterface, err := readExcelSheet(filePath, setSheet)
	if err != nil {
		log.Fatalf("Set 시트 읽기 실패: %v", err)
	}
	setSongs := setSongsInterface.([]setSong)

	// Read 시트 읽기
	readSongsInterface, err := readExcelSheet(filePath, readSheet)
	if err != nil {
		log.Fatalf("Read 시트 읽기 실패: %v", err)
	}
	readSongs := readSongsInterface.([]readSong)

	for _, setSong := range setSongs {
		query := fmt.Sprintf("%s %s 합창", setSong.s_bNm, setSong.s_sNm)
		mainUrl, err := searchYouTube(query)
		if err != nil {
			log.Printf("YouTube 검색 실패 (합창): %v", err)
		}

		query = fmt.Sprintf("%s %s 소프라노", setSong.s_bNm, setSong.s_sNm)
		sopranoUrl, err := searchYouTube(query)
		if err != nil {
			log.Printf("YouTube 검색 실패 (소프라노): %v", err)
		}

		query = fmt.Sprintf("%s %s 알토", setSong.s_bNm, setSong.s_sNm)
		altoUrl, err := searchYouTube(query)
		if err != nil {
			log.Printf("YouTube 검색 실패 (알토): %v", err)
		}

		query = fmt.Sprintf("%s %s 테너", setSong.s_bNm, setSong.s_sNm)
		tenorUrl, err := searchYouTube(query)
		if err != nil {
			log.Printf("YouTube 검색 실패 (테너): %v", err)
		}

		query = fmt.Sprintf("%s %s 베이스", setSong.s_bNm, setSong.s_sNm)
		bassUrl, err := searchYouTube(query)
		if err != nil {
			log.Printf("YouTube 검색 실패 (베이스): %v", err)
		}

		for i, readSong := range readSongs {
			if readSong.s_bNm == setSong.s_bNm && readSong.s_sNm == setSong.s_sNm {
				readSongs[i].s_allUrl = mainUrl
				readSongs[i].s_sUrl = sopranoUrl
				readSongs[i].s_aUrl = altoUrl
				readSongs[i].s_tUrl = tenorUrl
				readSongs[i].s_bUrl = bassUrl
				break
			}
		}
	}

	if err := updateExcelFile(filePath, readSheet, readSongs); err != nil {
		log.Fatalf("Read 시트 업데이트 실패: %v", err)
	}

	fmt.Println("Excel 파일이 성공적으로 업데이트되었습니다.")
}
