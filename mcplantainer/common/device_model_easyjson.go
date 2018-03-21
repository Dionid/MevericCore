// Code generated by easyjson for marshaling/unmarshaling. DO NOT EDIT.

package common

import (
	json "encoding/json"
	easyjson "github.com/mailru/easyjson"
	jlexer "github.com/mailru/easyjson/jlexer"
	jwriter "github.com/mailru/easyjson/jwriter"
	bson "gopkg.in/mgo.v2/bson"
	time "time"
)

// suppress unused package warning
var (
	_ *json.RawMessage
	_ *jlexer.Lexer
	_ *jwriter.Writer
	_ easyjson.Marshaler
)

func easyjson8042eaf4DecodeMevericcoreMcplantainerCommon(in *jlexer.Lexer, out *PlantainersList) {
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
func easyjson8042eaf4EncodeMevericcoreMcplantainerCommon(out *jwriter.Writer, in PlantainersList) {
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
	easyjson8042eaf4EncodeMevericcoreMcplantainerCommon(&w, v)
	return w.Buffer.BuildBytes(), w.Error
}

// MarshalEasyJSON supports easyjson.Marshaler interface
func (v PlantainersList) MarshalEasyJSON(w *jwriter.Writer) {
	easyjson8042eaf4EncodeMevericcoreMcplantainerCommon(w, v)
}

// UnmarshalJSON supports json.Unmarshaler interface
func (v *PlantainersList) UnmarshalJSON(data []byte) error {
	r := jlexer.Lexer{Data: data}
	easyjson8042eaf4DecodeMevericcoreMcplantainerCommon(&r, v)
	return r.Error()
}

// UnmarshalEasyJSON supports easyjson.Unmarshaler interface
func (v *PlantainersList) UnmarshalEasyJSON(l *jlexer.Lexer) {
	easyjson8042eaf4DecodeMevericcoreMcplantainerCommon(l, v)
}
func easyjson8042eaf4DecodeMevericcoreMcplantainerCommon1(in *jlexer.Lexer, out *PlantainerModelSt) {
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
		case "customData":
			easyjson8042eaf4DecodeMevericcoreMcplantainerCommon2(in, &out.CustomData)
		case "customAdminData":
			easyjson8042eaf4DecodeMevericcoreMcplantainerCommon3(in, &out.CustomAdminData)
		case "shadow":
			(out.Shadow).UnmarshalEasyJSON(in)
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
					var v4 bson.ObjectId
					if data := in.Raw(); in.Ok() {
						in.AddError((v4).UnmarshalJSON(data))
					}
					out.OwnersIds = append(out.OwnersIds, v4)
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
func easyjson8042eaf4EncodeMevericcoreMcplantainerCommon1(out *jwriter.Writer, in PlantainerModelSt) {
	out.RawByte('{')
	first := true
	_ = first
	if true {
		if !first {
			out.RawByte(',')
		}
		first = false
		out.RawString("\"customData\":")
		easyjson8042eaf4EncodeMevericcoreMcplantainerCommon2(out, in.CustomData)
	}
	if true {
		if !first {
			out.RawByte(',')
		}
		first = false
		out.RawString("\"customAdminData\":")
		easyjson8042eaf4EncodeMevericcoreMcplantainerCommon3(out, in.CustomAdminData)
	}
	if true {
		if !first {
			out.RawByte(',')
		}
		first = false
		out.RawString("\"shadow\":")
		(in.Shadow).MarshalEasyJSON(out)
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
			for v5, v6 := range in.OwnersIds {
				if v5 > 0 {
					out.RawByte(',')
				}
				out.Raw((v6).MarshalJSON())
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
	easyjson8042eaf4EncodeMevericcoreMcplantainerCommon1(&w, v)
	return w.Buffer.BuildBytes(), w.Error
}

// MarshalEasyJSON supports easyjson.Marshaler interface
func (v PlantainerModelSt) MarshalEasyJSON(w *jwriter.Writer) {
	easyjson8042eaf4EncodeMevericcoreMcplantainerCommon1(w, v)
}

// UnmarshalJSON supports json.Unmarshaler interface
func (v *PlantainerModelSt) UnmarshalJSON(data []byte) error {
	r := jlexer.Lexer{Data: data}
	easyjson8042eaf4DecodeMevericcoreMcplantainerCommon1(&r, v)
	return r.Error()
}

// UnmarshalEasyJSON supports easyjson.Unmarshaler interface
func (v *PlantainerModelSt) UnmarshalEasyJSON(l *jlexer.Lexer) {
	easyjson8042eaf4DecodeMevericcoreMcplantainerCommon1(l, v)
}
func easyjson8042eaf4DecodeMevericcoreMcplantainerCommon3(in *jlexer.Lexer, out *PlantainerCustomAdminData) {
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
func easyjson8042eaf4EncodeMevericcoreMcplantainerCommon3(out *jwriter.Writer, in PlantainerCustomAdminData) {
	out.RawByte('{')
	first := true
	_ = first
	out.RawByte('}')
}
func easyjson8042eaf4DecodeMevericcoreMcplantainerCommon2(in *jlexer.Lexer, out *PlantainerCustomData) {
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
func easyjson8042eaf4EncodeMevericcoreMcplantainerCommon2(out *jwriter.Writer, in PlantainerCustomData) {
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
