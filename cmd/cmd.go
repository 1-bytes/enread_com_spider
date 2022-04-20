package cmd

import (
	"context"
	"encoding/json"
	elasticsearch "enread_com/pkg/elastic"
	"fmt"
	"github.com/olivere/elastic/v7"
)

// SaveData 存储数据
func SaveData(index string, id string, data interface{}) error {
	j, err := json.Marshal(data)
	if err != nil {
		return err
	}

	var e *elastic.IndexService
	e = elasticsearch.GetInstance().Index()
	if id != "" {
		e.Id(id)
	}
	do, err := e.Index(index).BodyJson(string(j)).Do(context.Background())
	fmt.Printf("%+v: %+v\n", do.Result, do.Id)
	return err
}
