package utils

import (
	"fmt"
	"strconv"
)

// export function createInvoiceCode(ticketId: string) {
// const ticketIdStr = `${ticketId}`;
// const ticketIdStrLength = ticketIdStr.length;
// let iSum = 0;
// let iDigit = 0;
// for (let i = 0; i < ticketIdStrLength; i++) {
// iDigit = Number(`${ticketIdStr[i]}`);
// if (i % 2 === 0) {
// iSum += iDigit * 1;
// } else {
// iSum += iDigit * 3;
// }
// }
// const iCheckSum = (10 - (iSum % 10)) % 10;
// if (ticketIdStrLength > 8) {
// return ticketIdStr.slice(-8) + `${iCheckSum}`;
// } else {
// return ticketIdStr + `${iCheckSum}`;
// }
// }

// CreateInvoiceCode generates a specialized invoice code based on a ticket ID.
// It calculates a checksum digit using a weighted sum (1x for even indices, 3x for odd)
// and appends it to the last 8 characters of the original ID.
func CreateInvoiceCode(ticketID string) string {
	sum := 0

	// Calculate weighted sum for the checksum
	for i, char := range ticketID {
		// Convert rune to integer; ignore non-numeric characters to prevent panic
		digit, err := strconv.Atoi(string(char))
		if err != nil {
			continue
		}

		// Apply weights based on position parity
		if i%2 == 0 {
			sum += digit
		} else {
			sum += digit * 3
		}
	}

	// Calculate the check digit (complement to the next multiple of 10)
	checkSum := (10 - (sum % 10)) % 10

	// Extract the suffix (last 8 characters) or use the full string if shorter
	baseID := ticketID
	if len(ticketID) > 8 {
		baseID = ticketID[len(ticketID)-8:]
	}

	// Format result as {BaseID}{Checksum}
	return fmt.Sprintf("%s%d", baseID, checkSum)
}
