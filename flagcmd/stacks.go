package flagcmd

import (
	"fmt"

	"github.com/wednesday-solutions/picky/internal/utils"
)

type StackFlag struct {
	// all available stacks status
	A bool
	// all existing stacks status
	E bool
	// web stacks status
	W bool
	// mobile stacks status
	M bool
	// backend stacks status
	B bool
}

func (f StackFlag) FlagStacks() string {

	var stacks []string
	var userOutput string
	if f.A {
		f.E = false
		userOutput = utils.AllStacksString()
		userOutput = fmt.Sprintf("\nAll stacks: %s", userOutput)
		if f.B || f.M || f.W {
			userOutput = ""
		}
		if f.B {
			userOutput = fmt.Sprintf("%s%s", utils.AllBackendStacksString(), userOutput)
		}
		if f.M {
			userOutput = fmt.Sprintf("%s%s", utils.AllMobileStacksString(), userOutput)
		}
		if f.W {
			userOutput = fmt.Sprintf("%s%s", utils.AllWebStacksString(), userOutput)
		}
	}
	if f.E {
		stacks = utils.GetExistingStacks()
		userOutput = utils.ConvertStacksIntoString(stacks)
		userOutput = fmt.Sprintf("\nAll existing stacks: %s", userOutput)
		if f.B || f.M || f.W {
			userOutput = ""
		}
		if f.B {
			var backendStacks []string
			backendStack, status := "", false
			for _, stack := range stacks {
				backendStack, status = utils.IsBackendStack(stack)
				if status {
					backendStacks = append(backendStacks, backendStack)
				}
			}
			if len(backendStacks) == 0 {
				userOutput = fmt.Sprintf("\n\tNo backend stacks exist.%s", userOutput)
			} else {
				userOutput = fmt.Sprintf("%s%s", utils.ConvertStacksIntoString(backendStacks), userOutput)
			}
			userOutput = fmt.Sprintf("\nAll existing backend stacks: %s", userOutput)
		}
		if f.M {
			var mobileStacks []string
			mobileStack, status := "", false
			for _, stack := range stacks {
				mobileStack, status = utils.IsMobileStack(stack)
				if status {
					mobileStacks = append(mobileStacks, mobileStack)
				}
			}
			if len(mobileStacks) == 0 {
				userOutput = fmt.Sprintf("\n\tNo mobile stacks exist.%s", userOutput)
			} else {
				userOutput = fmt.Sprintf("%s%s", utils.ConvertStacksIntoString(mobileStacks), userOutput)
			}
			userOutput = fmt.Sprintf("\nAll existing mobile stacks: %s", userOutput)
		}
		if f.W {
			var webStacks []string
			webStack, status := "", false
			for _, stack := range stacks {
				webStack, status = utils.IsWebStack(stack)
				if status {
					webStacks = append(webStacks, webStack)
				}
			}
			if len(webStacks) == 0 {
				userOutput = fmt.Sprintf("\n\tNo web stacks exist.%s", userOutput)
			} else {
				userOutput = fmt.Sprintf("%s%s", utils.ConvertStacksIntoString(webStacks), userOutput)
			}
			userOutput = fmt.Sprintf("\nAll existing web stacks: %s", userOutput)
		}
	}
	return userOutput
}
