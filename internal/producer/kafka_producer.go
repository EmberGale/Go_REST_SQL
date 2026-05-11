package producer

import (
	"GoRestSQL/internal/model"
	"GoRestSQL/pkg/kafka"
	"strconv"
	"time"
)

//TODO: logger

func ProduceMessage(message *model.Payment) error {
	paymentMessage := kafka.PaymentCreatedMessage{
		EventID:   strconv.Itoa(payment.Id),
		EventType: "PaymentCreated",
		Timestamp: time.Now().String(),
		Payment:   model.Payment{},
	}
	err := p.kafkaProducer.SendPaymentCreated(&paymentMessage)
	if err != nil {
		return err
	}
}
