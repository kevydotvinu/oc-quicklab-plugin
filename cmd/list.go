/*
Copyright © 2021 NAME HERE <EMAIL ADDRESS>

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package cmd

import (
	"fmt"

	"github.com/spf13/cobra"

	"context"
	"log"
	"os"
	"strings"
	"text/tabwriter"

	"github.com/PuerkitoBio/goquery"
	"github.com/chromedp/cdproto/cdp"
	"github.com/chromedp/cdproto/dom"
	"github.com/chromedp/chromedp"
)

// listCmd represents the list command
var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List the quicklab shared clusters",
	Long:  `This command gives the list of available clusters in the quicklab shared environment`,
	Run: func(cmd *cobra.Command, args []string) {
		getClustersList()
	},
}

func init() {
	rootCmd.AddCommand(listCmd)
}

func getClustersList() {
	opts := append(chromedp.DefaultExecAllocatorOptions[:],
		chromedp.Flag("headless", false),
		chromedp.WindowSize(300, 300),
		chromedp.Flag("disable-gpu", false),
		chromedp.Flag("enable-automation", false),
		chromedp.Flag("disable-extensions", false),
		chromedp.UserDataDir(os.Getenv("HOME")+"/.config/chromium"),
	)

	allocCtx, cancel := chromedp.NewExecAllocator(context.Background(), opts...)
	defer cancel()

	// create context
	ctx, cancel := chromedp.NewContext(allocCtx, chromedp.WithLogf(log.Printf))
	defer cancel()
	var ids []cdp.NodeID
	var body string
	if err := chromedp.Run(ctx,
		chromedp.Navigate(`https://quicklab-quicklab.apps.ocp-c1.prod.psi.redhat.com/login`),
		chromedp.Click(`sharedclusters`, chromedp.NodeVisible),
		chromedp.NodeIDs(`document.querySelector("#main-container > div > main > div > section > article > div.pf-c-card__body")`, &ids, chromedp.ByJSPath),
		chromedp.ActionFunc(func(ctx context.Context) error {
			var erro error
			body, erro = dom.GetOuterHTML().WithNodeID(ids[0]).Do(ctx)
			return erro
		}),
	); err != nil {
		log.Fatal(err)
	}

	var headings, row []string
	var rows [][]string
	var rowLength int
	var rowsLength int

	doc, err := goquery.NewDocumentFromReader(strings.NewReader(body))
	if err != nil {
		fmt.Println("No url found")
		log.Fatal(err)
	}
	// Find each table
	doc.Find("table").Each(func(index int, tablehtml *goquery.Selection) {
		tablehtml.Find("tr").Each(func(indextr int, rowhtml *goquery.Selection) {
			rowhtml.Find("th").Each(func(indexth int, tableheading *goquery.Selection) {
				headings = append(headings, tableheading.Text())
			})
			rowhtml.Find("td").Each(func(indexth int, tablecell *goquery.Selection) {
				row = append(row, tablecell.Text())
				rowLength = len(row)
			})
			rows = append(rows, row)
			row = nil
		})
	})

	rowsLength = len(rows)
	link := []string{}
	links := [][]string{{"null"}}
	doc.Find("a").Each(func(in int, hreflink *goquery.Selection) {
		href, _ := hreflink.Attr("href")
		link = append(link, href)
		links = append(links, link)
		link = nil
	})

	headings = append(headings, "Links")

	w := new(tabwriter.Writer)
	w.Init(os.Stdout, 8, 8, 0, '\t', 0)
	// defer w.Flush()
	fmt.Fprintf(w, "%s\t%s\t%s\t%s\t\n", headings[0], headings[1], headings[2], headings[3])
	for i := 1; i < rowsLength; i++ {
		for j := 0; j < rowLength; j++ {
			fmt.Fprintf(w, "%s\t", rows[i][j])
		}
		fmt.Fprintf(w, "\n")
	}
	w.Flush()
}
