package rpcService

import (
	"fmt"
	"hmsList/crypt"
)

func (this *RegisterRequest) Encrypt(ckey []byte) error {
	if !this.IsEncrypt || len(ckey) == 0 {
		return nil
	}
	var err error
	this.Id, err = crypt.EncryptClientData(this.Id, ckey)
	if err != nil {
		return fmt.Errorf("RegisterRequest Encrypt Id error:%w", err)
	}
	this.Pid, err = crypt.EncryptClientData(this.Pid, ckey)
	if err != nil {
		return fmt.Errorf("RegisterRequest Encrypt Pid error:%w", err)
	}
	this.Version, err = crypt.EncryptClientData(this.Version, ckey)
	if err != nil {
		return fmt.Errorf("RegisterRequest Encrypt Version error:%w", err)
	}

	for idx := range this.Attrs {
		this.Attrs[idx].Key, err = crypt.EncryptClientData(this.Attrs[idx].Key, ckey)
		if err != nil {
			return fmt.Errorf("RegisterRequest Encrypt Attrs[%d].key error:%w", idx, err)
		}
		this.Attrs[idx].Value, err = crypt.EncryptClientData(this.Attrs[idx].Value, ckey)
		if err != nil {
			return fmt.Errorf("RegisterRequest Encrypt Attrs[%d].value error:%w", idx, err)
		}
	}
	return nil
}

func (this *AuthRequest) Encrypt(ckey []byte) error {
	if !this.IsEncrypt || len(ckey) == 0 {
		return nil
	}
	var err error
	this.Id, err = crypt.EncryptClientData(this.Id, ckey)
	if err != nil {
		return fmt.Errorf("AuthRequest Encrypt Id error:%w", err)
	}
	this.Pid, err = crypt.EncryptClientData(this.Pid, ckey)
	if err != nil {
		return fmt.Errorf("AuthRequest Encrypt Pid error:%w", err)
	}
	return nil
}

func (this *ReportRequest) Encrypt(ckey []byte) error {
	if !this.IsEncrypt || len(ckey) == 0 {
		return nil
	}
	var err error
	this.Id, err = crypt.EncryptClientData(this.Id, ckey)
	if err != nil {
		return fmt.Errorf("ReportRequest Encrypt Id error:%w", err)
	}
	this.Pid, err = crypt.EncryptClientData(this.Pid, ckey)
	if err != nil {
		return fmt.Errorf("ReportRequest Encrypt Pid error:%w", err)
	}
	this.Action, err = crypt.EncryptClientData(this.Action, ckey)
	if err != nil {
		return fmt.Errorf("ReportRequest Encrypt Action error:%w", err)
	}

	for idx := range this.Attrs {
		this.Attrs[idx].Key, err = crypt.EncryptClientData(this.Attrs[idx].Key, ckey)
		if err != nil {
			return fmt.Errorf("ReportRequest Encrypt Attrs[%d].key error:%w", idx, err)
		}
		this.Attrs[idx].Value, err = crypt.EncryptClientData(this.Attrs[idx].Value, ckey)
		if err != nil {
			return fmt.Errorf("ReportRequest Encrypt Attrs[%d].value error:%w", idx, err)
		}
	}
	return nil
}

func (this *AuthResponse) Decrypt(ckey []byte) error {
	if !this.IsEncrypt || len(ckey) == 0 {
		return nil
	}
	var err error
	this.Id, err = crypt.DecryptServerData(this.Id, ckey)
	if err != nil {
		return fmt.Errorf("AuthResponse Encrypt Id error:%w", err)
	}

	for idx := range this.Attrs {
		this.Attrs[idx].Key, err = crypt.DecryptServerData(this.Attrs[idx].Key, ckey)
		if err != nil {
			return fmt.Errorf("AuthResponse Decrypt Attrs[%d].key error:%w", idx, err)
		}
		this.Attrs[idx].Value, err = crypt.DecryptServerData(this.Attrs[idx].Value, ckey)
		if err != nil {
			return fmt.Errorf("AuthResponse Decrypt Attrs[%d].value error:%w", idx, err)
		}
	}
	return nil
}
