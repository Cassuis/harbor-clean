/*
@Time : 2020/2/13 15:53
@Author : Tux
@Description :
*/

package harbor

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"sort"
	"strconv"

	"clean-harbor/model"
)

// Client
type Client struct {
	Client  *http.Client //http.Clinet类型，结构体嵌套
	BaseUrl string
}

// NewClient
func NewClient(username, password, baseUrl string) *Client {
	client := &http.Client{ //定义client变量，是一个结构体
		Transport: &http.Transport{
			Proxy: func(req *http.Request) (*url.URL, error) { //匿名函数
				req.SetBasicAuth(username, password)
				return nil, nil
			},
		},
	}
	return &Client{
		Client:  client,
		BaseUrl: baseUrl,
	}
}

// getProjectID . 根据projectName 获取projectID
/* get返回内容格式如下
[
    {
        "project_id": 2,
        "owner_id": 1,
        "name": "public",
        "creation_time": "2019-09-24T15:32:08Z",
        "update_time": "2019-09-24T15:32:08Z",
        "deleted": false,
        "owner_name": "",
        "togglable": true,
        "current_user_role_id": 0,
        "repo_count": 36,
        "chart_count": 0,
        "metadata": {
            "public": "true"
        }
    }
]
*/
func (c *Client) getProjectID(projectName string) (projectId int, err error) {
	resp, err := c.Client.Get(c.BaseUrl + "/api/projects?name=" + projectName)
	if err != nil {
		return
	}
	if resp.StatusCode != 200 {
		err = fmt.Errorf("response code is:%v", resp.StatusCode)
		return
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {

		return
	}

	var projects []model.Project
	err = json.Unmarshal(body, &projects) //解析json
	if err != nil {
		return
	}
	for _, p := range projects { //循环struct,保存了name和id
		if p.Name == projectName {
			return p.ID, nil
		}
	}
	// 若前面出错则验证
	if string(body) == "null" {
		// log.Fatal("body解析异常", body)
		panic("body解析异常")

		return
	}

	return 0, errors.New("not found")
}

// func (c *Client) GetRepo(projectId int) (repoNames []string, err error)
/* 接口返回报文类似如下
[
    {
        "id": 91,
        "name": "xmc2/base-image-pytorch-test",
        "project_id": 6,
        "description": "",
        "pull_count": 0,
        "star_count": 0,
        "tags_count": 1,
        "labels": [],
        "creation_time": "2019-09-24T16:01:42.114206Z",
        "update_time": "2019-09-24T16:01:42.114206Z"
    },
    {
        "id": 100,
        "name": "xmc2/base-image-tx2-cuda",
        "project_id": 6,
        "description": "",
        "pull_count": 2,
        "star_count": 0,
        "tags_count": 2,
        "labels": [],
        "creation_time": "2019-09-24T16:04:00.92622Z",
        "update_time": "2019-10-18T08:09:24.697128Z"
	}
]
*/
func (c *Client) GetRepoNames(projectName string) (repoNames []string, err error) {
	projectId, err := c.getProjectID(projectName)
	if err != nil {
		panic("获取项目(getProjectID)ID失败")
		// return
	}
	fmt.Println(c.BaseUrl + "/api/repositories?project_id=" + strconv.Itoa(projectId))
	resp, err := c.Client.Get(c.BaseUrl + "/api/repositories?project_id=" + strconv.Itoa(projectId))
	if err != nil {
		panic("获取项目信息失败")
		// return
	}
	defer resp.Body.Close()
	// fmt.Println(string(resp.Body))
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic("body获取失败")
		// return
	}
	var repos []model.Repo
	// fmt.Println(string(body))
	// fmt.Println(body)
	err = json.Unmarshal(body, &repos) //json转结构体
	if err != nil {
		panic("json转换失败")
		// return
	}
	for _, repo := range repos { //结构体转切片
		repoNames = append(repoNames, repo.Name)
	}
	return repoNames, nil
}

//获取tag列表
/*返回报文类似如下
[
    {
        "digest": "sha256:83b764007fb2921738b6dc92a77cb205bc806335fab0f45713cd28d0d53fd013",
        "name": "0.0.2",
        "size": 1664527503,
        "architecture": "amd64",
        "os": "linux",
        "os.version": "",
        "docker_version": "18.09.2",
        "author": "",
        "created": "2019-08-17T04:52:25.276809597Z",
        "config": {
            "labels": null
        },
        "signature": null,
        "labels": []
    },
    {
        "digest": "sha256:277092a23f0e90b3e029da6f9e4bdfee53c5d8191d59124adbef942ae0163d75",
        "name": "0.0.1",
        "size": 5027236257,
        "architecture": "arm64",
        "os": "linux",
        "os.version": "",
        "docker_version": "18.06.1-ce",
        "author": "",
        "created": "2019-07-05T12:49:10.123174746Z",
        "config": {
            "labels": null
        },
        "signature": null,
        "labels": []
    }
]
*/
func (c *Client) GetRepoTags(repo string) (tags model.Tags, err error) {
	resp, err := c.Client.Get(c.BaseUrl + "/api/repositories/" + repo + "/tags")
	if err != nil {
		return
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return
	}

	err = json.Unmarshal(body, &tags) //json转结构体
	if err != nil {
		return
	}
	// 排序
	sort.Sort(tags)
	if !sort.IsSorted(tags) {
		return nil, errors.New("tags not sorted")
	}
	return
}

// DeleteRepoTag delete tags with repo name and tag.
func (c *Client) DeleteRepoTag(repo string, tag string) (err error) {
	request, err := http.NewRequest("DELETE", c.BaseUrl+"/api/repositories/"+repo+"/tags/"+tag, nil)
	if err != nil {
		return
	}
	resp, err := c.Client.Do(request)
	if err != nil {
		return
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		err = fmt.Errorf("resp code=%v", resp.StatusCode)
		return
	}
	return
}
