package getters

import "context"

func (r *Registry) tieba(ctx context.Context) ([]Item, error) {
	var res struct {
		Data struct {
			BangTopic struct {
				TopicList []struct {
					TopicID   string `json:"topic_id"`
					TopicName string `json:"topic_name"`
					TopicURL  string `json:"topic_url"`
				} `json:"topic_list"`
			} `json:"bang_topic"`
		} `json:"data"`
	}
	if err := r.getJSON(ctx, "https://tieba.baidu.com/hottopic/browse/topicList", nil, &res); err != nil {
		return nil, err
	}
	out := make([]Item, 0, len(res.Data.BangTopic.TopicList))
	for _, k := range res.Data.BangTopic.TopicList {
		out = append(out, Item{ID: k.TopicID, Title: k.TopicName, URL: k.TopicURL})
	}
	return out, nil
}
