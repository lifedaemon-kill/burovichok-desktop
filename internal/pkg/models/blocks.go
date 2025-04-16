package models

import "time"

// BlockOne модель для файла  Рзаб и Тзаб.xlsx.
type BlockOne struct {
	Timestamp   time.Time `xlsx:"Дата, время"`                     // колонка A
	Pressure    float64   `xlsx:"Рзаб на глубине замера, кгс/см2"` // колонка B
	Temperature float64   `xlsx:"Tзаб на глубине замера, оС"`      // колонка C
}

type BlockTwo struct {
}

type BlockThree struct {
}

type BlockFour struct {
}

type BlockFive struct {
}

type BlockSix struct {
}

type BlockSeven struct {
}
