package passkeeper

import (
	"context"
	"errors"
	"fmt"
	"os"
	"path"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/vindosVP/go-pass/internal/models"
	"github.com/vindosVP/go-pass/internal/services/passkeeper/mocks"
	"github.com/vindosVP/go-pass/internal/storage"
	"github.com/vindosVP/go-pass/pkg/logger/sl"
)

func TestKeeper_SavePassword(t *testing.T) {

	unexpected := errors.New("unexpected error")
	type w struct {
		id  int
		err error
	}
	type m struct {
		needed bool
		id     int
		err    error
	}

	tests := []struct {
		name string
		e    *models.Entity
		m    m
		w    w
	}{
		{
			name: "ok",
			e: &models.Entity{
				ID:       1,
				OwnerID:  1,
				Type:     models.TypePassword,
				Login:    "login",
				Password: "password",
				Metadata: "metadata",
			},
			m: m{
				needed: true,
				id:     1,
				err:    nil,
			},
			w: w{
				id:  1,
				err: nil,
			},
		},
		{
			name: "unexpected",
			e: &models.Entity{
				ID:       1,
				OwnerID:  1,
				Type:     models.TypePassword,
				Login:    "login",
				Password: "password",
				Metadata: "metadata",
			},
			m: m{
				needed: true,
				id:     0,
				err:    unexpected,
			},
			w: w{
				id:  0,
				err: unexpected,
			},
		},
		{
			name: "unknown",
			e: &models.Entity{
				ID:       1,
				OwnerID:  1,
				Type:     "UNKNOWN",
				Login:    "login",
				Password: "password",
				Metadata: "metadata",
			},
			m: m{
				needed: false,
			},
			w: w{
				id:  0,
				err: ErrUnknownEntity,
			},
		},
	}

	sl.SetupLogger("test")

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()
			ps := mocks.NewPasswordStorage(t)
			if tt.m.needed {
				ps.On("AddPassword", mock.Anything, tt.e.ToPassword()).Return(tt.m.id, tt.m.err)
			}
			k := New(ps, nil, nil, nil, "")
			id, err := k.Save(ctx, tt.e)
			if tt.w.err == nil {
				assert.Equal(t, tt.w.id, id)
			}
			assert.ErrorIs(t, tt.w.err, err)
		})
	}
}

func TestKeeper_Save(t *testing.T) {

	unexpected := errors.New("unexpected error")
	type w struct {
		id  int
		err error
	}
	type psM struct {
		needed bool
		id     int
		err    error
	}
	type csM struct {
		needed bool
		id     int
		err    error
	}
	type tsM struct {
		needed bool
		id     int
		err    error
	}
	type testCases []struct {
		name string
		e    models.Entity
		psM  psM
		csM  csM
		tsM  tsM
		w    w
	}

	template := testCases{
		{
			name: "ok",
			e: models.Entity{
				ID:         1,
				OwnerID:    1,
				Login:      "login",
				Password:   "password",
				CardNumber: "1234 1234 1234 1234",
				CardOwner:  "CARD OWNER",
				CardCVC:    "123",
				CardExp:    "06/22",
				Text:       "text",
				Metadata:   "metadata",
			},
			csM: csM{
				id:  1,
				err: nil,
			},
			psM: psM{
				id:  1,
				err: nil,
			},
			tsM: tsM{
				id:  1,
				err: nil,
			},
			w: w{
				id:  1,
				err: nil,
			},
		},
		{
			name: "unexpected error",
			e: models.Entity{
				ID:         1,
				OwnerID:    1,
				Login:      "login",
				Password:   "password",
				CardNumber: "1234 1234 1234 1234",
				CardOwner:  "CARD OWNER",
				CardCVC:    "123",
				CardExp:    "06/22",
				Text:       "text",
				Metadata:   "metadata",
			},
			csM: csM{
				id:  0,
				err: unexpected,
			},
			psM: psM{
				id:  0,
				err: unexpected,
			},
			tsM: tsM{
				id:  0,
				err: unexpected,
			},
			w: w{
				id:  0,
				err: unexpected,
			},
		},
	}

	cases := make(testCases, 0, len(template)*len(eTypes()))

	for _, et := range eTypes() {
		tc := make(testCases, len(template), len(template))
		copy(tc, template)
		for i, tt := range tc {
			tc[i].psM.needed = et == models.TypePassword
			tc[i].csM.needed = et == models.TypeCard
			tc[i].tsM.needed = et == models.TypeText
			tc[i].name = fmt.Sprintf("%s_%s", tt.name, et)
			tc[i].e.Type = et
			if et == models.TypeFile {
				tc[i].w.id = 0
				tc[i].w.err = ErrUnableToSaveFile
			}
		}
		cases = append(cases, tc...)
	}

	sl.SetupLogger("test")

	t.Parallel()
	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()
			ps := mocks.NewPasswordStorage(t)
			cs := mocks.NewCardStorage(t)
			ts := mocks.NewTextStorage(t)
			if tt.psM.needed {
				ps.On("AddPassword", mock.Anything, tt.e.ToPassword()).Return(tt.psM.id, tt.psM.err)
			}
			if tt.csM.needed {
				cs.On("AddCard", mock.Anything, tt.e.ToCard()).Return(tt.csM.id, tt.csM.err)
			}
			if tt.tsM.needed {
				ts.On("AddText", mock.Anything, tt.e.ToText()).Return(tt.tsM.id, tt.tsM.err)
			}
			k := New(ps, cs, ts, nil, "")
			id, err := k.Save(ctx, &tt.e)
			if tt.w.err == nil {
				assert.Equal(t, tt.w.id, id)
			}
			assert.ErrorIs(t, tt.w.err, err)
		})
	}
}

func TestKeeper_Update(t *testing.T) {

	unexpected := errors.New("unexpected error")
	type w struct {
		err error
	}
	type psM struct {
		needed bool
		err    error
	}
	type csM struct {
		needed bool
		err    error
	}
	type tsM struct {
		needed bool
		err    error
	}
	type testCases []struct {
		name string
		e    models.Entity
		psM  psM
		csM  csM
		tsM  tsM
		w    w
	}

	template := testCases{
		{
			name: "ok",
			e: models.Entity{
				ID:         1,
				OwnerID:    1,
				Login:      "login",
				Password:   "password",
				CardNumber: "1234 1234 1234 1234",
				CardOwner:  "CARD OWNER",
				CardCVC:    "123",
				CardExp:    "06/22",
				Text:       "text",
				Metadata:   "metadata",
			},
			csM: csM{
				err: nil,
			},
			psM: psM{
				err: nil,
			},
			tsM: tsM{
				err: nil,
			},
			w: w{
				err: nil,
			},
		},
		{
			name: "unexpected error",
			e: models.Entity{
				ID:         1,
				OwnerID:    1,
				Login:      "login",
				Password:   "password",
				CardNumber: "1234 1234 1234 1234",
				CardOwner:  "CARD OWNER",
				CardCVC:    "123",
				CardExp:    "06/22",
				Text:       "text",
				Metadata:   "metadata",
			},
			csM: csM{
				err: unexpected,
			},
			psM: psM{
				err: unexpected,
			},
			tsM: tsM{
				err: unexpected,
			},
			w: w{
				err: unexpected,
			},
		},
	}

	cases := make(testCases, 0, len(template)*len(eTypes()))

	for _, et := range eTypes() {
		tc := make(testCases, len(template), len(template))
		copy(tc, template)
		for i, tt := range tc {
			tc[i].psM.needed = et == models.TypePassword
			tc[i].csM.needed = et == models.TypeCard
			tc[i].tsM.needed = et == models.TypeText
			tc[i].name = fmt.Sprintf("%s_%s", tt.name, et)
			tc[i].e.Type = et
			if et == models.TypeFile {
				tc[i].w.err = ErrUnableToUpdateFile
			}
		}
		cases = append(cases, tc...)
	}

	sl.SetupLogger("test")

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()
			ps := mocks.NewPasswordStorage(t)
			cs := mocks.NewCardStorage(t)
			ts := mocks.NewTextStorage(t)
			if tt.psM.needed {
				ps.On("UpdatePassword", mock.Anything, tt.e.ToPassword()).Return(tt.psM.err)
			}
			if tt.csM.needed {
				cs.On("UpdateCard", mock.Anything, tt.e.ToCard()).Return(tt.csM.err)
			}
			if tt.tsM.needed {
				ts.On("UpdateText", mock.Anything, tt.e.ToText()).Return(tt.tsM.err)
			}
			k := New(ps, cs, ts, nil, "")
			err := k.Update(ctx, &tt.e)
			assert.ErrorIs(t, tt.w.err, err)
		})
	}
}

func TestKeeper_Delete(t *testing.T) {

	unexpected := errors.New("unexpected error")
	type w struct {
		err error
	}
	type psM struct {
		needed bool
		err    error
	}
	type csM struct {
		needed bool
		err    error
	}
	type tsM struct {
		needed bool
		err    error
	}

	type testCases []struct {
		name    string
		id      int
		ownerId int
		et      models.EntityType
		psM     psM
		csM     csM
		tsM     tsM
		w       w
	}

	template := testCases{
		{
			name:    "ok",
			id:      1,
			ownerId: 1,
			csM: csM{
				err: nil,
			},
			psM: psM{
				err: nil,
			},
			tsM: tsM{
				err: nil,
			},

			w: w{
				err: nil,
			},
		},
		{
			name:    "unexpected error",
			id:      1,
			ownerId: 1,
			csM: csM{
				err: unexpected,
			},
			psM: psM{
				err: unexpected,
			},
			tsM: tsM{
				err: unexpected,
			},
			w: w{
				err: unexpected,
			},
		},
	}

	cases := make(testCases, 0, len(template)*len(eTypes()))

	for _, et := range eTypes() {
		if et == models.TypeFile {
			continue
		}
		tc := make(testCases, len(template), len(template))
		copy(tc, template)
		for i, tt := range tc {
			tc[i].psM.needed = et == models.TypePassword
			tc[i].csM.needed = et == models.TypeCard
			tc[i].tsM.needed = et == models.TypeText
			tc[i].et = et
			tc[i].name = fmt.Sprintf("%s_%s", tt.name, et)
		}
		cases = append(cases, tc...)
	}

	sl.SetupLogger("test")

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()
			ps := mocks.NewPasswordStorage(t)
			cs := mocks.NewCardStorage(t)
			ts := mocks.NewTextStorage(t)
			if tt.psM.needed {
				ps.On("DeletePassword", mock.Anything, tt.id, tt.ownerId).Return(tt.psM.err)
			}
			if tt.csM.needed {
				cs.On("DeleteCard", mock.Anything, tt.id, tt.ownerId).Return(tt.csM.err)
			}
			if tt.tsM.needed {
				ts.On("DeleteText", mock.Anything, tt.id, tt.ownerId).Return(tt.tsM.err)
			}
			k := New(ps, cs, ts, nil, "")
			err := k.Delete(ctx, tt.id, tt.ownerId, tt.et)
			assert.ErrorIs(t, tt.w.err, err)
		})
	}
}

func TestKeeper_List(t *testing.T) {
	sl.SetupLogger("test")
	unexpected := errors.New("unexpected error")
	ownerID := 1
	ctx := context.Background()

	ps := mocks.NewPasswordStorage(t)
	cs := mocks.NewCardStorage(t)
	ts := mocks.NewTextStorage(t)
	fs := mocks.NewFileStorage(t)

	k := New(ps, cs, ts, fs, "")

	ps.On("GetPasswords", mock.Anything, ownerID).Return(nil, unexpected).Once()
	_, err := k.List(ctx, ownerID)
	assert.ErrorIs(t, err, unexpected)

	ps.On("GetPasswords", mock.Anything, ownerID).Return([]*models.Password{}, nil).Once()
	cs.On("GetCards", mock.Anything, ownerID).Return(nil, unexpected).Once()
	_, err = k.List(ctx, ownerID)
	assert.ErrorIs(t, err, unexpected)

	ps.On("GetPasswords", mock.Anything, ownerID).Return([]*models.Password{}, nil).Once()
	cs.On("GetCards", mock.Anything, ownerID).Return([]*models.Card{}, nil).Once()
	ts.On("GetTexts", mock.Anything, ownerID).Return(nil, unexpected).Once()
	_, err = k.List(ctx, ownerID)
	assert.ErrorIs(t, err, unexpected)

	ps.On("GetPasswords", mock.Anything, ownerID).Return([]*models.Password{}, nil).Once()
	cs.On("GetCards", mock.Anything, ownerID).Return([]*models.Card{}, nil).Once()
	ts.On("GetTexts", mock.Anything, ownerID).Return([]*models.Text{}, nil).Once()
	fs.On("GetFiles", mock.Anything, ownerID).Return(nil, unexpected).Once()
	_, err = k.List(ctx, ownerID)
	assert.ErrorIs(t, err, unexpected)

	ps.On("GetPasswords", mock.Anything, ownerID).Return([]*models.Password{}, nil).Once()
	cs.On("GetCards", mock.Anything, ownerID).Return([]*models.Card{}, nil).Once()
	ts.On("GetTexts", mock.Anything, ownerID).Return([]*models.Text{}, nil).Once()
	fs.On("GetFiles", mock.Anything, ownerID).Return([]*models.File{}, nil).Once()
	_, err = k.List(ctx, ownerID)
	assert.NoError(t, err)
}

func eTypes() []models.EntityType {
	return []models.EntityType{models.TypePassword, models.TypeCard, models.TypeText, models.TypeFile}
}

func TestKeeper_DeleteFile(t *testing.T) {

	sl.SetupLogger("test")
	fType := models.TypeFile
	unexpected := errors.New("unexpected error")
	fileName := "file.txt"
	fileLocation := "./files"
	id := 1
	ownerId := 1

	ctx := context.Background()
	fs := mocks.NewFileStorage(t)

	err := os.Mkdir("./files", 0777)
	require.NoError(t, err)

	file, err := os.Create(path.Join(fileLocation, fmt.Sprintf("%d_%s", id, fileName)))
	require.NoError(t, err)
	_, err = file.Write([]byte("text"))
	require.NoError(t, err)
	err = file.Close()
	require.NoError(t, err)

	k := New(nil, nil, nil, fs, fileLocation)
	fs.On("GetFile", mock.Anything, id, ownerId).Return(&models.File{FileName: fileName, ID: id}, nil).Once()
	fs.On("DeleteFile", mock.Anything, id, ownerId).Return(nil).Once()
	err = k.Delete(ctx, id, ownerId, fType)
	assert.NoError(t, err)

	fs.On("GetFile", mock.Anything, id, ownerId).Return(&models.File{FileName: fileName, ID: id}, nil).Once()
	err = k.Delete(ctx, id, ownerId, fType)
	assert.Error(t, err)

	fs.On("GetFile", mock.Anything, id, ownerId).Return(nil, storage.ErrFileNotExist).Once()
	err = k.Delete(ctx, id, ownerId, fType)
	assert.ErrorIs(t, err, storage.ErrFileNotExist)

	fs.On("GetFile", mock.Anything, id, ownerId).Return(nil, unexpected).Once()
	err = k.Delete(ctx, id, ownerId, fType)
	assert.ErrorIs(t, err, unexpected)

	err = os.Remove("./files")
	require.NoError(t, err)
}
