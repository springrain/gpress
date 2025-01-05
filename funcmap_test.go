// Copyright (c) 2023 gpress Authors.
//
// This file is part of gpress.
//
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU Affero General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// This program is distributed in the hope that it will be useful
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU Affero General Public License for more details.
//
// You should have received a copy of the GNU Affero General Public License
// along with this program.  If not, see <http://www.gnu.org/licenses/>.

package main

import (
	"reflect"
	"testing"
)

func TestDataSliceCategory2Tree(t *testing.T) {
	type args struct {
		categories []*Category
	}
	tests := []struct {
		name string
		args args
		want []*Category
	}{
		{
			name: "无节点",
			args: args{
				categories: nil,
			},
			want: []*Category{},
		},
		{
			name: "两级节点",
			args: args{
				categories: []*Category{
					{Id: "1", Name: "Category 1", Pid: ""},
					{Id: "2", Name: "Category 2", Pid: "1"},
					{Id: "3", Name: "Category 3", Pid: ""},
					{Id: "4", Name: "Category 4", Pid: "3"},
				},
			},
			want: []*Category{
				{
					Id:   "1",
					Name: "Category 1",
					Leaf: []*Category{{Id: "2", Name: "Category 2", Pid: "1"}},
				},
				{
					Id:   "3",
					Name: "Category 3",
					Leaf: []*Category{{Id: "4", Name: "Category 4", Pid: "3"}},
				},
			},
		},
		{
			name: "多级节点",
			args: args{
				categories: []*Category{
					{Id: "1", Name: "Category 1", Pid: ""},
					{Id: "2", Name: "Category 2", Pid: "1"},
					{Id: "3", Name: "Category 3", Pid: "1"},
					{Id: "4", Name: "Category 4", Pid: "2"},
					{Id: "5", Name: "Category 5", Pid: "2"},
					{Id: "6", Name: "Category 6", Pid: "3"},
				},
			},
			want: []*Category{
				{
					Id:   "1",
					Name: "Category 1",
					Leaf: []*Category{
						{
							Id:   "2",
							Pid:  "1",
							Name: "Category 2",
							Leaf: []*Category{
								{
									Id:   "4",
									Pid:  "2",
									Name: "Category 4",
								},
								{
									Id:   "5",
									Pid:  "2",
									Name: "Category 5",
								},
							},
						},
						{
							Id:   "3",
							Pid:  "1",
							Name: "Category 3",
							Leaf: []*Category{
								{
									Id:   "6",
									Pid:  "3",
									Name: "Category 6",
								},
							},
						},
					},
				},
			},
		},
		{
			name: "多颗树",
			args: args{
				categories: []*Category{
					{Id: "1", Name: "Category 1", Pid: ""},
					{Id: "2", Name: "Category 2", Pid: "1"},
					{Id: "3", Name: "Category 3", Pid: ""},
					{Id: "4", Name: "Category 4", Pid: "3"},
					{Id: "5", Name: "Category 5", Pid: ""},
					{Id: "6", Name: "Category 6", Pid: "5"},
				},
			},
			want: []*Category{
				{
					Id:   "1",
					Name: "Category 1",
					Leaf: []*Category{
						{
							Id:   "2",
							Pid:  "1",
							Name: "Category 2",
						},
					},
				},
				{
					Id:   "3",
					Name: "Category 3",
					Leaf: []*Category{
						{
							Id:   "4",
							Pid:  "3",
							Name: "Category 4",
						},
					},
				},
				{
					Id:   "5",
					Name: "Category 5",
					Leaf: []*Category{
						{
							Id:   "6",
							Pid:  "5",
							Name: "Category 6",
						},
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := sliceCategory2Tree(tt.args.categories); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("sliceCategory2Tree() = %v, want %v", got, tt.want)
			}
		})
	}
}
