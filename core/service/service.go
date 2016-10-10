package service

import (
	"errors"
	"github.com/mywebmyresponsibility.org/core/message"
)

type Service interface {
	CheckMessage(msg *message.Message) bool
	Process(msg *message.Message) error
}

type Services struct {
	on             map[string]Service
	DefaultService string
}

func (s *Services) Registry(name string, service Service) {
	s.on[name] = service
	if name == "DefaultService" {
		s.DefaultService = name
	}
}

func (s *Services) Route(msg *message.Message) error {
	for name, service := range s.on {
		if name == msg.Kind {
			if service.CheckMessage(msg) {
				return service.Process(msg)
			} else {
				errors.New("Message fail on check")
			}
		}
	}
	return s.on[s.DefaultService].Process(msg)
}

func NewServiceRouter() *Services {
	s := &Services{make(map[string]Service), ""}
	return s
}

type DefaultService struct {
}

func (s DefaultService) CheckMessage(msg *message.Message) bool {
	return false
}
func (s DefaultService) Process(msg *message.Message) error {
	return errors.New("Service DefaultService can't process")
}

func NewDefaultService() DefaultService {
	return DefaultService{}
}
