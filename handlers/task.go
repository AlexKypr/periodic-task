package handlers

import (
	"inaccess/ptask"
	"inaccess/utils"
	"log"
	"net/http"
)

type TaskHandler struct {
	l *log.Logger
}

func NewTaskHandler(l *log.Logger) *TaskHandler {
	return &TaskHandler{
		l: l,
	}
}

func (t *TaskHandler) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		t.getTsArray(rw, r)
		return
	}

	// catch http methods that are not supported
	rw.WriteHeader(http.StatusBadRequest)
}

func (t *TaskHandler) getTsArray(rw http.ResponseWriter, r *http.Request) {
	t.l.Println("Getting timestamp array")

	ptask, err := ptask.GetPTaskFromURLQueries(r, t.l)
	if err != nil {
		t.l.Println(err.Error())
		errDto := utils.NewErrorDTO("error", err.Error())
		utils.ToJSON(rw, errDto, http.StatusNotFound)
		return
	}

	tsarray, err := ptask.GetPTasks()
	if err != nil {
		t.l.Println(err.Error())
		errDto := utils.NewErrorDTO("error", err.Error())
		utils.ToJSON(rw, errDto, http.StatusNotFound)
		return
	}

	utils.ToJSON(rw, tsarray, http.StatusOK)
}
