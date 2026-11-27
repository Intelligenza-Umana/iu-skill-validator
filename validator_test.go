package main

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

const testFrontmatter = "---\nname: test-skill\ndescription: A test skill.\n---\n"

func TestValidateSkillDirFixtures(t *testing.T) {
	tests := []struct {
		name        string
		skillDir    string
		wantValid   bool
		wantFinding string
	}{
		{
			name:      "valid",
			skillDir:  "./fixtures/valid",
			wantValid: true,
		},
		{
			name:        "missing reference",
			skillDir:    "./fixtures/invalid-missing",
			wantFinding: "missing referenced path",
		},
		{
			name:        "path traversal",
			skillDir:    "./fixtures/invalid-traversal",
			wantFinding: "path traversal",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			valid, findings := ValidateSkillDir(tt.skillDir)
			if valid != tt.wantValid {
				t.Fatalf("expected valid=%v, got valid=%v findings=%v", tt.wantValid, valid, findings)
			}
			if tt.wantFinding != "" {
				for _, finding := range findings {
					if strings.Contains(finding, tt.wantFinding) {
						return
					}
				}
				t.Fatalf("expected a finding containing %q, got %v", tt.wantFinding, findings)
			}
			if len(findings) != 0 {
				t.Fatalf("expected no findings, got %v", findings)
			}
			if findings == nil {
				t.Fatal("expected findings to be a non-nil empty slice")
			}
		})
	}
}

func TestValidateSkillDirRejectsAbsoluteReferences(t *testing.T) {
	tests := []struct {
		name string
		ref  string
		link bool
	}{
		{name: "unix markdown link", ref: "/etc/passwd", link: true},
		{name: "windows markdown link", ref: `C:\Windows\System32\cmd.exe`, link: true},
		{name: "unix backtick", ref: "/etc/passwd"},
		{name: "windows backtick", ref: `C:\Windows\System32\cmd.exe`},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			skillDir := t.TempDir()
			ref := "[" + tt.ref + "](" + tt.ref + ")"
			if !tt.link {
				ref = "`" + tt.ref + "`"
			}
			skill := testFrontmatter + ref + "\n"
			if err := os.WriteFile(filepath.Join(skillDir, "SKILL.md"), []byte(skill), 0o644); err != nil {
				t.Fatal(err)
			}

			valid, findings := ValidateSkillDir(skillDir)
			if valid {
				t.Fatal("expected absolute reference to be invalid")
			}
			if len(findings) != 1 || !strings.Contains(findings[0], "absolute path reference not allowed") {
				t.Fatalf("expected absolute path finding, got %v", findings)
			}
		})
	}
}

func TestValidateSkillDirRejectsSymlinkEscape(t *testing.T) {
	skillDir := t.TempDir()
	outsideDir := t.TempDir()
	outsideFile := filepath.Join(outsideDir, "secret.txt")
	if err := os.WriteFile(outsideFile, []byte("secret"), 0o644); err != nil {
		t.Fatal(err)
	}
	link := filepath.Join(skillDir, "linked")
	if err := os.Symlink(outsideDir, link); err != nil {
		t.Skipf("symlinks are unavailable: %v", err)
	}
	skill := "[linked file](./linked/secret.txt)\n"
	skill = testFrontmatter + skill
	if err := os.WriteFile(filepath.Join(skillDir, "SKILL.md"), []byte(skill), 0o644); err != nil {
		t.Fatal(err)
	}

	valid, findings := ValidateSkillDir(skillDir)
	if valid {
		t.Fatal("expected symlink escaping the skill directory to be invalid")
	}
	if len(findings) != 1 || !strings.Contains(findings[0], "outside skill directory via symlink") {
		t.Fatalf("expected symlink escape finding, got %v", findings)
	}
}

func TestValidateSkillDirRequiresFrontmatterFields(t *testing.T) {
	tests := []struct {
		name        string
		content     string
		wantFinding string
	}{
		{
			name:        "missing frontmatter",
			content:     "# No metadata",
			wantFinding: "must start with YAML frontmatter",
		},
		{
			name:        "empty name",
			content:     "---\nname: \"\"\ndescription: A test skill.\n---\nBody",
			wantFinding: `non-empty "name" field`,
		},
		{
			name:        "empty description block",
			content:     "---\nname: test-skill\ndescription: |\n---\nBody",
			wantFinding: `non-empty "description" field`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			skillDir := t.TempDir()
			if err := os.WriteFile(filepath.Join(skillDir, "SKILL.md"), []byte(tt.content), 0o644); err != nil {
				t.Fatal(err)
			}
			valid, findings := ValidateSkillDir(skillDir)
			if valid {
				t.Fatalf("expected invalid frontmatter, got findings: %v", findings)
			}
			for _, finding := range findings {
				if strings.Contains(finding, tt.wantFinding) {
					return
				}
			}
			t.Fatalf("expected finding containing %q, got %v", tt.wantFinding, findings)
		})
	}
}

func TestValidateSkillDirWarnsForUnusedFilesWithoutFailing(t *testing.T) {
	skillDir := t.TempDir()
	skill := testFrontmatter + "[used](./assets/used.txt)\n"
	if err := os.WriteFile(filepath.Join(skillDir, "SKILL.md"), []byte(skill), 0o644); err != nil {
		t.Fatal(err)
	}
	for _, path := range []string{
		"assets/used.txt",
		"assets/unused.txt",
		".hidden",
		".cache/ignored.txt",
	} {
		fullPath := filepath.Join(skillDir, path)
		if err := os.MkdirAll(filepath.Dir(fullPath), 0o755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(fullPath, []byte("content"), 0o644); err != nil {
			t.Fatal(err)
		}
	}

	valid, findings := ValidateSkillDir(skillDir)
	if !valid {
		t.Fatalf("warnings should not invalidate the skill: %v", findings)
	}
	if len(findings) != 1 || findings[0] != "WARN: unused file: assets/unused.txt" {
		t.Fatalf("expected only the unused-file warning, got %v", findings)
	}
}

func TestValidateSkillDirRejectsEmptySkillFile(t *testing.T) {
	skillDir := t.TempDir()
	if err := os.WriteFile(filepath.Join(skillDir, "SKILL.md"), []byte(" \n\t "), 0o644); err != nil {
		t.Fatal(err)
	}

	valid, findings := ValidateSkillDir(skillDir)
	if valid || len(findings) != 1 {
		t.Fatalf("expected one finding for empty SKILL.md, got valid=%v findings=%v", valid, findings)
	}
}
