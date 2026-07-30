package infra

import (
	"path/filepath"
	"testing"
)

var wantAwsChapters = []wantContentChapter{
	{"secure-architectures", []string{
		"iam-users-roles-policies", "iam-policy-evaluation", "organizations-scp", "cognito",
		"kms-encryption", "secrets-vs-parameter-store", "sg-vs-nacl", "vpc-endpoints-privatelink",
		"waf-shield", "s3-security", "vpc-fundamentals", "hybrid-cross-vpc-connectivity",
		"security-detection-governance", "multi-account-access-governance", "acm-tls-in-transit-encryption",
	}},
	{"resilient-architectures", []string{
		"regions-az-edge", "elb-types", "auto-scaling-groups", "rds-multi-az-read-replicas",
		"aurora-ha", "sqs-sns-decoupling", "eventbridge", "route53-routing-policies",
		"dr-strategies", "backup-strategies", "observability-cloudwatch-cloudtrail",
		"api-gateway-step-functions", "managed-ai-services-overview",
	}},
	{"high-performing", []string{
		"ec2-families-purchasing", "ebs-efs-instance-store", "s3-performance", "cloudfront",
		"elasticache", "dynamodb-fundamentals", "dynamodb-advanced", "rds-performance",
		"kinesis", "athena-glue", "lambda-performance", "ecs-eks-fargate",
	}},
	{"cost-optimized", []string{
		"pricing-models-ri-sp-spot", "s3-storage-classes-lifecycle", "compute-cost-optimization",
		"data-transfer-costs", "cost-tools-budgets", "cost-optimized-databases-capacity",
		"cost-optimized-databases-storage-lifecycle", "cost-optimized-networking",
		"migration-and-transfer-services", "storage-cost-optimization-beyond-s3",
	}},
}

// TestLoadTopic_AwsContentFile guards content/curriculum/aws.json itself: it
// is data that can rot (a typo'd slug, a dropped concept, a reordered
// chapter) with no compiler to catch it.
func TestLoadTopic_AwsContentFile(t *testing.T) {
	topic, err := LoadTopic(filepath.Join("..", "..", "..", "content", "curriculum", "aws.json"))
	if err != nil {
		t.Fatalf("LoadTopic(aws.json) unexpected error: %v", err)
	}
	if topic.Track().String() != "aws" || topic.Slug().String() != "aws-saa-c03" {
		t.Fatalf("topic = {%q %q}, want {aws aws-saa-c03}", topic.Track().String(), topic.Slug().String())
	}

	chapters := topic.Chapters()
	if len(chapters) != len(wantAwsChapters) {
		t.Fatalf("len(Chapters()) = %d, want %d", len(chapters), len(wantAwsChapters))
	}

	for i, wantCh := range wantAwsChapters {
		ch := chapters[i]
		if ch.Slug().String() != wantCh.slug {
			t.Errorf("Chapters()[%d].Slug() = %q, want %q", i, ch.Slug().String(), wantCh.slug)
		}
		if ch.Position().Int() != i+1 {
			t.Errorf("Chapters()[%d].Position() = %d, want %d", i, ch.Position().Int(), i+1)
		}

		concepts := ch.Concepts()
		if len(concepts) != len(wantCh.concepts) {
			t.Fatalf("Chapters()[%d] %q has %d concepts, want %d", i, wantCh.slug, len(concepts), len(wantCh.concepts))
		}
		for j, wantSlug := range wantCh.concepts {
			c := concepts[j]
			if c.Slug().String() != wantSlug {
				t.Errorf("Chapters()[%d].Concepts()[%d].Slug() = %q, want %q", i, j, c.Slug().String(), wantSlug)
			}
			if c.Position().Int() != j+1 {
				t.Errorf("Chapters()[%d].Concepts()[%d].Position() = %d, want %d", i, j, c.Position().Int(), j+1)
			}
		}
	}
}
