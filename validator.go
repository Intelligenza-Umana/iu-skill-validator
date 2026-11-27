package main

import (
	"bytes"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

var (
	markdownLinkPattern = regexp.MustCompile(`\]\(\s*(?:<([^>]+)>|([^\s)]+))`)
	backtickPattern     = regexp.MustCompile("`([^`]+)`")
)

func ValidateSkillDir(skillDir string) (bool, []string) {
	findings := make([]string, 0)

	info, err := os.Stat(skillDir)
	if err != nil || !info.IsDir() {
		if err != nil {
			findings = append(findings, "not a valid directory: "+skillDir+": "+err.Error())
		} else {
			findings = append(findings, "not a valid directory: "+skillDir)
		}
		return false, findings
	}

	skillPath := filepath.Join(skillDir, "SKILL.md")
	data, err := os.ReadFile(skillPath)
	if err != nil {
		findings = append(findings, "could not read SKILL.md: "+err.Error())
		return false, findings
	}
	if len(bytes.TrimSpace(data)) == 0 {
		findings = append(findings, "SKILL.md is empty")
		return false, findings
	}
	findings = append(findings, validateFrontmatter(data)...)

	absRoot, err := filepath.Abs(skillDir)
	if err != nil {
		findings = append(findings, "could not resolve skill directory: "+err.Error())
		return false, findings
	}
	realRoot, err := filepath.EvalSymlinks(absRoot)
	if err != nil {
		findings = append(findings, "could not resolve skill directory symlinks: "+err.Error())
		return false, findings
	}

	references := extractLocalReferences(string(data))
	for _, ref := range references {
		if isAbsoluteReference(ref) {
			findings = append(findings, "absolute path reference not allowed: "+ref)
			continue
		}

		resolved := filepath.Join(absRoot, filepath.FromSlash(ref))
		rel, err := filepath.Rel(absRoot, resolved)
		if err != nil || rel == ".." || strings.HasPrefix(rel, ".."+string(filepath.Separator)) {
			findings = append(findings, "path traversal outside skill directory: "+ref)
			continue
		}

		_, err = os.Stat(resolved)
		if os.IsNotExist(err) {
			findings = append(findings, "missing referenced path: "+ref)
			continue
		}
		if err != nil {
			findings = append(findings, "could not inspect referenced path "+ref+": "+err.Error())
			continue
		}

		realPath, err := filepath.EvalSymlinks(resolved)
		if err != nil {
			findings = append(findings, "could not resolve referenced path symlinks "+ref+": "+err.Error())
			continue
		}
		realRel, err := filepath.Rel(realRoot, realPath)
		if err != nil || realRel == ".." || strings.HasPrefix(realRel, ".."+string(filepath.Separator)) {
			findings = append(findings, "path traversal outside skill directory via symlink: "+ref)
		}
	}

	warnings, err := findUnusedFiles(realRoot, references)
	if err != nil {
		findings = append(findings, "could not scan skill directory for unused files: "+err.Error())
	}
	valid := len(findings) == 0
	return valid, append(findings, warnings...)
}

func validateFrontmatter(data []byte) []string {
	lines := strings.Split(strings.TrimPrefix(string(data), "\uFEFF"), "\n")
	if len(lines) == 0 || strings.TrimSpace(lines[0]) != "---" {
		return []string{"SKILL.md must start with YAML frontmatter"}
	}

	end := -1
	for i := 1; i < len(lines); i++ {
		if strings.TrimSpace(lines[i]) == "---" {
			end = i
			break
		}
	}
	if end < 0 {
		return []string{"SKILL.md frontmatter is missing its closing ---"}
	}

	fields := map[string]bool{
		"name":        false,
		"description": false,
	}
	for i := 1; i < end; i++ {
		line := lines[i]
		if strings.TrimLeft(line, " \t") != line || strings.TrimSpace(line) == "" || strings.HasPrefix(strings.TrimSpace(line), "#") {
			continue
		}
		key, value, ok := strings.Cut(line, ":")
		if !ok {
			continue
		}
		key = strings.TrimSpace(key)
		if _, required := fields[key]; !required {
			continue
		}
		fields[key] = nonEmptyFrontmatterValue(value, lines, i)
	}

	findings := make([]string, 0)
	for _, field := range []string{"name", "description"} {
		if !fields[field] {
			findings = append(findings, `SKILL.md frontmatter requires a non-empty "`+field+`" field`)
		}
	}
	return findings
}

func nonEmptyFrontmatterValue(value string, lines []string, lineIndex int) bool {
	value = strings.TrimSpace(value)
	if len(value) >= 2 && ((value[0] == '"' && strings.Contains(value[1:], `"`)) ||
		(value[0] == '\'' && strings.Contains(value[1:], `'`))) {
		quote := value[0]
		end := strings.LastIndexByte(value, quote)
		return end > 0 && strings.TrimSpace(value[1:end]) != ""
	}
	if comment := strings.Index(value, " #"); comment >= 0 {
		value = value[:comment]
	}
	value = strings.TrimSpace(value)

	if value == "|" || value == "|-" || value == "|+" || value == ">" || value == ">-" || value == ">+" {
		for i := lineIndex + 1; i < len(lines); i++ {
			line := lines[i]
			if strings.TrimSpace(line) == "" {
				continue
			}
			return strings.HasPrefix(line, " ") || strings.HasPrefix(line, "\t")
		}
		return false
	}

	return value != "" && !strings.EqualFold(value, "null") && value != "~" && !strings.HasPrefix(value, "#")
}

func extractLocalReferences(markdown string) []string {
	references := make([]string, 0)
	seen := make(map[string]struct{})
	add := func(ref string, requirePathPrefix bool) {
		ref = strings.TrimSpace(ref)
		ref = strings.TrimRight(ref, ".,;:!?)]}")
		if ref == "" || isExternalReference(ref) {
			return
		}

		if requirePathPrefix {
			if !isAbsoluteReference(ref) && !strings.HasPrefix(ref, ".") {
				return
			}
		} else if !isAbsoluteReference(ref) && !strings.ContainsAny(ref, `/\.`) {
			return
		}

		if index := strings.IndexAny(ref, "?#"); index >= 0 {
			ref = ref[:index]
		}
		if ref == "" {
			return
		}
		if _, ok := seen[ref]; ok {
			return
		}
		seen[ref] = struct{}{}
		references = append(references, ref)
	}

	for _, match := range backtickPattern.FindAllStringSubmatch(markdown, -1) {
		for _, token := range strings.Fields(match[1]) {
			add(token, true)
		}
	}

	markdownWithoutCode := backtickPattern.ReplaceAllString(markdown, "")
	for _, match := range markdownLinkPattern.FindAllStringSubmatch(markdownWithoutCode, -1) {
		ref := match[1]
		if ref == "" {
			ref = match[2]
		}
		add(ref, false)
	}

	return references
}

func isAbsoluteReference(ref string) bool {
	if filepath.IsAbs(ref) || strings.HasPrefix(ref, "/") || strings.HasPrefix(ref, `\\`) {
		return true
	}
	return len(ref) >= 3 &&
		((ref[0] >= 'a' && ref[0] <= 'z') || (ref[0] >= 'A' && ref[0] <= 'Z')) &&
		ref[1] == ':' &&
		(ref[2] == '/' || ref[2] == '\\')
}

func findUnusedFiles(root string, references []string) ([]string, error) {
	referenced := make(map[string]struct{}, len(references))
	for _, ref := range references {
		if isAbsoluteReference(ref) {
			continue
		}
		resolved := filepath.Clean(filepath.Join(root, filepath.FromSlash(ref)))
		realPath, err := filepath.EvalSymlinks(resolved)
		if err != nil {
			continue
		}
		rel, err := filepath.Rel(root, realPath)
		if err != nil || rel == ".." || strings.HasPrefix(rel, ".."+string(filepath.Separator)) {
			continue
		}
		referenced[rel] = struct{}{}
	}

	warnings := make([]string, 0)
	err := filepath.WalkDir(root, func(path string, entry os.DirEntry, walkErr error) error {
		if walkErr != nil {
			return walkErr
		}
		if path == root {
			return nil
		}
		rel, err := filepath.Rel(root, path)
		if err != nil {
			return err
		}
		for _, part := range strings.Split(rel, string(filepath.Separator)) {
			if strings.HasPrefix(part, ".") {
				if entry.IsDir() {
					return filepath.SkipDir
				}
				return nil
			}
		}
		if entry.Type()&os.ModeSymlink != 0 {
			return nil
		}
		if entry.IsDir() {
			return nil
		}
		info, err := entry.Info()
		if err != nil {
			return err
		}
		if !info.Mode().IsRegular() || rel == "SKILL.md" {
			return nil
		}
		for ref := range referenced {
			if ref == "." || rel == ref || strings.HasPrefix(rel, ref+string(filepath.Separator)) {
				return nil
			}
		}
		warnings = append(warnings, "WARN: unused file: "+filepath.ToSlash(rel))
		return nil
	})
	return warnings, err
}

func isExternalReference(ref string) bool {
	lowerRef := strings.ToLower(ref)
	return strings.HasPrefix(lowerRef, "http://") ||
		strings.HasPrefix(lowerRef, "https://") ||
		strings.HasPrefix(lowerRef, "mailto:")
}
