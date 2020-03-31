package agent

import "encoding/json"

type validationReq struct {
	SvcID  string `json:"service"`
	Client string `json:"client"`
}

func (a *agent) req(r []byte) (validationReq, error) {
	b, err := a.crypto.Decrypt(r)
	if err != nil {
		return validationReq{}, err
	}

	var val validationReq

	if err := json.Unmarshal(b, &val); err != nil {
		return validationReq{}, err
	}
	return val, nil
}
