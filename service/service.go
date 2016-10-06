package service


type Service interface{
	CheckMessage(msg *Message) bool
	Process(msg *Message) error
}

type Services struct{
	map[string]*Service
}
func (s *Services)Registry(name string, service *Service){

}
/*
post [m1,m2,m3]
messages := decode(post)
foreach message {
	if services.Kind(message.Kind).Check(message) {
		return services.Kind(message.Kind).Process(message)
	}
}

--- 
passar tudo isso para channel, colocar um select e ser feliz.
*/

