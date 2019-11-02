package rabbitmq

import (
	"encoding/json"
	"testing"
)

const host = "amqp://linen:linen@localhost:5672"

func TestPublish(t *testing.T) {
	q := New(host)
	defer q.Close()
	q.Bind("test")

	q2 := New(host)
	defer q2.Close()
	q2.Bind("test")

	q3 := New(host)
	defer q3.Close()

	except := "fuck"
	q3.Publish("test", except)

	ch := q.Cosume()
	msg := <- ch
	var actual interface{}
	err := json.Unmarshal(msg.Body, &actual)
	if err != nil {
		t.Error(err)
	}
	if actual != except {
		t.Errorf("excepted %s, actual : %s", except, actual)
	}
	if msg.ReplyTo != q3.Name {
		t.Error(msg)
	}

	ch2 := q2.Cosume()
	msg2 := <-ch2
	var actual2 interface{}
	err = json.Unmarshal(msg2.Body, &actual2)
	if err != nil {
		t.Error(err)
	}
	if actual2 != except {
		t.Errorf("excepted %s, actual : %s", except, actual2)
	}
	if msg2.ReplyTo != q3.Name {
		t.Error(msg2)
	}
}

func TestSend(t *testing.T) {
	q := New(host)
	defer q.Close()

	q2 := New(host)
	defer q2.Close()

	except2 := "test2"
	except := "test"
	q.Send(q2.Name, except2)
	q.Send(q.Name, except)

	ch2 := q2.Cosume()
	msg2 := <-ch2
	var actual2 string
	err := json.Unmarshal(msg2.Body, &actual2)
	if err != nil {
		t.Error(err)
	}
	if actual2 != except2 {
		t.Errorf("except : %s, actual : %s", except2, actual2)
	}

	ch := q.Cosume()
	msg := <-ch
	var actual string
	err = json.Unmarshal(msg.Body, &actual)
	if err != nil {
		t.Error(err)
	}
	if actual != except {
		t.Errorf("except : %s, actual : %s", except, actual)
	}
}