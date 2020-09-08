// https://github.com/zou2699/clean-harbor
package main

import (
	"flag"
	"fmt"
	"log"
	"sort"

	"clean-harbor/pkg/harbor"
)

var (
	url         string
	user        string
	password    string
	projectName string
	keepNum     int
	help        bool
)

func init() {
	flag.BoolVar(&help, "h", false, "help message")
	flag.StringVar(&url, "url", "", "harbor地址")
	flag.StringVar(&user, "user", "", "harbor账号")
	flag.StringVar(&password, "password", "", "harbor密码")
	flag.StringVar(&projectName, "projectName", "", "projectName")
	flag.IntVar(&keepNum, "keepNum", 5, "每个repo保留的tag个数")
}

func main() {
	flag.Parse()
	if help {
		flag.Usage()
		return
	}
	if url == "" || user == "" || password == "" || projectName == "" {
		flag.Usage()
		return
	}
	harborClient := harbor.NewClient(user, password, url)

	// 获取某个project下面所有repo的名字
	repoNames, err := harborClient.GetRepoNames(projectName)
	if err != nil {
		log.Fatalf("GetRepoNames: %s\n", err)
	}

	var size int64
	for _, repoName := range repoNames {
		// 根据repo的获取其下所有的tag
		tags, err := harborClient.GetRepoTags(repoName)
		if err != nil {
			panic(err)
		}

		//tags内容包含3个元素，格式类似于：[{2190824484 v1 2020-08-28 02:14:13.009841239 +0000 UTC} {53504535 v2 2020-08-14 00:36:48.610531148 +0000 UTC}]
		if len(tags) > keepNum { //tag数量大于keepNum才需要执行删除
			fmt.Printf("repo: %s 当前的tag数量为: %d，需要保留的tag数量是%d，开始执行删除\n", repoName, len(tags), keepNum)
			sort.Sort(tags)                ////自定义排序，根据tag的创建时间戳正序排列
			toDeleteTags := tags[keepNum:] //需要删除的tag切片
			// fmt.Println(toDeleteTags)
			for _, tag := range toDeleteTags {
				fmt.Printf("     删除image: %s:%s, 创建时间为: %s\n", repoName, tag.Name, tag.Created)
				err := harborClient.DeleteRepoTag(repoName, tag.Name)
				if err != nil {
					fmt.Printf("image: %s:%s DeleteRepoTag: %s\n", repoName, tag.Name, err)
					continue
				}
				size += tag.Size
			}
			fmt.Printf("repo: %s共清理: %.2f MB\n", repoName, float64(size)/1024/1024)
		} else {
			fmt.Printf("repo: %s 当前的tag数量为: %d,需要保留的tag数量是%d,无需删除! \n", repoName, len(tags), keepNum)
		}

	}

}
