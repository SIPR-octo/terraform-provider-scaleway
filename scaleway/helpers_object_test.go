package scaleway

import (
	"testing"

	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/scaleway/scaleway-sdk-go/scw"
	"github.com/stretchr/testify/assert"
)

func TestExpandObjectBucketTags(t *testing.T) {
	tests := []struct {
		name string
		tags interface{}
		want []*s3.Tag
	}{
		{
			name: "no tags",
			tags: map[string]interface{}{},
			want: []*s3.Tag(nil),
		},
		{
			name: "single tag",
			tags: map[string]interface{}{
				"key1": "val1",
			},
			want: []*s3.Tag{
				{Key: scw.StringPtr("key1"), Value: scw.StringPtr("val1")},
			},
		},
		{
			name: "many tags",
			tags: map[string]interface{}{
				"key1": "val1",
				"key2": "val2",
				"key3": "val3",
			},
			want: []*s3.Tag{
				{Key: scw.StringPtr("key1"), Value: scw.StringPtr("val1")},
				{Key: scw.StringPtr("key2"), Value: scw.StringPtr("val2")},
				{Key: scw.StringPtr("key3"), Value: scw.StringPtr("val3")},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.ElementsMatch(t, tt.want, expandObjectBucketTags(tt.tags))
		})
	}
}

func TestExpandObjectBucketLifecycleRules(t *testing.T) {
	type args struct {
		rawLifecycleRules []interface{}
		bucket            string
	}
	tests := []struct {
		name string
		got  args
		want []*s3.LifecycleRule
	}{
		{
			name: "no rules",
			got: args{
				rawLifecycleRules: []interface{}{},
				bucket:            "toto",
			},
			want: []*s3.LifecycleRule{},
		},
		{
			name: "single rule",
			got: args{
				rawLifecycleRules: []interface{}{
					map[string]interface{}{
						"name":            "rule1",
						"status":          "Enabled",
						"expiration_days": 10,
					},
				},
				bucket: "toto",
			},
			want: []*s3.LifecycleRule{
				{
					Filter: &s3.LifecycleRuleFilter{Prefix: nil},
					Expiration: &s3.LifecycleExpiration{
						Days: scw.Int64Ptr(10),
						// Date: &time.Time{},
						// ExpiredObjectDeleteMarker: nil,
					},
					// Filter: &s3.LifecycleRuleFilter{
					// 	And: &s3.LifecycleRuleAndOperator{
					// 		Prefix: nil,
					// 		Tags:   nil,
					// 	},
					// 	Prefix: nil,
					// 	Tag: &s3.Tag{
					// 		Key:   nil,
					// 		Value: nil,
					// 	},
					// },
					ID:     scw.StringPtr("rule1"),
					Status: scw.StringPtr("Enabled"),
				},
			},
		},
		{
			name: "multiple rules",
			got: args{
				rawLifecycleRules: []interface{}{
					map[string]interface{}{
						"name":            "rule1",
						"status":          "Enabled",
						"expiration_days": 10,
					},
					map[string]interface{}{
						"name":            "rule2",
						"status":          "Enabled",
						"expiration_days": 25,
					},
					map[string]interface{}{
						"name":            "rule3",
						"status":          "Disabled",
						"expiration_days": 30,
					},
				},
				bucket: "toto",
			},
			want: []*s3.LifecycleRule{
				{
					Filter: &s3.LifecycleRuleFilter{Prefix: nil},
					Expiration: &s3.LifecycleExpiration{
						Days: scw.Int64Ptr(10),
					},
					ID:     scw.StringPtr("rule1"),
					Status: scw.StringPtr("Enabled"),
				},
				{
					Filter: &s3.LifecycleRuleFilter{Prefix: nil},
					Expiration: &s3.LifecycleExpiration{
						Days: scw.Int64Ptr(25),
					},
					ID:     scw.StringPtr("rule2"),
					Status: scw.StringPtr("Enabled"),
				},
				{
					Filter: &s3.LifecycleRuleFilter{Prefix: nil},
					Expiration: &s3.LifecycleExpiration{
						Days: scw.Int64Ptr(30),
					},
					ID:     scw.StringPtr("rule3"),
					Status: scw.StringPtr("Disabled"),
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.ElementsMatch(t, tt.want, expandBucketLifecycleRules(tt.got.rawLifecycleRules, tt.got.bucket))
		})
	}
}
