// Package mcqguess measures how guessable mcq recall_checks are without
// reading the lesson: it runs simple test-taking heuristics (longest option,
// fixed position) against the options actually stored in
// content/lessons/*/*.json and compares each heuristic's hit rate to the
// question's own baseline (1/len(options)) — never one hardcoded constant,
// because tracks in this repo use different option counts (3 for the
// existing corpus, 4 for the AWS SAA-C03 bank) and mixing them into one
// number is meaningless. See docs/tickets/aws-cert.md's "AWS-S3" section for
// why this package exists and docs/tickets/mcq-quality.md for the rules it
// measures.
package mcqguess
