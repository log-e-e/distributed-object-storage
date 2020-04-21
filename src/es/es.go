package es

import (
    "context"
    "encoding/json"
    "fmt"
    "github.com/olivere/elastic/v7"
    "log"
    "os"
    "reflect"
)

type Metadata struct {
    Name string `json:"name"`
    Version int `json:"version"`
    Size int64  `json:"size"`
    Hash string `json:"hash"`
}

var (
    serverEsClientMap = make(map[string]*elastic.Client)
)

func getEsClient(esServer string) *elastic.Client {
    if serverEsClientMap[esServer] != nil {
        return serverEsClientMap[esServer]
    }
    newClient, err := elastic.NewClient()
    if err != nil {
        log.Printf("ES Error: failed to new es client for server [%s]\n, details: %s\n", esServer, err.Error())
        return nil
    }
    serverEsClientMap[esServer] = newClient
    return newClient
}

func metadataExists(name string, version int, size int64, hash string) bool {
    esClient := getEsClient(os.Getenv("ES_SERVER"))
    nameQuery := elastic.NewTermQuery("name", name)
    versionQuery := elastic.NewTermQuery("version", version)
    sizeQuery := elastic.NewTermQuery("size", size)
    hashQuery := elastic.NewTermQuery("hash", hash)
    searchResult, err := esClient.Search().
        Index("metadata").
        Query(nameQuery).Query(versionQuery).Query(sizeQuery).Query(hashQuery).
        Pretty(true).
        Do(context.Background())
    if err != nil {
        log.Println(err)
        return false
    }
    return searchResult.TotalHits() != 0
}

func getMetadata(name string, version int) (metadata Metadata, err error) {
    esClient := getEsClient(os.Getenv("ES_SERVER"))
    nameQuery, versionQuery := elastic.NewTermQuery("name", name), elastic.NewTermQuery("version", version)
    searchResult, err := esClient.Search().
        Index("metadata").
        Query(nameQuery).Query(versionQuery).
        Pretty(true).
        Do(context.Background())
    if err != nil {
        return
    }

    dataBytes, err := searchResult.Hits.Hits[0].Source.MarshalJSON()
    if err != nil {
        return
    }
    json.Unmarshal(dataBytes, &metadata)
    return
}

func SearchLatestVersion(name string) (metadata Metadata, err error) {
    esClient := getEsClient(os.Getenv("ES_SERVER"))
    searchResult, err := esClient.Search().
        Index("metadata").
        Query(elastic.NewTermQuery("name", name)).
        Sort("version", false).
        Pretty(true).
        Do(context.Background())
    if err != nil {
        return
    }
    if searchResult.TotalHits() > 0 {
        dataBytes, err := searchResult.Hits.Hits[0].Source.MarshalJSON()
        if err != nil {
           return metadata, err
        }
        json.Unmarshal(dataBytes, &metadata)
    }
    return
}

func GetMetadata(name string, version int) (Metadata, error) {
    // version一般至少从1开始，若为0则说明未指定特定版本，默认返回最新版本
    if version == 0 {
        return SearchLatestVersion(name)
    }
    return getMetadata(name, version)
}

func PutMetadata(name string, version int, size int64, hash string) error {
    if metadataExists(name, version, size, hash) {
        return PutMetadata(name, version + 1, size, hash)
    }

    metadata := Metadata{
        Name:    name,
        Version: version,
        Size:    size,
        Hash:    hash,
    }
    esClient := getEsClient(os.Getenv("ES_SERVER"))
    _, err := esClient.Index().
        Index("metadata").
        Id(fmt.Sprintf("%s_%d", name, version)).
        BodyJson(metadata).
        Refresh("wait_for").
        Do(context.Background())
    return err
}

func AddVersion(name string, size int64, hash string) error {
    metadata, err := SearchLatestVersion(name)
    if err != nil {
        return err
    }
    return PutMetadata(name, metadata.Version + 1, size, hash)
}

func SearchAllVersions(name string, from, size int) ([]Metadata, error) {
    esClient := getEsClient(os.Getenv("ES_SERVER"))
    var searchResult *elastic.SearchResult
    var err error
    searchService := esClient.Search().
        Index("metadata").
        From(from).Size(size).
        Pretty(true)
    if name == "" {
        searchResult, err = searchService.Do(context.Background())
    } else {
        nameQuery := elastic.NewTermQuery("name", name)
        searchResult, err = searchService.
            Query(nameQuery).
            Do(context.Background())
    }
    if err != nil {
        return nil, err
    }

    var metadata Metadata
    var metadatas []Metadata
    for _, item := range searchResult.Each(reflect.TypeOf(metadata)) {
        if t, ok := item.(Metadata); ok {
            metadatas = append(metadatas, Metadata{
                Name:    t.Name,
                Version: t.Version,
                Size:    t.Size,
                Hash:    t.Hash,
            })
        }
    }
    return metadatas, nil
}

// 暂时没用上，本章的删除机制是创建一个新版本并将size和hash置空
func DelMetadata(name string, version int) {
   esClient := getEsClient(os.Getenv("ES_SERVER"))
   esClient.Update().
       Index("metadata").
       Id(fmt.Sprintf("%s_%d", name, version)).
       Doc(map[string]interface{}{
           "size": 0,
           "hash": "",
   })
}
