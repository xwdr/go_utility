package utils

import (
	"errors"
	"testing"

	gomock "github.com/golang/mock/gomock"
	. "github.com/smartystreets/goconvey/convey"
)

func TestFindUser(t *testing.T) {
	Convey("FindUser", t, func() {
		Convey("参数错误", func() {
			// 生成伪数据
			//id := new(int64)
			//_ = faker.FakeData(id)
			// 初始化一个mock controller
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			api := NewMockApi(ctrl)
			Convey("id错误", func() {
				// 数据打桩
				api.EXPECT().GetNameById(int64(100)).Return("", errors.New("id错误")).AnyTimes()
				// 测试程序
				service := &Service{api: api}
				name, err := service.FindUser(100)
				// 断言
				So(err, ShouldBeNil)
				So(name, ShouldEqual, "")
			})
			Convey("id正常", func() {
				// 数据打桩
				api.EXPECT().GetNameById(int64(100)).Return("张山", nil).AnyTimes()
				// 测试程序
				service := &Service{api: api}
				name, err := service.FindUser(100)
				// 断言
				So(err, ShouldBeNil)
				So(name, ShouldNotBeEmpty)
			})
		})
		Convey("数据未找到", func() {})
	})
}
