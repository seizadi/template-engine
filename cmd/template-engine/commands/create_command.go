package commands

import (
	"fmt"
	"github.com/seizadi/cmdb/client"
	"github.com/seizadi/template-engine/engine"
	"github.com/spf13/viper"
	"strings"
	
	"github.com/codegangsta/cli"
)

const app = "cmdb"
const fileUrl = "https://github.com/seizadi/cmdb/archive/v0.0.4.zip"
const srcFile = "./tmp/repo.zip"
const dstDir = "./tmp"

func createManifest(c *cli.Context) {
	
	//err := util.DownloadFile(srcFile, fileUrl )
	//if err != nil {
	//	log.Fatalf("Failed to download file %s\n", err)
	//}
	//files := util.UnzipFiles(srcFile, dstDir)
	//
	//path := getAppRepo(app, files)
	//
	//if len(path) == 0 {
	//	fmt.Printf("Did not find repo for app %s\n", app)
	//	return
	//}
	
	path := "tmp/cmdb-0.0.4/repo/cmdb"
	
	fmt.Printf("Path repo for app %s\n", path)
	
	conn, err := client.GetConn(viper.GetString("server.cmdbHost"))
	if err != nil {
		return
	}
	
	resp, err := client.GetRegions(conn)
	
	if (err != nil) {
		fmt.Printf("Error getting Regions %s\n", err)
		return
	}
	
	engine.CreateTemplate(resp)
	
	return
}

func getAppRepo(app string, files []string) string {
	for _, f := range files {
		if strings.HasSuffix(f, "/repo/"+ app) {
			fmt.Printf("found file %s\n", f)
			return f
		}
	}
	
	return ""
}