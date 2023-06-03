package main

import (
	"context"
	"math/rand"
	"strings"
	"time"
	// "github.com/TikTokTechImmersion/assignment_demo_2023/rpc-server/kitex_gen/rpc"
	"../kitex_gen/rpc"
)

// IMServiceImpl implements the last service interface defined in the IDL.
type IMServiceImpl struct{}

// func (s *IMServiceImpl) Send(ctx context.Context, req *rpc.SendRequest) (*rpc.SendResponse, error) {
// 	resp := rpc.NewSendResponse()
// 	resp.Code, resp.Msg = areYouLucky()
// 	return resp, nil
// }

// func (s *IMServiceImpl) Pull(ctx context.Context, req *rpc.PullRequest) (*rpc.PullResponse, error) {
// 	resp := rpc.NewPullResponse()
// 	resp.Code, resp.Msg = areYouLucky()
// 	return resp, nil
// }

// func areYouLucky() (int32, string) {
// 	if rand.Int31n(2) == 1 {
// 		return 0, "success"
// 	} else {
// 		return 500, "unlucky"
// 	}
// }

func GetRoomID(chat string) (string, error){
	var roomID string

	lowercase := strings.ToLower(chat)
	senders := strings.Split(lowercase, ":")
	if len(senders) != 2 {
		err := fmt.Errorf("Invalid Chat ID '%s', required format A1:A2, current format: ", chat)
		return "", err
	}
	sender1, sender2 := senders[0], senders[1]
	if comp := strings.compare(sender1, sender2); comp == 1 {
		roomID = fmt.Sprintf("%s:%s", sender2, sender1)
	}
	else {
		roomID = fmt.Sprintf("%s:%s", sender1, sender2)
	}
	return roomID, nil
}

func ValidateRequest(req *rpc.SendRequest){
	senders := strings.Split(req.Message.Chat, ":")
	if len(senders) != 2 {
		err := fmt.Errorf("Invalid Chat ID '%s', required format A1:A2, current format: ", req.Message.GetChat())
		return err
	}
	sender1, sender2 := senders[0], senders[1]

	if req.Message.GetSender() != sender1 && req.Message.GetSender() != sender2 {
		err := fmt.Errorf("sender '%s' not in this chat room", req.Message.GetSender())
		return err
	}
	return nil
}

func (s *IMServiceImpl) Send(ctx context.Context, req *rpc.SendRequest) (*rpc.SendResponse, error){
	if err := ValidateRequest(req); err != nil {
		return nil, err
	}
	timestamp := time.Now().Unix()
	message := &Message{
		Message: req.Message.GetText(),
		Sender: req.Message.GetSender(),
		Timestamp: timestamp
	}

	roomID := GetRoomID(req.Message.GetChat())

	err := rdb.SaveMessage(ctx, roomID, message)
	if err != nil {
		return nil, err
	}
	resp := rpc.NewSendResponse()
	resp.Code, resp.Msg = 0, "success"
    return resp, nil
}

func (s * IMServiceImpl) Pull(ctx context.Context, req rpc.PullRequest) (*rpc.PullResponse, error){
	roomID, err := GetRoomID(req.GetChat())
	if err != nil {
		return nil, err
	}

	start := req.GetCursor()
	end := start + int64(req.GetLimit())

	messages, err := rdb.GetMessage(ctx, roomID, start, end, req.GetReverse())
	if err != nil {
		return nil, err
	}

	respMessage := make([]*rpc.Message, 0)
	var counter int32 = 0
	var nextCursor int64 = 0
	hasMore := false
	for _, msg := range messages {
		if counter+1 > req.GetLimit() {
			hasMore = true
			nextCursor = end
			break
		}
		temp := &rpc.Message{
			Chat: req.GetChat(),
            Text: msg.Message,
            Sender: msg.Sender,
            SendTime: msg.Timestamp
       }
	   resp := rpc.NewPullResponse()
	   resp.Messages = respMessage
	   resp.Code = 0
	   resp.Msg = "success"
	   resp.HasMore = &hasMore
	   resp.NextCursor = &nextCursor

	   return resp, nil
	}
}