package uid

import(
	"github.com/google/uuid"
	"strings"
)

func NewUid() string {
	uid := uuid.New()
	uidstr := strings.Split(uid.String(), "-")[0]
	return uidstr
}