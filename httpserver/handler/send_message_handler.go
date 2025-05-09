package handler

import (
	"encoding/json"
	"github.com/li1553770945/openmcp-discord-bot/cogs"
	model2 "github.com/li1553770945/openmcp-discord-bot/cogs/model"
	"github.com/li1553770945/openmcp-discord-bot/httpserver/constant"
	"github.com/li1553770945/openmcp-discord-bot/httpserver/model"
	"github.com/li1553770945/openmcp-discord-bot/infra/config"
	"go.uber.org/zap"
	"net/http"
)

func SendMessageHandler(w http.ResponseWriter, req *http.Request) {
	response := &model.BasicResponse{
		Code:    0,
		Message: "",
	}

	if req.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		response.Code = constant.MethodNotAllowed
		response.Message = "不支持的请求方式"
		zap.S().Infof("不支持的请求方式:%v", req.Method)
		bytes, err := json.Marshal(response)
		if err != nil {
			zap.S().Errorf("构造http响应失败:%v", err)
			return
		}
		_, err = w.Write(bytes)
		if err != nil {
			zap.S().Errorf("发送http响应失败:%v", err)
		} else {
			return
		}
		return
	}

	token := req.Header.Get("Authorization")
	if token != config.GetConfig().MessageSendToken {
		w.WriteHeader(http.StatusForbidden)
		response.Code = constant.UnAuthorized
		response.Message = "token错误"
		zap.S().Infof("用户提交错误token为:%s", token)
		bytes, err := json.Marshal(response)
		if err != nil {
			zap.S().Errorf("构造http响应失败:%v", err)
			return
		}
		_, err = w.Write(bytes)
		if err != nil {
			zap.S().Errorf("发送http响应失败:%v", err)
		} else {
			return
		}
		return
	}

	// 解析请求体中的JSON数据
	var messageSendReq model2.MessageSendReq
	err := json.NewDecoder(req.Body).Decode(&messageSendReq)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		response.Code = constant.InvalidBody
		response.Message = "请求参数错误"
		bytes, err := json.Marshal(response)
		if err != nil {
			zap.S().Errorf("构造http响应失败:%v", err)
			return
		}
		_, err = w.Write(bytes)
		if err != nil {
			zap.S().Errorf("发送http响应失败:%v", err)
		}
		return
	}

	cogs.GetMessageSendReqChan() <- &messageSendReq

	response.Code = 0
	bytes, err := json.Marshal(response)
	_, err = w.Write(bytes)
	if err != nil {
		zap.S().Errorf("发送http响应失败:%v", err)
	} else {
		return
	}
	return
}
