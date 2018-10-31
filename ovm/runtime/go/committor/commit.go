package committor

import (
  "encoding/json"
  "encoding/base64"
  "log"
)

func (c Committor) Commit (returnValue string, transaction string) {
  preTransactionMap := map[string] string {"returnValue" : returnValue, "transaction" : transaction}
  blob, error := json.Marshal(preTransactionMap)
  if error == nil {
    transactionEncoded := base64.StdEncoding.EncodeToString(blob)
    log.Println("send transaction:",transactionEncoded)
  }else{
    log.Fatal(error)
  }
}
