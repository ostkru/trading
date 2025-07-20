package main

import (
fmt
log
net/http
github.com/gin-gonic/gin
)

func main() {
r := gin.Default()
r.GET(/, func(c *gin.Context) {
c.JSON(http.StatusOK, gin.H{message: API on 8095, status: ok})
})
fmt.Println(Starting on 8095...)
if err := r.Run(:8095); err != nil {
log.Fatal(err)
}
}
