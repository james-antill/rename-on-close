package renameonclose_test

import (
	"io/ioutil"
	"os"
	"testing"

	roc "github.com/james-antill/rename-on-close"
)

func TestDiff(t *testing.T) {
	content1 := []byte("one file content 1")
	content2 := []byte("two file content 22")
	content3 := []byte("three file content 333")
	content4 := []byte("four file content 4444")

	if ioutil.WriteFile("t1", content1, 0700) != nil {
		t.Errorf("Fail: Can't write t1")
	}
	if ioutil.WriteFile("t2", content2, 0700) != nil {
		t.Errorf("Fail: Can't write t2")
	}
	if ioutil.WriteFile("t3", content3, 0700) != nil {
		t.Errorf("Fail: Can't write t3")
	}
	if ioutil.WriteFile("t4", content4, 0700) != nil {
		t.Errorf("Fail: Can't write t4")
	}

	t1s, err := os.Stat("t1")
	if err != nil {
		t.Errorf("Fail: Can't stat t1")
	}
	if t1s.Size() != int64(len(content1)) {
		t.Errorf("Fail: T1 size %d != %d", t1s.Size(), len(content1))
	}

	t2s, err := os.Stat("t2")
	if err != nil {
		t.Errorf("Fail: Can't stat t2")
	}
	if t2s.Size() != int64(len(content2)) {
		t.Errorf("Fail: T2 size %d != %d", t2s.Size(), len(content2))
	}

	t3s, err := os.Stat("t3")
	if err != nil {
		t.Errorf("Fail: Can't stat t3")
	}
	if t3s.Size() != int64(len(content3)) {
		t.Errorf("Fail: T3 size %d != %d", t3s.Size(), len(content3))
	}

	t4s, err := os.Stat("t4")
	if err != nil {
		t.Errorf("Fail: Can't stat t4")
	}
	if t4s.Size() != int64(len(content4)) {
		t.Errorf("Fail: T4 size %d != %d", t4s.Size(), len(content4))
	}

	nf, err := roc.Create("t1")
	if err != nil {
		t.Errorf("Fail: Can't create nf t1")
	}
	if _, err := nf.Write(content1); err != nil {
		t.Errorf("Fail: Can't write nf")
	}
	if d, _ := nf.IsDifferent(); d {
		t.Errorf("Fail: nf is different t1")
	}
	if err := nf.CloseRename(); err != nil {
		t.Errorf("Fail: Can't CloseRename nf t1")
	}

	nt1s, err := os.Stat("t1")
	if err != nil {
		t.Errorf("Fail: Can't stat nt1")
	}
	if nt1s.Size() != int64(len(content1)) {
		t.Errorf("Fail: nT1 size %d != %d", nt1s.Size(), len(content1))
	}
	if os.SameFile(t1s, nt1s) {
		t.Errorf("Fail: nT1 is the same %v == %v", t1s, nt1s)
	}

	nf, err = roc.Create("t2")
	if err != nil {
		t.Errorf("Fail: Can't create nf2")
	}
	if _, err := nf.Write(content2); err != nil {
		t.Errorf("Fail: Can't write nf2")
	}
	if d, _ := nf.IsDifferent(); d {
		t.Errorf("Fail: nf is different t2")
	}
	if err := nf.Close(); err != nil {
		t.Errorf("Fail: Can't Close nf t2")
	}
	nf.Close()

	nt2s, err := os.Stat("t2")
	if err != nil {
		t.Errorf("Fail: Can't stat nt2")
	}
	if nt2s.Size() != int64(len(content2)) {
		t.Errorf("Fail: nT2 size %d != %d", nt2s.Size(), len(content2))
	}
	if !os.SameFile(t2s, nt2s) {
		t.Errorf("Fail: nT2 is not the same %v == %v", t2s, nt2s)
	}

	nf, err = roc.Create("t3")
	if err != nil {
		t.Errorf("Fail: Can't create nf t3")
	}
	if _, err := nf.Write(content3); err != nil {
		t.Errorf("Fail: Can't write nf t3")
	}
	nf.Sync()
	if d, _ := nf.IsDifferent(); d {
		t.Errorf("Fail: nf is different t3")
	}
	if err := nf.CloseRename(); err != nil {
		t.Errorf("Fail: Can't CloseRename nf t3")
	}

	nt3s, err := os.Stat("t3")
	if err != nil {
		t.Errorf("Fail: Can't stat nt3")
	}
	if nt3s.Size() != int64(len(content3)) {
		t.Errorf("Fail: nT3 size %d != %d", nt3s.Size(), len(content3))
	}
	if os.SameFile(t3s, nt3s) {
		t.Errorf("Fail: nT3 is the same %v == %v", t3s, nt3s)
	}

	nf, err = roc.Create("t4")
	if err != nil {
		t.Errorf("Fail: Can't create nf4")
	}
	if _, err := nf.Write(content4); err != nil {
		t.Errorf("Fail: Can't write nf4")
	}
	nf.Sync()
	if d, _ := nf.IsDifferent(); d {
		t.Errorf("Fail: nf is different t4")
	}
	if err := nf.Close(); err != nil {
		t.Errorf("Fail: Can't Close nf t4")
	}
	nf.Close()

	nt4s, err := os.Stat("t4")
	if err != nil {
		t.Errorf("Fail: Can't stat nt4")
	}
	if nt4s.Size() != int64(len(content4)) {
		t.Errorf("Fail: nT4 size %d != %d", nt4s.Size(), len(content4))
	}
	if !os.SameFile(t4s, nt4s) {
		t.Errorf("Fail: nT4 is not the same %v == %v", t4s, nt4s)
	}

	os.Remove("t1")
	os.Remove("t2")
	os.Remove("t3")
	os.Remove("t4")
}
