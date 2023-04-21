package main

// -----------------------------------------
import (
	"gogin/tools"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	docs "gogin/docs"

	swaggerfiles "github.com/swaggo/files"     // swagger embed files
	ginSwagger "github.com/swaggo/gin-swagger" // gin-swagger middleware
)

// Some Globals-----------------------------
var persons = tools.NewPersons()

// Playground-------------------------------

type msg struct {
	Message string `json:"message"`
}
type id struct {
	Id uuid.UUID `json:"id"`
}

// -----------------------------------------
// @Schemes
// @Description Hello From the API
// @Tags root
// @Accept json
// @Produce json
// // @Param   some_id     path    string     true        "Some ID"
// @Success 200 {object} msg
// @Router / [get]
func HelloFromTheApi(ctx *gin.Context) {
	nowStr := time.Now().Format("2006-01-02 15:04:05")
	myMsg := msg{Message: "Hello from the API [" + nowStr + "]"}
	ctx.JSON(http.StatusOK, myMsg)
}

// -------------------------
// @Schemes
// @Description Get list of person items
// @Tags persons
// @Accept json
// @Produce json
// @Success 200 {object} []tools.Person
// @Router /persons [get]
func PersonsGet(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, persons.Data)
}

// -------------------------
// @Schemes
// @Description Get person item by id
// @Tags persons
// @Accept json
// @Produce json
// @Param  id  path  string  true  "Person ID"
// @Success 200 {object} []tools.Person
// @Router /person/{id} [get]
func PersonGet(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}
	person, err := persons.Get(id)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"message": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, person)
}

// -------------------------
// @Schemes
// @Description Create new person item
// @Tags persons
// @Accept json
// @Produce json
// @Param  data  body  tools.Person  true  "Person struct"
// @Success 200 {object} id
// @Router /person [post]
func PersonPost(ctx *gin.Context) {
	person := new(tools.Person)
	if err := ctx.BindJSON(&person); err != nil {
		ctx.AbortWithError(http.StatusBadRequest, err)
		return
	}
	id := persons.Add(*person)

	// fmt.Println(id)
	ctx.JSON(http.StatusAccepted, gin.H{
		"id": id,
	})
}

// -------------------------
// @Schemes
// @Description Edit person item
// @Tags persons
// @Accept json
// @Produce json
// @Param  id    path    string        true  "Person ID"
// @Param  data  body    tools.Person  true  "Person struct"
// @Success 200 {object} tools.Person
// @Router /person/{id} [put]
func PersonPut(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}
	person := new(tools.Person)
	if err := ctx.BindJSON(&person); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}
	if _, err := persons.Change(id, *person); err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"message": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, person)
}

// -------------------------
// @Schemes
// @Description Edit person item
// @Tags persons
// @Accept json
// @Produce json
// @Param  id    path    string  true  "Person ID"
// @Success 200 {object} id
// @Router /person/{id} [delete]
func PersonDelete(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	if _, err := persons.Delete(id); err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"message": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"id": id})
}

// -------------------------
// -------------------------
// -------------------------

// ------------------------------------------------
func main() {
	router := gin.Default()
	gin.SetMode(gin.DebugMode)
	docs.SwaggerInfo.BasePath = "/api"

	// The Routes-------------------------------
	router.GET("/api", HelloFromTheApi)

	router.GET("/api/persons", PersonsGet)

	router.GET("/api/person/:id", PersonGet)

	router.POST("/api/person", PersonPost)

	router.PUT("/api/person/:id", PersonPut)

	router.DELETE("/api/person/:id", PersonDelete)

	// The Runner -----------------------------
	router.GET("/docs/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))
	router.Run()
}
