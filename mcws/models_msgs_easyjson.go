// Code generated by easyjson for marshaling/unmarshaling. DO NOT EDIT.

package mcws

import (
	json "encoding/json"
	easyjson "github.com/mailru/easyjson"
	jlexer "github.com/mailru/easyjson/jlexer"
	jwriter "github.com/mailru/easyjson/jwriter"
)

// suppress unused package warning
var (
	_ *json.RawMessage
	_ *jlexer.Lexer
	_ *jwriter.Writer
	_ easyjson.Marshaler
)

func easyjsonF0e85041DecodeMevericcoreMcws(in *jlexer.Lexer, out *WsResActionSingleErrorMsg) {
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
		case "error":
			out.Error = string(in.String())
		case "errorCode":
			out.ErrorCode = int(in.Int())
		case "status":
			out.Status = string(in.String())
		case "requestId":
			if in.IsNull() {
				in.Skip()
				out.RequestId = nil
			} else {
				if out.RequestId == nil {
					out.RequestId = new(string)
				}
				*out.RequestId = string(in.String())
			}
		case "action":
			out.Action = string(in.String())
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
func easyjsonF0e85041EncodeMevericcoreMcws(out *jwriter.Writer, in WsResActionSingleErrorMsg) {
	out.RawByte('{')
	first := true
	_ = first
	if in.Error != "" {
		if !first {
			out.RawByte(',')
		}
		first = false
		out.RawString("\"error\":")
		out.String(string(in.Error))
	}
	if in.ErrorCode != 0 {
		if !first {
			out.RawByte(',')
		}
		first = false
		out.RawString("\"errorCode\":")
		out.Int(int(in.ErrorCode))
	}
	if in.Status != "" {
		if !first {
			out.RawByte(',')
		}
		first = false
		out.RawString("\"status\":")
		out.String(string(in.Status))
	}
	if in.RequestId != nil {
		if !first {
			out.RawByte(',')
		}
		first = false
		out.RawString("\"requestId\":")
		if in.RequestId == nil {
			out.RawString("null")
		} else {
			out.String(string(*in.RequestId))
		}
	}
	if in.Action != "" {
		if !first {
			out.RawByte(',')
		}
		first = false
		out.RawString("\"action\":")
		out.String(string(in.Action))
	}
	out.RawByte('}')
}

// MarshalJSON supports json.Marshaler interface
func (v WsResActionSingleErrorMsg) MarshalJSON() ([]byte, error) {
	w := jwriter.Writer{}
	easyjsonF0e85041EncodeMevericcoreMcws(&w, v)
	return w.Buffer.BuildBytes(), w.Error
}

// MarshalEasyJSON supports easyjson.Marshaler interface
func (v WsResActionSingleErrorMsg) MarshalEasyJSON(w *jwriter.Writer) {
	easyjsonF0e85041EncodeMevericcoreMcws(w, v)
}

// UnmarshalJSON supports json.Unmarshaler interface
func (v *WsResActionSingleErrorMsg) UnmarshalJSON(data []byte) error {
	r := jlexer.Lexer{Data: data}
	easyjsonF0e85041DecodeMevericcoreMcws(&r, v)
	return r.Error()
}

// UnmarshalEasyJSON supports easyjson.Unmarshaler interface
func (v *WsResActionSingleErrorMsg) UnmarshalEasyJSON(l *jlexer.Lexer) {
	easyjsonF0e85041DecodeMevericcoreMcws(l, v)
}
func easyjsonF0e85041DecodeMevericcoreMcws1(in *jlexer.Lexer, out *WsResActionMsg) {
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
		case "status":
			out.Status = string(in.String())
		case "requestId":
			if in.IsNull() {
				in.Skip()
				out.RequestId = nil
			} else {
				if out.RequestId == nil {
					out.RequestId = new(string)
				}
				*out.RequestId = string(in.String())
			}
		case "action":
			out.Action = string(in.String())
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
func easyjsonF0e85041EncodeMevericcoreMcws1(out *jwriter.Writer, in WsResActionMsg) {
	out.RawByte('{')
	first := true
	_ = first
	if in.Status != "" {
		if !first {
			out.RawByte(',')
		}
		first = false
		out.RawString("\"status\":")
		out.String(string(in.Status))
	}
	if in.RequestId != nil {
		if !first {
			out.RawByte(',')
		}
		first = false
		out.RawString("\"requestId\":")
		if in.RequestId == nil {
			out.RawString("null")
		} else {
			out.String(string(*in.RequestId))
		}
	}
	if in.Action != "" {
		if !first {
			out.RawByte(',')
		}
		first = false
		out.RawString("\"action\":")
		out.String(string(in.Action))
	}
	out.RawByte('}')
}

// MarshalJSON supports json.Marshaler interface
func (v WsResActionMsg) MarshalJSON() ([]byte, error) {
	w := jwriter.Writer{}
	easyjsonF0e85041EncodeMevericcoreMcws1(&w, v)
	return w.Buffer.BuildBytes(), w.Error
}

// MarshalEasyJSON supports easyjson.Marshaler interface
func (v WsResActionMsg) MarshalEasyJSON(w *jwriter.Writer) {
	easyjsonF0e85041EncodeMevericcoreMcws1(w, v)
}

// UnmarshalJSON supports json.Unmarshaler interface
func (v *WsResActionMsg) UnmarshalJSON(data []byte) error {
	r := jlexer.Lexer{Data: data}
	easyjsonF0e85041DecodeMevericcoreMcws1(&r, v)
	return r.Error()
}

// UnmarshalEasyJSON supports easyjson.Unmarshaler interface
func (v *WsResActionMsg) UnmarshalEasyJSON(l *jlexer.Lexer) {
	easyjsonF0e85041DecodeMevericcoreMcws1(l, v)
}
func easyjsonF0e85041DecodeMevericcoreMcws2(in *jlexer.Lexer, out *WsResActionArrErrorMsg) {
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
		case "errors":
			if in.IsNull() {
				in.Skip()
				out.Errors = nil
			} else {
				in.Delim('[')
				if out.Errors == nil {
					if !in.IsDelim(']') {
						out.Errors = make([]string, 0, 4)
					} else {
						out.Errors = []string{}
					}
				} else {
					out.Errors = (out.Errors)[:0]
				}
				for !in.IsDelim(']') {
					var v1 string
					v1 = string(in.String())
					out.Errors = append(out.Errors, v1)
					in.WantComma()
				}
				in.Delim(']')
			}
		case "status":
			out.Status = string(in.String())
		case "requestId":
			if in.IsNull() {
				in.Skip()
				out.RequestId = nil
			} else {
				if out.RequestId == nil {
					out.RequestId = new(string)
				}
				*out.RequestId = string(in.String())
			}
		case "action":
			out.Action = string(in.String())
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
func easyjsonF0e85041EncodeMevericcoreMcws2(out *jwriter.Writer, in WsResActionArrErrorMsg) {
	out.RawByte('{')
	first := true
	_ = first
	if len(in.Errors) != 0 {
		if !first {
			out.RawByte(',')
		}
		first = false
		out.RawString("\"errors\":")
		if in.Errors == nil && (out.Flags&jwriter.NilSliceAsEmpty) == 0 {
			out.RawString("null")
		} else {
			out.RawByte('[')
			for v2, v3 := range in.Errors {
				if v2 > 0 {
					out.RawByte(',')
				}
				out.String(string(v3))
			}
			out.RawByte(']')
		}
	}
	if in.Status != "" {
		if !first {
			out.RawByte(',')
		}
		first = false
		out.RawString("\"status\":")
		out.String(string(in.Status))
	}
	if in.RequestId != nil {
		if !first {
			out.RawByte(',')
		}
		first = false
		out.RawString("\"requestId\":")
		if in.RequestId == nil {
			out.RawString("null")
		} else {
			out.String(string(*in.RequestId))
		}
	}
	if in.Action != "" {
		if !first {
			out.RawByte(',')
		}
		first = false
		out.RawString("\"action\":")
		out.String(string(in.Action))
	}
	out.RawByte('}')
}

// MarshalJSON supports json.Marshaler interface
func (v WsResActionArrErrorMsg) MarshalJSON() ([]byte, error) {
	w := jwriter.Writer{}
	easyjsonF0e85041EncodeMevericcoreMcws2(&w, v)
	return w.Buffer.BuildBytes(), w.Error
}

// MarshalEasyJSON supports easyjson.Marshaler interface
func (v WsResActionArrErrorMsg) MarshalEasyJSON(w *jwriter.Writer) {
	easyjsonF0e85041EncodeMevericcoreMcws2(w, v)
}

// UnmarshalJSON supports json.Unmarshaler interface
func (v *WsResActionArrErrorMsg) UnmarshalJSON(data []byte) error {
	r := jlexer.Lexer{Data: data}
	easyjsonF0e85041DecodeMevericcoreMcws2(&r, v)
	return r.Error()
}

// UnmarshalEasyJSON supports easyjson.Unmarshaler interface
func (v *WsResActionArrErrorMsg) UnmarshalEasyJSON(l *jlexer.Lexer) {
	easyjsonF0e85041DecodeMevericcoreMcws2(l, v)
}
func easyjsonF0e85041DecodeMevericcoreMcws3(in *jlexer.Lexer, out *WsActionMsgBaseSt) {
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
		case "requestId":
			if in.IsNull() {
				in.Skip()
				out.RequestId = nil
			} else {
				if out.RequestId == nil {
					out.RequestId = new(string)
				}
				*out.RequestId = string(in.String())
			}
		case "action":
			out.Action = string(in.String())
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
func easyjsonF0e85041EncodeMevericcoreMcws3(out *jwriter.Writer, in WsActionMsgBaseSt) {
	out.RawByte('{')
	first := true
	_ = first
	if in.RequestId != nil {
		if !first {
			out.RawByte(',')
		}
		first = false
		out.RawString("\"requestId\":")
		if in.RequestId == nil {
			out.RawString("null")
		} else {
			out.String(string(*in.RequestId))
		}
	}
	if in.Action != "" {
		if !first {
			out.RawByte(',')
		}
		first = false
		out.RawString("\"action\":")
		out.String(string(in.Action))
	}
	out.RawByte('}')
}

// MarshalJSON supports json.Marshaler interface
func (v WsActionMsgBaseSt) MarshalJSON() ([]byte, error) {
	w := jwriter.Writer{}
	easyjsonF0e85041EncodeMevericcoreMcws3(&w, v)
	return w.Buffer.BuildBytes(), w.Error
}

// MarshalEasyJSON supports easyjson.Marshaler interface
func (v WsActionMsgBaseSt) MarshalEasyJSON(w *jwriter.Writer) {
	easyjsonF0e85041EncodeMevericcoreMcws3(w, v)
}

// UnmarshalJSON supports json.Unmarshaler interface
func (v *WsActionMsgBaseSt) UnmarshalJSON(data []byte) error {
	r := jlexer.Lexer{Data: data}
	easyjsonF0e85041DecodeMevericcoreMcws3(&r, v)
	return r.Error()
}

// UnmarshalEasyJSON supports easyjson.Unmarshaler interface
func (v *WsActionMsgBaseSt) UnmarshalEasyJSON(l *jlexer.Lexer) {
	easyjsonF0e85041DecodeMevericcoreMcws3(l, v)
}