// Copyright 2026 The PipeCD Authors.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package notifier

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"

	"github.com/pipe-cd/pipecd/pkg/config"
	"github.com/pipe-cd/pipecd/pkg/model"
)

func TestWebhookSendEvent(t *testing.T) {
	t.Parallel()

	type received struct {
		contentType string
		signature   string
		body        []byte
	}
	ch := make(chan received, 1)

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, err := io.ReadAll(r.Body)
		assert.NoError(t, err)
		ch <- received{
			contentType: r.Header.Get("Content-Type"),
			signature:   r.Header.Get("PipeCD-Signature"),
			body:        body,
		}
	}))
	defer server.Close()

	sender := newWebhookSender("test", config.NotificationReceiverWebhook{
		URL:            server.URL,
		SignatureKey:   "PipeCD-Signature",
		SignatureValue: "secret",
	}, "", zap.NewNop())

	sender.sendEvent(context.Background(), model.NotificationEvent{
		Type: model.NotificationEventType_EVENT_DEPLOYMENT_TRIGGERED,
	})

	got := <-ch
	assert.Equal(t, "application/json", got.contentType)
	assert.Equal(t, "secret", got.signature)

	var decoded map[string]any
	require.NoError(t, json.Unmarshal(got.body, &decoded))
	assert.Contains(t, decoded, "Type")
}
