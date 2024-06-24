package main

import (
	"bytes"
	"encoding/csv"
	"encoding/gob"
	"encoding/json"
	"io"
	"log"
	"net/http"
)

type Data struct {
	Part [][]string
}

type Cluster struct {
	Centroid []float64
	Points   [][]float64
}

func LoadAndDivideDataset() [][][]string {
	url := "https://raw.githubusercontent.com/MrPepePollo/TF_Concurrente/master/SocialNetworkDataset.csv"

	resp, err := http.Get(url)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Fatalf("unexpected status code: %d", resp.StatusCode)
	}

	reader := csv.NewReader(resp.Body)

	var records [][]string

	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatal(err)
		}

		records = append(records, record)
	}

	numGoroutines := 10
	partSize := len(records) / numGoroutines
	parts := make([][][]string, numGoroutines)
	for i := 0; i < numGoroutines; i++ {
		start := i * partSize
		end := start + partSize
		if i == numGoroutines-1 {
			end = len(records)
		}
		parts[i] = records[start:end]
	}

	return parts
}

func handleStart(w http.ResponseWriter, r *http.Request) {
	log.Println("Received a request")
	parts := LoadAndDivideDataset()

	var totalClusters []Cluster
	for _, part := range parts {
		data := &Data{Part: part}
		buf := &bytes.Buffer{}
		gob.NewEncoder(buf).Encode(data)

		resp, err := http.Post("http://localhost:8080/api/kmeans", "application/gob", buf)
		if err != nil {
			log.Fatal(err)
		}
		defer resp.Body.Close()

		clusters := make([]Cluster, 0)
		json.NewDecoder(resp.Body).Decode(&clusters)
		totalClusters = append(totalClusters, clusters...)
	}

	for i, cluster := range totalClusters {
		log.Printf("Cluster %d: %v\n", i, cluster.Centroid)
	}

	w.Write([]byte("Process started"))
}

func main() {
	http.HandleFunc("/start", handleStart)
	log.Fatal(http.ListenAndServe(":8081", nil))
}
