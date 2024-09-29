/*
Copyright © 2024 Adam Presley

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in
all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
THE SOFTWARE.
*/
package services

import (
	"github.com/adampresley/adamgokit/slices"
	"github.com/adampresley/pusher/pkg/contextinfo"
	"github.com/adampresley/pusher/pkg/sshutils"
)

type ServiceItem struct {
	ServiceName string
	Description string
	Step        *sshutils.Step
	Collector   func(info *contextinfo.ContextInfo)
}

var (
	ServiceList = []ServiceItem{
		PostgresService,
	}
)

func GetServiceListNames() []string {
	result := slices.Map(ServiceList, func(input ServiceItem, index int) string {
		return input.ServiceName
	})

	return result
}

func GetServiceByName(name string) ServiceItem {
	return slices.Find(ServiceList, func(item ServiceItem) bool {
		if item.ServiceName == name {
			return true
		}

		return false
	})
}
