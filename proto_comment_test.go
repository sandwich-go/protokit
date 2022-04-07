package protokit

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestComments(t *testing.T) {
	Convey("different input for Filter", t, func() {
		for _, t := range []struct {
			// input
			in      []string
			comment *Comment
		}{
			{
				in: []string{"@actor @tell"},
				comment: &Comment{
					Tags:    map[string]string{
						"actor":"true",
						"tell":"true",
					},
				},
			},
			{
				in: []string{"@actor @tell some words"},
				comment: &Comment{
					Tags:    map[string]string{
						"actor":"true",
						"tell some words":"true",
					},
				},
			},
		} {
			c := NewComment(t.in)
			So(c.Tags,ShouldResemble,t.comment.Tags)
		}
	})
}
