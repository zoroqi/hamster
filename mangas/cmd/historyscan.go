package cmd

import (
	"encoding/csv"
	"errors"
	"fmt"
	"github.com/spf13/cobra"
	"github.com/zoroqi/rubbish/mangas/data"
	"github.com/zoroqi/rubbish/mangas/website"
	"io"
	"os"
	"strconv"
	"time"
)

var historyScanCmd = &cobra.Command{
	Use:   "historyscan",
	Short: "history scan",
	Long:  `history scan`,
	RunE:  historyScan,
}
var (
	of       string
	starPage int
	endPage  int
)

func init() {
	RootCmd.AddCommand(historyScanCmd)
	historyScanCmd.Flags().StringVarP(&outfile, "outfile", "f", "", "outfile")
	historyScanCmd.Flags().IntVarP(&starPage, "starPage", "s", 0, "starPage")
	historyScanCmd.Flags().IntVarP(&endPage, "endPage", "e", 0, "endPage")
}

func historyScan(cmd *cobra.Command, args []string) error {
	if outfile == "" {
		return errors.New("outfile is empty")
	}
	f, err := os.OpenFile(outfile, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		return err
	}
	mhg := &website.Mhg{}
	for i := starPage; i <= endPage; i++ {
		bs, err := mhg.HistoryList("https://www.manhuagui.com/list/index_p%s.html", strconv.Itoa(i), strconv.Itoa(i-1))
		if err != nil {
			return fmt.Errorf("download err, page %d, err %w", i, err)
		}
		mgs, err := mhg.ListParse2(bs)
		if err != nil {
			return fmt.Errorf("parser err, page %d, err %w", i, err)
		}
		fmt.Printf("page %d, len %d\n", i, len(mgs))
		if outfile != "" {
			if err := writeCsv(mgs, f); err != nil {
				return fmt.Errorf("write err, page %d, err %w", i, err)
			}
		}
		time.Sleep(time.Duration(crawlInterval) * time.Second)
	}
	return nil
}

func writeCsv(mgs []data.Manga, f io.Writer) error {
	writer := csv.NewWriter(f)
	for _, manga := range mgs {
		err := writer.Write(manga.ToSlice())
		if err != nil {
			fmt.Printf("marshal err, %+v, %v\n", manga, err)
			continue
		}
	}
	writer.Flush()
	return nil
}
