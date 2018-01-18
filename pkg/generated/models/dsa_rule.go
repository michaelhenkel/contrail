package models

// DsaRule

import "encoding/json"

// DsaRule
type DsaRule struct {
	DsaRuleEntry *DiscoveryServiceAssignmentType `json:"dsa_rule_entry,omitempty"`
	FQName       []string                        `json:"fq_name,omitempty"`
	IDPerms      *IdPermsType                    `json:"id_perms,omitempty"`
	Perms2       *PermType2                      `json:"perms2,omitempty"`
	ParentUUID   string                          `json:"parent_uuid,omitempty"`
	DisplayName  string                          `json:"display_name,omitempty"`
	Annotations  *KeyValuePairs                  `json:"annotations,omitempty"`
	UUID         string                          `json:"uuid,omitempty"`
	ParentType   string                          `json:"parent_type,omitempty"`
}

// String returns json representation of the object
func (model *DsaRule) String() string {
	b, _ := json.Marshal(model)
	return string(b)
}

// MakeDsaRule makes DsaRule
func MakeDsaRule() *DsaRule {
	return &DsaRule{
		//TODO(nati): Apply default
		ParentUUID:   "",
		DsaRuleEntry: MakeDiscoveryServiceAssignmentType(),
		FQName:       []string{},
		IDPerms:      MakeIdPermsType(),
		Perms2:       MakePermType2(),
		DisplayName:  "",
		Annotations:  MakeKeyValuePairs(),
		UUID:         "",
		ParentType:   "",
	}
}

// MakeDsaRuleSlice() makes a slice of DsaRule
func MakeDsaRuleSlice() []*DsaRule {
	return []*DsaRule{}
}
