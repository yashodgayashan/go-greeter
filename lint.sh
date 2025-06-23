#!/bin/bash

# Go Greeter Linting Script
# This script runs various linting tools and provides a summary

set -e

echo "🔍 Running Go Linting Tools..."
echo "================================"

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Function to print colored output
print_status() {
    local status=$1
    local message=$2
    case $status in
        "PASS")
            echo -e "${GREEN}✅ PASS${NC}: $message"
            ;;
        "WARN")
            echo -e "${YELLOW}⚠️  WARN${NC}: $message"
            ;;
        "FAIL")
            echo -e "${RED}❌ FAIL${NC}: $message"
            ;;
        "INFO")
            echo -e "${BLUE}ℹ️  INFO${NC}: $message"
            ;;
    esac
}

# Track overall status
overall_status=0

echo -e "\n${BLUE}1. Code Formatting (gofmt)${NC}"
echo "----------------------------"
if go fmt ./... | grep -q .; then
    print_status "WARN" "Code formatting issues found and fixed"
else
    print_status "PASS" "Code is properly formatted"
fi

echo -e "\n${BLUE}2. Go Vet Analysis${NC}"
echo "-------------------"
if go vet ./... 2>&1; then
    print_status "PASS" "No issues found by go vet"
else
    print_status "FAIL" "Issues found by go vet"
    overall_status=1
fi

echo -e "\n${BLUE}3. golangci-lint Analysis${NC}"
echo "-------------------------"
if golangci-lint run --out-format=colored-line-number 2>/dev/null; then
    print_status "PASS" "No issues found by golangci-lint"
else
    echo "Issues found by golangci-lint:"
    echo ""
    
    # Run golangci-lint and filter out test file issues
    echo "📋 Main code issues:"
    golangci-lint run --out-format=line-number 2>&1 | grep -v "_test.go" || echo "  No main code issues found!"
    
    echo ""
    echo "📋 Test file issues (informational):"
    golangci-lint run --out-format=line-number 2>&1 | grep "_test.go" | head -5 || echo "  No test file issues found!"
    
    test_issues_count=$(golangci-lint run --out-format=line-number 2>&1 | grep "_test.go" | wc -l | tr -d ' ')
    if [ "$test_issues_count" -gt 5 ]; then
        echo "  ... and $((test_issues_count - 5)) more test file issues"
    fi
    
    main_issues_count=$(golangci-lint run --out-format=line-number 2>&1 | grep -v "_test.go" | wc -l | tr -d ' ')
    if [ "$main_issues_count" -gt 0 ]; then
        print_status "FAIL" "$main_issues_count issues in main code"
        overall_status=1
    else
        print_status "WARN" "Only test file issues found (not critical)"
    fi
fi

echo -e "\n${BLUE}4. Module Tidiness${NC}"
echo "------------------"
if go mod tidy && git diff --quiet go.mod go.sum 2>/dev/null; then
    print_status "PASS" "Go modules are tidy"
else
    print_status "WARN" "Go modules may need tidying"
fi

echo -e "\n${BLUE}5. Security Check (gosec)${NC}"
echo "------------------------"
if command -v gosec >/dev/null 2>&1; then
    if gosec -quiet ./... 2>/dev/null; then
        print_status "PASS" "No security issues found"
    else
        print_status "WARN" "Potential security issues found (run 'gosec ./...' for details)"
    fi
else
    print_status "INFO" "gosec not installed (optional security scanner)"
fi

echo ""
echo "================================"
if [ $overall_status -eq 0 ]; then
    print_status "PASS" "All critical linting checks passed! 🎉"
    echo ""
    echo "💡 Note: Test file linting issues are typically acceptable"
    echo "   as test code has different quality standards."
else
    print_status "FAIL" "Some critical issues found that should be fixed"
    echo ""
    echo "🔧 Run 'make lint-fix' to auto-fix some issues"
    echo "🔧 Run 'golangci-lint run' for detailed output"
fi

echo ""
echo "📊 Quick Stats:"
echo "  - Go files: $(find . -name "*.go" -not -path "./vendor/*" | wc -l | tr -d ' ')"
echo "  - Lines of code: $(wc -l *.go | tail -1 | awk '{print $1}')"
echo "  - Test coverage: $(go test -cover 2>/dev/null | grep coverage | awk '{print $2}' || echo 'unknown')"

exit $overall_status 