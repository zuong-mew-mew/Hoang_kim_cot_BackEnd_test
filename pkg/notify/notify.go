package notify

import (
	"context"
	"fmt"
	"sync"

	firebase "firebase.google.com/go/v4"
	"firebase.google.com/go/v4/messaging"
	"google.golang.org/api/option"
)

// Client bọc Firebase Messaging — nil-safe: chưa cấu hình thì Send trả lỗi rõ ràng.
type Client struct {
	msg *messaging.Client
}

func New(ctx context.Context, credentialsPath string) (*Client, error) {
	if credentialsPath == "" {
		return &Client{}, nil // chưa cấu hình — BE vẫn chạy, chỉ không gửi được
	}
	app, err := firebase.NewApp(ctx, nil, option.WithCredentialsFile(credentialsPath))
	if err != nil {
		return nil, err
	}
	msg, err := app.Messaging(ctx)
	if err != nil {
		return nil, err
	}
	return &Client{msg: msg}, nil
}

func (c *Client) Ready() bool { return c.msg != nil }

func (c *Client) SendToCondition(ctx context.Context, condition, title, body string, data map[string]string) (string, error) {
	if c.msg == nil {
		return "", fmt.Errorf("FCM chưa được cấu hình (thiếu FIREBASE_CREDENTIALS)")
	}
	return c.msg.Send(ctx, &messaging.Message{
		Condition:    condition,
		Notification: &messaging.Notification{Title: title, Body: body},
		Data:         data,
		Android: &messaging.AndroidConfig{
			Priority: "high",
		},
	})
}

// hàm này gửi từng token bằng API v1 thay thế.
// Trả về số máy nhận thành công / lỗi thật.
func (c *Client) SendMulticast(ctx context.Context, tokens []string, title, body string, data map[string]string) (success, failure int, err error) {
	if c.msg == nil {
		return 0, 0, fmt.Errorf("FCM chưa được cấu hình (thiếu FIREBASE_CREDENTIALS)")
	}
	var mu sync.Mutex
	var wg sync.WaitGroup
	sem := make(chan struct{}, 20) // tối đa 20 request cùng lúc
	for _, token := range tokens {
		wg.Add(1)
		go func(token string) {
			defer wg.Done()
			sem <- struct{}{}
			defer func() { <-sem }()
			_, sendErr := c.msg.Send(ctx, &messaging.Message{
				Token:        token,
				Notification: &messaging.Notification{Title: title, Body: body},
				Data:         data,
				Android: &messaging.AndroidConfig{
					Priority: "high",
				},
			})
			mu.Lock()
			if sendErr != nil {
				failure++
			} else {
				success++
			}
			mu.Unlock()
		}(token)
	}
	wg.Wait()
	if success == 0 && failure > 0 {
		return 0, failure, fmt.Errorf("không gửi được tới máy nào (%d token lỗi)", failure)
	}
	return success, failure, nil
}
