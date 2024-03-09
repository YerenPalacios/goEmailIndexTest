package main

type SyncSearchBool struct {
	Must []map[string]map[string]string `json:"must"`
}

type SyncSearchQuery struct {
	Bool SyncSearchBool `json:"bool"`
}

type ZyncSearchPayload struct {
	From  int             `json:"from"`
	To    int             `json:"to"`
	Query SyncSearchQuery `json:"query"`
	Sort  []string        `json:"sort"`
}

type MessagesService struct{}

func (s *MessagesService) createPayload(
	search string,
	order []string,
	from int,
	to int,
) ZyncSearchPayload {
	boolData := SyncSearchBool{
		Must: []map[string]map[string]string{
			{"match_all": map[string]string{}},
		},
	}
	if search != "" {
		boolData = SyncSearchBool{
			Must: []map[string]map[string]string{
				{"query_string": map[string]string{
					"query": search,
				}},
			},
		}
	}

	payload := ZyncSearchPayload{
		From: from,
		To:   to,
		Query: SyncSearchQuery{
			Bool: boolData,
		},
		Sort: order,
	}
	return payload
}

func (s *MessagesService) getMessages(payload []byte) ([]byte, error) {
	return doPost(ZYNCSARCH_URL+"/es/Messages/_search", payload)
}
