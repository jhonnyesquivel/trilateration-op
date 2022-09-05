package quasar

import (
	"errors"
)

type Emissor struct {
	position EmmisorPosition
	message  EmissorMessage
}

type EmmisorPosition struct {
	x SatelliteAxis
	y SatelliteAxis
}

func NewEmmisor(message string, x float64, y float64) (Emissor, error) {

	positionVO, err := NewEmmisorPosition(x, y)
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

func NewEmmisorPosition(x float64, y float64) (EmmisorPosition, error) {
	xVO, err := NewAxis(x)
	if err != nil {
		return EmmisorPosition{}, err
	}
	yVO, err := NewAxis(y)
	if err != nil {
		return EmmisorPosition{}, err
	}

	return EmmisorPosition{
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

func (e Emissor) Position() EmmisorPosition {
	return e.position
}

func (e EmmisorPosition) AxisX() SatelliteAxis {
	return e.x
}

func (e EmmisorPosition) AxisY() SatelliteAxis {
	return e.y
}

var ErrXAxisEmpty = errors.New("the X axis can not be empty")
var ErrYAxisEmpty = errors.New("the Y axis can not be empty")
