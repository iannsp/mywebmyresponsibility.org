package service

import (
	"github.com/mywebmyresponsibility.org/core/message"
	"github.com/mywebmyresponsibility.org/core/service"
	"testing"
)

func TestServicesWithDefaultServiceOnly(t *testing.T) {
	msg := message.NewMessage("INVALID")
	services := service.NewServiceRouter()
	services.Registry("DefaultService", service.NewDefaultService())
	err := services.Route(msg)
	if err == nil || err.Error() != "Service DefaultService can't process" {
		t.Error(err)
	}
}
