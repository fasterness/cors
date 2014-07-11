cors
====

CORS middleware for Go


```
  corsHandler := cors.New(myHandler)
  corsHandler.RemoveMethod("PUT")
  corsHandler.RemoveMethod("DELETE")
  max_age := 60 * time.Second
  corsHandler.SetMaxAge(max_age.Seconds())

  http.Handle("/", corsHandler)

  http.ListenAndServe("localhost:9999", nil)

```
