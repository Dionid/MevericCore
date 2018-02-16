// Code generated by easyjson for marshaling/unmarshaling. DO NOT EDIT.

package mccommon

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

func easyjson5371cc0DecodeMevericcoreMccommon(in *jlexer.Lexer, out *UsersListModel) {
	isTopLevel := in.IsStart()
	if in.IsNull() {
		in.Skip()
		*out = nil
	} else {
		in.Delim('[')
		if *out == nil {
			if !in.IsDelim(']') {
				*out = make(UsersListModel, 0, 1)
			} else {
				*out = UsersListModel{}
			}
		} else {
			*out = (*out)[:0]
		}
		for !in.IsDelim(']') {
			var v1 UserModel
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
func easyjson5371cc0EncodeMevericcoreMccommon(out *jwriter.Writer, in UsersListModel) {
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
func (v UsersListModel) MarshalJSON() ([]byte, error) {
	w := jwriter.Writer{}
	easyjson5371cc0EncodeMevericcoreMccommon(&w, v)
	return w.Buffer.BuildBytes(), w.Error
}

// MarshalEasyJSON supports easyjson.Marshaler interface
func (v UsersListModel) MarshalEasyJSON(w *jwriter.Writer) {
	easyjson5371cc0EncodeMevericcoreMccommon(w, v)
}

// UnmarshalJSON supports json.Unmarshaler interface
func (v *UsersListModel) UnmarshalJSON(data []byte) error {
	r := jlexer.Lexer{Data: data}
	easyjson5371cc0DecodeMevericcoreMccommon(&r, v)
	return r.Error()
}

// UnmarshalEasyJSON supports easyjson.Unmarshaler interface
func (v *UsersListModel) UnmarshalEasyJSON(l *jlexer.Lexer) {
	easyjson5371cc0DecodeMevericcoreMccommon(l, v)
}
func easyjson5371cc0DecodeMevericcoreMccommon1(in *jlexer.Lexer, out *UserModel) {
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
		case "login":
			out.Login = string(in.String())
		case "email":
			if in.IsNull() {
				in.Skip()
				out.Email = nil
			} else {
				if out.Email == nil {
					out.Email = new(string)
				}
				*out.Email = string(in.String())
			}
		case "password":
			out.Password = string(in.String())
		case "isAdmin":
			out.IsAdmin = bool(in.Bool())
		case "phone":
			if in.IsNull() {
				in.Skip()
				out.Phone = nil
			} else {
				if out.Phone == nil {
					out.Phone = new(string)
				}
				*out.Phone = string(in.String())
			}
		case "companyId":
			if in.IsNull() {
				in.Skip()
				out.CompanyId = nil
			} else {
				if out.CompanyId == nil {
					out.CompanyId = new(bson.ObjectId)
				}
				if data := in.Raw(); in.Ok() {
					in.AddError((*out.CompanyId).UnmarshalJSON(data))
				}
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
func easyjson5371cc0EncodeMevericcoreMccommon1(out *jwriter.Writer, in UserModel) {
	out.RawByte('{')
	first := true
	_ = first
	if in.Login != "" {
		if !first {
			out.RawByte(',')
		}
		first = false
		out.RawString("\"login\":")
		out.String(string(in.Login))
	}
	if in.Email != nil {
		if !first {
			out.RawByte(',')
		}
		first = false
		out.RawString("\"email\":")
		if in.Email == nil {
			out.RawString("null")
		} else {
			out.String(string(*in.Email))
		}
	}
	if in.Password != "" {
		if !first {
			out.RawByte(',')
		}
		first = false
		out.RawString("\"password\":")
		out.String(string(in.Password))
	}
	if in.IsAdmin {
		if !first {
			out.RawByte(',')
		}
		first = false
		out.RawString("\"isAdmin\":")
		out.Bool(bool(in.IsAdmin))
	}
	if in.Phone != nil {
		if !first {
			out.RawByte(',')
		}
		first = false
		out.RawString("\"phone\":")
		if in.Phone == nil {
			out.RawString("null")
		} else {
			out.String(string(*in.Phone))
		}
	}
	if in.CompanyId != nil {
		if !first {
			out.RawByte(',')
		}
		first = false
		out.RawString("\"companyId\":")
		if in.CompanyId == nil {
			out.RawString("null")
		} else {
			out.Raw((*in.CompanyId).MarshalJSON())
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
func (v UserModel) MarshalJSON() ([]byte, error) {
	w := jwriter.Writer{}
	easyjson5371cc0EncodeMevericcoreMccommon1(&w, v)
	return w.Buffer.BuildBytes(), w.Error
}

// MarshalEasyJSON supports easyjson.Marshaler interface
func (v UserModel) MarshalEasyJSON(w *jwriter.Writer) {
	easyjson5371cc0EncodeMevericcoreMccommon1(w, v)
}

// UnmarshalJSON supports json.Unmarshaler interface
func (v *UserModel) UnmarshalJSON(data []byte) error {
	r := jlexer.Lexer{Data: data}
	easyjson5371cc0DecodeMevericcoreMccommon1(&r, v)
	return r.Error()
}

// UnmarshalEasyJSON supports easyjson.Unmarshaler interface
func (v *UserModel) UnmarshalEasyJSON(l *jlexer.Lexer) {
	easyjson5371cc0DecodeMevericcoreMccommon1(l, v)
}
func easyjson5371cc0DecodeMevericcoreMccommon2(in *jlexer.Lexer, out *CompanyModel) {
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
		case "employeesIds":
			if in.IsNull() {
				in.Skip()
				out.EmployeesIds = nil
			} else {
				in.Delim('[')
				if out.EmployeesIds == nil {
					if !in.IsDelim(']') {
						out.EmployeesIds = make([]bson.ObjectId, 0, 4)
					} else {
						out.EmployeesIds = []bson.ObjectId{}
					}
				} else {
					out.EmployeesIds = (out.EmployeesIds)[:0]
				}
				for !in.IsDelim(']') {
					var v4 bson.ObjectId
					if data := in.Raw(); in.Ok() {
						in.AddError((v4).UnmarshalJSON(data))
					}
					out.EmployeesIds = append(out.EmployeesIds, v4)
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
func easyjson5371cc0EncodeMevericcoreMccommon2(out *jwriter.Writer, in CompanyModel) {
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
	if len(in.EmployeesIds) != 0 {
		if !first {
			out.RawByte(',')
		}
		first = false
		out.RawString("\"employeesIds\":")
		if in.EmployeesIds == nil && (out.Flags&jwriter.NilSliceAsEmpty) == 0 {
			out.RawString("null")
		} else {
			out.RawByte('[')
			for v5, v6 := range in.EmployeesIds {
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
func (v CompanyModel) MarshalJSON() ([]byte, error) {
	w := jwriter.Writer{}
	easyjson5371cc0EncodeMevericcoreMccommon2(&w, v)
	return w.Buffer.BuildBytes(), w.Error
}

// MarshalEasyJSON supports easyjson.Marshaler interface
func (v CompanyModel) MarshalEasyJSON(w *jwriter.Writer) {
	easyjson5371cc0EncodeMevericcoreMccommon2(w, v)
}

// UnmarshalJSON supports json.Unmarshaler interface
func (v *CompanyModel) UnmarshalJSON(data []byte) error {
	r := jlexer.Lexer{Data: data}
	easyjson5371cc0DecodeMevericcoreMccommon2(&r, v)
	return r.Error()
}

// UnmarshalEasyJSON supports easyjson.Unmarshaler interface
func (v *CompanyModel) UnmarshalEasyJSON(l *jlexer.Lexer) {
	easyjson5371cc0DecodeMevericcoreMccommon2(l, v)
}
func easyjson5371cc0DecodeMevericcoreMccommon3(in *jlexer.Lexer, out *CompanyListModel) {
	isTopLevel := in.IsStart()
	if in.IsNull() {
		in.Skip()
		*out = nil
	} else {
		in.Delim('[')
		if *out == nil {
			if !in.IsDelim(']') {
				*out = make(CompanyListModel, 0, 1)
			} else {
				*out = CompanyListModel{}
			}
		} else {
			*out = (*out)[:0]
		}
		for !in.IsDelim(']') {
			var v7 CompanyModel
			(v7).UnmarshalEasyJSON(in)
			*out = append(*out, v7)
			in.WantComma()
		}
		in.Delim(']')
	}
	if isTopLevel {
		in.Consumed()
	}
}
func easyjson5371cc0EncodeMevericcoreMccommon3(out *jwriter.Writer, in CompanyListModel) {
	if in == nil && (out.Flags&jwriter.NilSliceAsEmpty) == 0 {
		out.RawString("null")
	} else {
		out.RawByte('[')
		for v8, v9 := range in {
			if v8 > 0 {
				out.RawByte(',')
			}
			(v9).MarshalEasyJSON(out)
		}
		out.RawByte(']')
	}
}

// MarshalJSON supports json.Marshaler interface
func (v CompanyListModel) MarshalJSON() ([]byte, error) {
	w := jwriter.Writer{}
	easyjson5371cc0EncodeMevericcoreMccommon3(&w, v)
	return w.Buffer.BuildBytes(), w.Error
}

// MarshalEasyJSON supports easyjson.Marshaler interface
func (v CompanyListModel) MarshalEasyJSON(w *jwriter.Writer) {
	easyjson5371cc0EncodeMevericcoreMccommon3(w, v)
}

// UnmarshalJSON supports json.Unmarshaler interface
func (v *CompanyListModel) UnmarshalJSON(data []byte) error {
	r := jlexer.Lexer{Data: data}
	easyjson5371cc0DecodeMevericcoreMccommon3(&r, v)
	return r.Error()
}

// UnmarshalEasyJSON supports easyjson.Unmarshaler interface
func (v *CompanyListModel) UnmarshalEasyJSON(l *jlexer.Lexer) {
	easyjson5371cc0DecodeMevericcoreMccommon3(l, v)
}
