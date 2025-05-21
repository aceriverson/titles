package models

type Settings struct {
	AutomaticTitle bool `json:"automatic_title"`
	Tone           int  `json:"tone"`
	Attribution    bool `json:"attribution"`
	Description    bool `json:"description"`
}
