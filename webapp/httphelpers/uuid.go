package httphelpers

import "github.com/twinj/uuid"

func GenerateUUID() string {
	return uuid.NewV4().String()
}
