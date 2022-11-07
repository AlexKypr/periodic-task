package handlers

import (
	"fmt"
	"net/http"
	"periodictask/ptask"
	"periodictask/utils"

	"go.uber.org/zap"
)

type TaskHandler struct {
	l *zap.Logger
}

func NewTaskHandler(l *zap.Logger) *TaskHandler {
	return &TaskHandler{
		l: l,
	}
}

func (t *TaskHandler) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	handleMsg := fmt.Sprintf("Handle %s %s", r.Method, r.URL.String())
	t.l.Info(handleMsg)
	defer func() {
		t.l.Sugar().Infof("Exit %s", handleMsg)
	}()

	if r.Method == http.MethodGet {
		t.getTsArray(rw, r)
		return
	}

	// catch http methods that are not supported
	t.l.Warn("Method not supported")
	errDto := utils.NewErrorDTO("error", "Method not supported")
	utils.ToJSON(rw, errDto, http.StatusBadRequest)
}

func (t *TaskHandler) getTsArray(rw http.ResponseWriter, r *http.Request) {
	t.l.Info("Entering getTsArray")
	defer func() {
		t.l.Info("Exiting getTsArray")
	}()

	ptask, err := ptask.GetPTaskFromURLQueries(r, t.l)
	if err != nil {
		t.l.Error(err.Error())
		errDto := utils.NewErrorDTO("error", err.Error())
		utils.ToJSON(rw, errDto, http.StatusBadRequest)
		t.l.Sugar().Infof("Return http code: %d", http.StatusBadRequest)
		return
	}

	tsarray, err := ptask.GetPTasks()
	if err != nil {
		t.l.Error(err.Error())
		errDto := utils.NewErrorDTO("error", err.Error())
		utils.ToJSON(rw, errDto, http.StatusBadRequest)
		t.l.Sugar().Infof("Return http code: %d", http.StatusBadRequest)
		return
	}

	utils.ToJSON(rw, tsarray, http.StatusOK)
	t.l.Sugar().Infof("Return http code: %d", http.StatusOK)
}
