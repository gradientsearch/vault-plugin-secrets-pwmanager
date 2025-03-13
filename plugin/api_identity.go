package secretsengine

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/hashicorp/vault/api"
)

// Identity is used to perform Identity operations on Vault.
type Identity struct {
	c *api.Client
}

// Identity is used to return the client for Identity API calls.
func (c *pwmanagerClient) Identity() *Identity {
	return &Identity{c: c.c}
}

func (c *Identity) EntityByID(entityID string) (Entity, error) {
	r := c.c.NewRequest("GET", fmt.Sprintf("/v1/identity/entity/id/%s", entityID))

	ctx, cancelFunc := context.WithCancel(context.Background())
	defer cancelFunc()
	resp, err := c.c.RawRequestWithContext(ctx, r)
	if err != nil {
		return Entity{}, err
	}
	defer resp.Body.Close()

	var result IdentityResponse
	err = json.NewDecoder(resp.Body).Decode(&result)
	if err != nil {
		return Entity{}, err
	}

	return result.Data, nil
}

type IdentityResponse struct {
	RequestID     string `json:"request_id"`
	LeaseID       string `json:"lease_id"`
	Renewable     bool   `json:"renewable"`
	LeaseDuration int    `json:"lease_duration"`
	Data          Entity `json:"data"`
	WrapInfo      any    `json:"wrap_info"`
	Warnings      any    `json:"warnings"`
	Auth          any    `json:"auth"`
	MountType     string `json:"mount_type"`
}

type Aliases struct {
	CanonicalID            string    `json:"canonical_id"`
	CreationTime           time.Time `json:"creation_time"`
	CustomMetadata         any       `json:"custom_metadata"`
	ID                     string    `json:"id"`
	LastUpdateTime         time.Time `json:"last_update_time"`
	Local                  bool      `json:"local"`
	MergedFromCanonicalIds any       `json:"merged_from_canonical_ids"`
	Metadata               any       `json:"metadata"`
	MountAccessor          string    `json:"mount_accessor"`
	MountPath              string    `json:"mount_path"`
	MountType              string    `json:"mount_type"`
	Name                   string    `json:"name"`
}
type Entity struct {
	Aliases           []Aliases `json:"aliases"`
	CreationTime      time.Time `json:"creation_time"`
	DirectGroupIds    []any     `json:"direct_group_ids"`
	Disabled          bool      `json:"disabled"`
	GroupIds          []any     `json:"group_ids"`
	ID                string    `json:"id"`
	InheritedGroupIds []any     `json:"inherited_group_ids"`
	LastUpdateTime    time.Time `json:"last_update_time"`
	MergedEntityIds   any       `json:"merged_entity_ids"`
	Metadata          any       `json:"metadata"`
	Name              string    `json:"name"`
	NamespaceID       string    `json:"namespace_id"`
	Policies          []any     `json:"policies"`
}
