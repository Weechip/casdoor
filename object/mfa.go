// Copyright 2023 The Casdoor Authors. All Rights Reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package object

import (
	"fmt"

	"github.com/casdoor/casdoor/util"

	"github.com/beego/beego/context"
)

type MfaSessionData struct {
	UserId string
}

type MfaProps struct {
	Id            string   `json:"id"`
	IsPreferred   bool     `json:"isPreferred"`
	AuthType      string   `json:"type" form:"type"`
	Secret        string   `json:"secret,omitempty"`
	CountryCode   string   `json:"countryCode,omitempty"`
	URL           string   `json:"url,omitempty"`
	RecoveryCodes []string `json:"recoveryCodes,omitempty"`
}

type MfaInterface interface {
	SetupVerify(ctx *context.Context, passCode string) error
	Verify(passCode string) error
	Initiate(ctx *context.Context, name1 string, name2 string) (*MfaProps, error)
	Enable(ctx *context.Context, user *User) error
}

const (
	SmsType  = "sms"
	TotpType = "app"
)

const (
	MfaSessionUserId = "MfaSessionUserId"
	NextMfa          = "NextMfa"
	RequiredMfa      = "RequiredMfa"
)

func GetMfaUtil(providerType string, config *MfaProps) MfaInterface {
	switch providerType {
	case SmsType:
		return NewSmsTwoFactor(config)
	case TotpType:
		return nil
	}

	return nil
}

func RecoverTfs(user *User, recoveryCode string) error {
	hit := false

	twoFactor := user.GetPreferMfa(false)
	if len(twoFactor.RecoveryCodes) == 0 {
		return fmt.Errorf("do not have recovery codes")
	}

	for _, code := range twoFactor.RecoveryCodes {
		if code == recoveryCode {
			hit = true
			break
		}
	}
	if !hit {
		return fmt.Errorf("recovery code not found")
	}

	affected, err := UpdateUser(user.GetId(), user, []string{"two_factor_auth"}, user.IsAdminUser())
	if err != nil {
		return err
	}

	if !affected {
		return fmt.Errorf("")
	}
	return nil
}

func GetMaskedProps(props *MfaProps) *MfaProps {
	maskedProps := &MfaProps{
		AuthType:    props.AuthType,
		Id:          props.Id,
		IsPreferred: props.IsPreferred,
	}

	if props.AuthType == SmsType {
		if !util.IsEmailValid(props.Secret) {
			maskedProps.Secret = util.GetMaskedPhone(props.Secret)
		} else {
			maskedProps.Secret = util.GetMaskedEmail(props.Secret)
		}
	}
	return maskedProps
}
