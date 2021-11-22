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
	"context"
	"fmt"
	"log"
	"os"
	"os/exec"
	"regexp"
	"strings"
	"text/tabwriter"

	"github.com/PuerkitoBio/goquery"
	"github.com/chromedp/cdproto/cdp"
	"github.com/chromedp/cdproto/dom"
	"github.com/chromedp/chromedp"
	"github.com/spf13/cobra"
)

// loginCmd represents the login command
var loginCmd = &cobra.Command{
	Use:   "login",
	Short: "Authenticate the quicklab shared cluster using oc binary.",
	Long:  `This command fetches the credentials from the portal and authenticates the quicklab shared cluster using oc binary. The username will be kubeadmin.`,
	Run: func(cmd *cobra.Command, args []string) {
		loginCluster()
	},
}

func init() {
	rootCmd.AddCommand(loginCmd)
	loginCmd.Flags().StringP("cluster", "c", "", "Pass quicklab shared cluster name")
	loginCmd.MarkFlagRequired("cluster")
}

// Maps cluster name and its link
func clusterNameAndLinks() map[string]string {
	url := `https://quicklab-quicklab.apps.ocp-c1.prod.psi.redhat.com/login`
	path := `sharedclusters`
	tag := `document.querySelector("#main-container > div > main > div > section > article > div.pf-c-card__body")`
	name, links, _ := getClustersList(url, path, tag)
	rowLength := len(name)
	nameLinkMap := make(map[string]string)
	for i := 1; i < rowLength; i++ {
		j := name[i][0]
		k := links[i][0]
		nameLinkMap[j] = k
	}
	return nameLinkMap
}

func getPath() string {
	clusterNameAndLink := clusterNameAndLinks()
	if val, ok := clusterNameAndLink[os.Args[3]]; ok {
		path := val
		return path
	} else {
		fmt.Println("The cluster named", os.Args[3], "does not exist in the list")
		os.Exit(1)
	}
	return ""
}

// Collects HTML body
func getHtmlBody(url, path, tag string) (body string) {
	checkDNS()
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

	if err := chromedp.Run(ctx,
		chromedp.Navigate(url),
		chromedp.Click(`sharedclusters`, chromedp.NodeVisible),
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
	return body
}

// Parse username, password and URL from HTML body
func parseHtmlBody(body string) (username string, password string, server string) {

	var s string
	username = "kubeadmin"
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(body))
	if err != nil {
		fmt.Println("HTML body not found")
		log.Fatal(err)
	}

	doc.Find("body").Each(func(indexdiv int, bodyhtml *goquery.Selection) {
		s = bodyhtml.Text()
	})
	r1, _ := regexp.Compile("([A-Za-z0-9]{5})+-+([A-Za-z0-9]{5})+-+([A-Za-z0-9]{5})+-+([A-Za-z0-9]{5})")
	r2, _ := regexp.Compile(`(apps)\.([\w_-]+(?:(?:\.[\w_-]+)+))([\w.,@?^=%&:\/~+#-]*[\w@?^=%&\/~+#-])\.(com)`)
	password = r1.FindString(s)
	server = r2.FindString(s)
	server = strings.Replace(server, "apps", "api", -1)
	server = `https://` + server + `:6443`
	return username, password, server
}

// Logging into the cluster
func loginCluster() {
	url := `https://quicklab-quicklab.apps.ocp-c1.prod.psi.redhat.com/login`
	path := getPath()
	tag := `document.querySelector("#main-container > div > main > div > section:nth-child(2) > div > div:nth-child(4) > article > div.pf-c-card__body")`
	body := getHtmlBody(url, path, tag)
	username, password, server := parseHtmlBody(body)
	console := strings.Replace(server, "api", "console-openshift-console.apps", -1)
	console = strings.Replace(console, ":6443", "", -1)
	if password == "" && server != "" {
		username = "administrator"
		password = "Password@123"
		server = strings.Replace(server, ":6443", "", -1)
		server = strings.Replace(server, "api", "openshift", -1)
		console = strings.Replace(server, "api", "openshift", -1)
	}

	w := new(tabwriter.Writer)
	w.Init(os.Stdout, 8, 8, 0, '\t', 0)
	fmt.Fprintf(w, "%s\t%s\t%s\t\n", strings.ToUpper("username"), strings.ToUpper("password"), strings.ToUpper("console"))
	fmt.Fprintf(w, "%s\t%s\t%s\t\n\n", strings.ToLower(username), password, strings.ToLower(console))
	w.Flush()

	cmd := exec.Command("oc", "login", "--insecure-skip-tls-verify=true", "-u", username, "-p", password, server)
	stdout, err := cmd.Output()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(string(stdout))
}
