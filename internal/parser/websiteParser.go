package parser

import (
	"context"
	"fmt"
	"net/http"
	"net/mail"
	"net/url"
	"regexp"
	"sort"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
)

var (
	emailRegex = regexp.MustCompile(`(?i)([a-z0-9._%+\-]+@[a-z0-9.\-]+\.[a-z]{2,})`)
	phoneRegex = regexp.MustCompile(`(?i)(\+?\d[\d\s()\-]{8,}\d)`)
	binRegex   = regexp.MustCompile(`\b\d{12}\b`)
)

type ParseResult struct {
	Website           string
	BIN               string
	Emails            []string
	Phones            []string
	Addresses         []string
	Linkedin          string
	Facebook          string
	ProcurementMethod string
	ProcurementEmail  string
	ProcurementPhone  string
	HRName            string
	HREmail           string
	HRPhone           string
	ESGName           string
	ESGEmail          string
	ESGPhone          string
	ESGReportURL      string
	HasESGDept        bool
	Source            string
}

type WebsiteParser struct {
	client *http.Client
}

func NewWebsiteParser() *WebsiteParser {
	return &WebsiteParser{
		client: &http.Client{Timeout: 15 * time.Second},
	}
}

func (p *WebsiteParser) Parse(ctx context.Context, rawURL string) (*ParseResult, error) {
	normalized, err := normalizeURL(rawURL)
	if err != nil {
		return nil, err
	}

	pages := []string{
		normalized,
		normalized + "/contact",
		normalized + "/contacts",
		normalized + "/about",
		normalized + "/procurement",
		normalized + "/careers",
		normalized + "/esg",
		normalized + "/sustainability",
		normalized + "/контакты",
		normalized + "/закупки",
		normalized + "/карьера",
		normalized + "/устойчивое-развитие",
	}

	result := &ParseResult{
		Website: normalized,
		Source:  normalized,
	}

	var fullText strings.Builder

	for _, pageURL := range pages {
		doc, text, err := p.fetchDoc(ctx, pageURL)
		if err != nil {
			continue
		}

		fullText.WriteString("\n")
		fullText.WriteString(text)

		html, _ := doc.Html()

		for _, e := range extractEmails(html + "\n" + text) {
			result.Emails = append(result.Emails, e)
		}

		if result.BIN == "" {
			result.BIN = extractBIN(text)
		}

		for _, ph := range extractPhones(doc, html+"\n"+text, result.BIN) {
			result.Phones = append(result.Phones, ph)
		}

		if result.Linkedin == "" {
			result.Linkedin = findSocial(html, "linkedin.com")
		}
		if result.Facebook == "" {
			result.Facebook = findSocial(html, "facebook.com")
		}
		if result.ESGReportURL == "" {
			result.ESGReportURL = findReportLink(html)
		}
	}

	text := strings.ToLower(fullText.String())
	result.Emails = uniqueValidEmails(result.Emails)
	result.Phones = rankPhones(uniqueStrings(result.Phones))

	if strings.Contains(text, "goszakup") || strings.Contains(text, "госзакуп") {
		result.ProcurementMethod = "Goszakup"
	} else if strings.Contains(text, "samruk") {
		result.ProcurementMethod = "Samruk"
	} else if strings.Contains(text, "nadloc") || strings.Contains(text, "nadloq") {
		result.ProcurementMethod = "Nadloc"
	}

	if strings.Contains(text, "esg") || strings.Contains(text, "устойчивое развитие") || result.ESGReportURL != "" {
		result.HasESGDept = true
	}

	if len(result.Emails) > 0 {
		result.ProcurementEmail = pickByKeyword(result.Emails, []string{"purchase", "procurement", "tender", "zakup"})
		result.HREmail = pickByKeyword(result.Emails, []string{"hr", "career", "job", "resume"})
		result.ESGEmail = pickByKeyword(result.Emails, []string{"esg", "sustain", "csr"})
	}

	if len(result.Phones) > 0 {
		result.ProcurementPhone = result.Phones[0]
		result.HRPhone = result.Phones[0]
		result.ESGPhone = result.Phones[0]
	}

	return result, nil
}

func (p *WebsiteParser) fetchDoc(ctx context.Context, pageURL string) (*goquery.Document, string, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, pageURL, nil)
	if err != nil {
		return nil, "", err
	}
	req.Header.Set("User-Agent", "CompanyESGBot/1.0")

	resp, err := p.client.Do(req)
	if err != nil {
		return nil, "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		return nil, "", fmt.Errorf("status %d", resp.StatusCode)
	}

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return nil, "", err
	}

	return doc, cleanWhitespace(doc.Text()), nil
}

func normalizeURL(raw string) (string, error) {
	raw = strings.TrimSpace(raw)
	if !strings.HasPrefix(raw, "http://") && !strings.HasPrefix(raw, "https://") {
		raw = "https://" + raw
	}

	u, err := url.Parse(raw)
	if err != nil || u.Host == "" {
		return "", fmt.Errorf("invalid url")
	}

	return strings.TrimRight(u.String(), "/"), nil
}

func extractEmails(s string) []string {
	return emailRegex.FindAllString(s, -1)
}

func extractPhones(doc *goquery.Document, s string, detectedBIN string) []string {
	var phones []string

	// 1. Надёжный источник: tel: ссылки
	doc.Find("a[href^='tel:'], a[href^='TEL:']").Each(func(i int, sel *goquery.Selection) {
		href, exists := sel.Attr("href")
		if !exists {
			return
		}

		phone := strings.TrimSpace(href)
		phone = strings.TrimPrefix(phone, "tel:")
		phone = strings.TrimPrefix(phone, "TEL:")

		if normalized, ok := normalizePhone(phone, detectedBIN, true); ok {
			phones = append(phones, normalized)
		}
	})

	// 2. Текстовый поиск — только строгая фильтрация
	matches := phoneRegex.FindAllString(s, -1)
	for _, match := range matches {
		if normalized, ok := normalizePhone(match, detectedBIN, false); ok {
			phones = append(phones, normalized)
		}
	}

	return uniqueStrings(phones)
}

func normalizePhone(phone string, detectedBIN string, fromTel bool) (string, bool) {
	original := strings.TrimSpace(phone)
	if original == "" {
		return "", false
	}

	var b strings.Builder
	hasPlus := false

	for i, r := range original {
		if r == '+' && i == 0 {
			hasPlus = true
			continue
		}
		if r >= '0' && r <= '9' {
			b.WriteRune(r)
		}
	}

	digits := b.String()

	if len(digits) < 10 || len(digits) > 15 {
		return "", false
	}

	// Совпал с BIN
	if detectedBIN != "" && digits == detectedBIN {
		return "", false
	}

	// Просто 12 цифр подряд без форматирования — скорее BIN/ID
	if len(digits) == 12 &&
		!fromTel &&
		!hasPlus &&
		!strings.ContainsAny(original, " ()-") {
		return "", false
	}

	// Для обычного текста отбрасываем "голые" длинные числа без форматирования
	if !fromTel && !hasPlus && !strings.ContainsAny(original, " ()-") {
		return "", false
	}

	// Отсекаем мусор типа 0000000000
	allSame := true
	for i := 1; i < len(digits); i++ {
		if digits[i] != digits[0] {
			allSame = false
			break
		}
	}
	if allSame {
		return "", false
	}

	if hasPlus {
		return "+" + digits, true
	}

	return digits, true
}

func rankPhones(in []string) []string {
	type scoredPhone struct {
		Phone string
		Score int
	}

	var scored []scoredPhone
	for _, ph := range in {
		score := 0

		switch {
		case strings.HasPrefix(ph, "+7"):
			score += 30
		case strings.HasPrefix(ph, "8"):
			score += 25
		case strings.HasPrefix(ph, "+1"):
			score += 20
		}

		if len(ph) == 11 || len(ph) == 12 {
			score += 10
		}

		scored = append(scored, scoredPhone{Phone: ph, Score: score})
	}

	sort.SliceStable(scored, func(i, j int) bool {
		return scored[i].Score > scored[j].Score
	})

	var out []string
	for _, item := range scored {
		out = append(out, item.Phone)
	}
	return out
}

func extractBIN(s string) string {
	return binRegex.FindString(s)
}

func findSocial(html, domain string) string {
	re := regexp.MustCompile(`https?://[^\s\"'<>]+`)
	for _, link := range re.FindAllString(html, -1) {
		if strings.Contains(link, domain) {
			return strings.TrimRight(strings.Trim(link, `"'<>),.;`), "/")
		}
	}
	return ""
}

func findReportLink(html string) string {
	re := regexp.MustCompile(`https?://[^\s\"'<>]+`)
	for _, link := range re.FindAllString(html, -1) {
		lower := strings.ToLower(link)
		if strings.Contains(lower, "esg") ||
			strings.Contains(lower, "sustain") ||
			strings.Contains(lower, "report") ||
			strings.Contains(lower, "csr") {
			return strings.TrimRight(strings.Trim(link, `"'<>),.;`), "/")
		}
	}
	return ""
}

func pickByKeyword(items, keys []string) string {
	for _, item := range items {
		lower := strings.ToLower(item)
		for _, key := range keys {
			if strings.Contains(lower, key) {
				return item
			}
		}
	}

	if len(items) > 0 {
		return items[0]
	}

	return ""
}

func uniqueValidEmails(in []string) []string {
	seen := map[string]bool{}
	var out []string

	for _, e := range in {
		e = strings.Trim(strings.ToLower(e), " .;,<>()[]{}\"'")
		if e == "" || seen[e] {
			continue
		}
		if _, err := mail.ParseAddress(e); err == nil {
			seen[e] = true
			out = append(out, e)
		}
	}

	sort.Strings(out)
	return out
}

func uniqueStrings(in []string) []string {
	seen := map[string]bool{}
	var out []string

	for _, item := range in {
		item = cleanWhitespace(item)
		if item == "" || seen[item] {
			continue
		}
		seen[item] = true
		out = append(out, item)
	}

	return out
}

func cleanWhitespace(s string) string {
	re := regexp.MustCompile(`\s+`)
	return strings.TrimSpace(re.ReplaceAllString(s, " "))
}
