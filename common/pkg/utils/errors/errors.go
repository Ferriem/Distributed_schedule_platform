package errors

import "errors"

var (
	ErrClientNotFound   = errors.New("mysql client not found")
	ErrClientDbNameNull = errors.New("mysql dbname is null")
	ErrValueMayChanged  = errors.New("the value has been changed by others on this time")
	ErrEtcdNotInit      = errors.New("etcd is not initialized")

	ErrNotFound = errors.New("not found")

	ErrEmptyJobName        = errors.New("name of job is empty")
	ErrEmptyJobCommand     = errors.New("command of job is empty")
	ErrIllegalJobId        = errors.New("invalid id that includes illegal characters such as '/' '\\'")
	ErrIllegalJobGroupName = errors.New("invalid job group name that includes illegal characters such as '/' '\\'")

	ErrEmptyScriptName    = errors.New("name of script is empty")
	ErrEmptyScriptCommand = errors.New("command of script is empty")
	ErrEmptyNodeGroupName = errors.New("name of node group is empty")
	ErrIllegalNodeGroupId = errors.New("invalid node group id that includes illegal characters such as '/'")
)
