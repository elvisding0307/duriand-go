package api

import (
	"duriand/internal/dao"
	"duriand/internal/model"
	"errors"
	"time"
)

func InsertAccountService(uid uint64, website string, account string, password string) error {
	// Get current timestamp
	now := time.Now().Unix()

	// Create account record
	accountRecord := model.Account{
		Uid:        uid,
		Website:    website,
		Account:    account,
		Password:   password,
		UpdateTime: now,
	}

	// Start transaction
	tx := dao.DB_INSTANCE.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// Insert account record
	if err := tx.Create(&accountRecord).Error; err != nil {
		tx.Rollback()
		return err
	}

	// Update timestamp using the user's uid
	timestamp := model.Timestamp{
		Uid:              uid,
		LatestUpdateTime: now,
	}

	if err := tx.Model(&model.Timestamp{}).Where("uid = ?", uid).
		Assign(map[string]interface{}{"latest_update_time": now}).
		FirstOrCreate(&timestamp).Error; err != nil {
		tx.Rollback()
		return err
	}

	// Commit transaction
	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		return err
	}
	return nil
}

type PullMode string
type QueriedAccount struct {
	Rid      uint64 `json:"rid"`
	Website  string `json:"website"`
	Account  string `json:"account"`
	Password string `json:"password"`
}

func QueryAccountService(uid uint64, updateTime int64) (PullMode, int64, []QueriedAccount, error) {
	const (
		PULL_ALL     PullMode = "PULL_ALL"
		PULL_UPDATED PullMode = "PULL_UPDATED"
		PULL_NOTHING PullMode = "PULL_NOTHING"
	)

	// Get timestamp info
	var timestamp model.Timestamp
	if err := dao.DB_INSTANCE.Where("uid = ?", uid).First(&timestamp).Error; err != nil {
		return "", 0, nil, err
	}

	pullMode := PULL_NOTHING
	latestUpdateTime := timestamp.LatestUpdateTime
	queriedAccounts := []QueriedAccount{}

	// Check update time
	// 不需要更新
	if updateTime >= timestamp.LatestUpdateTime {
		pullMode = PULL_NOTHING
		return pullMode, latestUpdateTime, queriedAccounts, nil
	}
	var accounts []model.Account
	if updateTime < timestamp.LatestDeleteTime || updateTime == 0 {
		pullMode = PULL_ALL
		// Return all accounts
		if err := dao.DB_INSTANCE.Where("uid = ?", uid).Find(&accounts).Error; err != nil {
			return "", 0, nil, err
		}
	} else {
		pullMode = PULL_UPDATED
		// Return only updated accounts
		if err := dao.DB_INSTANCE.Where("uid = ? AND update_time > ?", uid, updateTime).Find(&accounts).Error; err != nil {
			return "", 0, nil, err
		}
	}
	for _, acc := range accounts {
		queriedAccounts = append(queriedAccounts, QueriedAccount{
			Rid:      acc.Rid,
			Website:  acc.Website,
			Account:  acc.Account,
			Password: acc.Password,
		})
	}
	return pullMode, latestUpdateTime, queriedAccounts, nil
}

func UpdateAccountService(uid uint64, rid uint64, website, account, password string) error {
	tx := dao.DB_INSTANCE.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// 验证记录是否存在且uid匹配
	var accountModel model.Account
	if err := tx.Where("rid = ?", rid).First(&accountModel).Error; err != nil {
		tx.Rollback()
		return err
	}

	if accountModel.Uid != uid {
		tx.Rollback()
		return errors.New("Permission denied")
	}

	// 更新记录
	currentTime := time.Now().Unix()
	updateData := map[string]interface{}{
		"update_time": currentTime,
	}
	if website != "" {
		updateData["website"] = website
	}
	if account != "" {
		updateData["account"] = account
	}
	if password != "" {
		updateData["password"] = password
	}

	if err := tx.Model(&model.Account{}).Where("rid = ?", rid).Updates(updateData).Error; err != nil {
		tx.Rollback()
		return err
	}

	// 更新timestamp表
	if err := tx.Model(&model.Timestamp{}).Where("uid = ?", uid).Update("latest_update_time", currentTime).Error; err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit().Error
}

func DeleteAccountService(uid uint64, rid uint64) error {
	tx := dao.DB_INSTANCE.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// 验证记录是否存在且uid匹配
	var accountModel model.Account
	if err := tx.Where("rid = ?", rid).First(&accountModel).Error; err != nil {
		tx.Rollback()
		return err
	}

	if accountModel.Uid != uid {
		tx.Rollback()
		return errors.New("Permission denied")
	}

	now := time.Now().Unix()

	// 删除账户记录
	if err := tx.Where("rid = ? AND uid = ?", rid, uid).Delete(&model.Account{}).Error; err != nil {
		tx.Rollback()
		return err
	}

	// 更新timestamp
	timestamp := model.Timestamp{
		Uid:              uid,
		LatestUpdateTime: now,
		LatestDeleteTime: now,
	}
	if err := tx.Model(&model.Timestamp{}).Where("uid = ?", uid).
		Updates(timestamp).Error; err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit().Error
}
