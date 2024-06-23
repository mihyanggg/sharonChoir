package main

import "fmt"

func clearScreen() {
	fmt.Print("\033[2J")
	fmt.Print("\033[H")
}

func moveCursor(row, col int) {
	fmt.Printf("\033[%d;%dH", row, col)
}

//func songList() {
//	clearScreen()
//	fmt.Print("Enter the song of the month.")
//	fmt.Print("yy    mm    dd    no    page    name")
//}

func allSongList(songs []readSong) {
	//cnt := len(songs)
	//clearScreen()
	fmt.Println("----- Song List ----")
	for i, song := range songs {
		fmt.Printf("%d.   \n", i+1)
		fmt.Printf("%-9s: %02s , \n", "Year", song.y)
		fmt.Printf("%-9s: %02s , \n", "Month", song.m)
		fmt.Printf("%-9s: %02s , \n", "Day", song.d)
		fmt.Printf("%-9s: %s , \n", "Part", song.dn)
		fmt.Printf("%-9s: %-20s , \n", "BookName", song.s_bNm)
		fmt.Printf("%-9s: %02s , \n", "SongNo", song.s_sNo)
		fmt.Printf("%-9s: %-20s , \n", "SongName", song.s_sNm)
		fmt.Printf("%-9s: %02s \n", "SongPage", song.s_sPg)
		fmt.Printf("%-9s: %s , \n", "mainUrl", song.s_mainUrl)
		fmt.Printf("%-9s: %s , \n", "allUrl", song.s_allUrl)
		fmt.Printf("%-9s: %s , \n", "sUrl", song.s_sUrl)
		fmt.Printf("%-9s: %s , \n", "aUrl", song.s_aUrl)
		fmt.Printf("%-9s: %s , \n", "tUrl", song.s_tUrl)
		fmt.Printf("%-9s: %s , \n", "bUrl", song.s_bUrl)
	}
}

func printSongList(songs []readSong) {
	//cnt := len(songs)
	//clearScreen()
	fmt.Printf("----- Song List ----\n\n")
	for i, song := range songs {
		if song.dn == "Day" {
			fmt.Printf("%s년 %s월 %s일 %d주차 찬양곡 올려드립니다 ☀️\n", song.y, song.m, song.d, i+1)
		} else {
			fmt.Printf("%s년 %s월 %s일 %d주차 저녁 찬양곡 올려드립니다 ☀️\n", song.y, song.m, song.d, i+1)
		}

		fmt.Printf("👉 %s. %s - %s(p%s) 👈\n", song.s_sNo, song.s_sNm, song.s_bNm, song.s_sPg)
		fmt.Printf("%s\n", song.s_mainUrl)
		fmt.Println()
		fmt.Printf("합창 🎼 : %s\n", song.s_allUrl)
		fmt.Printf("소프라노 🎼 : %s\n", song.s_sUrl)
		fmt.Printf("알토 🎼 : %s\n", song.s_aUrl)
		fmt.Printf("테너 🎼 : %s\n", song.s_tUrl)
		fmt.Printf("베이스 🎼 : %s\n", song.s_bUrl)
		fmt.Printf("________________________________________\n\n")
	}
}
