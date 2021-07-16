package api

import (
	"encoding/json"
	"fmt"
	"github.com/amrHassanAbdallah/notificationaway/persistence"
	"github.com/amrHassanAbdallah/notificationaway/service"
	"github.com/amrHassanAbdallah/notificationaway/utils"
	"github.com/go-chi/render"
	"net/http"
)

type ServerInterfaceGlobal interface {
	ServerInterface
}
type server struct {
	NotificationawayService *service.Service
}

func (u *NewMessage) toServiceMessage() (*persistence.Message, error) {
	errorMsg := "invalid message object format"
	jsonbody, err := json.Marshal(u)
	if err != nil {
		return nil, fmt.Errorf(errorMsg)
	}
	val := persistence.Message{}
	if err := json.Unmarshal(jsonbody, &val); err != nil {
		return nil, fmt.Errorf(errorMsg)
	}
	return &val, nil
}
func (s server) AddMessage(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	var nentity NewMessage
	if err := json.NewDecoder(r.Body).Decode(&nentity); err != nil {
		HandleError(w, r, &ValidationError{
			Cause:  err,
			Detail: nil,
			Status: 400,
		})
		return
	}
	serviceEntity, err := nentity.toServiceMessage()
	if err != nil {
		HandleError(w, r, &ValidationError{
			Cause:  err,
			Detail: nil,
			Status: http.StatusBadRequest,
		})
		return
	}
	err = utils.Validator.Struct(serviceEntity)
	if err != nil {
		HandleError(w, r, &ValidationError{
			Cause:  err,
			Detail: nil,
			Status: http.StatusBadRequest,
		})
		return
	}
	resultEntity, err := s.NotificationawayService.AddMessage(ctx, *serviceEntity)
	if err != nil {
		switch err.(type) {
		case *persistence.DuplicateEntityException:
			HandleError(w, r, &service.ServiceError{
				Cause: err,
				Type:  http.StatusConflict,
			})
		default:
			HandleError(w, r, &service.ServiceError{
				Cause: err,
				Type:  http.StatusInternalServerError,
			})
		}
		return
	}

	render.Status(r, http.StatusCreated)
	render.JSON(w, r, MapMessageToResponse(resultEntity))
	return
}
func MapMessageToResponse(b *persistence.Message) MessageResponse {
	return MessageResponse{
		NewMessage: NewMessage{
			Language:     b.Language,
			ProviderType: b.ProviderType,
			Template:     b.Template,
			TemplateKeys: b.TemplateKeys,
			Type:         b.Type,
		},
		CreatedAt: b.CreatedAt,
		Id:        b.Id,
	}
}

func (s server) TriggerMessage(w http.ResponseWriter, r *http.Request) {
	panic("implement me")
}

func (s server) GetMessage(w http.ResponseWriter, r *http.Request, messageId string) {
	panic("implement me")
}

func NewServer(svc *service.Service) ServerInterfaceGlobal {
	return &server{
		NotificationawayService: svc,
	}
}
