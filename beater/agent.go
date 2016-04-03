package beater

import (
	"os"
	"net/http"
	"io/ioutil"
	"encoding/json"
	"path/filepath"
	"github.com/elastic/beats/libbeat/logp"
)

func (mb *Mesosbeat) GetAgentStatistics(u string) (map[string]float64, error) {
	statistics := make(map[string]float64)

	resp, err := http.Get(u)
	defer resp.Body.Close()

	if err != nil {
		logp.Err("An error occured while executing HTTP request: %v", err)
		return statistics, err
	}

	// read json http response
    jsonDataFromHttp, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		logp.Err("An error occured while reading HTTP response: %v", err)
		return statistics, err
	}

    err = json.Unmarshal([]byte(jsonDataFromHttp), &statistics)

	if err != nil {
		logp.Err("An error occured while unmarshaling agent statistics: %v", err)
		return statistics, err
	}
	return statistics, nil
}

func (mb *Mesosbeat) GetAgentAttributes(dir string) (map[string]string, error) {
	// need to find files rooted at dir/attributes/ and create a map
	// filename -> contents

	attributes := make(map[string]string)
    fileList := []string{}

    err := filepath.Walk(dir + "/attributes", func(path string, f os.FileInfo, err error) error {
    	if f.IsDir() {
    		// continue
    	} else {
    		fileList = append(fileList, f.Name())
    	}

        return err
    })

	if err != nil {
		logp.Err("An error occured while executing attributes search: %v", err)
		return attributes, err
	}

    for _, file := range fileList {
        buf,err := ioutil.ReadFile(dir + "/attributes/" + file)

     	if err != nil {
     		logp.Err("An error occurred reading file: %v", err)
     		return attributes, err
     	}

     	content := string(buf)

     	attributes[file] = content
    }

	return attributes, nil
}