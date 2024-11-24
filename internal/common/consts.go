package common

// FollowCategory 追番分类
type FollowCategory int

// 追番分类
const (
	FollowCategoryClassic FollowCategory = iota
	FollowCategoryHighQuality
	FollowCategoryNew
	FollowCategoryToiletPaper
	FollowCategoryMasterpiece
	FollowCategoryUnknown = 999
)

// String 返回追番分类的字符串表示
func (fc FollowCategory) String() string {
	return [...]string{"旧时代的残党", "我们仍未知道那天所看见的番剧的名字", "新番妙妙屋", "厕纸", "神！"}[fc]
}

// Description 返回追番分类的描述
func (fc FollowCategory) Description() string {
	return [...]string{"经典老番或长篇番", "听说高质量的番", "看新番导视比较感兴趣的番", "专门用来找乐子杀时间的番", "无需多言"}[fc]
}

// IsValid 检查追番分类是否合法
func (fc FollowCategory) IsValid() bool {
	switch fc {
	case FollowCategoryClassic, FollowCategoryHighQuality, FollowCategoryNew, FollowCategoryToiletPaper, FollowCategoryMasterpiece:
		return true
	default:
		return false
	}
}

// AllFollowCategories 返回所有追番分类
func AllFollowCategories() []FollowCategory {
	return []FollowCategory{
		FollowCategoryClassic,
		FollowCategoryHighQuality,
		FollowCategoryNew,
		FollowCategoryToiletPaper,
		FollowCategoryMasterpiece,
	}
}

// FollowStatus 追番状态
type FollowStatus int

const (
	FollowStatusWantToWatch FollowStatus = iota
	FollowStatusWatching
	FollowStatusWatched
	FollowStatusUnknown = 999
)

func (fs FollowStatus) String() string {
	return [...]string{"want_to_watch", "watching", "watched"}[fs]
}

func (fs FollowStatus) IsValid() bool {
	switch fs {
	case FollowStatusWantToWatch, FollowStatusWatching, FollowStatusWatched:
		return true
	default:
		return false
	}
}
