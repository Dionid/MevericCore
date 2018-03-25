// Code generated by easyjson for marshaling/unmarshaling. DO NOT EDIT.

package mcplantainer

import (
	json "encoding/json"
	easyjson "github.com/mailru/easyjson"
	jlexer "github.com/mailru/easyjson/jlexer"
	jwriter "github.com/mailru/easyjson/jwriter"
	bson "gopkg.in/mgo.v2/bson"
	mclightmodule "mevericcore/mcmodules/mclightmodule"
	time "time"
)

// suppress unused package warning
var (
	_ *json.RawMessage
	_ *jlexer.Lexer
	_ *jwriter.Writer
	_ easyjson.Marshaler
)

func easyjson6b0df247DecodeMevericcoreMcplantainer(in *jlexer.Lexer, out *PlantainersList) {
	isTopLevel := in.IsStart()
	if in.IsNull() {
		in.Skip()
		*out = nil
	} else {
		in.Delim('[')
		if *out == nil {
			if !in.IsDelim(']') {
				*out = make(PlantainersList, 0, 1)
			} else {
				*out = PlantainersList{}
			}
		} else {
			*out = (*out)[:0]
		}
		for !in.IsDelim(']') {
			var v1 PlantainerModelSt
			(v1).UnmarshalEasyJSON(in)
			*out = append(*out, v1)
			in.WantComma()
		}
		in.Delim(']')
	}
	if isTopLevel {
		in.Consumed()
	}
}
func easyjson6b0df247EncodeMevericcoreMcplantainer(out *jwriter.Writer, in PlantainersList) {
	if in == nil && (out.Flags&jwriter.NilSliceAsEmpty) == 0 {
		out.RawString("null")
	} else {
		out.RawByte('[')
		for v2, v3 := range in {
			if v2 > 0 {
				out.RawByte(',')
			}
			(v3).MarshalEasyJSON(out)
		}
		out.RawByte(']')
	}
}

// MarshalJSON supports json.Marshaler interface
func (v PlantainersList) MarshalJSON() ([]byte, error) {
	w := jwriter.Writer{}
	easyjson6b0df247EncodeMevericcoreMcplantainer(&w, v)
	return w.Buffer.BuildBytes(), w.Error
}

// MarshalEasyJSON supports easyjson.Marshaler interface
func (v PlantainersList) MarshalEasyJSON(w *jwriter.Writer) {
	easyjson6b0df247EncodeMevericcoreMcplantainer(w, v)
}

// UnmarshalJSON supports json.Unmarshaler interface
func (v *PlantainersList) UnmarshalJSON(data []byte) error {
	r := jlexer.Lexer{Data: data}
	easyjson6b0df247DecodeMevericcoreMcplantainer(&r, v)
	return r.Error()
}

// UnmarshalEasyJSON supports easyjson.Unmarshaler interface
func (v *PlantainersList) UnmarshalEasyJSON(l *jlexer.Lexer) {
	easyjson6b0df247DecodeMevericcoreMcplantainer(l, v)
}
func easyjson6b0df247DecodeMevericcoreMcplantainer1(in *jlexer.Lexer, out *PlantainerShadowStateSt) {
	isTopLevel := in.IsStart()
	if in.IsNull() {
		if isTopLevel {
			in.Consumed()
		}
		in.Skip()
		return
	}
	in.Delim('{')
	for !in.IsDelim('}') {
		key := in.UnsafeString()
		in.WantColon()
		if in.IsNull() {
			in.Skip()
			in.WantComma()
			continue
		}
		switch key {
		case "reported":
			(out.Reported).UnmarshalEasyJSON(in)
		case "desired":
			if in.IsNull() {
				in.Skip()
				out.Desired = nil
			} else {
				if out.Desired == nil {
					out.Desired = new(PlantainerShadowStatePieceSt)
				}
				(*out.Desired).UnmarshalEasyJSON(in)
			}
		case "delta":
			if in.IsNull() {
				in.Skip()
				out.Delta = nil
			} else {
				if out.Delta == nil {
					out.Delta = new(PlantainerShadowStatePieceSt)
				}
				(*out.Delta).UnmarshalEasyJSON(in)
			}
		default:
			in.SkipRecursive()
		}
		in.WantComma()
	}
	in.Delim('}')
	if isTopLevel {
		in.Consumed()
	}
}
func easyjson6b0df247EncodeMevericcoreMcplantainer1(out *jwriter.Writer, in PlantainerShadowStateSt) {
	out.RawByte('{')
	first := true
	_ = first
	if true {
		if !first {
			out.RawByte(',')
		}
		first = false
		out.RawString("\"reported\":")
		(in.Reported).MarshalEasyJSON(out)
	}
	if in.Desired != nil {
		if !first {
			out.RawByte(',')
		}
		first = false
		out.RawString("\"desired\":")
		if in.Desired == nil {
			out.RawString("null")
		} else {
			(*in.Desired).MarshalEasyJSON(out)
		}
	}
	if in.Delta != nil {
		if !first {
			out.RawByte(',')
		}
		first = false
		out.RawString("\"delta\":")
		if in.Delta == nil {
			out.RawString("null")
		} else {
			(*in.Delta).MarshalEasyJSON(out)
		}
	}
	out.RawByte('}')
}

// MarshalJSON supports json.Marshaler interface
func (v PlantainerShadowStateSt) MarshalJSON() ([]byte, error) {
	w := jwriter.Writer{}
	easyjson6b0df247EncodeMevericcoreMcplantainer1(&w, v)
	return w.Buffer.BuildBytes(), w.Error
}

// MarshalEasyJSON supports easyjson.Marshaler interface
func (v PlantainerShadowStateSt) MarshalEasyJSON(w *jwriter.Writer) {
	easyjson6b0df247EncodeMevericcoreMcplantainer1(w, v)
}

// UnmarshalJSON supports json.Unmarshaler interface
func (v *PlantainerShadowStateSt) UnmarshalJSON(data []byte) error {
	r := jlexer.Lexer{Data: data}
	easyjson6b0df247DecodeMevericcoreMcplantainer1(&r, v)
	return r.Error()
}

// UnmarshalEasyJSON supports easyjson.Unmarshaler interface
func (v *PlantainerShadowStateSt) UnmarshalEasyJSON(l *jlexer.Lexer) {
	easyjson6b0df247DecodeMevericcoreMcplantainer1(l, v)
}
func easyjson6b0df247DecodeMevericcoreMcplantainer2(in *jlexer.Lexer, out *PlantainerShadowStatePieceSt) {
	isTopLevel := in.IsStart()
	if in.IsNull() {
		if isTopLevel {
			in.Consumed()
		}
		in.Skip()
		return
	}
	in.Delim('{')
	for !in.IsDelim('}') {
		key := in.UnsafeString()
		in.WantColon()
		if in.IsNull() {
			in.Skip()
			in.WantComma()
			continue
		}
		switch key {
		case "lightModule":
			easyjson6b0df247DecodeMevericcoreMcplantainer3(in, &out.LightModule)
		default:
			in.SkipRecursive()
		}
		in.WantComma()
	}
	in.Delim('}')
	if isTopLevel {
		in.Consumed()
	}
}
func easyjson6b0df247EncodeMevericcoreMcplantainer2(out *jwriter.Writer, in PlantainerShadowStatePieceSt) {
	out.RawByte('{')
	first := true
	_ = first
	if true {
		if !first {
			out.RawByte(',')
		}
		first = false
		out.RawString("\"lightModule\":")
		easyjson6b0df247EncodeMevericcoreMcplantainer3(out, in.LightModule)
	}
	out.RawByte('}')
}

// MarshalJSON supports json.Marshaler interface
func (v PlantainerShadowStatePieceSt) MarshalJSON() ([]byte, error) {
	w := jwriter.Writer{}
	easyjson6b0df247EncodeMevericcoreMcplantainer2(&w, v)
	return w.Buffer.BuildBytes(), w.Error
}

// MarshalEasyJSON supports easyjson.Marshaler interface
func (v PlantainerShadowStatePieceSt) MarshalEasyJSON(w *jwriter.Writer) {
	easyjson6b0df247EncodeMevericcoreMcplantainer2(w, v)
}

// UnmarshalJSON supports json.Unmarshaler interface
func (v *PlantainerShadowStatePieceSt) UnmarshalJSON(data []byte) error {
	r := jlexer.Lexer{Data: data}
	easyjson6b0df247DecodeMevericcoreMcplantainer2(&r, v)
	return r.Error()
}

// UnmarshalEasyJSON supports easyjson.Unmarshaler interface
func (v *PlantainerShadowStatePieceSt) UnmarshalEasyJSON(l *jlexer.Lexer) {
	easyjson6b0df247DecodeMevericcoreMcplantainer2(l, v)
}
func easyjson6b0df247DecodeMevericcoreMcplantainer3(in *jlexer.Lexer, out *PlantainerLightModuleStateSt) {
	isTopLevel := in.IsStart()
	if in.IsNull() {
		if isTopLevel {
			in.Consumed()
		}
		in.Skip()
		return
	}
	in.Delim('{')
	for !in.IsDelim('}') {
		key := in.UnsafeString()
		in.WantColon()
		if in.IsNull() {
			in.Skip()
			in.WantComma()
			continue
		}
		switch key {
		case "mode":
			if in.IsNull() {
				in.Skip()
				out.Mode = nil
			} else {
				if out.Mode == nil {
					out.Mode = new(string)
				}
				*out.Mode = string(in.String())
			}
		case "lightTurnedOn":
			if in.IsNull() {
				in.Skip()
				out.LightTurnedOn = nil
			} else {
				if out.LightTurnedOn == nil {
					out.LightTurnedOn = new(bool)
				}
				*out.LightTurnedOn = bool(in.Bool())
			}
		case "lightLvlCheckActive":
			if in.IsNull() {
				in.Skip()
				out.LightLvlCheckActive = nil
			} else {
				if out.LightLvlCheckActive == nil {
					out.LightLvlCheckActive = new(bool)
				}
				*out.LightLvlCheckActive = bool(in.Bool())
			}
		case "lightLvlCheckInterval":
			if in.IsNull() {
				in.Skip()
				out.LightLvlCheckInterval = nil
			} else {
				if out.LightLvlCheckInterval == nil {
					out.LightLvlCheckInterval = new(int)
				}
				*out.LightLvlCheckInterval = int(in.Int())
			}
		case "lightLvlCheckLastIntervalCallTimestamp":
			if in.IsNull() {
				in.Skip()
				out.LightLvlCheckLastIntervalCallTimestamp = nil
			} else {
				if out.LightLvlCheckLastIntervalCallTimestamp == nil {
					out.LightLvlCheckLastIntervalCallTimestamp = new(int)
				}
				*out.LightLvlCheckLastIntervalCallTimestamp = int(in.Int())
			}
		case "lightIntervalsArr":
			if in.IsNull() {
				in.Skip()
				out.LightIntervalsArr = nil
			} else {
				if out.LightIntervalsArr == nil {
					out.LightIntervalsArr = new([]mclightmodule.LightModuleInterval)
				}
				if in.IsNull() {
					in.Skip()
					*out.LightIntervalsArr = nil
				} else {
					in.Delim('[')
					if *out.LightIntervalsArr == nil {
						if !in.IsDelim(']') {
							*out.LightIntervalsArr = make([]mclightmodule.LightModuleInterval, 0, 1)
						} else {
							*out.LightIntervalsArr = []mclightmodule.LightModuleInterval{}
						}
					} else {
						*out.LightIntervalsArr = (*out.LightIntervalsArr)[:0]
					}
					for !in.IsDelim(']') {
						var v4 mclightmodule.LightModuleInterval
						easyjson6b0df247DecodeMevericcoreMclightmodule(in, &v4)
						*out.LightIntervalsArr = append(*out.LightIntervalsArr, v4)
						in.WantComma()
					}
					in.Delim(']')
				}
			}
		case "lightIntervalsRestTimeTurnedOn":
			if in.IsNull() {
				in.Skip()
				out.LightIntervalsRestTimeTurnedOn = nil
			} else {
				if out.LightIntervalsRestTimeTurnedOn == nil {
					out.LightIntervalsRestTimeTurnedOn = new(bool)
				}
				*out.LightIntervalsRestTimeTurnedOn = bool(in.Bool())
			}
		case "lightIntervalsCheckingInterval":
			if in.IsNull() {
				in.Skip()
				out.LightIntervalsCheckingInterval = nil
			} else {
				if out.LightIntervalsCheckingInterval == nil {
					out.LightIntervalsCheckingInterval = new(int)
				}
				*out.LightIntervalsCheckingInterval = int(in.Int())
			}
		case "lightLvl":
			if in.IsNull() {
				in.Skip()
				out.LightLvl = nil
			} else {
				if out.LightLvl == nil {
					out.LightLvl = new(int)
				}
				*out.LightLvl = int(in.Int())
			}
		default:
			in.SkipRecursive()
		}
		in.WantComma()
	}
	in.Delim('}')
	if isTopLevel {
		in.Consumed()
	}
}
func easyjson6b0df247EncodeMevericcoreMcplantainer3(out *jwriter.Writer, in PlantainerLightModuleStateSt) {
	out.RawByte('{')
	first := true
	_ = first
	if in.Mode != nil {
		if !first {
			out.RawByte(',')
		}
		first = false
		out.RawString("\"mode\":")
		if in.Mode == nil {
			out.RawString("null")
		} else {
			out.String(string(*in.Mode))
		}
	}
	if in.LightTurnedOn != nil {
		if !first {
			out.RawByte(',')
		}
		first = false
		out.RawString("\"lightTurnedOn\":")
		if in.LightTurnedOn == nil {
			out.RawString("null")
		} else {
			out.Bool(bool(*in.LightTurnedOn))
		}
	}
	if in.LightLvlCheckActive != nil {
		if !first {
			out.RawByte(',')
		}
		first = false
		out.RawString("\"lightLvlCheckActive\":")
		if in.LightLvlCheckActive == nil {
			out.RawString("null")
		} else {
			out.Bool(bool(*in.LightLvlCheckActive))
		}
	}
	if in.LightLvlCheckInterval != nil {
		if !first {
			out.RawByte(',')
		}
		first = false
		out.RawString("\"lightLvlCheckInterval\":")
		if in.LightLvlCheckInterval == nil {
			out.RawString("null")
		} else {
			out.Int(int(*in.LightLvlCheckInterval))
		}
	}
	if in.LightLvlCheckLastIntervalCallTimestamp != nil {
		if !first {
			out.RawByte(',')
		}
		first = false
		out.RawString("\"lightLvlCheckLastIntervalCallTimestamp\":")
		if in.LightLvlCheckLastIntervalCallTimestamp == nil {
			out.RawString("null")
		} else {
			out.Int(int(*in.LightLvlCheckLastIntervalCallTimestamp))
		}
	}
	if in.LightIntervalsArr != nil {
		if !first {
			out.RawByte(',')
		}
		first = false
		out.RawString("\"lightIntervalsArr\":")
		if in.LightIntervalsArr == nil {
			out.RawString("null")
		} else {
			if *in.LightIntervalsArr == nil && (out.Flags&jwriter.NilSliceAsEmpty) == 0 {
				out.RawString("null")
			} else {
				out.RawByte('[')
				for v5, v6 := range *in.LightIntervalsArr {
					if v5 > 0 {
						out.RawByte(',')
					}
					easyjson6b0df247EncodeMevericcoreMclightmodule(out, v6)
				}
				out.RawByte(']')
			}
		}
	}
	if in.LightIntervalsRestTimeTurnedOn != nil {
		if !first {
			out.RawByte(',')
		}
		first = false
		out.RawString("\"lightIntervalsRestTimeTurnedOn\":")
		if in.LightIntervalsRestTimeTurnedOn == nil {
			out.RawString("null")
		} else {
			out.Bool(bool(*in.LightIntervalsRestTimeTurnedOn))
		}
	}
	if in.LightIntervalsCheckingInterval != nil {
		if !first {
			out.RawByte(',')
		}
		first = false
		out.RawString("\"lightIntervalsCheckingInterval\":")
		if in.LightIntervalsCheckingInterval == nil {
			out.RawString("null")
		} else {
			out.Int(int(*in.LightIntervalsCheckingInterval))
		}
	}
	if in.LightLvl != nil {
		if !first {
			out.RawByte(',')
		}
		first = false
		out.RawString("\"lightLvl\":")
		if in.LightLvl == nil {
			out.RawString("null")
		} else {
			out.Int(int(*in.LightLvl))
		}
	}
	out.RawByte('}')
}
func easyjson6b0df247DecodeMevericcoreMclightmodule(in *jlexer.Lexer, out *mclightmodule.LightModuleInterval) {
	isTopLevel := in.IsStart()
	if in.IsNull() {
		if isTopLevel {
			in.Consumed()
		}
		in.Skip()
		return
	}
	in.Delim('{')
	for !in.IsDelim('}') {
		key := in.UnsafeString()
		in.WantColon()
		if in.IsNull() {
			in.Skip()
			in.WantComma()
			continue
		}
		switch key {
		case "fromTimeHours":
			out.FromTimeHours = int(in.Int())
		case "fromTimeMinutes":
			out.FromTimeMinutes = int(in.Int())
		case "toTimeHours":
			out.ToTimeHours = int(in.Int())
		case "toTimeMinutes":
			out.ToTimeMinutes = int(in.Int())
		case "turnedOn":
			out.TurnedOn = bool(in.Bool())
		default:
			in.SkipRecursive()
		}
		in.WantComma()
	}
	in.Delim('}')
	if isTopLevel {
		in.Consumed()
	}
}
func easyjson6b0df247EncodeMevericcoreMclightmodule(out *jwriter.Writer, in mclightmodule.LightModuleInterval) {
	out.RawByte('{')
	first := true
	_ = first
	if in.FromTimeHours != 0 {
		if !first {
			out.RawByte(',')
		}
		first = false
		out.RawString("\"fromTimeHours\":")
		out.Int(int(in.FromTimeHours))
	}
	if in.FromTimeMinutes != 0 {
		if !first {
			out.RawByte(',')
		}
		first = false
		out.RawString("\"fromTimeMinutes\":")
		out.Int(int(in.FromTimeMinutes))
	}
	if in.ToTimeHours != 0 {
		if !first {
			out.RawByte(',')
		}
		first = false
		out.RawString("\"toTimeHours\":")
		out.Int(int(in.ToTimeHours))
	}
	if in.ToTimeMinutes != 0 {
		if !first {
			out.RawByte(',')
		}
		first = false
		out.RawString("\"toTimeMinutes\":")
		out.Int(int(in.ToTimeMinutes))
	}
	if in.TurnedOn {
		if !first {
			out.RawByte(',')
		}
		first = false
		out.RawString("\"turnedOn\":")
		out.Bool(bool(in.TurnedOn))
	}
	out.RawByte('}')
}
func easyjson6b0df247DecodeMevericcoreMcplantainer4(in *jlexer.Lexer, out *PlantainerShadowSt) {
	isTopLevel := in.IsStart()
	if in.IsNull() {
		if isTopLevel {
			in.Consumed()
		}
		in.Skip()
		return
	}
	in.Delim('{')
	for !in.IsDelim('}') {
		key := in.UnsafeString()
		in.WantColon()
		if in.IsNull() {
			in.Skip()
			in.WantComma()
			continue
		}
		switch key {
		case "id":
			out.Id = string(in.String())
		case "state":
			(out.State).UnmarshalEasyJSON(in)
		case "metadata":
			easyjson6b0df247DecodeMevericcoreMcplantainer5(in, &out.Metadata)
		default:
			in.SkipRecursive()
		}
		in.WantComma()
	}
	in.Delim('}')
	if isTopLevel {
		in.Consumed()
	}
}
func easyjson6b0df247EncodeMevericcoreMcplantainer4(out *jwriter.Writer, in PlantainerShadowSt) {
	out.RawByte('{')
	first := true
	_ = first
	if in.Id != "" {
		if !first {
			out.RawByte(',')
		}
		first = false
		out.RawString("\"id\":")
		out.String(string(in.Id))
	}
	if true {
		if !first {
			out.RawByte(',')
		}
		first = false
		out.RawString("\"state\":")
		(in.State).MarshalEasyJSON(out)
	}
	if true {
		if !first {
			out.RawByte(',')
		}
		first = false
		out.RawString("\"metadata\":")
		easyjson6b0df247EncodeMevericcoreMcplantainer5(out, in.Metadata)
	}
	out.RawByte('}')
}

// MarshalJSON supports json.Marshaler interface
func (v PlantainerShadowSt) MarshalJSON() ([]byte, error) {
	w := jwriter.Writer{}
	easyjson6b0df247EncodeMevericcoreMcplantainer4(&w, v)
	return w.Buffer.BuildBytes(), w.Error
}

// MarshalEasyJSON supports easyjson.Marshaler interface
func (v PlantainerShadowSt) MarshalEasyJSON(w *jwriter.Writer) {
	easyjson6b0df247EncodeMevericcoreMcplantainer4(w, v)
}

// UnmarshalJSON supports json.Unmarshaler interface
func (v *PlantainerShadowSt) UnmarshalJSON(data []byte) error {
	r := jlexer.Lexer{Data: data}
	easyjson6b0df247DecodeMevericcoreMcplantainer4(&r, v)
	return r.Error()
}

// UnmarshalEasyJSON supports easyjson.Unmarshaler interface
func (v *PlantainerShadowSt) UnmarshalEasyJSON(l *jlexer.Lexer) {
	easyjson6b0df247DecodeMevericcoreMcplantainer4(l, v)
}
func easyjson6b0df247DecodeMevericcoreMcplantainer5(in *jlexer.Lexer, out *PlantainerShadowMetadataSt) {
	isTopLevel := in.IsStart()
	if in.IsNull() {
		if isTopLevel {
			in.Consumed()
		}
		in.Skip()
		return
	}
	in.Delim('{')
	for !in.IsDelim('}') {
		key := in.UnsafeString()
		in.WantColon()
		if in.IsNull() {
			in.Skip()
			in.WantComma()
			continue
		}
		switch key {
		case "version":
			out.Version = int(in.Int())
		default:
			in.SkipRecursive()
		}
		in.WantComma()
	}
	in.Delim('}')
	if isTopLevel {
		in.Consumed()
	}
}
func easyjson6b0df247EncodeMevericcoreMcplantainer5(out *jwriter.Writer, in PlantainerShadowMetadataSt) {
	out.RawByte('{')
	first := true
	_ = first
	if in.Version != 0 {
		if !first {
			out.RawByte(',')
		}
		first = false
		out.RawString("\"version\":")
		out.Int(int(in.Version))
	}
	out.RawByte('}')
}
func easyjson6b0df247DecodeMevericcoreMcplantainer6(in *jlexer.Lexer, out *PlantainerModelSt) {
	isTopLevel := in.IsStart()
	if in.IsNull() {
		if isTopLevel {
			in.Consumed()
		}
		in.Skip()
		return
	}
	in.Delim('{')
	for !in.IsDelim('}') {
		key := in.UnsafeString()
		in.WantColon()
		if in.IsNull() {
			in.Skip()
			in.WantComma()
			continue
		}
		switch key {
		case "shadow":
			(out.Shadow).UnmarshalEasyJSON(in)
		case "customData":
			easyjson6b0df247DecodeMevericcoreMcplantainer7(in, &out.CustomData)
		case "customAdminData":
			easyjson6b0df247DecodeMevericcoreMcplantainer8(in, &out.CustomAdminData)
		case "srcId":
			out.Src = string(in.String())
		case "type":
			out.Type = string(in.String())
		case "firstActivation":
			if in.IsNull() {
				in.Skip()
				out.FirstActivation = nil
			} else {
				if out.FirstActivation == nil {
					out.FirstActivation = new(time.Time)
				}
				if data := in.Raw(); in.Ok() {
					in.AddError((*out.FirstActivation).UnmarshalJSON(data))
				}
			}
		case "lastSeenOnline":
			if in.IsNull() {
				in.Skip()
				out.LastSeenOnline = nil
			} else {
				if out.LastSeenOnline == nil {
					out.LastSeenOnline = new(time.Time)
				}
				if data := in.Raw(); in.Ok() {
					in.AddError((*out.LastSeenOnline).UnmarshalJSON(data))
				}
			}
		case "isOnline":
			if in.IsNull() {
				in.Skip()
				out.IsOnline = nil
			} else {
				if out.IsOnline == nil {
					out.IsOnline = new(bool)
				}
				*out.IsOnline = bool(in.Bool())
			}
		case "ownersIds":
			if in.IsNull() {
				in.Skip()
				out.OwnersIds = nil
			} else {
				in.Delim('[')
				if out.OwnersIds == nil {
					if !in.IsDelim(']') {
						out.OwnersIds = make([]bson.ObjectId, 0, 4)
					} else {
						out.OwnersIds = []bson.ObjectId{}
					}
				} else {
					out.OwnersIds = (out.OwnersIds)[:0]
				}
				for !in.IsDelim(']') {
					var v7 bson.ObjectId
					if data := in.Raw(); in.Ok() {
						in.AddError((v7).UnmarshalJSON(data))
					}
					out.OwnersIds = append(out.OwnersIds, v7)
					in.WantComma()
				}
				in.Delim(']')
			}
		case "id":
			if data := in.Raw(); in.Ok() {
				in.AddError((out.ID).UnmarshalJSON(data))
			}
		case "updatedAt":
			if data := in.Raw(); in.Ok() {
				in.AddError((out.UpdatedAt).UnmarshalJSON(data))
			}
		case "deletedAt":
			if in.IsNull() {
				in.Skip()
				out.DeletedAt = nil
			} else {
				if out.DeletedAt == nil {
					out.DeletedAt = new(time.Time)
				}
				if data := in.Raw(); in.Ok() {
					in.AddError((*out.DeletedAt).UnmarshalJSON(data))
				}
			}
		default:
			in.SkipRecursive()
		}
		in.WantComma()
	}
	in.Delim('}')
	if isTopLevel {
		in.Consumed()
	}
}
func easyjson6b0df247EncodeMevericcoreMcplantainer6(out *jwriter.Writer, in PlantainerModelSt) {
	out.RawByte('{')
	first := true
	_ = first
	if true {
		if !first {
			out.RawByte(',')
		}
		first = false
		out.RawString("\"shadow\":")
		(in.Shadow).MarshalEasyJSON(out)
	}
	if true {
		if !first {
			out.RawByte(',')
		}
		first = false
		out.RawString("\"customData\":")
		easyjson6b0df247EncodeMevericcoreMcplantainer7(out, in.CustomData)
	}
	if true {
		if !first {
			out.RawByte(',')
		}
		first = false
		out.RawString("\"customAdminData\":")
		easyjson6b0df247EncodeMevericcoreMcplantainer8(out, in.CustomAdminData)
	}
	if in.Src != "" {
		if !first {
			out.RawByte(',')
		}
		first = false
		out.RawString("\"srcId\":")
		out.String(string(in.Src))
	}
	if in.Type != "" {
		if !first {
			out.RawByte(',')
		}
		first = false
		out.RawString("\"type\":")
		out.String(string(in.Type))
	}
	if in.FirstActivation != nil {
		if !first {
			out.RawByte(',')
		}
		first = false
		out.RawString("\"firstActivation\":")
		if in.FirstActivation == nil {
			out.RawString("null")
		} else {
			out.Raw((*in.FirstActivation).MarshalJSON())
		}
	}
	if in.LastSeenOnline != nil {
		if !first {
			out.RawByte(',')
		}
		first = false
		out.RawString("\"lastSeenOnline\":")
		if in.LastSeenOnline == nil {
			out.RawString("null")
		} else {
			out.Raw((*in.LastSeenOnline).MarshalJSON())
		}
	}
	if in.IsOnline != nil {
		if !first {
			out.RawByte(',')
		}
		first = false
		out.RawString("\"isOnline\":")
		if in.IsOnline == nil {
			out.RawString("null")
		} else {
			out.Bool(bool(*in.IsOnline))
		}
	}
	if len(in.OwnersIds) != 0 {
		if !first {
			out.RawByte(',')
		}
		first = false
		out.RawString("\"ownersIds\":")
		if in.OwnersIds == nil && (out.Flags&jwriter.NilSliceAsEmpty) == 0 {
			out.RawString("null")
		} else {
			out.RawByte('[')
			for v8, v9 := range in.OwnersIds {
				if v8 > 0 {
					out.RawByte(',')
				}
				out.Raw((v9).MarshalJSON())
			}
			out.RawByte(']')
		}
	}
	if in.ID != "" {
		if !first {
			out.RawByte(',')
		}
		first = false
		out.RawString("\"id\":")
		out.Raw((in.ID).MarshalJSON())
	}
	if true {
		if !first {
			out.RawByte(',')
		}
		first = false
		out.RawString("\"updatedAt\":")
		out.Raw((in.UpdatedAt).MarshalJSON())
	}
	if in.DeletedAt != nil {
		if !first {
			out.RawByte(',')
		}
		first = false
		out.RawString("\"deletedAt\":")
		if in.DeletedAt == nil {
			out.RawString("null")
		} else {
			out.Raw((*in.DeletedAt).MarshalJSON())
		}
	}
	out.RawByte('}')
}

// MarshalJSON supports json.Marshaler interface
func (v PlantainerModelSt) MarshalJSON() ([]byte, error) {
	w := jwriter.Writer{}
	easyjson6b0df247EncodeMevericcoreMcplantainer6(&w, v)
	return w.Buffer.BuildBytes(), w.Error
}

// MarshalEasyJSON supports easyjson.Marshaler interface
func (v PlantainerModelSt) MarshalEasyJSON(w *jwriter.Writer) {
	easyjson6b0df247EncodeMevericcoreMcplantainer6(w, v)
}

// UnmarshalJSON supports json.Unmarshaler interface
func (v *PlantainerModelSt) UnmarshalJSON(data []byte) error {
	r := jlexer.Lexer{Data: data}
	easyjson6b0df247DecodeMevericcoreMcplantainer6(&r, v)
	return r.Error()
}

// UnmarshalEasyJSON supports easyjson.Unmarshaler interface
func (v *PlantainerModelSt) UnmarshalEasyJSON(l *jlexer.Lexer) {
	easyjson6b0df247DecodeMevericcoreMcplantainer6(l, v)
}
func easyjson6b0df247DecodeMevericcoreMcplantainer8(in *jlexer.Lexer, out *PlantainerCustomAdminData) {
	isTopLevel := in.IsStart()
	if in.IsNull() {
		if isTopLevel {
			in.Consumed()
		}
		in.Skip()
		return
	}
	in.Delim('{')
	for !in.IsDelim('}') {
		key := in.UnsafeString()
		in.WantColon()
		if in.IsNull() {
			in.Skip()
			in.WantComma()
			continue
		}
		switch key {
		default:
			in.SkipRecursive()
		}
		in.WantComma()
	}
	in.Delim('}')
	if isTopLevel {
		in.Consumed()
	}
}
func easyjson6b0df247EncodeMevericcoreMcplantainer8(out *jwriter.Writer, in PlantainerCustomAdminData) {
	out.RawByte('{')
	first := true
	_ = first
	out.RawByte('}')
}
func easyjson6b0df247DecodeMevericcoreMcplantainer7(in *jlexer.Lexer, out *PlantainerCustomData) {
	isTopLevel := in.IsStart()
	if in.IsNull() {
		if isTopLevel {
			in.Consumed()
		}
		in.Skip()
		return
	}
	in.Delim('{')
	for !in.IsDelim('}') {
		key := in.UnsafeString()
		in.WantColon()
		if in.IsNull() {
			in.Skip()
			in.WantComma()
			continue
		}
		switch key {
		case "name":
			out.Name = string(in.String())
		default:
			in.SkipRecursive()
		}
		in.WantComma()
	}
	in.Delim('}')
	if isTopLevel {
		in.Consumed()
	}
}
func easyjson6b0df247EncodeMevericcoreMcplantainer7(out *jwriter.Writer, in PlantainerCustomData) {
	out.RawByte('{')
	first := true
	_ = first
	if in.Name != "" {
		if !first {
			out.RawByte(',')
		}
		first = false
		out.RawString("\"name\":")
		out.String(string(in.Name))
	}
	out.RawByte('}')
}
