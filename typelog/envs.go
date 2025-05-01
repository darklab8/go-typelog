package typelog

import (
	"os"
	"strings"
)

const (
	TOOL_NAME = "typelog"
)

type TypelogEnvs struct {
	EnableJson        bool
	EnableFileShown   bool
	EnableScopesShown bool
}

var EnvPresent struct {
	EnableJson        bool
	EnableFileShown   bool
	EnableScopesShown bool
} = struct {
	EnableJson        bool
	EnableFileShown   bool
	EnableScopesShown bool
}{
	EnableJson:        os.Getenv(strings.ToUpper(TOOL_NAME)+"_LOG_JSON") != "" || os.Getenv(strings.ToUpper(TOOL_NAME)+"_JSON") != "",
	EnableFileShown:   os.Getenv(strings.ToUpper(TOOL_NAME)+"_LOG_FILE_SHOWING") != "" || os.Getenv(strings.ToUpper(TOOL_NAME)+"_FILES") != "",
	EnableScopesShown: os.Getenv(strings.ToUpper(TOOL_NAME)+"_SCOPES") != "",
}

var Env TypelogEnvs = TypelogEnvs{
	EnableJson:        os.Getenv(strings.ToUpper(TOOL_NAME)+"_LOG_JSON") == "true" || os.Getenv(strings.ToUpper(TOOL_NAME)+"_JSON") == "true",
	EnableFileShown:   os.Getenv(strings.ToUpper(TOOL_NAME)+"_LOG_FILE_SHOWING") == "true" || os.Getenv(strings.ToUpper(TOOL_NAME)+"_FILES") == "true",
	EnableScopesShown: os.Getenv(strings.ToUpper(TOOL_NAME)+"_SCOPES") == "true",
}
