package parse_file

import "fmt"

const extractSystemPrompt = "You extract data from documents. Your entire response MUST be a single raw JSON object and nothing else: start with { and end with }. Do NOT wrap it in markdown code fences (no ```json, no ```). Do NOT add prose, explanations, backticks, or any text before or after the JSON."

// rpaPrompt is the v1 single-field (purchase_price) prompt, migrated as-is.
// US2 extends it to the 5-field rpaFields set + arbitratePrompt.
func rpaPrompt(text string) string {
	return fmt.Sprintf(`Extract the purchase price from the following document text.
"purchase_price" MUST be a JSON number (float), not a string. Strip thousands separators (commas), currency symbols, and whitespace. Examples: "1,750,000" -> 1750000, "$1,750,000.00" -> 1750000.
Respond as JSON: {"purchase_price": <number>}
If not found, respond: {"purchase_price": 0}

Document text:
%s`, text)
}

func escrowPrompt() string {
	return `Extract these fields from the attached first page of a wire/escrow instruction document.
- escrow_company: the escrow/title company name in standard Title Case (e.g. "Mission Escrow, Inc.").
- officer_name: the escrow officer's full name only (e.g. "Renee Lo").
- officer_email: the escrow officer's email address, if printed on the page.
- officer_phone: the escrow officer's phone number as DIGITS ONLY — strip spaces, parentheses, dashes, dots, plus signs (e.g. "(408) 513-6005" becomes "4085136005").
- escrow_number: the value next to the "Escrow No." label (e.g. "202656-RL").
Only extract what is actually printed on the page. If a field is absent, unreadable, or uncertain, use null — never guess.
Respond as JSON: {"escrow_company": <string|null>, "officer_name": <string|null>, "officer_email": <string|null>, "officer_phone": <string|null>, "escrow_number": <string|null>}`
}
