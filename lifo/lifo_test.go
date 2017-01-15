// Copyright © 2016 Abcum Ltd
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package lifo

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestInit(t *testing.T) {

	c, e := New(0)

	Convey("Initial size should be greater than 0", t, func() {
		So(c, ShouldBeNil)
		So(e, ShouldNotBeNil)
	})

}

func TestMain(t *testing.T) {

	c, _ := New(10)

	Convey("------------------------------", t, nil)

	Convey("Size should be stored", t, func() {
		So(c.size, ShouldEqual, 10)
		So(c.bytes, ShouldEqual, 0)
	})

	Convey("------------------------------", t, nil)

	Convey("Can insert 1st item", t, func() {
		c.Put("test1", []byte("Hello"))
		So(c.bytes, ShouldEqual, 5)
		So(c.Has("test1"), ShouldBeTrue)
	})

	Convey("Can insert 2nd item", t, func() {
		c.Put("test2", []byte("Testing"))
		So(c.bytes, ShouldEqual, 7)
		So(c.Has("test1"), ShouldBeFalse)
		So(c.Has("test2"), ShouldBeTrue)
	})

	Convey("Can insert 3rd item", t, func() {
		c.Put("test3", []byte("This is really long."))
		So(c.bytes, ShouldEqual, 7)
		So(c.Has("test1"), ShouldBeFalse)
		So(c.Has("test2"), ShouldBeTrue)
		So(c.Has("test3"), ShouldBeFalse)
	})

	Convey("Can insert 4th item", t, func() {
		c.Put("test4", []byte("Bye"))
		So(c.bytes, ShouldEqual, 10)
		So(c.Has("test1"), ShouldBeFalse)
		So(c.Has("test2"), ShouldBeTrue)
		So(c.Has("test3"), ShouldBeFalse)
		So(c.Has("test4"), ShouldBeTrue)
	})

	Convey("Can insert duplicate", t, func() {
		c.Put("test4", []byte("Bye"))
		So(c.bytes, ShouldEqual, 10)
		So(c.Has("test1"), ShouldBeFalse)
		So(c.Has("test2"), ShouldBeTrue)
		So(c.Has("test3"), ShouldBeFalse)
		So(c.Has("test4"), ShouldBeTrue)
		So(c.Get("test4"), ShouldResemble, []byte("Bye"))
	})

	Convey("------------------------------", t, nil)

	Convey("Can get 1st item", t, func() {
		i := c.Get("test1")
		So(i, ShouldBeNil)
	})

	Convey("Can get 2nd item", t, func() {
		i := c.Get("test2")
		So(i, ShouldResemble, []byte("Testing"))
	})

	Convey("------------------------------", t, nil)

	Convey("Can del 1st item", t, func() {
		i := c.Del("test1")
		So(i, ShouldBeNil)
	})

	Convey("Can del 2nd item", t, func() {
		i := c.Del("test2")
		So(i, ShouldResemble, []byte("Testing"))
	})

	Convey("------------------------------", t, nil)

	Convey("Can clear and empty cache", t, func() {
		c.Clr()
		So(c.bytes, ShouldEqual, 0)
		So(len(c.items), ShouldEqual, 0)
	})

}
