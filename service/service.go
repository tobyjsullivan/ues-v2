package service

import (
    "github.com/tobyjsullivan/ues-v2/events"
    "encoding/json"
    "log"
    "os"
)

const (
    eventTypeAccountOpened = "AccountOpened"
)

var (
    logger *log.Logger
)

func init() {
    logger = log.New(os.Stdout, "[service] ", 0)
}

type ServiceAggregate struct {
    ID string
    AccountIDs []string
}

type accountOpenedEvent struct {
    AccountID string `json:"accountId"`
}

func LoadAggregate(entityId string) (*ServiceAggregate, error) {
    svcEvents, err := events.GetEvents(entityId)
    if err != nil {
        return false, err
    }

    var aggregate ServiceAggregate

    for _, e := range svcEvents {
        switch e.Type {
        case eventTypeAccountOpened:
            content, err := json.Marshal(e.Data)
            if err != nil {
                logger.Println("Error marshalling event data.", err.Error())
                return nil, err
            }

            var event accountOpenedEvent
            err = json.Unmarshal(content, &event)
            if err != nil {
                logger.Println("Error unmarshalling event data.", err.Error())
                return nil, err
            }

            aggregate.AccountIDs = append(aggregate.AccountIDs, event.AccountID)
        default:
            continue
        }
    }
}
