package quasar

import (
	"errors"
)

type Emissor struct {
	position Position
	message  EmissorMessage
}

type Position struct {
	x SatelliteAxis
	y SatelliteAxis
}

func NewEmmisor(message string, x float64, y float64) (Emissor, error) {

	positionVO, err := NewPosition(x, y)
	if err != nil {
		return Emissor{}, err
	}

	msg, err := NewEmissorMessage(message)
	if err != nil {
		return Emissor{}, err
	}

	return Emissor{
		position: positionVO,
		message:  msg,
	}, nil
}

func NewPosition(x float64, y float64) (Position, error) {
	xVO, err := NewAxis(x)
	if err != nil {
		return Position{}, err
	}
	yVO, err := NewAxis(y)
	if err != nil {
		return Position{}, err
	}

	return Position{
		x: xVO,
		y: yVO,
	}, nil
}

type SatelliteAxis struct {
	value float64
}

type EmissorMessage struct {
	value string
}

func NewEmissorMessage(msg string) (EmissorMessage, error) {
	return EmissorMessage{
		value: msg,
	}, nil
}

func NewAxis(value float64) (SatelliteAxis, error) {
	return SatelliteAxis{
		value: value,
	}, nil
}

func (axis SatelliteAxis) Value() float64 {
	return axis.value
}

func (msg EmissorMessage) Value() string {
	return msg.value
}

func (e Emissor) Message() EmissorMessage {
	return e.message
}

func (e Emissor) Position() Position {
	return e.position
}

func (e Position) AxisX() SatelliteAxis {
	return e.x
}

func (e Position) AxisY() SatelliteAxis {
	return e.y
}

var ErrXAxisEmpty = errors.New("the X axis can not be empty")
var ErrYAxisEmpty = errors.New("the Y axis can not be empty")
