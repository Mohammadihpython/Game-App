package protobufEncoder

import (
	"GameApp/contract/goproto/matching"
	"GameApp/entity"
	"GameApp/pkg/slice"
	"encoding/base64"
	"google.golang.org/protobuf/proto"
	"log"
)

func EncodeEvent(event entity.Event, data any) string {
	var payload []byte
	switch event {
	case entity.MatchingUsersMatchedEvent:
		mu, ok := data.(entity.MatchedUsers)
		if ok {
			//	TODO  - log
			//	TODO - update metrics
			return ""
		}
		protobufMU := matching.MatchedUsers{
			Category: string(mu.Category),
			UserIds:  slice.MapFromUintToUint64(mu.UserIDs),
		}
		var err error
		payload, err = proto.Marshal(&protobufMU)
		if err != nil {
			//	TODO  - log
			//	TODO - update metrics
			return ""
		}
		return base64.StdEncoding.EncodeToString(payload)
	default:
		return ""

	}
}

func DecoderEvent(event entity.Event, data string) any {
	payload, err := base64.StdEncoding.DecodeString(data)

	if err != nil {
		log.Println(err)
		return nil
	}

	switch event {
	case entity.MatchingUsersMatchedEvent:
		pbMu := matching.MatchedUsers{}
		if err := proto.Unmarshal(payload, &pbMu); err != nil {
			log.Println(err)
		}
		return entity.MatchedUsers{
			Category: entity.Category(pbMu.Category),
			UserIDs:  slice.MapFromUint64ToUint(pbMu.UserIds),
		}
	}
	return nil

}
