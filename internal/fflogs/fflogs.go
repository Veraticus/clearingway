package fflogs

import ()

func Init(clientId, clientSecret string) *Fflogs {
	f := &Fflogs{
		clientId:     clientId,
		clientSecret: clientSecret,
	}
	f.SetGraphqlClient()
	return f
}
