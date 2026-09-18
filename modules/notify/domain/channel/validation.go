package channel

import (
	"strings"

	notifyerr "nfxnews/errors/src/notify"
)

var supportedKinds = map[string]struct{}{
	"feishu": {}, "dingtalk": {}, "wework": {}, "telegram": {},
	"email": {}, "ntfy": {}, "bark": {}, "slack": {},
}

func Kinds() []string {
	return []string{"feishu", "dingtalk", "wework", "telegram", "email", "ntfy", "bark", "slack"}
}

func NormalizeKind(kind string) (string, error) {
	kind = strings.ToLower(strings.TrimSpace(kind))
	if _, ok := supportedKinds[kind]; !ok {
		return "", notifyerr.ErrNotifyKindInvalid
	}
	return kind, nil
}
