package mongo

import (
	"github.com/AdhityaRamadhanus/gopatrol"
	"github.com/AdhityaRamadhanus/gopatrol/config"
	"gopkg.in/mgo.v2"
)

type LoggingService struct {
	session  *mgo.Session
	CollName string
}

func NewLoggingService(session *mgo.Session, collName string) *LoggingService {
	LogColl := session.DB(config.DatabaseName).C(collName)

	LogColl.Create(&mgo.CollectionInfo{
		Capped: true,
		// Set Max Size in bytes to 5 GB (just a guess number)
		MaxBytes: 5000 * 1000,
		MaxDocs:  100000,
	})

	// Ensure Index
	LogColl.EnsureIndex(mgo.Index{
		Key:        []string{"slug"},
		Unique:     false,
		DropDups:   false,
		Background: true,
		Sparse:     true,
	})

	return &LoggingService{
		session:  session,
		CollName: collName,
	}
}

// for more flexible use
func (p *LoggingService) CopySession() *mgo.Session {
	return p.session.Copy()
}

func (p *LoggingService) InsertLog(result gopatrol.Result) error {
	copySession := p.session.Copy()
	defer copySession.Close()
	EndpointColl := copySession.DB(config.DatabaseName).C(p.CollName)
	return EndpointColl.Insert(result)
}

func (p *LoggingService) GetAllLogs(q map[string]interface{}) ([]gopatrol.Result, error) {
	copySession := p.session.Copy()
	defer copySession.Close()

	LogColl := copySession.DB(config.DatabaseName).C(p.CollName)
	logs := []gopatrol.Result{}
	MongoQuery := LogColl.Find(q["query"])
	if ok, val := q["pagination"].(bool); ok && val {
		MongoQuery.
			Skip(q["page"].(int) * q["limit"].(int)).
			Limit(q["limit"].(int))
	}

	if err := MongoQuery.All(&logs); err != nil {
		return nil, err
	}
	return logs, nil
}
