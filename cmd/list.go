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
	"io/ioutil"
	"net"
	"text/tabwriter"

	"github.com/spf13/cobra"

	"context"
	"log"
	"os"
	"strings"

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
		printClusterList()
	},
}

func init() {
	rootCmd.AddCommand(listCmd)
}

// Check if Red Hat VPN is connected
func checkDNS() {
	_, err := net.LookupIP("quicklab.upshift.redhat.com")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Please connect VPN!\n")
		os.Exit(1)
	}
}

// Collects clusters list
func getClustersList(url, path, tag string) (rows, links [][]string, headings []string) {
	checkDNS()
	dir, err := ioutil.TempDir("/tmp/", "quicklab")
	if err != nil {
		log.Fatal(err)
	}
	defer os.RemoveAll(dir)
	opts := append(chromedp.DefaultExecAllocatorOptions[:],
		chromedp.Flag("headless", false),
		chromedp.WindowSize(300, 300),
		chromedp.Flag("disable-gpu", false),
		chromedp.Flag("enable-automation", false),
		chromedp.Flag("disable-extensions", false),
		chromedp.UserDataDir(dir),
	)

	allocCtx, cancel := chromedp.NewExecAllocator(context.Background(), opts...)
	defer cancel()

	// create context
	ctx, cancel := chromedp.NewContext(allocCtx, chromedp.WithLogf(log.Printf))
	defer cancel()
	var ids []cdp.NodeID
	var body string

	if err := chromedp.Run(ctx,
		chromedp.Navigate(url),
		chromedp.Click(path, chromedp.NodeVisible),
		chromedp.NodeIDs(tag, &ids, chromedp.ByJSPath),
		chromedp.ActionFunc(func(ctx context.Context) error {
			var erro error
			body, erro = dom.GetOuterHTML().WithNodeID(ids[0]).Do(ctx)
			return erro
		}),
	); err != nil {
		log.Fatal(err)
	}

	var row []string

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
			})
			rows = append(rows, row)
			row = nil
		})
	})

	var link []string
	null := []string{"null"}
	links = append(links, null)
	doc.Find("a").Each(func(in int, hreflink *goquery.Selection) {
		href, _ := hreflink.Attr("href")
		link = append(link, href)
		links = append(links, link)
		link = nil
	})
	return
}

// Prints clusters list in table format
func printClusterList() {
	url := `https://quicklab-quicklab.apps.ocp-c1.prod.psi.redhat.com/login`
	path := `sharedclusters`
	tag := `document.querySelector("#main-container > div > main > div > section > article > div.pf-c-card__body")`
	name, _, headings := getClustersList(url, path, tag)
	singleName := name[1]
	rowLength := len(name)
	columnLength := len(singleName)
	w := new(tabwriter.Writer)
	w.Init(os.Stdout, 8, 8, 0, '\t', 0)
	fmt.Fprintf(w, "%s\t%s\t%s\t%s\t\n", strings.ToUpper(headings[0]), strings.ToUpper(headings[1]), strings.ToUpper(headings[2]), strings.ToUpper(headings[3]))
	for i := 1; i < rowLength; i++ {
		for j := 0; j < columnLength; j++ {
			fmt.Fprintf(w, "%s\t", name[i][j])
		}
		fmt.Fprintf(w, "\n")
	}
	w.Flush()
}
