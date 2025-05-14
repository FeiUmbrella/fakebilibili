package live

type ReqGetRoom struct {
	Status int    `json:"status" binding:"required"`
	Data   string `json:"data" binding:"required"`
}

type GetLiveRoomInfoReceiveStruct struct {
	RoomID uint `json:"room_id"`
}

type LivestatRes struct {
	//Status int `json:"status"`
	//Data   struct {
	//	Publishers []struct {
	//		Key             string `json:"key"`
	//		Url             string `json:"url"`
	//		StreamId        int    `json:"stream_id"`
	//		VideoTotalBytes int64  `json:"video_total_bytes"`
	//		VideoSpeed      int    `json:"video_speed"`
	//		AudioTotalBytes int    `json:"audio_total_bytes"`
	//		AudioSpeed      int    `json:"audio_speed"`
	//	} `json:"publishers"`
	//	Players interface{} `json:"players"`
	//} `json:"data"`
	HTTPFLV struct {
		Servers []struct {
			Applications []struct {
				Live struct {
					Streams []struct {
						Name string `json:"name"` // 这个就是推流码 room-${uid}
					} `json:"streams"`
				} `json:"live"`
			} `json:"applications"`
		} `json:"servers"`
	} `json:"http-flv"`
}
