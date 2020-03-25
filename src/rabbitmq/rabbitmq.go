package rabbitmq

import (
    "distributed-object-storage/src/err_utils"
    "encoding/json"
    "github.com/streadway/amqp"
)

type RabbitMQ struct {
    conn *amqp.Connection
    channel *amqp.Channel
    exchangeName string
    queueName string
}

func NewRabbitMQ(mqUrl string) *RabbitMQ {
    // 连接MQ服务器，打开一个channel，客户端通过channel来执行命令
    conn, err := amqp.Dial(mqUrl)
    err_utils.PanicNonNilError(err)
    channel, err := conn.Channel()
    err_utils.PanicNonNilError(err)
    // 创建队列
    queue, err := channel.QueueDeclare(
        "",
        false,
        true,
        false,
        false,
        nil,
        )
    err_utils.PanicNonNilError(err)
    // 创建RabbitMQ结构体
    mq := new(RabbitMQ)
    mq.conn = conn
    mq.channel = channel
    mq.queueName = queue.Name

    return mq
}

func (mq *RabbitMQ) BindExchange(exchangeName string) {
    err := mq.channel.QueueBind(mq.queueName, "", exchangeName, false, nil)
    err_utils.PanicNonNilError(err)
    // 绑定成功再赋值
    mq.exchangeName = exchangeName
}

// Send()用于向特定队列发送消息
func (mq *RabbitMQ) Send(queueName string, messageBody interface{}) {
    msgJsonData, err := json.Marshal(messageBody)
    err_utils.PanicNonNilError(err)
    // 发送消息
    err = mq.channel.Publish("", queueName, false, false,
        amqp.Publishing{
            ReplyTo:         mq.queueName, // ACK确认以便队列删除对应的msg
            Body:            msgJsonData,
        })
    err_utils.PanicNonNilError(err)
}

// Publish()通过向交换器投递新消息，发送给所有绑定了该交换器的队列
func (mq *RabbitMQ) Publish(exchangeName string, messageBody interface{}) {
    msgJsonData, err := json.Marshal(messageBody)
    err_utils.PanicNonNilError(err)
    // 投递消息
    err = mq.channel.Publish(exchangeName, "", false, false,
        amqp.Publishing{
            ReplyTo:         mq.queueName,
            Body:            msgJsonData,
        })
    err_utils.PanicNonNilError(err)
}

func (mq *RabbitMQ) Consume() <-chan amqp.Delivery {
    channel, err := mq.channel.Consume(mq.queueName, "", true, false, false, false, nil)
    err_utils.PanicNonNilError(err)
    return channel
}

func (mq *RabbitMQ) Close() {
    mq.channel.Close()
    mq.conn.Close()
}
