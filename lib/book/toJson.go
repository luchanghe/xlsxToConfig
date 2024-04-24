package book

import (
	"encoding/json"
	"fmt"
	"github.com/iancoleman/orderedmap"
)

func (b *Book) WriteToJson() {
	b.write(b.toJson(), "json")
}

func (b *Book) toJson() []byte {
	o := orderedmap.New()
	o.Set("a", 1.01)
	jsonData, err := json.MarshalIndent(b.dataMap, "", "    ")
	if err != nil {
		fmt.Println("Error marshalling JSON:", err)
	}
	return jsonData
}
