package indexing
import (
  "net/http"
  "encoding/json"
  "gopkg.in/mgo.v2"
  "gopkg.in/mgo.v2/bson"
  "fmt"
  "strings"
  "math"
)

type Object struct{
  Id string `bson:"Id"`
  Text []byte `bson:"Text"`
}


func AddIndex(w http.ResponseWriter, r *http.Request){
    if r.Method == http.MethodGet {
        fmt.Fprint(w, "Please send POST request with params for create index")
    }else if r.Method == http.MethodPost {

      var id string = r.FormValue("id")
      var name_table string = r.FormValue("table-name")
      var text string = r.FormValue("text")

      error_counter := 0
      if len(id) == 0{
        error_counter += 1
        fmt.Fprint(w, "Incorrect 'id' value, please add 'id' field\n")
      }
      if len(name_table) == 0{
        error_counter += 1
        fmt.Fprint(w, "Incorrect 'table-name' value, please add 'table-name' field\n")
      }
      if len(text) == 0{
        error_counter += 1
        fmt.Fprint(w, "Incorrect 'text' value, please add 'text' field\n")
      }

      session, err := mgo.Dial("mongodb://127.0.0.1")

      if err == nil && error_counter == 0{
          collection := session.DB("DataBase").C(name_table)

          var data interface{}
          err := json.Unmarshal([]byte(text), &data)
          data1 := data.(map[string]interface{})

          if err != nil{
            fmt.Println(err)
          }else{
            new_data := prepare_for_indexing(data1)
            err1 := collection.Update(bson.M{"Id": id}, bson.M{"$set":bson.M{"Text": new_data}})

            if err1 != nil{
              obj1 := &Object{Id: id, Text: new_data}
              err := collection.Insert(obj1)

              if err != nil{
                fmt.Fprint(w, err)
              }else{
                fmt.Fprint(w, "{'result': 'success'}")
              }
            }else{
              fmt.Fprint(w, "{'result': 'success'}")
            }
          }
      }
      defer session.Close()
    }else {
        fmt.Fprint(w, "Method " + r.Method + " is not suported")
    }
}

func RemoveIndex(w http.ResponseWriter, r *http.Request){
  id := r.FormValue("id")
  name_table := r.FormValue("table-name")

  session, err := mgo.Dial("mongodb://127.0.0.1")
  if err == nil{
    collection := session.DB("DataBase").C(name_table)
    _, err = collection.RemoveAll(bson.M{"Id": id})
    if err != nil{
        fmt.Fprint(w, "Invalid table-name, please add new table-name")
    }else{
      fmt.Fprint(w, "{'result': 'success'}")
    }
  }
  defer session.Close()

}
func string_to_tanh(title string) []float64{
  title_byte := []byte(title)
  write_array := make([]float64, len(title_byte))

  counter := 0
  for i := range title_byte{
    write_array[counter] = norm(float64(title_byte[i]))
    counter += 1
  }

  return write_array
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
func prepare_for_indexing(data map[string]interface {}) []byte{
  return_data := make(map[string]interface{})

  for key := range data{
    if strings.Contains(key, "__indexing"){
      return_data[key[:len(key)-10]] = string_to_tanh(data[key].(string))
    }else{
      return_data[key] = data[key]
    }
  }

  content, _ := json.Marshal(return_data)

  return content
}
