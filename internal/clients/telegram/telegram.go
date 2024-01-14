package telegram

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"income_expense_bot/internal/lib"
	"income_expense_bot/internal/lib/e"
	"io"
	"net/http"
	"net/url"
	"path"
	"strconv"
	"time"
)

type Client struct {
	host     string
	basePath string
	client   http.Client
}

const (
	getUpdatesMethod    = "getUpdates"
	sendMessageMethod   = "sendMessage"
	answerCallbackQuery = "answerCallbackQuery"
)

func New(host string, token string) *Client {
	return &Client{
		host:     host,
		basePath: newBasePath(token),
		client:   http.Client{},
	}
}

func (c *Client) Updates(ctx context.Context, offset int, limit int) ([]Update, error) {
	q := url.Values{}
	q.Add("offset", strconv.Itoa(offset))
	q.Add("limit", strconv.Itoa(limit))

	data, err := c.doRequest(ctx, getUpdatesMethod, q, nil)
	if err != nil {
		return nil, err
	}

	var res UpdatesResponse

	if err := json.Unmarshal(data, &res); err != nil {
		return nil, err
	}

	return res.Result, nil
}

func (c *Client) SendMessage(ctx context.Context, outgoingMessage OutgoingMessage) error {
	q, err := outgoingMessage.AsUrlValues()
	if err != nil {
		return e.Wrap("can't send message", err)
	}

	res, err := c.doRequest(ctx, sendMessageMethod, *q, lib.Pointer(2*time.Second))
	if err != nil {
		return e.Wrap("can't send message", err)
	}

	var response TelegramResponse
	err = json.Unmarshal(res, &response)
	if err != nil {
		return e.Wrap("can't send message", err)
	}

	if !response.Ok {
		return errors.New(fmt.Sprintf("can't send message %v", *response.Description))
	}

	return nil
}

func (c *Client) AnswerCallbackQuery(ctx context.Context, callbackMessage OutgoingCallbackMessage) (err error) {
	defer func() {
		if err != nil {
			err = e.Wrap("can't answer on callback query", err)
		}
	}()

	q, err := callbackMessage.AsUrlValues()
	if err != nil {
		return err
	}

	res, err := c.doRequest(ctx, answerCallbackQuery, *q, lib.Pointer(2*time.Second))
	if err != nil {
		return err
	}

	var response TelegramResponse
	err = json.Unmarshal(res, &response)
	if err != nil {
		return err
	}

	if !response.Ok {
		return errors.New(*response.Description)
	}

	return nil
}

func (c *Client) doRequest(ctx context.Context, method string, query url.Values, timeout *time.Duration) (data []byte, err error) {
	defer func() {
		err = e.WrapIfErr("can't do request: %w", err)
	}()

	u := url.URL{
		Scheme: "https",
		Host:   c.host,
		Path:   path.Join(c.basePath, method),
	}

	req, err := http.NewRequest(http.MethodGet, u.String(), nil)
	if err != nil {
		return nil, err
	}

	if timeout != nil {
		var cancel func()
		ctx, cancel = context.WithTimeout(ctx, *timeout)
		defer cancel()
	}

	req = req.WithContext(ctx)

	req.URL.RawQuery = query.Encode()

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer func() {
		_ = resp.Body.Close()
	}()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return body, nil
}

func newBasePath(token string) string {
	return "bot" + token
}
