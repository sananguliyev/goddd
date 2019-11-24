package postgres

import (
	"fmt"
	"github.com/SananGuliyev/goddd/domain/throwable"
	rqlParser "github.com/tbaud0n/go-rql-parser"
	"strings"
)

type parser struct {
}

func (p *parser) Parse(filter string, limit int, offset int) (string, error) {
	rParser := rqlParser.NewParser()

	whereNode, err := rParser.Parse(strings.NewReader(filter))

	if err != nil {
		return "", throwable.NewInvalidFilter("Invalid filter string")
	}

	where, err := rqlParser.NewSqlTranslator(whereNode).Where()
	if err != nil {
		return "", throwable.NewInvalidFilter("Invalid where operation")
	} else if len(where) > 0 {
		where = "WHERE " + where
	}

	limitString := fmt.Sprintf("limit(%d,%d)", limit, offset)
	limitNode, err := rParser.Parse(strings.NewReader(limitString))

	if err != nil {
		return "", throwable.NewInvalidFilter("Invalid limit operation")
	}

	limitOffset, err := rqlParser.NewSqlTranslator(limitNode).Sql()
	if err != nil {
		return "", throwable.NewInvalidFilter("Invalid pagination parameters")
	}

	return where + limitOffset, nil
}

func NewParser() *parser {
	return &parser{}
}
