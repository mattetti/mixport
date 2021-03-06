package exports

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/erik/mixport/mixpanel"
	"io/ioutil"
	"testing"
)

func TestJSONtreamer(t *testing.T) {
	var expected, output bytes.Buffer

	records := make(chan mixpanel.EventData, 3)

	for i := 0; i < 3; i++ {
		event := make(mixpanel.EventData)
		event[mixpanel.EventIDKey] = fmt.Sprintf("%d", i)
		event["foo"] = "bar,baz"

		records <- event

		var buf bytes.Buffer
		enc := json.NewEncoder(&buf)
		enc.Encode(event)

		expected.Write(buf.Bytes())
	}

	close(records)

	JSONStreamer(&output, records)

	if !bytes.Equal(output.Bytes(), expected.Bytes()) {
		t.Errorf("got (%s), expected(%s)", output.Bytes(), expected.Bytes())
	}
}

func BenchmarkJSONStreamer(b *testing.B) {
	records := make(chan mixpanel.EventData, b.N)

	event := make(mixpanel.EventData)
	event[mixpanel.EventIDKey] = "id"
	event["foo"] = "bar,baz"

	for i := 0; i < b.N; i++ {
		records <- event
	}

	close(records)

	b.ResetTimer()
	JSONStreamer(ioutil.Discard, records)
}
