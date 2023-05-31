// Copyright 2022 The Casdoor Authors. All Rights Reserved.
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

package controllers

import (
	"encoding/json"

	"github.com/casdoor/casdoor/object"
	"github.com/casdoor/casdoor/util"
)

func (c *ApiController) Enforce() {
	permissionId := c.Input().Get("permissionId")
	modelId := c.Input().Get("modelId")
	resourceId := c.Input().Get("resourceId")

	var request object.CasbinRequest
	err := json.Unmarshal(c.Ctx.Input.RequestBody, &request)
	if err != nil {
		c.ResponseError(err.Error())
		return
	}

	if permissionId != "" {
		c.ResponseOk(object.Enforce(permissionId, &request))
		return
	}

	permissions := make([]*object.Permission, 0)
	res := []bool{}

	if modelId != "" {
		owner, modelName := util.GetOwnerAndNameFromId(modelId)
		permissions, err = object.GetPermissionsByModel(owner, modelName)
		if err != nil {
			panic(err)
		}
	} else {
		permissions, err = object.GetPermissionsByResource(resourceId)
		if err != nil {
			panic(err)
		}
	}

	for _, permission := range permissions {
		res = append(res, object.Enforce(permission.GetId(), &request))
	}
	c.Data["json"] = res
	c.ServeJSON()
}

func (c *ApiController) BatchEnforce() {
	permissionId := c.Input().Get("permissionId")
	modelId := c.Input().Get("modelId")

	var requests []object.CasbinRequest
	err := json.Unmarshal(c.Ctx.Input.RequestBody, &requests)
	if err != nil {
		panic(err)
	}

	if permissionId != "" {
		c.Data["json"] = object.BatchEnforce(permissionId, &requests)
		c.ServeJSON()
	} else {
		owner, modelName := util.GetOwnerAndNameFromId(modelId)
		permissions, err := object.GetPermissionsByModel(owner, modelName)
		if err != nil {
			panic(err)
		}

		res := [][]bool{}
		for _, permission := range permissions {
			res = append(res, object.BatchEnforce(permission.GetId(), &requests))
		}

		c.ResponseOk(res)
	}
}

func (c *ApiController) GetAllObjects() {
	userId := c.GetSessionUsername()
	if userId == "" {
		c.ResponseError(c.T("general:Please login first"))
		return
	}

	c.ResponseOk(object.GetAllObjects(userId))
}

func (c *ApiController) GetAllActions() {
	userId := c.GetSessionUsername()
	if userId == "" {
		c.ResponseError(c.T("general:Please login first"))
		return
	}

	c.ResponseOk(object.GetAllActions(userId))
}

func (c *ApiController) GetAllRoles() {
	userId := c.GetSessionUsername()
	if userId == "" {
		c.ResponseError(c.T("general:Please login first"))
		return
	}

	c.ResponseOk(object.GetAllRoles(userId))
}
