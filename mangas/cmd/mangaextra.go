package cmd

import (
	"encoding/csv"
	"errors"
	"fmt"
	"github.com/spf13/cobra"
	"github.com/zoroqi/hamster/mangas/data"
	"github.com/zoroqi/hamster/mangas/website"
	"os"
	"time"
)

var mangaextraCmd = &cobra.Command{
	Use:   "extra",
	Short: `manga extra`,
	Long:  `manga extra`,
	RunE:  mangaExtra,
}
var (
	inputCsvFile         string
	outputCsvFile        string
	outputChapterCsvFile string
	startLine            uint
	endLine              uint
)

func init() {
	RootCmd.AddCommand(mangaextraCmd)
	mangaextraCmd.Flags().StringVar(&inputCsvFile, "icsv", "", "input csv")
	mangaextraCmd.Flags().StringVar(&outputCsvFile, "ocsv", "", "output csv")
	mangaextraCmd.Flags().StringVar(&outputChapterCsvFile, "ochapter", "", "chapter csv")
	mangaextraCmd.Flags().UintVar(&startLine, "start", 0, "start line")
	mangaextraCmd.Flags().UintVar(&endLine, "end", 0, "end line")

}

func mangaExtra(cmd *cobra.Command, args []string) error {
	if inputCsvFile == "" || outputCsvFile == "" || outputChapterCsvFile == "" {
		return errors.New("icsv or ocsv or ochapter is empty")
	}
	if endLine == 0 {
		return errors.New("endLine is 0")
	}
	icsv, err := os.OpenFile(inputCsvFile, os.O_RDONLY, 0644)
	if err != nil {
		return err
	}
	defer icsv.Close()
	ocsv, err := os.OpenFile(outputCsvFile, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer ocsv.Close()
	chaperCsv, err := os.OpenFile(outputChapterCsvFile, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer ocsv.Close()
	csvr := csv.NewReader(icsv)
	csvw := csv.NewWriter(ocsv)
	chaperCsvW := csv.NewWriter(chaperCsv)
	defer csvw.Flush()
	defer chaperCsvW.Flush()
	for i := uint(0); ; i++ {
		line, err := csvr.Read()
		if err != nil {
			break
		}
		if i < startLine || i >= endLine {
			continue
		}
		mg := data.ParseSlice(line)
		if err != nil {
			return err
		}
		if mg.Mid == "" {
			mg.Mid = website.ParseMhgMid(mg.Link)
		}
		mhg := website.Mhg{}
		bs, err := mhg.Manga(fmt.Sprintf("https://www.manhuagui.com/comic/%s/", mg.Mid))
		fmt.Printf("%d %s", i, mg.Title)
		if err != nil {
			return err
		}
		mge, err := mhg.ParseManga(bs)
		if err != nil {
			return err
		}
		fmt.Printf(", chapter size:%d\n", len(mge.ChapterList))
		mge.Manga = mg
		err = csvw.Write(mge.ToSlice())
		if err != nil {
			return err
		}
		for _, c := range mge.ChapterList {
			c.Mid = mg.Mid
			err = chaperCsvW.Write(c.ToSlice())
			if err != nil {
				return err
			}
		}
		csvw.Flush()
		chaperCsvW.Flush()
		time.Sleep(time.Duration(crawlInterval) * time.Second)
	}
	return nil
}
