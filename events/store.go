package events

import (
    "sync"
    "github.com/aws/aws-sdk-go/service/dynamodb"
    "github.com/aws/aws-sdk-go/aws/session"
    "github.com/aws/aws-sdk-go/aws"
    "errors"
    "strconv"
    "github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
    "os"
    "github.com/aws/aws-sdk-go/aws/credentials"
)

const (
    columnEntityID = "Entity ID"
    columnVersion = "Version"
    columnType = "Type"
    columnData = "Data"
)

var (
    dbTable string
    locks map[string]*sync.RWMutex
    db *dynamodb.DynamoDB
    lockLock sync.Mutex
)

func init() {
    require("AWS_ACCESS_KEY_ID")
    require("AWS_SECRET_ACCESS_KEY")
    awsRegion := require("AWS_REGION")
    dbTable = require("DYNAMODB_TABLE")

    sess := session.Must(session.NewSession(aws.Config{
        Credentials: credentials.NewEnvCredentials(),
        Region: aws.String(awsRegion),
    }))
    db = dynamodb.New(sess)
}

func require(envkey string) string {
    v := os.Getenv(envkey)
    if v == "" {
        panic("Required Env Var not found: "+envkey)
    }
    return v
}

type Event struct {
    EntityID string
    Version int
    Type string
    Data map[string]interface{}
}

func GetEvents(entityId string) ([]*Event, error) {
    lock := lock(entityId)
    lock.RLock()
    defer lock.RUnlock()

    res, err := db.Query(dynamodb.QueryInput{
        TableName:                  aws.String(dbTable),
        ConsistentRead:             aws.Bool(true),
        KeyConditionExpression:     aws.String("#entityId = :entityId"),
        ProjectionExpression:       aws.String("#version, #type, #data"),
        ExpressionAttributeNames:   map[string]*string{
            "#entityId":            aws.String(columnEntityID),
            "#version":             aws.String(columnVersion),
            "#type":                aws.String(columnType),
            "#data":                aws.String(columnData),
        },
        ExpressionAttributeValues:  map[string]*dynamodb.AttributeValue{
            ":entityId":            {S: aws.String(entityId)},
        },
    })
    if err != nil {
        logger.Println("Error querying events.", err.Error())
        return []*Event{}, err
    }

    out := make([]*Event, 0)
    for _, i := range res.Items {
        ver := aws.StringValue(i["Version"].N)
        if ver == "" {
            logger.Println("Version was empty.", entityId)
            return []*Event{}, errors.New("Version was empty.")
        }

        parsedVer, err := strconv.Atoi(ver)
        if err != nil {
            logger.Println("Error parsing version.", err.Error())
            return []*Event{}, err
        }

        data := make(map[string]interface{})
        err = dynamodbattribute.Unmarshal(i["Data"], data)
        if err != nil {
            logger.Println("Error parsing data.", err.Error())
            return []*Event{}, err
        }

        event := &Event{
            EntityID: entityId,
            Version: parsedVer,
            Type: aws.StringValue(i["Type"].S),
            Data: data,
        }

        out = append(out, event)
    }

    return out, nil
}

func Commit(event *Event) error {
    lock := lock(event.EntityID)
    lock.Lock()
    defer lock.Unlock()

    data, err := dynamodbattribute.Marshal(event.Data)
    if err != nil {
        logger.Println("Error marshalling data.", err.Error())
        return err
    }

    _, err = db.PutItem(&dynamodb.PutItemInput{
        TableName: aws.String(dbTable),
        Item: map[string]*dynamodb.AttributeValue {
            columnEntityID: {S: aws.String(event.EntityID)},
            columnVersion: {N: aws.String(strconv.Itoa(event.Version))},
            columnType: {S: aws.String(event.Type)},
            columnData: data,
        },
    })

    if err != nil {
        logger.Println("Error committing event.", err.Error())
        return err
    }

    logger.Println("Event committed successfully")
    return nil
}

func lock(entityId string) *sync.RWMutex {
    lockLock.Lock()
    defer lockLock.Unlock()

    lock := locks[entityId]
    if lock == nil {
        lock = &sync.RWMutex{}
        locks[entityId] = lock
    }
    return lock
}
