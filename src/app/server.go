// server.go
package main

import (
	"encoding/gob"
	"encoding/json"
	"log"
	"math/rand"
	"net"
	"net/http"
	"strconv"
)

type Data struct {
	Part [][]string
}

type Cluster struct {
	Centroid []float64
	Sum      []float64
	Count    int
}

func handleRequest(conn net.Conn) {
	decoder := gob.NewDecoder(conn)
	data := &Data{}
	decoder.Decode(data)

	clusters := kmeans(data.Part, 3)

	encoder := gob.NewEncoder(conn)
	encoder.Encode(clusters)
}

func handleKmeansRequest(w http.ResponseWriter, r *http.Request) {
	decoder := gob.NewDecoder(r.Body)
	data := &Data{}
	decoder.Decode(data)

	clusters := kmeans(data.Part, 3)

	json.NewEncoder(w).Encode(clusters)
}

func main() {
	http.HandleFunc("/api/kmeans", handleKmeansRequest)
	http.ListenAndServe(":8080", nil)
}

func kmeans(part [][]string, k int) []Cluster {
	if len(part) == 0 {
		log.Fatalf("The part slice is empty.")
	}
	clusters := make([]Cluster, k)
	for i := range clusters {
		clusters[i].Centroid = make([]float64, len(part[0])-1)
		clusters[i].Sum = make([]float64, len(part[0])-1)
		for j := range clusters[i].Centroid {
			switch j {
			case 1: // Age
				clusters[i].Centroid[j] = rand.Float64() * 100
			case 3, 4, 5, 6, 7, 8:
				clusters[i].Centroid[j] = rand.Float64() * 10000
			default:
				clusters[i].Centroid[j] = rand.Float64()
			}
		}
	}

	for i := 0; i < 1000; i++ {
		for _, row := range part {
			point := make([]float64, len(row)-1)
			for j := range point {
				point[j], _ = strconv.ParseFloat(row[j+1], 64)
			}

			minDistance := distance(point, clusters[0].Centroid)
			minIndex := 0
			for j := 1; j < k; j++ {
				d := distance(point, clusters[j].Centroid)
				if d < minDistance {
					minDistance = d
					minIndex = j
				}
			}

			clusters[minIndex].Count++
			for j, x := range point {
				clusters[minIndex].Sum[j] += x
			}
		}

		for i := range clusters {
			if clusters[i].Count > 0 {
				for j := range clusters[i].Centroid {
					clusters[i].Centroid[j] = clusters[i].Sum[j] / float64(clusters[i].Count)
				}
			} else {
				for j := range clusters[i].Centroid {
					switch j {
					case 1: // Age
						clusters[i].Centroid[j] = rand.Float64() * 100
					case 3, 4, 5, 6, 7, 8:
						clusters[i].Centroid[j] = rand.Float64() * 10000
					default:
						clusters[i].Centroid[j] = rand.Float64()
					}
				}
			}
		}
	}

	return clusters
}

func distance(a, b []float64) float64 {
	sum := 0.0
	for i := range a {
		d := a[i] - b[i]
		sum += d * d
	}
	return sum
}
