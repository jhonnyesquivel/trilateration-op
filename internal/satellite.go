package quasar

import (
	"errors"
)

// ---------------------------------------------------------

type Satellite struct {
	name     SateliteName
	distance SateliteDistance
	message  SateliteMessage
}

func NewSatellite(name string, distance float64, message []string) (Satellite, error) {
	nameVO, err := NewSateliteName(name)
	if err != nil {
		return Satellite{}, err
	}

	distanceVO, err := NewSateliteDistance(distance)
	if err != nil {
		return Satellite{}, err
	}

	messageVO, err := NewSateliteMessage(message)
	if err != nil {
		return Satellite{}, err
	}

	return Satellite{
		name:     nameVO,
		distance: distanceVO,
		message:  messageVO,
	}, nil
}

// ---------------------------------------------------------

type SateliteDistance struct {
	value float64
}

type SateliteName struct {
	value string
}

type SateliteMessage struct {
	value []string
}

func NewSateliteDistance(value float64) (SateliteDistance, error) {
	return SateliteDistance{
		value: value,
	}, nil
}

func NewSateliteName(value string) (SateliteName, error) {
	return SateliteName{
		value: value,
	}, nil
}

func NewSateliteMessage(message []string) (SateliteMessage, error) {
	return SateliteMessage{
		value: message,
	}, nil
}

func (distance SateliteDistance) Value() float64 {
	return distance.value
}

func (name SateliteName) Value() string {
	return name.value
}

func (message SateliteMessage) Value() []string {
	return message.value
}

func (s Satellite) Name() SateliteName {
	return s.name
}

func (s Satellite) Distance() SateliteDistance {
	return s.distance
}

func (s Satellite) Message() SateliteMessage {
	return s.message
}

var ErrDistanceEmpty = errors.New("the field distance can not be empty")
var ErrNameEmpty = errors.New("the field name can not be empty")
var ErrMessageEmpty = errors.New("the message parts can not be empty")
