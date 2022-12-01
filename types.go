package lget

type LgateLoginInfo struct {
	LoginId  string `json:"login_id"`
	Password string `json:"password"`
}

type LoginedResp struct {
	IsSuccessful bool `json:"is_successful"`
	Code         int  `json:"code"`
	Gmt          int  `json:"gmt"`
	User         struct {
		LastName               string `json:"last_name"`
		FirstName              string `json:"first_name"`
		LoginID                string `json:"login_id"`
		UUID                   string `json:"uuid"`
		LastNameKana           string `json:"last_name_kana"`
		FirstNameKana          string `json:"first_name_kana"`
		Email                  string `json:"email"`
		IsLocal                bool   `json:"is_local"`
		IsForcedChangeName     bool   `json:"is_forced_change_name"`
		IsForcedChangeClass    bool   `json:"is_forced_change_class"`
		IsForcedChangePassword bool   `json:"is_forced_change_password"`
		IsActive               bool   `json:"is_active"`
		IsPlannedToDelete      bool   `json:"is_planned_to_delete"`
		IsAccountLocked        bool   `json:"is_account_locked"`
		CreatedAt              int    `json:"created_at"`
		UpdatedAt              int    `json:"updated_at"`
		Belongs                []struct {
			UUID        string `json:"uuid"`
			Number      string `json:"number"`
			IsCurrent   bool   `json:"is_current"`
			StartAt     string `json:"start_at"`
			EndAt       string `json:"end_at"`
			CreatedAt   int    `json:"created_at"`
			UpdatedAt   int    `json:"updated_at"`
			SchoolClass struct {
				UUID         string `json:"uuid"`
				Name         string `json:"name"`
				IsClass      bool   `json:"is_class"`
				IsInitial    bool   `json:"is_initial"`
				Number       string `json:"number"`
				CreatedAt    int    `json:"created_at"`
				UpdatedAt    int    `json:"updated_at"`
				Organization struct {
					UUID       string      `json:"uuid"`
					Name       string      `json:"name"`
					Code       string      `json:"code"`
					SchoolCode interface{} `json:"school_code"`
					Type       string      `json:"type"`
					TypeText   string      `json:"type_text"`
					CreatedAt  int         `json:"created_at"`
					UpdatedAt  int         `json:"updated_at"`
				} `json:"organization"`
				Term struct {
					UUID      string `json:"uuid"`
					Name      string `json:"name"`
					IsCurrent bool   `json:"is_current"`
					StartAt   int    `json:"start_at"`
					EndAt     int64  `json:"end_at"`
					CreatedAt int    `json:"created_at"`
					UpdatedAt int    `json:"updated_at"`
				} `json:"term"`
				Grade struct {
					Code      string `json:"code"`
					Name      string `json:"name"`
					IsLower   bool   `json:"is_lower"`
					CreatedAt int    `json:"created_at"`
					UpdatedAt int    `json:"updated_at"`
				} `json:"grade"`
			} `json:"school_class"`
			Role struct {
				UUID           string `json:"uuid"`
				PermissionCode string `json:"permission_code"`
				Name           string `json:"name"`
				IsInitial      bool   `json:"is_initial"`
				CreatedAt      int    `json:"created_at"`
				UpdatedAt      int    `json:"updated_at"`
				Permission     struct {
					Code      string `json:"code"`
					Name      string `json:"name"`
					CreatedAt int    `json:"created_at"`
					UpdatedAt int    `json:"updated_at"`
				} `json:"permission"`
				Organization struct {
					UUID       string      `json:"uuid"`
					Name       string      `json:"name"`
					Code       string      `json:"code"`
					SchoolCode interface{} `json:"school_code"`
					Type       string      `json:"type"`
					TypeText   string      `json:"type_text"`
					CreatedAt  int         `json:"created_at"`
					UpdatedAt  int         `json:"updated_at"`
				} `json:"organization"`
			} `json:"role"`
		} `json:"belongs"`
	} `json:"user"`
	Result struct {
		Message string `json:"message"`
		Tenant  struct {
			UUID           string `json:"uuid"`
			Type           string `json:"type"`
			MunicipalityID string `json:"municipality_id"`
			Code           string `json:"code"`
			Name           string `json:"name"`
			Domain         string `json:"domain"`
			PermissionURL  string `json:"permission_url"`
			Features       struct {
				Tao           bool `json:"tao"`
				Information   bool `json:"information"`
				Application   bool `json:"application"`
				Questionnaire bool `json:"questionnaire"`
				ActionLog     bool `json:"action-log"`
				Mexcbt        bool `json:"mexcbt"`
				Sync          bool `json:"sync"`
				SyncGroup     bool `json:"sync-group"`
			} `json:"features"`
		} `json:"tenant"`
	} `json:"result"`
}

type GetDataResp struct {
	IsSuccessful bool `json:"is_successful"`
	Code         int  `json:"code"`
	Gmt          int  `json:"gmt"`
	User         struct {
		LastName               string `json:"last_name"`
		FirstName              string `json:"first_name"`
		LoginID                string `json:"login_id"`
		UUID                   string `json:"uuid"`
		LastNameKana           string `json:"last_name_kana"`
		FirstNameKana          string `json:"first_name_kana"`
		Email                  string `json:"email"`
		IsLocal                bool   `json:"is_local"`
		IsForcedChangeName     bool   `json:"is_forced_change_name"`
		IsForcedChangeClass    bool   `json:"is_forced_change_class"`
		IsForcedChangePassword bool   `json:"is_forced_change_password"`
		IsActive               bool   `json:"is_active"`
		IsPlannedToDelete      bool   `json:"is_planned_to_delete"`
		IsAccountLocked        bool   `json:"is_account_locked"`
		CreatedAt              int    `json:"created_at"`
		UpdatedAt              int    `json:"updated_at"`
		Belongs                []struct {
			UUID        string `json:"uuid"`
			Number      string `json:"number"`
			IsCurrent   bool   `json:"is_current"`
			StartAt     string `json:"start_at"`
			EndAt       string `json:"end_at"`
			CreatedAt   int    `json:"created_at"`
			UpdatedAt   int    `json:"updated_at"`
			SchoolClass struct {
				UUID         string `json:"uuid"`
				Name         string `json:"name"`
				IsClass      bool   `json:"is_class"`
				IsInitial    bool   `json:"is_initial"`
				Number       string `json:"number"`
				CreatedAt    int    `json:"created_at"`
				UpdatedAt    int    `json:"updated_at"`
				Organization struct {
					UUID       string      `json:"uuid"`
					Name       string      `json:"name"`
					Code       string      `json:"code"`
					SchoolCode interface{} `json:"school_code"`
					Type       string      `json:"type"`
					TypeText   string      `json:"type_text"`
					CreatedAt  int         `json:"created_at"`
					UpdatedAt  int         `json:"updated_at"`
				} `json:"organization"`
				Term struct {
					UUID      string `json:"uuid"`
					Name      string `json:"name"`
					IsCurrent bool   `json:"is_current"`
					StartAt   int    `json:"start_at"`
					EndAt     int64  `json:"end_at"`
					CreatedAt int    `json:"created_at"`
					UpdatedAt int    `json:"updated_at"`
				} `json:"term"`
				Grade struct {
					Code      string `json:"code"`
					Name      string `json:"name"`
					IsLower   bool   `json:"is_lower"`
					CreatedAt int    `json:"created_at"`
					UpdatedAt int    `json:"updated_at"`
				} `json:"grade"`
			} `json:"school_class"`
			Role struct {
				UUID           string `json:"uuid"`
				PermissionCode string `json:"permission_code"`
				Name           string `json:"name"`
				IsInitial      bool   `json:"is_initial"`
				CreatedAt      int    `json:"created_at"`
				UpdatedAt      int    `json:"updated_at"`
				Permission     struct {
					Code      string `json:"code"`
					Name      string `json:"name"`
					CreatedAt int    `json:"created_at"`
					UpdatedAt int    `json:"updated_at"`
				} `json:"permission"`
				Organization struct {
					UUID       string      `json:"uuid"`
					Name       string      `json:"name"`
					Code       string      `json:"code"`
					SchoolCode interface{} `json:"school_code"`
					Type       string      `json:"type"`
					TypeText   string      `json:"type_text"`
					CreatedAt  int         `json:"created_at"`
					UpdatedAt  int         `json:"updated_at"`
				} `json:"organization"`
			} `json:"role"`
		} `json:"belongs"`
	} `json:"user"`
	Result struct {
		UUID      string `json:"uuid"`
		Type      string `json:"type"`
		TypeName  string `json:"type_name"`
		IsSuccess bool   `json:"is_success"`
		Message   string `json:"message"`
		Result    struct {
			FileUUID string `json:"file_uuid"`
		} `json:"result"`
		CreatedAt int `json:"created_at"`
		UpdatedAt int `json:"updated_at"`
		DoneAt    int `json:"done_at"`
	} `json:"result"`
}

type LgateOpener struct {
	SessId     string
	ResultUuid string
	JobUuid    string
}
