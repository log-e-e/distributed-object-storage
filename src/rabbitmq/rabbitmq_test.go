package rabbitmq

import (
    "encoding/json"
    "testing"
)

const (
    mqUrl = "amqp://test:test@localhost:5672"
    testExchange = "testServers"
)

func TestRabbitMQ_Publish(t *testing.T) {
    q1 := NewRabbitMQ(mqUrl)
    defer q1.Close()
    q1.BindExchange(testExchange)

    q2 := NewRabbitMQ(mqUrl)
    defer q2.Close()
    q2.BindExchange(testExchange)

    q3 := NewRabbitMQ(mqUrl)
    defer q3.Close()

    // 发送消息
    expectMsg := "msg delivery by " + testExchange
    q3.Publish(testExchange, expectMsg)
    // 接收消息并回复消息
    c1 := q1.Consume()
    // 阻塞，等待MQ服务器发送消息
    msg := <- c1
    var actualMsg interface{}
    err := json.Unmarshal(msg.Body, &actualMsg)
    if err != nil {
        t.Error(err)
    }
    // 校验消息数据
    if actualMsg != expectMsg {
        t.Errorf("Message error: received '%s', expected '%s'\n", actualMsg, expectMsg)
    } else {
        println("Received message:", actualMsg.(string))
    }
    // 校验回复目标
    if msg.ReplyTo != q3.queueName {
        t.Errorf("ReplyTo error: recieved '%s', expected '%s'\n", msg.ReplyTo, q3.queueName)
    }
    // 回复消息
    replyMsg := "reply message"
    q1.Send(msg.ReplyTo, replyMsg)
    // 接收回复的消息
    c3 := q3.Consume()
    msg = <- c3
    if string(msg.Body) != `"reply message"` {
        t.Errorf("Reply Error: recieved '%s', expected '%s'\n", string(msg.Body), replyMsg)
    }
}

func TestRabbitMQ_Send(t *testing.T) {
    q1 := NewRabbitMQ(mqUrl)
    defer q1.Close()

    q2 := NewRabbitMQ(mqUrl)
    defer q2.Close()

    // 定向发送测试
    expect1 := "test1"
    expect2 := "test2"
    q2.Send(q1.queueName, expect1)
    q2.Send(q2.queueName, expect2)

    c1 := q1.Consume()
    msg := <- c1
    var actualMsg interface{}
    err := json.Unmarshal(msg.Body, &actualMsg)
    if err != nil {
        t.Error(err)
    }
    if actualMsg.(string) != expect1 {
        t.Errorf("Message Error: received '%s', expected '%s'\n", actualMsg.(string), expect1)
    }

    c2 := q2.Consume()
    msg = <- c2
    err = json.Unmarshal(msg.Body, &actualMsg)
    if err != nil {
        t.Error(err)
    }
    if actualMsg.(string) != expect2 {
        t.Errorf("Message Error: received '%s', expected '%s'\n", actualMsg.(string), expect2)
    }
}
