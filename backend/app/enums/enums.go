package enums

const (
	// 관리자 타입
	AdminTypeManager = 1
	//AdminTypeRegular = 2

	//휴가 타입
	VacationTypeNormal    = 1
	VacationTypePromotion = 2

	//휴가 발생 타입
	VacationGenerateTypeAnnualNormal           = 1
	VacationGenerateTypeAnnualThisYearPreGiven = 2
	VacationGenerateTypeAnnualOneYearPreGiven  = 3

	VacationGenerateTypePreAccountingNormal           = 4
	VacationGenerateTypePreAccountingThisYearPreGiven = 5
	VacationGenerateTypePreAccountingOneYearPreGiven  = 6

	VacationGenerateTypeProAccountingNormal           = 7
	VacationGenerateTypeProAccountingThisYearPreGiven = 8
	VacationGenerateTypeProAccountingOneYearPreGiven  = 9

	//휴가 촉진 상태
	VacationPromotionStateNone           = 1
	VacationPromotionStateFirstNoti      = 2
	VacationPromotionStateFirstComplete  = 3
	VacationPromotionStateSecondNeed     = 4
	VacationPromotionStateSecondNoti     = 5
	VacationPromotionStateSecondComplete = 6

	//휴가 처리 상태
	VacationProcessStateApplied       = 1
	VacationProcessStateFirstApproved = 2
	VacationProcessStateFinalApproved = 3
	VacationProcessStateRejected      = 4

	//휴가 취소 상태
	VacationCancelStateDefault   = 1
	VacationCancelStateRequested = 2
	VacationCancelStateCompleted = 3

	//알림 타입
	NotificationTypeNormal                        = 1
	NotificationTypeVacationApplied               = 2
	NotificationTypeVacationFirstPromotion        = 3
	NotificationTypeVacationFirstPromotionAccept  = 4
	NotificationTypeVacationSecondPromotion       = 5
	NotificationTypeVacationSecondPromotionAccept = 6
	NotificationTypeVacationDenyWork              = 7
	NotificationTypeVacationDenyWorkAccept        = 8
)
