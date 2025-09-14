package iasiutils

import "fmt"

// Recipe handles prompt building and related logic for LLMs.
type Recipe struct{}

// BuildLLMPrompt creates a prompt for the LLM using the problem statement and solution
func (r *Recipe) BuildLLMPrompt(statement, solution string) string {
	if len(statement) == 0 {
		statement = "(Problem statement could not be fetched)"
	}
	if len(solution) == 0 {
		solution = "(Solution code could not be fetched)"
	}
	return fmt.Sprintf(`You are an expert competitive programming assistant. Given the following problem statement and its solution, generate:
- 3 helpful hints for a student (in Romanian, do not give away the full solution)
- a detailed editorial (in Romanian, explaining the solution and key ideas)

Problem statement:
%s

Solution:
%s

Return a JSON object with two fields: "hints" (an array of 3 strings) and "editorial" (a string).`, statement, solution)
}
