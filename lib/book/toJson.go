package book

import (
	"encoding/json"
	"fmt"
)

func (b *Book) WriteToJson() {
	b.write(b.toJson(), "json")
}

func (b *Book) toJson() []byte {
	jsonData, err := json.MarshalIndent(b.dataMap, "", "    ")
	if err != nil {
		fmt.Println("Error marshalling JSON:", err)
	}
	return jsonData
}
