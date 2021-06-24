package clipboard

import (
	"bufio"
	"crypto/md5"
	"fmt"
	"io"
	"regexp"
	"strings"
)

type Mark struct {
	Id      string
	Text    string
	Title   string
	Link    string
	LinkStr string
	Source  string
	Tags    []string
	Opinion string
}

func (s Mark) tagString() string {
	if len(s.Tags) == 0 {
		return ""
	}
	sb := strings.Builder{}
	for _, t := range s.Tags {
		if strings.HasPrefix(t, "#") {
			sb.WriteString(strings.TrimSpace(t))
		} else {
			sb.WriteString("#")
			sb.WriteString(strings.TrimSpace(t))
		}
		sb.WriteString(" ")
	}
	return sb.String()[0 : sb.Len()-1]
}

func (s Mark) idString() string {
	if len(s.Id) == 0 {
		return NewId(s.Text)
	}
	return s.Id
}

func NewId(t string) string {
	data := []byte(t)
	s := fmt.Sprintf("%x", md5.Sum(data))
	return s[0:8] + s[24:32]
}

func (s Mark) String() string {
	sb := strings.Builder{}

	sb.WriteString(fmt.Sprintf("id: %s\n\n", s.idString()))
	sb.WriteString(fmt.Sprintf("text:\n```\n%s\n```\n\n", s.Text))
	if s.LinkStr == "" {
		sb.WriteString(fmt.Sprintf("link: [%s](%s)\n\n", s.Title, s.Link))
	} else {
		sb.WriteString(fmt.Sprintf("link: %s\n\n", s.LinkStr))
	}
	sb.WriteString(fmt.Sprintf("source: %s\n\n", s.Source))
	sb.WriteString(fmt.Sprintf("tag: %s\n\n", s.tagString()))
	sb.WriteString(fmt.Sprintf("opinion: %s\n\n", s.Opinion))

	return sb.String()
}

func ParseTag(s string) []string {
	r := regexp.MustCompile("\\s+|,")
	ss := r.Split(s, -1)
	var tags []string
	for _, t := range ss {
		var tt string
		if strings.HasPrefix(t, "#") {
			tt = strings.TrimSpace(t[1:])
		} else {
			tt = t
		}
		if len(tt) == 0 {
			continue
		}
		tags = append(tags, tt)
	}
	return tags
}

func ParseStoreFile(in io.Reader) ([]Mark, string) {
	br := bufio.NewReader(in)
	var ss Mark
	r := make([]Mark, 0)
	textflag := false
	opinionflag := false
	changeFlag := func(text, opinion bool) {
		textflag = text
		opinionflag = opinion
	}
	fileTitle := ""
	clean := func() {
		ss.Text = strings.TrimSpace(ss.Text)
		if strings.HasPrefix(ss.Text, "```") {
			ss.Text = ss.Text[4:]
			ss.Text = ss.Text[:len(ss.Text)-4]
			ss.Text = strings.TrimSpace(ss.Text)
		}
		ss.Opinion = strings.TrimSpace(ss.Opinion)
	}
	for {
		line, eof := br.ReadString('\n')
		if eof == io.EOF || eof != nil {
			break
		}
		line = strings.TrimSpace(line)
		if len(fileTitle) == 0 {
			if strings.HasPrefix(line, "# ") {
				fileTitle = line[2:]
				continue
			}
		}

		if !textflag && strings.HasPrefix(line, "---") {
			clean()
			r = append(r, ss)
			ss = Mark{}
			changeFlag(false, false)
			continue
		}

		if line == "" && !textflag && !opinionflag {
			continue
		}

		switch prefix(line) {
		case "id":
			s := line[3:]
			ss.Id = strings.TrimSpace(s)
			changeFlag(false, false)
		case "link":
			s := line[5:]
			ss.LinkStr = strings.TrimSpace(s)
			changeFlag(false, false)
		case "source":
			s := line[7:]
			ss.Source = strings.TrimSpace(s)
			changeFlag(false, false)
		case "tag":
			s := line[4:]
			ss.Tags = ParseTag(strings.TrimSpace(s))
			changeFlag(false, false)
		case "text":
			s := line[5:]
			ss.Text = strings.TrimSpace(s)
			changeFlag(true, false)
		case "opinion":
			s := line[8:]
			ss.Opinion = strings.TrimSpace(s)
			changeFlag(false, true)
		default:
			if textflag {
				ss.Text += strings.TrimSpace(line) + "\n"
			}
			if opinionflag {
				ss.Opinion += strings.TrimSpace(line) + "\n"
			}
		}
	}

	// 对最后一个分组进行处理
	if strings.TrimSpace(ss.Text) != "" {
		if len(r) > 1 && ss.idString() != r[len(r)-1].idString() {
			clean()
			r = append(r, ss)
		}
	}
	if len(r) == 0 {
		return nil, fileTitle
	}

	return r[1:], fileTitle
}

func Merge(oldMarks []Mark, newMarks []Mark) []Mark {
	index := make(map[string]int)
	for i, o := range oldMarks {
		index[o.idString()] = i
	}
	for _, m := range newMarks {
		if i, ok := index[m.idString()]; ok {
			oldMarks[i] = m
		} else {
			oldMarks = append(oldMarks, m)
		}
	}
	return oldMarks
}

func FileString(title string, marks []Mark) *strings.Builder {
	sb := strings.Builder{}
	sb.WriteString(fmt.Sprintf("# %s\n\n---\n\n", title))
	for _, m := range marks {
		sb.WriteString(m.String())
		sb.WriteString("\n---\n\n")
	}
	return &sb
}

func prefix(a string) string {
	ss := strings.SplitN(a, ":", 2)
	if len(ss) == 0 {
		return ""
	}
	return ss[0]
}
