package middleware

import "strings"

func containsIgnoreCase(slice []string, target string) bool {
    for _, s := range slice {
        if strings.EqualFold(s, target) {
            return true
        }
    }
    return false
}