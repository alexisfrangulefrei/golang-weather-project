package main

import (
	"fmt"
	"time"
)

func ParseStringToTime(stringToParse string, stringLayout string) time.Time {
	formattedTime, err := time.Parse(stringLayout, stringToParse)

	if err != nil {
		fmt.Println("Error:", err)
		panic(err)
	}

	return formattedTime
}
