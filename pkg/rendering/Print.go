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
package rendering

import "github.com/pterm/pterm"

func BlankLine() {
	pterm.Println()
}

func Error(message string, args ...any) {
	pterm.Error.Printfln(message, args...)
}

func Header(message string, args ...any) {
	pterm.DefaultHeader.WithFullWidth().Printfln(
		message,
		args...,
	)
}

func Paragraph(message string, args ...any) {
	pterm.DefaultParagraph.Printfln(message, args...)
}

func Print(message string, args ...any) {
	pterm.DefaultBasicText.Printfln(message, args...)
}

func Spinner(message string) *pterm.SpinnerPrinter {
	result, _ := pterm.DefaultSpinner.Start(message)
	return result
}

func Success(message string, args ...any) {
	pterm.Success.Printfln(message, args...)
}

func Warning(message string, args ...any) {
	pterm.Warning.Printfln(message, args...)
}
