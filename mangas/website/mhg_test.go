package website

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"testing"
)

func TestMhg_ListParse2(t *testing.T) {
	mhg := Mhg{}
	bs, err := os.ReadFile("./list.html")
	if err != nil {
		t.Fatal(err)
	}
	mgs, err := mhg.ListParse2(bs)
	if err != nil {
		t.Fatal(err)
	}
	for _, mg := range mgs {
		t.Log(mg)
	}
}

func TestMhg_ParseManga(t *testing.T) {
	mhg := Mhg{}
	bs, err := os.ReadFile("./manga.html")
	if err != nil {
		t.Fatal(err)
	}
	mg, err := mhg.ParseManga(bs)
	if err != nil {
		t.Fatal(err)
	}
	empJSON, err := json.MarshalIndent(mg, "", "  ")
	if err != nil {
		log.Fatalf(err.Error())
	}
	fmt.Printf("%s\n", empJSON)
}
