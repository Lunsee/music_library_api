package testData

import (
	"music_library_api/internal/models"
	"time"
)

// test data
var TestData = []models.Song{
	{Group: "Muse", Song: "Supermassive Black Hole", ReleaseDate: time.Date(2006, 7, 16, 0, 0, 0, 0, time.UTC), Text: "Some lyrics", Link: "https://www.youtube.com/watch?v=Xsp3_a-PMTw"},
	{Group: "Radiohead", Song: "Creep", ReleaseDate: time.Date(1992, 9, 21, 0, 0, 0, 0, time.UTC), Text: "Some lyrics", Link: "https://www.youtube.com/watch?v=d1J5kz5T6Bo"},
	{Group: "Nirvana", Song: "Smells Like Teen Spirit", ReleaseDate: time.Date(1991, 9, 10, 0, 0, 0, 0, time.UTC), Text: "Some lyrics", Link: "https://www.youtube.com/watch?v=hTWKbfoikeg"},
	{Group: "Queen", Song: "Bohemian Rhapsody", ReleaseDate: time.Date(1975, 10, 31, 0, 0, 0, 0, time.UTC), Text: "Some lyrics", Link: "https://www.youtube.com/watch?v=fJ9rUzIMcZQ"},
	{Group: "The Beatles", Song: "Hey Jude", ReleaseDate: time.Date(1968, 8, 26, 0, 0, 0, 0, time.UTC), Text: "Some lyrics", Link: "https://www.youtube.com/watch?v=A_MjCqQoLLA"},
	{Group: "Pink Floyd", Song: "Comfortably Numb", ReleaseDate: time.Date(1979, 11, 30, 0, 0, 0, 0, time.UTC), Text: "Some lyrics", Link: "https://www.youtube.com/watch?v=JwYX52BP2Sk"},
	{Group: "Led Zeppelin", Song: "Stairway to Heaven", ReleaseDate: time.Date(1971, 11, 8, 0, 0, 0, 0, time.UTC), Text: "Some lyrics", Link: "https://www.youtube.com/watch?v=QkF3oxziUI4"},
	{Group: "AC/DC", Song: "Back In Black", ReleaseDate: time.Date(1980, 7, 25, 0, 0, 0, 0, time.UTC), Text: "Some lyrics", Link: "https://www.youtube.com/watch?v=pAgnJDJN4VA"},
	{Group: "The Rolling Stones", Song: "Paint It Black", ReleaseDate: time.Date(1966, 5, 7, 0, 0, 0, 0, time.UTC), Text: "Some lyrics", Link: "https://www.youtube.com/watch?v=O4irXQhgMqg"},
	{Group: "U2", Song: "With or Without You", ReleaseDate: time.Date(1987, 3, 21, 0, 0, 0, 0, time.UTC), Text: "Some lyrics", Link: "https://www.youtube.com/watch?v=XmSdTa9kaiQ"},
	{Group: "Metallica", Song: "Enter Sandman", ReleaseDate: time.Date(1991, 7, 29, 0, 0, 0, 0, time.UTC), Text: "Some lyrics", Link: "https://www.youtube.com/watch?v=CD-E-LDc384"},
	{Group: "Oasis", Song: "Wonderwall", ReleaseDate: time.Date(1995, 10, 30, 0, 0, 0, 0, time.UTC), Text: "Some lyrics", Link: "https://www.youtube.com/watch?v=bx1Bh8ZvH84"},
	{
		Group:       "test",
		Song:        "test",
		ReleaseDate: time.Date(1995, 10, 30, 0, 0, 0, 0, time.UTC),
		Text: `Test test, this is the first verse.
				The melody flows like a gentle breeze.

				Test test, this is the second verse.
				The rhythm keeps you on your toes, dancing in the streets.

				Test test, this is the third verse.
				We sing out loud under the starry sky.

				Test test, this is the final verse.
				The music fades, but the memories remain forever.`,
		Link: "https://www.youtube.com/watch?TESTTEST",
	},
}
