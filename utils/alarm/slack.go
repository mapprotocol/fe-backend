package alarm

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/mapprotocol/fe-backend/resource/log"
	"io"
	blog "log"
	"net/http"
	"os"
	"time"
)

const (
	Local   = "local"
	Project = "fe-task"
)

var (
	env = ""
)

const (
	Env   = "FE_TASK_ENV"
	Hooks = "FE_TASK_SLACK_HOOKS"
)

const defaultHTTPTimeout = 20 * time.Second

func ValidateEnv() {
	env = os.Getenv(Env)
	if env == "" {
		blog.Fatal("FE_TASK_ENV is empty")
	}
	hooks := os.Getenv(Hooks)
	if hooks == "" {
		blog.Fatal("FE_TASK_SLACK_HOOKS is empty")
	}
	blog.Printf("env: %s", env)
}

func Slack(ctx context.Context, msg string) {
	hooks := os.Getenv(Hooks)
	if hooks == "" {
		log.Logger().Error("hooks is empty")
		return
	}
	if env == Local {
		log.Logger().Debug("evn is local does not send msg to slack", "msg", msg)
		return
	}

	body, err := json.Marshal(map[string]interface{}{
		"text": fmt.Sprintf("Project: %s\nEnv: %s\nMsg: %s", Project, env, msg),
	})
	if err != nil {
		params := map[string]interface{}{
			"env":   Env,
			"msg":   msg,
			"error": err.Error(),
		}
		log.Logger().WithFields(params).Error("failed to json marshal")
		return
	}

	client := http.Client{
		Timeout: defaultHTTPTimeout,
	}
	req, err := http.NewRequestWithContext(ctx, "POST", hooks, io.NopCloser(bytes.NewReader(body)))
	if err != nil {
		params := map[string]interface{}{
			"hooks": hooks,
			"body":  string(body),
			"error": err.Error(),
		}
		log.Logger().WithFields(params).Error("failed to new request with context")
		return
	}
	req.Header.Set("Content-type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		log.Logger().WithField("error", err.Error()).Error("failed to send request")
		return
	}
	if resp == nil {
		log.Logger().Error("response is nil")
		return
	}
	if resp.Body != nil {
		defer resp.Body.Close()
	}

	if resp.StatusCode != http.StatusOK {
		log.Logger().Error("response status code is not 200", "code", resp.StatusCode)
		return
	}
	log.Logger().Info("sent alarm message")
}
