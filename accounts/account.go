package accounts

import "github.com/tobyjsullivan/ues-v2/events"

type AccountAggregate struct {
    ID string
    Email string
    PasswordHash string
}

func LoadAggregate(accountId string) (*AccountAggregate, error) {
    events, err := events.GetEvents(accountId)
    if err != nil {
        return nil, err
    }

    agg := AccountAggregate{
        ID: accountId,
    }
    for _, e := range events {
        switch e.Type {
        default:
            continue
        }
    }

    return &agg, nil
}
