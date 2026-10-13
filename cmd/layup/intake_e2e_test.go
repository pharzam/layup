//go:build e2e

package main

import (
	"fmt"
	"strings"
	"testing"
	"time"
)

// questionsOf gives the rows of questions.tsv on layup-records of the world.
func questionsOf(t *testing.T, w runWorld) [][]string {
	t.Helper()
	out, err := gitRun(w.bare, "show", "refs/heads/layup-records:questions.tsv")
	if err != nil {
		t.Fatalf("questions.tsv: %v\n%s", err, out)
	}
	var rows [][]string
	for i, l := range strings.Split(strings.TrimSuffix(out, "\n"), "\n") {
		if i > 0 {
			rows = append(rows, strings.Split(l, "\t"))
		}
	}
	return rows
}

// answer adds a comment of the idea owner (carol, 303) on the Intake issue, a
// minute after the asked of the question row q, the time read from the
// record (condition 1 of the plan review of #226).
func answer(t *testing.T, w runWorld, id int, q []string, body string) {
	t.Helper()
	answerBy(t, w, id, 303, "carol", q, body)
}

// answerBy is answer for another approver: the Operator (pharzam, 101) answers
// the setup question of the Operator's block (row 53, task T-4c5c).
func answerBy(t *testing.T, w runWorld, id, author int, login string, q []string, body string) {
	t.Helper()
	asked, err := time.Parse(time.RFC3339, q[9])
	if err != nil {
		t.Fatalf("the asked of %q: %v", q, err)
	}
	at := asked.Add(time.Minute).UTC().Format(time.RFC3339)
	w.forge.mu.Lock()
	defer w.forge.mu.Unlock()
	w.forge.people = append(w.forge.people, map[string]any{"id": id, "user": map[string]any{"id": author, "login": login},
		"created_at": at, "updated_at": at, "body": body})
}

// The demo of row 48 (task T-mo9w, #226; intake.md, The acceptance tests of
// M2c, "The steps of Intake"): with the built binary, a local fake forge and
// the world's scripted harness programs, no secret: a restart runs Intake to
// the batch and waits; after a scripted answer comment that leaves one
// question unanswered, the next restart reviews the answers, posts one
// follow-up for that question and waits; after a scripted answer to the
// follow-up, the next restart records it and ends Intake.
func TestTheStepsOfIntake(t *testing.T) {
	binary, _ := buildRunBinary(t)
	w, first := startWorld(t, binary)
	if first.code != 0 {
		t.Fatalf("layup run --new: exit %d\n%s\n%s", first.code, first.stdout, first.stderr)
	}
	restart := func() result { return w.run(t, binary, "run", "acme/target", "--host", ".") }

	// 1. To the batch: G1, the three fields of the intent form, and S01-stack in the Operator's block.
	r := restart()
	if r.code != 0 || !strings.HasSuffix(r.stdout, "\nbatch\tdone\t5 questions, in questions.tsv and the comment of the batch\n"+
		"answers\tdone\twaiting: no answer yet in the block idea-owner and the block operator\n") || len(w.forge.intake) != 1 {
		t.Fatalf("the first restart: exit %d, %d comments of the App\n%s\n%s", r.code, len(w.forge.intake), r.stdout, r.stderr)
	}
	questions := questionsOf(t, w)
	if len(questions) != 5 || questions[0][5] != "G1" || questions[4][1] != "setup" || questions[4][2] != "S01-stack" || questions[4][3] != "operator" {
		t.Fatalf("questions.tsv %q; want G1, the intent form and S01-stack of the Operator", questions)
	}

	// 2. The idea owner answers each question but Q-001; the Operator answers S01-stack.
	answer(t, w, 901, questions[0], fmt.Sprintf("%s: 5.0\n%s: 10.0\n%s: 40.0\n", questions[1][0], questions[2][0], questions[3][0]))
	answerBy(t, w, 903, 101, "pharzam", questions[4], questions[4][0]+": go\n")
	r = restart()
	for _, want := range []string{
		"\nanswers\tdone\t4 answers, in intake/answers.tsv; budget.tsv written\n",
		"\nanswers-review\tdone\t4 answers reviewed, 0 open, in intake/answers-review.tsv\n",
		"\nfollow-up\tdone\t1 questions asked again, in questions.tsv and the comment of the follow-up\n",
	} {
		if !strings.Contains(r.stdout, want) {
			t.Errorf("the second restart has no %q", want)
		}
	}
	if r.code != 0 || !strings.HasSuffix(r.stdout, "\nfollow-up-answers\tdone\twaiting: no answer yet to the follow-up in the block idea-owner\n") {
		t.Fatalf("the second restart: exit %d\n%s\n%s", r.code, r.stdout, r.stderr)
	}
	questions = questionsOf(t, w)
	follow := questions[len(questions)-1]
	if len(questions) != 6 || follow[1] != "follow-up" || follow[2] != questions[0][0] || len(w.forge.intake) != 2 ||
		!strings.HasPrefix(w.forge.intake[1], "LAYUP Intake: the follow-up\n") || !strings.Contains(w.forge.intake[1], "\n"+follow[0]+": ") {
		t.Fatalf("the follow-up: the questions %q, the App's comments %q; want one follow-up of %s", questions, w.forge.intake, questions[0][0])
	}

	// 3. The answer to the follow-up ends Intake.
	answer(t, w, 902, follow, follow[0]+": Go\n")
	r = restart()
	// Intake ends, and the step setup runs after it: in this world, whose baseline is
	// no baseline of a setup and whose harness writes no prose, it is fail (row 54, task T-eahr).
	if r.code != 1 || !strings.Contains(r.stdout, "\nfollow-up-answers\tdone\t1 answers to the follow-up, in intake/answers.tsv; budget.tsv kept; Intake is done\nsetup\tfail\t") {
		t.Fatalf("the third restart: exit %d\n%s\n%s", r.code, r.stdout, r.stderr)
	}
	for path, want := range map[string]string{
		"intake/answers.tsv": "\n" + follow[0] + "\tGo\t902\t303\t",
		"budget.tsv":         "\nappetite.hours\t40.0\t" + questions[3][0] + "\n",
		"questions.tsv":      "\tlayup run\t",
	} {
		if out, err := gitRun(w.bare, "show", "refs/heads/layup-records:"+path); err != nil || !strings.Contains(out, want) {
			t.Errorf("%s of the records branch: %v; want %q\n%s", path, err, want, out)
		}
	}
	if len(w.forge.intake) != 2 {
		t.Errorf("the App's comments on the Intake issue: %d; want the batch and the follow-up", len(w.forge.intake))
	}
}
