package view

import (
	"log"
	"path"
	"strings"

	"github.com/iainjreid/source/git"
	"github.com/iainjreid/source/view/nav"
	"github.com/alecthomas/chroma/v2"
	"github.com/alecthomas/chroma/v2/formatters/html"
	"github.com/alecthomas/chroma/v2/lexers"
	"github.com/alecthomas/chroma/v2/styles"
)

type View struct {
	Branches     []git.Branch
	Nav          nav.Nav
	TreeEntryMap *git.TreeEntryMap
	DirPath      string
	Commit       git.Commit
	repo         *git.Repo
	File         []*git.Line
	FileName     string
	Contents     string
}

func New(repo *git.Repo) *View {
	branches, err := repo.GetBranches()

	if err != nil {
		panic(err)
	}

	return &View{
		Branches: branches,
		Nav: nav.Nav{
			Items: []*nav.NavItem{
				nav.NewItem("Code", "/", true),
				nav.NewItem("Issues", "/issues", false),
			},
		},
		repo: repo,
	}
}

// Clean this hack up!
func (v *View) LoadCommit(revision string) *View {
	v.Commit = *git.Must(v.repo.GetCommit(revision))
	return v
}

// Clean this hack up!
func (v *View) LoadDir(dirpath string) *View {
	var err error
	v.DirPath = dirpath
	v.TreeEntryMap, err = v.Commit.GetTree(dirpath[1:], false)
	if err != nil {
		panic(err)
	}
	return v
}

func (v *View) LoadBlob(filename string, blame bool) (*View, error) {
	file, err := v.Commit.GetFileContents(path.Join(v.DirPath, filename)[1:], blame)

	if err != nil {
		return nil, err
	}

	lexer := lexers.Match(filename)

	if lexer != nil {
		lexer = chroma.Coalesce(lexer)
	} else {
		lexer = lexers.Fallback
	}

	style := styles.Get("github")
	if style == nil {
		log.Fatalf("no style")
		style = styles.Fallback
	}

	formatter := html.New(
		html.WithLineNumbers(true),
		html.WithClasses(true),
		html.WithAllClasses(true),
		html.LineNumbersInTable(true),
		html.WithLinkableLineNumbers(true, "L"),
		html.HighlightLines([][2]int{{1, 2}}),
	)

	var source string
	for _, line := range file {
		source += line.Text + "\n"
	}

	iter, err := lexer.Tokenise(nil, source)

	if err != nil {
		return nil, err
	}

	lineBuilder := new(strings.Builder)

	err = formatter.Format(lineBuilder, style, iter)
	if err != nil {
		return nil, err
	}

	// cssBuilder := new(strings.Builder)
	// err = formatter.WriteCSS(cssBuilder, style)
	// if err != nil {
	// 	return nil, err
	// }
	// log.Println(cssBuilder.String())

	v.FileName = filename
	v.File = file
	v.Contents = lineBuilder.String()

	log.Println(v.Contents)

	return v, nil
}
