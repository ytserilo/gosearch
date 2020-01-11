package search
import (
  "net/http"
  "encoding/json"
  "fmt"
  "math"
  "gopkg.in/mgo.v2"
  "gopkg.in/mgo.v2/bson"
  "sort"
  "strings"
)

type SearchObject struct{
  Id string
  Score float64
}

type Object struct{
  Id string `bson:"Id"`
  Text []byte `bson:"Text"`
}


func Search(w http.ResponseWriter, r *http.Request){
    var text string = r.FormValue("text")
    var name_table string = r.FormValue("table-name")

    var data interface{}
    err := json.Unmarshal([]byte(text), &data)
    if err != nil{
      fmt.Println(err)
    }
    search_data := data.(map[string]interface{})

    session, err := mgo.Dial("mongodb://127.0.0.1")

    if err == nil{
        collection := session.DB("DataBase").C(name_table)

        query := bson.M{}
        products := []Object{}
        collection.Find(query).All(&products)

        objects := first_filter(search_data, products)

        var sort_array []SearchObject

        for search_param := range search_data{
          if strings.Contains(string(search_param), "__search"){
            var maxLength int = find_max(objects, string(search_param[:len(search_param)-8]))

            var input_text []float64

            var textLength int = len(search_data[search_param].(string))

            if maxLength < textLength{
              input_text = string_to_tanh(strings.ToLower(search_data[search_param].(string)), textLength)
            }else{
              input_text = string_to_tanh(strings.ToLower(search_data[search_param].(string)), maxLength)
            }

            var objectsLength int = len(objects)

            fin := make(chan []SearchObject, 10000)
            var return_result []SearchObject


            for i := 0; i < objectsLength; i += 20{
              if i+20 > objectsLength{
                go search_objects(input_text, products[i: objectsLength], string(search_param[:len(search_param)-8]), fin)
              }else{
                go search_objects(input_text, products[i: i+20], string(search_param[:len(search_param)-8]), fin)
              }
            }
            for i := 0; i < objectsLength; i += 20{
              return_result = append(return_result, <-fin...)
            }

            if len(sort_array) == 0{
              sort_array = return_result
            }else{
              for i := range return_result{
                for j := range sort_array{
                  if return_result[i].Id == sort_array[j].Id{
                    sort_array[j].Score = sort_array[j].Score + return_result[i].Score
                    break
                  }else if len(sort_array)-1 == j{
                    sort_array = append(sort_array, return_result[i])
                  }
                }
              }
            }
          }
        }
        sort.Slice(sort_array, func(i, j int) bool { return sort_array[i].Score < sort_array[j].Score })
        json_result, _ := json.Marshal(sort_array)
        fmt.Fprint(w, string(json_result))
    }
    defer session.Close()
}



func search_objects(input_title []float64, objects []Object, key string, fin chan []SearchObject){
  var result []SearchObject
  var lengthInputTitle int = len(input_title)

  for i := range objects{
    var standart []float64 = make([]float64, lengthInputTitle)
    var f interface{}
    json.Unmarshal(objects[i].Text, &f)
    data := f.(map[string]interface{})
    var get_title []float64

    for i := range data[key].([]interface{}){
      get_title = append(get_title, data[key].([]interface{})[i].(float64))
    }

    for char := range get_title{
      standart[char] = get_title[char]
    }

    var avg float64 = 0

    for j := range input_title{
      avg += math.Atan((standart[j] - input_title[j]) * (standart[j] - input_title[j]))
    }

    avg = avg / float64(lengthInputTitle)
    if avg < 0.01{
      result = append(result, SearchObject{Id: objects[i].Id, Score: avg})
    }
  }
  fin <- result
}

func norm(val float64) float64{
  var minval float64 = 0.0
  var maxval float64 = 150.0
  var newmin float64 = -1.0
  var newmax float64 = 1.0

  result := newmin + (val - minval) * (newmax - newmin) / (maxval - minval)

  result = math.Tanh(result)
  return result
}


func string_to_tanh(title string, max int) []float64{
  var title_byte []byte = []byte(title)
  var write_array []float64 = make([]float64, max)

  for i := range title_byte{
    write_array[i] = norm(float64(title_byte[i]))
  }
  return write_array
}

func find_max(objects []Object, key string) int{
  var max int = 0

  for i := range objects{
    var f interface{}
    json.Unmarshal([]byte(objects[i].Text), &f)
    data := f.(map[string]interface{})

    text := data[key].([]interface{})

    var maxLength int = len(text)

    if maxLength > max{
      max = maxLength
    }
  }
  return max
}

func from_string_to_tanh(str []interface{}) []float64{
  var new_byte []float64

  for i := range str{
    new_byte = append(new_byte, str[i].(float64))
  }
  return new_byte
}
func first_filter(search_data map[string]interface{}, objects []Object) []Object{
  var return_array []Object

  for key := range search_data{

    for obj := range objects{
        var data interface{}
        json.Unmarshal(objects[obj].Text, &data)
        ndata := data.(map[string]interface{})

        if strings.Contains(key, "__lte"){
          data_key := key[:len(key)-5]
          if ndata[data_key].(float64) < search_data[key].(float64){
            return_array = append(return_array, objects[obj])
          }
        }else if strings.Contains(key, "__gte"){
          data_key := key[:len(key)-5]

          if ndata[data_key].(float64) > search_data[key].(float64){
            return_array = append(return_array, objects[obj])
          }
        }else if strings.Contains(key, "__lt"){
          data_key := key[:len(key)-4]
          if ndata[data_key].(float64) <= search_data[key].(float64){
            return_array = append(return_array, objects[obj])
          }
        }else if strings.Contains(key, "__gt"){
          data_key := key[:len(key)-4]
          if ndata[data_key].(float64) >= search_data[key].(float64){
            return_array = append(return_array, objects[obj])
          }
        }else if strings.Contains(key, "__in"){
          data_key := key[:len(key)-4]
          if strings.Contains(ndata[data_key].(string), search_data[key].(string)){
            return_array = append(return_array, objects[obj])
          }
        }else{
          return_array = append(return_array, objects[obj])
        }
    }
  }
  return return_array
}
