package iasiutils

import "fmt"

// Recipe handles prompt building and related logic for LLMs.
type Recipe struct {
	SystemPrompt string
}

// BuildLLMPrompt creates a prompt for the LLM using the problem statement and solution
func (r *Recipe) BuildLLMPrompt(statement, solution string) (prompt string, systemPrompt string) {
	if len(statement) == 0 {
		statement = "(Problem statement could not be fetched)"
	}
	if len(solution) == 0 {
		solution = "(Solution code could not be fetched)"
	}
	prompt = fmt.Sprintf(`You are an expert competitive programming assistant. Given the following problem statement and its solution, generate:
	- some helpful hints for a student (in English, do not give away the full solution). Make them so that the student can understand the key ideas and approach to solve the problem on their own. They should gradually lead the student to the solution, without revealing it directly. Provide around 3 hints. Adjust the number based on the complexity and difficulty of the problem. Keep the hints concise and to the point, rather short, don't give away too much.
	- a detailed editorial (in English, explaining the solution and key ideas). Don't include snippets of code from the solution. Do an editoril like on Codeforces. Please structure it in markdown format with the necessary sections. Use the solution only as guidance, do not use any namings from the solution at all. You can use names from the task itself.

Problem statement:
%s

Solution (this is not the official solution):
%s

Return a JSON object with two fields: "hints" (an array of strings) and "editorial" (a string).`, statement, solution)
	systemPrompt = r.SystemPrompt
	return
}
