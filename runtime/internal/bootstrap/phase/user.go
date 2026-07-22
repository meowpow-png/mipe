package phase

import (
	"fmt"
	"os"
	"strconv"
)

var currentIdentity = func() (int, int) {
	return os.Geteuid(), os.Getegid()
}

func commandAsUser(uid, gid, name string, args ...string) (string, []string, error) {
	targetUID, err := strconv.Atoi(uid)
	if err != nil {
		return "", nil, fmt.Errorf("local_uid must be a numeric user id: %w", err)
	}
	targetGID, err := strconv.Atoi(gid)
	if err != nil {
		return "", nil, fmt.Errorf("local_gid must be a numeric group id: %w", err)
	}
	currentUID, currentGID := currentIdentity()
	if currentUID == targetUID && currentGID == targetGID {
		return name, args, nil
	}
	return "gosu", append([]string{uid + ":" + gid, name}, args...), nil
}
