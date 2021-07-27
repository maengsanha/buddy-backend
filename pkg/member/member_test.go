package member_test

import (
	"testing"

	"github.com/kmu-kcc/buddy-backend/pkg/member"
)

func TestSignUp(t *testing.T) {
	guests := []*member.Member{
		member.New("20210001", "Test1", "Department1", "1", "010-2021-0001", "testmail1", member.Attending),
		member.New("20190002", "Test2", "Department2", "2", "010-2019-0002", "testmail2", member.Absent),
		member.New("20190003", "Test3", "Department3", "3", "010-2019-0003", "testmail3", member.Attending),
		member.New("20160004", "Test4", "Department2", "4", "010-2016-0004", "testmail4", member.Graduate),
	}

	for _, guest := range guests {
		if err := guest.SignUp(); err != nil {
			t.Error(err)
		}
	}
}

func TestSignUps(t *testing.T) {
	guests, err := member.SignUps()
	if err != nil {
		t.Error(err)
	}

	for _, guest := range guests {
		if guest.Approved {
			t.Error(member.ErrAlreadyMember)
		}
		t.Log(guest)
	}
}

func TestApprove(t *testing.T) {
	ids := []string{"20210001", "20190003"}
	if err := member.Approve(ids); err != nil {
		t.Error(err)
	}
}

func TestSignIn(t *testing.T) {
	memb := member.Member{ID: "20210001", Password: "20210001"}
	guest := member.Member{ID: "20190002", Password: "20190002"}

	if err := memb.SingIn(); err != nil {
		t.Error(err)
	}
	if err := guest.SingIn(); err == member.ErrUnderReview {
		t.Log(err)
	} else if err != nil {
		t.Error(err)
	}
}

func TestExit(t *testing.T) {
	memb := member.Member{ID: "20210001"}
	if err := memb.Exit(); err != nil {
		t.Error(err)
	}
	if err := memb.Exit(); err == member.ErrOnDelete {
		t.Log(err)
	} else if err != nil {
		t.Error(err)
	}

	memb.ID = "20190003"
	if err := memb.Exit(); err != nil {
		t.Error(err)
	}
}

func TestCancelExit(t *testing.T) {
	memb := member.Member{ID: "20210001"}
	if err := memb.CancelExit(); err != nil {
		t.Error(err)
	}
	if err := memb.CancelExit(); err == member.ErrNotOnDelete {
		t.Log(err)
	} else if err != nil {
		t.Error(err)
	}
}

func TestExits(t *testing.T) {
	if membs, err := member.Exits(); err != nil {
		t.Error(err)
	} else {
		for _, memb := range membs {
			t.Log(memb)
		}
	}
}

func TestDelete(t *testing.T) {
	if err := member.Delete([]string{"20190003"}); err != nil {
		t.Error(err)
	}
}

func TestUpdate(t *testing.T) {
	if err := member.Approve([]string{"20190002"}); err != nil {
		t.Error(err)
	}

	memb := member.Member{ID: "20190002"}
	if err := memb.Update(map[string]interface{}{
		"attendance": member.Attending,
		"password":   "00000000"}); err != nil {
		t.Error(err)
	}
}

func TestSearch(t *testing.T) {
	if membs, err := member.Search(map[string]interface{}{"attendance": member.Attending}); err != nil {
		t.Error(err)
	} else {
		for _, memb := range membs {
			t.Log(memb)
		}
	}
}

func TestGraduate(t *testing.T) {
	memb := member.Member{ID: "20210001"}
	if err := memb.Graduate(); err != nil {
		t.Error(err)
	}

	if membs, err := member.Search(map[string]interface{}{"attendance": member.Graduate}); err != nil {
		t.Error(err)
	} else {
		for _, memb := range membs {
			t.Log(memb)
		}
	}
}
